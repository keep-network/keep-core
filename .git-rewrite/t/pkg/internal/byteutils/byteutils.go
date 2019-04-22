// Package byteutils provides helper utilities for working with bytes
package byteutils

import (
	"crypto/sha256"
	"fmt"
	"math/big"
)

// LeftPadTo32Bytes appends zeros to bytes slice to make it exactly 32 bytes long
func LeftPadTo32Bytes(bytes []byte) ([]byte, error) {
	expectedByteLen := 32
	if len(bytes) > expectedByteLen {
		return nil, fmt.Errorf(
			"cannot pad %v byte array to %v bytes", len(bytes), expectedByteLen,
		)
	}

	result := make([]byte, 0)
	if len(bytes) < expectedByteLen {
		result = make([]byte, expectedByteLen-len(bytes))
	}
	result = append(result, bytes...)

	return result, nil
}

// Sha256Sum calculates sha256 hash for passed `bytes` and converts it to `big.Int`.
func Sha256Sum(bytes []byte) *big.Int {
	hash := sha256.Sum256(bytes)

	return new(big.Int).SetBytes(hash[:])
}
