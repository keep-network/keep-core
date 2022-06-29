package ethereum

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"github.com/keep-network/keep-common/pkg/rate"

	"github.com/keep-network/keep-common/pkg/chain/ethlike"

	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/keep-network/keep-common/pkg/chain/ethereum"
	"github.com/keep-network/keep-common/pkg/chain/ethereum/ethutil"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/chain/gen/contract"
)

// Definitions of contract names.
const (
	KeepRandomBeaconOperatorContractName = "KeepRandomBeaconOperator"
	TokenStakingContractName             = "TokenStaking"
	KeepRandomBeaconServiceContractName  = "KeepRandomBeaconService"
)

type ethereumChain struct {
	config                           ethereum.Config
	accountKey                       *keystore.Key
	client                           ethutil.EthereumClient
	clientRPC                        *rpc.Client
	clientWS                         *rpc.Client
	chainID                          *big.Int
	keepRandomBeaconOperatorContract *contract.KeepRandomBeaconOperator
	stakingContract                  *contract.TokenStaking
	blockCounter                     *ethlike.BlockCounter
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

func connect(
	ctx context.Context,
	config ethereum.Config,
) (*ethereumChain, error) {
	client, clientWS, clientRPC, err := ethutil.ConnectClients(config.URL, config.URLRPC)
	if err != nil {
		return nil, fmt.Errorf(
			"error connecting to Ethereum server: %s [%v]",
			config.URL,
			err,
		)
	}

	return connectWithClient(ctx, config, client, clientWS, clientRPC)
}

func connectWithClient(
	ctx context.Context,
	config ethereum.Config,
	client *ethclient.Client,
	clientWS *rpc.Client,
	clientRPC *rpc.Client,
) (*ethereumChain, error) {
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf(
			"failed to resolve Ethereum chain id: [%v]",
			err,
		)
	}

	ec := &ethereumChain{
		config:           config,
		client:           addClientWrappers(config, client),
		clientRPC:        clientRPC,
		clientWS:         clientWS,
		chainID:          chainID,
		transactionMutex: &sync.Mutex{},
	}

	blockCounter, err := ethutil.NewBlockCounter(ec.client)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create Ethereum blockcounter: [%v]",
			err,
		)
	}
	ec.blockCounter = blockCounter

	if ec.accountKey == nil {
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
		ec.accountKey = key
	}

	miningWaiter := ethutil.NewMiningWaiter(ec.client, config)

	address, err := config.ContractAddress(KeepRandomBeaconOperatorContractName)
	if err != nil {
		return nil, fmt.Errorf("error resolving KeepRandomBeaconOperator contract: [%v]", err)
	}

	nonceManager := ethutil.NewNonceManager(
		ec.client,
		ec.accountKey.Address,
	)

	keepRandomBeaconOperatorContract, err :=
		contract.NewKeepRandomBeaconOperator(
			address,
			ec.chainID,
			ec.accountKey,
			ec.client,
			nonceManager,
			miningWaiter,
			blockCounter,
			ec.transactionMutex,
		)
	if err != nil {
		return nil, fmt.Errorf("error attaching to KeepRandomBeaconOperator contract: [%v]", err)
	}
	ec.keepRandomBeaconOperatorContract = keepRandomBeaconOperatorContract

	address, err = config.ContractAddress(TokenStakingContractName)
	if err != nil {
		return nil, fmt.Errorf("error resolving TokenStaking contract: [%v]", err)
	}

	stakingContract, err :=
		contract.NewTokenStaking(
			address,
			ec.chainID,
			ec.accountKey,
			ec.client,
			nonceManager,
			miningWaiter,
			blockCounter,
			ec.transactionMutex,
		)
	if err != nil {
		return nil, fmt.Errorf("error attaching to TokenStaking contract: [%v]", err)
	}
	ec.stakingContract = stakingContract

	chainConfig, err := fetchChainConfig(ec)
	if err != nil {
		return nil, fmt.Errorf("could not fetch chain config: [%v]", err)
	}
	ec.chainConfig = chainConfig

	ec.initializeBalanceMonitoring(ctx)

	return ec, nil
}

func addClientWrappers(
	config ethereum.Config,
	backend ethutil.EthereumClient,
) ethutil.EthereumClient {
	loggingBackend := ethutil.WrapCallLogging(logger, backend)

	if config.RequestsPerSecondLimit > 0 || config.ConcurrencyLimit > 0 {
		logger.Infof(
			"enabled ethereum client request rate limiter; "+
				"rps limit [%v]; "+
				"concurrency limit [%v]",
			config.RequestsPerSecondLimit,
			config.ConcurrencyLimit,
		)

		return ethutil.WrapRateLimiting(
			loggingBackend,
			&rate.LimiterConfig{
				RequestsPerSecondLimit: config.RequestsPerSecondLimit,
				ConcurrencyLimit:       config.ConcurrencyLimit,
			},
		)
	}

	return loggingBackend
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

	base, err := connectWithClient(
		context.Background(),
		config,
		client,
		clientWS,
		clientRPC,
	)
	if err != nil {
		return nil, err
	}

	blockCounter, err := ethutil.NewBlockCounter(client)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create Ethereum blockcounter: [%v]",
			err,
		)
	}

	miningWaiter := ethutil.NewMiningWaiter(client, config)

	address, err := config.ContractAddress(KeepRandomBeaconServiceContractName)
	if err != nil {
		return nil, fmt.Errorf("error resolving KeepRandomBeaconService contract: [%v]", err)
	}

	nonceManager := ethutil.NewNonceManager(
		client,
		base.accountKey.Address,
	)

	keepRandomBeaconServiceContract, err :=
		contract.NewKeepRandomBeaconService(
			address,
			base.chainID,
			base.accountKey,
			base.client,
			nonceManager,
			miningWaiter,
			blockCounter,
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
func Connect(
	ctx context.Context,
	config ethereum.Config,
) (chain.Handle, error) {
	return connect(ctx, config)
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
		ResultPublicationBlockStep: resultPublicationBlockStep.Uint64(),
		RelayEntryTimeout:          relayEntryTimeout.Uint64(),
	}, nil
}
