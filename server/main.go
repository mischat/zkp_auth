package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/big"
	"net"
	"time"

	pb "github.com/mischat/zkp_auth/pb"
	zkpautils "github.com/mischat/zkp_auth/utils"
	"google.golang.org/grpc"
)

var (
	portFlag = flag.Int("port", 50051, "The server port")

	// Public variables needed for the auth system to work
	pFlag = flag.Int64("p", 23, "the prime number we group from")
	qFlag = flag.Int64("q", 11, "for prime order calculation")
	gFlag = flag.Int64("g", 12, "first in group")
	hFlag = flag.Int64("h", 13, "second in group")
)

// server is used to implement zkp_auth.server
type server struct {
	pb.UnimplementedAuthServer
	userRegData map[string]UserRegistration
	// This datastructure doesn't have any timeout in place
	// One for the future, would be good for these not to be valid for long
	authenticationData map[string]Authentication

	sessionData map[string]Session
}

func newServer() *server {
	return &server{
		userRegData:        make(map[string]UserRegistration),
		authenticationData: make(map[string]Authentication),
		sessionData:        make(map[string]Session),
	}
}

// This stores the user registration data against the user ID
type UserRegistration struct {
	y1 *big.Int
	y2 *big.Int
}

// This stores the authentication data against the auth ID
type Authentication struct {
	user string
	r1   *big.Int
	r2   *big.Int
	c    *big.Int
}

// This stores the session data against the session ID
type Session struct {
	user      string
	createdAt time.Time
}

// This implements the Register gRPC call
// Note that this implementation does not support updating a user's registration info
func (srv *server) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	log.Printf("Received UserID: %v", in.GetUser())
	log.Printf("Received Y1: %v", in.GetY1())
	log.Printf("Received Y2: %v", in.GetY2())

	y1 := new(big.Int)
	y1.SetString(in.GetY1(), 10)

	y2 := new(big.Int)
	y2.SetString(in.GetY2(), 10)

	// Retrieve User from the map
	_, exists := srv.userRegData[in.GetUser()]
	if exists {
		return &pb.RegisterResponse{}, fmt.Errorf("user '%v' already exists", in.GetUser())
	}

	// TODO: validate that y1 and y2 are in the group

	// Store Y1 and Y2 in the userRegData map
	srv.userRegData[in.GetUser()] = UserRegistration{
		y1: y1,
		y2: y2,
	}

	log.Printf("Stored UserID: %v", in.GetUser())

	return &pb.RegisterResponse{}, nil
}

// This is the second step in the authentication process
// The client will send the r1 and r2 values based on a random value k
func (srv *server) CreateAuthenticationChallenge(ctx context.Context, in *pb.AuthenticationChallengeRequest) (*pb.AuthenticationChallengeResponse, error) {
	log.Printf("Received UserID: %v", in.GetUser())
	log.Printf("Received R1: %v", in.GetR1())
	log.Printf("Received R2: %v", in.GetR2())

	// Retrieve User from the map
	_, exists := srv.userRegData[in.GetUser()]
	if !exists {
		return &pb.AuthenticationChallengeResponse{}, fmt.Errorf("user doesn't exists")
	}

	r1 := new(big.Int)
	r1.SetString(in.GetR1(), 10)
	r2 := new(big.Int)
	r2.SetString(in.GetR2(), 10)

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

	srv.authenticationData[authId] = Authentication{
		user: in.GetUser(),
		r1:   r1,
		r2:   r2,
		c:    c,
	}

	return &pb.AuthenticationChallengeResponse{AuthId: authId, C: c.String()}, nil

}

// This is the third step in the authentication process
// This is where the verifier proofs authentication with no knowledge of the secret x
func (srv *server) VerifyAuthentication(ctx context.Context, in *pb.AuthenticationAnswerRequest) (*pb.AuthenticationAnswerResponse, error) {
	log.Printf("Received AuthID: %v", in.GetAuthId())
	log.Printf("Received S: %v", in.GetS())

	s := new(big.Int)
	s.SetString(in.GetS(), 10)
	// Retrieve Auth object from map
	auth, exists := srv.authenticationData[in.GetAuthId()]
	if !exists {
		return &pb.AuthenticationAnswerResponse{}, fmt.Errorf("authId doesn't exists: %v", in.GetAuthId())
	}

	// Retrieve User from the map
	user, exists := srv.userRegData[auth.user]
	if !exists {
		return &pb.AuthenticationAnswerResponse{}, fmt.Errorf("user doesn't exists: %v", auth.user)
	}

	// Now we have all the data we need to validate the proof
	// Now the verifier needs to verify the proof
	// r1 = g^s . y1^c mod p
	_, err := zkpautils.VerifyProof(auth.r1, big.NewInt(*gFlag), s, user.y1, auth.c, big.NewInt(*pFlag))
	if err != nil {
		return &pb.AuthenticationAnswerResponse{}, fmt.Errorf("r1 does not match: %v", err)
	}

	// r2 = h^s . y2^c mod p
	_, err = zkpautils.VerifyProof(auth.r2, big.NewInt(*hFlag), s, user.y2, auth.c, big.NewInt(*pFlag))
	if err != nil {
		log.Fatal("r2 does not match", err)
		return &pb.AuthenticationAnswerResponse{}, fmt.Errorf("r2 does not match: %v", err)
	}

	log.Println("Proof verified!")

	// Now we mint a sessionID
	sessionId := zkpautils.RandomString(20)

	// Now we store the sessionID against the user, with a createdAt timestamp
	// for the future
	srv.sessionData[sessionId] = Session{
		user:      auth.user,
		createdAt: time.Now(),
	}

	// This deletes the old authentication data object
	// as we don't want to use it again.
	delete(srv.authenticationData, in.GetAuthId())

	return &pb.AuthenticationAnswerResponse{SessionId: sessionId}, nil
}

func main() {
	flag.Parse()

	// creating bigInts from the flags
	p := big.NewInt(*pFlag)
	q := big.NewInt(*qFlag)
	g := big.NewInt(*gFlag)
	h := big.NewInt(*hFlag)

	log.Printf("p: %v q: %v g: %v h: %v\n", p, q, g, h)

	// This makes sure that we validate the public variables passed in
	_, err := zkpautils.ValidatePublicVariables(p, q, g, h)
	if err != nil {
		log.Fatalf("could not validate public variables: %v", err)
	}
	// The config is now validated and in good shape

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *portFlag))
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
