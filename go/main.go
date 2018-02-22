package main

import (
	"fmt"

	"github.com/dfinity/go-dfinity-crypto/bls"
	"github.com/dfinity/go-dfinity-crypto/rand"
)

func main() {
	bls.Init(bls.CurveFp254BNb)
	r := rand.NewRand()
	id := bls.ID{}
	id.SetHexString(r.String())
	fmt.Printf("%s %v\n", r, id)
}
