package state

import (
	"context"
	"log"
	"os"

	floodsub "github.com/libp2p/go-floodsub"
)

type onChainMesssage struct{}

type groupMessage struct {
	*floodsub.Message
}

type State struct {
	ctx context.Context
	// we have events from our contract on chain
	chainC chan *onChainMesssage

	groupC chan *groupMessage

	ps *floodsub.PubSub

	// list of subscriptions (groups which a node belongs to)
	// TODO: add a lock here
	subs map[string]*floodsub.Subscription
}

// Lol globals
var NodeState *State

// TODO: this shouldn't all be done in here
func init() {
	// Welcome to Keep!
	// TODO: pull in environment variables ie. staking address
	// TODO: add cli!
	_ = mustGetenv("KEEP_STAKING_ADDR")

	var err error
	NodeState, err = NewState(context.Background())
	if err != nil {
		panic("Something bad happened")
	}

	// TODO: sync with network, by calling abi to see if a user is staked
}

// called only on init
func NewState(ctx context.Context) (*State, error) {
	st := &State{}
	go st.eventLoop(ctx)
	return st, nil
}

func (st *State) eventLoop(ctx context.Context) {
	for {
		select {
		// case join group
		case <-st.groupC:
			// we have a message from a group:
			// send off to a async busy loop that's processing group messages
		case <-st.chainC:
			// we have a message from a chain:
			// send off to a async busy loop that's processing chain messages
		case <-ctx.Done():
			// shutdown from server - could make this group dissolution as well
		default:
			// block main thread
			// TODO: does this just eat up memory? maybe restructuring this to select{}
			// is better?
		}
	}
}

func (st *State) handleSubscriptions() {
	for {
		// FIXME: find a better way to do this than to enumerate all subs in a busy loop
		for _, group := range st.subs {
			msg, err := group.Next(st.ctx)
			if err != nil {
				// TODO: handle errors
				// TODO: better logging
				log.Println("Error: ", err)
				return
			}

			// TODO: process the message
			// TODO: better logging
			log.Println("Message: ", msg)
			st.groupC <- &groupMessage{msg}
		}
	}

}

func mustGetenv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("%s must be set", key)
	}
	return v
}
