package chain

import (
	"bytes"
	"testing"

	"github.com/keep-network/go-ethereum/common"
)

func TestDKGResultEquals(t *testing.T) {
	var tests = map[string]struct {
		result1        *DKGResult
		result2        *DKGResult
		expectedResult bool
	}{
		"both nil": {
			result1:        nil,
			result2:        nil,
			expectedResult: true,
		},
		"both empty": {
			result1:        &DKGResult{},
			result2:        &DKGResult{},
			expectedResult: true,
		},
		"nil and empty": {
			result1:        nil,
			result2:        &DKGResult{},
			expectedResult: false,
		},
		"empty and nil": {
			result1:        &DKGResult{},
			result2:        nil,
			expectedResult: false,
		},
		"success - equal": {
			result1:        &DKGResult{Success: true},
			result2:        &DKGResult{Success: true},
			expectedResult: true,
		},
		"success - not equal": {
			result1:        &DKGResult{Success: true},
			result2:        &DKGResult{Success: false},
			expectedResult: false,
		},
		"group public keys - equal": {
			result1:        &DKGResult{GroupPublicKey: []byte{2}},
			result2:        &DKGResult{GroupPublicKey: []byte{2}},
			expectedResult: true,
		},
		"group public keys - not equal": {
			result1:        &DKGResult{GroupPublicKey: []byte{3}},
			result2:        &DKGResult{GroupPublicKey: []byte{4}},
			expectedResult: false,
		},
		"disqualified - equal": {
			result1:        &DKGResult{Disqualified: []byte{0x00, 0x00, 0x01}},
			result2:        &DKGResult{Disqualified: []byte{0x00, 0x00, 0x01}},
			expectedResult: true,
		},
		"disqualified - other members DQ": {
			result1:        &DKGResult{Disqualified: []byte{0x00, 0x00, 0x01}},
			result2:        &DKGResult{Disqualified: []byte{0x00, 0x01, 0x00}},
			expectedResult: false,
		},
		"disqualified - different length of DQ members": {
			result1:        &DKGResult{Disqualified: []byte{0x00, 0x00, 0x00}},
			result2:        &DKGResult{Disqualified: []byte{0x00, 0x00}},
			expectedResult: false,
		},
		"inactive - equal": {
			result1:        &DKGResult{Inactive: []byte{0x01, 0x01, 0x00}},
			result2:        &DKGResult{Inactive: []byte{0x01, 0x01, 0x00}},
			expectedResult: true,
		},
		"inactive - other members IA": {
			result1:        &DKGResult{Inactive: []byte{0x01, 0x01, 0x00}},
			result2:        &DKGResult{Inactive: []byte{0x01, 0x01, 0x01}},
			expectedResult: false,
		},
		"inactive - different length of IA members": {
			result1:        &DKGResult{Inactive: []byte{0x00, 0x00}},
			result2:        &DKGResult{Inactive: []byte{0x00, 0x00, 0x00}},
			expectedResult: false,
		},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			actualResult := test.result1.Equals(test.result2)
			if test.expectedResult != actualResult {
				t.Fatalf("\nexpected: %v\nactual:   %v\n", test.expectedResult, actualResult)
			}
		})
	}
}

