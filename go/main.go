package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	relay "github.com/keep-network/keep-core/go/relayclient"
)

func main() {
	initLibP2P()
}

func initLibP2P() {
	// Parse options from the command line
	var listenF = flag.Int("l", 0, "wait for incoming connections")
	var seed = flag.Int64("seed", 0, "set random seed for id generation")

	flag.Parse()

	ctx := context.Background()
	n, err := relay.NewRelayClient(ctx, *listenF, *seed)
	if err != nil {
		log.Fatalf("Failed to initialize relay node with: ", err)
	}

	log.Printf("New node: %+v", n)
	log.Printf("Node is operational.")

	// Subscribe all peers to topic
	log.Printf("Current Group state: %+v\n", n.Groups.GetActiveGroups())

	topic := "x"

	// Join the Group "x"
	err = n.Groups.JoinGroup(ctx, topic)
	if err != nil {
		log.Fatalf("Failed to subscribe to channel with err: ", err)
	}

	// wait for heartbeats to build mesh
	time.Sleep(time.Second * 2)

	go func(ctx context.Context, n *relay.RelayClient, topic string) {
		// first tick happens immediately
		t := time.NewTimer(1)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				// Make a message to send
				r := rand.Intn(100 + 1)
				msg := fmt.Sprintf("keep group message %d from %s", r, n.Identity.PeerID)
				m := &relay.Message{
					Data: msg,
				}
				// Get the group we care about
				g, err := n.Groups.GetGroup(ctx, topic)
				if err != nil {
					log.Fatalf("Failed to group with err: ", err)
				}
				// Send a message to the group...
				if ok := g.Send(m); !ok {
					log.Fatalf("Failed to send message: %#v\n", m)
				}
				// ...every 5 seconds
				t.Reset(5 * time.Second)
			case <-ctx.Done():
				return
			}
		}
	}(ctx, n, topic)

	go func(ctx context.Context, n *relay.RelayClient, topic string) {
		// t := time.NewTimer(1)
		// defer t.Stop()

		// Get the group we care about
		g, err := n.Groups.GetGroup(ctx, topic)
		if err != nil {
			log.Fatalf("Failed to group with err: ", err)
		}

		for {
			select {
			// Get a message from the group...
			case msg := <-g.RecvChan():
				log.Printf("Group msg: %#v\n", msg)
			case <-ctx.Done():
				return
			}
		}
	}(ctx, n, topic)

	select {}
}
