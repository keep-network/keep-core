package chain

import (
	"bytes"
	"testing"
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

func TestHash(t *testing.T) {
	r1 := &DKGResult{
		Success: true,
		GroupPublicKey: [32]byte{
			101, 102, 103, 104, 105, 106, 107, 108, 109, 110,
			101, 102, 103, 104, 105, 106, 107, 108, 109, 110,
			101, 102, 103, 104, 105, 106, 107, 108, 109, 110,
			101, 102,
		},
		Disqualified: []byte{1, 0, 1, 0},
		Inactive:     []byte{0, 0, 1, 0},
	}
	actualResult := r1.Hash()
	expectedResult := []byte{
		80, 164, 35, 227, 202, 182, 90, 75, 211, 124, 161, 54, 169, 20, 65, 149,
		33, 135, 254, 232, 161, 217, 97, 65, 216, 193, 251, 217, 126, 78, 137, 114,
	}
	if !bytes.Equal(actualResult, expectedResult) {
		t.Errorf("\nexpected: %v\nactual:   %v\n", expectedResult, actualResult)
	}
}

/*
Should we add some more cases to check serialization?
E.g.
	`Disqualified` empty,
	`Inactive` empty,
*/

func TestSerialize(t *testing.T) {
	var tests = map[string]struct {
		result         *DKGResult // r1
		expectedResult []byte
	}{
		"1st Test - General Case": {
			result: &DKGResult{
				Success: true,
				GroupPublicKey: [32]byte{
					101, 102, 103, 104, 105, 106, 107, 108, 109, 110,
					101, 102, 103, 104, 105, 106, 107, 108, 109, 110,
					101, 102, 103, 104, 105, 106, 107, 108, 109, 110,
					101, 102,
				},
				Disqualified: []byte{1, 0, 1, 0},
				Inactive:     []byte{0, 0, 1, 0},
			},
			expectedResult: []byte{1, 0, 0, 0, 32, 101, 102, 103, 104, 105, 106, 107, 108, 109,
				110, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 101, 102, 103, 104, 105,
				106, 107, 108, 109, 110, 101, 102, 0, 0, 0, 4, 1, 0, 1, 0, 0, 0, 0, 4, 0, 0, 1, 0,
			},
		},
		"Empty Group Public Key ": {
			result: &DKGResult{
				Success: true,
				GroupPublicKey: [32]byte{
					0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
					0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
					0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
					0, 0,
				},
				Disqualified: []byte{1, 0, 1, 0},
				Inactive:     []byte{0, 0, 1, 0},
			},
			expectedResult: []byte{
				1, 0, 0, 0, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 1, 0, 1, 0, 0, 0, 0, 4, 0, 0, 1, 0,
			},
		},
		"Disqualified Empty": {
			result: &DKGResult{
				Success: true,
				GroupPublicKey: [32]byte{
					101, 102, 103, 104, 105, 106, 107, 108, 109, 110,
					101, 102, 103, 104, 105, 106, 107, 108, 109, 110,
					101, 102, 103, 104, 105, 106, 107, 108, 109, 110,
					101, 102,
				},
				Disqualified: []byte{},
				Inactive:     []byte{0, 0, 1, 0},
			},
			expectedResult: []byte{
				1, 0, 0, 0, 32, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 101, 102, 103,
				104, 105, 106, 107, 108, 109, 110, 101, 102, 103, 104, 105, 106, 107, 108, 109,
				110, 101, 102, 0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 1, 0,
			},
		},
		"Inactive Empty": {
			result: &DKGResult{
				Success: true,
				GroupPublicKey: [32]byte{
					101, 102, 103, 104, 105, 106, 107, 108, 109, 110,
					101, 102, 103, 104, 105, 106, 107, 108, 109, 110,
					101, 102, 103, 104, 105, 106, 107, 108, 109, 110,
					101, 102,
				},
				Disqualified: []byte{1, 0, 1, 0},
				Inactive:     []byte{},
			},
			expectedResult: []byte{
				1, 0, 0, 0, 32, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 101, 102, 103,
				104, 105, 106, 107, 108, 109, 110, 101, 102, 103, 104, 105, 106, 107, 108, 109,
				110, 101, 102, 0, 0, 0, 4, 1, 0, 1, 0, 0, 0, 0, 0,
			},
		},
		"Both Disqualified and Inactive Empty": {
			result: &DKGResult{
				Success: true,
				GroupPublicKey: [32]byte{
					101, 102, 103, 104, 105, 106, 107, 108, 109, 110,
					101, 102, 103, 104, 105, 106, 107, 108, 109, 110,
					101, 102, 103, 104, 105, 106, 107, 108, 109, 110,
					101, 102,
				},
				Disqualified: []byte{},
				Inactive:     []byte{},
			},
			expectedResult: []byte{
				1, 0, 0, 0, 32, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 101, 102, 103,
				104, 105, 106, 107, 108, 109, 110, 101, 102, 103, 104, 105, 106, 107, 108, 109,
				110, 101, 102, 0, 0, 0, 0, 0, 0, 0, 0,
			},
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			actualResult := test.result.serialize()
			if !bytes.Equal(actualResult, test.expectedResult) {
				t.Errorf("\nexpected: %v\nactual:   %v\n", test.expectedResult, actualResult)
			}
		})
	}
}
