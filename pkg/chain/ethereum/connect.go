package ethereum

import (
	"fmt"

	"github.com/ethereum/go-ethereum/rpc"
)

type EthereumConnection struct {
	Config EthereumConfig
	Client *rpc.Client
}

// Connect makes the network connection to the Ethereum network.  Note: for other things
// to work correctly this will need to be a websocket "ws://" or a local IPC connection.
func Connect(cfg EthereumConfig) (*EthereumConnection, error) {
	client, err := rpc.Dial(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("FAIL: Error Connecting to Geth Server: %s server %s\n",
			err, cfg.URL)
	}

	return &EthereumConnection{
		Config: cfg,
		Client: client,
	}, nil
}
