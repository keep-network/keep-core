package ethereum

import "testing"

func TestByteSliceToSliceOf1Byte(t *testing.T) {
	// func ByteSliceToSliceOf1Byte(in []byte) (rv [][1]byte) {
	var b []byte
	b = make([]byte, 3, 3)
	b[0] = 'a'
	b[1] = 'b'
	b[2] = 'c'
	rv := ByteSliceToSliceOf1Byte(b)
	if len(rv) != 3 {
		t.Errorf("Expected length of 3 got %d\n", len(rv))
	}
	if rv[2][0] != 'c' {
		t.Errorf("Expected 'c' got %v\n", rv[2][0])
	}

	// test the converstion back to to byte slice
	// func SliceOf1ByteToByteSlice(in [][1]byte) (rv []byte) {
	n := SliceOf1ByteToByteSlice(rv)

	if string(n) != string(b) {
		t.Errorf("Expected original [%s] to match with converted [%s]\n", b, n)
	}
}

func TestToByte32(t *testing.T) {
	var b []byte
	b = make([]byte, 32, 32)
	b[0] = 'a'
	b[1] = 'b'
	b[2] = 'c'
	// func ToByte32(in []byte) (tmp [32]byte) {
	rv, err := ToByte32(b)
	_ = err
	if len(rv) != 32 {
		t.Errorf("Expected length of 32 got %d\n", len(rv))
	}
}
