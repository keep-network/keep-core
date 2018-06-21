package entry

import (
	"math/big"
	"time"
)

// Entry represents one entry in the threshold relay.
type Entry struct {
	Value     [8]byte
	Timestamp time.Time
}

// Request represents a request for an entry in the threshold relay.
type Request struct {
	id            *big.Int
	previousEntry Entry
}
