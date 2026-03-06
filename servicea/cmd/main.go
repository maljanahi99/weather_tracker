package main

import (
	"log"
	"os"

	"github.com/maljanahi99/weather_tracker/servicea"
)

func main() {
	serviceBAddr := os.Getenv("SERVICEB_ADDR")
	if serviceBAddr == "" {
		serviceBAddr = "serviceb:50051"
	}

	client, conn, err := servicea.NewServiceBClient(serviceBAddr)
	if err != nil {
		log.Fatal("failed to connect to service B: ", err)
	}
	defer conn.Close()

	log.Println("Service A (HTTP) running on :8080")
	log.Fatal(servicea.StartServer(client))
}
