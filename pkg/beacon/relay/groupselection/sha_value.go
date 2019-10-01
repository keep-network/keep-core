package groupselection

import (
	"crypto/sha256"
	"fmt"
	"math/big"
)

// shaValue is a wrapper type for a fixed-size byte array that contains an SHA
// signature. It can be represented as a byte slice (Bytes()), *big.Int (Int()),
// or the raw underlying fixed-size array (Raw()).
type shaValue [sha256.Size]byte

// bytes returns a byte slice of a copy of the SHAValue byte array.
func (v shaValue) bytes() []byte {
	var byteSlice []byte
	for _, byte := range v {
		byteSlice = append(byteSlice, byte)
	}
	return byteSlice
}

// int returns a version of the byte array interpreted as a big.Int.
func (v shaValue) int() *big.Int {
	return new(big.Int).SetBytes(v.bytes())
}

// raw returns the underlying fixed sha256.Size-size byte array.
func (v shaValue) raw() [sha256.Size]byte {
	return v
}

// setBytes takes 32 bytes from the provided byte slice and sets them as an
// internal value. If slice length is different than 32 bytes it returns an error.
func (v shaValue) setBytes(bytes []byte) (shaValue, error) {
	var container [sha256.Size]byte

	if len(bytes) != sha256.Size {
		return container, fmt.Errorf("%v bytes expected for SHA value", sha256.Size)
	}

	copy(container[:], bytes[0:sha256.Size])

	return container, nil
}
