package bitcoin

import (
	"encoding/binary"
	"math/big"
)

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

// Deserialize deserializes a byte array to a BlockHeader using the block header
// serialization format:
// [Version][PreviousBlockHeaderHash][MerkleRootHash][Time][Bits][Nonce].
func (bh *BlockHeader) Deserialize(rawBlockHeader [BlockHeaderByteLength]byte) {
	offset := 0

	// Version
	bh.Version = int32(binary.LittleEndian.Uint32(rawBlockHeader[offset:]))
	offset += 4

	// PreviousBlockHeaderHash
	copy(
		bh.PreviousBlockHeaderHash[:],
		rawBlockHeader[offset:offset+HashByteLength],
	)
	offset += HashByteLength

	// MerkleRootHash
	copy(bh.MerkleRootHash[:], rawBlockHeader[offset:offset+HashByteLength])
	offset += HashByteLength

	// Time
	bh.Time = binary.LittleEndian.Uint32(rawBlockHeader[offset:])
	offset += 4

	// Bits
	bh.Bits = binary.LittleEndian.Uint32(rawBlockHeader[offset:])
	offset += 4

	// Nonce
	bh.Nonce = binary.LittleEndian.Uint32(rawBlockHeader[offset:])
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

// Target calculates the difficulty target of a block header.
func (bh *BlockHeader) Target() *big.Int {
	// A serialized 80-byte block header stores the `bits` value as a 4-byte
	// little-endian hexadecimal value in a slot including bytes 73, 74, 75, and
	// 76. This function's input argument is expected to be a numerical
	// representation of that 4-byte value reverted to the big-endian order.
	// For example, if the `bits` little-endian value in the header is
	// `0xcb04041b`, it must be reverted to the big-endian form `0x1b0404cb` and
	// turned to a decimal number `453248203` in order to be used as this
	// function's input.
	//
	// The `bits` 4-byte big-endian representation is a compact value that works
	// like a base-256 version of scientific notation. It encodes the target
	// exponent in the first byte and the target mantissa in the last three bytes.
	// Referring to the previous example, if `bits = 453248203`, the hexadecimal
	// representation is `0x1b0404cb` so the exponent is `0x1b` while the mantissa
	// is `0x0404cb`.
	//
	// To extract the exponent, we need to shift right by 3 bytes (24 bits),
	// extract the last byte of the result, and subtract 3 (because of the
	// mantissa length):
	// - 0x1b0404cb >>> 24 = 0x0000001b
	// - 0x0000001b & 0xff = 0x1b
	// - 0x1b - 3 = 24 (decimal)
	//
	// To extract the mantissa, we just need to take the last three bytes:
	// - 0x1b0404cb & 0xffffff = 0x0404cb = 263371 (decimal)
	//
	// The final difficulty can be computed as mantissa * 256^exponent:
	// - 263371 * 256^24 =
	// 1653206561150525499452195696179626311675293455763937233695932416 (decimal)
	//
	// Sources:
	// - https://developer.bitcoin.org/reference/block_chain.html#target-nbits
	// - https://wiki.bitcoinsv.io/index.php/Target
	exponent := (int(bh.Bits>>24) & 0xff) - 3
	mantissa := bh.Bits & 0xffffff

	// Compute 256^exponent.
	pow := new(big.Int).Exp(
		big.NewInt(256),
		big.NewInt(int64(exponent)),
		nil,
	)

	// Compute mantissa * (256^exponent).
	target := new(big.Int).Mul(big.NewInt(int64(mantissa)), pow)

	return target
}

// Difficulty calculates the difficulty of a block header.
func (bh *BlockHeader) Difficulty() *big.Int {
	maxTarget := new(big.Int)
	maxTarget.SetString(
		"00000000ffff0000000000000000000000000000000000000000000000000000",
		16,
	)

	target := bh.Target()

	difficulty := new(big.Int)
	difficulty.Div(maxTarget, target)

	return difficulty
}
