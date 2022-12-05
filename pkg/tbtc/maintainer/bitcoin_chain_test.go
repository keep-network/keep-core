package maintainer

import (
	"fmt"

	"github.com/keep-network/keep-core/pkg/bitcoin"
)

// localBitcoinChain represents a local Bitcoin chain.
type localBitcoinChain struct {
	blockHeaders map[uint]*bitcoin.BlockHeader
}

// GetTransaction gets the transaction with the given transaction hash.
// If the transaction with the given hash was not found on the chain,
// this function returns an error.
func (lc *localBitcoinChain) GetTransaction(
	transactionHash bitcoin.Hash,
) (*bitcoin.Transaction, error) {
	panic("unsupported")
}

// GetTransactionConfirmations gets the number of confirmations for the
// transaction with the given transaction hash. If the transaction with the
// given hash was not found on the chain, this function returns an error.
func (lc *localBitcoinChain) GetTransactionConfirmations(
	transactionHash bitcoin.Hash,
) (uint, error) {
	panic("unsupported")
}

// BroadcastTransaction broadcasts the given transaction over the
// network of the Bitcoin chain nodes. If the broadcast action could not be
// done, this function returns an error. This function does not give any
// guarantees regarding transaction mining. The transaction may be mined or
// rejected eventually.
func (lc *localBitcoinChain) BroadcastTransaction(
	transaction *bitcoin.Transaction,
) error {
	panic("unsupported")
}

// GetCurrentBlockNumber gets the number of the current block. If the
// current block was not determined, this function returns an error.
func (lc *localBitcoinChain) GetCurrentBlockNumber() (uint, error) {
	blockchainTip := uint(0)
	for blockHeaderHeight := range lc.blockHeaders {
		if blockHeaderHeight > blockchainTip {
			blockchainTip = blockHeaderHeight
		}
	}

	if blockchainTip == 0 {
		return 0, fmt.Errorf("blockchain does not contain any blocks")
	}

	return blockchainTip, nil
}

// GetBlockHeader gets the block header for the given block number. If the
// block with the given number was not found on the chain, this function
// returns an error.
func (lc *localBitcoinChain) GetBlockHeader(
	blockNumber uint,
) (*bitcoin.BlockHeader, error) {
	blockHeader, found := lc.blockHeaders[blockNumber]
	if !found {
		return nil, fmt.Errorf(
			"block header at height %v does not exist",
			blockNumber,
		)
	}

	return blockHeader, nil
}

// SetBlockHeaders sets internal headers for testing purposes.
func (lc *localBitcoinChain) SetBlockHeaders(
	blockHeaders map[uint]*bitcoin.BlockHeader,
) {
	lc.blockHeaders = blockHeaders
}

// connectLocalBitcoinChain connects to the local Bitcoin chain and returns
// a chain handle.
func connectLocalBitcoinChain() *localBitcoinChain {
	return &localBitcoinChain{}
}
