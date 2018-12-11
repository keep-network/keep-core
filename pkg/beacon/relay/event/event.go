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
	GroupID       *big.Int
	PreviousEntry *big.Int
	Timestamp     time.Time
}

// Request represents a request for an entry in the threshold relay.
type Request struct {
	RequestID   *big.Int
	Payment     *big.Int
	BlockReward *big.Int
	Seed        *big.Int

	PreviousValue *big.Int
}

// GroupRegistration represents a registered group in the threshold relay with a
// public key, that is considered active at ActivationBlockHeight, and was
// spawned by the relay request with id, RequestID.
type GroupRegistration struct {
	GroupPublicKey        []byte
	RequestID             *big.Int
	ActivationBlockHeight *big.Int
}

// StakerRegistration is the data for the OnStakerAdded event.  This type may
// only be needed in Milestone 1 - it may change at Milestone 2.
type StakerRegistration struct {
	Index         int
	GroupMemberID string
}

// DKGResultPublication represents a DKG result publication event.
type DKGResultPublication struct {
	RequestID *big.Int
}

// DKGResultVote - PHASE 14
type DKGResultVote struct {
	RequestID *big.Int
}
