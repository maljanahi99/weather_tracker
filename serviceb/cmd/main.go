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

	grpcServer := grpc.NewServer()
	pb.RegisterWeatherServiceServer(grpcServer, &sv.Server{})

	log.Println("Service B (gRPC) running on :50051")
	log.Fatal(grpcServer.Serve(lis))
}
