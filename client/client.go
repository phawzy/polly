package main

import (
	"context"
	"flag"
	"log"

	pb "github.com/phawzy/polly/polly"
	"google.golang.org/grpc"
)

var (
	action     = flag.String("action", "", "add_poll or upvote_poll or downvote_poll or list_polls")
	pollName   = flag.String("poll_name", "", "string poll_name")
	serverAddr = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
)

func main() {
	flag.Parse()
	println(*serverAddr)
	conn, err := grpc.Dial(*serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewPollyClient(conn)
	ctx := context.Background()
	if *action == "add_poll" {
		response, _ := client.AddPoll(ctx, &pb.Poll{PollName: *pollName})
		if response.Done {
			println("done")
		}
	} else if *action == "upvote_poll" {
		response, _ := client.UpdatePoll(ctx, &pb.UpdatePoll{PollName: *pollName, PollAction: "up"})
		if response.Done {
			println("done")
		}
	} else if *action == "downvote_poll" {
		response, _ := client.UpdatePoll(ctx, &pb.UpdatePoll{PollName: *pollName, PollAction: "down"})
		if response.Done {
			println("done")
		}
	} else if *action == "list_polls" {
		response, _ := client.ListPolls(ctx, &pb.Empty{})
		if len(response.Polls) == 0 {
			println("no polls")
		} else {
			println(response.Polls)
		}

	} else {
		response, _ := client.AddPoll(ctx, &pb.Poll{PollName: *pollName})
		if response.Done {
			println("done")
		}
	}

}
