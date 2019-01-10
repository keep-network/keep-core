package chain

import (
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
		"group public keys - nil and set": {
			result1:        &DKGResult{GroupPublicKey: nil},
			result2:        &DKGResult{GroupPublicKey: []byte{1}},
			expectedResult: false,
		},
		"group public keys - set and nil": {
			result1:        &DKGResult{GroupPublicKey: []byte{1}},
			result2:        &DKGResult{GroupPublicKey: nil},
			expectedResult: false,
		},
		"group public keys - not equal": {
			result1:        &DKGResult{GroupPublicKey: []byte{3}},
			result2:        &DKGResult{GroupPublicKey: []byte{4}},
			expectedResult: false,
		},
		"disqualified - equal": {
			result1:        &DKGResult{Disqualified: []bool{false, false, true}},
			result2:        &DKGResult{Disqualified: []bool{false, false, true}},
			expectedResult: true,
		},
		"disqualified - not equal": {
			result1:        &DKGResult{Disqualified: []bool{false, false, true}},
			result2:        &DKGResult{Disqualified: []bool{false, true, false}},
			expectedResult: false,
		},
		"inactive - equal": {
			result1:        &DKGResult{Inactive: []bool{true, true, false}},
			result2:        &DKGResult{Inactive: []bool{true, true, false}},
			expectedResult: true,
		},
		"inactive - not equal": {
			result1:        &DKGResult{Inactive: []bool{true, true, false}},
			result2:        &DKGResult{Inactive: []bool{true, true}},
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
