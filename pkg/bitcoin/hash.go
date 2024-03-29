package bitcoin

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/keep-network/keep-core/pkg/internal/byteutils"
)

// HashByteLength is the byte length of the Hash type.
const HashByteLength = 32

// Hash represents the double SHA-256 of some arbitrary data using the
// InternalByteOrder.
type Hash [HashByteLength]byte

// NewHashFromString creates a new Hash instance using the given string.
// The string is interpreted according to the given ByteOrder. That is, the
// string is taken as is if the ByteOrder is InternalByteOrder and reversed if
// the ByteOrder is ReversedByteOrder. The string's length must be equal
// to 2*HashByteLength.
func NewHashFromString(hash string, byteOrder ByteOrder) (Hash, error) {
	if len(hash) != 2*HashByteLength {
		return Hash{}, fmt.Errorf("wrong hash string size")
	}

	hashBytes, err := hex.DecodeString(hash)
	if err != nil {
		return Hash{}, fmt.Errorf(
			"cannot decode hash string: [%w]",
			err,
		)
	}

	return NewHash(hashBytes, byteOrder)
}

// NewHash creates a new Hash instance using the given byte slice.
// The byte slice is interpreted according to the given ByteOrder. That is, the
// byte slice is taken as is if the ByteOrder is InternalByteOrder and reversed
// if the ByteOrder is ReversedByteOrder. The byte slice's length must be equal
// to HashByteLength.
func NewHash(hash []byte, byteOrder ByteOrder) (Hash, error) {
	if len(hash) != HashByteLength {
		return Hash{}, fmt.Errorf("wrong hash size")
	}

	var result Hash

	switch byteOrder {
	case InternalByteOrder:
		copy(result[:], hash[:])
	case ReversedByteOrder:
		reversed := byteutils.Reverse(hash)
		copy(result[:], reversed)
	default:
		panic("unknown byte order")
	}

	return result, nil
}

// ComputeHash computes the Hash for the provided data.
func ComputeHash(data []byte) Hash {
	first := sha256.Sum256(data)
	return sha256.Sum256(first[:])
}

// String returns the unprefixed hexadecimal string representation of the Hash
// in the InternalByteOrder.
func (h Hash) String() string {
	return h.Hex(InternalByteOrder)
}

// Hex returns the unprefixed hexadecimal string representation of the Hash
// in the given ByteOrder.
func (h Hash) Hex(byteOrder ByteOrder) string {
	switch byteOrder {
	case InternalByteOrder:
		return hex.EncodeToString(h[:])
	case ReversedByteOrder:
		for i := 0; i < HashByteLength/2; i++ {
			h[i], h[HashByteLength-1-i] = h[HashByteLength-1-i], h[i]
		}
		return hex.EncodeToString(h[:])
	default:
		panic("unknown byte order")
	}
}
