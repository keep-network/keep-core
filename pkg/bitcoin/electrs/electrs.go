package electrs

import (
	"io"
	"net/http"
	"time"

	"github.com/keep-network/keep-core/pkg/bitcoin"
)

// Connection is a handle for interactions with Electrum.
type Connection struct {
	url     string
	client  httpClient
	timeout time.Duration
}

type httpClient interface {
	Get(url string) (*http.Response, error)
	Post(url string, contentType string, body io.Reader) (*http.Response, error)
}

func (c *Connection) GetTransaction(
	transactionHash bitcoin.Hash,
) (*bitcoin.Transaction, error) {
	// TODO: Implementation.
	panic("not implemented")
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
