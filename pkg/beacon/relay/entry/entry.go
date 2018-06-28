package entry

import (
	"time"
)

// Request represents a request for an entry in the threshold relay.
type Request struct {
	PreviousEntry Entry
}

// Entry represents one entry in the threshold relay.
type Entry struct {
	Value     [8]byte
	Timestamp time.Time
}
