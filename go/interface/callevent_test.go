package RelayContract

// Note: setup ./exampels/setup.json first

// TODO: figure out how this should be automated - in a reasonable auto-test fashion.

import (
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/keep-network/keep-core/go/BeaconConfig"
	"github.com/keep-network/keep-core/go/interface/lib/KeepRelayBeacon"
	"github.com/pschlump/godebug"
)

func Test_CallEvent(t *testing.T) {
	// setup a call, then catch the event, then a 2nd call, then catch 2nd event

	fmt.Printf("Note: test takes about 40 seconds to run\n")

	quit := make(chan struct{}, 2)
	ctx := &RelayContractContext{}
	ctx.SetDebug(false)
	// cfg := excfg.ReadCfg("./test_setup.json")
	requestID := int64(-1)

	gcfg := BeaconConfig.GetBeaconConfig("./test/config2.json")
	conn, err := BeaconConfig.OpenEthConnection(gcfg) // xyzzy RPC or ETH?
	if err != nil {
		t.Errorf("Failed to read config file")
	}

	// ----------------------------------------------------------------------------
	// First Call - Setup to watch for events!
	// This is essentially the client top level code.
	// ----------------------------------------------------------------------------
	go func() {

		// ev, err := NewKeepRelayBeaconEvents(ctx, cfg.GethServer, cfg.ContractAddress)
		ev, err := NewKeepRelayBeaconEvents(ctx, conn, gcfg)
		if err != nil {
			t.Errorf("Error connecing to contract: %s\n", err)
			return
		}

		sink := make(chan *KeepRelayBeacon.KeepRelayBeaconRelayEntryRequested, 10) // xyzzy - should be Requested!
		event, err := ev.WatchKeepRelayBeaconRelayEntryRequested(ctx, sink)

		for {
			select {
			case rn := <-sink:
				fmt.Printf("Success Event Data: %s\n", godebug.SVarI(rn))
				requestID = rn.RequestID.Int64()
				// This is the place where the **signature** genration should start.
				// SignatureGenerated, err := host.GenerateANumber( rn.RequestID, rn.Seed )
				// After the **signature** is generated call (See AAA_444) below.

			case ee := <-event.Err():
				err = fmt.Errorf("Error watching for KeepRelayBeacon.RelayEntryRequested: %s", ee)
				// process the error - note - an EOF error will not wait - so you need to exit loop on an error
				t.Errorf("%s", err)
				return

			case <-quit:
				return
			}
		}

	}()

	// ----------------------------------------------------------------------------
	// Catch 2nd event **signature** complete.
	// ----------------------------------------------------------------------------

	go func() {

		//func NewKeepRelayBeaconEvents(ctx *RelayContractContext, conn *ethclient.Client,
		//	gcfg BeaconConfig.BeaconConfig) (ev *KeepRelayBeaconEvents, err error) {
		// ev, err := NewKeepRelayBeaconEvents(ctx, cfg.GethServer, cfg.ContractAddress)
		ev, err := NewKeepRelayBeaconEvents(ctx, conn, gcfg)
		if err != nil {
			t.Errorf("Error connecing to contract: %s\n", err)
			return
		}

		err = ev.CallbackKeepRelayBeaconRelayEntryGenerated(ctx,
			func(data *KeepRelayBeacon.KeepRelayBeaconRelayEntryGenerated, errIn error) (err error) {
				if errIn != nil {
					t.Errorf("Error: %s\n", errIn)
				} else {
					fmt.Printf("Success Event Data: %s\n", godebug.SVarI(data))
				}
				return
			})
		if err != nil {
			t.Errorf("Error on event callback: %s\n", err)
		}

	}()

	// ----------------------------------------------------------------------------
	// Now Call - to KeepRelayBeacon.RelayRequest to simulate a user making a
	// request for a **signature**.
	// ----------------------------------------------------------------------------

	// OLD: ri, err := NewKeepRelayBeaconContract(ctx, cfg.GethServer, cfg.ContractAddress, cfg.KeyFile, cfg.KeyFilePassword)
	ri, err := NewKeepRelayBeaconContract(ctx, conn, gcfg)
	if err != nil {
		t.Errorf("Error connecing to contract: %s\n", err)
		return
	}

	// from := common.Address{}
	from := common.HexToAddress(gcfg.FromAddress)

	tx, err := ri.RequestRelay(ctx, from, 21, 42, []byte("aabbccddee"))
	if err != nil {
		t.Errorf("Error call contract: %s\n", err)
		return
	}

	fmt.Printf("KeepRelayBeacon.RequestRelay called: Tx = %s\n", godebug.SVarI(tx))

	// sleep for 14 sec, give it time to process the block
	time.Sleep(14 * time.Second)

	// ----------------------------------------------------------------------------
	// This is the section of code that simulates the generated **signature**
	// has been successfully competed (Note: AAA_444 mentioned above)
	// ----------------------------------------------------------------------------

	// generate/call to generate 2nd event - register **signature** complete.  The
	// dummy number tha is used is 'aabbccddee'.
	tx, err = ri.RelayEntry(ctx, requestID, []byte("aabbccddee"), []byte("aabcdefghi"), []byte("xxuuvv"))
	if err != nil {
		t.Errorf("Error call contract: %s\n", err)
		return
	}

	time.Sleep(14 * time.Second)

	// send event on "quit" channel to end test.
	quit <- struct{}{}
	quit <- struct{}{}
	return

}
