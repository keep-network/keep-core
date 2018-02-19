package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/keep.network/keep-core/go/interface/lib/KStart"
	"github.com/pschlump/MiscLib"
	"github.com/pschlump/godebug"
)

/*
// TODO
/Users/corwin/go/src/github.com/keep.network/keep-core/go/interface/WatchForEvent/watchForEvents.go
	1. Add in CLI
	2. 	-w Name - watch for this event
		-s Server - mac default: "/Users/corwin/Projects/eth/data/geth.ipc"		-> $HOME/... -- Normal geth ipc
				  - mac default: "ws://192.168.0.157:8546"						-> 127.0.0.1:8546			ws
		-f contractAddress - The address of the contract to call				-> for testnet: *ContractAddress := "0xe705ab560794cf4912960e5069d23ad6420acde7"
*/
var Name = flag.String("watch", "RequetRelayEvent", "Name of contract method to call")                         // 0
var GethServer = flag.String("server", "ws://127.0.0.1:8546", "Geth server to talk to")                        // 2
var ContractAddress = flag.String("conaddr", "0xe705ab560794cf4912960e5069d23ad6420acde7", "Contract Address") // 3
func init() {
	flag.StringVar(Name, "w", "RequestRelayEvent", "Name of contract method to call")                      // 0
	flag.StringVar(GethServer, "s", "ws://127.0.0.1:8546", "Geth server to talk to")                       // 2
	flag.StringVar(ContractAddress, "a", "0xe705ab560794cf4912960e5069d23ad6420acde7", "Contract Address") // 3
}

func main() {

	flag.Parse()

	if *Name == "RequestRelayEvent" || *Name == "RelayEntryEvent" ||
		*Name == "RelayResetEvent" || *Name == "SubmitGroupPublicKeyEvent" {
	} else {
		fmt.Printf("Invalid Name=%s Must be name of contract call.\n")
		os.Exit(1)
	}

	// Create RPC connection to a remote node
	//
	// Examples of connection strings - note - this must be an `.ipc` or a websocket "ws://" connection.
	// A HTTP/HTTPS connection will not work.
	//
	// GethServer := "http://10.51.245.75:8545"	// will not work - since this is "push" from server data.
	// GethServer := "http://10.51.245.75:8546" 	// WebSocket connect
	// GethServer := "/Users/corwin/Projects/eth/data/geth.ipc"
	// GethServer := "ws://192.168.0.157:8546" // Based on running a local "geth"
	conn, err := ethclient.Dial(*GethServer)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v at address: %s", err, *GethServer)
	}

	// Address of the contract - extracted from the output from "truffle migrate --reset --network testnet"
	// *ContractAddress := "0xe705ab560794cf4912960e5069d23ad6420acde7"
	fmt.Printf("Connected - v15 (watch for event: %s) - on testnet(%s)\n", *Name, *GethServer)

	fmt.Printf("----------------------------------------------------------------------------------\n")
	fmt.Printf("Event Capture\n")
	fmt.Printf("----------------------------------------------------------------------------------\n")

	// Create a object to talk to the contract.
	ks, err := KStart.NewKStart(common.HexToAddress(*ContractAddress), conn)
	if err != nil {
		log.Fatalf("Failed to instantiate contract object: %v at address: %s", err, *GethServer)
	}

	switch *Name {
	case "RequestRelayEvent":
		// event RequestRelayEvent(uint256 RequestID, uint256 Payment, uint256 BlockReward, uint256 Seed);
		// TODO/Question? do we want it to buffer? don't know (128 is buffer size in go-ethereum
		// See: /Users/corwin/go/src/github.com/ethereum/go-ethereum/accounts/abi/bind/base.go
		sink := make(chan *KStart.KStartRequestRelayEvent, 10)
		fmt.Printf("Calling Watch for %s, %s\n", *Name, godebug.LF())
		event, err := ks.WatchRequestRelayEvent(nil, sink)
		if err != nil {
			fmt.Printf("%sError creating watch for relay events: %s%s\n", MiscLib.ColorRed, err, MiscLib.ColorReset)
		} else {
			// Event is an event.Subscription - that is a chanel and the "Unsubscribe" -- So, Unsubscribe will get called automatically.
			for {
				// TODO/Question? do we need a "quit" event in this?
				select {
				case rn := <-sink:
					fmt.Printf("%sGot a relay request event! Yea!%s, %+v\n", MiscLib.ColorGreen, MiscLib.ColorReset, rn)
					fmt.Printf("%s        Decoded into JSON data!%s, %s\n", MiscLib.ColorGreen, MiscLib.ColorReset, godebug.SVarI(rn))

				case ee := <-event.Err():
					fmt.Printf("%sGot an error: %s%s\n", MiscLib.ColorYellow, ee, MiscLib.ColorReset)
					os.Exit(1)
				}
			}
		}

	case "RelayEntryEvent":
		// event RelayEntryEvent(uint256 RequestID, uint256 Signature, uint256 GroupID, uint256 PreviousEntry );
		sink := make(chan *KStart.KStartRelayEntryEvent, 10)
		fmt.Printf("Calling Watch for %s, %s\n", *Name, godebug.LF())
		event, err := ks.WatchRelayEntryEvent(nil, sink)
		if err != nil {
			fmt.Printf("%sError creating watch for relay events: %s%s\n", MiscLib.ColorRed, err, MiscLib.ColorReset)
		} else {
			for {
				select {
				case rn := <-sink:
					fmt.Printf("%sGot a relay entry event! Yea!%s, %+v\n", MiscLib.ColorGreen, MiscLib.ColorReset, rn)
					fmt.Printf("%s      Decoded into JSON data!%s, %s\n", MiscLib.ColorGreen, MiscLib.ColorReset, godebug.SVarI(rn))

				case ee := <-event.Err():
					fmt.Printf("%sGot an error: %s%s\n", MiscLib.ColorYellow, ee, MiscLib.ColorReset)
					os.Exit(1)
				}
			}
		}

	case "RelayResetEvent":
		// event RelayResetEvent(uint256 LastValidRelayEntry, uint256 LastValidRelayTxHash, uint256 LastValidRelayBlock);
		fmt.Printf("TODO 0004\n")

	case "SubmitGroupPublicKeyEvent":
		// event SubmitGroupPublicKeyEvent(uint256 _PK_G_i, uint256 _id, uint256 _activationBlockHeight);
		fmt.Printf("TODO 0005\n")

	}

}
