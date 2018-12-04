package event

import (
	"math/big"
	"testing"
)

func TestPublishedResultEquals(t *testing.T) {
	var tests = map[string]struct {
		result1        *PublishedResult
		result2        *PublishedResult
		expectedResult bool
	}{
		"both nil": {
			result1:        nil,
			result2:        nil,
			expectedResult: true,
		},
		"both empty": {
			result1:        &PublishedResult{},
			result2:        &PublishedResult{},
			expectedResult: true,
		},
		"nil and empty": {
			result1:        nil,
			result2:        &PublishedResult{},
			expectedResult: false,
		},
		"success - equal": {
			result1:        &PublishedResult{Success: true},
			result2:        &PublishedResult{Success: true},
			expectedResult: true,
		},
		"success - not equal": {
			result1:        &PublishedResult{Success: true},
			result2:        &PublishedResult{Success: false},
			expectedResult: false,
		},
		"group public keys - equal": {
			result1:        &PublishedResult{GroupPublicKey: big.NewInt(2)},
			result2:        &PublishedResult{GroupPublicKey: big.NewInt(2)},
			expectedResult: true,
		},
		"group public keys - nil and set": {
			result1:        &PublishedResult{GroupPublicKey: nil},
			result2:        &PublishedResult{GroupPublicKey: big.NewInt(1)},
			expectedResult: false,
		},
		"group public keys - not equal": {
			result1:        &PublishedResult{GroupPublicKey: big.NewInt(3)},
			result2:        &PublishedResult{GroupPublicKey: big.NewInt(4)},
			expectedResult: false,
		},
		"disqualified - equal": {
			result1:        &PublishedResult{Disqualified: []bool{false, false, true}},
			result2:        &PublishedResult{Disqualified: []bool{false, false, true}},
			expectedResult: true,
		},
		"disqualified - not equal": {
			result1:        &PublishedResult{Disqualified: []bool{false, false, true}},
			result2:        &PublishedResult{Disqualified: []bool{false, true, false}},
			expectedResult: false,
		},
		"inactive - equal": {
			result1:        &PublishedResult{Inactive: []bool{true, true, false}},
			result2:        &PublishedResult{Inactive: []bool{true, true, false}},
			expectedResult: true,
		},
		"inactive - not equal": {
			result1:        &PublishedResult{Inactive: []bool{true, true, false}},
			result2:        &PublishedResult{Inactive: []bool{true, true}},
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
