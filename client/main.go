// Package main implements a client for Greeter service.
package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "github.com/mischat/zkp_auth/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewAuthClient(conn)

	// Send a Register request to the server
	user := "alice@example.com"
	y1 := int64(2)
	y2 := int64(3)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = c.Register(ctx, &pb.RegisterRequest{User: user, Y1: y1, Y2: y2})
	if err != nil {
		log.Fatalf("could not register: %v", err)
	}
	log.Printf("Registered user %s with Y1=%d and Y2=%d", user, y1, y2)
}
