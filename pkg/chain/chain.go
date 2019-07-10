package chain

import (
	"math/big"

	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/gen/async"
)

// BlockCounter is an interface that provides the ability to wait for a certain
// number of abstract blocks. It provides for two ways to wait, one blocking and
// one chan-based. Block height is expected to increase monotonically, though
// the time between blocks will depend on the underlying implementation. See
// LocalBlockCounter() for a local implementation.
type BlockCounter interface {
	// WaitForBlockHeight blocks at the caller until the given block height is
	// reached. If the number of blocks is zero or negative or if the given
	// block height has been already reached, it returns immediately.
	WaitForBlockHeight(blockNumber uint64) error

	// BlockHeightWaiter returns a channel that will emit the block number after
	// the given block height is reached and then immediately close.
	// Reading from the returned channel immediately will effectively behave the
	// same way as calling WaitForBlockHeight.
	BlockHeightWaiter(blockNumber uint64) (<-chan uint64, error)

	// CurrentBlock returns the current block height.
	CurrentBlock() (uint64, error)
}

// StakeMonitor is an interface that provides ability to check and monitor
// the stake for the provided address.
type StakeMonitor interface {
	// HasMinimumStake checks if the provided address staked enough to become
	// a network operator. The minimum stake is an on-chain parameter.
	HasMinimumStake(address string) (bool, error)
	// StakerFor returns a Staker for the given address.
	StakerFor(address string) (Staker, error)
}

// Handle represents a handle to a blockchain that provides access to the core
// operator functionality needed for Keep network interactions.
type Handle interface {
	BlockCounter() (BlockCounter, error)
	StakeMonitor() (StakeMonitor, error)
	ThresholdRelay() relaychain.Interface
}

// Utility represents a handle to a blockchain that provides access to certain
// utility functions for Keep network interactions. Notably, these functions can
// either be application or operator functionality, and they are generally not
// part of the day-to-day process of operating a Keep node.
type Utility interface {
	Handle

	Genesis() *async.RelayEntryPromise
	RequestRelayEntry(seed *big.Int) *async.RelayRequestPromise
}
