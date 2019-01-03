package chain

import (
	"math/big"
	"testing"
)

func TestSerialize(t *testing.T) {
	r1 := &DKGResult{
		Success:        true,
		GroupPublicKey: big.NewInt(100),
		Disqualified:   []bool{true, false, true, false},
		Inactive:       []bool{false, false, true, false},
	}
	actualResult := r1.serialize()
	expectedResult := []byte{1, 100, 1, 0, 1, 0, 0, 0, 1, 0}
	if !bytesEqual(actualResult, expectedResult) {
		t.Errorf("\nexpected: %v\nactual:   %v\n", expectedResult, actualResult)
	}
}

func TestHash(t *testing.T) {
	r1 := &DKGResult{
		Success:        true,
		GroupPublicKey: big.NewInt(100),
		Disqualified:   []bool{true, false, true, false},
		Inactive:       []bool{false, false, true, false},
	}
	actualResult := r1.Hash()
	expectedResult := []byte{8, 164, 224, 67, 206, 144, 73, 72, 162, 22, 148, 136, 241,
		243, 2, 210, 221, 121, 31, 208, 51, 144, 48, 23, 142, 126, 12, 7, 222, 185, 107, 98}
	if !bytesEqual(actualResult, expectedResult) {
		t.Errorf("\nexpected: %v\nactual:   %v\n", expectedResult, actualResult)
	}
}

// Equal tells whether a and b contain the same elements.
// A nil argument is equivalent to an empty slice.
func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
