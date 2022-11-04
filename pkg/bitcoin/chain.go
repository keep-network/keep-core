package bitcoin

// Chain defines an interface meant to be used for interaction with the
// Bitcoin chain.
type Chain interface {
	// GetTransaction gets the transaction with the given transaction hash.
	// If the transaction with the given hash was not found on the chain,
	// this function returns an error.
	GetTransaction(transactionHash Hash) (*Transaction, error)

	// GetTransactionConfirmations gets the number of confirmations for the
	// transaction with the given transaction hash. If the transaction with the
	// given hash was not found on the chain, this function returns an error.
	GetTransactionConfirmations(transactionHash Hash) (uint, error)

	// BroadcastTransaction broadcasts the given transaction over the
	// network of the Bitcoin chain nodes. If the broadcast action could not be
	// done, this function returns an error. This function does not give any
	// guarantees regarding transaction mining. The transaction may be mined or
	// rejected eventually.
	BroadcastTransaction(transaction *Transaction) error

	// GetLatestBlockHeight gets the height of the latest block (tip). If the
	// latest block was not determined, this function returns an error.
	GetLatestBlockHeight() (uint, error)

	// GetBlockHeader gets the block header for the given block height. If the
	// block with the given height was not found on the chain, this function
	// returns an error.
	GetBlockHeader(blockHeight uint) (*BlockHeader, error)
}
