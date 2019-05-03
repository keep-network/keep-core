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
	config                           Config
	client                           *ethclient.Client
	clientRPC                        *rpc.Client
	clientWS                         *rpc.Client
	requestID                        *big.Int
	keepRandomBeaconBackendContract  *keepRandomBeaconBackend
	keepRandomBeaconFrontendContract *KeepRandomBeaconFrontend
	stakingContract                  *staking
	accountKey                       *keystore.Key
}

// Connect makes the network connection to the Ethereum network. Note: for
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

	keepRandomBeaconFrontendContract, err := newKeepRandomBeaconFrontend(pv)
	if err != nil {
		return nil, fmt.Errorf(
			"error attaching to KeepRandomBeaconFrontend contract: [%v]",
			err,
		)
	}
	pv.keepRandomBeaconFrontendContract = keepRandomBeaconFrontendContract

	keepRandomBeaconBackendContract, err := newKeepRandomBeaconBackend(pv)
	if err != nil {
		return nil, fmt.Errorf("error attaching to keepRandomBeaconBackend contract: [%v]", err)
	}
	pv.keepRandomBeaconBackendContract = keepRandomBeaconBackendContract

	stakingContract, err := newStaking(pv)
	if err != nil {
		return nil, fmt.Errorf("error attaching to TokenStaking contract: [%v]", err)
	}
	pv.stakingContract = stakingContract

	return pv, nil
}
