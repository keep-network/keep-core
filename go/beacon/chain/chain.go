package chain

import (
	"sync"
	"time"
)

// BlockCounter is an interface that provides the ability to wait for a certain
// number of abstract blocks. It provides for two ways to wait, one blocking and
// one chan-based. Block height is expected to increase monotonically, though
// the time between blocks will depend on the underlying implementation. See
// LocalBlockCounter() for a local implementation.
type BlockCounter interface {
	// WaitForBlocks blocks at the caller until numBlocks new blocks have been
	// seen.
	WaitForBlocks(numBlocks int)
	// BlockWaiter returns a channel that will emit the current block height
	// after the given number of blocks has elapsed and then immediately close.
	BlockWaiter(numBlocks int) <-chan int
}

type localBlockCounter struct {
	structMutex sync.Mutex
	blockHeight int
	waiters     map[int][]chan int
}

func (counter *localBlockCounter) WaitForBlocks(numBlocks int) {
	waiter := counter.BlockWaiter(numBlocks)
	<-waiter
	return
}

func (counter *localBlockCounter) BlockWaiter(numBlocks int) <-chan int {
	newWaiter := make(chan int)

	counter.structMutex.Lock()
	defer counter.structMutex.Unlock()
	notifyBlockHeight := counter.blockHeight + numBlocks

	if notifyBlockHeight <= counter.blockHeight {
		newWaiter <- notifyBlockHeight
	} else {
		waiterList, exists := counter.waiters[notifyBlockHeight]
		if !exists {
			waiterList = make([]chan int, 0)
		}

		counter.waiters[notifyBlockHeight] = append(waiterList, newWaiter)
	}

	return newWaiter
}

func (counter *localBlockCounter) count() {
	ticker := time.NewTicker(time.Duration(time.Second / 2))

	for _ = range ticker.C {
		counter.structMutex.Lock()
		counter.blockHeight++
		height := counter.blockHeight
		waiters, exists := counter.waiters[height]
		delete(counter.waiters, height)
		counter.structMutex.Unlock()

		if exists {
			for _, waiter := range waiters {
				go func(w chan int) { w <- height }(waiter)
			}
		}
	}
}

// LocalBlockCounter creates a BlockCounter that runs completely locally. It is
// designed to simply increase block height at a set time interval in the
// background.
func LocalBlockCounter() BlockCounter {
	counter := localBlockCounter{blockHeight: 0, waiters: make(map[int][]chan int)}

	go counter.count()

	return &counter
}

// BeaconConfig contains configuration for the threshold relay beacon, typically
// from the underlying blockchain.
type BeaconConfig struct {
	GroupSize int
	Threshold int
}

// GetBeaconConfig Get the latest threshold relay beacon configuration.
// TODO Make this actually look up/update from chain information.
func GetBeaconConfig() BeaconConfig {
	return BeaconConfig{10, 4}
}
