package electrum

import (
	"github.com/keep-network/keep-core/pkg/bitcoin"
)

type Client struct{}

func (c *Client) GetTransaction(
	transactionHash bitcoin.Hash,
) (*bitcoin.Transaction, error) {
	// TODO: Implementation.
	panic("not implemented")
}

func (c *Client) GetTransactionConfirmations(
	transactionHash bitcoin.Hash,
) (uint, error) {
	// TODO: Implementation.
	panic("not implemented")
}

func (c *Client) BroadcastTransaction(
	transaction *bitcoin.Transaction,
) error {
	// TODO: Implementation.
	panic("not implemented")
}

func (c *Client) GetCurrentBlockNumber() (uint, error) {
	// TODO: Implementation.
	panic("not implemented")
}

func (c *Client) GetBlockHeader(
	blockNumber uint,
) (*bitcoin.BlockHeader, error) {
	// TODO: Implementation.
	panic("not implemented")
}
