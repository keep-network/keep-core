package ethereum

import (
	"fmt"
	"math/big"
	"sync"
	"time"

	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/keep-network/keep-common/pkg/chain/ethereum"
	"github.com/keep-network/keep-common/pkg/chain/ethereum/blockcounter"
	"github.com/keep-network/keep-common/pkg/chain/ethereum/ethutil"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/chain/gen/contract"
)

var (
	// DefaultMiningCheckInterval is the default interval in which transaction
	// mining status is checked. If the transaction is not mined within this
	// time, the gas price is increased and transaction is resubmitted.
	// This value can be overwritten in the configuration file.
	DefaultMiningCheckInterval = 60 * time.Second

	// DefaultMaxGasPrice specifies the default maximum gas price the client is
	// willing to pay for the transaction to be mined. The offered transaction
	// gas price can not be higher than the max gas price value. If the maximum
	// allowed gas price is reached, no further resubmission attempts are
	// performed. This value can be overwritten in the configuration file.
	DefaultMaxGasPrice = big.NewInt(70000000000) // 70 Gwei
)

type ethereumChain struct {
	config                           ethereum.Config
	client                           bind.ContractBackend
	clientRPC                        *rpc.Client
	clientWS                         *rpc.Client
	keepRandomBeaconOperatorContract *contract.KeepRandomBeaconOperator
	stakingContract                  *contract.TokenStaking
	accountKey                       *keystore.Key
	blockCounter                     *blockcounter.EthereumBlockCounter
	chainConfig                      *relaychain.Config

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

type ethereumUtilityChain struct {
	ethereumChain

	keepRandomBeaconServiceContract *contract.KeepRandomBeaconService
}

func connect(config ethereum.Config) (*ethereumChain, error) {
	client, clientWS, clientRPC, err := ethutil.ConnectClients(config.URL, config.URLRPC)
	if err != nil {
		return nil, fmt.Errorf(
			"error connecting to Ethereum server: %s [%v]",
			config.URL,
			err,
		)
	}

	return connectWithClient(config, client, clientWS, clientRPC)
}

func connectWithClient(
	config ethereum.Config,
	client *ethclient.Client,
	clientWS *rpc.Client,
	clientRPC *rpc.Client,
) (*ethereumChain, error) {

	blockCounter, err := blockcounter.CreateBlockCounter(client)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create Ethereum blockcounter: [%v]",
			err,
		)
	}

	pv := &ethereumChain{
		config:           config,
		client:           ethutil.WrapCallLogging(logger, client),
		clientRPC:        clientRPC,
		clientWS:         clientWS,
		transactionMutex: &sync.Mutex{},
		blockCounter:     blockCounter,
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

	checkInterval := DefaultMiningCheckInterval
	maxGasPrice := DefaultMaxGasPrice
	if config.MiningCheckInterval != 0 {
		checkInterval = time.Duration(config.MiningCheckInterval) * time.Second
	}
	if config.MaxGasPrice != 0 {
		maxGasPrice = new(big.Int).SetUint64(config.MaxGasPrice)
	}

	logger.Infof("using [%v] mining check interval", checkInterval)
	logger.Infof("using [%v] wei max gas price", maxGasPrice)
	miningWaiter := ethutil.NewMiningWaiter(client, checkInterval, maxGasPrice)

	address, err := addressForContract(config, "KeepRandomBeaconOperator")
	if err != nil {
		return nil, fmt.Errorf("error resolving KeepRandomBeaconOperator contract: [%v]", err)
	}

	nonceManager := ethutil.NewNonceManager(
		pv.accountKey.Address,
		pv.client,
	)

	keepRandomBeaconOperatorContract, err :=
		contract.NewKeepRandomBeaconOperator(
			*address,
			pv.accountKey,
			pv.client,
			nonceManager,
			miningWaiter,
			pv.transactionMutex,
		)
	if err != nil {
		return nil, fmt.Errorf("error attaching to KeepRandomBeaconOperator contract: [%v]", err)
	}
	pv.keepRandomBeaconOperatorContract = keepRandomBeaconOperatorContract

	address, err = addressForContract(config, "TokenStaking")
	if err != nil {
		return nil, fmt.Errorf("error resolving TokenStaking contract: [%v]", err)
	}

	stakingContract, err :=
		contract.NewTokenStaking(
			*address,
			pv.accountKey,
			pv.client,
			nonceManager,
			miningWaiter,
			pv.transactionMutex,
		)
	if err != nil {
		return nil, fmt.Errorf("error attaching to TokenStaking contract: [%v]", err)
	}
	pv.stakingContract = stakingContract

	chainConfig, err := fetchChainConfig(pv)
	if err != nil {
		return nil, fmt.Errorf("could not fetch chain config: [%v]", err)
	}
	pv.chainConfig = chainConfig

	return pv, nil
}

