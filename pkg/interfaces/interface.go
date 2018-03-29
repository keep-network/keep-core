package net

import (
	"github.com/dfinity/go-dfinity-crypto/bls"
)

// GroupState captures on of the many states a client may be in while
// in a peer-to-peer network on Keep.
type GroupState int

const (
	// Null state; does not apply if you're already in a group
	WaitingForGroup GroupState = iota
	// The process of setting up a group
	JoiningGroup
	// The group has successfully formed and is elligeble to handle a relay request
	ProcessingRequests
	// The group is in the process of disolving or has been dissolved
	GroupDissolved
)

// Network is our interface to the underlying p2p network that the client leverages
// TODO: Consider renaming to P2PNetwork
type Network interface {
	// Given a name for a Group, return the channel the group communicates over
	GetChannel(name string) BroadcastChannel
	// Given a name for a Group, return the state of the group as defined by an enum
	GroupStatus(name string) GroupState
	// For initialization; call Bootstap() to initiate a handshake and connection to
	// predefined bootstrap nodes
	Bootstrap() error
}

// TODO: move to a more appropriate package
type GroupIdentity struct {
	// Public key for the group
	GroupPublicKey *bls.PublicKey
	// The final list of qualified Group Members; empty if not yet computed
	qualifiedMembers []bls.ID
	// A map of group member ids and their group signature share
	receivedGroupSignatureShares map[bls.ID][]byte
}

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
