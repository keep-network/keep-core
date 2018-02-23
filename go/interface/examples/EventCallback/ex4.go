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

	err = ev.CallbackKeepRelayBeaconRelayEntryRequested(ctx,
		func(data *KeepRelayBeacon.KeepRelayBeaconRelayEntryRequested, errIn error) (err error) {
			if errIn != nil {
				fmt.Printf("Error: %s\n", errIn)
			} else {
				fmt.Printf("Success Event Data: %s\n", godebug.SVarI(data))
			}
			return
		})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error on event callback: %s\n", err)
		os.Exit(1)
	}

}
