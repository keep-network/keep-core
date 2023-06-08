package bitcoin

import (
	"fmt"
)

var errNoBlocksSet = fmt.Errorf("blockchain does not contain any blocks")

// localSpvChain represents a local Bitcoin chain for SPV proof testing.
type localSpvChain struct {
	blockHeaders             map[uint]*BlockHeader
	transactions             map[Hash]*Transaction
	transactionMerkleProof   *TransactionMerkleProof
	transactionConfirmations map[Hash]uint
}

func newLocalSpvChain() *localSpvChain {
	return &localSpvChain{}
}

func (lsc *localSpvChain) GetTransaction(transactionHash Hash) (
	*Transaction,
	error,
) {
	transaction, found := lsc.transactions[transactionHash]
	if !found {
		return nil, fmt.Errorf(
			"transaction with hash %v does not exist", transactionHash,
		)
	}

	return transaction, nil
}

func (lsc *localSpvChain) GetTransactionConfirmations(
	transactionHash Hash,
) (uint, error) {
	return lsc.transactionConfirmations[transactionHash], nil
}

func (lsc *localSpvChain) BroadcastTransaction(
	transaction *Transaction,
) error {
	panic("unsupported")
}

func (lsc *localSpvChain) GetLatestBlockHeight() (uint, error) {
	blockchainTip := uint(0)
	for blockHeaderHeight := range lsc.blockHeaders {
		if blockHeaderHeight > blockchainTip {
			blockchainTip = blockHeaderHeight
		}
	}

	if blockchainTip == 0 {
		return 0, errNoBlocksSet
	}

	return blockchainTip, nil
}

func (lsc *localSpvChain) GetBlockHeader(blockNumber uint) (
	*BlockHeader,
	error,
) {
	blockHeader, found := lsc.blockHeaders[blockNumber]
	if !found {
		return nil, fmt.Errorf(
			"block header at height %v does not exist",
			blockNumber,
		)
	}

	return blockHeader, nil
}

func (lsc *localSpvChain) GetTransactionMerkleProof(
	transactionHash Hash,
	blockHeight uint,
) (*TransactionMerkleProof, error) {
	return lsc.transactionMerkleProof, nil
}

func (lsc *localSpvChain) GetTransactionsForPublicKeyHash(
	publicKeyHash [20]byte,
	limit int,
) ([]*Transaction, error) {
	panic("unsupported")
}

func (lsc *localSpvChain) GetMempoolForPublicKeyHash(
	publicKeyHash [20]byte,
) ([]*Transaction, error) {
	panic("unsupported")
}

func (lsc *localSpvChain) setBlockHeaders(blockHeaders map[uint]*BlockHeader) {
	lsc.blockHeaders = blockHeaders
}

func (lsc *localSpvChain) setTransactions(transactions map[Hash]*Transaction) {
	lsc.transactions = transactions
}

func (lsc *localSpvChain) setTransactionMerkleProof(
	transactionMerkleProof *TransactionMerkleProof,
) {
	lsc.transactionMerkleProof = transactionMerkleProof
}

func (lsc *localSpvChain) setTransactionConfirmations(
	transactionConfirmations map[Hash]uint,
) {
	lsc.transactionConfirmations = transactionConfirmations
}

func (lsc *localSpvChain) EstimateSatPerVByteFee(
	blocks uint32,
) (int64, error) {
	panic("unsupported")
}
