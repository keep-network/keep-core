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
	Value         [32]byte
	GroupID       *big.Int
	PreviousEntry *big.Int
	Timestamp     time.Time
}

// Request represents a request for an entry in the threshold relay.
type Request struct {
	previousEntry Entry

	RequestID   *big.Int
	Payment     *big.Int
	BlockReward *big.Int
	Seed        *big.Int
}

// GroupRegistration represents a registered group in the threshold relay with a
// public key, that is considered active at ActivationBlockHeight, and was
// spawned by the relay request with id, RequestID.
type GroupRegistration struct {
	GroupPublicKey        []byte
	RequestID             *big.Int
	ActivationBlockHeight *big.Int
}
