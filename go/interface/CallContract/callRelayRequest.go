package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/keep-network/keep-core/go/interface/lib/KeepRelayBeacon"
	"github.com/pschlump/MiscLib"
)

// "github.com/ethereum/go-ethereum/ethclient"

/*
/Users/corwin/go/src/github.com/keep-network/keep-core/go/interface/CallContract
	1. Add in CLI
	2. 	-c Name - call this item method in the contract -
		-d Data - with this set of data -
		-s GethServer - mac default: "/Users/corwin/Projects/eth/data/geth.ipc"
			-> $HOME/... -- Normal geth ipc
			  - mac default: "ws://192.168.0.157:8546"
			-> 127.0.0.1:8546			ws
		-a contractAddress - The address of the contract to call
			-> for testnet: ContractAddress := "0x639deb0dd975af8e4cc91fe9053a37e4faf37649"
		-k KeyFile
			-> keyFile :=
			"./UTC--2018-02-15T19-57-35.216297214Z--6ffba2d0f4c8fd7961f516af43c55fe2d56f6044"
		-p Password
			-> KeyFilePassword := "password"
*/
var Name = flag.String("name", "RequetRelay", "Name of contract method to call")                               // 0
var Data = flag.String("data", "./RequestRelay-data.json", "Contract parameters (data)")                       // 1
var GethServer = flag.String("server", "ws://127.0.0.1:8546", "Geth server to talk to")                        // 2
var ContractAddress = flag.String("conaddr", "0x639deb0dd975af8e4cc91fe9053a37e4faf37649", "Contract Address") // 3
var KeyFile = flag.String("keyfile",
	"./UTC--2018-02-15T19-57-35.216297214Z--6ffba2d0f4c8fd7961f516af43c55fe2d56f6044",
	"Account Private Key File") // 4

var KeyFilePassword = flag.String("password", "password", "Private key encryption password") // 5
func init() {
	flag.StringVar(Name, "n", "RequestRelay", "Name of contract method to call")                           // 0
	flag.StringVar(Data, "d", "./RequestRelay-data.json", "Contract parameters (data)")                    // 1
	flag.StringVar(GethServer, "s", "ws://127.0.0.1:8546", "Geth server to talk to")                       // 2
	flag.StringVar(ContractAddress, "a", "0x639deb0dd975af8e4cc91fe9053a37e4faf37649", "Contract Address") // 3
	flag.StringVar(KeyFile, "k",
		"./UTC--2018-02-15T19-57-35.216297214Z--6ffba2d0f4c8fd7961f516af43c55fe2d56f6044",
		"Account Private Key File") // 4
	flag.StringVar(KeyFilePassword, "p", "password", "Private key encryption password") // 5
}

