package ethereum

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/keep-network/keep-core/pkg/chain"
)

type ethereumChain struct {
	config                   Config
	client                   *ethclient.Client
	clientRPC                *rpc.Client
	clientWS                 *rpc.Client
	requestID                *big.Int
	keepGroupContract        *keepGroup
	keepRandomBeaconContract *KeepRandomBeacon
	tx                       *types.Transaction
	handlerMutex             sync.Mutex
}

// Connect makes the network connection to the Ethereum network.  Note: for
// other things to work correctly the configuration will need to reference a
// websocket, "ws://", or local IPC connection.
func Connect(cfg Config) (chain.Handle, error) {
	client, err := ethclient.Dial(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf(
			"error Connecting to Geth Server: %s [%v]",
			cfg.URL,
			err,
		)
	}

	clientws, err := rpc.Dial(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf(
			"error Connecting to Geth Server: %s [%v]",
			cfg.URL,
			err,
		)
	}

	clientrpc, err := rpc.Dial(cfg.URLRPC)
	if err != nil {
		return nil, fmt.Errorf(
			"error Connecting to Geth Server: %s [%v]",
			cfg.URL,
			err,
		)
	}

	pv := &ethereumChain{
		config:    cfg,
		client:    client,
		clientRPC: clientrpc,
		clientWS:  clientws,
	}

	keepRandomBeaconContract, err := newKeepRandomBeacon(pv)
	if err != nil {
		return nil, fmt.Errorf(
			"error attaching to KeepRandomBeacon contract: [%v]",
			err,
		)
	}
	pv.keepRandomBeaconContract = keepRandomBeaconContract

	keepGroupContract, err := newKeepGroup(pv)
	if err != nil {
		return nil, fmt.Errorf("error attaching to KeepGroup contract: [%v]", err)
	}
	pv.keepGroupContract = keepGroupContract

	return pv, nil
}
