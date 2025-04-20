package main

import (
	"log"
	"net"

	grpcserver "example.com/trading-app/grpc_server"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}
	go grpcserver.RunSimulation()
	s := grpcserver.New()
	log.Printf("gRPC server listening on %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
