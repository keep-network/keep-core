package electrs

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ipfs/go-log"

	"github.com/keep-network/keep-core/pkg/bitcoin"
)

const (
	// DefaultRequestTimeout is a default timeout used for HTTP requests.
	DefaultRequestTimeout = 5 * time.Second
	// DefaultRetryTimeout is a default timeout used for retries.
	DefaultRetryTimeout = 60 * time.Second
)

var logger = log.Logger("keep-electrs")

// Config holds configurable properties.
type Config struct {
	URL            string
	RequestTimeout time.Duration
	RetryTimeout   time.Duration
}

// Connection is a handle for interactions with Electrum.
type Connection struct {
	url          string
	client       httpClient
	retryTimeout time.Duration
}

type httpClient interface {
	Get(url string) (*http.Response, error)
	Post(url string, contentType string, body io.Reader) (*http.Response, error)
}

// Connect initializes handle with provided Config.
func Connect(config Config) (bitcoin.Chain, error) {
	var netClient = &http.Client{
		Timeout: config.RequestTimeout,
	}

	if config.URL == "" {
		return nil, fmt.Errorf("URL not set")
	}

	return &Connection{
		url:          strings.TrimSuffix(config.URL, "/"),
		client:       netClient,
		retryTimeout: config.RetryTimeout,
	}, nil
}

// GetTransaction gets the transaction with the given transaction hash.
// If the transaction with the given hash was not found on the chain,
// this function returns an error.
func (c *Connection) GetTransaction(
	transactionHash bitcoin.Hash,
) (*bitcoin.Transaction, error) {
	txID := transactionHash.String(bitcoin.ReversedByteOrder)

	transactionJSON, err := c.httpGet(fmt.Sprintf("tx/%s", txID))
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get the transaction [%s]: [%w]",
			txID,
			err,
		)
	}

	tx, err := decodeJSON[transaction](transactionJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: [%w]", err)
	}

	result, err := tx.convert()
	if err != nil {
		return nil, fmt.Errorf("failed to convert transaction: %w", err)
	}

	return &result, nil
}

// GetTransactionConfirmations gets the number of confirmations for the
// transaction with the given transaction hash. If the transaction with the
// given hash was not found on the chain, this function returns an error.
func (c *Connection) GetTransactionConfirmations(
	transactionHash bitcoin.Hash,
) (uint, error) {
	txID := transactionHash.String(bitcoin.ReversedByteOrder)

	transactionJSON, err := c.httpGet(fmt.Sprintf("tx/%s", txID))
	if err != nil {
		return 0, fmt.Errorf(
			"failed to get the transaction [%s]: [%w]",
			txID,
			err,
		)
	}

	tx, err := decodeJSON[transaction](transactionJSON)
	if err != nil {
		return 0, fmt.Errorf("failed to decode response: [%w]", err)
	}

	if !tx.Status.Confirmed {
		return 0, nil
	}

	txBlockHeight := tx.Status.BlockHeight

	currentBlockHeight, err := c.GetCurrentBlockNumber()
	if err != nil {
		return 0, fmt.Errorf("failed to get current block height: [%w]", err)
	}

	confirmations := uint(1)

	if currentBlockHeight > txBlockHeight {
		confirmations += currentBlockHeight - txBlockHeight
	}
	return confirmations, nil
}

func (c *Connection) BroadcastTransaction(
	transaction *bitcoin.Transaction,
) error {
	// TODO: Implementation.
	panic("not implemented")
}

// GetCurrentBlockNumber gets the number of the current block. If the
// current block was not determined, this function returns an error.
func (c *Connection) GetCurrentBlockNumber() (uint, error) {
	blockHeight, err := c.httpGet("blocks/tip/height")
	if err != nil {
		return 0, fmt.Errorf("failed to get the blocks tip height: [%w]", err)
	}

	result, err := strconv.ParseUint(string(blockHeight), 10, 32)
	if err != nil {
		return 0, fmt.Errorf("failed to parse block height [%s]: [%w]", blockHeight, err)
	}

	return uint(result), nil
}

// GetBlockHeader gets the block header for the given block number. If the
// block with the given number was not found on the chain, this function
// returns an error.
func (c *Connection) GetBlockHeader(
	blockNumber uint,
) (*bitcoin.BlockHeader, error) {
	blockHash, err := c.httpGet(fmt.Sprintf("block-height/%d", blockNumber))
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get hash of block at height [%d]: [%w]",
			blockNumber,
			err,
		)
	}

	blockHeaderJSON, err := c.httpGet(fmt.Sprintf("block/%s", blockHash))
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get header of block [%s]: [%w]",
			blockHash,
			err,
		)
	}

	blockHeader, err := decodeJSON[blockHeader](blockHeaderJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to decode block header response: [%w]", err)
	}

	result, err := blockHeader.convert()
	if err != nil {
		return nil, fmt.Errorf("failed to convert block header: %w", err)
	}

	return &result, nil
}
