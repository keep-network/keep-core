package chain

import (
	"fmt"
	"math/big"
	"testing"
)

func Test_serialize(t *testing.T) {
	r1 := &DKGResult{
		Success:        true,
		GroupPublicKey: big.NewInt(100),
		Disqualified:   []bool{true, false, true, false},
		Inactive:       []bool{false, false, true, false},
	}
	actualResult := r1.serialize()
	expectedResult := "013130300100010000000100"
	if fmt.Sprintf("%x", actualResult) != expectedResult {
		t.Errorf("\nexpected: %v\nactual:   %v\n", expectedResult, actualResult)
	}
}

func Test_Hash(t *testing.T) {
	r1 := &DKGResult{
		Success:        true,
		GroupPublicKey: big.NewInt(100),
		Disqualified:   []bool{true, false, true, false},
		Inactive:       []bool{false, false, true, false},
	}
	actualResult := r1.Hash()
	expectedResult := "cc4140787d3c888c67575792bc377020def9ec654c7fc19d8dbfe6c19becbe6d"
	if fmt.Sprintf("%x", actualResult) != expectedResult {
		t.Errorf("\nexpected: %v\nactual:   %v\n", expectedResult, actualResult)
	}
}
