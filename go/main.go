package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
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

	// Subscribe all peers to topic
	subch, err := n.Network.Sub.Subscribe("x")
	if err != nil {
		log.Fatalf("Failed to subscribe to channel with err: ", err)
	}

	go func(n *node.Node) {
		t := time.NewTimer(1)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				peers := n.Network.Sub.ListPeers("x")
				for _, peer := range peers {
					log.Printf("Connected to peer: %s\n", peer)
				}
				t.Reset(1 * time.Second)
			}
		}
	}(n)

	go func(n *node.Node) {
		t := time.NewTimer(1)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				r := rand.Intn(100-1) + 1
				msg := []byte(fmt.Sprintf("keep group message %d from %s", r, n.Identity.PeerID))
				n.Network.Sub.Publish("x", msg)
				got, err := subch.Next(context.Background())
				if err != nil {
					log.Fatalf("Failed to get message with err: ", err)
				}
				log.Println(got)
				t.Reset(5 * time.Second)
			}
		}
	}(n)

	select {}
}
