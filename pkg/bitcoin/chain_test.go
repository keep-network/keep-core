package bitcoin

import (
	"fmt"
	"sync"
)

type localChain struct {
	transactionsMutex sync.Mutex
	transactions      map[Hash]*Transaction

	satPerVByteFeeMutex sync.Mutex
	satPerVByteFee      int64
}

func newLocalChain() *localChain {
	return &localChain{
		transactions: make(map[Hash]*Transaction),
	}
}

func (lc *localChain) GetTransaction(
	transactionHash Hash,
) (*Transaction, error) {
	lc.transactionsMutex.Lock()
	defer lc.transactionsMutex.Unlock()

	if transaction, exists := lc.transactions[transactionHash]; exists {
		return transaction, nil
	}

	return nil, fmt.Errorf("transaction not found")
}

func (lc *localChain) GetTransactionConfirmations(
	transactionHash Hash,
) (uint, error) {
	panic("not implemented")
}

func (lc *localChain) BroadcastTransaction(
	transaction *Transaction,
) error {
	panic("not implemented")
}

func (lc *localChain) GetLatestBlockHeight() (uint, error) {
	panic("not implemented")
}

func (lc *localChain) GetBlockHeader(
	blockNumber uint,
) (*BlockHeader, error) {
	panic("not implemented")
}

func (lc *localChain) GetTransactionMerkleProof(
	transactionHash Hash,
	blockHeight uint,
) (*TransactionMerkleProof, error) {
	panic("not implemented")
}

func (lc *localChain) GetTransactionsForPublicKeyHash(
	publicKeyHash [20]byte,
	limit int,
) ([]*Transaction, error) {
	panic("not implemented")
}

func (lc *localChain) GetMempoolForPublicKeyHash(
	publicKeyHash [20]byte,
) ([]*Transaction, error) {
	panic("not implemented")
}

func (lc *localChain) EstimateSatPerVByteFee(
	blocks uint32,
) (int64, error) {
	lc.satPerVByteFeeMutex.Lock()
	defer lc.satPerVByteFeeMutex.Unlock()

	return lc.satPerVByteFee, nil
}

func (lc *localChain) setSatPerVByteFee(
	satPerVByteFee int64,
) {
	lc.satPerVByteFeeMutex.Lock()
	defer lc.satPerVByteFeeMutex.Unlock()

	lc.satPerVByteFee = satPerVByteFee
}

func (lc *localChain) addTransaction(
	transaction *Transaction,
) error {
	lc.transactionsMutex.Lock()
	defer lc.transactionsMutex.Unlock()

	transactionHash := transaction.Hash()

	if _, exists := lc.transactions[transactionHash]; exists {
		return fmt.Errorf("transaction already exists")
	}

	lc.transactions[transactionHash] = transaction

	return nil
}
