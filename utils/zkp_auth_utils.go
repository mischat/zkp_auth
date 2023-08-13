package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"math/big"
)

// ValidatePublicVariables takes the public variables and validates them
func ValidatePublicVariables(p *big.Int, q *big.Int, g *big.Int, h *big.Int) (bool, error) {
	// p is a prime number
	if !p.ProbablyPrime(20) {
		return false, fmt.Errorf("p:'%d' is not prime", p)
	}

	// q is a prime number
	if !q.ProbablyPrime(20) {
		return false, fmt.Errorf("q:'%d' is not prime", q)
	}

	// This validates that q divides p - 1 evenly
	pDivQ := new(big.Int).Div(new(big.Int).Sub(p, big.NewInt(1)), q)
	if pDivQ.Mod(big.NewInt(0), big.NewInt(2)) == big.NewInt(0) {
		return false, fmt.Errorf("q:'%d' needs to divide evenly to p-1 where p:'%d'", q, p)
	}

	// g and h must have the same prime order
	// g^q mod p = 1
	gpq := new(big.Int).Exp(g, q, p)
	if gpq.Cmp(big.NewInt(1)) != 0 {
		return false, fmt.Errorf("g:'%d' does not have order q:'%d'", g, q)
	}

	// h^q mod p = 1
	hpq := new(big.Int).Exp(h, q, p)
	if hpq.Cmp(big.NewInt(1)) != 0 {
		return false, fmt.Errorf("h:'%d' does not have order q:'%d'", h, q)
	}

	return true, nil
}

// This function is used to calculate the initial (y1, y2) = g^x, h^x and (r1, r2) = g^k, h^k
func CalculateExp(gh *big.Int, xk *big.Int, p *big.Int) *big.Int {
	return new(big.Int).Mod(new(big.Int).Exp(gh, xk, nil), p)
}

// Prover needs to compute S with their random k and the challenger's c
// s = (k - c .x) mod q
func CalculateS(k *big.Int, c *big.Int, x *big.Int, q *big.Int) *big.Int {
	return new(big.Int).Mod(new(big.Int).Sub(k, new(big.Int).Mul(c, x)), q)
}

// This method is used to verfiy proof
// r1 = g^s . y1^c mod p
// r2 = h^s . y2^c mod p
func VerifyProof(r *big.Int, gh *big.Int, s *big.Int, y *big.Int, c *big.Int, p *big.Int) (bool, error) {
	if r.Cmp(new(big.Int).Mod(new(big.Int).Mul(new(big.Int).Exp(gh, s, nil), new(big.Int).Exp(y, c, nil)), p)) != 0 {
		log.Fatalf("r1 does not match")
		return false, fmt.Errorf("r:'%d' does not match gh:'%d' s:'%d' y:'%d' c:'%d' p:'%d'", r, gh, s, y, c, p)
	}

	return true, nil
}

// We're using the Rand function from the crypto/rand package
// This number should be big enough to be unique for this exercise
func RandomBigInt() *big.Int {
	// Max random value, a 16-bits integer, i.e 2^16 - 1
	max := new(big.Int)
	max.Exp(big.NewInt(2), big.NewInt(16), nil).Sub(max, big.NewInt(1))

	randInt, err := rand.Int(rand.Reader, max)
	if err != nil {
		panic(err)
	}
	return randInt
}

// Generate a random string of length n
// We will use these for string IDs
func RandomString(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(b)[:n]
}
