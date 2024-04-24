package main

import (
	"context"
	"log"
	"time"

	namespace "github.com/sinisa-andric/meridian/pkg/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	connection, err := grpc.Dial("localhost:9001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to gRPC server at localhost:9001: %v", err)
	}
	defer connection.Close()

	c := namespace.NewNameSpaceServiceClient(connection)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := c.SayHi(ctx, &namespace.NameSpaceRequest{})
	if err != nil {
		log.Fatalf("error calling function SayHello: %v", err)
	}

	log.Printf("Response from gRPC server's SayHi function: %s", response.GetData())

}