func main() {

	flag.Parse()

	if *Name == "RequestRelay" || *Name == "IsStaked" || *Name == "RelayEntry" ||
		*Name == "RelayEntryAcustation" || *Name == "SubmitGroupPublicKey" {
	} else {
		fmt.Printf("Invalid Name=%s Must be name of contract call.\n")
		os.Exit(1)
	}

	// Create RPC connection to a remote node
	//
	// Examples of connection strings - you can use .ipc, ws: or http:
	// connections for this.
	//
	// GethServer := "http://10.51.245.75:8545"
	// GethServer := "http://10.51.245.75:8546" // WebSocket connect
	// GethServer := "/Users/corwin/Projects/eth/data/geth.ipc" // IPC local connection
	conn, err := ethclient.Dial(*GethServer)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v at address: %s", err, *GethServer)
	}

	// Address of the contract - extracted from the output from "truffle migrate --reset --network testnet"
	// ContractAddress := "0x639deb0dd975af8e4cc91fe9053a37e4faf37649"
	fmt.Printf("Connected - v15 %s(Call %s to create event)%s - test on testnet(%s)\n",
		MiscLib.ColorCyan, Name, MiscLib.ColorReset, *GethServer)

	// Instantiate the contract and call a method
	kstart, err := KeepRelayBeacon.NewKeepRelayBeaconTransactor(common.HexToAddress(*ContractAddress), conn)
	if err != nil {
		log.Fatalf("Failed to instantiate a KeepRelayBeaconTranactor contract: %v", err)
	}

	kstartCaller, err := KeepRelayBeacon.NewKeepRelayBeaconCaller(common.HexToAddress(*ContractAddress), conn)
	if err != nil {
		log.Fatalf("Failed to instantiate a KeepRelayBeaconCaller contract: %v", err)
	}

	fmt.Printf("----------------------------------------------------------------------------------\n")
	fmt.Printf("Connect to KeepRelayBeacon contract successful\n")
	fmt.Printf("----------------------------------------------------------------------------------\n")

	// Account with "gas" to spend so that you can call the contract.
	// KeyFile := "./UTC--2018-02-15T19-57-35.216297214Z--6ffba2d0f4c8fd7961f516af43c55fe2d56f6044"
	// KeyFilePassword := "password"
	file, err := os.Open(*KeyFile) // Read in file
	if err != nil {
		log.Fatalf("Failed to open keyfile: %v, %s", err, *KeyFile)
	}

	// Create a new trasaction call (will mine) for the contract.
	opts, err := bind.NewTransactor(bufio.NewReader(file), *KeyFilePassword)
	if err != nil {
		log.Fatalf("Failed to read keyfile: %v, %s", err, *KeyFile)
	}

	var tx *types.Transaction
	switch *Name {
	case "RequestRelay":
		// function requestRelay(uint256 _payment, uint256 _blockReward, uint256 _seed)
		// public returns ( uint256 RequestID ) {
		rawseed := []byte("abcdef") // TODO - variable data - data for this call
		payment := big.NewInt(12)
		blockReward := big.NewInt(12)
		seed := big.NewInt(0).SetBytes(rawseed)
		tx, err = kstart.RequestRelay(opts, payment, blockReward, seed)

	case "IsStaked":
		// function isStaked(uint256 _UserPublicKey) pure public returns(bool) {
		var isStakedNow bool
		userPublicKey := big.NewInt(12)
		isStakedNow, err = kstartCaller.IsStaked(nil, userPublicKey)

		if err == nil {
			fmt.Printf("----------------------------------------------------------------------------------\n")
			fmt.Printf("%sKeepRelayBeacon.%s called%s Results Are=%v\n", MiscLib.ColorGreen, Name,
				isStakedNow, MiscLib.ColorReset, tx)
			fmt.Printf("----------------------------------------------------------------------------------\n")
		}

	case "RelayEntry":
		// function relayEntry(uint256 _RequestID, uint256 _groupSignature, uint256 _groupID,
		// uint256 _previousEntry) public {
		requestID := big.NewInt(12)
		groupSignature := big.NewInt(12)
		groupID := big.NewInt(12)
		previousEntry := big.NewInt(12)
		tx, err = kstart.RelayEntry(opts, requestID, groupSignature, groupID, previousEntry)

	case "RelayEntryAcustation":
		// function relayEntryAcustation( uint256 _LastValidRelayTxHash, uint256 _LastValidRelayBlock) public {
		fmt.Printf("TODO 0001\n")

	case "SubmitGroupPublicKey":
		// function submitGroupPublicKey (uint256 _PK_G_i, uint256 _id) public {
		fmt.Printf("TODO 0002\n")
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "----------------------------------------------------------------------------------\n")
		fmt.Fprintf(os.Stderr, "%sFailed to call KeepRelayBeacon.%s: %s%s\n", MiscLib.ColorRed, Name, err,
			MiscLib.ColorReset)
		fmt.Fprintf(os.Stderr, "----------------------------------------------------------------------------------\n")
		os.Exit(1)
	}

	if tx != nil {
		fmt.Printf("----------------------------------------------------------------------------------\n")
		fmt.Printf("%sKeepRelayBeacon.%s called%s tx=%s\n", MiscLib.ColorGreen, Name, MiscLib.ColorReset, tx)
		fmt.Printf("----------------------------------------------------------------------------------\n")
	}

}

/*


 */
