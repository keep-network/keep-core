package ethereum

import (
	"fmt"
	"reflect"
	"testing"
)

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
	tests := map[string]struct {
		nOfBytes    int
		expectError bool
		errorFormat string
	}{
		"test expected length of 32 bytes": {
			nOfBytes:    32,
			expectError: false,
			errorFormat: "",
		},
		"test to short, only 12 long": {
			nOfBytes:    12,
			expectError: true,
			errorFormat: "cannot convert slice of length %d to [32]byte, must be of length 32",
		},
		"test too long, more than 32 length": {
			nOfBytes:    42,
			expectError: true,
			errorFormat: "cannot convert slice of length %d to [32]byte, must be of length 32",
		},
	}

	var b []byte
	var expectedError error

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			if test.nOfBytes < 3 {
				t.Fatalf("Test data incorrect, must be length 3 or more, have %d", test.nOfBytes)
			}
			b = make([]byte, test.nOfBytes, test.nOfBytes)
			b[0] = 'a'
			b[1] = 'b'
			b[2] = 'c'
			rv, err := toByte32(b)
			expectedError = nil
			if test.errorFormat != "" {
				expectedError = fmt.Errorf(test.errorFormat, test.nOfBytes)
			}
			if !reflect.DeepEqual(err, expectedError) {
				t.Fatalf("\nexpected: %v\nactual:   %v", expectedError, err)
			}
			if err == nil {
				if len(rv) != 32 {
					t.Errorf("\nexpected: 32 \nactual:   %d", len(rv))
				}
				if rv[0] != 'a' {
					t.Errorf("\nexpected: 'a' \nactual:   %v", rv[0])
				}
				if rv[1] != 'b' {
					t.Errorf("\nexpected: 'b' \nactual:   %v", rv[1])
				}
				if rv[2] != 'c' {
					t.Errorf("\nexpected: 'c' \nactual:   %v", rv[2])
				}
			}
		})
	}
}
