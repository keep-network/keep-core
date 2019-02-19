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
		Success:        true,
		GroupPublicKey: []byte{100},
		Disqualified:   []byte{1, 0, 1, 0},
		Inactive:       []byte{0, 0, 1, 0},
	}
	actualResult := r1.Hash()
	expectedResult := []byte{253, 157, 252, 65, 254, 175, 165, 131, 87, 84, 98, 128, 37, 49,
		106, 147, 230, 64, 192, 22, 193, 71, 150, 55, 101, 1, 110, 4, 10, 189, 121, 38}
	if !bytes.Equal(actualResult, expectedResult) {
		t.Errorf("\nexpected: %v\nactual:   %v\n", expectedResult, actualResult)
	}
}

func TestSerialize(t *testing.T) {
	r1 := &DKGResult{
		Success:        true,
		GroupPublicKey: []byte{100},
		Disqualified:   []byte{1, 0, 1, 0},
		Inactive:       []byte{0, 0, 1, 0},
	}
	actualResult := r1.serialize()
	expectedResult := []byte{1, 0, 0, 0, 1, 100, 0, 0, 0, 4, 1, 0, 1, 0, 0, 0, 0, 4, 0, 0, 1, 0}
	if !bytes.Equal(actualResult, expectedResult) {
		t.Errorf("\nexpected: %v\nactual:   %v\n", expectedResult, actualResult)
	}
}
