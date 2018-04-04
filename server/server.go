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

type Poll struct {
	up   int
	down int
}

var db = sql.DB{}
var polls = map[string]Poll{}

//PollyServer struct for server
type pollyServer struct{}

//AddPoll add poll
func (grpcServer *pollyServer) AddPoll(ctx context.Context, in *pb.Poll) (*pb.Response, error) {
	_, noError := polls[in.PollName]
	if !noError {
		polls[in.PollName] = Poll{up: 0, down: 0}
		return &pb.Response{Done: true}, nil
	}
	return &pb.Response{Done: false}, errors.New("poll already exists")

}

func (grpcServer *pollyServer) UpdatePoll(ctx context.Context, in *pb.UpdatePoll) (*pb.Response, error) {
	currentPoll, noError := polls[in.PollName]
	if !noError {
		return &pb.Response{Done: false}, errors.New("poll doesn't exists")
	}
	up := currentPoll.up
	if in.PollAction == "up" {
		up = up + 1
	}
	down := currentPoll.down
	if in.PollAction == "down" {
		down = down + 1
	}

	polls[in.PollName] = Poll{up: up, down: down}
	println("updated : " + in.PollAction + " on poll " + in.PollName)
	return &pb.Response{Done: true}, nil
}

func (grpcServer *pollyServer) ListPolls(ctx context.Context, in *pb.Empty) (*pb.Polls, error) {
	pollsList := []string{}
	for pollname, pollVotes := range polls {
		println(pollname)
		pollsList = append(pollsList, pollname+"{up: "+strconv.Itoa(pollVotes.up)+", down: "+strconv.Itoa(pollVotes.down)+"}")
	}
	return &pb.Polls{Polls: strings.Join(pollsList[:], ",")}, nil
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	db, err := sql.Open("sqlite3", "./polls.db")
	checkErr(err)

	sqlStmt := `
	create table IF NOT EXISTS polls (pollName text not null primary key, upvotes integer, downvotes integer);
	delete from polls;
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	if err != nil {
		log.Fatal(err)
	}

	// insert
	addPoll, err := db.Prepare("INSERT INTO polls(pollName, upvotes, downvotes) values(?,?,?)")
	checkErr(err)
	res, err := addPoll.Exec("hell", 0, 0)
	checkErr(err)
	id, err := res.LastInsertId()
	checkErr(err)

	fmt.Println(id)

	upvotePoll, err := db.Prepare("update polls set upvotes=? where pollName=?")
	checkErr(err)

	res, err = upvotePoll.Exec(1, "hell")
	checkErr(err)

	downvotePoll, err := db.Prepare("update polls set downvotes=? where pollName=?")
	checkErr(err)

	res, err = downvotePoll.Exec(1, "hell")
	checkErr(err)

	rows, err := db.Query("SELECT pollName, upvotes, downvotes FROM polls")
	checkErr(err)
	var pollName string
	var upvotes int
	var downvotes int

	for rows.Next() {
		err = rows.Scan(&pollName, &upvotes, &downvotes)
		checkErr(err)
		fmt.Println(pollName)
		fmt.Println(upvotes)
		fmt.Println(downvotes)
	}

	rows.Close() //good habit to close

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
