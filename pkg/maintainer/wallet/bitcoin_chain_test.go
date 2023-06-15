package wallet

import (
	"github.com/keep-network/keep-core/pkg/bitcoin"
)

type localBitcoinChain struct {
}

func newLocalBitcoinChain() *localBitcoinChain {
	return &localBitcoinChain{}
}

func (lbc *localBitcoinChain) GetTransaction(
	transactionHash bitcoin.Hash,
) (*bitcoin.Transaction, error) {
	panic("unsupported")
}

func (lbc *localBitcoinChain) GetTransactionConfirmations(
	transactionHash bitcoin.Hash,
) (uint, error) {
	panic("unsupported")
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
	panic("unsupported")
}
