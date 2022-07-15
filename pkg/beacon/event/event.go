// Package event contains data structures that are attached to events in the
// relay. Though many of these events are triggered on-chain, that is not an
// inherent requirement of structures in this package.
package event

import (
	"math/big"
)

// RelayEntrySubmitted indicates that valid relay entry has been submitted to
// the chain for the currently processed relay request. This event is intended
// to be used by operators for tracking entry generation and submission progress.
// TODO: Adjust to the v2 RandomBeacon contract.
type RelayEntrySubmitted struct {
	BlockNumber uint64
}

// RelayEntryRequested represents a request for an entry in the threshold relay.
// TODO: Adjust to the v2 RandomBeacon contract.
type RelayEntryRequested struct {
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
// TODO: Adjust to the v2 RandomBeacon contract and rename to GroupRegistered.
type GroupRegistration struct {
	GroupPublicKey []byte

	BlockNumber uint64
}

// DKGResultSubmission represents a DKG result submission event. It is emitted
// after a submitted DKG result is positively validated on the chain. It contains
// the index of the member who submitted the result and a final public key of
// the group.
// TODO: Adjust to the v2 RandomBeacon contract and rename to DKGResultSubmitted.
type DKGResultSubmission struct {
	MemberIndex    uint32
	GroupPublicKey []byte
	Misbehaved     []uint8

	BlockNumber uint64
}
