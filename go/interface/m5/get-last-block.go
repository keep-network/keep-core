package main

// From: https://ethereum.stackexchange.com/questions/13341/ethereum-go-how-to-get-the-latest-block

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/rpc"
)

type Block struct {
	Number string
}

func main() {
	// Connect the client
	// client, err := rpc.Dial("http://127.0.0.1:9545")
	client, err := rpc.Dial("http://192.168.0.139:8545")
	if err != nil {
		log.Fatalf("could not create ipc client: %v", err)
	}

	var lastBlock Block
	err = client.Call(&lastBlock, "eth_getBlockByNumber", "latest", true)
	if err != nil {
		fmt.Println("can't get latest block:", err)
		return
	}

	// Print events from the subscription as they arrive.
	fmt.Printf("latest block: %v\n", lastBlock.Number)
}
