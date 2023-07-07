package spv

import (
	"bytes"
	"fmt"
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"sync"
)

type localBitcoinChain struct {
	mutex sync.Mutex

	transactions []*bitcoin.Transaction
}

func newLocalBitcoinChain() *localBitcoinChain {
	return &localBitcoinChain{
		transactions: make([]*bitcoin.Transaction, 0),
	}
}

func (lbc *localBitcoinChain) GetTransaction(transactionHash bitcoin.Hash) (
	*bitcoin.Transaction,
	error,
) {
	lbc.mutex.Lock()
	defer lbc.mutex.Unlock()

	for _, transaction := range lbc.transactions {
		if transaction.Hash() == transactionHash {
			return transaction, nil
		}
	}

	return nil, fmt.Errorf("transaction not found")
}

func (lbc *localBitcoinChain) GetTransactionConfirmations(transactionHash bitcoin.Hash) (
	uint,
	error,
) {
	panic("unsupported")
}

func (lbc *localBitcoinChain) BroadcastTransaction(transaction *bitcoin.Transaction) error {
	lbc.mutex.Lock()
	defer lbc.mutex.Unlock()

	transactionHash := transaction.Hash()

	for _, existingTransaction := range lbc.transactions {
		if transactionHash == existingTransaction.Hash() {
			return fmt.Errorf("transaction already exists")
		}
	}

	lbc.transactions = append(lbc.transactions, transaction)

	return nil
}

func (lbc *localBitcoinChain) GetLatestBlockHeight() (uint, error) {
	panic("unsupported")
}

func (lbc *localBitcoinChain) GetBlockHeader(blockHeight uint) (
	*bitcoin.BlockHeader,
	error,
) {
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
	lbc.mutex.Lock()
	defer lbc.mutex.Unlock()

	p2pkh, err := bitcoin.PayToPublicKeyHash(publicKeyHash)
	if err != nil {
		return nil, err
	}

	p2wpkh, err := bitcoin.PayToWitnessPublicKeyHash(publicKeyHash)
	if err != nil {
		return nil, err
	}

	matchingTransactions := make([]*bitcoin.Transaction, 0)

	for _, transaction := range lbc.transactions {
		for _, output := range transaction.Outputs {
			script := output.PublicKeyScript
			if bytes.Equal(script, p2pkh) || bytes.Equal(script, p2wpkh) {
				matchingTransactions = append(matchingTransactions, transaction)
				break
			}
		}
	}

	if len(matchingTransactions) > limit {
		return matchingTransactions[len(matchingTransactions)-limit:], nil
	}

	return matchingTransactions, nil
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
