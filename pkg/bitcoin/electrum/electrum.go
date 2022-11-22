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
	requestRetryTimeout time.Duration
}

// Connect initializes handle with provided Config.
func Connect(ctx context.Context, config Config) (bitcoin.Chain, error) {
	if config.KeepAliveInterval == 0 {
		config.KeepAliveInterval = DefaultKeepAliveInterval
	}
	if config.RequestRetryTimeout == 0 {
		config.RequestRetryTimeout = DefaultRequestRetryTimeout
	}

	var client *electrum.Client
	var err error
	switch config.Protocol {
	case TCP:
		// TODO: Add retry to connection establishment.
		client, err = electrum.NewClientTCP(ctx, config.URL)
	case SSL:
		// TODO: Implement certificate verification to be able to disable the `InsecureSkipVerify: true` workaround.
		// #nosec G402 (TLS InsecureSkipVerify set true)
		tlsConfig := &tls.Config{InsecureSkipVerify: true}
		// TODO: Add retry to connection establishment.
		client, err = electrum.NewClientSSL(ctx, config.URL, tlsConfig)
	default:
		err = fmt.Errorf("unsupported protocol: [%s]", config.Protocol)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to initialize electrum client: [%w]", err)
	}

	// Get the server's details.
	err = wrappers.DoWithDefaultRetry(
		config.RequestRetryTimeout,
		func(ctx context.Context) error {
			serverVersion, protocolVersion, err := client.ServerVersion(ctx)
			if err != nil {
				return fmt.Errorf("ServerVersion failed: [%w]", err)
			}
			logger.Infof(
				"connected to electrum server [version: [%s], protocol: [%s]]",
				serverVersion,
				protocolVersion,
			)

			// Log a warning if connected to a server running an unsupported protocol version.
			if !slices.Contains(supportedProtocolVersions, protocolVersion) {
				logger.Warnf(
					"electrum server [%s] runs an unsupported protocol version: [%s]; expected one of: [%s]",
					config.URL,
					protocolVersion,
					strings.Join(supportedProtocolVersions, ","),
				)
			}

			return nil
		})
	if err != nil {
		return nil, fmt.Errorf("failed to get server version: [%w]", err)
	}

	// Keep the connection alive and check the connection health.
	go func() {
		ticker := time.NewTicker(config.KeepAliveInterval)

		for {
			select {
			case <-ticker.C:
				err := wrappers.DoWithDefaultRetry(config.RequestRetryTimeout, client.Ping)
				if err != nil {
					logger.Errorf(
						"failed to ping the electrum server; "+
							"please verify health of the electrum server: [%v]",
						err,
					)
				}
			case <-ctx.Done():
				ticker.Stop()
				client.Shutdown()
				return
			}
		}
	}()

	// TODO: Add reconnects on lost connection.

	return &Connection{
		ctx:                 ctx,
		client:              client,
		requestRetryTimeout: config.RequestRetryTimeout,
	}, nil
}

// GetTransaction gets the transaction with the given transaction hash.
// If the transaction with the given hash was not found on the chain,
// this function returns an error.
func (c *Connection) GetTransaction(
	transactionHash bitcoin.Hash,
) (*bitcoin.Transaction, error) {
	txID := transactionHash.Hex(bitcoin.ReversedByteOrder)

	var rawTransaction string
	err := wrappers.DoWithDefaultRetry(c.requestRetryTimeout, func(ctx context.Context) error {
		// We cannot use `GetTransaction` to get the the transaction details
		// as Esplora/Electrs doesn't support verbose transactions.
		// See: https://github.com/Blockstream/electrs/pull/36
		rawTx, err := c.client.GetRawTransaction(c.ctx, txID)
		if err != nil {
			return fmt.Errorf("GetRawTransaction failed: [%w]", err)
		}

		rawTransaction = rawTx

		return nil
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

	var rawTransaction string
	err := wrappers.DoWithDefaultRetry(c.requestRetryTimeout, func(ctx context.Context) error {
		// We cannot use `GetTransaction` to get the the transaction details
		// as Esplora/Electrs doesn't support verbose transactions.
		// See: https://github.com/Blockstream/electrs/pull/36
		rawTx, err := c.client.GetRawTransaction(c.ctx, txID)
		if err != nil {
			return fmt.Errorf(
				"GetRawTransaction failed: [%w]",
				err,
			)
		}
		rawTransaction = rawTx
		return nil
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

		var scriptHashHistory []*electrum.GetMempoolResult
		err := wrappers.DoWithDefaultRetry(c.requestRetryTimeout, func(ctx context.Context) error {
			history, err := c.client.GetHistory(c.ctx, reversedScriptHashString)
			if err != nil {
				return fmt.Errorf("GetHistory failed: [%w]", err)
			}

			scriptHashHistory = history

			return nil
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
	err := wrappers.DoWithDefaultRetry(c.requestRetryTimeout, func(ctx context.Context) error {
		var err error
		response, err = c.client.BroadcastTransaction(c.ctx, string(rawTx))
		if err != nil {
			return fmt.Errorf("BroadcastTransaction failed: [%w]", err)
		}

		rawTxLogger.Infof("transaction broadcast successful: [%s]", response)

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to broadcast the transaction: [%w]", err)
	}

	return nil
}

// GetLatestBlockHeight gets the height of the latest block (tip). If the
// latest block was not determined, this function returns an error.
func (c *Connection) GetLatestBlockHeight() (uint, error) {
	var tip *electrum.SubscribeHeadersResult
	err := wrappers.DoWithDefaultRetry(c.requestRetryTimeout, func(ctx context.Context) error {
		headersChan, err := c.client.SubscribeHeaders(c.ctx)
		if err != nil {
			return fmt.Errorf("failed to get the blocks tip height: [%w]", err)
		}

		tip = <-headersChan
		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("failed to subscribe for headers: [%w]", err)
	}

	blockHeight := tip.Height

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
	var getBlockHeaderResult *electrum.GetBlockHeaderResult
	err := wrappers.DoWithDefaultRetry(c.requestRetryTimeout, func(ctx context.Context) error {
		result, err := c.client.GetBlockHeader(c.ctx, uint32(blockHeight), 0)
		if err != nil {
			return fmt.Errorf(
				"GetBlockHeader failed: [%w]",
				err,
			)
		}
		getBlockHeaderResult = result
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get block header: [%w]", err)
	}

	blockHeader, err := convertBlockHeader(getBlockHeaderResult)
	if err != nil {
		return nil, fmt.Errorf("failed to convert block header: %w", err)
	}

	return blockHeader, nil
}
