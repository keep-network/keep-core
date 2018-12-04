package gjkr

import (
	"math/big"
	"testing"
)

func TestResultEquals(t *testing.T) {
	var tests = map[string]struct {
		result1        *Result
		result2        *Result
		expectedResult bool
	}{
		"both nil": {
			result1:        nil,
			result2:        nil,
			expectedResult: true,
		},
		"both empty": {
			result1:        &Result{},
			result2:        &Result{},
			expectedResult: true,
		},
		"nil and empty": {
			result1:        nil,
			result2:        &Result{},
			expectedResult: false,
		},
		"success - equal": {
			result1:        &Result{Success: true},
			result2:        &Result{Success: true},
			expectedResult: true,
		},
		"success - not equal": {
			result1:        &Result{Success: true},
			result2:        &Result{Success: false},
			expectedResult: false,
		},
		"group public keys - equal": {
			result1:        &Result{GroupPublicKey: big.NewInt(2)},
			result2:        &Result{GroupPublicKey: big.NewInt(2)},
			expectedResult: true,
		},
		"group public keys - nil and set": {
			result1:        &Result{GroupPublicKey: nil},
			result2:        &Result{GroupPublicKey: big.NewInt(1)},
			expectedResult: false,
		},
		"group public keys - not equal": {
			result1:        &Result{GroupPublicKey: big.NewInt(3)},
			result2:        &Result{GroupPublicKey: big.NewInt(4)},
			expectedResult: false,
		},
		"disqualified - equal": {
			result1:        &Result{Disqualified: []MemberID{1, 2, 3}},
			result2:        &Result{Disqualified: []MemberID{1, 2, 3}},
			expectedResult: true,
		},
		"disqualified - not equal": {
			result1:        &Result{Disqualified: []MemberID{1, 2, 3}},
			result2:        &Result{Disqualified: []MemberID{1, 4, 3}},
			expectedResult: false,
		},
		"inactive - equal": {
			result1:        &Result{Inactive: []MemberID{3, 2, 1}},
			result2:        &Result{Inactive: []MemberID{3, 2, 1}},
			expectedResult: true,
		},
		"inactive - not equal": {
			result1:        &Result{Inactive: []MemberID{1, 2, 3}},
			result2:        &Result{Inactive: []MemberID{1, 2}},
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

func TestMemberIDSlicesEqual(t *testing.T) {
	var tests = map[string]struct {
		slice1         []MemberID
		slice2         []MemberID
		expectedResult bool
	}{
		"equal - nil": {
			slice1:         nil,
			slice2:         nil,
			expectedResult: true,
		},
		"equal - empty": {
			slice1:         []MemberID{},
			slice2:         []MemberID{},
			expectedResult: true,
		},
		"equal": {
			slice1:         []MemberID{1, 2, 3},
			slice2:         []MemberID{1, 2, 3},
			expectedResult: true,
		},

		"not equal - changed order": {
			slice1:         []MemberID{1, 2, 3},
			slice2:         []MemberID{1, 3, 2},
			expectedResult: false,
		},
		"not equal - different one entry": {
			slice1:         []MemberID{1, 2, 3},
			slice2:         []MemberID{1, 4, 3},
			expectedResult: false,
		},
		"not equal - different length": {
			slice1:         []MemberID{1, 2},
			slice2:         []MemberID{1, 2, 3},
			expectedResult: false,
		},
		"not equal - nil": {
			slice1:         nil,
			slice2:         []MemberID{1, 2, 3},
			expectedResult: false,
		},
		"not equal - empty": {
			slice1:         []MemberID{},
			slice2:         []MemberID{1, 2, 3},
			expectedResult: false,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			actualResult := MemberIDSlicesEqual(test.slice1, test.slice2)

			if test.expectedResult != actualResult {
				t.Fatalf("\nexpected: %v\nactual:   %v\n", test.expectedResult, actualResult)
			}
		})
	}
}
