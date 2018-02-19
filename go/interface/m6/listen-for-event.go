package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/keep.network/keep-core/go/interface/m4/KStart" // "github.com/keep.network/keep-core/go/interface/m4/KStart"
	"github.com/pschlump/godebug"
)

//	"github.com/miguelmota/go-web3-example/greeter"

func main() {
	// conn, err := ethclient.Dial("wss://rinkeby.infura.io/ws")
	addr := "ws://192.168.0.157:8546"
	conn, err := ethclient.Dial(addr)

	if err != nil {
		log.Fatal(err)
	}

	// From original example
	// greeterAddress := "a7b2eb1b9fff7c9625373a6a6d180e36b552fc4c"
	// priv := "abcdbcf6bdc3a8e57f311a2b4f513c25b20e3ad4606486d7a927d8074872cefg"

	greeterAddress := "0xe705ab560794cf4912960e5069d23ad6420acde7" // contract address (not greeter)
	// priv := "ba3887bec7c1cd7b2a2cda2a84509e49b715b1ff815432681ca91c1764b35b79"	// contrct private key, the wrong thing to use
	// (0) 0x627306090abab3a6e1400e9345bc60c78a8bef57

	//0	priv := "c87509a1c067bbde78beb793e6fa76530b6382a4c0241e5e4a9ec0a0f44dc0d3" // Private key of account with funds, the corret thing // xyzzy!
	//0	// (0) c87509a1c067bbde78beb793e6fa76530b6382a4c0241e5e4a9ec0a0f44dc0d3
	//0	key, err := crypto.HexToECDSA(priv)

	contractAddress := common.HexToAddress(greeterAddress)
	// greeterClient, err := greeter.NewGreeter(contractAddress, conn)
	//0 greeterClient, err := NewGreeter(contractAddress, conn)
	//0 if err != nil {
	//0 	log.Fatal(err)
	//0 }

	fmt.Printf("AT: %s\n", godebug.LF())

	//0	{
	//		auth := bind.NewKeyedTransactor(key)
	//
	//		fmt.Printf("AT: %s\n", godebug.LF())
	//		// not sure why I have to set this when using testrpc
	//		var nonce int64 = 8
	//		auth.Nonce = big.NewInt(nonce)
	//		fmt.Printf("AT: %s\n", godebug.LF())
	//
	//		tx, err := greeterClient.Greet(auth, "hello")
	//		if err != nil {
	//			log.Fatal(err)
	//		}
	//
	//		fmt.Printf("AT: %s\n", godebug.LF())
	//
	//		fmt.Printf("Pending TX: 0x%x\n", tx.Hash())
	//0	}

	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	fmt.Printf("AT: %s\n", godebug.LF())

	var ch = make(chan types.Log)
	ctx := context.Background()

	fmt.Printf("AT: %s\n", godebug.LF()) // last working line with truffle, "Subscribe: notifications not supported"

	sub, err := conn.SubscribeFilterLogs(ctx, query, ch)
	if err != nil {
		log.Println("Subscribe:", err)
		return
	}
	defer sub.Unsubscribe()

	fmt.Printf("AT: %s\n", godebug.LF())
	contract, err := bindKStart(contractAddress, conn, conn, conn)

	for {
		fmt.Printf("AT: %s\n", godebug.LF())
		select {
		case err := <-sub.Err():
			fmt.Printf("AT: %s\n", godebug.LF())
			log.Fatal(err)
		case logData := <-ch:
			fmt.Printf("AT: %s\n", godebug.LF())
			// fmt.Printf("Log: %+v %T\n", logData, logData)
			// --------------------------------------------------------------------------------------
			data := new(KStart.KStartRequestRelayEvent)
			if err := contract.UnpackLog(data, "RequestRelayEvent", logData); err != nil {
				fmt.Printf("Error: unpacking event: %s\n", err)
			} else {
				fmt.Printf("\tEvent: %s\n", godebug.SVarI(data))
			}
			// --------------------------------------------------------------------------------------
		}
		fmt.Printf("AT: %s\n", godebug.LF())
	}

	fmt.Printf("AT: %s\n", godebug.LF())

}

// bindKStart binds a generic wrapper to an already deployed contract.
func bindKStart(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(KStart.KStartABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}
