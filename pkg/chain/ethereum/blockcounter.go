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
	latestBlockHeight   uint64
	subscriptionChannel chan block
	config              *ethereumChain
	waiters             map[uint64][]chan uint64
	waitingCancelFuncs  map[uint64][]context.CancelFunc
}

type block struct {
	Number string
}

func (ebc *ethereumBlockCounter) WaitForBlockHeight(blockNumber uint64) error {
	waiter, err := ebc.BlockHeightWaiter(blockNumber)
	if err != nil {
		return err
	}
	<-waiter
	return nil
}

func (ebc *ethereumBlockCounter) BlockHeightWaiter(
	blockNumber uint64,
) (<-chan uint64, error) {
	newWaiter := make(chan uint64)

	ebc.structMutex.Lock()
	defer ebc.structMutex.Unlock()

	if blockNumber <= ebc.latestBlockHeight {
		go func() { newWaiter <- blockNumber }()
	} else {
		waiterList, exists := ebc.waiters[blockNumber]
		if !exists {
			waiterList = make([]chan uint64, 0)
		}

		ebc.waiters[blockNumber] = append(waiterList, newWaiter)
	}

	return newWaiter, nil
}

func (ebc *ethereumBlockCounter) CurrentBlock() (uint64, error) {
	return ebc.latestBlockHeight, nil
}

// receiveBlocks gets each new block back from Geth and extracts the
// block height (topBlockNumber) form it. For each block height that is being
// waited on a message will be sent.
func (ebc *ethereumBlockCounter) receiveBlocks() {
	for block := range ebc.subscriptionChannel {
		topBlockNumber, err := strconv.ParseInt(block.Number, 0, 32)
		if err != nil {
			// FIXME Consider the right thing to do here.
			fmt.Printf("Error receiving a new block: %v", err)
		}

		// receivedBlockHeight is the current blockchain height as just
		// received in the notification. latestBlockHeightSeen is the
		// blockchain height as observed in the previous invocation of
		// receiveBlocks().
		//
		// If we have already received notification about this block,
		// we do nothing. All handlers were already called for this block
		// height.
		receivedBlockHeight := uint64(topBlockNumber)
		if receivedBlockHeight == ebc.latestBlockHeight {
			continue
		}

		// We have already seen latestBlockHeightSeen during the previous
		// execution of receiveBlocks() function and all handlers for
		// latestBlockHeightSeen were called. Now we start from the next block
		// after it and that's latestBlockHeightSeen + 1.
		for unseenBlockNumber := ebc.latestBlockHeight + 1; unseenBlockNumber <= receivedBlockHeight; unseenBlockNumber++ {
			ebc.structMutex.Lock()
			height := unseenBlockNumber
			ebc.latestBlockHeight++
			waiters := ebc.waiters[height]
			delete(ebc.waiters, height)
			ebc.structMutex.Unlock()

			for _, waiter := range waiters {
				go func(w chan uint64) { w <- height }(waiter)
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

	startupBlockHeight, err := strconv.ParseInt(startupBlock.Number, 0, 32)
	if err != nil {
		return nil,
			fmt.Errorf(
				"failed to get initial number of blocks from the chain [%v]",
				err,
			)
	}

	blockCounter := &ethereumBlockCounter{
		latestBlockHeight:   uint64(startupBlockHeight),
		waiters:             make(map[uint64][]chan uint64),
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
