package electrum

import (
	"context"
	"crypto/sha256"
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
	clientMutex *sync.Mutex
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
		clientMutex: &sync.Mutex{},
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
		},
		"GetRawTransaction",
	)
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
		},
		"GetRawTransaction",
	)
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
			},
			"GetHistory",
		)
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
		},
		"BroadcastTransaction",
	)
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
		},
		"SubscribeHeaders",
	)
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
		"GetBlockHeader",
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

// GetTransactionMerkleProof gets the Merkle proof for a given transaction.
// The transaction's hash and the block the transaction was included in the
// blockchain need to be provided.
func (c *Connection) GetTransactionMerkleProof(
	transactionHash bitcoin.Hash,
	blockHeight uint,
) (*bitcoin.TransactionMerkleProof, error) {
	txID := transactionHash.Hex(bitcoin.ReversedByteOrder)

	getMerkleProofResult, err := requestWithRetry(
		c,
		func(
			ctx context.Context,
			client *electrum.Client,
		) (*electrum.GetMerkleProofResult, error) {
			return client.GetMerkleProof(
				ctx,
				txID,
				uint32(blockHeight),
			)
		},
		"GetMerkleProof",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get merkle proof: [%w]", err)
	}

	return convertMerkleProof(getMerkleProofResult), nil
}

// GetTransactionsForPublicKeyHash gets confirmed transactions that pays the
// given public key hash using either a P2PKH or P2WPKH script. The returned
// transactions are ordered by block height in the ascending order, i.e.
// the latest transaction is at the end of the list. The returned list does
// not contain unconfirmed transactions living in the mempool at the moment
// of request. The returned transactions list can be limited using the
// `limit` parameter. For example, if `limit` is set to `5`, only the
// latest five transactions will be returned. Note that taking an unlimited
// transaction history may be time-consuming as this function fetches
// complete transactions with all necessary data.
func (c *Connection) GetTransactionsForPublicKeyHash(
	publicKeyHash [20]byte,
	limit int,
) ([]*bitcoin.Transaction, error) {
	txHashes, err := c.GetTxHashesForPublicKeyHash(publicKeyHash)
	if err != nil {
		return nil, err
	}

	var selectedTxHashes []bitcoin.Hash
	if len(txHashes) > limit {
		selectedTxHashes = txHashes[len(txHashes)-limit:]
	} else {
		selectedTxHashes = txHashes
	}

	transactions := make([]*bitcoin.Transaction, len(selectedTxHashes))
	for i, txHash := range selectedTxHashes {
		transaction, err := c.GetTransaction(txHash)
		if err != nil {
			return nil, fmt.Errorf("cannot get transaction: [%v]", err)
		}

		transactions[i] = transaction
	}

	return transactions, nil
}

// GetTxHashesForPublicKeyHash gets hashes of confirmed transactions that pays
// the given public key hash using either a P2PKH or P2WPKH script. The returned
// transactions hashes are ordered by block height in the ascending order, i.e.
// the latest transaction hash is at the end of the list. The returned list does
// not contain unconfirmed transactions hashes living in the mempool at the
// moment of request.
func (c *Connection) GetTxHashesForPublicKeyHash(
	publicKeyHash [20]byte,
) ([]bitcoin.Hash, error) {
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

	txHashes := make([]bitcoin.Hash, len(items))
	for i, item := range items {
		txHashes[i] = item.txHash
	}

	return txHashes, nil
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
		},
		"GetHistory",
	)
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

// GetMempoolForPublicKeyHash gets the unconfirmed mempool transactions
// that pays the given public key hash using either a P2PKH or P2WPKH script.
// The returned transactions are in an indefinite order.
func (c *Connection) GetMempoolForPublicKeyHash(
	publicKeyHash [20]byte,
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

	p2pkhItems, err := c.getScriptMempool(p2pkh)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot get P2PKH mempool items for public key hash [0x%x]: [%v]",
			publicKeyHash,
			err,
		)
	}

	p2wpkhItems, err := c.getScriptMempool(p2wpkh)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot get P2WPKH mempool items for public key hash [0x%x]: [%v]",
			publicKeyHash,
			err,
		)
	}

	items := append(p2pkhItems, p2wpkhItems...)

	transactions := make([]*bitcoin.Transaction, len(items))
	for i, item := range items {
		transaction, err := c.GetTransaction(item.txHash)
		if err != nil {
			return nil, fmt.Errorf("cannot get transaction: [%v]", err)
		}

		transactions[i] = transaction
	}

	return transactions, nil
}

type scriptMempoolItem struct {
	txHash      bitcoin.Hash
	blockHeight int32
	fee         uint32
}

