package main

import (
	"fmt"

	"github.com/dfinity/go-dfinity-crypto/rand"
	"github.com/keep-network/go-dfinity-crypto/bls"
	"github.com/keep-network/keep-core/go/beacon/broadcast"
	"github.com/keep-network/keep-core/go/beacon/chain"
	"github.com/keep-network/keep-core/go/beacon/dkg"
)

func main() {
	r := rand.NewRand()
	fmt.Printf("%v\n", r)
}
