package electrum

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"math"
	"sort"
	"strings"
	"sync"
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
	parentCtx   context.Context
	client      *electrum.Client
	clientMutex *sync.RWMutex
	config      Config
}

// Connect initializes handle with provided Config.
func Connect(parentCtx context.Context, config Config) (bitcoin.Chain, error) {
	if config.ConnectTimeout == 0 {
		config.ConnectTimeout = DefaultConnectTimeout
	}
	if config.ConnectRetryTimeout == 0 {
		config.ConnectRetryTimeout = DefaultConnectRetryTimeout
	}
	if config.RequestTimeout == 0 {
		config.RequestTimeout = DefaultRequestTimeout
	}
	if config.RequestRetryTimeout == 0 {
		config.RequestRetryTimeout = DefaultRequestRetryTimeout
	}
	if config.KeepAliveInterval == 0 {
		config.KeepAliveInterval = DefaultKeepAliveInterval
	}

	c := &Connection{
		parentCtx:   parentCtx,
		config:      config,
		clientMutex: &sync.RWMutex{},
	}

	if err := c.electrumConnect(); err != nil {
		return nil, fmt.Errorf("failed to initialize electrum client: [%w]", err)
	}

	if err := c.verifyServer(); err != nil {
		return nil, fmt.Errorf("failed to verify electrum server: [%w]", err)
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
	rawTx := hex.EncodeToString(transaction.Serialize())

	rawTxLogger := logger.With(
		zap.String("rawTx", rawTx),
	)
	rawTxLogger.Debugf("broadcasting transaction")

	var response string

	response, err := requestWithRetry(
		c,
		func(ctx context.Context, client *electrum.Client) (string, error) {
			return client.BroadcastTransaction(ctx, rawTx)
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

// GetTransactionsForPublicKeyHash get confirmed transactions that pays the
// given public key hash using either a P2PKH or P2WPKH script. The returned
// transactions are ordered by block height in the ascending order, i.e.
// the latest transaction is at the end of the list. The returned list does
// not contain unconfirmed transactions living in the mempool at the moment
// of request. The returned transactions list can be limited using the
// `limit` parameter. For example, if `limit` is set to `5`, only the
// latest five transactions will be returned.
func (c *Connection) GetTransactionsForPublicKeyHash(
	publicKeyHash [20]byte,
	limit int,
) ([]*bitcoin.Transaction, error) {
	p2pkh, err := bitcoin.PayToPublicKeyHash(publicKeyHash)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot build P2PKH for public key hash [0x%x]: [%v]",
			publicKeyHash,
			err,
		)
	}

	p2wpkh, err := bitcoin.PayToWitnessPublicKeyHash(publicKeyHash)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot build P2WPKH for public key hash [0x%x]: [%v]",
			publicKeyHash,
			err,
		)
	}

	p2pkhItems, err := c.getConfirmedScriptHistory(p2pkh)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot get P2PKH history for public key hash [0x%x]: [%v]",
			publicKeyHash,
			err,
		)
	}

	p2wpkhItems, err := c.getConfirmedScriptHistory(p2wpkh)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot get P2WPKH history for public key hash [0x%x]: [%v]",
			publicKeyHash,
			err,
		)
	}

	items := append(p2pkhItems, p2wpkhItems...)

	sort.SliceStable(
		items,
		func(i, j int) bool {
			return items[i].blockHeight < items[j].blockHeight
		},
	)

	var selectedItems []*scriptHistoryItem
	if len(items) > limit {
		selectedItems = items[len(items)-limit:]
	} else {
		selectedItems = items
	}

	transactions := make([]*bitcoin.Transaction, len(selectedItems))
	for i, item := range selectedItems {
		transaction, err := c.GetTransaction(item.txHash)
		if err != nil {
			return nil, fmt.Errorf("cannot get transaction: [%v]", err)
		}

		transactions[i] = transaction
	}

	return transactions, nil
}

type scriptHistoryItem struct {
	txHash      bitcoin.Hash
	blockHeight int32
}

