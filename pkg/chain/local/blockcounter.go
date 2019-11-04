package local

import (
	"sync"
	"time"

	"github.com/keep-network/keep-core/pkg/chain"
)

type localBlockCounter struct {
	structMutex sync.Mutex
	blockHeight uint64
	waiters     map[uint64][]chan uint64
}

var blockTime = time.Duration(500 * time.Millisecond)

func (lbc *localBlockCounter) WaitForBlockHeight(blockNumber uint64) error {
	waiter, err := lbc.BlockHeightWaiter(blockNumber)
	if err != nil {
		return err
	}
	<-waiter
	return nil
}

func (lbc *localBlockCounter) BlockHeightWaiter(
	blockNumber uint64,
) (<-chan uint64, error) {
	newWaiter := make(chan uint64)

	lbc.structMutex.Lock()
	defer lbc.structMutex.Unlock()

	if blockNumber <= lbc.blockHeight {
		go func() { newWaiter <- blockNumber }()
	} else {
		waiterList, exists := lbc.waiters[blockNumber]
		if !exists {
			waiterList = make([]chan uint64, 0)
		}

		lbc.waiters[blockNumber] = append(waiterList, newWaiter)
	}

	return newWaiter, nil
}

func (lbc *localBlockCounter) CurrentBlock() (uint64, error) {
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
				go func(w chan uint64) { w <- height }(waiter)
			}
		}
	}
}

// BlockCounter creates a BlockCounter that runs completely locally. It is
// designed to simply increase block height at a set time interval in the
// background.
func BlockCounter() (chain.BlockCounter, error) {
	counter := localBlockCounter{blockHeight: 0, waiters: make(map[uint64][]chan uint64)}

	go counter.count()

	return &counter, nil
}
