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

// MockRandReader is an implementation of `io.Reader` allowing to get
// predictable random numbers in your tests. Each new generated number is larger
// by 1 from the previous one starting from counter seed provided when
// constructing MockRandReader.
//
// We use `MockRandReader` to test commitment phase of ZKPs defined in this
// package where we need predictable values instead of random ones.
//
// mockRandom := &MockRandReader{ counter: big.NewInt(1) }
// r1, _ := rand.Int(mockRandom, big.NewInt(10000)) // r1=1
// r2, _ := rand.Int(mockRandom, big.NewInt(10000)) // r2=2
// r3, _ := rand.Int(mockRandom, big.NewInt(10000)) // r3=3
type MockRandReader struct {
	counter *big.Int
}

// NewMockRandReader returns new MockRandReader instance.
func NewMockRandReader(counter *big.Int) *MockRandReader {
	return &MockRandReader{counter}
}
func (r *MockRandReader) Read(b []byte) (int, error) {
	cb := r.counter.Bytes()
	for i := range b {
		// iterate backwards
		bIdx := len(b) - i - 1
		cbIdx := len(cb) - i - 1
		if cbIdx >= 0 {
			b[bIdx] = cb[cbIdx]
		}
	}
	r.counter = new(big.Int).Add(r.counter, big.NewInt(1))
	return len(b), nil
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
