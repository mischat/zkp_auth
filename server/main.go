package main

import (
	"context"
	"flag"
	"fmt"
	"log"
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
}

// SayHello implements helloworld.GreeterServer
func (s *server) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	log.Printf("Received Y1: %v", in.GetY1())
	log.Printf("Received Y2: %v", in.GetY2())
	return &pb.RegisterResponse{}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterAuthServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
