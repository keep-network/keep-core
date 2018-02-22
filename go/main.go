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

	// wait for heartbeats to build mesh
	time.Sleep(time.Second * 2)

	go func(n *node.Node) {
		// first tick happens immediately
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
				log.Printf("GOT: %+v", got)
				log.Printf("GOT FROM: %+v", got.GetFrom())
				log.Printf("GOT Data: %s", got.GetData())
				log.Printf("GOT Seqno: %d", got.GetSeqno())
				log.Printf("GOT TopicIDs: %d", got.GetTopicIDs())
				t.Reset(5 * time.Second)
			}
		}
	}(n)

	go func(n *node.Node) {
		t := time.NewTimer(1)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				peers := n.Network.Sub.ListPeers("")
				for _, peer := range peers {
					log.Printf("Connected to peer: %s\n", peer)
				}
				t.Reset(5 * time.Second)
			}
		}
	}(n)

	select {}
}
