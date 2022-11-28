package bitcoin

import "fmt"

// LocalChain represents a local Bitcoin chain.
type LocalChain struct {
	blockHeaders map[uint]*BlockHeader
}

// GetTransaction gets the transaction with the given transaction hash.
// If the transaction with the given hash was not found on the chain,
// this function returns an error.
func (lc *LocalChain) GetTransaction(transactionHash Hash) (*Transaction, error) {
	panic("unsupported")
}

// GetTransactionConfirmations gets the number of confirmations for the
// transaction with the given transaction hash. If the transaction with the
// given hash was not found on the chain, this function returns an error.
func (lc *LocalChain) GetTransactionConfirmations(transactionHash Hash) (uint, error) {
	panic("unsupported")
}

// BroadcastTransaction broadcasts the given transaction over the
// network of the Bitcoin chain nodes. If the broadcast action could not be
// done, this function returns an error. This function does not give any
// guarantees regarding transaction mining. The transaction may be mined or
// rejected eventually.
func (lc *LocalChain) BroadcastTransaction(transaction *Transaction) error {
	panic("unsupported")
}

// GetCurrentBlockNumber gets the number of the current block. If the
// current block was not determined, this function returns an error.
func (lc *LocalChain) GetCurrentBlockNumber() (uint, error) {
	blockchainTip := uint(0)
	for blockHeaderHeight := range lc.blockHeaders {
		if blockHeaderHeight > blockchainTip {
			blockchainTip = blockHeaderHeight
		}
	}

	if blockchainTip == 0 {
		return 0, fmt.Errorf("could not get current block block number")
	}

	return blockchainTip, nil
}

// GetBlockHeader gets the block header for the given block number. If the
// block with the given number was not found on the chain, this function
// returns an error.
func (lc *LocalChain) GetBlockHeader(blockNumber uint) (*BlockHeader, error) {
	blockHeader, found := lc.blockHeaders[blockNumber]
	if !found {
		return nil, fmt.Errorf(
			"could not find block header at height %v",
			blockNumber,
		)
	}

	return blockHeader, nil
}

// SetBlockHeaders sets internal headers for testing purposes.
func (lc *LocalChain) SetBlockHeaders(blockHeaders map[uint]*BlockHeader) {
	lc.blockHeaders = blockHeaders
}

// ConnectLocal connects to the local Bitcoin chain and returns a chain handle.
func ConnectLocal() *LocalChain {
	return &LocalChain{}
}
