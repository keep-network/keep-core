package spv

import (
	"fmt"
	"sync"

	"github.com/keep-network/keep-core/pkg/bitcoin"
)

//lint:ignore U1000 Ignore unused type temporarily.
type localBitcoinChain struct {
	mutex sync.Mutex

	transactionConfirmations map[bitcoin.Hash]uint
	blockHeaders             map[uint]*bitcoin.BlockHeader
}

//lint:ignore U1000 Ignore unused function temporarily.
func newLocalBitcoinChain() *localBitcoinChain {
	return &localBitcoinChain{
		transactionConfirmations: make(map[bitcoin.Hash]uint),
		blockHeaders:             make(map[uint]*bitcoin.BlockHeader),
	}
}

func (lbc *localBitcoinChain) GetTransaction(transactionHash bitcoin.Hash) (
	*bitcoin.Transaction,
	error,
) {
	panic("unsupported")
}

func (lbc *localBitcoinChain) GetTransactionConfirmations(transactionHash bitcoin.Hash) (
	uint,
	error,
) {
	lbc.mutex.Lock()
	defer lbc.mutex.Unlock()

	if transactionConfirmations, exists :=
		lbc.transactionConfirmations[transactionHash]; exists {
		return transactionConfirmations, nil
	}

	return 0, fmt.Errorf("transaction not found")
}

func (lbc *localBitcoinChain) BroadcastTransaction(transaction *bitcoin.Transaction) error {
	panic("unsupported")
}

func (lbc *localBitcoinChain) GetLatestBlockHeight() (uint, error) {
	lbc.mutex.Lock()
	defer lbc.mutex.Unlock()

	// Return the highest block header's height.
	blockchainTip := uint(0)
	for blockHeaderHeight := range lbc.blockHeaders {
		if blockHeaderHeight > blockchainTip {
			blockchainTip = blockHeaderHeight
		}
	}

	if blockchainTip == 0 {
		return 0, fmt.Errorf("block headers not found")
	}

	return blockchainTip, nil
}

func (lbc *localBitcoinChain) GetBlockHeader(blockHeight uint) (
	*bitcoin.BlockHeader,
	error,
) {
	lbc.mutex.Lock()
	defer lbc.mutex.Unlock()

	if blockHeader, exists := lbc.blockHeaders[blockHeight]; exists {
		return blockHeader, nil
	}

	return nil, fmt.Errorf("block header does not exist")
}

func (lbc *localBitcoinChain) GetTransactionMerkleProof(
	transactionHash bitcoin.Hash,
	blockHeight uint,
) (*bitcoin.TransactionMerkleProof, error) {
	panic("unsupported")
}

func (lbc *localBitcoinChain) GetTransactionsForPublicKeyHash(
	publicKeyHash [20]byte,
	limit int,
) ([]*bitcoin.Transaction, error) {
	panic("unsupported")
}

func (lbc *localBitcoinChain) GetMempoolForPublicKeyHash(publicKeyHash [20]byte) (
	[]*bitcoin.Transaction,
	error,
) {
	panic("unsupported")
}

func (lbc *localBitcoinChain) EstimateSatPerVByteFee(blocks uint32) (
	int64,
	error,
) {
	panic("unsupported")
}

func (lbc *localBitcoinChain) addBlockHeader(
	blockNumber uint,
	blockHeader *bitcoin.BlockHeader,
) error {
	lbc.mutex.Lock()
	defer lbc.mutex.Unlock()

	if _, exists := lbc.blockHeaders[blockNumber]; exists {
		return fmt.Errorf("block header already exists")
	}

	lbc.blockHeaders[blockNumber] = blockHeader

	return nil
}

func (lbc *localBitcoinChain) addTransactionConfirmations(
	transactionHash bitcoin.Hash,
	transactionConfirmations uint,
) error {
	lbc.mutex.Lock()
	defer lbc.mutex.Unlock()

	if _, exists := lbc.transactionConfirmations[transactionHash]; exists {
		return fmt.Errorf("transaction confirmations already set")
	}

	lbc.transactionConfirmations[transactionHash] = transactionConfirmations

	return nil
}
