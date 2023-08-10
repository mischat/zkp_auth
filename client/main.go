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
	// ultimately these need to be strings, as bigInts are bigger than int64s
	p = flag.Int64("p", 23, "the prime number we group from")
	q = flag.Int64("q", 11, "for prime order calculation")
	g = flag.Int64("g", 4, "first in group")
	h = flag.Int64("h", 9, "second in group")

	// This is the client id and secret
	u = flag.String("u", "alice@example.com", "the client id")
	x = flag.Int64("x", 6, "the client secret")
)

func main() {
	flag.Parse()

	// creating bigInts from the flags
	bP := big.NewInt(*p)
	bQ := big.NewInt(*q)
	bG := big.NewInt(*g)
	bH := big.NewInt(*h)

	bX := big.NewInt(*x)

	log.Printf("p: %v q: %v g: %v h: %v\n", bP, bQ, bG, bH)

	// This makes sure that we validate the public variables passed in
	_, err := zkpautils.ValidatePublicVariables(bP, bQ, bG, bH)
	if err != nil {
		log.Fatalf("could not validate public variables: %v", err)
	}
	// The config is now validated and in good shape

	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewAuthClient(conn)

	// Now to calculate y1
	y1 := zkpautils.CalculateExp(bG, bX, bP)
	// Now to calculate y2
	y2 := zkpautils.CalculateExp(bH, bX, bP)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// Note that i am doing a conversion from big.Int to int64 here
	// Ideally we would change the proto to use strings instead of int64s but i didn't want to do this for this excercise
	// bigInt has a marshaltext and unmarshaltext method that we could use to do this
	_, err = c.Register(ctx, &pb.RegisterRequest{User: *u, Y1: zkpautils.BigIntToInt64(y1), Y2: zkpautils.BigIntToInt64(y2)})
	if err != nil {
		log.Fatalf("could not register: %v", err)
	}
	log.Printf("Registered user %s with Y1=%d and Y2=%d", *u, y1, y2)

	// now we need to generate a random k
	// it is important that these are unique for each request
	// ideally we store the used ones somewhere, but for this excercise we will just generate a new random one each time
	// using a big(ish) number to ensure some randomness
	k := zkpautils.RandomBigInt()
	log.Printf("Generated random k: %d", k)

	// Now to calculate (r1, r2) = g^k, h^k
	r1 := zkpautils.CalculateExp(bG, k, bP)
	r2 := zkpautils.CalculateExp(bH, k, bP)

	resp, err := c.CreateAuthenticationChallenge(ctx, &pb.AuthenticationChallengeRequest{User: *u, R1: zkpautils.BigIntToInt64(r1), R2: zkpautils.BigIntToInt64(r2)})
	if err != nil {
		log.Fatalf("failed to create auth challenge: %v", err)
	}

	authId := resp.AuthId
	chal := big.NewInt(resp.C)
	log.Printf("authId: %s c: %d", authId, chal)

	// Not to calculate s = (k - c .x) mod q
	s := zkpautils.CalculateS(k, chal, bX, bQ)

	verResp, err := c.VerifyAuthentication(ctx, &pb.AuthenticationAnswerRequest{AuthId: authId, S: zkpautils.BigIntToInt64(s)})
	if err != nil {
		log.Fatalf("Failed to auth: %v", err)
	}

	log.Printf("Success, this is our session ID: '%s'", verResp.SessionId)
}
