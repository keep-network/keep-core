package bitcoin

import "fmt"

type localChain struct {
	transactions map[Hash]*Transaction
}

func newLocalChain() *localChain {
	return &localChain{
		transactions: make(map[Hash]*Transaction),
	}
}

func (lc *localChain) GetTransaction(
	transactionHash Hash,
) (*Transaction, error) {
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

func (lc *localChain) addTransaction(
	transaction *Transaction,
) error {
	transactionHash := transaction.Hash()

	if _, exists := lc.transactions[transactionHash]; exists {
		return fmt.Errorf("transaction already exists")
	}

	lc.transactions[transactionHash] = transaction

	return nil
}
