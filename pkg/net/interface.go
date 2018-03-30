package net

import (
	"github.com/dfinity/go-dfinity-crypto/bls"
)

// GroupIdentity contains the Group's public key as created by the dkg process
// and a list of Members that belong to the group.
// TODO: move to a more appropriate package; revise which fields we need
type GroupIdentity struct {
	Name string
	// Public key for the group; created through DKG.
	// Verified from the on-chain Group Registry
	GroupPublicKey *bls.PublicKey
	// The final list of qualified Group Members; empty if not yet computed
	// Verified from the on-chain Group Registry
	Members []bls.PublicKey
}

// Message corresponds with our proto Envelope type.
type Message struct{}

// HandleMessageFunc is the type of function called for each Message m furnished by
// the BroadcastChannel. If there is a problem handling the Message, the incoming error will
// describe the problem and the function can decide how to handle that error. If an error is returned,
// processing stops.
type HandleMessageFunc func(m Message) error

// BroadcastChannel represents a named pubsub channel. It allows Group Members
// to send messages on the channel (via Send), and to access a low-level receive chan
// that furnishes messages sent onto the BroadcastChannel.
type BroadcastChannel interface {
	// Return the member list and identifiying information for the Group
	GroupIdentity(name string) *GroupIdentity
	// Given a Message m, broadcast m to members of the Group through the BroadcastChannel
	Send(m Message) error
	// Recv takes a HandleMessageFunc and returns an error. This function should be retried.
	Recv(h HandleMessageFunc) error
}
