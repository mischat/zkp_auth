package main

import (
	"fmt"
	"log"
	"math/big"

	"github.com/fxtlabs/primes"
	zkal "github.com/mischat/zkp_auth/lib"
)

func main() {
	fmt.Println("Hello, this is a walk through of the algorithm!")
	//Firstly let's get a prime number, and then we need to find a group.
	// Generate the prime numbers less than or equal to 20
	ps := primes.Sieve(20000)
	fmt.Println(ps)

	// We will then go ahead and create a group from the prime number
	// p is a prime number
	p := big.NewInt(23)
	q := big.NewInt(11)

	g := big.NewInt(4)
	h := big.NewInt(9)

	valid, err := zkal.ValidatePublicVariables(p, q, g, h)
	if err != nil {
		log.Fatal(err)
	}
	// At this point all of our input values are setup and valid
	fmt.Printf("Valid: %t\n", valid)

	// We need to set a secret value for Peggy
	x := big.NewInt(6)

	// Now to calculate y1
	y1 := zkal.CalculateExp(g, x, p)

	// Now to calculate y2
	y2 := zkal.CalculateExp(h, x, p)

	// This needs to be sent to Victor in the registration phase
	fmt.Printf("Peggy sends y1: %d and y2: %d\n", y1, y2)

	// Peggy needs to pick a random k
	k := big.NewInt(7)

	// Now to calculate (r1, r2) = g^k, h^k
	r1 := zkal.CalculateExp(g, k, p)
	r2 := zkal.CalculateExp(h, k, p)

	fmt.Printf("Peggy sends r1: %d and r2: %d\n", r1, r2)

	// Now the challenger picks a value c
	c := big.NewInt(4)

	// The prover needs to then compute s
	// s = (k - c .x) mod q
	s := zkal.CalculateS(k, c, x, q)

	fmt.Printf("Peggy sends s: %d \n", s)

	// Now the challenger needs to verify the proof
	// r1 = g^s . y1^c mod p
	// r2 = h^s . y2^c mod p

	if r1.Cmp(new(big.Int).Mod(new(big.Int).Mul(new(big.Int).Exp(g, s, nil), new(big.Int).Exp(y1, c, nil)), p)) != 0 {
		log.Fatalf("r1 does not match")
	}

	if r2.Cmp(new(big.Int).Mod(new(big.Int).Mul(new(big.Int).Exp(h, s, nil), new(big.Int).Exp(y2, c, nil)), p)) != 0 {
		log.Fatalf("r2 does not match")
	}

	fmt.Println("Proof verified")

}
