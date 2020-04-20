package service

import (
	aPb "github.com/c12s/scheme/apollo"
	"google.golang.org/grpc"
	"log"
)

func NewApolloClient(address string) aPb.ApolloServiceClient {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to start gRPC connection to apollo service: %v", err)
	}

	return aPb.NewApolloServiceClient(conn)
}
