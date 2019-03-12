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
	RequestID     *big.Int
	Value         *big.Int
	GroupPubKey   []byte
	PreviousEntry *big.Int
	Timestamp     time.Time
	Seed          *big.Int

	BlockNumber uint64
}

// Request represents a request for an entry in the threshold relay.
type Request struct {
	RequestID     *big.Int
	Payment       *big.Int
	BlockReward   *big.Int
	Seed          *big.Int
	PreviousValue *big.Int

	BlockNumber uint64
}

// GroupTicketSubmission represents a group selection ticket submission event.
type GroupTicketSubmission struct {
	TicketValue *big.Int

	BlockNumber uint64
}

// GroupRegistration represents a registered group in the threshold relay with a
// public key, that is considered active at ActivationBlockHeight, and was
// spawned by the relay request with id, RequestID.
type GroupRegistration struct {
	GroupPublicKey []byte
	RequestID      *big.Int

	BlockNumber uint64
}

// DKGResultSubmission represents a DKG result submission event. It is emitted
// after a submitted DKG result is positively validated on the chain. It contains
// the index of the member who submitted the result and a final public key of
// the group.
type DKGResultSubmission struct {
	RequestID      *big.Int
	MemberIndex    uint32
	GroupPublicKey []byte

	BlockNumber uint64
}
