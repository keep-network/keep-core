package RelayContract

// Note: setup ./exampels/setup.json first

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/keep.network/keep-core/go/interface/examples/excfg"
	"github.com/keep.network/keep-core/go/interface/lib/KeepRelayBeacon"
	"github.com/pschlump/godebug"
)

func Test_CallEvent(t *testing.T) {
	// setup a call, then catch the event, then a 2nd call, then catch 2nd event

	fmt.Printf("Note: test takes about 40 seconds to run\n")

	quit := make(chan struct{}, 2)
	ctx := &RelayContractContext{}
	ctx.SetDebug(false)
	cfg := excfg.ReadCfg("./test_setup.json")
	requestID := int64(-1)

	// ----------------------------------------------------------------------------
	// First Call - Setup to watch for events!
	// ----------------------------------------------------------------------------
	go func() {

		ev, err := NewKeepRelayBeaconEvents(ctx, cfg.GethServer, cfg.ContractAddress)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error connecing to contract: %s\n", err)
			os.Exit(1)
		}

		sink := make(chan *KeepRelayBeacon.KeepRelayBeaconRelayEntryGenerated, 10)
		event, err := ev.WatchKeepRelayBeaconRelayEntryGenerated(ctx, sink)

		for {
			select {
			case rn := <-sink:
				fmt.Printf("Success Event Data: %s\n", godebug.SVarI(rn))
				requestID = rn.RequestID.Int64()

			case ee := <-event.Err():
				err = fmt.Errorf("Error watching for KeepRelayBeacon.RelayEntryRequested: %s", ee)
				// process the error - note - an EOF error will not wait - so you need to exit loop on an error
				return
			case <-quit:
				return
			}
		}

	}()

	go func() {

		ev, err := NewKeepRelayBeaconEvents(ctx, cfg.GethServer, cfg.ContractAddress)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error connecing to contract: %s\n", err)
			os.Exit(1)
		}

		err = ev.CallbackKeepRelayBeaconRelayEntryGenerated(ctx, func(data *KeepRelayBeacon.KeepRelayBeaconRelayEntryGenerated, errIn error) (err error) {
			if errIn != nil {
				fmt.Printf("Error: %s\n", errIn)
			} else {
				fmt.Printf("Success Event Data: %s\n", godebug.SVarI(data))
			}
			return
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error on event callback: %s\n", err)
		}

	}()

	// ----------------------------------------------------------------------------
	// Now Call - to KeepRelayBeacon.RelayRequest
	// ----------------------------------------------------------------------------

	ri, err := NewKeepRelayBeaconContract(ctx, cfg.GethServer, cfg.ContractAddress, cfg.KeyFile, cfg.KeyFilePassword)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecing to contract: %s\n", err)
		os.Exit(1)
	}

	tx, err := ri.RequestRelay(ctx, 21, 42, []byte("aabbccddee"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error call contract: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("KeepRelayBeacon.RequestRelay called: Tx = %s\n", godebug.SVarI(tx))

	// sleep for 14 sec, give it time to process the block
	time.Sleep(14 * time.Second)

	// generate/call for 2nd event
	tx, err = ri.RelayEntry(ctx, requestID, []byte("aabbccddee"), []byte("aabcdefghi"), []byte("xxuuvv"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error call contract: %s\n", err)
		os.Exit(1)
	}

	time.Sleep(14 * time.Second)

	// send event on "quit" channel to end test.
	quit <- struct{}{}
	quit <- struct{}{}
	os.Exit(0)

}
