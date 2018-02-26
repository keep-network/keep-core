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

func initLibP2P() {
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
	log.Printf("Current Group state: %+v\n", n.Groups.GetActiveGroups())

	ctx := context.Background()
	topic := "x"
	err = n.Groups.JoinGroup(ctx, topic)
	if err != nil {
		log.Fatalf("Failed to subscribe to channel with err: ", err)
	}

	// wait for heartbeats to build mesh
	time.Sleep(time.Second * 2)

	go func(ctx context.Context, n *node.Node, topic string) {
		// first tick happens immediately
		t := time.NewTimer(1)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				r := rand.Intn(100 + 1)
				msg := fmt.Sprintf("keep group message %d from %s", r, n.Identity.PeerID)
				err := n.Groups.BroadcastGroupMessage(ctx, n.Identity.PrivKey, topic, []byte(msg))
				if err != nil {
					log.Fatalf("Failed to get message with err: ", err)
				}
				t.Reset(5 * time.Second)
			}
		}
	}(ctx, n, topic)

	go func(n *node.Node) {
		t := time.NewTimer(1)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				for _, group := range n.Groups.GetActiveGroups() {
					log.Printf("Current Group state: %#v\n", group)
				}
				t.Reset(3 * time.Second)
			}
		}
	}(n)

	select {}
}
