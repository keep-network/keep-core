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
	blockHeight         int
	subscriptionChannel chan block
	config              *ethereumChain
	waiters             map[int][]chan int
}

type block struct {
	Number string
}

// WaitForBlocks waits for a minimum of 1 block before returing.
func (blockWait *ethereumBlockCounter) WaitForBlocks(numBlocks int) error {
	waiter, err := blockWait.BlockWaiter(numBlocks)
	if err != nil {
		return err
	}
	<-waiter
	return nil
}

// BlockWaiter returns the block number as a chanel with a minimum of 1 block
// wait. 0 and negative numBlocks are converted to 1.
func (blockWait *ethereumBlockCounter) BlockWaiter(numBlocks int) (<-chan int, error) {
	newWaiter := make(chan int)

	blockWait.structMutex.Lock()
	defer blockWait.structMutex.Unlock()

	notifyBlockHeight := blockWait.blockHeight + numBlocks

	if notifyBlockHeight <= blockWait.blockHeight {
		go func() { newWaiter <- notifyBlockHeight }()
	} else {
		waiterList, exists := blockWait.waiters[notifyBlockHeight]
		if !exists {
			waiterList = make([]chan int, 0)
		}

		blockWait.waiters[notifyBlockHeight] = append(waiterList, newWaiter)
	}

	return newWaiter, nil
}

// receiveBlocks gets each new block back from Geth and extracts the
// block height (top) form it.  For each block height that is being
// weighted on a message will be sent.
func (blockWait *ethereumBlockCounter) receiveBlocks() {
	for block := range blockWait.subscriptionChannel {
		if top, err := strconv.ParseInt(block.Number, 0, 64); err == nil {
			if err != nil {
				fmt.Printf("Error: %s\n", err)
			}
			for ii := blockWait.blockHeight; ii <= int(top); ii++ {

				blockWait.structMutex.Lock()
				height := blockWait.blockHeight
				blockWait.blockHeight++
				waiters, exists := blockWait.waiters[height]
				delete(blockWait.waiters, height)
				blockWait.structMutex.Unlock()

				if exists {
					for _, waiter := range waiters {
						go func(w chan int) { w <- height }(waiter)
					}
				}
			}
		}
	}

}

// subscribeBlocks creates a subscription to Geth to get each block.
func (blockWait *ethereumBlockCounter) subscribeBlocks() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := blockWait.config.clientWS.EthSubscribe(
		ctx,
		blockWait.subscriptionChannel,
		"newHeads",
	)
	if err != nil {
		return
	}

	var lastBlock block
	err = blockWait.config.clientRPC.Call(&lastBlock, "eth_getBlockByNumber",
		"latest", true)
	if err != nil {
		return
	}

	blockWait.subscriptionChannel <- lastBlock

}

// BlockCounter creates a BlockCounter that uses the block number in ethereum.
// func BlockCounter(config chain.Provider) chain.BlockCounter {
func (ec *ethereumChain) BlockCounter() (chain.BlockCounter, error) {
	blockWait := ethereumBlockCounter{
		blockHeight: 0,
		waiters:     make(map[int][]chan int),
		config:      ec,
	}

	var lastBlock block
	err := blockWait.config.clientRPC.Call(&lastBlock, "eth_getBlockByNumber",
		"latest", true)
	if err != nil {
		return nil,
			fmt.Errorf("Failed to get initial number of blocks from the chain: %s",
				err)
	}

	var ii int64
	if ii, err = strconv.ParseInt(lastBlock.Number, 0, 64); err != nil {
		return nil,
			fmt.Errorf("Failed to get initial number of blocks from the chain, %s",
				err)
	}

	blockWait.blockHeight = int(ii)
	blockWait.subscriptionChannel = make(chan block)

	go func() {
		for i := 0; ; i++ {
			if i > 0 {
				time.Sleep(2 * time.Second)
			}
			blockWait.subscribeBlocks()
		}
	}()

	go blockWait.receiveBlocks()

	return &blockWait, nil
}
