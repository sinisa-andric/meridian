package main

import (
	"context"
	"log"
	"net"

	namespace "github.com/sinisa-andric/meridian/pkg/model"
	"google.golang.org/grpc"
)

type server struct {
	namespace.UnimplementedNameSpaceServiceServer
}

func (s *server) SayHi(ctx context.Context, in *namespace.NameSpaceRequest) (*namespace.NameSpaceResponse, error) {
	return &namespace.NameSpaceResponse{Data: "Hi! "}, nil
}

func main() {

	listen, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatalf("failed to listen on port 9001: %v", err)
	}

	s := grpc.NewServer()

	namespace.RegisterNameSpaceServiceServer(s, &server{})

	log.Printf("gRPC server listening at %v", listen.Addr())

	if err := s.Serve(listen); err != nil {
		log.Fatalf("failed to register server %v", err)
	}

}
