package main

import (
	"fmt"
	"os"

	RelayContract "github.com/keep.network/keep-core/go/interface"
	"github.com/keep.network/keep-core/go/interface/examples/excfg"
	"github.com/keep.network/keep-core/go/interface/lib/KeepRelayBeacon"
	"github.com/pschlump/godebug"
)

func main() {

	cfg := excfg.ReadCfg("../setup.json")

	ctx := &RelayContract.RelayContractContext{}
	ctx.SetDebug(true)

	ev, err := RelayContract.NewKeepRelayBeaconEvents(ctx, cfg.GethServer, cfg.ContractAddress)
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

		case ee := <-event.Err():
			err = fmt.Errorf("Error watching for KeepRelayBeacon.RelayEntryRequested: %s", ee)
			// process the error - note - an EOF error will not wait - so you need to exit loop on an error
			os.Exit(1)
		}
	}

}
