package spv

import (
	"time"
)

const (
	// Default value for history depth which is the number of blocks to look
	// back from the current block when searching for past deposit sweep
	// proposal submitted events. The value is the approximate number of
	// Ethereum blocks in a week, assuming one block is 12s.
	DefaultHistoryDepth = 50400

	// Default value for the limit of transactions returned for a given wallet
	// public key hash. The value is based on the frequency of how often deposit
	// sweep and redemption transaction will happen. Deposit sweep transactions
	// are assumed to happen every 48h. Redemption transactions are assumed to
	// happen every 3h. The wallet should refuse any proposals from the
	// coordinator if the previously executed Bitcoin transaction was not proved
	// to the Bridge yet so in theory, the value of 1 should be enough. We make
	// it a bit higher - better to be safe than sorry.
	DefaultTransactionLimit = 20

	// Default value for restart back-off time.
	DefaultRestartBackoffTime = 30 * time.Minute

	// Default value for idle back-off time.
	DefaultIdleBackOffTime = 10 * time.Minute
)

// Config holds configurable properties.
type Config struct {
	// Enabled indicates whether the SPV maintainer should be started.
	Enabled bool

	// HistoryDepth is the number of blocks to look back from the current block
	// when searching for past deposit sweep proposal submitted events.
	// To find Bitcoin transactions for which the SPV proof should be submitted,
	// the maintainer first inspects the sweep proposal events from the wallet
	// coordinator contract. This depth determines how far into the past the
	// system will consider events for processing. This value must not be too
	// high so that the event lookup is efficient. At the same time, this value
	// can not be too low to make sure all performed and not yet proven sweeps
	// can be found.
	HistoryDepth uint64

	// TransactionLimit sets the maximum number of confirmed transactions
	// returned when getting transactions for a public key hash. Once the
	// maintainer establishes the list of proposed sweeps based on the sweep
	// proposal events from the wallet coordinator contract, it needs to check
	// Bitcoin transactions executed by the wallet. Then, it tries to find the
	// sweep among those transactions. For example, if set to `20`, only the
	// latest twenty transactions will be returned. This value must not be too
	// high so that the transaction lookup is efficient. At the same time, this
	// value can not be too low to make sure the performed sweep can be found in
	// case the wallet decided to execute some other Bitcoin transaction after
	// the yet-not-proven sweep.
	TransactionLimit int

	// RestartBackoffTime is a restart backoff which should be applied when the
	// SPV maintainer is restarted. It helps to avoid being flooded with error
	// logs in case of a permanent error in the SPV maintainer.
	RestartBackoffTime time.Duration

	// IdleBackoffTime is a wait time which should be applied when there are no
	// more transaction proofs to submit.
	IdleBackoffTime time.Duration
}
