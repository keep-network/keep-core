package ethereum

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/chain/ethereum/ethutil"
	"github.com/keep-network/keep-core/pkg/chain/gen/contract"
)

type ethereumChain struct {
	config                   Config
	client                   *ethclient.Client
	clientRPC                *rpc.Client
	clientWS                 *rpc.Client
	requestID                *big.Int
	keepGroupContract        *contract.KeepGroup
	keepRandomBeaconContract *contract.KeepRandomBeacon
	stakingContract          *contract.StakingProxy
	accountKey               *keystore.Key
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

	if pv.accountKey == nil {
		key, err := ethutil.DecryptKeyFile(
			cfg.Account.KeyFile,
			cfg.Account.KeyFilePassword,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to read KeyFile: %s: [%v]",
				cfg.Account.KeyFile,
				err,
			)
		}
		pv.accountKey = key
	}

	address, err := addressForContract(cfg, "KeepRandomBeacon")
	if err != nil {
		return nil, fmt.Errorf("error resolving KeepRandomBeacon contract: [%v]", err)
	}

	keepRandomBeaconContract, err :=
		contract.NewKeepRandomBeacon(
			*address,
			pv.accountKey,
			pv.client,
		)
	if err != nil {
		return nil, fmt.Errorf(
			"error attaching to KeepRandomBeacon contract: [%v]",
			err,
		)
	}
	pv.keepRandomBeaconContract = keepRandomBeaconContract

	address, err = addressForContract(cfg, "KeepGroup")
	if err != nil {
		return nil, fmt.Errorf("error resolving KeepGroup contract: [%v]", err)
	}

	keepGroupContract, err :=
		contract.NewKeepGroup(
			*address,
			pv.accountKey,
			pv.client,
		)
	if err != nil {
		return nil, fmt.Errorf("error attaching to KeepGroup contract: [%v]", err)
	}
	pv.keepGroupContract = keepGroupContract

	address, err = addressForContract(cfg, "Staking")
	if err != nil {
		return nil, fmt.Errorf("error resolving TokenStaking contract: [%v]", err)
	}

	stakingContract, err :=
		contract.NewStakingProxy(
			*address,
			pv.accountKey,
			pv.client,
		)
	if err != nil {
		return nil, fmt.Errorf("error attaching to TokenStaking contract: [%v]", err)
	}
	pv.stakingContract = stakingContract

	return pv, nil
}

func addressForContract(cfg Config, contractName string) (*common.Address, error) {
	addressString, exists := cfg.ContractAddresses["KeepRandomBeacon"]
	if !exists {
		return nil, fmt.Errorf(
			"no address information for [%v] in configuration",
			contractName,
		)
	}

	if !common.IsHexAddress(addressString) {
		return nil, fmt.Errorf(
			"configured address [%v] for contract [%v] is not valid hex address",
			addressString,
			contractName,
		)
	}

	address := common.HexToAddress(addressString)
	return &address, nil
}