// getScriptMempool returns unconfirmed mempool transactions for
// the given script (P2PKH, P2WPKH, P2SH, P2WSH, etc.). The returned list
// is in an indefinite order.
func (c *Connection) getScriptMempool(
	script []byte,
) ([]*scriptMempoolItem, error) {
	scriptHash := sha256.Sum256(script)
	reversedScriptHash := byteutils.Reverse(scriptHash[:])
	reversedScriptHashString := hex.EncodeToString(reversedScriptHash)

	items, err := requestWithRetry(
		c,
		func(
			ctx context.Context,
			client *electrum.Client,
		) ([]*electrum.GetMempoolResult, error) {
			return client.GetMempool(ctx, reversedScriptHashString)
		},
		"GetMempool",
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get mempool for script [0x%x]: [%v]",
			script,
			err,
		)
	}

	convertedItems := make([]*scriptMempoolItem, len(items))
	for i, item := range items {
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

		convertedItems[i] = &scriptMempoolItem{
			txHash:      txHash,
			blockHeight: item.Height,
			fee:         item.Fee,
		}
	}

	return convertedItems, nil
}

// GetUtxosForPublicKeyHash gets unspent outputs of confirmed transactions that
// are controlled by the given public key hash (either a P2PKH or P2WPKH script).
// The returned UTXOs are ordered by block height in the ascending order, i.e.
// the latest UTXO is at the end of the list. The returned list does not contain
// unspent outputs of unconfirmed transactions living in the mempool at the
// moment of request. Outputs used as inputs of confirmed or mempool
// transactions are not returned as well because they are no longer UTXOs.
func (c *Connection) GetUtxosForPublicKeyHash(
	publicKeyHash [20]byte,
) ([]*bitcoin.UnspentTransactionOutput, error) {
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

	p2pkhItems, err := c.getScriptUtxos(p2pkh, true)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot get P2PKH UTXOs for public key hash [0x%x]: [%v]",
			publicKeyHash,
			err,
		)
	}

	p2wpkhItems, err := c.getScriptUtxos(p2wpkh, true)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot get P2WPKH UTXOs for public key hash [0x%x]: [%v]",
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

	utxos := make([]*bitcoin.UnspentTransactionOutput, len(items))
	for i, item := range items {
		utxos[i] = &bitcoin.UnspentTransactionOutput{
			Outpoint: &bitcoin.TransactionOutpoint{
				TransactionHash: item.txHash,
				OutputIndex:     item.outputIndex,
			},
			Value: int64(item.value),
		}
	}

	return utxos, nil
}

// GetMempoolUtxosForPublicKeyHash gets unspent outputs of unconfirmed transactions
// that are controlled by the given public key hash (either a P2PKH or P2WPKH script).
// The returned UTXOs are in an indefinite order. The returned list does not
// contain unspent outputs of confirmed transactions. Outputs used as inputs of
// confirmed or mempool transactions are not returned as well because they are
// no longer UTXOs.
func (c *Connection) GetMempoolUtxosForPublicKeyHash(
	publicKeyHash [20]byte,
) ([]*bitcoin.UnspentTransactionOutput, error) {
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

	p2pkhItems, err := c.getScriptUtxos(p2pkh, false)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot get P2PKH UTXOs for public key hash [0x%x]: [%v]",
			publicKeyHash,
			err,
		)
	}

	p2wpkhItems, err := c.getScriptUtxos(p2wpkh, false)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot get P2WPKH UTXOs for public key hash [0x%x]: [%v]",
			publicKeyHash,
			err,
		)
	}

	items := append(p2pkhItems, p2wpkhItems...)

	utxos := make([]*bitcoin.UnspentTransactionOutput, len(items))
	for i, item := range items {
		utxos[i] = &bitcoin.UnspentTransactionOutput{
			Outpoint: &bitcoin.TransactionOutpoint{
				TransactionHash: item.txHash,
				OutputIndex:     item.outputIndex,
			},
			Value: int64(item.value),
		}
	}

	return utxos, nil
}

type scriptUtxoItem struct {
	txHash      bitcoin.Hash
	outputIndex uint32
	value       uint64
	blockHeight uint32
}

