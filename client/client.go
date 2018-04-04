/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a simple gRPC client that demonstrates how to use gRPC-Go libraries
// to perform unary, client streaming, server streaming and full duplex RPCs.
//
// It interacts with the route guide service whose definition can be found in routeguide/route_guide.proto.
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
