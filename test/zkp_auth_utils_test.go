package utils_test

import (
	"math/big"
	"testing"

	zkutils "github.com/mischat/zkp_auth/utils"
)

func TestCalculateExp(t *testing.T) {
	g := big.NewInt(2)
	x := big.NewInt(12345)
	p := big.NewInt(67890)
	expectedY := big.NewInt(1234)

	y := zkutils.CalculateExp(g, x, p)

	if y.Cmp(expectedY) != 0 {
		t.Errorf("CalculateExp(%v, %v, %v) = %v, expected %v", g, x, p, y, expectedY)
	}
}
