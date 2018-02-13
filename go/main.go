package main

import (
	"context"
	"flag"
	"log"

	"github.com/keep-network/keep-core/go/node"
)

func main() {
	// Parse options from the command line
	listenF := flag.Int("l", 0, "wait for incoming connections")
	seed := flag.Int64("seed", 0, "set random seed for id generation")
	flag.Parse()

	n := node.NewNode(context.Background(), *listenF, *seed)
	log.Printf("New node: %+v", n)

	log.Printf("Node is operational.")
	select {}
	// TODO: defer func() { node.GracefulShutdown() }()
}
