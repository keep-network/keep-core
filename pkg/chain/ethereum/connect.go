package ethereum

import (
	"fmt"

	"github.com/ethereum/go-ethereum/rpc"
)

type provider struct {
	Config Config
	Client *rpc.Client
}

// Connect makes the network connection to the Ethereum network.  Note: for other things
// to work correctly the configuration will need to reference a websocket, "ws://", or
// local IPC connection.
func Connect(cfg Config) (*provider, error) {
	client, err := rpc.Dial(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("error Connecting to Geth Server: %s server %s",
			err, cfg.URL)
	}

	return &provider{
		Config: cfg,
		Client: client,
	}, nil
}
