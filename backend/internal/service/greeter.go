package service

import (
	"context"
	"fmt"

	pb "github.com/og3og/meche/backend/proto/v1"
)

type GreeterServer struct {
	pb.UnimplementedGreeterServiceServer
}

func NewGreeterServer() *GreeterServer {
	return &GreeterServer{}
}

func (s *GreeterServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	message := fmt.Sprintf("Hello, %s!", req.Name)
	return &pb.HelloResponse{
		Message: message,
	}, nil
}
