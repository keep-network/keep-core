package ethereum

import (
	"fmt"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/keep-network/keep-core/pkg/config"
)

type EthereumConnection struct {
	Config config.Config
	Client *rpc.Client
}

// Connect makes the network connection to the Ethereum network.
func Connect(cfg config.Config) (*EthereumConnection, error) {
	client, err := rpc.Dial(cfg.Ethereum.URL)
	if err != nil {
		return nil, fmt.Errorf("FAIL: Error Connecting to Geth Server: %s server %s\n", err, cfg.Ethereum.URL)
	}

	return &EthereumConnection{
		Config: cfg,
		Client: client,
	}, nil
}
