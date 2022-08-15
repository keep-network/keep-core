// Package testutils contains general utilities for testing to help ensure
// consistency in output style.
package testutils

import (
	"errors"
	"fmt"
	"math/big"
	"testing"

	crand "crypto/rand"
)

// AssertErrorsSame checks if two errors are the same error. If not, it reports
// a test failure. Note that this function doesn't check errors deep equality
// but just checks whether both arguments point to the same underlying instance.
func AssertErrorsSame(t *testing.T, expected error, actual error) {
	if expected != actual {
		t.Errorf(
			"unexpected error\nexpected: %v\nactual:   %v\n",
			expected,
			actual,
		)
	}
}

// AssertAnyErrorInChainMatchesTarget checks if any error in the error chain matches the target.
// If not, it reports a test failure.
func AssertAnyErrorInChainMatchesTarget(t *testing.T, target error, actual error) {
	if !errors.Is(actual, target) {
		t.Errorf(
			"no error in the chain matches the target\ntarget: %v\nactual:   %v\n",
			target,
			actual,
		)
	}
}

// AssertStringsEqual checks if two strings are equal. If not, it reports a test
// failure.
func AssertStringsEqual(t *testing.T, description string, expected string, actual string) {
	if expected != actual {
		t.Errorf(
			"unexpected %s\nexpected: %s\nactual:   %s\n",
			description,
			expected,
			actual,
		)
	}
}

// AssertBigIntsEqual checks if two not-nil big integers are equal. If not, it
// reports a test failure.
func AssertBigIntsEqual(t *testing.T, description string, expected *big.Int, actual *big.Int) {
	if expected.Cmp(actual) != 0 {
		t.Errorf(
			"unexpected %s\nexpected: %v\nactual:   %v\n",
			description,
			expected,
			actual,
		)
	}
}

// AssertBigIntNonZero checks if the provided not-nil big integer is non-zero.
// If the provided big integer is zero, it reports a test failure.
func AssertBigIntNonZero(t *testing.T, description string, actual *big.Int) {
	if actual.Cmp(big.NewInt(0)) == 0 {
		t.Errorf("expected %s to be non-zero", description)
	}
}

// AssertBigIntsNotEqual checks if two not-nil big integers are not equal.
// If they are equal, reports a test failure.
func AssertBigIntsNotEqual(t *testing.T, description string, a *big.Int, b *big.Int) {
	if a.Cmp(b) == 0 {
		t.Errorf("expected %s to not be equal", description)
	}
}

// AssertBytesEqual checks if the two bytes array are equal. If not, it reports
// a test failure.
func AssertBytesEqual(t *testing.T, expectedBytes []byte, actualBytes []byte) {
	err := testBytesEqual(expectedBytes, actualBytes)

	if err != nil {
		t.Error(err)
	}
}

func testBytesEqual(expectedBytes []byte, actualBytes []byte) error {
	minLen := len(expectedBytes)
	diffCount := 0
	if actualLen := len(actualBytes); actualLen < minLen {
		diffCount = minLen - actualLen
		minLen = actualLen
	} else {
		diffCount = actualLen - minLen
	}

	for i := 0; i < minLen; i++ {
		if expectedBytes[i] != actualBytes[i] {
			diffCount++
		}
	}

	if diffCount != 0 {
		return fmt.Errorf(
			"Byte slices differ in %v places\nexpected: [%v]\nactual:   [%v]",
			diffCount,
			expectedBytes,
			actualBytes,
		)
	}

	return nil
}

// AssertIntsEqual checks if two integers are equal. If not, it reports a test
// failure.
func AssertIntsEqual(t *testing.T, description string, expected int, actual int) {
	if expected != actual {
		t.Errorf(
			"unexpected %s\nexpected: %v\nactual:   %v\n",
			description,
			expected,
			actual,
		)
	}
}

// AssertBoolsEqual checks if two booleans are equal. If not, it reports a test
// failure.
func AssertBoolsEqual(t *testing.T, description string, expected bool, actual bool) {
	if expected != actual {
		t.Errorf(
			"unexpected %s\nexpected: %v\nactual:   %v\n",
			description,
			expected,
			actual,
		)
	}
}

// NewRandInt generates a random value in range [0, max), different from the
// passed current value.
func NewRandInt(currentValue, max *big.Int) *big.Int {
	newValue := currentValue
	for currentValue.Cmp(newValue) == 0 {
		newValue, _ = crand.Int(crand.Reader, max)
	}
	return newValue
}
