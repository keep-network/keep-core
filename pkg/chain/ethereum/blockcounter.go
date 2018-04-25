package ethereum

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/keep-network/keep-core/pkg/chain"
)

type ethereumBlockCounter struct {
	structMutex sync.Mutex
	blockHeight int
	subch       chan Block
	conn        *rpc.Client
	waiters     map[int][]chan int
	Debug01     bool
}

// compile time test - at compile time to verify that the local satisfies the interface.
var _ chain.BlockCounter = (*ethereumBlockCounter)(nil)

type Block struct {
	Number string
}

// WaitForBlocks waits for a minimum of 1 block before returing.
func (blockWait *ethereumBlockCounter) WaitForBlocks(numBlocks int) {
	waiter := blockWait.BlockWaiter(numBlocks)
	<-waiter
	return
}

// BlockWaiter returns the block number as a chanel with a minimum of 1 block wait. 0 and negative numBlocks are converted to 1.
func (blockWait *ethereumBlockCounter) BlockWaiter(numBlocks int) <-chan int {
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

	return newWaiter
}

// receiveBlocks gets each new block back from Geth and extracts the block height (top) form it.
// For each block height that is being weighted on a message will be sent.
func (blockWait *ethereumBlockCounter) receiveBlocks() {

	for block := range blockWait.subch {
		if top, err := strconv.ParseInt(block.Number, 0, 64); err == nil {
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

	blockWait.structMutex.Lock()
	defer blockWait.structMutex.Unlock()

	blockSubscription, err := blockWait.conn.EthSubscribe(ctx, blockWait.subch, "newHeads")
	if err != nil {
		if blockWait.Debug01 {
			fmt.Println("Subscription Failed => ", err)
		}
		return
	}

	var lastBlock Block
	err = blockWait.conn.Call(&lastBlock, "eth_getBlockByNumber", "latest", true)
	if err != nil {
		if blockWait.Debug01 {
			fmt.Println("can't get latest block:", err)
		}
		return
	}

	blockWait.subch <- lastBlock

	if blockWait.Debug01 {
		fmt.Println("connection lost", <-blockSubscription.Err())
	}
}

// BlockCounter creates a BlockCounter that uses the block number in ethereum.
func BlockCounter(client *rpc.Client) chain.BlockCounter {
	blockWait := ethereumBlockCounter{blockHeight: 0, waiters: make(map[int][]chan int), conn: client}

	// -----------------------------------------------------------------------------------------------------------------------------
	var lastBlock Block
	err := blockWait.conn.Call(&lastBlock, "eth_getBlockByNumber", "latest", true)
	if err != nil {
		if blockWait.Debug01 {
			fmt.Println("can't get latest block:", err)
		}
		log.Fatal("Failed to get initial number of blocks from the chain")
		// return
	}

	// TODO : must set blockHeight to correct value: ---------------------- > notifyBlockHeight := blockWait.blockHeight + numBlocks
	if ii, err := strconv.ParseInt(lastBlock.Number, 0, 64); err == nil {
		// fmt.Printf("%T, %v\n", ii, ii)
		blockWait.blockHeight = int(ii)
	} else {
		log.Fatal("Failed to get initial number of blocks from the chain")
	}
	// -----------------------------------------------------------------------------------------------------------------------------

	blockWait.subch = make(chan Block)

	go func() {
		for i := 0; ; i++ {
			if i > 0 {
				time.Sleep(2 * time.Second)
			}
			blockWait.subscribeBlocks()
		}
	}()

	go blockWait.receiveBlocks()

	return &blockWait
}
