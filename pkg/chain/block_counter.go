package chain

import "context"

// BlockCounter is an interface that provides the ability to wait for a certain
// number of abstract blocks or watch as they are mined.
// Block height is expected to increase monotonically, though
// the time between blocks depends on the underlying implementation.
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

	// WatchBlocks returns a channel that will emit new block numbers as they
	// are mined. When the context provided as the parameter ends, new blocks
	// are no longer pushed to the channel and the channel is closed. If there
	// is no reader for the channel or reader is too slow, block updates can be
	// dropped.
	WatchBlocks(ctx context.Context) <-chan uint64
}
