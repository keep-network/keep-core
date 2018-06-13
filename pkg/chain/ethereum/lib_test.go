package ethereum

import "testing"

func TestByteSliceToSliceOf1Byte(t *testing.T) {
	var b []byte
	b = make([]byte, 3, 3)
	b[0] = 'a'
	b[1] = 'b'
	b[2] = 'c'
	rv := byteSliceToSliceoF1Byte(b)
	if len(rv) != 3 {
		t.Errorf("Expected length of 3 got %d\n", len(rv))
	}
	if rv[2][0] != 'c' {
		t.Errorf("Expected 'c' got %v\n", rv[2][0])
	}

	// test the converstion back to to byte slice.
	n := sliceOf1ByteToByteSlice(rv)

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
	rv, err := toByte32(b)
	if err != nil {
		t.Errorf("function toByte32 reported an error [%v]\n", err)
	}
	if len(rv) != 32 {
		t.Errorf("expected length of 32 got %d\n", len(rv))
	}
}

func TestToByte32LessThan32(t *testing.T) {
	var b []byte
	b = make([]byte, 3, 3)
	b[0] = 'a'
	b[1] = 'b'
	b[2] = 'c'
	_, err := toByte32(b)
	if err == nil {
		t.Errorf("function toByte32 failed to report an error\n")
	}
}
