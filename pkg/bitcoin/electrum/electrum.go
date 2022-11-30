package electrum

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/checksum0/go-electrum/electrum"
	"github.com/ipfs/go-log"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"

	"github.com/keep-network/keep-common/pkg/wrappers"
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/internal/byteutils"
)

var (
	supportedProtocolVersions = []string{"1.4"}
	logger                    = log.Logger("keep-electrum")
)

// Connection is a handle for interactions with Electrum server.
type Connection struct {
	ctx                 context.Context
	client              *electrum.Client
	requestRetryTimeout time.Duration //TODO: REMOVE
	config              Config
}

// Connect initializes handle with provided Config.
func Connect(parentCtx context.Context, config Config) (bitcoin.Chain, error) {
	c := &Connection{
		ctx:                 parentCtx,
		requestRetryTimeout: config.RequestRetryTimeout,
		config:              config,
		clientMutex:         &sync.RWMutex{},
	}

	if err := c.electrumConnect(); err != nil {
		return nil, fmt.Errorf("failed to initialize electrum client: [%w]", err)
	}

	// Get the server's details.
	type Server struct {
		version  string
		protocol string
	}
	server, err := requestWithRetry(
		c,
		func(ctx context.Context, client *electrum.Client) (*Server, error) {
			serverVersion, protocolVersion, err := client.ServerVersion(ctx)
			if err != nil {
				return nil, err
			}
			return &Server{serverVersion, protocolVersion}, nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get server version: [%w]", err)
	}

	logger.Infof(
		"connected to electrum server [version: [%s], protocol: [%s]]",
		server.version,
		server.protocol,
	)

	// Log a warning if connected to a server running an unsupported protocol version.
	if !slices.Contains(supportedProtocolVersions, server.protocol) {
		logger.Warnf(
			"electrum server [%s] runs an unsupported protocol version: [%s]; expected one of: [%s]",
			config.URL,
			server.protocol,
			strings.Join(supportedProtocolVersions, ","),
		)
	}

	// Keep the connection alive and check the connection health.
	go c.keepAlive()

	// TODO: Add reconnects on lost connection.

	return c, nil
}

// GetTransaction gets the transaction with the given transaction hash.
// If the transaction with the given hash was not found on the chain,
// this function returns an error.
func (c *Connection) GetTransaction(
	transactionHash bitcoin.Hash,
) (*bitcoin.Transaction, error) {
	txID := transactionHash.Hex(bitcoin.ReversedByteOrder)

	rawTransaction, err := requestWithRetry(
		c,
		func(ctx context.Context, client *electrum.Client) (string, error) {
			// We cannot use `GetTransaction` to get the the transaction details
			// as Esplora/Electrs doesn't support verbose transactions.
			// See: https://github.com/Blockstream/electrs/pull/36
			return client.GetRawTransaction(ctx, txID)
		})
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get raw transaction with ID [%s]: [%w]",
			txID,
			err,
		)
	}

	result, err := convertRawTransaction(rawTransaction)
	if err != nil {
		return nil, fmt.Errorf("failed to convert transaction: [%w]", err)
	}

	return result, nil
}

// GetTransactionConfirmations gets the number of confirmations for the
// transaction with the given transaction hash. If the transaction with the
// given hash was not found on the chain, this function returns an error.
func (c *Connection) GetTransactionConfirmations(
	transactionHash bitcoin.Hash,
) (uint, error) {
	txID := transactionHash.Hex(bitcoin.ReversedByteOrder)

	txLogger := logger.With(
		zap.String("txID", txID),
	)

	rawTransaction, err := requestWithRetry(
		c,
		func(ctx context.Context, client *electrum.Client) (string, error) {
			// We cannot use `GetTransaction` to get the the transaction details
			// as Esplora/Electrs doesn't support verbose transactions.
			// See: https://github.com/Blockstream/electrs/pull/36
			return client.GetRawTransaction(ctx, txID)
		})
	if err != nil {
		return 0,
			fmt.Errorf(
				"failed to get raw transaction with ID [%s]: [%w]",
				txID,
				err,
			)
	}

	tx, err := decodeTransaction(rawTransaction)
	if err != nil {
		return 0, fmt.Errorf(
			"failed to decode the transaction [%s]: [%w]",
			rawTransaction,
			err,
		)
	}

	// As a workaround for the problem described in https://github.com/Blockstream/electrs/pull/36
	// we need to calculate the number of confirmations based on the latest
	// block height and block height of the transaction.
	// Electrum protocol doesn't expose a function to get the transaction's block
	// height (other that the `GetTransaction` that is unsupported by Esplora/Electrs).
	// To get the block height of the transaction we query the history of transactions
	// for the output script hash, as the history contains the transaction's block
	// height.

	// Initialize txBlockHeigh with minimum int32 value to identify a problem when
	// a block height was not found in a history of any of the script hashes.
	//
	// The history is expected to return a block height for confirmed transaction.
	// If a transaction is unconfirmed (is still in the mempool) the height will
	// have a value of `0` or `-1`.
	txBlockHeight := int32(math.MinInt32)
txOutLoop:
	for _, txOut := range tx.TxOut {
		script := txOut.PkScript
		scriptHash := sha256.Sum256(script)
		reversedScriptHash := byteutils.Reverse(scriptHash[:])
		reversedScriptHashString := hex.EncodeToString(reversedScriptHash)

		scriptHashHistory, err := requestWithRetry(
			c,
			func(
				ctx context.Context,
				client *electrum.Client,
			) ([]*electrum.GetMempoolResult, error) {
				return client.GetHistory(ctx, reversedScriptHashString)
			})
		if err != nil {
			// Don't return an error, but continue to the next TxOut entry.
			txLogger.Errorf("failed to get history for script hash: [%v]", err)
			continue txOutLoop
		}

		for _, transaction := range scriptHashHistory {
			if transaction.Hash == txID {
				txBlockHeight = transaction.Height
				break txOutLoop
			}
		}
	}

	// History querying didn't come up with the transaction's block height. Return
	// an error.
	if txBlockHeight == math.MinInt32 {
		return 0, fmt.Errorf(
			"failed to find the transaction block height in script hashes' histories",
		)
	}

	// If the block height is greater than `0` the transaction is confirmed.
	if txBlockHeight > 0 {
		latestBlockHeight, err := c.GetLatestBlockHeight()
		if err != nil {
			return 0, fmt.Errorf(
				"failed to get the latest block height: [%w]",
				err,
			)
		}

		if latestBlockHeight >= uint(txBlockHeight) {
			// Add `1` to the calculated difference as if the transaction block
			// height equals the latest block height the transaction is already
			// confirmed, so it has one confirmation.
			return latestBlockHeight - uint(txBlockHeight) + 1, nil
		}
	}

	return 0, nil
}

