package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/oziev02/checklist-microservices/internal/api/infrastructure/grpc_client"
	"github.com/oziev02/checklist-microservices/internal/api/infrastructure/http/handlers"
	"log"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Создаём gRPC-клиент
	grpcClient, err := grpc_client.NewClient()
	if err != nil {
		log.Fatalf("Failed to create gRPC client: %v", err)
	}
	defer grpcClient.Close()

	r := gin.Default()
	handlers.RegisterRoutes(r, grpcClient)

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
