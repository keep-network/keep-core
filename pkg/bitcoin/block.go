package bitcoin

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
	// Time is a Unix epoch time when the block was created.
	Time uint32
	// Bits determines the target threshold this block's header hash must be
	// less than or equal to.
	Bits uint32
	// Nonce is the nonce used to generate the block.
	Nonce uint32
}

// Serialize serializes the block header to a byte array using the block header
// serialization format:
// [Version][PreviousBlockHeaderHash][MerkleRootHash][Time][Bits][Nonce].
func (bh *BlockHeader) Serialize() [BlockHeaderByteLength]byte {
	// TODO: Implementation of the Serialize function that concatenates all
	//       serialized fields of the block header into a single byte array.
	//       All numbers should be serialized to an InternalByteOrder byte array.
	return [BlockHeaderByteLength]byte{}
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
