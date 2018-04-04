# Description
 simple Golang cli client and server using grpc to add polls and votes.

Please refer to [gRPC Basics: Go](https://grpc.io/docs/tutorials/basic/go.html) for more information about installing grpc for golang.

# Run the code
To compile and run the server, assuming you are in the root of the polly
folder, i.e., .../github/phawzy/polly/, simply:

```sh
$ go run server/server.go
```

Likewise, to run the client:

```sh
$ go run client/client.go [options]
```

Options are [-add_pull pollname] or [-upvote_poll pollname] or [-downvote_poll pollname] or [-list_polls]

problems and enhancements
=========================
* error handling for db queries responses
* using streams or repeated instead of joined string in listing polls
* seperating db queries from other logic or maybe use orm
* detailed documentation
* using mutex with shared resources like db driver singlton



# Optional command line flags
The server and client both take optional command line flags. For example, the
address/port that server uses
