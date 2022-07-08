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
		"misbehaved - equal": {
			result1:        &DKGResult{Misbehaved: []byte{0x01, 0x02, 0x03}},
			result2:        &DKGResult{Misbehaved: []byte{0x01, 0x02, 0x03}},
			expectedResult: true,
		},
		"misbehaved - other members, same length": {
			result1:        &DKGResult{Misbehaved: []byte{0x01, 0x02, 0x04}},
			result2:        &DKGResult{Misbehaved: []byte{0x01, 0x02, 0x05}},
			expectedResult: false,
		},
		"misbehaved - different length": {
			result1:        &DKGResult{Misbehaved: []byte{0x01, 0x02, 0x03}},
			result2:        &DKGResult{Misbehaved: []byte{0x01, 0x02}},
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
