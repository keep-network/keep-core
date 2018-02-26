package EventLoop

// This file is not complete - it is our notes from Friday 23rd talks.

// Create a "main" event loop that precesses all of the work in the application

import (
	"fmt"
	"math/big"

	"github.com/keep-network/keep-core/go/interface/examples/excfg"
)

// Main event loop code

type P2PSetup struct {
	// all the stuff that goes with P2P communication
}

type ContrctSetup struct {
	ctx RelayContractContext
	cfg interface{}
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
var quit chan struct{}

func init() {
	RelayReadyChannel = make(chan<- RelayReady)
	P2PMessageChannel = make(chan<- P2PMessage)
	quit = make(chan struct{}, 2)
}

func StartInit() {
	InitRelayChannel <- InitRelay{}
	InitContractChannel <- InitContract{}
}

func MainEventLoop() {

	p2p := initializeP2P()
	eth := initializeContract()

	for {
		select {
		case rd <- RelayReadyChannel:
			// Generate *signature*
			// send message on channel that signature is ready

		case sig <- SignatureReadyChannel:
			// Got a signature - do something

		case msg <- P2PMessageChannel:
			// got a broadcast or a one-on-one message
		}
	}

}

func initializeP2P() P2PSetup {
	// do a bunch of stuff to set up P2P communication
	return
}

func initializeContract() (rv ContrctSetup) {
	// Setup to talk to contract
	rv.ctx.SetDebug(false)
	rv.cfg = excfg.ReadCfg("./user_setup.json")
	ev, err := NewKeepRelayBeaconEvents(ctx, cfg.GethServer, cfg.ContractAddress)
	if err != nil {
		fmt.Printf("Error connecing to contract: %s\n", err)
		return
	}
	rv.ev = ev
	return
}
