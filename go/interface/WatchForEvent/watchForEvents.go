package main

import (
	"fmt"
	"log"

	// "www.2c-why.com/job/ethrpc/ethABI/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/keep.network/keep-core/go/interface/lib/KStart" // "github.com/keep.network/keep-core/go/interface/m4/KStart"
	"github.com/pschlump/MiscLib"
	"github.com/pschlump/godebug"
)

// "github.com/ethereum/go-ethereum/ethclient"

func main() {

	// Create RPC connection to a remote node
	// addr := "http://10.51.245.75:8545"
	// addr := "http://10.51.245.75:8546" // WebSocket connect
	// addr := "/Users/corwin/Projects/eth/data/geth.ipc"
	addr := "ws://192.168.0.157:8546"
	conn, err := ethclient.Dial(addr)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v at address: %s", err, addr)
	}

	caddr := "0xe705ab560794cf4912960e5069d23ad6420acde7"
	fmt.Printf("Connected - v14 (watch for event) - test on testnet(%s)\n", addr)

	// See ../m7 for event generation

	fmt.Printf("----------------------------------------------------------------------------------\n")
	fmt.Printf("Event Capture\n")
	fmt.Printf("----------------------------------------------------------------------------------\n")
	fmt.Printf("\n")

	// ks, err := NewKStart(common.HexToAddress(caddr), backend bind.ContractBackend)
	ks, err := KStart.NewKStart(common.HexToAddress(caddr), conn)

	//	var opts2 bind.WatchOpts
	//	var x uint64
	//	x = 20
	//	opts2.Start = &x

	sink := make(chan *KStart.KStartRequestRelayEvent, 10) // TODO/Question? do we want it to buffer? don't know (128 is buffer size in go-ethereum
	// /Users/corwin/go/src/github.com/ethereum/go-ethereum/accounts/abi/bind/base.go

	fmt.Printf("Calling Watch, %s\n", godebug.LF())

	event, err := ks.WatchRequestRelayEvent(nil, sink)
	// /Users/corwin/go/src/github.com/ethereum/go-ethereum/accounts/abi/bind/base.go:298
	// 		func (c *BoundContract) WatchLogs(opts *WatchOpts, name string, query ...[]interface{}) (chan types.Log, event.Subscription, error) {
	if err != nil {
		fmt.Printf("%sError creating watch for relay events: %s%s\n", MiscLib.ColorRed, err, MiscLib.ColorReset)
	} else {
		// Event is an event.Subscription - that is a chanel and the "Unsubscribe" -- So, Unsubscribe will get called,

		for {
			// TODO/Question? do we need a "quit" event in this?
			select {
			case rn := <-sink:
				fmt.Printf("%sGot a relay request event! Yea!%s, %+v\n", MiscLib.ColorGreen, MiscLib.ColorReset, rn)

			case ee := <-event.Err():
				fmt.Printf("%sGot an error: %s%s\n", MiscLib.ColorYellow, ee, MiscLib.ColorReset)
				break
			}
		}

	}

}
