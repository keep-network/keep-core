package ethereum

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/keystore"
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
	accountKey               *keystore.Key
	blockCounter             chain.BlockCounter
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

	blockCounter, err := pv.BlockCounter()
	if err != nil {
		return nil, fmt.Errorf(
			"error setting up BlockCounter: [%v]",
			err,
		)
	}
	pv.blockCounter = blockCounter

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
