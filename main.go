package main

import (
	"log"
	"net"

	pb "github.com/maljanahi99/weather_tracker/protos"
	sv "github.com/maljanahi99/weather_tracker/serviceb/grpc"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcserver := grpc.NewServer()

	service := &sv.Server{}

	pb.RegisterWeatherServiceServer(grpcserver, service)

	log.Println("Service B running on :50051")

	if err := grpcserver.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
