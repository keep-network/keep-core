// Package RelayContract provides an interface in Go to the KeepRelayBeacon Ethereum
// contract.  The contrct souce can be found in ../contract/Threshold/src/KeepRelayBeacon.src.sol
// The contract uses a `.m4` file to get the address of the RequestID generator.
// The RequestID generator is in ../contract/GenID/contracts/GenRequestID.go
//
// TODO: - Add in "logging" - how - where - what system are we going to use?
// Probably add in some configuration with the context on how/where the
// logging is to take place.
//

package RelayContract

import (
	"bufio"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/keep-network/keep-core/go/interface/lib/KeepRelayBeacon"
	"github.com/pschlump/godebug"
)

// KeepRelayBeaconContract contains the connection information to call the contract.
// Call NewKeepRelayBeaconContract to build this.
type KeepRelayBeaconContract struct {
	conn            *ethclient.Client
	kstart          *KeepRelayBeacon.KeepRelayBeaconTransactor
	kstartCaller    *KeepRelayBeacon.KeepRelayBeaconCaller
	opts            *bind.TransactOpts
	optsCall        *bind.CallOpts
	ContractAddress string
	GethServer      string
}

// NewKeepRelayBeaconContract connects to a "geth" server and sets up for calls into the
// KeepRelayBeacon contract.  The 'GethServer' can be a "http://", "ws://" or an IPC file
// path to geth.  The 'ContractAddress' needs to be the hex address of the KeepRelayBeacon
// contrct. 0x in front of the address is required.  A full path to, and the password for, an
// account file is required.  The account private key file is "UTC--" a date/time and the hash
// of the account number.
func NewKeepRelayBeaconContract(ctx *RelayContractContext, GethServer, ContractAddress, KeyFile,
	KeyFilePassword string) (ri *KeepRelayBeaconContract, err error) {

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

	ri.optsCall = new(bind.CallOpts)

	return
}

// RequestRelay calls into the KeepRelayBeacon.requestRelay - parameters are converted from 64 bit to 256 bit integers.
// The contract generates an event, RelayEntryRequested with the generated RequestID in it.
func (ri *KeepRelayBeaconContract) RequestRelay(ctx *RelayContractContext, payment, blockReward int64,
	seed []byte) (tx *types.Transaction, err error) {
	pPayment := big.NewInt(payment)
	pBlockReward := big.NewInt(blockReward)
	pSeed := big.NewInt(0).SetBytes(seed)
	// Contract: function requestRelay(uint256 _payment, uint256 _blockReward,
	// uint256 _seed) public returns ( uint256 RequestID ) {
	tx, err = ri.kstart.RequestRelay(ri.opts, pPayment, pBlockReward, pSeed)
	if err != nil {
		err = fmt.Errorf("Error: %s on call to (geth=%s) KeepRelayBeacon.at(0x%s).requestRelay(%d,%d,%x)\n",
			err, ri.GethServer, ri.ContractAddress, payment, blockReward, seed)
		// LOG at this point
		return
	}
	if ctx.dbOn() {
		fmt.Printf("KeepRelayBeacon.RequestRelay called tx=%s from=%s\n", godebug.SVar(tx), godebug.LF())
		// TODO: log the success and tx that it is in
	}
	return
}

// IsStaked returns true if the publicKey is associated with a sufficiently large
// stake in the registry.  KeepRelayBeacon.isStaked is called.
// No gas is required to call this.
func (ri *KeepRelayBeaconContract) IsStaked(ctx *RelayContractContext, publicKey []byte) (memberIsStaked bool) {
	pPubKey := big.NewInt(0).SetBytes(publicKey)
	// Contract: function isStaked(uint256 _UserPublicKey) pure public returns(bool) {
	memberIsStaked, err := ri.kstartCaller.IsStaked(ri.optsCall, pPubKey)
	if err != nil {
		err = fmt.Errorf("Error: %s on call to (geth=%s) KeepRelayBeacon.at(0x%s).isStaked(%x)\n",
			err, ri.GethServer, ri.ContractAddress, publicKey)
		// LOG at this point
		return false
	}
	if ctx.dbOn() {
		fmt.Printf("KeepRelayBeacon.isStaked called pubKey=%x isStaked()=%v, from=%s\n",
			publicKey, memberIsStaked, godebug.LF())
		// TODO: log the success and tx that it is in
	}
	return
}

