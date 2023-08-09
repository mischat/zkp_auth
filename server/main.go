package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/big"
	"net"

	pb "github.com/mischat/zkp_auth/pb"
	zkpautils "github.com/mischat/zkp_auth/utils"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")

	// Public variables needed for the auth system to work
	p = flag.Int64("p", 23, "the prime number we group from")
	q = flag.Int64("q", 11, "for prime order calculation")
	g = flag.Int64("g", 4, "first in group")
	h = flag.Int64("h", 9, "second in group")
)

// server is used to implement zkp_auth.server
type server struct {
	pb.UnimplementedAuthServer
	userRegData map[string]UserRegistration
	// This datastructure doesn't have any timeout in place
	// One for the future, would be good for these not to be valid for long
	authenticationData map[string]Authentication
}

func newServer() *server {
	return &server{
		userRegData:        make(map[string]UserRegistration),
		authenticationData: make(map[string]Authentication),
	}
}

// This stores the user registration data agains the user ID
type UserRegistration struct {
	y1 *big.Int
	y2 *big.Int
}

type Authentication struct {
	user    string
	r1      *big.Int
	r2      *big.Int
	c       *big.Int
	session string
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

func (s *server) CreateAuthenticationChallenge(ctx context.Context, in *pb.AuthenticationChallengeRequest) (*pb.AuthenticationChallengeResponse, error) {
	log.Printf("Received UserID: %v", in.GetUser())
	log.Printf("Received R1: %v", in.GetR1())
	log.Printf("Received R2: %v", in.GetR2())

	// Retrieve User from the map
	_, exists := s.userRegData[in.GetUser()]
	if !exists {
		return &pb.AuthenticationChallengeResponse{}, fmt.Errorf("user doesn't exists")
	}

	// Now the challenger picks a random value c
	// it is important that these are unique for each user
	// ideally we store the used ones somewhere, like in an associative array or something.
	// but for this excercise we will just generate a new random one each time
	// using a big(ish) number to ensure some randomness
	c := zkpautils.RandomBigInt()
	log.Printf("Generated random c: %d", c)

	// Store c in the authenticationmap
	// Note that this will only allow for a user to authenticate in one place at a time
	authId := zkpautils.RandomString(20)

	s.authenticationData[authId] = Authentication{
		user:    in.GetUser(),
		r1:      new(big.Int).SetInt64(in.GetR1()),
		r2:      new(big.Int).SetInt64(in.GetR2()),
		c:       c,
		session: "",
	}

	return &pb.AuthenticationChallengeResponse{AuthId: authId, C: zkpautils.BigIntToInt64(c)}, nil

}

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
	// The config is now validated and in good shape

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