// getConfirmedScriptHistory returns a history of confirmed transactions for
// the given script (P2PKH, P2WPKH, P2SH, P2WSH, etc.). The returned list
// is sorted by the block height in the ascending order, i.e. the latest
// transaction is at the end of the list. The resulting list does not contain
// unconfirmed transactions living in the mempool at the moment of request.
func (c *Connection) getConfirmedScriptHistory(
	script []byte,
) ([]*scriptHistoryItem, error) {
	scriptHash := sha256.Sum256(script)
	reversedScriptHash := byteutils.Reverse(scriptHash[:])
	reversedScriptHashString := hex.EncodeToString(reversedScriptHash)

	items, err := requestWithRetry(
		c,
		func(
			ctx context.Context,
			client *electrum.Client,
		) ([]*electrum.GetMempoolResult, error) {
			return client.GetHistory(ctx, reversedScriptHashString)
		})
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get history for script [0x%x]: [%v]",
			script,
			err,
		)
	}

	// According to https://electrumx.readthedocs.io/en/latest/protocol-methods.html#blockchain-scripthash-get-history
	// unconfirmed items living in the mempool are appended at the end of the
	// returned list and their height value is either -1 or 0. That means
	// we need to take all items with height >0 to obtain a confirmed txs
	// history.
	confirmedItems := make([]*scriptHistoryItem, 0)
	for _, item := range items {
		if item.Height > 0 {
			txHash, err := bitcoin.NewHashFromString(
				item.Hash,
				bitcoin.ReversedByteOrder,
			)
			if err != nil {
				return nil, fmt.Errorf(
					"cannot parse hash [%s]: [%v]",
					item.Hash,
					err,
				)
			}

			confirmedItems = append(
				confirmedItems, &scriptHistoryItem{
					txHash:      txHash,
					blockHeight: item.Height,
				},
			)
		}
	}

	// The list returned from client.GetHistory is sorted by the block height
	// in the ascending order though we are sorting it again just in case
	// (e.g. API contract changes).
	sort.SliceStable(
		confirmedItems,
		func(i, j int) bool {
			return confirmedItems[i].blockHeight < confirmedItems[j].blockHeight
		},
	)

	return confirmedItems, nil
}

func (c *Connection) electrumConnect() error {
	var client *electrum.Client
	var err error
	switch c.config.Protocol {
	case TCP:
		logger.Debug("establishing TCP connection to electrum server...")
		client, err = connectWithRetry(
			c,
			func(ctx context.Context) (*electrum.Client, error) {
				return electrum.NewClientTCP(ctx, c.config.URL)
			},
		)
	case SSL:
		// TODO: Implement certificate verification to be able to disable the `InsecureSkipVerify: true` workaround.
		// #nosec G402 (TLS InsecureSkipVerify set true)
		tlsConfig := &tls.Config{InsecureSkipVerify: true}

		logger.Debug("establishing SSL connection to electrum server...")
		client, err = connectWithRetry(
			c,
			func(ctx context.Context) (*electrum.Client, error) {
				return electrum.NewClientSSL(ctx, c.config.URL, tlsConfig)
			},
		)
	default:
		err = fmt.Errorf("unsupported protocol: [%s]", c.config.Protocol)
	}

	if err == nil {
		c.client = client
	}

	return err
}

func (c *Connection) verifyServer() error {
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
		return fmt.Errorf("failed to get server version: [%w]", err)
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
			c.config.URL,
			server.protocol,
			strings.Join(supportedProtocolVersions, ","),
		)
	}

	return nil
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
			} else {
				// Adjust ticker starting at the time of the latest successful ping.
				ticker = time.NewTicker(c.config.KeepAliveInterval)
			}
		case <-c.parentCtx.Done():
			ticker.Stop()
			c.client.Shutdown()
			return
		}
	}
}

func connectWithRetry(
	c *Connection,
	newClientFn func(ctx context.Context) (*electrum.Client, error),
) (*electrum.Client, error) {
	var result *electrum.Client
	err := wrappers.DoWithDefaultRetry(
		c.parentCtx,
		c.config.ConnectRetryTimeout,
		func(ctx context.Context) error {
			connectCtx, connectCancel := context.WithTimeout(
				ctx,
				c.config.ConnectTimeout,
			)
			defer connectCancel()

			client, err := newClientFn(connectCtx)
			if err == nil {
				result = client
			}

			return err
		},
	)

	return result, err
}

func requestWithRetry[K interface{}](
	c *Connection,
	requestFn func(ctx context.Context, client *electrum.Client) (K, error),
) (K, error) {
	var result K

	err := wrappers.DoWithDefaultRetry(
		c.parentCtx,
		c.config.RequestRetryTimeout,
		func(ctx context.Context) error {
			if err := c.reconnectIfShutdown(); err != nil {
				return err
			}

			requestCtx, requestCancel := context.WithTimeout(ctx, c.config.RequestTimeout)
			defer requestCancel()

			c.clientMutex.RLock()
			r, err := requestFn(requestCtx, c.client)
			c.clientMutex.RUnlock()

			if err != nil {
				return fmt.Errorf("request failed: [%w]", err)
			}

			result = r
			return nil
		})

	return result, err
}

func (c *Connection) reconnectIfShutdown() error {
	c.clientMutex.Lock()
	defer c.clientMutex.Unlock()

	isClientShutdown := c.client.IsShutdown()
	if isClientShutdown {
		logger.Warn("connection to electrum server is down; reconnecting...")
		err := c.electrumConnect()
		if err != nil {
			return fmt.Errorf("failed to reconnect to electrum server: [%w]", err)
		}
		logger.Info("reconnected to electrum server")
	}

	return nil
}
