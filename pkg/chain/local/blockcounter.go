package local

import (
	"sync"
	"time"

	"github.com/keep-network/keep-core/pkg/chain"
)

type localBlockCounter struct {
	structMutex sync.Mutex
	blockHeight int
	waiters     map[int][]chan int
}

// WaitForBlocks waits for the specified number of blocks before returning. If
// the number of blocks is zero or negative, it should return immediately.
func (counter *localBlockCounter) WaitForBlocks(numBlocks int) {
	waiter := counter.BlockWaiter(numBlocks)
	<-waiter
	return
}

// BlockWaiter immediately returns a channel that will receive the block number
// after the specified number of blocks. Reading from the returned channel
// immediately will effectively behave the same way as calling WaitForBlocks.
func (counter *localBlockCounter) BlockWaiter(numBlocks int) <-chan int {
	newWaiter := make(chan int)

	counter.structMutex.Lock()
	defer counter.structMutex.Unlock()
	notifyBlockHeight := counter.blockHeight + numBlocks

	if notifyBlockHeight <= counter.blockHeight {
		go func() { newWaiter <- notifyBlockHeight }()
	} else {
		waiterList, exists := counter.waiters[notifyBlockHeight]
		if !exists {
			waiterList = make([]chan int, 0)
		}

		counter.waiters[notifyBlockHeight] = append(waiterList, newWaiter)
	}

	return newWaiter
}

// count is an internal function that counts up time to simulate the generation
// of blocks.
func (counter *localBlockCounter) count() {
	ticker := time.NewTicker(time.Duration(500 * time.Millisecond))

	for range ticker.C {
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

// BlockCounter creates a BlockCounter that runs completely locally. It is
// designed to simply increase block height at a set time interval in the
// background.
func BlockCounter() chain.BlockCounter {
	counter := localBlockCounter{blockHeight: 0, waiters: make(map[int][]chan int)}

	go counter.count()

	return &counter
}
