package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/dfinity/go-dfinity-crypto/rand"
	"github.com/keep-network/keep-core/go/node"
)

func main() {
	// Parse options from the command line
	listenF := flag.Int("l", 0, "wait for incoming connections")
	seed := flag.Int64("seed", 0, "set random seed for id generation")
	flag.Parse()

	r := rand.NewRand()
	fmt.Printf("%v\n", r)

	n := node.NewNode(context.Background(), *listenF, *seed)
	log.Printf("New node: %+v", n)

	log.Printf("Node is operational.")
	select {}
	// TODO: defer func() { node.GracefulShutdown() }()
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
