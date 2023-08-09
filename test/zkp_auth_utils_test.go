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
