package EthBlockCounter

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/keep-network/keep-core/go/BeaconConfig"
	"github.com/keep-network/keep-core/go/beacon/chain"
	"github.com/pschlump/MiscLib"
	"github.com/pschlump/godebug"
)

// "github.com/keep-network/keep-core/go/BeaconConfig"

type EthBlockCounter struct {
	structMutex sync.Mutex         //
	blockHeight int                //
	waiters     map[int][]chan int //
	gethServer  string             //
	timeout     int                // can be set to 0 to disable the timeout - if you do this and run on truffle you may never get a block count
	client      *rpc.Client        //
	debugFlag   bool               //
	ticker      *time.Ticker       //
	ticker2     *time.Ticker       //
}

// At compile time verify that this meats the interface requirement
var _ chain.BlockCounter = (*EthBlockCounter)(nil)

func (eth *EthBlockCounter) WaitForBlocks(numBlocks int) {
	waiter := eth.BlockWaiter(numBlocks)
	<-waiter
	return
}

func (eth *EthBlockCounter) BlockWaiter(numBlocks int) <-chan int {
	newWaiter := make(chan int)

	eth.structMutex.Lock()
	defer eth.structMutex.Unlock()
	notifyBlockHeight := eth.blockHeight + numBlocks
	if eth.debugFlag {
		fmt.Printf("%sBlockWaiter called, notifyBlockHeight=%d AT:%s%s\n", MiscLib.ColorYellow, notifyBlockHeight, godebug.LF(), MiscLib.ColorReset)
	}

	if notifyBlockHeight <= eth.blockHeight {
		newWaiter <- notifyBlockHeight
	} else {
		waiterList, exists := eth.waiters[notifyBlockHeight]
		if !exists {
			waiterList = make([]chan int, 0)
		}

		eth.waiters[notifyBlockHeight] = append(waiterList, newWaiter)
	}

	return newWaiter
}

func (eth *EthBlockCounter) count() {

	var curNumber int64
	nSec := 0
	init2 := true
	var e0, err error
	lastHeightSeen := eth.blockHeight

	run1 := func(top int) {
		if eth.debugFlag {
			fmt.Printf("AT:%s\n", godebug.LF())
		}
		// This loop is because occasionally the block height skips a number.  With the
		// data structure this could cause a skipped chanel write - that would be bad.
		for height := top; height <= lastHeightSeen; height++ {
			if eth.debugFlag {
				fmt.Printf("run1 called(), height=%d AT:%s\n", height, godebug.LF())
			}

			eth.structMutex.Lock()
			waiters, exists := eth.waiters[height]
			if exists {
				delete(eth.waiters, height)
			}
			eth.structMutex.Unlock()

			if exists {
				// fmt.Printf("AT:%s\n", godebug.LF())
				for _, waiter := range waiters {
					// fmt.Printf("AT:%s\n", godebug.LF())
					go func(w chan int) { w <- height }(waiter)
				}
			}
		}
		lastHeightSeen = top
	}

	eth.ticker = time.NewTicker(time.Millisecond * 1000)  // check 1ce a second for a new block
	eth.ticker2 = time.NewTicker(time.Millisecond * 1000) // count seconds for a simulated new block
	go func() {
		for t := range eth.ticker.C {
			if eth.debugFlag {
				fmt.Printf("Check For Block Number Now %v\n", t)
			}
			curNumber, e0 = eth.getBlockNumber()
			if err != nil {
				err = e0
				return
			}
			run1(int(curNumber))
		}
	}()
	if err != nil {
		if eth.debugFlag {
			fmt.Printf("Error from go routine %s\n", err)
		}
	}

	if eth.timeout > 0 {
		// Setup a "guard" if more than N seconds have passed w/o a channel, then send a
		// message on channel
		go func() {
			for t := range eth.ticker2.C {
				_ = t
				if init2 {
					init2 = false
					nSec = 0
				} else {
					nSec++
					if nSec > eth.timeout {
						nSec = 0
						run1(int(curNumber))
					}
				}
			}
		}()
	}

}

func (eth *EthBlockCounter) getBlockNumber() (n int64, err error) {

	type Block struct {
		Number string `json:"number"`
	}

	var lastBlock Block

	err = eth.client.Call(&lastBlock, "eth_getBlockByNumber", "latest", true)
	if err != nil {
		fmt.Println("can't get latest block:", err)
		return
	}

	curNumber, _ := strconv.ParseInt(lastBlock.Number[2:], 16, 64)
	if eth.debugFlag {
		fmt.Printf("%slatest block: %d %s raw: ->%s<-\n", MiscLib.ColorGreen, curNumber, MiscLib.ColorReset, lastBlock.Number)
	}
	return curNumber, nil
}

// EthBlockCounter creates a BlockCounter that pols the Ethereum chain. It hadles
// weird stuff like missing block numbers.   You should specify a timeout in
// seconds as a backup so if this is run on "testrpc"/"truffle" where no
// mining occures it will still timeout.   Call BeaconConfig.NewEthConnection,
// or BeaconConfig.OpenConnection to get the "client".
func NewEthBlockCounter(client *rpc.Client, cfg BeaconConfig.BeaconConfig) chain.BlockCounter {

	// OLD: func NewEthBlockCounter(GethServer string, timeout int) chain.BlockCounter {

	eth := EthBlockCounter{
		blockHeight: 0,
		waiters:     make(map[int][]chan int),
		gethServer:  cfg.GethServer,
		timeout:     cfg.BlockTimeout,
		debugFlag:   false,
		client:      client,
	}

	curNumber, err := eth.getBlockNumber()
	if err != nil {
		log.Fatalf("could not setup initial block number for server %s client: %v", cfg.GethServer, err)
		return &eth
	}
	eth.blockHeight = int(curNumber)

	go eth.count()

	return &eth
}
