package ethereum

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/keep-network/keep-core/pkg/chain"
)

const (
	// Subscription identifier to receive blocks as they are mined.
	newHeadsSubscription = "newHeads"
	// RPC call identifier to get a block's information given its number.
	getBlockByNumber = "eth_getBlockByNumber"
	// RPC call parameter to reference the latest block rather than a particular
	// number.
	latestBlock = "latest"
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

func (ebc *ethereumBlockCounter) WaitForBlocks(numBlocks int) error {
	waiter, err := ebc.BlockWaiter(numBlocks)
	if err != nil {
		return err
	}
	<-waiter
	return nil
}

func (ebc *ethereumBlockCounter) BlockWaiter(
	numBlocks int,
) (<-chan int, error) {
	notifyBlockHeight := ebc.latestBlockHeight + numBlocks
	return ebc.BlockHeightWaiter(notifyBlockHeight)
}

func (ebc *ethereumBlockCounter) WaitForBlockHeight(blockNumber int) error {
	waiter, err := ebc.BlockHeightWaiter(blockNumber)
	if err != nil {
		return err
	}
	<-waiter
	return nil
}

func (ebc *ethereumBlockCounter) BlockHeightWaiter(
	blockNumber int,
) (<-chan int, error) {
	newWaiter := make(chan int)

	ebc.structMutex.Lock()
	defer ebc.structMutex.Unlock()

	if blockNumber <= ebc.latestBlockHeight {
		go func() { newWaiter <- blockNumber }()
	} else {
		waiterList, exists := ebc.waiters[blockNumber]
		if !exists {
			waiterList = make([]chan int, 0)
		}

		ebc.waiters[blockNumber] = append(waiterList, newWaiter)
	}

	return newWaiter, nil
}

func (ebc *ethereumBlockCounter) CurrentBlock() (int, error) {
	return ebc.latestBlockHeight, nil
}

// receiveBlocks gets each new block back from Geth and extracts the
// block height (topBlockNumber) form it.  For each block height that is being
// waited on a message will be sent.
func (ebc *ethereumBlockCounter) receiveBlocks() {
	for block := range ebc.subscriptionChannel {
		topBlockNumber, err := strconv.ParseInt(block.Number, 0, 32)
		if err != nil {
			// FIXME Consider the right thing to do here.
			fmt.Printf("Error receiving a new block: %v", err)
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

// subscribeBlocks creates a subscription to Geth to get each block.
func (ebc *ethereumBlockCounter) subscribeBlocks() error {
	subscribeContext, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	_, err := ebc.config.clientWS.EthSubscribe(
		subscribeContext,
		ebc.subscriptionChannel,
		newHeadsSubscription,
	)
	if err != nil {
		return err
	}

	var lastBlock block
	err = ebc.config.clientRPC.Call(
		&lastBlock,
		getBlockByNumber,
		latestBlock,
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
		getBlockByNumber,
		latestBlock,
		true,
	)
	if err != nil {
		return nil,
			fmt.Errorf(
				"failed to get initial block from the chain: [%v]",
				err,
			)
	}

	startupBlockNumber, err := strconv.ParseInt(startupBlock.Number, 0, 32)
	if err != nil {
		return nil,
			fmt.Errorf(
				"failed to get initial number of blocks from the chain [%v]",
				err,
			)
	}

	blockCounter := &ethereumBlockCounter{
		latestBlockHeight:   int(startupBlockNumber),
		waiters:             make(map[int][]chan int),
		config:              ec,
		subscriptionChannel: make(chan block),
	}

	go blockCounter.receiveBlocks()
	err = blockCounter.subscribeBlocks()
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to new blocks: [%v]", err)
	}

	return blockCounter, nil
}
