package group

import (
	"math/big"
)

// MemberIndex is an index of a member in a group.
type MemberIndex uint8

// Int converts `MemberIndex` to `big.Int`.
func (id MemberIndex) Int() *big.Int {
	return new(big.Int).SetUint64(uint64(id))
}

// Equals checks if `MemberIndex` equals the passed int value.
func (id MemberIndex) Equals(value int) bool {
	return id == MemberIndex(value)
}
