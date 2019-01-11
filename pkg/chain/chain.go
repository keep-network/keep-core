package chain

import (
	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
)

// BlockCounter is an interface that provides the ability to wait for a certain
// number of abstract blocks. It provides for two ways to wait, one blocking and
// one chan-based. Block height is expected to increase monotonically, though
// the time between blocks will depend on the underlying implementation. See
// LocalBlockCounter() for a local implementation.
type BlockCounter interface {
	// WaitForBlocks blocks at the caller until numBlocks new blocks have been
	// seen. If the number of blocks is zero or negative, it returns
	// immediately.
	WaitForBlocks(numBlocks int) error
	// BlockWaiter returns a channel that will emit the current block height
	// after the given number of blocks has elapsed and then immediately close.
	// Reading from the returned channel immediately will effectively behave the
	// same way as calling WaitForBlocks.
	BlockWaiter(numBlocks int) (<-chan int, error)
	// CurrentBlock returns the current block height.
	CurrentBlock() (int, error)
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
// functionality needed for Keep network interactions.
type Handle interface {
	BlockCounter() (BlockCounter, error)
	StakeMonitor() (StakeMonitor, error)
	ThresholdRelay() relaychain.Interface
}
