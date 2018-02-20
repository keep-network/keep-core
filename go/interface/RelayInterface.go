package RelayContract

// ++ xyzzy - Calls to contract functions -- finish											1hr
// ++ xyzzy - Event captures -- finish														1hr

// ++ xyzzy - Run auto tests (truffle test --network testnet)								1hr
// ++ xyzzy - Run test of events with ./* files.											1hr

// ++ xyzzy - Create a set of tests for this												2hrs

// ++ xyzzy - go lint																		15min
// ++ xyzzy - go doc																		1hr

// ++ xyzzy - Add in "logging"

import (
	"bufio"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/keep.network/keep-core/go/interface/lib/KeepRelayBeacon"
	"github.com/pschlump/godebug"
)

type KeepRelayBeaconContract struct {
	conn            *ethclient.Client
	kstart          *KeepRelayBeacon.KeepRelayBeaconTransactor
	kstartCaller    *KeepRelayBeacon.KeepRelayBeaconCaller
	opts            *bind.TransactOpts
	ContractAddress string
	GethServer      string
}

func NewKeepRelayBeaconContract(ctx *RelayContractContext, GethServer, ContractAddress, KeyFile, KeyFilePassword string) (ri *KeepRelayBeaconContract, err error) {

	ri = &KeepRelayBeaconContract{
		GethServer:      GethServer,
		ContractAddress: ContractAddress,
	}

	// Create JSON RPC connection to a remote node
	//
	// Examples of connection strings - you can use .ipc, ws: or http:
	// connections for this.
	//
	// GethServer := "http://10.51.245.75:8545"
	// GethServer := "http://10.51.245.75:8546" // WebSocket connect
	// GethServer := "/Users/corwin/Projects/eth/data/geth.ipc" // IPC local connection
	conn, err := ethclient.Dial(GethServer)
	if err != nil {
		err = fmt.Errorf("Failed to connect to the Ethereum client: %s at address: %s", err, GethServer)
		return
	}

	ri.conn = conn

	if ctx.dbOn() {
		fmt.Printf("Connected - v16 - on net(%s)\n", GethServer)
	}

	// Address of the contract - extracted from the output from "truffle migrate --reset --network testnet"
	// ContractAddress := "0xe705ab560794cf4912960e5069d23ad6420acde7"
	hexAddr := common.HexToAddress(ContractAddress)

	// Instantiate the contract and call a method
	kstart, err := KeepRelayBeacon.NewKeepRelayBeaconTransactor(hexAddr, conn)
	if err != nil {
		err = fmt.Errorf("Failed to instantiate a KeepRelayBeaconTranactor contract: %s", err)
		// LOG at this point
		return
	}

	kstartCaller, err := KeepRelayBeacon.NewKeepRelayBeaconCaller(hexAddr, conn)
	if err != nil {
		err = fmt.Errorf("Failed to instantiate a KeepRelayBeaconCaller contract: %s", err)
		// LOG at this point
		return
	}

	ri.kstart = kstart
	ri.kstartCaller = kstartCaller

	if ctx.dbOn() {
		fmt.Printf("Connect to KeepRelayBeacon contract successful\n")
	}

	// Account with "gas" to spend so that you can call the contract.
	// KeyFile := "./UTC--2018-02-15T19-57-35.216297214Z--6ffba2d0f4c8fd7961f516af43c55fe2d56f6044"
	// KeyFilePassword := "password"
	file, err := os.Open(KeyFile) // Read in file
	if err != nil {
		err = fmt.Errorf("Failed to open keyfile: %s, %s", err, KeyFile)
		return
	}

	// Create a new trasaction call (will mine) for the contract.
	opts, err := bind.NewTransactor(bufio.NewReader(file), KeyFilePassword)
	if err != nil {
		err = fmt.Errorf("Failed to read keyfile: %s, %s", err, KeyFile)
		return
	}

	ri.opts = opts

	return
}

