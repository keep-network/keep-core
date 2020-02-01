package group

import (
	"math/big"
)

// MemberIndex is an index of a member in a group. The maximum member index
// value is 255.
type MemberIndex uint8

// Uint8 converts MemberIndex to uint8 without losing any precision.
func (id MemberIndex) Uint8() uint8 {
	return uint8(id)
}

// Int converts `MemberIndex` to `big.Int`.
func (id MemberIndex) Int() *big.Int {
	return new(big.Int).SetUint64(uint64(id))
}

// Equals checks if `MemberIndex` equals the passed int value.
func (id MemberIndex) Equals(value int) bool {
	return id == MemberIndex(value)
}