func TestEncodeAndHash(t *testing.T) {
	var tests = map[string]struct {
		dkgResult                *DKGResult
		expectedSerializedResult string
		expectedHash             string
	}{
		"dkg result has only success provided": {
			dkgResult: &DKGResult{
				Success: true,
			},
			expectedSerializedResult: "0000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000c0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			expectedHash:             "3550dcd2ae8c05d5ca00d5fcd004b6aaf7d3b90f64e0e6078ae62f1f7a2e3f27",
		},
		"dkg result has only group public key provided": {
			dkgResult: &DKGResult{
				GroupPublicKey: []byte{100},
			},
			expectedSerializedResult: "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000000e00000000000000000000000000000000000000000000000000000000000000001640000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			expectedHash:             "132d7294d5d6a5434295d5f5d611bd210261f4e1360cfeb1054e30e69747ea16",
		},
		"dkg result has only disqualified provided": {
			dkgResult: &DKGResult{
				Disqualified: []byte{1, 0, 1, 0},
			},
			expectedSerializedResult: "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000e00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000401000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			expectedHash:             "19120029436cfa37c800dce5ca006c1392b7c1d3d3dce50b91e44d1dc2e41c82",
		},
		"dkg result has only inactive provided": {
			dkgResult: &DKGResult{
				Inactive: []byte{0, 1, 1, 1},
			},
			expectedSerializedResult: "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000040001010100000000000000000000000000000000000000000000000000000000",
			expectedHash:             "941ef3e2fb3c006ea88533e1f3adf7ce974d114d30161e6a6209428a39747009",
		},
		"dkg result has all parameters provided": {
			dkgResult: &DKGResult{
				Success:        true,
				GroupPublicKey: []byte{3, 40, 200},
				Disqualified:   []byte{1, 0, 1, 0},
				Inactive:       []byte{0, 1, 1, 0},
			},
			expectedSerializedResult: "0000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000c0000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000030328c800000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004010001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000040001010000000000000000000000000000000000000000000000000000000000",
			expectedHash:             "d89ec37a4636214269a8aa8639691dc3c39157111a7321dddfcc6f6146b77fd3",
		},
		"dkg result has disqualified longer than 32 bytes": {
			dkgResult: &DKGResult{
				Disqualified: []byte{
					1, 0, 1, 0, 0, 0, 0, 1, 0, 0,
					1, 0, 1, 0, 1, 0, 1, 0, 1, 0,
					1, 0, 1, 0, 0, 0, 0, 1, 0, 0,
					1, 0, 1, 0, 1, 0, 1, 0, 1, 0,
					1, 0, 1, 0, 0, 0, 0, 1, 0, 0,
					1, 0, 1, 0, 1, 0, 1, 0, 1, 0,
				},
			},
			expectedSerializedResult: "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000003c010001000000000100000100010001000100010001000100000000010000010001000100010001000100010000000001000001000100010001000100000000000000000000000000000000000000000000000000000000000000000000000000",
			expectedHash:             "b390bfe240f6cf522b4d866b069f08a96bf9d2c6a7cac6574291a8353e05562a",
		},
		"dkg result has group public key longer than 64 bytes": {
			dkgResult: &DKGResult{
				GroupPublicKey: []byte{
					33, 249, 72, 108, 111, 44, 64, 58, 107, 112,
					108, 74, 214, 170, 149, 99, 212, 2, 48, 137,
					146, 12, 128, 8, 103, 47, 13, 161, 14, 126,
					5, 151, 0, 199, 90, 57, 31, 29, 175, 197,
					158, 45, 138, 205, 82, 95, 171, 104, 246, 8,
					203, 130, 138, 115, 72, 108, 232, 87, 129, 161,
					39, 228, 55, 222, 94, 238, 85, 128, 137, 187,
					27, 252, 25, 38, 201, 41, 127, 179, 75, 112,
				},
			},
			expectedSerializedResult: "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000120000000000000000000000000000000000000000000000000000000000000005021f9486c6f2c403a6b706c4ad6aa9563d4023089920c8008672f0da10e7e059700c75a391f1dafc59e2d8acd525fab68f608cb828a73486ce85781a127e437de5eee558089bb1bfc1926c9297fb34b700000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			expectedHash:             "86706e29471a9a32ca95525a8f50ec9dbd0751e6eff09f0bf3630023ea464545",
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			expectedEncodedResult := common.Hex2Bytes(test.expectedSerializedResult)
			expectedHash := common.Hex2Bytes(test.expectedHash)

			actualEncodedResult, err := test.dkgResult.encode()
			if err != nil {
				t.Fatal(err)
			}

			actualHash, err := test.dkgResult.Hash()
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(expectedEncodedResult, actualEncodedResult) {
				t.Errorf(
					"\nexpected: %v\nactual:   %x\n",
					test.expectedSerializedResult,
					actualEncodedResult,
				)
			}

			if !bytes.Equal(expectedHash, actualHash[:]) {
				t.Errorf(
					"\nexpected: %v\nactual:   %x\n",
					test.expectedHash,
					actualHash,
				)
			}
		})
	}
}
