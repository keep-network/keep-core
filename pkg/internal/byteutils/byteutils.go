// Package byteutils provides helper utilities for working with bytes
package byteutils

import (
	"fmt"
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

// Reverse reverses bytes order in the slice.
func Reverse(bytes []byte) []byte {
	result := make([]byte, 0, len(bytes))
	for i := len(bytes) - 1; i >= 0; i-- {
		result = append(result, bytes[i])
	}
	return result
}
