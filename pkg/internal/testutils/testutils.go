// Package testutils contains general utilities for testing to help ensure
// consistency in output style.
package testutils

import "testing"

// AssertBytesEqual takes a testing.T and two byte slices and reports an error
// if the two bytes are not equal.
func AssertBytesEqual(t *testing.T, expectedBytes []byte, actualBytes []byte) {
	maxLen := len(expectedBytes)
	diffCount := 0
	if actualLen := len(actualBytes); actualLen > maxLen {
		diffCount = maxLen - actualLen
		maxLen = actualLen
	} else {
		diffCount = actualLen - maxLen
	}

	for i := 0; i < maxLen; i++ {
		if expectedBytes[i] != actualBytes[i] {
			diffCount++
		}
	}

	if diffCount != 0 {
		t.Errorf(
			"Byte slices differ in %v places\nexpected: [%v]\nactual:   [%v]",
			diffCount,
			expectedBytes,
			actualBytes,
		)
	}
}
