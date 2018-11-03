package chain

import relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"

// BlockCounter is an interface that provides the ability to wait for a certain
// number of abstract blocks. It provides for two ways to wait, one blocking and
// one chan-based. Block height is expected to increase monotonically, though
// the time between blocks will depend on the underlying implementation. See
// LocalBlockCounter() for a local implementation.
type BlockCounter interface {
	// WaitForBlocks blocks at the caller until numBlocks new blocks have been
	// seen.
	WaitForBlocks(numBlocks int) error
	// BlockWaiter returns a channel that will emit the current block height
	// after the given number of blocks has elapsed and then immediately close.
	BlockWaiter(numBlocks int) (<-chan int, error)
	// GetBlockNo returns the block number.
	GetBlockNo() int
}

// Handle represents a handle to a blockchain that provides access to the core
// functionality needed for Keep network interactions.
type Handle interface {
	BlockCounter() (BlockCounter, error)
	ThresholdRelay() relaychain.Interface
}
