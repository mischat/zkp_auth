package main

import (
	"fmt"
	"log"
	"math/big"
)

// This script generates some bigger number that work
func main() {

	p := new(big.Int)
	p.SetString("115792089237316195423570985008687907852837564279074904382605163141518161494337", 10)

	q := new(big.Int)
	q.SetString("341948486974166000522343609283189", 10)

	r := new(big.Int)
	r.SetString("338624364920977752681389262317185522840540224", 10)

	h := new(big.Int)
	h.SetString("3141592653589793238462643383279502884197", 10)

	g := new(big.Int).Exp(h, r, p)

	fmt.Println("Here is g: ", g)

	if g.Cmp(big.NewInt(1)) == 0 {
		log.Fatalf("g should not be 1")
	}

	if new(big.Int).Exp(g, q, p).Cmp(big.NewInt(1)) != 0 {
		log.Fatalf("should be 1")
	}

	for i := big.NewInt(1); i.Cmp(big.NewInt(20)) < 0; i.Add(i, big.NewInt(1)) {
		group := new(big.Int).Exp(g, i, p)
		fmt.Println("HEre is a group: ", group)
	}

}
