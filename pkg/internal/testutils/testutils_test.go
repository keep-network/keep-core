package testutils

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"

	crand "crypto/rand"
)

var byteArray = [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

func TestAssertBytesEqual(t *testing.T) {
	var tests = map[string]struct {
		expectedBytes []byte
		actualBytes   []byte
		expectedError error
	}{
		"same bytes": {
			expectedBytes: make([]byte, 5, 5),
			actualBytes:   make([]byte, 5, 5),
			expectedError: nil,
		},
		"empty bytes": {
			expectedBytes: make([]byte, 0),
			actualBytes:   make([]byte, 0),
			expectedError: nil,
		},
		"empty vs present bytes": {
			expectedBytes: make([]byte, 0),
			actualBytes:   byteArray[0:5],
			expectedError: fmt.Errorf(
				"%s\nexpected: [%s]\nactual:   [%s]",
				"Byte slices differ in 5 places",
				"[]",
				"[0 1 2 3 4]",
			),
		},
		"present vs empty bytes": {
			expectedBytes: byteArray[0:5],
			actualBytes:   make([]byte, 0),
			expectedError: fmt.Errorf(
				"%s\nexpected: [%s]\nactual:   [%s]",
				"Byte slices differ in 5 places",
				"[0 1 2 3 4]",
				"[]",
			),
		},
		"same but longer bytes": {
			expectedBytes: make([]byte, 3, 3),
			actualBytes:   make([]byte, 6, 6),
			expectedError: fmt.Errorf(
				"%s\nexpected: [%s]\nactual:   [%s]",
				"Byte slices differ in 3 places",
				"[0 0 0]",
				"[0 0 0 0 0 0]",
			),
		},
		"same but shorter bytes": {
			expectedBytes: make([]byte, 6, 6),
			actualBytes:   make([]byte, 3, 3),
			expectedError: fmt.Errorf(
				"%s\nexpected: [%s]\nactual:   [%s]",
				"Byte slices differ in 3 places",
				"[0 0 0 0 0 0]",
				"[0 0 0]",
			),
		},
		"different bytes, same length": {
			expectedBytes: byteArray[0:5],
			actualBytes:   byteArray[5:10],
			expectedError: fmt.Errorf(
				"%s\nexpected: [%s]\nactual:   [%s]",
				"Byte slices differ in 5 places",
				"[0 1 2 3 4]",
				"[5 6 7 8 9]",
			),
		},
		"different, longer bytes": {
			expectedBytes: byteArray[0:5],
			actualBytes:   byteArray[5:8],
			expectedError: fmt.Errorf(
				"%s\nexpected: [%s]\nactual:   [%s]",
				"Byte slices differ in 5 places",
				"[0 1 2 3 4]",
				"[5 6 7]",
			),
		},
		"different, shorter bytes": {
			expectedBytes: byteArray[0:3],
			actualBytes:   byteArray[5:10],
			expectedError: fmt.Errorf(
				"%s\nexpected: [%s]\nactual:   [%s]",
				"Byte slices differ in 5 places",
				"[0 1 2]",
				"[5 6 7 8 9]",
			),
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			err := testBytesEqual(test.expectedBytes, test.actualBytes)
			if !reflect.DeepEqual(err, test.expectedError) {
				t.Errorf(
					"\nexpected: [%v]\nactual:   [%v]",
					test.expectedError,
					err,
				)
			}
		})
	}
}

func TestMockRandReader(t *testing.T) {
	counter := big.NewInt(10)
	mockRandReader := NewMockRandReader(counter)
	for i := 0; i <= 5; i++ {
		expectedNumber := new(big.Int).Add(counter, big.NewInt(int64(i)))
		randomNumber, err := crand.Int(mockRandReader, big.NewInt(1000))
		if err != nil {
			t.Fatalf("unexpected error [%s]", err)
		}
		if randomNumber.Cmp(expectedNumber) != 0 {
			t.Fatalf("generated random number %s doesn't match expected %s", randomNumber, expectedNumber)
		}
	}
}
