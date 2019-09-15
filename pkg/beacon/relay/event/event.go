// Package event contains data structures that are attached to events in the
// relay. Though many of these events are triggered on-chain, that is not an
// inherent requirement of structures in this package.
package event

import (
	"math/big"
	"time"
)

// Entry represents one entry in the threshold relay.
type Entry struct {
	Value         *big.Int
	GroupPubKey   []byte
	PreviousEntry *big.Int
	Timestamp     time.Time
	Seed          *big.Int

	BlockNumber uint64
}

// Request represents a request for an entry in the threshold relay.
type Request struct {
	PreviousEntry  *big.Int
	Seed           *big.Int
	GroupPublicKey []byte
	BlockNumber    uint64
}

// GroupSelectionStart represents a group selection start event.
type GroupSelectionStart struct {
	NewEntry    *big.Int
	BlockNumber uint64
}

// GroupTicketSubmission represents a group selection ticket submission event.
type GroupTicketSubmission struct {
	TicketValue *big.Int

	BlockNumber uint64
}

// GroupRegistration represents an event of registering a new group with the
// given public key.
type GroupRegistration struct {
	GroupPublicKey []byte

	BlockNumber uint64
}

// DKGResultSubmission represents a DKG result submission event. It is emitted
// after a submitted DKG result is positively validated on the chain. It contains
// the index of the member who submitted the result and a final public key of
// the group.
type DKGResultSubmission struct {
	MemberIndex    uint32
	GroupPublicKey []byte

	BlockNumber uint64
}
