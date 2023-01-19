package local_v1

import (
	"context"
	"sync"
	"time"

	"github.com/keep-network/keep-core/pkg/chain"
)

type localBlockCounter struct {
	structMutex sync.Mutex
	blockHeight uint64
	waiters     map[uint64][]chan uint64
	watchers    []*watcher
}

type watcher struct {
	ctx     context.Context
	channel chan uint64
}

var blockTime = time.Duration(12 * time.Second)

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

func (lbc *localBlockCounter) WatchBlocks(ctx context.Context) <-chan uint64 {
	watcher := &watcher{
		ctx:     ctx,
		channel: make(chan uint64, 1),
	}

	lbc.structMutex.Lock()
	lbc.watchers = append(lbc.watchers, watcher)
	lbc.structMutex.Unlock()

	go func() {
		<-ctx.Done()

		lbc.structMutex.Lock()
		for i, w := range lbc.watchers {
			if w == watcher {
				lbc.watchers[i] = lbc.watchers[len(lbc.watchers)-1]
				lbc.watchers = lbc.watchers[:len(lbc.watchers)-1]
				break
			}
		}
		lbc.structMutex.Unlock()
	}()

	return watcher.channel
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

		lbc.structMutex.Lock()
		watchers := make([]*watcher, len(lbc.watchers))
		copy(watchers, lbc.watchers)
		lbc.structMutex.Unlock()

		for _, watcher := range watchers {
			if watcher.ctx.Err() != nil {
				close(watcher.channel)
				continue
			}

			select {
			case watcher.channel <- height: // perfect
			default: // we don't care, let's drop it
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
