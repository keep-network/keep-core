package ethereum

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

type FxError func(err error) (eout error)

const db1 = true
