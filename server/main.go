package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/big"
	"net"

	pb "github.com/mischat/zkp_auth/pb"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement zkp_auth.server
type server struct {
	pb.UnimplementedAuthServer
	userRegData map[string]UserRegistration
}

func newServer() *server {
	return &server{
		userRegData: make(map[string]UserRegistration),
	}
}

// This stores the user registration data agains the user ID
type UserRegistration struct {
	y1 *big.Int
	y2 *big.Int
}

// This implements the Register gRPC call
// Note that this implementation does not support updating a user's registration info
func (s *server) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	log.Printf("Received UserID: %v", in.GetUser())
	log.Printf("Received Y1: %v", in.GetY1())
	log.Printf("Received Y2: %v", in.GetY2())

	// Retrieve User from the map
	_, exists := s.userRegData[in.GetUser()]
	if exists {
		return &pb.RegisterResponse{}, fmt.Errorf("user already exists")
	}
	// Store Y1 and Y2 in the userRegData map
	s.userRegData[in.GetUser()] = UserRegistration{
		y1: new(big.Int).SetInt64(in.GetY1()),
		y2: new(big.Int).SetInt64(in.GetY2()),
	}

	log.Printf("Stored UserID: %v", in.GetUser())

	return &pb.RegisterResponse{}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterAuthServer(s, newServer())
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
