package ethereum_v1

import (
	"fmt"
)

// byteSliceToSliceOf1Byte convert from slice to slice of arrays of 1
// long bytes.
func byteSliceToSliceOf1Byte(in []byte) (rv [][1]byte) {
	if len(in) == 0 {
		rv = make([][1]byte, 0, 1)
		return
	}
	rv = make([][1]byte, 0, len(in))
	for _, vv := range in {
		rv = append(rv, [1]byte{vv})
	}
	return
}

// sliceOf1ByteToByteSlice convert from solidity type, slice of 1 long
// array of bytes to byte slice.
func sliceOf1ByteToByteSlice(in [][1]byte) (rv []byte) {
	if len(in) == 0 {
		rv = make([]byte, 0, 1)
		return
	}
	rv = make([]byte, 0, len(in))
	for _, vv := range in {
		rv = append(rv, vv[0])
	}
	return
}

// toByte32 convert from byte slice to fixed length array of 32 long.
func toByte32(in []byte) (tmp [32]byte, err error) {
	if len(in) != 32 {
		return tmp, fmt.Errorf(
			"cannot convert slice of length %d to [32]byte, must be of length 32",
			len(in),
		)
	}
	for i := 0; i < len(in); i++ {
		tmp[i] = in[i]
	}
	return
}
