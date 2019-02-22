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

var blockTime = time.Duration(500 * time.Millisecond)

func (counter *localBlockCounter) WaitForBlocks(numBlocks int) error {
	waiter, err := counter.BlockWaiter(numBlocks)
	if err != nil {
		return err
	}
	<-waiter
	return nil
}

func (counter *localBlockCounter) BlockWaiter(numBlocks int) (<-chan int, error) {
	notifyBlockHeight := counter.blockHeight + numBlocks
	return counter.BlockHeightWaiter(notifyBlockHeight)
}

func (counter *localBlockCounter) WaitForBlockHeight(blockNumber int) error {
	waiter, err := counter.BlockHeightWaiter(blockNumber)
	if err != nil {
		return err
	}
	<-waiter
	return nil
}

func (counter *localBlockCounter) BlockHeightWaiter(
	blockNumber int,
) (<-chan int, error) {
	newWaiter := make(chan int)

	counter.structMutex.Lock()
	defer counter.structMutex.Unlock()

	if blockNumber <= counter.blockHeight {
		go func() { newWaiter <- blockNumber }()
	} else {
		waiterList, exists := counter.waiters[blockNumber]
		if !exists {
			waiterList = make([]chan int, 0)
		}

		counter.waiters[blockNumber] = append(waiterList, newWaiter)
	}

	return newWaiter, nil
}

func (counter *localBlockCounter) CurrentBlock() (int, error) {
	counter.structMutex.Lock()
	defer counter.structMutex.Unlock()
	return counter.blockHeight, nil
}

// count is an internal function that counts up time to simulate the generation
// of blocks.
func (counter *localBlockCounter) count() {
	ticker := time.NewTicker(blockTime)

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

// blockCounter creates a BlockCounter that runs completely locally. It is
// designed to simply increase block height at a set time interval in the
// background.
func blockCounter() (chain.BlockCounter, error) {
	counter := localBlockCounter{blockHeight: 0, waiters: make(map[int][]chan int)}

	go counter.count()

	return &counter, nil
}
