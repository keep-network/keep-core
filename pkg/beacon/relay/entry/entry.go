package entry

import (
	"math/big"
	"time"
)

// Request represents a request for an entry in the threshold relay.
type Request struct {
	previousEntry Entry

	RequestID   *big.Int
	Payment     *big.Int
	BlockReward *big.Int
	Seed        *big.Int
}

// Entry represents one entry in the threshold relay.
type Entry struct {
	RequestID     *big.Int
	Value         [32]byte
	GroupID       *big.Int
	PreviousEntry *big.Int
	Timestamp     time.Time
}
