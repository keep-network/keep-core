package ethereum

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/chain/ethereum/ethutil"
	"github.com/keep-network/keep-core/pkg/chain/gen/contract"
)

type ethereumChain struct {
	config                           Config
	client                           *ethclient.Client
	clientRPC                        *rpc.Client
	clientWS                         *rpc.Client
	requestID                        *big.Int
	keepRandomBeaconOperatorContract *contract.KeepRandomBeaconOperator
	keepRandomBeaconServiceContract  *contract.KeepRandomBeaconService
	stakingContract                  *contract.StakingProxy
	accountKey                       *keystore.Key

	// transactionMutex allows interested parties to forcibly serialize
	// transaction submission.
	//
	// When transactions are submitted, they require a valid nonce. The nonce is
	// equal to the count of transactions the account has submitted so far, and
	// for a transaction to be accepted it should be monotonically greater than
	// any previous submitted transaction. To do this, transaction submission
	// asks the Ethereum client it is connected to for the next pending nonce,
	// and uses that value for the transaction. Unfortunately, if multiple
	// transactions are submitted in short order, they may all get the same
	// nonce. Serializing submission ensures that each nonce is requested after
	// a previous transaction has been submitted.
	transactionMutex *sync.Mutex
}

// Connect makes the network connection to the Ethereum network.  Note: for
// other things to work correctly the configuration will need to reference a
// websocket, "ws://", or local IPC connection.
func Connect(config Config) (chain.Handle, error) {
	client, err := ethclient.Dial(config.URL)
	if err != nil {
		return nil, fmt.Errorf(
			"error Connecting to Geth Server: %s [%v]",
			config.URL,
			err,
		)
	}

	clientws, err := rpc.Dial(config.URL)
	if err != nil {
		return nil, fmt.Errorf(
			"error Connecting to Geth Server: %s [%v]",
			config.URL,
			err,
		)
	}

	clientrpc, err := rpc.Dial(config.URLRPC)
	if err != nil {
		return nil, fmt.Errorf(
			"error Connecting to Geth Server: %s [%v]",
			config.URL,
			err,
		)
	}

	pv := &ethereumChain{
		config:           config,
		client:           client,
		clientRPC:        clientrpc,
		clientWS:         clientws,
		transactionMutex: &sync.Mutex{},
	}

	if pv.accountKey == nil {
		key, err := ethutil.DecryptKeyFile(
			config.Account.KeyFile,
			config.Account.KeyFilePassword,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to read KeyFile: %s: [%v]",
				config.Account.KeyFile,
				err,
			)
		}
		pv.accountKey = key
	}

	address, err := addressForContract(config, "KeepRandomBeaconService")
	if err != nil {
		return nil, fmt.Errorf("error resolving KeepRandomBeaconService contract: [%v]", err)
	}

	keepRandomBeaconServiceContract, err :=
		contract.NewKeepRandomBeaconService(
			*address,
			pv.accountKey,
			pv.client,
			pv.transactionMutex,
		)
	if err != nil {
		return nil, fmt.Errorf(
			"error attaching to KeepRandomBeaconService contract: [%v]",
			err,
		)
	}
	pv.keepRandomBeaconServiceContract = keepRandomBeaconServiceContract

	address, err = addressForContract(config, "KeepRandomBeaconOperator")
	if err != nil {
		return nil, fmt.Errorf("error resolving KeepRandomBeaconOperator contract: [%v]", err)
	}

	keepRandomBeaconOperatorContract, err :=
		contract.NewKeepRandomBeaconOperator(
			*address,
			pv.accountKey,
			pv.client,
			pv.transactionMutex,
		)
	if err != nil {
		return nil, fmt.Errorf("error attaching to KeepRandomBeaconOperator contract: [%v]", err)
	}
	pv.keepRandomBeaconOperatorContract = keepRandomBeaconOperatorContract

	address, err = addressForContract(config, "StakingProxy")
	if err != nil {
		return nil, fmt.Errorf("error resolving StakingProxy contract: [%v]", err)
	}

	stakingContract, err :=
		contract.NewStakingProxy(
			*address,
			pv.accountKey,
			pv.client,
			pv.transactionMutex,
		)
	if err != nil {
		return nil, fmt.Errorf("error attaching to TokenStaking contract: [%v]", err)
	}
	pv.stakingContract = stakingContract

	return pv, nil
}

func addressForContract(config Config, contractName string) (*common.Address, error) {
	addressString, exists := config.ContractAddresses[contractName]
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
