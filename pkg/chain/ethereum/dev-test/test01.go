// +build ignore

package main

// favorite
// I have followed this guide and wrote the same code but instead of getting the latest block, I am getting this,

// Subscription Failed =>  The method newBlocks_newBlocks does not exist/is not available
// Code :

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
)

func main() {
	clientSubscription()
}

type Block struct {
	Number string
}

func clientSubscription() {
	//client, err := rpc.Dial("ws://115.186.176.139:8547")
	client, err := rpc.Dial("ws://192.168.0.158:8546")
	if err != nil {
		fmt.Println("Error Connecting to Server", err)

	}
	// fmt.Printf("type= %T\n", client)
	// os.Exit(0)

	subch := make(chan Block)
	go func() {
		for i := 0; ; i++ {
			if i > 0 {
				time.Sleep(2 * time.Second)
			}
			subscribeBlocks(client, subch)
		}
	}()

	for block := range subch {
		if s, err := strconv.ParseInt(block.Number, 0, 32); err == nil {
			// fmt.Printf("%T, %v\n", s, s)
			fmt.Printf("Latest Block => %s, in decimal=%d\n", block.Number, s)
		}
	}

}

func subscribeBlocks(client *rpc.Client, subch chan Block) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// From: https://github.com/ethereum/go-ethereum/wiki/RPC-PUB-SUB
	blockSubscription, err := client.EthSubscribe(ctx, subch, "newHeads")
	if err != nil {
		fmt.Println("Subscription Failed => ", err)
		return
	}

	var lastBlock Block
	err = client.Call(&lastBlock, "eth_getBlockByNumber", "latest", true)
	if err != nil {
		fmt.Println("can't get latest block:", err)
		return
	}

	subch <- lastBlock

	fmt.Println("connection lost", <-blockSubscription.Err())
}

/*
	and for the latest block:
	web3.eth.getBlock(web3.eth.blockNumber).hash
*/
