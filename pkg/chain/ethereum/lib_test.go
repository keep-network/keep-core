package ethereum

import "testing"

func TestByteSliceToSliceOf1Byte(t *testing.T) {
	var b []byte
	b = make([]byte, 3, 3)
	b[0] = 'a'
	b[1] = 'b'
	b[2] = 'c'
	rv := byteSliceToSliceOf1Byte(b)
	if len(rv) != 3 {
		t.Errorf("Expected length of 3 got %d\n", len(rv))
	}
	if rv[2][0] != 'c' {
		t.Errorf("Expected 'c' got %v\n", rv[2][0])
	}

	// test the converstion back to to byte slice.
	n := sliceOf1ByteToByteSlice(rv)

	if string(n) != string(b) {
		t.Errorf("Expected original [%s] to match with converted [%s]", b, n)
	}
}

func TestToByte32(t *testing.T) {
	tests := []struct {
		nOfBytes    int
		expectError bool
	}{
		{
			nOfBytes:    32,
			expectError: false,
		},
		{
			nOfBytes:    12,
			expectError: true,
		},
		{
			nOfBytes:    42,
			expectError: true,
		},
	}

	var b []byte

	for testIndex, test := range tests {
		b = make([]byte, test.nOfBytes, test.nOfBytes)
		b[0] = 'a'
		b[1] = 'b'
		b[2] = 'c'
		rv, err := toByte32(b)
		if test.expectError {
			if err == nil {
				t.Errorf("function toByte32 failed to report an error, test %d", testIndex)
			}
		} else {
			if err != nil {
				t.Errorf("function toByte32 reported an error [%v], test %d", err, testIndex)
			}
			if len(rv) != 32 {
				t.Errorf("expected length of 32 got %d, test %d\n", len(rv), testIndex)
			}
		}
	}
}
