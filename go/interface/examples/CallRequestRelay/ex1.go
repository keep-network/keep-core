package main

import (
	"fmt"
	"os"

	RelayContract "github.com/keep.network/keep-core/go/interface"
	"github.com/keep.network/keep-core/go/interface/examples/excfg"
	"github.com/pschlump/godebug"
)

func main() {

	cfg := excfg.ReadCfg("../setup.json")

	ctx := &RelayContract.RelayContractContext{}
	ctx.SetDebug(true)

	ri, err := RelayContract.NewKeepRelayBeaconContract(ctx, cfg.GethServer, cfg.ContractAddress,
		cfg.KeyFile, cfg.KeyFilePassword)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecing to contract: %s\n", err)
		os.Exit(1)
	}

	tx, err := ri.RequestRelay(ctx, 21, 42, []byte("aabbccddee"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error call contract: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Success: Tx = %s\n", godebug.SVarI(tx))

}
