package main

import (
	"fmt"
	"github.com/dfinity/go-dfinity-crypto/rand"
)

func main() {
	r := rand.NewRand()
	fmt.Printf("%v\n", r)
}
