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

func (lbc *localBlockCounter) WaitForBlocks(numBlocks int) error {
	waiter, err := lbc.BlockWaiter(numBlocks)
	if err != nil {
		return err
	}
	<-waiter
	return nil
}

func (lbc *localBlockCounter) BlockWaiter(numBlocks int) (<-chan int, error) {
	notifyBlockHeight := lbc.blockHeight + numBlocks
	return lbc.BlockHeightWaiter(notifyBlockHeight)
}

func (lbc *localBlockCounter) WaitForBlockHeight(blockNumber int) error {
	waiter, err := lbc.BlockHeightWaiter(blockNumber)
	if err != nil {
		return err
	}
	<-waiter
	return nil
}

func (lbc *localBlockCounter) BlockHeightWaiter(
	blockNumber int,
) (<-chan int, error) {
	newWaiter := make(chan int)

	lbc.structMutex.Lock()
	defer lbc.structMutex.Unlock()

	if blockNumber <= lbc.blockHeight {
		go func() { newWaiter <- blockNumber }()
	} else {
		waiterList, exists := lbc.waiters[blockNumber]
		if !exists {
			waiterList = make([]chan int, 0)
		}

		lbc.waiters[blockNumber] = append(waiterList, newWaiter)
	}

	return newWaiter, nil
}

func (lbc *localBlockCounter) CurrentBlock() (int, error) {
	lbc.structMutex.Lock()
	defer lbc.structMutex.Unlock()
	return lbc.blockHeight, nil
}

// count is an internal function that counts up time to simulate the generation
// of blocks.
func (lbc *localBlockCounter) count() {
	ticker := time.NewTicker(blockTime)

	for range ticker.C {
		lbc.structMutex.Lock()
		lbc.blockHeight++
		height := lbc.blockHeight
		waiters, exists := lbc.waiters[height]
		delete(lbc.waiters, height)
		lbc.structMutex.Unlock()

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
