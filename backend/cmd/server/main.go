package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/og3og/meche/backend/internal/middleware"
	"github.com/og3og/meche/backend/internal/service"
	pb "github.com/og3og/meche/backend/proto/v1"
)

func main() {
	// Load .env file
	if err := godotenv.Load(filepath.Join("backend", ".env")); err != nil {
		// Try loading from current directory if backend/.env fails
		if err := godotenv.Load(".env"); err != nil {
			log.Printf("Warning: Error loading .env file: %v", err)
		}
	}

	// Load Auth0 configuration
	if os.Getenv("AUTH0_DOMAIN") == "" || os.Getenv("AUTH0_AUDIENCE") == "" {
		log.Fatal("Please set AUTH0_DOMAIN and AUTH0_AUDIENCE environment variables")
	}

	// Create a listener for gRPC
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create a gRPC server with auth interceptor
	s := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.UnaryAuthInterceptor),
	)
	pb.RegisterGreeterServiceServer(s, service.NewGreeterServer())

	// Start gRPC server in a goroutine
	go func() {
		log.Printf("gRPC server listening at %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	// Create a client connection to the gRPC server
	conn, err := grpc.DialContext(
		context.Background(),
		"localhost:50051",
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()

	// Create a new ServeMux for the HTTP server
	gwmux := runtime.NewServeMux()

	// Register Greeter handler
	err = pb.RegisterGreeterServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalf("Failed to register gateway: %v", err)
	}

	// Create HTTP handler with auth middleware
	handler := middleware.EnsureValidToken()(gwmux)

	// Start HTTP server
	gwServer := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	log.Printf("HTTP server listening at :8080")
	if err := gwServer.ListenAndServe(); err != nil {
		log.Fatalf("failed to serve HTTP: %v", err)
	}
}
