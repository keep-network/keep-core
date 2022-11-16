package electrum

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"
	"time"

	"github.com/checksum0/go-electrum/electrum"
	"github.com/ipfs/go-log"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"

	"github.com/keep-network/keep-common/pkg/wrappers"
	"github.com/keep-network/keep-core/pkg/bitcoin"
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
		client, err = electrum.NewClientTCP(ctx, config.URL)
	case SSL:
		// TODO: Implement certificate verification to be able to disable the `InsecureSkipVerify: true` workaround.
		tlsConfig := &tls.Config{InsecureSkipVerify: true}
		client, err = electrum.NewClientSSL(ctx, config.URL, tlsConfig)
	default:
		return nil, fmt.Errorf("unsupported protocol: [%s]", config.Protocol)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to initialize electrum client: [%w]", err)
	}

	// Get the server's details.
	err = wrappers.DoWithDefaultRetry(
		config.RequestRetryTimeout,
		func(ctx context.Context,
		) error {
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
		return nil, fmt.Errorf("failed to get raw transaction [%s]: [%w]", txID, err)
	}

	result, err := convertRawTransaction(rawTransaction)
	if err != nil {
		return nil, fmt.Errorf("failed to convert transaction: [%w]", err)
	}

	return result, nil
}

func (c *Connection) GetTransactionConfirmations(
	transactionHash bitcoin.Hash,
) (uint, error) {
	// TODO: Implementation.
	panic("not implemented")
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
