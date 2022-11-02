package electrs

import (
	"fmt"
	"io"
	"net/http"
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

func (c *Connection) GetTransaction(
	transactionHash bitcoin.Hash,
) (*bitcoin.Transaction, error) {
	txID := transactionHash.String(bitcoin.ReversedByteOrder)

	responseBody, err := c.httpGet(fmt.Sprintf("tx/%s", txID))
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get the transaction [%s]: [%w]",
			txID,
			err,
		)
	}

	tx, err := decodeJSON[transaction](responseBody)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: [%w]", err)
	}

	result, err := tx.convert()
	if err != nil {
		return nil, fmt.Errorf("failed to convert transaction: %w", err)
	}

	return &result, nil
}

func (c *Connection) GetTransactionConfirmations(
	transactionHash bitcoin.Hash,
) (uint, error) {
	// TODO: Implementation.
	panic("not implemented")
}

func (c *Connection) BroadcastTransaction(
	transaction *bitcoin.Transaction,
) error {
	// TODO: Implementation.
	panic("not implemented")
}

func (c *Connection) GetCurrentBlockNumber() (uint, error) {
	// TODO: Implementation.
	panic("not implemented")
}

func (c *Connection) GetBlockHeader(
	blockNumber uint,
) (*bitcoin.BlockHeader, error) {
	// TODO: Implementation.
	panic("not implemented")
}
