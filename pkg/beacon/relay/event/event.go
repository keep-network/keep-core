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

// DKGResultPublication represents a DKG result publication event.
// TODO: Change name to: DKGResultSubmission
type DKGResultPublication struct {
	RequestID *big.Int
	// 	MemberIndex   int TODO: add member index
	GroupPublicKey []byte

	BlockNumber uint64
}

// DKGResultVote represents a DKG result voting event.
type DKGResultVote struct {
	RequestID     *big.Int
	MemberIndex   int
	DKGResultHash [32]byte

	BlockNumber uint64
}

// DKGResultElected represents event fired after result publication and voting
// phase when one of the results have been elected as the final one.
type DKGResultElected struct {
	RequestID      *big.Int
	GroupPublicKey []byte
	Success        bool

	BlockNumber uint64
}
