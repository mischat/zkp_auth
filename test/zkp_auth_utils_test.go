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
	// The schnorr group for p 23 and 1 11 is:
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

	// Test case 2: Invalid p
	p = big.NewInt(22)
	q = big.NewInt(11)
	g = big.NewInt(5) // not valid
	h = big.NewInt(7) // not valid
	valid, err = zkutils.ValidatePublicVariables(p, q, g, h)
	if valid || err == nil {
		t.Errorf("ValidatePublicVariables(%d, %d, %d, %d) = (%t, %v), expected (false, error)", p, q, g, h, valid, err)
	}

	// Test case 3: Invalid q
	p = big.NewInt(23)
	q = big.NewInt(10)
	g = big.NewInt(2)
	h = big.NewInt(7) // not valid
	valid, err = zkutils.ValidatePublicVariables(p, q, g, h)
	if valid || err == nil {
		t.Errorf("ValidatePublicVariables(%d, %d, %d, %d) = (%t, %v), expected (false, error)", p, q, g, h, valid, err)
	}

	// // Test case 4: Invalid q not dividing p-1 evenly
	// p = big.NewInt(23)
	// q = big.NewInt(7)
	// g = big.NewInt(2)
	// h = big.NewInt(8)
	// valid, err = zkutils.ValidatePublicVariables(p, q, g, h)
	// if valid || err == nil {
	// 	t.Errorf("ValidatePublicVariables(%d, %d, %d, %d) = (%t, %v), expected (false, error)", p, q, g, h, valid, err)
	// }

	// // Test case 5: Invalid g
	// p = big.NewInt(23)
	// q = big.NewInt(11)
	// g = big.NewInt(3)
	// h = big.NewInt(8)
	// valid, err = zkutils.ValidatePublicVariables(p, q, g, h)
	// if valid || err == nil {
	// 	t.Errorf("ValidatePublicVariables(%d, %d, %d, %d) = (%t, %v), expected (false, error)", p, q, g, h, valid, err)
	// }

	// // Test case 6: Invalid h
	// p = big.NewInt(23)
	// q = big.NewInt(11)
	// g = big.NewInt(2)
	// h = big.NewInt(9)
	// valid, err = zkutils.ValidatePublicVariables(p, q, g, h)
	// if valid || err == nil {
	// 	t.Errorf("ValidatePublicVariables(%d, %d, %d, %d) = (%t, %v), expected (false, error)", p, q, g, h, valid, err)
	// }
}
