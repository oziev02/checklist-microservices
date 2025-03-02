package main

import (
	"github.com/joho/godotenv"
	"github.com/oziev02/checklist-microservices/internal/db/infrastructure/grpc_server"
	"log"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	server, err := grpc_server.NewServer()
	if err != nil {
		log.Fatalf("Failed to create gRPC server: %v", err)
	}
	defer server.Stop()

	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}
