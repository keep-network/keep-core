package maintainer

import (
	"fmt"

	"github.com/keep-network/keep-core/pkg/bitcoin"
)

var errNoBlocksSet = fmt.Errorf("blockchain does not contain any blocks")

// localBitcoinChain represents a local Bitcoin chain.
type localBitcoinChain struct {
	blockHeaders             map[uint]*bitcoin.BlockHeader
	transactions             map[bitcoin.Hash]*bitcoin.Transaction
	transactionMerkleProof   *bitcoin.TransactionMerkleProof
	transactionConfirmations map[bitcoin.Hash]uint
}

// GetTransaction gets the transaction with the given transaction hash.
// If the transaction with the given hash was not found on the chain,
// this function returns an error.
func (lc *localBitcoinChain) GetTransaction(
	transactionHash bitcoin.Hash,
) (*bitcoin.Transaction, error) {
	transaction, found := lc.transactions[transactionHash]
	if !found {
		return nil, fmt.Errorf(
			"transaction with hash %v does not exist", transactionHash,
		)
	}

	return transaction, nil
}

// GetTransactionConfirmations gets the number of confirmations for the
// transaction with the given transaction hash. If the transaction with the
// given hash was not found on the chain, this function returns an error.
func (lc *localBitcoinChain) GetTransactionConfirmations(
	transactionHash bitcoin.Hash,
) (uint, error) {
	return lc.transactionConfirmations[transactionHash], nil
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

// GetLatestBlockHeight gets the height of the latest block (tip). If the
// latest block was not determined, this function returns an error.
func (lc *localBitcoinChain) GetLatestBlockHeight() (uint, error) {
	blockchainTip := uint(0)
	for blockHeaderHeight := range lc.blockHeaders {
		if blockHeaderHeight > blockchainTip {
			blockchainTip = blockHeaderHeight
		}
	}

	if blockchainTip == 0 {
		return 0, errNoBlocksSet
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

func (lc *localBitcoinChain) GetTransactionMerkleProof(
	transactionHash bitcoin.Hash,
	blockHeight uint,
) (*bitcoin.TransactionMerkleProof, error) {
	return lc.transactionMerkleProof, nil
}

func (lc *localBitcoinChain) GetTransactionsForPublicKeyHash(
	publicKeyHash [20]byte,
	limit int,
) ([]*bitcoin.Transaction, error) {
	panic("unsupported")
}

func (lc *localBitcoinChain) GetMempoolForPublicKeyHash(
	publicKeyHash [20]byte,
) ([]*bitcoin.Transaction, error) {
	panic("unsupported")
}

// SetBlockHeaders sets internal headers for testing purposes.
func (lc *localBitcoinChain) SetBlockHeaders(
	blockHeaders map[uint]*bitcoin.BlockHeader,
) {
	lc.blockHeaders = blockHeaders
}

// SetTransactions sets internal transactions for testing purposes.
func (lc *localBitcoinChain) SetTransactions(
	transactions map[bitcoin.Hash]*bitcoin.Transaction,
) {
	lc.transactions = transactions
}

// SetTransactionMerkleProof sets internal transaction Merkle proof for testing
// purposes.
func (lc *localBitcoinChain) SetTransactionMerkleProof(
	transactionMerkleProof *bitcoin.TransactionMerkleProof,
) {
	lc.transactionMerkleProof = transactionMerkleProof
}

// SetTransactionConfirmations sets internal amounts of transaction
// confirmations for testing purposes.
func (lc *localBitcoinChain) SetTransactionConfirmations(
	transactionConfirmations map[bitcoin.Hash]uint,
) {
	lc.transactionConfirmations = transactionConfirmations
}

// connectLocalBitcoinChain connects to the local Bitcoin chain and returns
// a chain handle.
func connectLocalBitcoinChain() *localBitcoinChain {
	return &localBitcoinChain{}
}
