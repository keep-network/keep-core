package spv

import (
	"time"
)

// Config holds configurable properties.
type Config struct {
	// Enabled indicates whether the SPV maintainer should be started.
	Enabled bool

	// HistoryDepth is the number of blocks to look back from the current block
	// when searching for past deposit sweep proposal submitted events.
	// This depth determines how far into the past the system will consider
	// events for processing.
	HistoryDepth uint64

	// TransactionLimit sets the maximum number of confirmed transactions
	// returned when getting transactions for a public key hash. For example,
	// if set to `20`, only the latest twenty transactions will be returned.
	TransactionLimit int

	// RestartBackOffTime is a restart backoff which should be applied when the
	// SPV maintainer is restarted. It helps to avoid being flooded with error
	// logs in case of a permanent error in the SPV maintainer.
	RestartBackOffTime time.Duration

	// IdleBackOffTime is a wait time which should be applied when there are no
	// more transaction proofs to submit.
	IdleBackOffTime time.Duration
}
