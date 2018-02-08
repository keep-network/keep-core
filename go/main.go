package main

import (
	"fmt"

	"github.com/dfinity/go-dfinity-crypto/rand"
)

func main() {
	r := rand.NewRand()
	fmt.Printf("%v\n", r)
	// 1. initialize a node
	// 2. start a node
	// 3. create a gatewayListener
	// 4. block with waitgroups and channels
	// 5. if we get an error, shutdown
	// extra: go routine that listens for a sigint of somesort
}

// func gatewayListener(ctx context.Context) (<-chan error, error) {
// 	mm, err := ma.NewMultiaddr("some meaningful addr")
// 	if err != nil {
// 		return nil, fmt.Error("some meaningful err")
// 	}

// 	listener, err := manet.Listen(mm)
// 	if err != nil {
// 		return nil, fmt.Error("some meaningful err")
// 	}
// 	mm = listener.Multiaddr()
// 	// ref the node here, start with global, shift to a field on some command struct
// 	node, err := node.NewNode(...)
// 	if err != nil {
// 		return nil, fmt.Error("some meaningful err")
// 	}
// 	errc := make(chan err)
// 	go func() {
// 		errc <- someBlockingFunction()
// 	        close(errc)
// 	}()
// 	return errc, nil
// }
