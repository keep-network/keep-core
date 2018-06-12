package ethereum

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/keep-network/keep-core/pkg/beacon/relay"

	"github.com/keep-network/keep-core/pkg/beacon"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/keep-network/keep-core/pkg/chain"
)

type ethereumChain struct {
	config                           Config
	client                           *ethclient.Client
	clientRPC                        *rpc.Client
	clientWS                         *rpc.Client
	requestID                        *big.Int
	tx                               *types.Transaction
	handlerMutex                     sync.Mutex
	groupPublicKeyFailureHandlers    []func(groupID string, errorMessage string)
	groupPublicKeySubmissionHandlers []func(groupID string, activationBlock *big.Int)
}

func (ec *ethereumChain) RandomBeacon() beacon.ChainInterface {
	return nil
}

func (ec *ethereumChain) ThresholdRelay() relay.ChainInterface {
	return nil
}

// Connect makes the network connection to the Ethereum network.  Note: for
// other things to work correctly the configuration will need to reference a
// websocket, "ws://", or local IPC connection.
func Connect(cfg Config) (chain.Handle, error) {
	client, err := ethclient.Dial(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("error Connecting to Geth Server: %s server %s",
			err, cfg.URL)
	}

	clientws, err := rpc.Dial(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("error Connecting to Geth Server: %s server %s",
			err, cfg.URL)
	}

	clientrpc, err := rpc.Dial(cfg.URLRPC)
	if err != nil {
		return nil, fmt.Errorf("error Connecting to Geth Server: %s server %s",
			err, cfg.URL)
	}

	pv := &ethereumChain{
		config:    cfg,
		client:    client,
		clientRPC: clientrpc,
		clientWS:  clientws,
	}

	return pv, nil
}
