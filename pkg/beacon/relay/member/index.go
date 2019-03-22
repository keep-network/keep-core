package member

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"strconv"
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

// IndexFromBytes returns a `Index` created from provided bytes.
func IndexFromBytes(bytes []byte) Index {
	return Index(binary.LittleEndian.Uint32(bytes))
}

// Bytes converts `Index` to bytes representation.
func (id Index) Bytes() []byte {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, uint32(id))
	return bytes
}

// IndexFromHex returns a `Index` created from the hex `string`
// representation.
func IndexFromHex(hex string) (Index, error) {
	id, err := strconv.ParseUint(hex, 16, 32)
	if err != nil {
		return 0, err
	}

	return Index(id), nil
}

// HexString converts `Index` to hex `string` representation.
func (id Index) HexString() string {
	return strconv.FormatInt(int64(id), 16)
}
