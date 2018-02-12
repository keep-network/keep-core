package chain

import (
	"sync"
	"time"
)

type BlockCounter interface {
	// WaitForBlocks blocks at the caller until numBlocks new blocks have been
	// seen.
	WaitForBlocks(numBlocks int)
	// BlockWaiter returns a channel that will emit the current block height
	// after the given number of blocks has elapsed and then immediately close.
	BlockWaiter(numBlocks int) <-chan int
}

type localBlockCounter struct {
	blockHeight int
	heightMutex sync.Mutex
	waiters     map[int][]chan int
}

func (counter *localBlockCounter) WaitForBlocks(numBlocks int) {
	waiter := counter.BlockWaiter(numBlocks)
	<-waiter
	return
}

func (counter *localBlockCounter) BlockWaiter(numBlocks int) <-chan int {
	newWaiter := make(chan int)

	counter.heightMutex.Lock()
	defer counter.heightMutex.Unlock()
	notifyBlockHeight := counter.blockHeight + numBlocks

	if notifyBlockHeight == counter.blockHeight {
		newWaiter <- notifyBlockHeight
	} else {
		waiterList, exists := counter.waiters[notifyBlockHeight]
		if !exists {
			waiterList = make([]chan int, 0)
			counter.waiters[notifyBlockHeight] = waiterList
		}

		counter.waiters[notifyBlockHeight] = append(waiterList, newWaiter)
	}

	return newWaiter
}

func (counter *localBlockCounter) count() {
	ticker := time.NewTicker(time.Duration(time.Second / 2))

	for _ = range ticker.C {
		counter.heightMutex.Lock()
		counter.blockHeight++
		waiters, exists := counter.waiters[counter.blockHeight]
		if exists {
			for _, waiter := range waiters {
				waiter <- counter.blockHeight
			}
			delete(counter.waiters, counter.blockHeight)
		}
		counter.heightMutex.Unlock()
	}
}

func LocalBlockCounter() BlockCounter {
	counter := localBlockCounter{blockHeight: 0, waiters: make(map[int][]chan int)}

	go counter.count()

	return &counter
}

type BeaconConfig struct {
	GroupSize int
	Threshold int
}

func GetBeaconConfig() BeaconConfig {
	return BeaconConfig{10, 4}
}
