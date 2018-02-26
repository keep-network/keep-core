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
	ri  *RelayContract.KeepRelayBeaconContract
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
	Signature     *big.Int // maybee should be []byte type
	GroupID       *big.Int // maybee should be []byte type
	PreviousEntry *big.Int // maybee should be []byte type
}

// Send a P2P message
type P2PSendBroadcast struct {
	From string
	Msg  string
}

// Send a P2P message
type P2PSendTo struct {
	To   string
	From string
	Msg  string
}

// Receive a P2P message
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

	// TODO: do we need a go routine that serves up current "status" on some port at this point?

	for {
		select {
		case rd := <-RelayReadyChannel:
			// Generate *signature*
			// send message on channel that signature is ready
			_ = rd

		case sig := <-SignatureReadyChannel:
			// Got a signature - Write it to the blockchain

			// Call contrct to save data
			tx, err := eth.ri.RelayEntryBI(&eth.ctx, sig.RequestID, sig.Signature, sig.GroupID, sig.PreviousEntry)

			_, _ = tx, err
			// TODO report status - log stuff ???

		case msg := <-P2PMessageChannel:
			// got a broadcast or a one-on-one message
			_ = msg

			// Other Chanels for Sending Messages across P2P

			// Other Chanels for responding to accusation of invalid *signature*

		}
	}

}

func initializeP2P() (rv P2PSetup) {
	// TODO: do a bunch of stuff to set up P2P communication
	return
}

func initializeContract() (rv ContrctSetup) {

	// Pull in a little bit of configuration stuff
	rv.cfg = excfg.ReadCfg("./user_setup.json")
	rv.ctx.SetDebug(rv.cfg.DebugFlag)

	// Setup to talk to contract
	ev, err := RelayContract.NewKeepRelayBeaconEvents(&rv.ctx, rv.cfg.GethServer, rv.cfg.ContractAddress)
	if err != nil {
		fmt.Printf("Error connecing to contract: %s\n", err)
		return
	}
	rv.ev = ev

	ri, err := RelayContract.NewKeepRelayBeaconContract(&rv.ctx, rv.cfg.GethServer,
		rv.cfg.ContractAddress, rv.cfg.KeyFile, rv.cfg.KeyFilePassword)
	if err != nil {
		fmt.Printf("Error connecing to contract: %s\n", err)
		return
	}
	rv.ri = ri

	return
}
