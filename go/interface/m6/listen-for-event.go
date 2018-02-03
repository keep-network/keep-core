package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pschlump/godebug"
)

//	"github.com/miguelmota/go-web3-example/greeter"

func main() {
	// client, err := ethclient.Dial("wss://rinkeby.infura.io/ws")
	client, err := ethclient.Dial("http://127.0.0.1:9545")

	if err != nil {
		log.Fatal(err)
	}

	// From original example
	// greeterAddress := "a7b2eb1b9fff7c9625373a6a6d180e36b552fc4c"
	// priv := "abcdbcf6bdc3a8e57f311a2b4f513c25b20e3ad4606486d7a927d8074872cefg"

	greeterAddress := "9fbda871d559710256a2502a2517b794b482db40"
	// priv := "ba3887bec7c1cd7b2a2cda2a84509e49b715b1ff815432681ca91c1764b35b79"	// contrct private key, the wrong thing to use
	// (0) 0x627306090abab3a6e1400e9345bc60c78a8bef57
	priv := "c87509a1c067bbde78beb793e6fa76530b6382a4c0241e5e4a9ec0a0f44dc0d3" // Private key of account with funds, the corret thing
	// (0) c87509a1c067bbde78beb793e6fa76530b6382a4c0241e5e4a9ec0a0f44dc0d3

	key, err := crypto.HexToECDSA(priv)

	contractAddress := common.HexToAddress(greeterAddress)
	// greeterClient, err := greeter.NewGreeter(contractAddress, client)
	greeterClient, err := NewGreeter(contractAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("AT: %s\n", godebug.LF())

	auth := bind.NewKeyedTransactor(key)

	fmt.Printf("AT: %s\n", godebug.LF())
	// not sure why I have to set this when using testrpc
	var nonce int64 = 8
	auth.Nonce = big.NewInt(nonce)
	fmt.Printf("AT: %s\n", godebug.LF())

	tx, err := greeterClient.Greet(auth, "hello")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("AT: %s\n", godebug.LF())

	fmt.Printf("Pending TX: 0x%x\n", tx.Hash())

	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	fmt.Printf("AT: %s\n", godebug.LF())

	var ch = make(chan types.Log)
	ctx := context.Background()

	fmt.Printf("AT: %s\n", godebug.LF()) // last working line with truffle, "Subscribe: notifications not supported"

	sub, err := client.SubscribeFilterLogs(ctx, query, ch)
	if err != nil {
		log.Println("Subscribe:", err)
		return
	}

	fmt.Printf("AT: %s\n", godebug.LF())

	for {
		fmt.Printf("AT: %s\n", godebug.LF())
		select {
		case err := <-sub.Err():
			fmt.Printf("AT: %s\n", godebug.LF())
			log.Fatal(err)
		case log := <-ch:
			fmt.Printf("AT: %s\n", godebug.LF())
			fmt.Println("Log:", log)
		}
		fmt.Printf("AT: %s\n", godebug.LF())
	}

	fmt.Printf("AT: %s\n", godebug.LF())

}
