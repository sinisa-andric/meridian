package main

import (
	"context"

	namespace "github.com/sinisa-andric/meridian/pkg/model"
)

func (s *server) SayHi(ctx context.Context, in *namespace.NameSpaceRequest) (*namespace.NameSpaceResponse, error) {
	return &namespace.NameSpaceResponse{Data: "Hi!"}, nil
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

func (s *server) Describe(ctx context.Context, in *namespace.NameSpaceRequest) (*namespace.NameSpaceResponse, error) {
	return &namespace.NameSpaceResponse{Data: "Describe!"}, nil
}
