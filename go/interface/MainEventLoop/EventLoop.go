package EventLoop

// This file is not complete - it is our notes from Friday 23rd talks.

// Create a "main" event loop that precesses all of the work in the application

import (
	"fmt"
	"math/big"

	RelayContract "github.com/keep-network/keep-core/go/interface"
	"github.com/keep-network/keep-core/go/interface/examples/excfg"
)

// Main event loop code

type P2PSetup struct {
	// all the stuff that goes with P2P communication
}

type ContrctSetup struct {
	ctx RelayContract.RelayContractContext
	cfg excfg.Cfg
	ev  *RelayContract.KeepRelayBeaconEvents
}

/* Based on:
type KeepRelayBeaconRelayEntryGenerated struct {
	RequestID     *big.Int
	Signature     *big.Int
	GroupID       *big.Int
	PreviousEntry *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}
*/
type RelayReady struct {
	RequestID   *big.Int
	Payment     *big.Int
	BlockReward *big.Int
	Seed        *big.Int
}

type SignatureReady struct {
	RequestID     *big.Int
	Signature     *big.Int
	GroupID       *big.Int
	PreviousEntry *big.Int
}

type P2PSendBroadcast struct {
	From string
	Msg  string
}
type P2PSendTo struct {
	To   string
	From string
	Msg  string
}

type P2PMessage struct {
	Msg string
}

var RelayReadyChannel chan RelayReady
var P2PMessageChannel chan P2PMessage
var SignatureReadyChannel chan SignatureReady
var quit chan struct{}

func init() {
	RelayReadyChannel = make(chan RelayReady)
	P2PMessageChannel = make(chan P2PMessage)
	SignatureReadyChannel = make(chan SignatureReady)
	quit = make(chan struct{}, 2)
}

func MainEventLoop() {

	p2p := initializeP2P()
	eth := initializeContract()

	_ = p2p
	_ = eth

	for {
		select {
		case rd := <-RelayReadyChannel:
			// Generate *signature*
			// send message on channel that signature is ready
			_ = rd

		case sig := <-SignatureReadyChannel:
			// Got a signature - do something
			_ = sig

		case msg := <-P2PMessageChannel:
			// got a broadcast or a one-on-one message
			_ = msg
		}
	}

}

func initializeP2P() (rv P2PSetup) {
	// do a bunch of stuff to set up P2P communication
	return
}

func initializeContract() (rv ContrctSetup) {
	// Setup to talk to contract
	rv.ctx.SetDebug(false)
	rv.cfg = excfg.ReadCfg("./user_setup.json")
	ev, err := RelayContract.NewKeepRelayBeaconEvents(&rv.ctx, rv.cfg.GethServer, rv.cfg.ContractAddress)
	if err != nil {
		fmt.Printf("Error connecing to contract: %s\n", err)
		return
	}
	rv.ev = ev
	return
}
