package main

import (
	"log"
	"net"

	pb "github.com/maljanahi99/weather_tracker/protos"
	"github.com/maljanahi99/weather_tracker/servicea"
	sv "github.com/maljanahi99/weather_tracker/serviceb/grpc"
	"google.golang.org/grpc"
)

func main() {
	// 1) Start Service B (gRPC) in a goroutine
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		grpcServer := grpc.NewServer()
		pb.RegisterWeatherServiceServer(grpcServer, &sv.Server{})

		log.Println("Service B (gRPC) running on :50051")
		log.Fatal(grpcServer.Serve(lis))
	}()

	// 2) Create Service B client for Service A
	client, conn, err := servicea.NewServiceBClient("localhost:50051")
	if err != nil {
		log.Fatal("failed to connect to service B:", err)
	}
	defer conn.Close()

	// 3) Start Service A (HTTP) using that client
	log.Println("Service A (HTTP) running on :8080")
	log.Fatal(servicea.StartServer(client))
}