// RelayEntry saves the signature to the blockchain on completion threshold signature generation.
// KeepRelayBeacon.relayEntry is called.
// An event, RelayEntryGenerated, is generated in the contract.
func (ri *KeepRelayBeaconContract) RelayEntry(ctx *RelayContractContext,
	requestID int64, groupSignature, groupID,
	previousEntry []byte) (tx *types.Transaction, err error) {
	pRequestID := big.NewInt(requestID)
	pGroupSignature := big.NewInt(0).SetBytes(groupSignature)
	pGroupID := big.NewInt(0).SetBytes(groupID)
	pPreviousEntry := big.NewInt(0).SetBytes(previousEntry)
	// Contract: function relayEntry(uint256 _RequestID, uint256 _groupSignature,
	// uint256 _groupID, uint256 _previousEntry) public {
	tx, err = ri.kstart.RelayEntry(ri.opts, pRequestID, pGroupSignature, pGroupID, pPreviousEntry)
	if err != nil {
		err = fmt.Errorf("Error: %s on call to (geth=%s) KeepRelayBeacon.at(0x%s).relayEntry(%d,%x,%x,%x)\n",
			err, ri.GethServer, ri.ContractAddress, requestID, groupSignature,
			groupID, previousEntry)
		// LOG at this point
		return
	}
	if ctx.dbOn() {
		fmt.Printf("KeepRelayBeacon.relayEntry called tx=%s from=%s\n", godebug.SVar(tx), godebug.LF())
		// TODO: log the success and tx that it is in
	}
	return
}

// RelayEntryAccusation calls KeepRelayBeacon.relayEntryAccusation to make the claim that
// the generated signature is deceptive.
// The underlying contract code has not been implemented yet.
func (ri *KeepRelayBeaconContract) RelayEntryAccusation(ctx *RelayContractContext,
	lastValidRelayTxHash, lastValidRelayBlock []byte) (tx *types.Transaction, err error) {
	pLastValidRelayTxHash := big.NewInt(0).SetBytes(lastValidRelayTxHash)
	pLastValidRelayBlock := big.NewInt(0).SetBytes(lastValidRelayBlock)
	// Contrct: function relayEntryAccusation( uint256 _LastValidRelayTxHash, uint256 _LastValidRelayBlock) public {
	tx, err = ri.kstart.RelayEntryAccusation(ri.opts, pLastValidRelayTxHash, pLastValidRelayBlock)
	if err != nil {
		err = fmt.Errorf("Error: %s on call to (geth=%s) KeepRelayBeacon.at(0x%s).relayEntryAccusation(%x,%x)\n",
			err, ri.GethServer, ri.ContractAddress, pLastValidRelayTxHash, pLastValidRelayBlock)
		// LOG at this point
		return
	}
	if ctx.dbOn() {
		fmt.Printf("KeepRelayBeacon.relayEntryAccusation called tx=%s from=%s\n", godebug.SVar(tx), godebug.LF())
		// TODO: log the success and tx that it is in
	}
	return
}

// SubmitGroupPublicKey calls KeepRelayBeacon.submitGroupPublicKey
// The contract generates an event, RelayEntryRequested with the generated RequestID in it.
func (ri *KeepRelayBeaconContract) SubmitGroupPublicKey(ctx *RelayContractContext,
	PK_G_i, id []byte) (tx *types.Transaction, err error) {
	pPK_G_i := big.NewInt(0).SetBytes(PK_G_i)
	pId := big.NewInt(0).SetBytes(id)
	// Contract: function submitGroupPublicKey(uint256 _PK_G_i, uint256 _id) public {
	tx, err = ri.kstart.SubmitGroupPublicKey(ri.opts, pPK_G_i, pId)
	if err != nil {
		err = fmt.Errorf("Error: %s on call to (geth=%s) KeepRelayBeacon.at(0x%s).submitGroupPublicKey(%x,%x)\n",
			err, ri.GethServer, ri.ContractAddress, PK_G_i, id)
		// LOG at this point
		return
	}
	if ctx.dbOn() {
		fmt.Printf("KeepRelayBeacon.RequestRelay called tx=%s from=%s\n", godebug.SVar(tx), godebug.LF())
		// TODO: log the success and tx that it is in
	}
	return
}

/*
 */
