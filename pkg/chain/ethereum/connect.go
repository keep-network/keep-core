package ethereum

import (
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// Provider contains the connections to the ethereum server.
type Provider struct {
	Config    Config
	Client    *ethclient.Client // connection with ws:// or .ipc for event filering calls
	ClientRPC *rpc.Client       // JSON RPC connection
}

// Connect makes the network connection to the Ethereum network.  Note: for other things
// to work correctly the configuration will need to reference a websocket, "ws://", or
// local IPC connection.
func Connect(cfg Config) (*Provider, error) {
	client, err := ethclient.Dial(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("error Connecting to Geth Server: %s server %s",
			err, cfg.URL)
	}

	clientrpc, err := rpc.Dial(cfg.URLRPC)
	if err != nil {
		return nil, fmt.Errorf("error Connecting to Geth Server: %s server %s",
			err, cfg.URL)
	}

	return &Provider{
		Config:    cfg,
		Client:    client,
		ClientRPC: clientrpc,
	}, nil
}
