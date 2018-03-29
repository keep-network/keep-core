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
	"crypto/rand"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/keep-network/keep-core/go/BeaconConfig"
	"github.com/keep-network/keep-core/go/interface/lib/KeepRelayBeacon"
	"github.com/pschlump/godebug"
)

// KeepRelayBeaconContract contains the connection information to call the contract.
// Call NewKeepRelayBeaconContract to build this.
type KeepRelayBeaconContract struct {
	conn            *ethclient.Client
	kstart          *KeepRelayBeacon.KeepRelayBeaconTransactor
	kstartCaller    *KeepRelayBeacon.KeepRelayBeaconCaller
	optsTransact    *bind.TransactOpts
	optsCall        *bind.CallOpts
	ContractAddress string
	GethServer      string
}

func OpenConnection(gcfg BeaconConfig.BeaconConfig) (conn *ethclient.Client, err error) {

	// Create JSON RPC connection to a remote node
	//
	// Examples of connection strings - you can use .ipc, ws: or http:
	// connections for this.
	//
	// GethServer := "http://10.51.245.75:8545"
	// GethServer := "ws://10.51.245.75:8546" // WebSocket connect
	// GethServer := "/Users/corwin/Projects/eth/data/geth.ipc" // IPC local connection
	// GethServer := "ws://192.168.0.157:8546" // Based on running a local "geth"
	conn, err = ethclient.Dial(gcfg.GethServer)
	if err != nil {
		err = fmt.Errorf("Failed to connect to the Ethereum client: %s at address: %s", err, gcfg.GethServer)
		return
	}
	return
}

// NewKeepRelayBeaconContract uses a connection to a "geth" server and sets up for calls into the
// KeepRelayBeacon contract.  The 'GethServer' can be a "http://", "ws://" or an IPC file
// path to geth.  The 'ContractAddress' needs to be the hex address of the KeepRelayBeacon
// contrct. 0x in front of the address is required.  A full path to, and the password for, an
// account file is required.  The account private key file is "UTC--" a date/time and the hash
// of the account number.
func NewKeepRelayBeaconContract(ctx *RelayContractContext, conn *ethclient.Client,
	gcfg BeaconConfig.BeaconConfig) (ri *KeepRelayBeaconContract, err error) {

	// OLD: func NewKeepRelayBeaconContract(ctx *RelayContractContext, GethServer, ContractAddress, KeyFile,
	//	KeyFilePassword string) (ri *KeepRelayBeaconContract, err error) {

	ri = &KeepRelayBeaconContract{
		GethServer:      gcfg.GethServer,
		ContractAddress: gcfg.BeaconRelayContractAddress,
		conn:            conn,
	}

	// Address of the contract - extracted from the output from "truffle migrate --reset --network testnet"
	// ContractAddress := "0xe705ab560794cf4912960e5069d23ad6420acde7"
	hexAddr := common.HexToAddress(gcfg.BeaconRelayContractAddress)

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
	file, err := os.Open(gcfg.KeyFile) // Read in file
	if err != nil {
		err = fmt.Errorf("Failed to open keyfile: %s, %s", err, gcfg.KeyFile)
		return
	}

	// Create a new trasaction call (will mine) for the contract.
	optsTransact, err := bind.NewTransactor(bufio.NewReader(file), gcfg.KeyFilePassword)
	if err != nil {
		err = fmt.Errorf("Failed to read keyfile: %s, %s", err, gcfg.KeyFile)
		return
	}

	ri.optsTransact = optsTransact

	ri.optsCall = new(bind.CallOpts)

	return
}

// RequestRelay calls into the KeepRelayBeacon.requestRelay - parameters are converted from 64 bit to 256 bit integers.
// The contract generates an event, RelayEntryRequested with the generated RequestID in it.
func (ri *KeepRelayBeaconContract) RequestRelay(ctx *RelayContractContext, requestFrom common.Address, payment, blockReward int64,
	seed []byte) (tx *types.Transaction, err error) {
	pBlockReward := big.NewInt(blockReward)
	pPayment := big.NewInt(payment)
	pSeed := big.NewInt(0).SetBytes(seed)
	// Contract: function requestRelay(uint256 _payment, uint256 _blockReward,
	// uint256 _seed) public returns ( uint256 RequestID ) {

	/*
		// CallOpts is the collection of options to fine tune a contract call request.
		type CallOpts struct {
			Pending bool           // Whether to operate on the pending state or the last known one
			From    common.Address // Optional the sender address, otherwise the first account is used

			Context context.Context // Network context to support cancellation and timeouts (nil = no timeout)
		}

		// TransactOpts is the collection of authorization data required to create a
		// valid Ethereum transaction.
		type TransactOpts struct {
			From   common.Address // Ethereum account to send the transaction from
			Nonce  *big.Int       // Nonce to use for the transaction execution (nil = use pending state)
			Signer SignerFn       // Method to use for signing the transaction (mandatory)

			Value    *big.Int // Funds to transfer along along the transaction (nil = 0 = no funds)
			GasPrice *big.Int // Gas price to use for the transaction execution (nil = gas price oracle)
			GasLimit uint64   // Gas limit to set for the transaction execution (0 = estimate)

			Context context.Context // Network context to support cancellation and timeouts (nil = no timeout)
		}
	*/

	// payment needs to be sent as "msg.value" - as a part of the standard message
	ri.optsTransact.From = requestFrom
	ri.optsTransact.Value = pPayment
	cRand, err := GenerateRandomBytes(32)
	if err != nil {
		err = fmt.Errorf("Error: %s on call to generate random Nonce (geth=%s) KeepRelayBeacon.at(0x%s).requestRelay(%d,%d,%x)\n",
			err, ri.GethServer, ri.ContractAddress, payment, blockReward, seed)
		// LOG at this point
		return
	}
	ri.optsTransact.Nonce = big.NewInt(0).SetBytes(cRand)

	// func (_KeepRelayBeacon *KeepRelayBeaconTransactor) RequestRelay(opts *bind.TransactOpts, _blockReward *big.Int, _seed *big.Int) (*types.Transaction, error) {
	tx, err = ri.kstart.RequestRelay(ri.optsTransact, pBlockReward, pSeed)
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

// RelayEntry saves the signature to the blockchain on completion threshold
// signature generation.  KeepRelayBeacon.relayEntry is called.
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
	tx, err = ri.kstart.RelayEntry(ri.optsTransact, pRequestID, pGroupSignature, pGroupID, pPreviousEntry)
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

func (ri *KeepRelayBeaconContract) RelayEntryBI(ctx *RelayContractContext,
	RequestID, GroupSignature, GroupID, PreviousEntry *big.Int) (tx *types.Transaction, err error) {

	// Contract: function relayEntry(uint256 _RequestID, uint256 _groupSignature,
	// uint256 _groupID, uint256 _previousEntry) public {
	tx, err = ri.kstart.RelayEntry(ri.optsTransact, RequestID, GroupSignature, GroupID, PreviousEntry)
	if err != nil {
		err = fmt.Errorf("Error: %s on call to (geth=%s) KeepRelayBeacon.at(0x%s).relayEntry(%d,%s,%s,%s)\n",
			err, ri.GethServer, ri.ContractAddress, RequestID, GroupSignature, GroupID, PreviousEntry)
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
	tx, err = ri.kstart.RelayEntryAccusation(ri.optsTransact, pLastValidRelayTxHash, pLastValidRelayBlock)
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
	tx, err = ri.kstart.SubmitGroupPublicKey(ri.optsTransact, pPK_G_i, pId)
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

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

/*
 */
