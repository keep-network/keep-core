package bitcoin

import (
	"encoding/hex"
	"fmt"
	"github.com/keep-network/keep-core/internal/testutils"
	"reflect"
	"testing"
)

func TestReadCompactSizeUint(t *testing.T) {
	fromHex := func(hexString string) []byte {
		bytes, err := hex.DecodeString(hexString)
		if err != nil {
			t.Fatal(err)
		}
		return bytes
	}

	var tests = map[string]struct {
		data               []byte
		expectedValue      CompactSizeUint
		expectedByteLength int
		expectedErr        error
	}{
		"1-byte compact size uint": {
			data:               fromHex("bb"),
			expectedValue:      187,
			expectedByteLength: 1,
		},
		"3-byte compact size uint": {
			data:               fromHex("fd0302"),
			expectedValue:      515,
			expectedByteLength: 3,
		},
		"5-byte compact size uint": {
			data:               fromHex("fe703a0f00"),
			expectedValue:      998000,
			expectedByteLength: 5,
		},
		"9-byte compact size uint": {
			data:               fromHex("ff57284e56dab40000"),
			expectedValue:      198849843832919,
			expectedByteLength: 9,
		},
		"malformed compact size uint": {
			data:        fromHex("fd01"),
			expectedErr: fmt.Errorf("unexpected EOF"),
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			value, byteLength, err := readCompactSizeUint(test.data)

			if !reflect.DeepEqual(test.expectedErr, err) {
				t.Errorf(
					"unexpected error\nexpected: %+v\nactual:   %+v\n",
					test.expectedErr,
					err,
				)
			}

			testutils.AssertIntsEqual(
				t,
				"value",
				int(test.expectedValue),
				int(value),
			)

			testutils.AssertIntsEqual(
				t,
				"byte length",
				test.expectedByteLength,
				byteLength,
			)
		})
	}
}

func TestWriteCompactSizeUint(t *testing.T) {
	fromHex := func(hexString string) []byte {
		bytes, err := hex.DecodeString(hexString)
		if err != nil {
			t.Fatal(err)
		}
		return bytes
	}

	var tests = map[string]struct {
		csu           CompactSizeUint
		expectedBytes []byte
	}{
		"1-byte compact size uint": {
			csu:           187,
			expectedBytes: fromHex("bb"),
		},
		"3-byte compact size uint": {
			csu:           515,
			expectedBytes: fromHex("fd0302"),
		},
		"5-byte compact size uint": {
			csu:           998000,
			expectedBytes: fromHex("fe703a0f00"),
		},
		"9-byte compact size uint": {
			csu:           198849843832919,
			expectedBytes: fromHex("ff57284e56dab40000"),
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			bytes, err := writeCompactSizeUint(test.csu)
			if err != nil {
				t.Fatal(err)
			}

			testutils.AssertBytesEqual(t, test.expectedBytes, bytes)
		})
	}
}
