package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/mischat/zkp-auth/pb"
)

type server struct{}

func (s *server) Authenticate(ctx context.Context, in *pb.AuthRequest) (*pb.AuthResponse, error) {
	log.Printf("Received: %v", in.GetUsername())
	return &pb.AuthResponse{Authenticated: true}, nil
}

func main() {
	fmt.Println("Hello, world!")
	// Output: Hello, world!
}
