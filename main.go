package main

import (
	"fmt"
	"log"
	"math/big"

	"github.com/fxtlabs/primes"
)

func main() {
	fmt.Println("Hello, World!")

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

	// This validates that q divides p - 1 evenly
	pDivQ := new(big.Int).Div(new(big.Int).Sub(p, big.NewInt(1)), q)
	if pDivQ.Mod(big.NewInt(0), big.NewInt(2)) == big.NewInt(0) {
		log.Fatalf("q:'%d' needs to divide evenly to p-1 where p:'%d'", q, p)
	}

	// We need to validate that g and h are in the same group
	// And that they both have order q

	// g and h are have the same prime order
	// g^q mod p = 1 AND h^q mod p = 1

	gPowQ := new(big.Int).Exp(g, q, nil)

	fmt.Println(gPowQ)

	gPowQModp := new(big.Int).Mod(gPowQ, p)

	fmt.Println(gPowQModp)

	if new(big.Int).Mod(new(big.Int).Exp(g, q, nil), p).Cmp(big.NewInt(1)) != 0 {
		log.Fatalf("g:'%d' does not have order q:'%d'", g, q)
	}

	if new(big.Int).Mod(new(big.Int).Exp(h, q, nil), p).Cmp(big.NewInt(1)) != 0 {
		log.Fatalf("h:'%d' does not have order q:'%d'", g, q)
	}

	// At this point all of our input values are setup and valid

	// We need to set a secret value for Peggy
	x := big.NewInt(6)

	// Now to calculate y1
	y1 := new(big.Int).Mod(new(big.Int).Exp(g, x, nil), p)

	// Now to calculate y2
	y2 := new(big.Int).Mod(new(big.Int).Exp(h, x, nil), p)

	// This needs to be sent to Victor in the registration phase
	fmt.Printf("Peggy sends y1: %d and y2: %d\n", y1, y2)

	// Peggy needs to pick a random k
	k := big.NewInt(7)

	// Now to calculate (r1, r2) = g^k, h^k
	r1 := new(big.Int).Mod(new(big.Int).Exp(g, k, nil), p)
	r2 := new(big.Int).Mod(new(big.Int).Exp(h, k, nil), p)

	fmt.Printf("Peggy sends r1: %d and r2: %d\n", r1, r2)

	// Now the challenger picks a value c
	c := big.NewInt(4)

	// The prover needs to then compute s
	// s = (k - c .x) mod q
	s := new(big.Int).Mod(new(big.Int).Sub(k, new(big.Int).Mul(c, x)), q)

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