// BroadcastTransaction broadcasts the given transaction over the
// network of the Bitcoin chain nodes. If the broadcast action could not be
// done, this function returns an error. This function does not give any
// guarantees regarding transaction mining. The transaction may be mined or
// rejected eventually.
func (c *Connection) BroadcastTransaction(
	transaction *bitcoin.Transaction,
) error {
	rawTx := transaction.Serialize()

	rawTxLogger := logger.With(
		zap.String("rawTx", hex.EncodeToString(rawTx)),
	)
	rawTxLogger.Debugf("broadcasting transaction")

	var response string

	response, err := requestWithRetry(
		c,
		func(ctx context.Context, client *electrum.Client) (string, error) {
			return client.BroadcastTransaction(ctx, string(rawTx))
		})
	if err != nil {
		return fmt.Errorf("failed to broadcast the transaction: [%w]", err)
	}

	rawTxLogger.Infof("transaction broadcast successful: [%s]", response)

	return nil
}

// GetLatestBlockHeight gets the height of the latest block (tip). If the
// latest block was not determined, this function returns an error.
func (c *Connection) GetLatestBlockHeight() (uint, error) {
	blockHeight, err := requestWithRetry(
		c,
		func(ctx context.Context, client *electrum.Client) (int32, error) {
			headersChan, err := client.SubscribeHeaders(ctx)
			if err != nil {
				return 0, fmt.Errorf("failed to get the blocks tip height: [%w]", err)
			}
			tip := <-headersChan
			return tip.Height, nil
		})
	if err != nil {
		return 0, fmt.Errorf("failed to subscribe for headers: [%w]", err)
	}

	if blockHeight > 0 {
		return uint(blockHeight), nil
	}

	return 0, nil
}

// GetBlockHeader gets the block header for the given block height. If the
// block with the given height was not found on the chain, this function
// returns an error.
func (c *Connection) GetBlockHeader(
	blockHeight uint,
) (*bitcoin.BlockHeader, error) {
	getBlockHeaderResult, err := requestWithRetry(
		c,
		func(
			ctx context.Context,
			client *electrum.Client,
		) (*electrum.GetBlockHeaderResult, error) {
			return client.GetBlockHeader(ctx, uint32(blockHeight), 0)
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get block header: [%w]", err)
	}

	blockHeader, err := convertBlockHeader(getBlockHeaderResult)
	if err != nil {
		return nil, fmt.Errorf("failed to convert block header: %w", err)
	}

	return blockHeader, nil
}

func (c *Connection) electrumConnect() error {
	var client *electrum.Client
	var err error
	switch c.config.Protocol {
	case TCP:
		// TODO: Add retry to connection establishment.
		client, err = electrum.NewClientTCP(c.ctx, c.config.URL)
	case SSL:
		// TODO: Implement certificate verification to be able to disable the `InsecureSkipVerify: true` workaround.
		// #nosec G402 (TLS InsecureSkipVerify set true)
		tlsConfig := &tls.Config{InsecureSkipVerify: true}
		// TODO: Add retry to connection establishment.
		client, err = electrum.NewClientSSL(c.ctx, c.config.URL, tlsConfig)
	default:
		err = fmt.Errorf("unsupported protocol: [%s]", c.config.Protocol)

	}

	if err == nil {
		c.client = client
	}

	return err
}

func (c *Connection) keepAlive() {
	ticker := time.NewTicker(c.config.KeepAliveInterval)

	for {
		select {
		case <-ticker.C:
			_, err := requestWithRetry(
				c,
				func(ctx context.Context, client *electrum.Client) (interface{}, error) {
					return nil, client.Ping(ctx)
				},
			)
			if err != nil {
				logger.Errorf(
					"failed to ping the electrum server; "+
						"please verify health of the electrum server: [%v]",
					err,
				)
			}
		case <-c.ctx.Done():
			ticker.Stop()
			c.client.Shutdown()
			return
		}
	}
}

func requestWithRetry[K interface{}](
	c *Connection,
	requestFn func(ctx context.Context, client *electrum.Client) (K, error),
) (K, error) {
	var result K

	err := wrappers.DoWithDefaultRetry(
		c.requestRetryTimeout,
		func(ctx context.Context) error {
			requestCtx, requestCancel := context.WithTimeout(ctx, c.config.RequestTimeout)
			defer requestCancel()
			r, err := requestFn(requestCtx, c.client)

			if err != nil {
				return fmt.Errorf("request failed: [%w]", err)
			}

			result = r
			return nil
		})

	return result, err
}
