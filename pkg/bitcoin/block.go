package bitcoin

import "encoding/binary"

// BlockHeaderByteLength is the byte length of a serialized block header.
const BlockHeaderByteLength = 80

// BlockHeader represents the header of a Bitcoin block. For reference, see:
// https://developer.bitcoin.org/reference/block_chain.html#block-headers
type BlockHeader struct {
	// Version is the block version number that indicates which set of block
	// validation rules to follow.
	Version int32
	// PreviousBlockHeaderHash is the hash of the previous block's header.
	PreviousBlockHeaderHash Hash
	// MerkleRootHash is a hash derived from the hashes of all transactions
	// included in this block.
	MerkleRootHash Hash
	// Time is a Unix epoch time when the miner started hashing the header.
	Time uint32
	// Bits determines the target threshold this block's header hash must be
	// less than or equal to.
	Bits uint32
	// Nonce is an arbitrary number miners change to modify the header hash
	// in order to produce a hash less than or equal to the target threshold.
	Nonce uint32
}

// Serialize serializes the block header to a byte array using the block header
// serialization format:
// [Version][PreviousBlockHeaderHash][MerkleRootHash][Time][Bits][Nonce].
func (bh *BlockHeader) Serialize() [BlockHeaderByteLength]byte {
	var result [BlockHeaderByteLength]byte
	offset := 0

	// Version
	binary.LittleEndian.PutUint32(result[offset:], uint32(bh.Version))
	offset += 4

	// PreviousBlockHeaderHash
	copy(result[offset:], bh.PreviousBlockHeaderHash[:])
	offset += len(bh.PreviousBlockHeaderHash)

	// MerkleRootHash
	copy(result[offset:], bh.MerkleRootHash[:])
	offset += len(bh.MerkleRootHash)

	// Time
	binary.LittleEndian.PutUint32(result[offset:], bh.Time)
	offset += 4

	// Bits
	binary.LittleEndian.PutUint32(result[offset:], bh.Bits)
	offset += 4

	// Nonce
	binary.LittleEndian.PutUint32(result[offset:], bh.Nonce)

	return result
}

// Hash calculates the block header's hash as the double SHA-256 of the
// block header serialization format:
// [Version][PreviousBlockHeaderHash][MerkleRootHash][Time][Bits][Nonce].
func (bh *BlockHeader) Hash() Hash {
	// TODO: Implementation of the Hash function that consists of the following:
	//       1. Call bh.Serialize() to get the serialized block hash.
	//       2. Compute the double SHA-256 over the serialized  block hash.
	//       3. Construct the Hash instance appropriately.
	return Hash{}
}

// Proof contains data required to perform a proof that a given transaction was
// included in the Bitcoin blockchain.
type Proof struct { // TODO: Possibly make it contain transaction data
	// MerkleProof is the Merkle proof of transaction inclusion in a block.
	MerkleProof []byte

	// TxIndexInBlock is the transaction index in the block (0-indexed).
	TxIndexInBlock uint

	// BitcoinHeaders is a chain of block headers that form confirmations of
	// blockchain inclusion.
	BitcoinHeaders []byte
}
