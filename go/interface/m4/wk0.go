package main

import (
	"bufio"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind" // "www.2c-why.com/job/ethrpc/ethABI/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/keep.network/keep-core/go/interface/m4/KStart" // "www.2c-why.com/job/m4/KStart"
)

// "www.2c-why.com/job/m4/KStart"
// "github.com/ethereum/go-ethereum/ethclient"
// "github.com/ethereum/go-ethereum/common"

type KeepGenNumber struct {
	contractAddress string
	conn            *ethclient.Client
	contract        *KStart.KStart
	keyFile         string
	opts            *bind.TransactOpts
}

type RelayRequest struct {
	RequestID   int64  //
	Payment     int64  // Eth to run (gas?)
	BlockReward int64  // Keep tokens as payment
	Seed        []byte //
}

type RelayEntry struct {
	RequestID      int    //
	RandomNumber   []byte // generated output, the Signature, the Public Key for this, the Output
	GroupPublicKey []byte //
	PrivateEntry   string //	???
}

func TstCall() {
	ethNode := "http://127.0.0.1:9545"
	caddr := "0xaa588d3737b611bafd7bd713445b314bd453a5c8"
	// ethNode = "http://192.168.0.139:8545"

	// keyFile := "/Users/corwin/go/src/www.2c-why.com/job/m4/UTC--2018-01-03T02-17-58.695623282Z--9980ecddef53089390136fde20feb7e03125c441"
	keyFile := "./UTC--2018-01-03T02-17-58.695623282Z--9980ecddef53089390136fde20feb7e03125c441"
	keyFilePassword := "password"
	kc, err := NewGenNumber(ethNode, caddr, keyFile, keyFilePassword)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	} else {
		fmt.Printf("Success: %+v\n", kc)
	}
}

// NewKeepConnection creates the connection to a contract
// caddr is the contract address
// ethNode is "http://127.0.0.1:9545", or http://192.168.0.139:8545 -- it could also be
// an IPC path ( ipc://home/pschlump/.ethereum/geth.ipc )
// if the etherium node is on the same machine and you have enabled IPC communication.
// https://medium.com/@PasschainBlog/jumping-into-truffle-and-rinkeby-3acf6a2d9bef
/*
pschlump@dev3:~$ geth --rinkeby account new
Your new account is locked with a password. Please give a password. Do not forget this password.
Passphrase:
Repeat passphrase:
Address: {9c3c3c8741f3a22b6728a19df0614b2638f8777a}
*/
/*
ethNode := "http://127.0.0.1:9545"
caddr := "0xaa588d3737b611bafd7bd713445b314bd453a5c8"
keyFile := "/Users/corwin/go/src/www.2c-why.com/job/m4/UTC--2018-01-03T02-17-58.695623282Z--9980ecddef53089390136fde20feb7e03125c441"
keyFilePassword := "password"
func NewGenNumber(ethNode, caddr, keyFile, keyFilePassword) (kc *KeepGenNumber, err error) {
*/
func NewGenNumber(ethNode, caddr, keyFile, keyFilePassword string) (kc *KeepGenNumber, err error) {
	kc = &KeepGenNumber{
		contractAddress: caddr,
		keyFile:         keyFile,
	}
	conn, err := ethclient.Dial(ethNode)
	if err != nil {
		err = fmt.Errorf("Failed to connect to the Ethereum client(%s): %v contract address: %s", ethNode, err, caddr)
		return
	}

	kc.conn = conn

	contract, err := KStart.NewKStart(common.HexToAddress(caddr), conn)
	if err != nil {
		log.Fatalf("Failed to instantiate the contract: %v", err)
	}

	kc.contract = contract

	// Question: should this come in from the user (not a file on local system?)

	file, err := os.Open(keyFile)
	if err != nil {
		log.Fatalf("Failed to open keyfile: %v, %s", err, keyFile)
	}

	opts, err := bind.NewTransactor(bufio.NewReader(file), keyFilePassword)
	if err != nil {
		log.Fatalf("Failed to read keyfile: %v, %s", err, keyFile)
	}

	kc.opts = opts

	return
}

func (kc *KeepGenNumber) CreateRelayRequest(rq *RelayRequest, seed []byte) (err error) {
	// TODO:
	// 1. Create the RequestID (Based on? as a TXid from ETH?)
	// 2. Call ETH to put request on chain?
	// 3. Lots of stuff at this point like check for valid inputs, convert data etc.
	//    This just compiles now.
	// var _RequestID *big.Int
	_RequestID := big.NewInt(rq.RequestID) // I am REALLY (did I say REALLY!) not fond of Eth's style with _ in front of parameters.
	_payment := big.NewInt(rq.Payment)
	_blockReward := big.NewInt(rq.BlockReward)
	_seed := big.NewInt(0).SetBytes(seed)
	// var opts *bind.TransactOpts
	// func (_KStart *KStartTransactor) RequestRandomNumber(opts *bind.TransactOpts, _RequestID *big.Int, _payment *big.Int, _blockReward *big.Int, _seed *big.Int) (*types.Transaction, error) {
	tx, err := kc.contract.RequestRandomNumber(kc.opts, _RequestID, _payment, _blockReward, _seed)
	_, _ = tx, err
	// xyzzy - what am I supposed to do with the output at this point?  Log the transaction? hm...
	return
}

// ReqlayRequestReturn -> RelayRequest (An Event?)

func (kc *KeepGenNumber) RelayEntryPublished(re RelayEntry, tipBowl string) (err error) {
	return
}

// func (_KStart *KStartFilterer) WatchRandomNumberReady(opts *bind.WatchOpts, sink chan<- *KStartRandomNumberReady) (event.Subscription, error) {
func (kc *KeepGenNumber) ProcessRandomNumbers() (err error) {
	var _KStart *KStart.KStartFilterer // xyzzy
	var opts *bind.WatchOpts           // xyzzy
	var sink chan *KStart.KStartRandomNumberReady
	event, err := _KStart.WatchRandomNumberReady(opts, sink)
	if err != nil {
		fmt.Printf("Error creating watch: %s\n", err)
	} else {
		_ = event // xyzzy - don't we need to do "go event()"???
		// Event is an event.Subscription - that is a chanel and the "Unsubscribe" -- So, Unsubscribe will get called,
		// do we need a "quit" event in this, do we need to catch/report errors from the err chanel?
		select {
		case rn := <-sink:
			fmt.Printf("Got a random number! Yea!, %+v\n", rn)

			//		case ee := <-event.Err()
			//			fmt.Printf ( "Got an error: %s\n", ee )
		}
	}

	return
}
