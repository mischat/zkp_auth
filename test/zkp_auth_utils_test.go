package utils_test

import (
	"math/big"
	"testing"

	zkutils "github.com/mischat/zkp_auth/utils"
)

func TestCalculateExp(t *testing.T) {
	// should pass
	g := big.NewInt(4)
	x := big.NewInt(6)
	p := big.NewInt(23)
	expectedY := big.NewInt(2)

	y := zkutils.CalculateExp(g, x, p)

	if y.Cmp(expectedY) != 0 {
		t.Errorf("CalculateExp(%v, %v, %v) = %v, expected %v", g, x, p, y, expectedY)
	}

	// should pass
	g = big.NewInt(3)
	x = big.NewInt(300)
	p = big.NewInt(10009)
	expectedY = big.NewInt(6419)

	y = zkutils.CalculateExp(g, x, p)

	if y.Cmp(expectedY) != 0 {
		t.Errorf("CalculateExp(%v, %v, %v) = %v, expected %v", g, x, p, y, expectedY)
	}

	// should pass
	g = big.NewInt(2892)
	x = big.NewInt(300)
	p = big.NewInt(10009)
	expectedY = big.NewInt(4984)

	y = zkutils.CalculateExp(g, x, p)

	if y.Cmp(expectedY) != 0 {
		t.Errorf("CalculateExp(%v, %v, %v) = %v, expected %v", g, x, p, y, expectedY)
	}
}

func TestValidatePublicVariables(t *testing.T) {
	// The schnorr group for p 23 and q 11 is:
	// {1, 3, 9, 4, 12, 13, 16, 2, 6, 18, 8}
	// Test case 1: Valid public variables
	p := big.NewInt(23)
	q := big.NewInt(11)
	g := big.NewInt(1)
	h := big.NewInt(3)
	valid, err := zkutils.ValidatePublicVariables(p, q, g, h)
	if !valid || err != nil {
		t.Errorf("ValidatePublicVariables(%d, %d, %d, %d) = (%t, %v), expected (%t, nil)", p, q, g, h, valid, err, true)
	}

	// The schnorr group for:
	// p: 115792089237316195423570985008687907852837564279074904382605163141518161494337
	// q: 341948486974166000522343609283189
	// run script "scripts/bignumbers/main.go" to generate Schnorr group for these values
	// Test case 2: Valid big public variables
	p = new(big.Int)
	p.SetString("115792089237316195423570985008687907852837564279074904382605163141518161494337", 10)

	q = new(big.Int)
	q.SetString("341948486974166000522343609283189", 10)

	g = new(big.Int)
	g.SetString("74446558554923317135296388588396736831887322850186029432124219757485062736903", 10)

	h = new(big.Int)
	h.SetString("79726485623116979445189935890227226532411986477410367519098002861237945910855", 10)

	valid, err = zkutils.ValidatePublicVariables(p, q, g, h)
	if !valid || err != nil {
		t.Errorf("ValidatePublicVariables(%d, %d, %d, %d) = (%t, %v), expected (%t, nil)", p, q, g, h, valid, err, true)
	}

	// Test case 3: Invalid g and h
	p = big.NewInt(23)
	q = big.NewInt(11)
	g = big.NewInt(5) // not valid
	h = big.NewInt(7) // not valid
	valid, err = zkutils.ValidatePublicVariables(p, q, g, h)
	if valid || err == nil {
		t.Errorf("ValidatePublicVariables(%d, %d, %d, %d) = (%t, %v), expected (false, error)", p, q, g, h, valid, err)
	}

	// Test case 4: Invalid h
	p = big.NewInt(23)
	q = big.NewInt(10)
	g = big.NewInt(2)
	h = big.NewInt(7) // not valid
	valid, err = zkutils.ValidatePublicVariables(p, q, g, h)
	if valid || err == nil {
		t.Errorf("ValidatePublicVariables(%d, %d, %d, %d) = (%t, %v), expected (false, error)", p, q, g, h, valid, err)
	}

	// Test case 4: Invalid g
	p = big.NewInt(23)
	q = big.NewInt(10)
	g = big.NewInt(7) // not valid
	h = big.NewInt(2)
	valid, err = zkutils.ValidatePublicVariables(p, q, g, h)
	if valid || err == nil {
		t.Errorf("ValidatePublicVariables(%d, %d, %d, %d) = (%t, %v), expected (false, error)", p, q, g, h, valid, err)
	}

	// Test case 5: Invalid p
	p = big.NewInt(22)
	q = big.NewInt(11)
	g = big.NewInt(1)
	h = big.NewInt(3)
	valid, err = zkutils.ValidatePublicVariables(p, q, g, h)
	if valid || err == nil {
		t.Errorf("ValidatePublicVariables(%d, %d, %d, %d) = (%t, %v), expected (false, error)", p, q, g, h, valid, err)
	}

	// Test case 6: Invalid q not dividing p-1 evenly
	p = big.NewInt(23)
	q = big.NewInt(7)
	g = big.NewInt(2)
	h = big.NewInt(8)
	valid, err = zkutils.ValidatePublicVariables(p, q, g, h)
	if valid || err == nil {
		t.Errorf("ValidatePublicVariables(%d, %d, %d, %d) = (%t, %v), expected (false, error)", p, q, g, h, valid, err)
	}
}