// ConnectUtility makes the network connection to the Ethereum network and
// returns a utility handle to the chain interface with additional methods for
// non- standard client interactions. Note: for other things to work correctly
// the configuration will need to reference a websocket, "ws://", or local IPC
// connection.
func ConnectUtility(config ethereum.Config) (chain.Utility, error) {
	client, clientWS, clientRPC, err := ethutil.ConnectClients(config.URL, config.URLRPC)
	if err != nil {
		return nil, fmt.Errorf(
			"error connecting to Ethereum server: %s [%v]",
			config.URL,
			err,
		)
	}

	base, err := connectWithClient(config, client, clientWS, clientRPC)
	if err != nil {
		return nil, err
	}

	checkInterval := DefaultMiningCheckInterval
	maxGasPrice := DefaultMaxGasPrice
	if config.MiningCheckInterval != 0 {
		checkInterval = time.Duration(config.MiningCheckInterval) * time.Second
	}
	if config.MaxGasPrice != 0 {
		maxGasPrice = new(big.Int).SetUint64(config.MaxGasPrice)
	}

	miningWaiter := ethutil.NewMiningWaiter(client, checkInterval, maxGasPrice)

	address, err := addressForContract(config, "KeepRandomBeaconService")
	if err != nil {
		return nil, fmt.Errorf("error resolving KeepRandomBeaconService contract: [%v]", err)
	}

	nonceManager := ethutil.NewNonceManager(
		base.accountKey.Address,
		base.client,
	)

	keepRandomBeaconServiceContract, err :=
		contract.NewKeepRandomBeaconService(
			*address,
			base.accountKey,
			base.client,
			nonceManager,
			miningWaiter,
			base.transactionMutex,
		)
	if err != nil {
		return nil, fmt.Errorf("error attaching to KeepRandomBeaconService contract: [%v]", err)
	}

	return &ethereumUtilityChain{
		*base,
		keepRandomBeaconServiceContract,
	}, nil
}

// Connect makes the network connection to the Ethereum network and returns a
// standard handle to the chain interface. Note: for other things to work
// correctly the configuration will need to reference a websocket, "ws://", or
// local IPC connection.
func Connect(config ethereum.Config) (chain.Handle, error) {
	return connect(config)
}

func addressForContract(config ethereum.Config, contractName string) (*common.Address, error) {
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

// BlockCounter creates a BlockCounter that uses the block number in ethereum.
func (ec *ethereumChain) BlockCounter() (chain.BlockCounter, error) {
	return ec.blockCounter, nil
}

func fetchChainConfig(ec *ethereumChain) (*relaychain.Config, error) {
	logger.Infof("fetching relay chain config")

	groupSize, err := ec.keepRandomBeaconOperatorContract.GroupSize()
	if err != nil {
		return nil, fmt.Errorf("error calling GroupSize: [%v]", err)
	}

	threshold, err := ec.keepRandomBeaconOperatorContract.GroupThreshold()
	if err != nil {
		return nil, fmt.Errorf("error calling GroupThreshold: [%v]", err)
	}

	ticketSubmissionTimeout, err :=
		ec.keepRandomBeaconOperatorContract.TicketSubmissionTimeout()
	if err != nil {
		return nil, fmt.Errorf(
			"error calling TicketSubmissionTimeout: [%v]",
			err,
		)
	}

	resultPublicationBlockStep, err := ec.keepRandomBeaconOperatorContract.ResultPublicationBlockStep()
	if err != nil {
		return nil, fmt.Errorf(
			"error calling ResultPublicationBlockStep: [%v]",
			err,
		)
	}

	relayEntryTimeout, err := ec.keepRandomBeaconOperatorContract.RelayEntryTimeout()
	if err != nil {
		return nil, fmt.Errorf("error calling RelayEntryTimeout: [%v]", err)
	}

	return &relaychain.Config{
		GroupSize:                  int(groupSize.Int64()),
		HonestThreshold:            int(threshold.Int64()),
		TicketSubmissionTimeout:    ticketSubmissionTimeout.Uint64(),
		ResultPublicationBlockStep: resultPublicationBlockStep.Uint64(),
		RelayEntryTimeout:          relayEntryTimeout.Uint64(),
	}, nil
}
