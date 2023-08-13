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
	addrFlag = flag.String("addr", "localhost:50051", "the address to connect to")

	// Public variables needed for the auth system to work
	// String so that we can handle big numbers
	pFlag = flag.String("p", "23", "the prime number we group from")
	qFlag = flag.String("q", "11", "for prime order calculation")
	gFlag = flag.String("g", "12", "first in group")
	hFlag = flag.String("h", "13", "second in group")

	// This is the client id and secret
	uFlag = flag.String("u", "alice@example.com", "the client id")
	xFlag = flag.String("x", "6", "the client secret")
)

func main() {
	flag.Parse()

	// creating bigInts from the flags
	p := new(big.Int)
	p.SetString(*pFlag, 10)

	q := new(big.Int)
	q.SetString(*qFlag, 10)

	g := new(big.Int)
	g.SetString(*gFlag, 10)

	h := new(big.Int)
	h.SetString(*hFlag, 10)

	x := new(big.Int)
	x.SetString(*xFlag, 10)

	log.Printf("p: %v q: %v g: %v h: %v\n", p, q, g, h)

	// This makes sure that we validate the public variables passed in
	_, err := zkpautils.ValidatePublicVariables(p, q, g, h)
	if err != nil {
		log.Fatalf("could not validate public variables: %v", err)
	}
	// The config is now validated and in good shape

	// Set up a connection to the server.
	conn, err := grpc.Dial(*addrFlag, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewAuthClient(conn)

	// Now to calculate y1
	y1 := zkpautils.CalculateExp(g, x, p)
	// Now to calculate y2
	y2 := zkpautils.CalculateExp(h, x, p)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// Note that i am doing a conversion from big.Int to int64 here
	// Ideally we would change the proto to use strings instead of int64s but i didn't want to do this for this excercise
	// bigInt has a marshaltext and unmarshaltext method that we could use to do this
	_, err = c.Register(ctx, &pb.RegisterRequest{User: *uFlag, Y1: y1.String(), Y2: y2.String()})
	if err != nil {
		log.Fatalf("could not register: %v", err)
	}
	log.Printf("Registered user %s with Y1=%d and Y2=%d", *uFlag, y1, y2)

	// now we need to generate a random k
	// it is important that these are unique for each request
	// ideally we store the used ones somewhere, but for this excercise we will just generate a new random one each time
	// using a big(ish) number to ensure some randomness
	k := zkpautils.RandomBigInt()
	log.Printf("Generated random k: %d", k)

	// Now to calculate (r1, r2) = g^k, h^k
	r1 := zkpautils.CalculateExp(g, k, p)
	r2 := zkpautils.CalculateExp(h, k, p)

	resp, err := c.CreateAuthenticationChallenge(ctx, &pb.AuthenticationChallengeRequest{User: *uFlag, R1: r1.String(), R2: r2.String()})
	if err != nil {
		log.Fatalf("failed to create auth challenge: %v", err)
	}

	authId := resp.AuthId
	chal := new(big.Int)
	chal.SetString(resp.C, 10)
	log.Printf("authId: %s c: %d", authId, chal)

	// Not to calculate s = (k - c .x) mod q
	s := zkpautils.CalculateS(k, chal, x, q)

	verResp, err := c.VerifyAuthentication(ctx, &pb.AuthenticationAnswerRequest{AuthId: authId, S: s.String()})
	if err != nil {
		log.Fatalf("Failed to auth: %v", err)
	}

	log.Printf("Success, this is our session ID: '%s'", verResp.SessionId)
}
