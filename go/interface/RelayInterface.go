package RelayContract

// ++ xyzzy - Calls to contract funtions -- finish
// ++ xyzzy - Event captures -- finish

import (
	"bufio"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/event"
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

	ri = &KeepRelayBeaconContract{}

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
	ri.GethServer = GethServer

	if ctx.dbOn() {
		// Address of the contract - extracted from the output from "truffle migrate --reset --network testnet"
		// ContractAddress := "0xe705ab560794cf4912960e5069d23ad6420acde7"
		fmt.Printf("Connected - v16 - on net(%s)\n", GethServer)
	}

	hexAddr := common.HexToAddress(ContractAddress)
	ri.ContractAddress = ContractAddress

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

type KeepRelayBeaconEvents struct {
	ks              *KeepRelayBeacon.KeepRelayBeacon
	conn            *ethclient.Client
	ContractAddress string
	GethServer      string
}

func NewKeepRelayBeaconEvents(ctx *RelayContractContext, GethServer, ContractAddress string) (ev *KeepRelayBeaconEvents, err error) {

	ev = &KeepRelayBeaconEvents{}

	// Create RPC connection to a remote node
	//
	// Examples of connection strings - note - this must be an `.ipc` or a websocket "ws://" connection.
	// A HTTP/HTTPS connection will not work.
	//
	// GethServer := "http://10.51.245.75:8545"	// will not work - since this is "push" from server data.
	// GethServer := "http://10.51.245.75:8546" 	// WebSocket connect
	// GethServer := "/Users/corwin/Projects/eth/data/geth.ipc"
	// GethServer := "ws://192.168.0.157:8546" // Based on running a local "geth"
	conn, err := ethclient.Dial(GethServer)
	if err != nil {
		err = fmt.Errorf("Failed to connect to the Ethereum client: %s at address: %s", err, GethServer)
		return
	}

	ev.conn = conn
	ev.GethServer = GethServer

	if ctx.dbOn() {
		// Address of the contract - extracted from the output from "truffle migrate --reset --network testnet"
		// ContractAddress := "0xe705ab560794cf4912960e5069d23ad6420acde7"
		fmt.Printf("Connected - v16 (watch for events) - on net(%s)\n", GethServer)
	}

	hexAddr := common.HexToAddress(ContractAddress)
	ev.ContractAddress = ContractAddress

	// Create a object to talk to the contract.
	ks, err := KeepRelayBeacon.NewKeepRelayBeacon(hexAddr, conn)
	if err != nil {
		log.Fatalf("Failed to instantiate contract object: %v at address: %s", err, GethServer)
	}

	ev.ks = ks

	return
}

// TODO
//  	function isStaked(uint256 _UserPublicKey) pure public returns(bool) {
//		function getRandomNumber( uint256 _RequestID ) view public returns ( uint256 theRandomNumber ) {
//  	function relayEntryAcustation( uint256 _LastValidRelayTxHash, uint256 _LastValidRelayBlock) public {
//  	function submitGroupPublicKey(uint256 _PK_G_i, uint256 _id) public {
//		function getRandomNumber( uint256 _RequestID ) view public returns ( uint256 theRandomNumber ) {

// ---------------------------------------------------------------------------------------------------------------------------------------------------------------------

//
// To Create 'sink':
//
//	sink  := make(chan *KeepRelayBeacon.KeepRelayBeaconRelayEntryRequested, 10)
//
// You may nee to increates the size of the chanel from 10 to 128 - that is the size used inside bind in the
// go-ethereum code.
//
// Then:
//
// 	event, err := ev.WatchKeepRelayBeaconRelayEntryRequested(ctx, sink)
//
//	for {
//		select {
//		case rn := <-sink:
//			// 'rn' has the data in it from the event
//
//		case ee := <-event.Err():
//			err = fmt.Errorf("Error watching for KeepRelayBeacon.RelayEntryRequested: %s", ee)
//			// process the error - note - an EOF error will not wait - so you need to exit loop on an error
//			return
//		}
//	}
//

func (ev *KeepRelayBeaconEvents) WatchKeepRelayBeaconRelayEntryRequested(ctx *RelayContractContext, sink chan<- *KeepRelayBeacon.KeepRelayBeaconRelayEntryRequested) (es event.Subscription, err error) {

	es, err = ev.ks.WatchRelayEntryRequested(nil, sink)
	if err != nil {
		err = fmt.Errorf("Error creating watch for relay events: %s\n", err)
		// LOG at this point
		return
	}
	return

}

type RelayRequestEventCallback func(data *KeepRelayBeacon.KeepRelayBeaconRelayEntryRequested, errIn error) error

func (ev *KeepRelayBeaconEvents) CallbackKeepRelayBeaconRelayEntryRequested(ctx *RelayContractContext, fx RelayRequestEventCallback) (err error) {

	sink := make(chan *KeepRelayBeacon.KeepRelayBeaconRelayEntryRequested, 10)
	if ctx.dbOn() {
		fmt.Printf("Calling Watch for KeepRelayBeacon.RelayEntryRequested, from=%s\n", godebug.LF())
	}
	event, err := ev.ks.WatchRelayEntryRequested(nil, sink)
	if err != nil {
		err = fmt.Errorf("Error creating watch for relay events: %s\n", err)
		// LOG at this point
		err = fx(nil, err)
		if err != nil {
			return
		}
		return
	}

	for {
		select {
		case rn := <-sink:
			if ctx.dbOn() {
				fmt.Printf("Got a relay request event! %+v from=%s\n", rn, godebug.LF())
				fmt.Printf("           In JSON format! %s\n", godebug.SVarI(rn))
			}
			// do some processing (Callback) or have the chanel exposed?
			err = fx(rn, nil)
			if err != nil {
				return
			}

		case ee := <-event.Err():
			err = fmt.Errorf("Error watching for KeepRelayBeacon.RelayEntryRequested: %s", ee)
			// LOG at this point
			err = fx(nil, err)
			if err != nil {
				return
			}
			return
		}
	}
}

// TODO
//    event RelayEntryGenerated(uint256 RequestID, uint256 Signature, uint256 GroupID, uint256 PreviousEntry ); // xyzzy - RelayEntryGenerated.
//    event RelayResetEvent(uint256 LastValidRelayEntry, uint256 LastValidRelayTxHash, uint256 LastValidRelayBlock);	// xyzzy - data types on TxHash, Block
//    event SubmitGroupPublicKeyEvent(uint256 _PK_G_i, uint256 _id, uint256 _activationBlockHeight);

type RelayContractContext struct {
	dbOnFlag bool
}

func (ctx *RelayContractContext) dbOn() bool {
	return ctx.dbOnFlag
}

func (ctx *RelayContractContext) SetDebug(b bool) {
	ctx.dbOnFlag = b
}
