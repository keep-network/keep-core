package member

import (
	"fmt"
	"math/big"
)

// Index is an index of a member in a group.
type Index uint32

// Int converts `Index` to `big.Int`.
func (id Index) Int() *big.Int {
	return new(big.Int).SetUint64(uint64(id))
}

// Equals checks if `Index` equals the passed int value.
func (id Index) Equals(value int) bool {
	return id == Index(value)
}

// Validate checks if `Index` has a valid value. `Index` is expected to be equal
// or greater than `1`.
func (id Index) Validate() error {
	if id < 1 {
		return fmt.Errorf("member index must be >= 1")
	}
	return nil
}
