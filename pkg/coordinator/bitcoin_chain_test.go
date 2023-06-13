package coordinator_test

import (
	"fmt"
	"sync"

	"github.com/keep-network/keep-core/pkg/bitcoin"
)

type localBitcoinChain struct {
	transactionsMutex         sync.Mutex
	transactions              map[bitcoin.Hash]*bitcoin.Transaction
	transactionsConfirmations map[bitcoin.Hash]uint

	feeMutex                 sync.Mutex
	satPerVByteFeeEstimation map[uint32]int64
}

func newLocalBitcoinChain() *localBitcoinChain {
	return &localBitcoinChain{
		transactions:              make(map[bitcoin.Hash]*bitcoin.Transaction),
		transactionsConfirmations: make(map[bitcoin.Hash]uint),
		satPerVByteFeeEstimation:  make(map[uint32]int64),
	}
}

func (lbc *localBitcoinChain) GetTransaction(
	transactionHash bitcoin.Hash,
) (*bitcoin.Transaction, error) {
	lbc.transactionsMutex.Lock()
	defer lbc.transactionsMutex.Unlock()

	transaction, ok := lbc.transactions[transactionHash]
	if !ok {
		return nil, fmt.Errorf("transaction not found")
	}
	return transaction, nil
}

func (lbc *localBitcoinChain) setTransaction(
	transactionHash bitcoin.Hash,
	transaction *bitcoin.Transaction,
) {
	lbc.transactionsMutex.Lock()
	defer lbc.transactionsMutex.Unlock()

	lbc.transactions[transactionHash] = transaction
}

func (lbc *localBitcoinChain) GetTransactionConfirmations(
	transactionHash bitcoin.Hash,
) (uint, error) {
	lbc.transactionsMutex.Lock()
	defer lbc.transactionsMutex.Unlock()

	if confirmations, ok := lbc.transactionsConfirmations[transactionHash]; ok {
		return confirmations, nil
	}

	return 0, fmt.Errorf("transaction not found")
}

func (lbc *localBitcoinChain) setTransactionConfirmations(
	transactionHash bitcoin.Hash,
	confirmations uint,
) {
	lbc.transactionsMutex.Lock()
	defer lbc.transactionsMutex.Unlock()

	lbc.transactionsConfirmations[transactionHash] = confirmations
}

func (lbc *localBitcoinChain) BroadcastTransaction(
	transaction *bitcoin.Transaction,
) error {
	panic("unsupported")
}

func (lbc *localBitcoinChain) GetLatestBlockHeight() (uint, error) {
	panic("unsupported")
}

func (lbc *localBitcoinChain) GetBlockHeader(
	blockNumber uint,
) (*bitcoin.BlockHeader, error) {
	panic("unsupported")
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

func (lbc *localBitcoinChain) GetMempoolForPublicKeyHash(
	publicKeyHash [20]byte,
) ([]*bitcoin.Transaction, error) {
	panic("unsupported")
}

func (lbc *localBitcoinChain) EstimateSatPerVByteFee(
	blocks uint32,
) (int64, error) {
	lbc.feeMutex.Lock()
	defer lbc.feeMutex.Unlock()

	return lbc.satPerVByteFeeEstimation[blocks], nil
}

func (lbc *localBitcoinChain) setEstimateSatPerVByteFee(
	blocks uint32,
	fee int64,
) {
	lbc.feeMutex.Lock()
	defer lbc.feeMutex.Unlock()

	lbc.satPerVByteFeeEstimation[blocks] = fee
}
