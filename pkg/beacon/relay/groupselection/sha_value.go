package groupselection

import (
	"crypto/sha256"
	"math/big"
)

// SHAValue is a wrapper type for a fixed-size byte array that contains an SHA
// signature. It can be represented as a byte slice (Bytes()), *big.Int (Int()),
// or the raw underlying fixed-size array (Raw()).
type SHAValue [sha256.Size]byte

// Bytes returns a byte slice of a copy of the SHAValue byte array.
func (v SHAValue) Bytes() []byte {
	var byteSlice []byte
	for _, byte := range v {
		byteSlice = append(byteSlice, byte)
	}
	return byteSlice
}

// Int returns a version of the byte array interpreted as a big.Int.
func (v SHAValue) Int() *big.Int {
	return big.NewInt(0).SetBytes(v.Bytes())
}

// Raw returns the underlying fixed sha256.Size-size byte array.
func (v SHAValue) Raw() [sha256.Size]byte {
	return v
}

// SetBytes takes the first 32 bytes from the provided byte slice and sets them
// as an internal value.
func (v SHAValue) SetBytes(bytes []byte) SHAValue {
	var container [sha256.Size]byte
	copy(container[:], bytes[0:sha256.Size])
	return container
}
