package main

import (
	"fmt"
	"log"
	"math/big"

	"github.com/fxtlabs/primes"
	zkpautils "github.com/mischat/zkp_auth/utils"
)

// Ultimately we are going to use this script to generate the public data we need
// To setup our ZKP Auth system
func main() {
	fmt.Println("Hello, this is a walk through of the algorithm!")
	// Firstly let's get a prime number, and then we need to find a group.
	// Generate the prime numbers less than or equal to 20
	ps := primes.Sieve(20000)
	fmt.Println(ps)

	// We will then go ahead and create a group from the prime number
	// p is a prime number
	p := big.NewInt(23)
	q := big.NewInt(11)

	g := big.NewInt(4)
	h := big.NewInt(9)

	valid, err := zkpautils.ValidatePublicVariables(p, q, g, h)
	if err != nil {
		log.Fatal(err)
	}
	// At this point all of our input values are setup and valid
	fmt.Printf("Valid: %t\n", valid)

	// We need to set a secret value for Peggy
	x := big.NewInt(6)

	// Now to calculate y1
	y1 := zkpautils.CalculateExp(g, x, p)

	// Now to calculate y2
	y2 := zkpautils.CalculateExp(h, x, p)

	// This needs to be sent to Victor in the registration phase
	fmt.Printf("Peggy sends y1: %d and y2: %d\n", y1, y2)

	// Peggy needs to pick a random k
	k := big.NewInt(7)

	// Now to calculate (r1, r2) = g^k, h^k
	r1 := zkpautils.CalculateExp(g, k, p)
	r2 := zkpautils.CalculateExp(h, k, p)

	fmt.Printf("Peggy sends r1: %d and r2: %d\n", r1, r2)

	// Now the challenger picks a value c
	c := big.NewInt(4)

	// The prover needs to then compute s
	// s = (k - c .x) mod q
	s := zkpautils.CalculateS(k, c, x, q)

	fmt.Printf("Peggy sends s: %d \n", s)

	// Now the challenger needs to verify the proof
	// r1 = g^s . y1^c mod p
	_, err = zkpautils.VerifyProof(r1, g, s, y1, c, p)
	if err != nil {
		log.Fatal("r1 does not match", err)
	}

	// r2 = h^s . y2^c mod p
	_, err = zkpautils.VerifyProof(r2, h, s, y2, c, p)
	if err != nil {
		log.Fatal("r2 does not match", err)
	}

	fmt.Println("Proof verified")

}
