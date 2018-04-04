package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	"golang.org/x/net/context"

	pb "github.com/phawzy/polly/polly"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	port = flag.Int("port", 10000, "The server port")
)

type DB struct {
	dbDriver *sql.DB
}

func CreateDBInstance() *sql.DB {
	mydb, err := sql.Open("sqlite3", "./polls.db")
	checkErr(err)
	return mydb
}

var dbInstance *DB

func GetDBInstance() *DB {
	if dbInstance == nil {
		dbInstance = &DB{dbDriver: CreateDBInstance()} // <--- NOT THREAD SAFE
	}
	return dbInstance
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

//PollyServer struct for server
type pollyServer struct{}

//AddPoll add poll
func (grpcServer *pollyServer) AddPoll(ctx context.Context, in *pb.Poll) (*pb.Response, error) {
	var db = GetDBInstance().dbDriver
	println(in.PollName)

	// insert
	addPoll, err := db.Prepare("INSERT INTO polls(pollName, upvotes, downvotes) values(?,?,?)")
	if err != nil {
		checkErr(errors.New("fuck"))
	}

	_, err = addPoll.Exec(in.PollName, 0, 0)

	if err != nil {
		return &pb.Response{Done: false}, errors.New("poll already exists")
	}
	return &pb.Response{Done: true}, nil

}

func (grpcServer *pollyServer) UpdatePoll(ctx context.Context, in *pb.UpdatePoll) (*pb.Response, error) {
	var db = GetDBInstance().dbDriver

	rows, err := db.Query("SELECT pollName, upvotes, downvotes FROM polls where pollName='" + in.PollName + "'")
	checkErr(err)
	var pollName string
	var upvotes int
	var downvotes int
	if rows.Next() {
		err = rows.Scan(&pollName, &upvotes, &downvotes)
		checkErr(err)
	} else {
		return &pb.Response{Done: false}, errors.New("poll doesn't exists")
	}
	rows.Close()

	if in.PollAction == "up" {
		upvotes = upvotes + 1
		upvotePoll, err := db.Prepare("update polls set upvotes=? where pollName=?")
		checkErr(err)

		_, err = upvotePoll.Exec(upvotes, in.PollName)
		checkErr(err)
	}
	if in.PollAction == "down" {
		downvotes = downvotes + 1
		downvotePoll, err := db.Prepare("update polls set downvotes=? where pollName=?")
		checkErr(err)

		_, err = downvotePoll.Exec(downvotes, in.PollName)
		checkErr(err)
	}
	println("updated : " + in.PollAction + " on poll " + in.PollName)
	return &pb.Response{Done: true}, nil
}

func (grpcServer *pollyServer) ListPolls(ctx context.Context, in *pb.Empty) (*pb.Polls, error) {
	pollsList := []string{}
	var db = GetDBInstance().dbDriver
	rows, err := db.Query("SELECT pollName, upvotes, downvotes FROM polls")
	checkErr(err)
	var pollName string
	var upvotes int
	var downvotes int

	for rows.Next() {
		err = rows.Scan(&pollName, &upvotes, &downvotes)
		checkErr(err)
		println(pollName)
		pollsList = append(pollsList, pollName+"{up: "+strconv.Itoa(upvotes)+", down: "+strconv.Itoa(downvotes)+"}")
	}

	rows.Close() //good habit to close
	return &pb.Polls{Polls: strings.Join(pollsList[:], ",")}, nil
}

func main() {

	var db = GetDBInstance().dbDriver
	defer db.Close()
	sqlStmt := `
	create table IF NOT EXISTS polls (pollName text not null primary key, upvotes integer, downvotes integer);
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	if err != nil {
		log.Fatal(err)
	}

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterPollyServer(grpcServer, &pollyServer{})
	println("started")
	reflection.Register(grpcServer)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