// getScriptUtxos returns unspent outputs of confirmed/unconfirmed transactions
// that are locked using the given script (P2PKH, P2WPKH, P2SH, P2WSH, etc.).
//
// If the `confirmed` flag is true, the returned list contains unspent outputs
// of confirmed transactions, sorted by the block height in the ascending order,
// i.e. the latest UTXO is at the end of the list. The resulting list does not
// contain unspent outputs of unconfirmed transactions living in the mempool
// at the moment of request.
//
// If the `confirmed` flag is false, the returned list contains unspent outputs
// of unconfirmed transactions, in an indefinite order. The resulting list
// does not contain unspent outputs of confirmed transactions.
//
// In both cases, the resulted list DOES NOT CONTAIN outputs already used as
// inputs of confirmed or mempool transactions because they are no longer UTXOs.
func (c *Connection) getScriptUtxos(
	script []byte,
	confirmed bool,
) ([]*scriptUtxoItem, error) {
	scriptHash := sha256.Sum256(script)
	reversedScriptHash := byteutils.Reverse(scriptHash[:])
	reversedScriptHashString := hex.EncodeToString(reversedScriptHash)

	items, err := requestWithRetry(
		c,
		func(
			ctx context.Context,
			client *electrum.Client,
		) ([]*electrum.ListUnspentResult, error) {
			return client.ListUnspent(ctx, reversedScriptHashString)
		},
		"ListUnspent",
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get UTXOs for script [0x%x]: [%v]",
			script,
			err,
		)
	}

	// According to https://electrumx.readthedocs.io/en/stable/protocol-methods.html#blockchain-scripthash-listunspent
	// unconfirmed items living in the mempool are appended at the end of the
	// returned list and their height value is either -1 or 0. That means
	// we need to take all items with height >0 to obtain confirmed UTXO
	// items and <=0 if we want to take the unconfirmed ones.
	var filterFn func(item *electrum.ListUnspentResult) bool
	if confirmed {
		filterFn = func(item *electrum.ListUnspentResult) bool {
			return item.Height > 0
		}
	} else {
		filterFn = func(item *electrum.ListUnspentResult) bool {
			return item.Height <= 0
		}
	}

	filteredItems := make([]*scriptUtxoItem, 0)
	for _, item := range items {
		if filterFn(item) {
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

			filteredItems = append(
				filteredItems, &scriptUtxoItem{
					txHash:      txHash,
					outputIndex: item.Position,
					value:       item.Value,
					blockHeight: item.Height,
				},
			)
		}
	}

	if confirmed {
		// The list returned from client.ListUnspent is sorted by the block height
		// in the ascending order though we are sorting it again just in case
		// (e.g. API contract changes). Sorting makes sense only for confirmed
		// items as unconfirmed ones have a block height of -1 or 0.
		sort.SliceStable(
			filteredItems,
			func(i, j int) bool {
				return filteredItems[i].blockHeight < filteredItems[j].blockHeight
			},
		)
	}

	return filteredItems, nil
}

// EstimateSatPerVByteFee returns the estimated sat/vbyte fee for a
// transaction to be confirmed within the given number of blocks.
func (c *Connection) EstimateSatPerVByteFee(blocks uint32) (int64, error) {
	// According to Electrum protocol docs, the returned fee is BTC/KB.
	btcPerKbFee, err := requestWithRetry(
		c,
		func(
			ctx context.Context,
			client *electrum.Client,
		) (float32, error) {
			// TODO: client.GetFee calls Electrum's blockchain.estimatefee underneath.
			//       According to https://electrumx.readthedocs.io/en/latest/protocol-methods.html#blockchain-estimatefee,
			//       the blockchain.estimatefee function will be deprecated
			//       since version 1.4.2 of the protocol. We need to replace it
			//       somehow once it disappears from Electrum implementations.
			return client.GetFee(ctx, blocks)
		},
		"GetFee",
	)
	if err != nil {
		return 0, fmt.Errorf("failed to get fee: [%v]", err)
	}

	// According to Electrum protocol docs, if the daemon does not have
	// enough information to make an estimate, the integer -1 is returned.
	if btcPerKbFee < 0 {
		return 0, fmt.Errorf(
			"daemon does not have enough information to make an estimate",
		)
	}

	return convertBtcKbToSatVByte(btcPerKbFee), nil
}

func convertBtcKbToSatVByte(btcPerKbFee float32) int64 {
	// To convert from BTC/KB to sat/vbyte, we need to multiply by 1e8/1e3.
	satPerVByte := (1e8 / 1e3) * float64(btcPerKbFee)
	// Make sure the minimum returned sat/vbyte fee is always 1.
	satPerVByte = math.Max(satPerVByte, 1)
	// Round the returned fee to be an integer.
	return int64(math.Round(satPerVByte))
}

func (c *Connection) electrumConnect() error {
	var client *electrum.Client
	var err error

	logger.Debug("establishing connection to electrum server...")
	client, err = connectWithRetry(
		c,
		func(ctx context.Context) (*electrum.Client, error) {
			return electrum.NewClient(ctx, c.config.URL, nil)
		},
	)

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
		"ServerVersion",
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
				"Ping",
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
	requestName string,
) (K, error) {
	startTime := time.Now()
	logger.Debugf("starting [%s] request to Electrum server", requestName)

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

			c.clientMutex.Lock()
			r, err := requestFn(requestCtx, c.client)
			c.clientMutex.Unlock()

			if err != nil {
				return fmt.Errorf("request failed: [%w]", err)
			}

			result = r
			return nil
		})

	solveRequestOutcome := func(err error) string {
		if err != nil {
			return fmt.Sprintf("error: [%v]", err)
		}
		return "success"
	}

	logger.Debugf("[%s] request to Electrum server completed with [%s] after [%s]",
		requestName,
		solveRequestOutcome(err),
		time.Since(startTime),
	)

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
