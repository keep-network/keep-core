package bitcoin

// SpvProof contains data required to perform a proof that a given transaction
// was included in the Bitcoin blockchain.
type SpvProof struct {
	// MerkleProof is the Merkle proof of transaction inclusion in a block.
	MerkleProof []byte

	// TxIndexInBlock is the transaction index in the block (0-indexed).
	TxIndexInBlock uint

	// BitcoinHeaders is a chain of block headers that form confirmations of
	// blockchain inclusion.
	BitcoinHeaders []byte
}
