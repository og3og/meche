package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	pb "github.com/og3og/meche/backend/proto/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func main() {
	// Load .env file
	if err := godotenv.Load(filepath.Join("backend", ".env")); err != nil {
		// Try loading from current directory if backend/.env fails
		if err := godotenv.Load(".env"); err != nil {
			log.Printf("Warning: Error loading .env file: %v", err)
		}
	}

	// Get Auth0 token from environment variable
	token := os.Getenv("AUTH0_TOKEN")
	if token == "" && os.Getenv("AUTH_DISABLED") != "true" {
		log.Fatal("Please set AUTH0_TOKEN environment variable or enable AUTH_DISABLED")
	}

	// Set up a connection to the server
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Create a client
	c := pb.NewGreeterServiceClient(conn)

	// Set up a context with timeout and auth token
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if token != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token)
	}

	// Make the gRPC call
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: "World"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
}
