// Package testutils contains general utilities for testing to help ensure
// consistency in output style.
package testutils

import (
	"fmt"
	"math/big"
	"testing"

	crand "crypto/rand"
)

// AssertBytesEqual takes a testing.T and two byte slices and reports an error
// if the two bytes are not equal.
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

// NewRandInt generates a random value in range [0, max), different from the
// passed current value.
func NewRandInt(currentValue, max *big.Int) *big.Int {
	newValue := currentValue
	for currentValue.Cmp(newValue) == 0 {
		newValue, _ = crand.Int(crand.Reader, max)
	}
	return newValue
}
