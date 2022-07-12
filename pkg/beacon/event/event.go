// Package event contains data structures that are attached to events in the
// relay. Though many of these events are triggered on-chain, that is not an
// inherent requirement of structures in this package.
package event

import (
	"math/big"
)

// EntrySubmitted indicates that valid relay entry has been submitted to the
// chain for the currently processed relay request. This event is intended to
// be used by operators for tracking entry generation and submission progress.
type EntrySubmitted struct {
	BlockNumber uint64
}

// Request represents a request for an entry in the threshold relay.
type Request struct {
	PreviousEntry  []byte
	GroupPublicKey []byte
	BlockNumber    uint64
}

// DKGStarted represents a DKG start event.
type DKGStarted struct {
	Seed        *big.Int
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
	Misbehaved     []byte

	BlockNumber uint64
}
