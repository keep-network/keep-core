package ethereum

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/keep-network/keep-core/pkg/provider_init"
)

//
// TODO
//
// Put in Main
// import "github.com/keep-network/keep-core/pkg/chain/ethereum"
//...
//  provider := ethereum.Connect(config)
//

type EthereumConnection struct {
	Config provider_init.ProviderInit
	Client *rpc.Client
}

func Connect(cfg provider_init.ProviderInit) (rv *EthereumConnection) {
	rv = &EthereumConnection{}

	// Use configuration data to setup connection to Geth/Ethereum
	client, err := rpc.Dial(cfg.EthereumProvider.GethConnectionString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "FAIL: Error Connecting to Geth Server: %s server %s\n", err, cfg.EthereumProvider.GethConnectionString)
		os.Exit(1)
	}

	rv.Config = cfg
	rv.Client = client
	return rv
}
