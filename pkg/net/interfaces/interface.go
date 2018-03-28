package net

import (
	"github.com/dfinity/go-dfinity-crypto/bls"
)

// GroupState captures on of the many states a client may be in while
// in a peer-to-peer network on Keep.
type GroupState int

const (
	JoinGroup GroupState = iota
	WaitForGroup
)

type Network interface {
	// Given a name for a Group, return the channel the group communicates over
	GetChannel(name string) Channel
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

type Channel interface {
	GroupIdentity(name string) *GroupIdentity
	Send(m Message) error
	Recv() Message
}

type Group struct{}

type GroupManager interface {
	GetGroup(name string) (*Group, error)
	JoinGroup(name string) error
	GetActiveGroups() []*Group
	DissolveGroup(name string) error
}
