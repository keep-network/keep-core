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
