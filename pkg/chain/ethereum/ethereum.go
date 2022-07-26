package ethereum

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"github.com/keep-network/keep-core/pkg/operator"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/ipfs/go-log"
	"github.com/keep-network/keep-common/pkg/chain/ethereum"
	"github.com/keep-network/keep-common/pkg/chain/ethereum/ethutil"
	"github.com/keep-network/keep-common/pkg/rate"
	"github.com/keep-network/keep-core/pkg/chain/ethereum/threshold/gen/contract"
)

// Definitions of contract names.
const (
	TokenStakingContractName = "TokenStaking"
)

var logger = log.Logger("keep-chain-ethereum")

// Chain represents a base, non-application-specific chain handle. It
// provides the implementation of generic features like balance monitor,
// block counter and similar.
type Chain struct {
	key     *keystore.Key
	client  ethutil.EthereumClient
	chainID *big.Int

	blockCounter *ethereum.BlockCounter
	nonceManager *ethereum.NonceManager
	miningWaiter *ethutil.MiningWaiter

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

	tokenStaking *contract.TokenStaking
}

// Connect creates Random Beacon and TBTC Ethereum chain handles.
func Connect(
	ctx context.Context,
	config ethereum.Config,
) (*BeaconChain, *TbtcChain, error) {
	client, err := ethclient.Dial(config.URL)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"error Connecting to Ethereum Server: %s [%v]",
			config.URL,
			err,
		)
	}

	baseChain, err := newChain(ctx, config, client)
	if err != nil {
		return nil, nil, fmt.Errorf("could not create base chain handle: [%v]", err)
	}

	beaconChain, err := newBeaconChain(config, baseChain)
	if err != nil {
		return nil, nil, fmt.Errorf("could not create beacon chain handle: [%v]", err)
	}

	tbtcChain, err := newTbtcChain(config, baseChain)
	if err != nil {
		return nil, nil, fmt.Errorf("could not create TBTC chain handle: [%v]", err)
	}

	return beaconChain, tbtcChain, nil
}

// newChain construct a new instance of the Ethereum chain handle.
func newChain(
	ctx context.Context,
	config ethereum.Config,
	client *ethclient.Client,
) (*Chain, error) {
	chainID, err := client.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to resolve Ethereum chain id: [%v]",
			err,
		)
	}

	key, err := decryptKey(config)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to decrypt Ethereum key: [%v]",
			err,
		)
	}

	clientWithAddons := wrapClientAddons(config, client)

	blockCounter, err := ethutil.NewBlockCounter(clientWithAddons)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create Ethereum blockcounter: [%v]",
			err,
		)
	}

	nonceManager := ethutil.NewNonceManager(
		clientWithAddons,
		key.Address,
	)

	miningWaiter := ethutil.NewMiningWaiter(clientWithAddons, config)

	transactionMutex := &sync.Mutex{}

	// TODO: Consider adding the balance monitoring.

	tokenStakingAddress, err := config.ContractAddress(TokenStakingContractName)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to resolve %s contract address: [%v]",
			TokenStakingContractName,
			err,
		)
	}

	tokenStaking, err :=
		contract.NewTokenStaking(
			tokenStakingAddress,
			chainID,
			key,
			client,
			nonceManager,
			miningWaiter,
			blockCounter,
			transactionMutex,
		)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to attach to TokenStaking contract: [%v]",
			err,
		)
	}

	return &Chain{
		key:              key,
		client:           clientWithAddons,
		chainID:          chainID,
		blockCounter:     blockCounter,
		nonceManager:     nonceManager,
		miningWaiter:     miningWaiter,
		transactionMutex: transactionMutex,
		tokenStaking:     tokenStaking,
	}, nil
}

// OperatorKeyPair returns the key pair of the operator assigned to this
// chain handle.
func (c *Chain) OperatorKeyPair() (
	*operator.PrivateKey,
	*operator.PublicKey,
	error,
) {
	privateKey, publicKey, err := ChainPrivateKeyToOperatorKeyPair(
		c.key.PrivateKey,
	)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"cannot convert chain private key to operator key pair: [%v]",
			err,
		)
	}

	return privateKey, publicKey, nil
}

// wrapClientAddons wraps the client instance with add-ons like logging, rate
// limiting and so on.
func wrapClientAddons(
	config ethereum.Config,
	client ethutil.EthereumClient,
) ethutil.EthereumClient {
	loggingClient := ethutil.WrapCallLogging(logger, client)

	if config.RequestsPerSecondLimit > 0 || config.ConcurrencyLimit > 0 {
		logger.Infof(
			"enabled ethereum client request rate limiter; "+
				"rps limit [%v]; "+
				"concurrency limit [%v]",
			config.RequestsPerSecondLimit,
			config.ConcurrencyLimit,
		)

		return ethutil.WrapRateLimiting(
			loggingClient,
			&rate.LimiterConfig{
				RequestsPerSecondLimit: config.RequestsPerSecondLimit,
				ConcurrencyLimit:       config.ConcurrencyLimit,
			},
		)
	}

	return loggingClient
}

// decryptKey decrypts the chain key pointed by the config.
func decryptKey(config ethereum.Config) (*keystore.Key, error) {
	return ethutil.DecryptKeyFile(
		config.Account.KeyFile,
		config.Account.KeyFilePassword,
	)
}
