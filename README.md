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
$ go run client/client.go
```

# Optional command line flags
The server and client both take optional command line flags. For example, the
address/port that server uses
