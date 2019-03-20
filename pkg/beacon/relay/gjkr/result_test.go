package gjkr

import (
	"math/big"
	"testing"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
)

func TestResultEquals(t *testing.T) {
	key1 := new(bn256.G2).ScalarBaseMult(big.NewInt(12))
	key2 := new(bn256.G2).ScalarBaseMult(big.NewInt(13))

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
		"empty and nil": {
			result1:        &Result{},
			result2:        nil,
			expectedResult: false,
		},
		"group public keys - equal": {
			result1:        &Result{GroupPublicKey: key1},
			result2:        &Result{GroupPublicKey: key1},
			expectedResult: true,
		},
		"group public keys - not equal": {
			result1:        &Result{GroupPublicKey: key1},
			result2:        &Result{GroupPublicKey: key2},
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

func TestPublicKeysEqual(t *testing.T) {
	var tests = map[string]struct {
		key1           *bn256.G2
		key2           *bn256.G2
		expectedResult bool
	}{
		"both nil": {
			key1:           nil,
			key2:           nil,
			expectedResult: true,
		},
		"nil and set": {
			key1:           nil,
			key2:           new(bn256.G2).ScalarBaseMult(big.NewInt(13)),
			expectedResult: false,
		},
		"set and nil": {
			key1:           new(bn256.G2).ScalarBaseMult(big.NewInt(13)),
			key2:           nil,
			expectedResult: false,
		},
		"equal": {
			key1:           new(bn256.G2).ScalarBaseMult(big.NewInt(13)),
			key2:           new(bn256.G2).ScalarBaseMult(big.NewInt(13)),
			expectedResult: true,
		},
		"not equal": {
			key1:           new(bn256.G2).ScalarBaseMult(big.NewInt(13)),
			key2:           new(bn256.G2).ScalarBaseMult(big.NewInt(14)),
			expectedResult: false,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			actualResult := publicKeysEqual(test.key1, test.key2)

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
		"both nil": {
			slice1:         nil,
			slice2:         nil,
			expectedResult: true,
		},
		"both empty": {
			slice1:         []MemberID{},
			slice2:         []MemberID{},
			expectedResult: true,
		},
		"both equal": {
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
		"not equal - nil and set": {
			slice1:         nil,
			slice2:         []MemberID{1, 2, 3},
			expectedResult: false,
		},
		"not equal - set and nil": {
			slice1:         []MemberID{1, 2, 3},
			slice2:         nil,
			expectedResult: false,
		},
		"not equal - empty and filled": {
			slice1:         []MemberID{},
			slice2:         []MemberID{1, 2, 3},
			expectedResult: false,
		},
		"not equal - filled and empty": {
			slice1:         []MemberID{1, 2, 3},
			slice2:         []MemberID{},
			expectedResult: false,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			actualResult := memberIDSlicesEqual(test.slice1, test.slice2)

			if test.expectedResult != actualResult {
				t.Fatalf("\nexpected: %v\nactual:   %v\n", test.expectedResult, actualResult)
			}
		})
	}
}
