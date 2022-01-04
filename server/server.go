package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	m "mock/proto"

	"google.golang.org/grpc"
)

type Server struct {
	m.UnimplementedMockServer
	Id       int32
	Critical int32
	clients  map[int32]bool
}

func main() {

	//var port = os.Args[1]
	var port = ":8080"
	//husk at der skal laves flere servere med forskellige porte og id'er
	//var port = ":8080"

	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on port %v", listen.Addr())
	}

	grpcServer := grpc.NewServer()
	server := Server{}

	//overvej noget setup/startup server
	server.Id = 1

	//register the server as a grpc server
	m.RegisterMockServer(grpcServer, &server)
	log.Printf("Server listening on %v", listen.Addr())

	//go server.RunServer()

	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("Failed to serve %v", err)
	}
}

//just testing
func (s *Server) RunServer() {
	for {
		time.Sleep(2 * time.Second)
		fmt.Print("Serving brudda")
	}
}

func (s *Server) Increment(ctx context.Context, in *m.Request) (*m.Value, error) {
	log.Printf("Increment has been called on server %v by client %v", s.Id, in.Id)
	//prepare the output
	var output = m.Value{Value: s.Critical}

	//increment the value
	s.Critical++

	//what to do with error?

	return &output, nil
}
