package ethereum

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/keep-network/keep-core/pkg/chain"
)

type ethereumBlockCounter struct {
	structMutex         sync.Mutex
	latestBlockHeight   int
	subscriptionChannel chan block
	config              *ethereumChain
	waiters             map[int][]chan int
}

type block struct {
	Number string
}

// WaitForBlocks waits for a minimum of 1 block before returing.
func (ebc *ethereumBlockCounter) WaitForBlocks(numBlocks int) error {
	waiter, err := ebc.BlockWaiter(numBlocks)
	if err != nil {
		return err
	}
	<-waiter
	return nil
}

// BlockWaiter returns the block number as a channel with a minimum of 1 block
// wait. 0 and negative numBlocks are converted to 1.
func (ebc *ethereumBlockCounter) BlockWaiter(numBlocks int) (<-chan int, error) {
	newWaiter := make(chan int)

	ebc.structMutex.Lock()
	defer ebc.structMutex.Unlock()

	notifyBlockHeight := ebc.latestBlockHeight + numBlocks

	if notifyBlockHeight <= ebc.latestBlockHeight {
		go func() { newWaiter <- notifyBlockHeight }()
	} else {
		waiterList, exists := ebc.waiters[notifyBlockHeight]
		if !exists {
			waiterList = make([]chan int, 0)
		}

		ebc.waiters[notifyBlockHeight] = append(waiterList, newWaiter)
	}

	return newWaiter, nil
}

// receiveBlocks gets each new block back from Geth and extracts the
// block height (topBlockNumber) form it.  For each block height that is being
// waited on a message will be sent.
func (ebc *ethereumBlockCounter) receiveBlocks() {
	for block := range ebc.subscriptionChannel {
		if topBlockNumber, err := strconv.ParseInt(block.Number, 0, 64); err == nil {
			if err != nil {
				fmt.Printf("Error: %s\n", err)
			}
			latestBlockNumber := int(topBlockNumber)
			if latestBlockNumber == ebc.latestBlockHeight {
				continue
			}

			for unseenBlockNumber := ebc.latestBlockHeight; unseenBlockNumber <= latestBlockNumber; unseenBlockNumber++ {
				ebc.structMutex.Lock()
				height := ebc.latestBlockHeight
				ebc.latestBlockHeight++
				waiters := ebc.waiters[height]
				delete(ebc.waiters, height)
				ebc.structMutex.Unlock()

				for _, waiter := range waiters {
					go func(w chan int) { w <- height }(waiter)
				}
			}
		}
	}

}

// subscribeBlocks creates a subscription to Geth to get each block.
func (ebc *ethereumBlockCounter) subscribeBlocks() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := ebc.config.clientWS.EthSubscribe(
		ctx,
		ebc.subscriptionChannel,
		"newHeads",
	)
	if err != nil {
		return
	}

	var lastBlock block
	err = ebc.config.clientRPC.Call(
		&lastBlock,
		"eth_getBlockByNumber",
		"latest",
		true,
	)
	if err != nil {
		return err
	}

	ebc.subscriptionChannel <- lastBlock

	return nil
}

// BlockCounter creates a BlockCounter that uses the block number in ethereum.
func (ec *ethereumChain) BlockCounter() (chain.BlockCounter, error) {
	var startupBlock block
	err := ec.clientRPC.Call(
		&startupBlock,
		"eth_getBlockByNumber",
		"latest",
		true,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"Failed to get initial number of blocks from the chain: %s",
			err,
		)
	}

	startupBlockNumber, err := strconv.ParseInt(startupBlock.Number, 0, 32)
	if err != nil {
		return nil,
			fmt.Errorf("Failed to get initial number of blocks from the chain, %s",
				err)
	}

	blockCounter := &ethereumBlockCounter{
		latestBlockHeight:   int(startupBlockNumber),
		waiters:             make(map[int][]chan int),
		config:              ec,
		subscriptionChannel: make(chan block),
	}

	go func() {
		for i := 0; ; i++ {
			if i > 0 {
				time.Sleep(2 * time.Second)
			}
			blockCounter.subscribeBlocks()
		}
	}()

	go blockCounter.receiveBlocks()

	return blockCounter, nil
}
