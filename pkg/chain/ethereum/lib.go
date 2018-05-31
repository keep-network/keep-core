package ethereum

// ByteSliceToSliceOf1Byte convert from slice to slice of arrays of 1 long bytes.
func ByteSliceToSliceOf1Byte(in []byte) (rv [][1]byte) {
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

// SliceOf1ByteToByteSlice convert form solidity type, slice of 1 long array of bytes to byte slice.
func SliceOf1ByteToByteSlice(in [][1]byte) (rv []byte) {
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

// ToByte32 convert from byte slice to fixed length array of 32 long
func ToByte32(in []byte) (tmp [32]byte) {
	var i int
	for ; i < len(in); i++ {
		tmp[i] = in[i]
	}
	for ; i < 32; i++ {
		tmp[i] = byte(0)
	}
	return
}

// FxError is the function that is call when an error occures in a callback
type FxError func(err error) (eout error)

const db1 = true
