package ethereum

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
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
	waiters             map[uint64][]chan uint64
	watchers            []*watcher
}

type block struct {
	Number string
}

type watcher struct {
	ctx     context.Context
	channel chan uint64
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

func (ebc *ethereumBlockCounter) WatchBlocks(ctx context.Context) <-chan uint64 {
	watcher := &watcher{
		ctx:     ctx,
		channel: make(chan uint64),
	}

	ebc.structMutex.Lock()
	ebc.watchers = append(ebc.watchers, watcher)
	ebc.structMutex.Unlock()

	go func() {
		<-ctx.Done()

		ebc.structMutex.Lock()
		for i, w := range ebc.watchers {
			if w == watcher {
				ebc.watchers[i] = ebc.watchers[len(ebc.watchers)-1]
				ebc.watchers = ebc.watchers[:len(ebc.watchers)-1]
				break
			}
		}
		ebc.structMutex.Unlock()
	}()

	return watcher.channel
}

// receiveBlocks gets each new block back from Geth and extracts the
// block height (topBlockNumber) form it. For each block height that is being
// waited on a message will be sent.
func (ebc *ethereumBlockCounter) receiveBlocks() {
	for block := range ebc.subscriptionChannel {
		topBlockNumber, err := strconv.ParseInt(block.Number, 0, 32)
		if err != nil {
			// FIXME Consider the right thing to do here.
			logger.Errorf("error receiving a new block: [%v]", err)
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

			ebc.structMutex.Lock()
			watchers := make([]*watcher, len(ebc.watchers))
			copy(watchers, ebc.watchers)
			ebc.structMutex.Unlock()

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
}

// subscribeBlocks creates a subscription to Geth to get each block.
func (ebc *ethereumBlockCounter) subscribeBlocks(
	clientWS *rpc.Client,
	clientRPC *rpc.Client,
) error {
	subscribeContext, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	_, err := clientWS.EthSubscribe(
		subscribeContext,
		ebc.subscriptionChannel,
		newHeadsSubscription,
	)
	if err != nil {
		return err
	}

	var lastBlock block
	err = clientRPC.Call(
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
	return ec.blockCounter, nil
}

func createBlockCounter(clientWS *rpc.Client, clientRPC *rpc.Client) (*ethereumBlockCounter, error) {
	var startupBlock block
	err := clientRPC.Call(
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
		subscriptionChannel: make(chan block),
	}

	go blockCounter.receiveBlocks()
	err = blockCounter.subscribeBlocks(clientWS, clientRPC)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to new blocks: [%v]", err)
	}

	return blockCounter, nil
}