//
//	var tx *types.Transaction
//	switch *Name {
//	case "RequestRelay":
//		// function requestRelay(uint256 _payment, uint256 _blockReward, uint256 _seed) public returns ( uint256 RequestID ) {
//		rawseed := []byte("abcdef") // TODO - variable data - data for this call
//		payment := big.NewInt(12)
//		blockReward := big.NewInt(12)
//		seed := big.NewInt(0).SetBytes(rawseed)
//		tx, err = kstart.RequestRelay(opts, payment, blockReward, seed)
//
//	case "IsStaked":
//		// function isStaked(uint256 _UserPublicKey) pure public returns(bool) {
//		var isStakedNow bool
//		userPublicKey := big.NewInt(12)
//		isStakedNow, err = kstartCaller.IsStaked(nil, userPublicKey)
//
//		if err == nil {
//			fmt.Printf("----------------------------------------------------------------------------------\n")
//			fmt.Printf("%sKeepRelayBeacon.%s called%s Results Are=%v\n", MiscLib.ColorGreen, Name, isStakedNow, MiscLib.ColorReset, tx)
//			fmt.Printf("----------------------------------------------------------------------------------\n")
//		}
//
//	case "RelayEntry":
//		// function relayEntry(uint256 _RequestID, uint256 _groupSignature, uint256 _groupID, uint256 _previousEntry) public {
//		requestID := big.NewInt(12)
//		groupSignature := big.NewInt(12)
//		groupID := big.NewInt(12)
//		previousEntry := big.NewInt(12)
//		tx, err = kstart.RelayEntry(opts, requestID, groupSignature, groupID, previousEntry)
//
//	case "RelayEntryAcustation":
//		// function relayEntryAcustation( uint256 _LastValidRelayTxHash, uint256 _LastValidRelayBlock) public {
//		fmt.Printf("TODO 0001\n")
//
//	case "SubmitGroupPublicKey":
//		// function submitGroupPublicKey (uint256 _PK_G_i, uint256 _id) public {
//		fmt.Printf("TODO 0002\n")
//	}
//
//	if err != nil {
//		fmt.Fprintf(os.Stderr, "----------------------------------------------------------------------------------\n")
//		fmt.Fprintf(os.Stderr, "%sFailed to call KeepRelayBeacon.%s: %s%s\n", MiscLib.ColorRed, Name, err, MiscLib.ColorReset)
//		fmt.Fprintf(os.Stderr, "----------------------------------------------------------------------------------\n")
//		os.Exit(1)
//	}
//
//	if tx != nil {
//		fmt.Printf("----------------------------------------------------------------------------------\n")
//		fmt.Printf("%sKeepRelayBeacon.%s called%s tx=%s\n", MiscLib.ColorGreen, Name, MiscLib.ColorReset, tx)
//		fmt.Printf("----------------------------------------------------------------------------------\n")
//	}
//
//}

func (ri *KeepRelayBeaconContract) RequestRelay(ctx *RelayContractContext, payment, blockReward int64, seed []byte) (tx *types.Transaction, err error) {
	pPayment := big.NewInt(payment)
	pBlockReward := big.NewInt(blockReward)
	pSeed := big.NewInt(0).SetBytes(seed)
	tx, err = ri.kstart.RequestRelay(ri.opts, pPayment, pBlockReward, pSeed)
	if err != nil {
		err = fmt.Errorf("Error: %s on call to (geth=%s) KeepRelayBeacon.at(0x%s).RequestRelay(%d,%d,%x)\n", err, ri.GethServer, ri.ContractAddress, payment, blockReward, seed)
		// LOG at this point
		return
	}
	if ctx.dbOn() {
		fmt.Printf("KeepRelayBeacon.RequestRelay called tx=%s from=%s\n", godebug.SVar(tx), godebug.LF())
		// TODO: log the success and tx that it is in
	}
	return
}

// TODO
//  	function isStaked(uint256 _UserPublicKey) pure public returns(bool) {
//		function getRandomNumber( uint256 _RequestID ) view public returns ( uint256 theRandomNumber ) {
//  	function relayEntryAcustation( uint256 _LastValidRelayTxHash, uint256 _LastValidRelayBlock) public {
//  	function submitGroupPublicKey(uint256 _PK_G_i, uint256 _id) public {
//		function getRandomNumber( uint256 _RequestID ) view public returns ( uint256 theRandomNumber ) {
