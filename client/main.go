// Package main implements a client for Greeter service.
package main

import (
	"context"
	"flag"
	"log"
	"math/big"
	"time"

	pb "github.com/mischat/zkp_auth/pb"
	zkpautils "github.com/mischat/zkp_auth/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")

	// Public variables needed for the auth system to work
	p = flag.Int64("p", 23, "the prime number we group from")
	q = flag.Int64("q", 11, "for prime order calculation")
	g = flag.Int64("g", 4, "first in group")
	h = flag.Int64("h", 9, "second in group")
)

func main() {
	flag.Parse()

	// creating bigInts from the flags
	bP := big.NewInt(*p)
	bQ := big.NewInt(*q)
	bG := big.NewInt(*g)
	bH := big.NewInt(*h)

	log.Printf("p: %v q: %v g: %v h: %v\n", bP, bQ, bG, bH)

	// This makes sure that we validate the public variables passed in
	_, err := zkpautils.ValidatePublicVariables(bP, bQ, bG, bH)
	if err != nil {
		log.Fatalf("could not validate public variables: %v", err)
	}

	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewAuthClient(conn)

	// Send a Register request to the server
	user := "alice1@example.com"
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
