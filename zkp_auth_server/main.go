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

// server is used to implement zkp_auth.zkp_auth_server
type server struct {
	pb.UnimplementedAuthServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
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

//type server struct{}
//
//func (s *server) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
//	log.Printf("Received\n")
//	return &pb.RegisterResponse{}, nil
//}

//func main() {
//	fmt.Println("Hello, world!")
//
//	// Start the gRPC server
//	lis, err := net.Listen("tcp", ":50051")
//	if err != nil {
//		log.Fatalf("failed to listen: %v", err)
//	}
//	s := grpc.NewServer()
//
//	pb.RegisterAuthServer(s, &server{})
//	if err := s.Serve(lis); err != nil {
//		log.Fatalf("failed to serve: %v", err)
//	}
//
//	// Create a gRPC client
//	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
//	if err != nil {
//		log.Fatalf("failed to dial: %v", err)
//	}
//	defer conn.Close()
//
//	client := pb.NewAuthClient(conn)
//
//	// Call an RPC method
//	resp, err := client.Register(context.Background(), &pb.RegisterRequest{Y1: 13, Y2: 1})
//	if err != nil {
//		log.Fatalf("failed to call Authenticate: %v", err)
//	}
//	log.Printf("response: %v", resp)
//
//}
//
