package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/keep-network/keep-core/go/node"
)

func main() {
	// Parse options from the command line
	listenF := flag.Int("l", 0, "wait for incoming connections")
	seed := flag.Int64("seed", 0, "set random seed for id generation")
	flag.Parse()

	n, err := node.NewNode(context.Background(), *listenF, *seed)
	if err != nil {
		log.Fatalf("Failed to initialize relay node with: ", err)
	}

	log.Printf("New node: %+v", n)
	log.Printf("Node is operational.")
	go func(n *node.Node) {
		t := time.NewTimer(1)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				n.Routing.PrintRoutingTable()
				t.Reset(1 * time.Second)
			}
		}
	}(n)

	// Subscribe all peers to topic
	subch, err := n.Floodsub.Subscribe("x")
	if err != nil {
		log.Fatalf("Failed to subscribe to channel with err: ", err)
	}
	msg := []byte(fmt.Sprintf("keep group message from %s", n.Identity.PeerID))
	n.Floodsub.Publish("x", msg)
	got, err := subch.Next(context.Background())
	if err != nil {
		log.Fatalf("Failed to get message with err: ", err)
	}
	log.Println(got)
	select {}
}
