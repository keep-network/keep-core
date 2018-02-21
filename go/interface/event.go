package RelayContract

// TODO
//    event RelayResetEvent(uint256 LastValidRelayEntry, uint256 LastValidRelayTxHash, uint256 LastValidRelayBlock);	// xyzzy - data types on TxHash, Block
//    event SubmitGroupPublicKeyEvent(uint256 _PK_G_i, uint256 _id, uint256 _activationBlockHeight);
import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/event"
	"github.com/keep.network/keep-core/go/interface/lib/KeepRelayBeacon"
	"github.com/pschlump/godebug"
)

type KeepRelayBeaconEvents struct {
	ks              *KeepRelayBeacon.KeepRelayBeacon
	conn            *ethclient.Client
	ContractAddress string
	GethServer      string
}

// NewKeepRelayBeaconEvents performs the necessary setup to watch/capture events from the KeepRelayBeacon contract.
// NOTE! GethServer must be a "ws://" or IPC connection to the geth server.  A http connection will not support "push" events.
func NewKeepRelayBeaconEvents(ctx *RelayContractContext, GethServer, ContractAddress string) (ev *KeepRelayBeaconEvents, err error) {

	ev = &KeepRelayBeaconEvents{
		GethServer:      GethServer,
		ContractAddress: ContractAddress,
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
	conn, err := ethclient.Dial(GethServer)
	if err != nil {
		err = fmt.Errorf("Failed to connect to the Ethereum client: %s at address: %s", err, GethServer)
		return
	}

	ev.conn = conn

	if ctx.dbOn() {
		// Address of the contract - extracted from the output from "truffle migrate --reset --network testnet"
		// ContractAddress := "0xe705ab560794cf4912960e5069d23ad6420acde7"
		fmt.Printf("Connected - v16 (watch for events) - on net(%s)\n", GethServer)
	}

	hexAddr := common.HexToAddress(ContractAddress)

	// Create a object to talk to the contract.
	ks, err := KeepRelayBeacon.NewKeepRelayBeacon(hexAddr, conn)
	if err != nil {
		log.Fatalf("Failed to instantiate contract object: %v at address: %s", err, GethServer)
	}

	ev.ks = ks

	return
}

// ---------------------------------------------------------------------------------------------------------------------------------------------------------------------
//
// Setup:
//
//	ctx := &RelayContractContext{}
// 	ev, err :=  NewKeepRelayBeaconEvents(ctx, "ws://127.0.0.1:8546", // Based on running a local "geth"
// 		 "0xe705ab560794cf4912960e5069d23ad6420acde7" )
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

// WatchKeepRelayBeaconRelayEntryRequested is a wrapper around the abigen WatchRelayEntryRequested
func (ev *KeepRelayBeaconEvents) WatchKeepRelayBeaconRelayEntryRequested(ctx *RelayContractContext, sink chan<- *KeepRelayBeacon.KeepRelayBeaconRelayEntryRequested) (es event.Subscription, err error) {

	// the 'nil' in the next line is really really imporant!
	es, err = ev.ks.WatchRelayEntryRequested(nil, sink)
	if err != nil {
		err = fmt.Errorf("Error creating watch for relay events: %s\n", err)
		// LOG at this point
		return
	}
	return

}

// Type of callback function
type RelayEntryRequestedCallback func(data *KeepRelayBeacon.KeepRelayBeaconRelayEntryRequested, errIn error) error

// CallbackKeepRelayBeaconRelayEntryRequested watches for RelayEntryRequested event from the contract and calls 'fx' when the event occures.
func (ev *KeepRelayBeaconEvents) CallbackKeepRelayBeaconRelayEntryRequested(ctx *RelayContractContext, fx RelayEntryRequestedCallback) (err error) {

	sink := make(chan *KeepRelayBeacon.KeepRelayBeaconRelayEntryRequested, 10)
	if ctx.dbOn() {
		fmt.Printf("Calling Watch for KeepRelayBeacon.RelayEntryRequested, from=%s\n", godebug.LF())
	}
	// the 'nil' in the next line is really really imporant!
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

// WatchKeepRelayBeaconRelayEntryRequested is a wrapper around the abigen WatchRelayEntryGenerated
func (ev *KeepRelayBeaconEvents) WatchKeepRelayBeaconRelayEntryGenerated(ctx *RelayContractContext, sink chan<- *KeepRelayBeacon.KeepRelayBeaconRelayEntryGenerated) (es event.Subscription, err error) {

	// the 'nil' in the next line is really really imporant!
	es, err = ev.ks.WatchRelayEntryGenerated(nil, sink)
	if err != nil {
		err = fmt.Errorf("Error creating watch for relay events: %s\n", err)
		// LOG at this point
		return
	}
	return

}

// Type of callback function
type RelayEntryGeneratedCallback func(data *KeepRelayBeacon.KeepRelayBeaconRelayEntryGenerated, errIn error) error

// CallbackKeepRelayBeaconRelayEntryGenerated watches for RelayEntryGenerated event from the contract and calls 'fx' when the event occures.
func (ev *KeepRelayBeaconEvents) CallbackKeepRelayBeaconRelayEntryGenerated(ctx *RelayContractContext, fx RelayEntryGeneratedCallback) (err error) {

	sink := make(chan *KeepRelayBeacon.KeepRelayBeaconRelayEntryGenerated, 10)
	if ctx.dbOn() {
		fmt.Printf("Calling Watch for KeepRelayBeacon.RelayEntryGenerated, from=%s\n", godebug.LF())
	}
	// the 'nil' in the next line is really really imporant!
	event, err := ev.ks.WatchRelayEntryGenerated(nil, sink)
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
			err = fmt.Errorf("Error watching for KeepRelayBeacon.RelayEntryGenerated: %s", ee)
			// LOG at this point
			err = fx(nil, err)
			if err != nil {
				return
			}
			return
		}
	}
}
