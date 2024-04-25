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
	return &namespace.NameSpaceResponse{Data: "Hi!"}, nil
}

func (s *server) Check(ctx context.Context, in *namespace.NameSpaceRequest) (*namespace.NameSpaceResponse, error) {
	return &namespace.NameSpaceResponse{Data: "Check!"}, nil
}

func (s *server) Create(ctx context.Context, in *namespace.NameSpaceRequest) (*namespace.NameSpaceResponse, error) {
	return &namespace.NameSpaceResponse{Data: "Create!"}, nil
}

func (s *server) Read(ctx context.Context, in *namespace.NameSpaceRequest) (*namespace.NameSpaceResponse, error) {
	return &namespace.NameSpaceResponse{Data: "Read!"}, nil
}

func (s *server) Update(ctx context.Context, in *namespace.NameSpaceRequest) (*namespace.NameSpaceResponse, error) {
	return &namespace.NameSpaceResponse{Data: "Update!"}, nil
}

func (s *server) Delete(ctx context.Context, in *namespace.NameSpaceRequest) (*namespace.NameSpaceResponse, error) {
	return &namespace.NameSpaceResponse{Data: "Delete!"}, nil
}

func (s *server) SoftDelete(ctx context.Context, in *namespace.NameSpaceRequest) (*namespace.NameSpaceResponse, error) {
	return &namespace.NameSpaceResponse{Data: "Soft Delete!"}, nil
}

func (s *server) Describe(ctx context.Context, in *namespace.NameSpaceRequest) (*namespace.NameSpaceResponse, error) {
	return &namespace.NameSpaceResponse{Data: "Describe!"}, nil
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
