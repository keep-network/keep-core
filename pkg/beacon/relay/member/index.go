package member

import (
	"fmt"
	"math/big"
)

// MemberIndex is an index of a member in a group.
type MemberIndex uint32

// Int converts `MemberIndex` to `big.Int`.
func (id MemberIndex) Int() *big.Int {
	return new(big.Int).SetUint64(uint64(id))
}

// Equals checks if `MemberIndex` equals the passed int value.
func (id MemberIndex) Equals(value int) bool {
	return id == MemberIndex(value)
}

// Validate checks if `MemberIndex` has a valid value. `MemberIndex` is expected
// to be equal or greater than `1`.
func (id MemberIndex) Validate() error {
	if id < 1 {
		return fmt.Errorf("member index must be >= 1")
	}
	return nil
}
