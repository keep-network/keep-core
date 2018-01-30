package main

import (
	"fmt"
	"log"

	// "github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {

	// Create an IPC based RPC connection to a remote node
	// addr := "http://192.168.0.139:8545"
	addr := "http://127.0.0.1:9545"
	conn, err := ethclient.Dial(addr)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v at address: %s", err, addr)
	}

	fmt.Printf("Connected - V8\n")
	// Instantiate the contract and display its name
	// token, err := common.NewToken(common.HexToAddress("0x21e6fc92f93c8a1bb41e2be64b4e1f88a54d3576"), conn)
	// token, err := NewTokenRecipientType(common.HexToAddress("0x627306090abab3a6e1400e9345bc60c78a8bef57"), conn)
	// token, err := NewTokenType(common.HexToAddress("0xaa588d3737b611bafd7bd713445b314bd453a5c8"), conn)
	// token, err := NewTokenType(common.HexToAddress("0x75c35c980c0d37ef46df04d31a140b65503c0eed"), conn)
	token, err := NewTokenType(common.HexToAddress("0xaa588d3737b611bafd7bd713445b314bd453a5c8"), conn)
	if err != nil {
		log.Fatalf("Failed to instantiate a Token contract: %v", err)
	}

	// Access a public variable
	fmt.Printf("NewTokenRecipientType Called\n")
	name, err := token.Name(nil) // value of this.name in contract public variable
	if err != nil {
		log.Fatalf("Failed to retrieve token name: %v", err)
	}
	fmt.Printf("Token name: %s\n", name)

	// Call a constant/pure function

	// Create a tranaction to call a function

}
