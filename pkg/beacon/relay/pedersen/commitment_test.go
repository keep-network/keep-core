package pedersen

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"
)

func TestGenerateAndValidateCommitment(t *testing.T) {
	committedValue := "eeyore"
	vss, err := NewVSS()
	if err != nil {
		t.Fatalf("vss creation failed [%v]", err)
	}

	var tests = map[string]struct {
		verificationValue     string
		modifyDecommitmentKey func(key *DecommitmentKey)
		modifyCommitment      func(commitment *Commitment)
		modifyVSS             func(parameters *VSS)
		expectedResult        bool
	}{
		"positive validation": {
			verificationValue: committedValue,
			expectedResult:    true,
		},
		"negative validation - verification value is not the committed one": {
			verificationValue: "pooh",
			expectedResult:    false,
		},
		"negative validation - incorrect decommitment key `r`": {
			verificationValue: committedValue,
			modifyDecommitmentKey: func(key *DecommitmentKey) {
				key.r = big.NewInt(3)
			},
			expectedResult: false,
		},
		"negative validation - incorrect `commitment`": {
			verificationValue: committedValue,
			modifyCommitment: func(commitment *Commitment) {
				commitment.commitment = big.NewInt(13)
			},
			expectedResult: false,
		},
		"negative validation - incorrect `h`": {
			verificationValue: committedValue,
			modifyVSS: func(vss *VSS) {
				vss.h = big.NewInt(23)
			},
			expectedResult: false,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {

			commitment, decommitmentKey, err := vss.CommitmentTo([]byte(committedValue))
			if err != nil {
				t.Fatalf("generation error [%v]", err)
			}

			if test.modifyCommitment != nil {
				test.modifyCommitment(commitment)
			}

			if test.modifyDecommitmentKey != nil {
				test.modifyDecommitmentKey(decommitmentKey)
			}

			if test.modifyVSS != nil {
				test.modifyVSS(vss)
			}

			result := commitment.Verify(decommitmentKey, []byte(test.verificationValue))

			if result != test.expectedResult {
				t.Fatalf(
					"expected: %v\nactual: %v\n",
					test.expectedResult,
					result,
				)
			}
		})
	}
}

func TestNewVSSpqValidation(t *testing.T) {
	p := big.NewInt(17)
	q := big.NewInt(4)
	expectedError := fmt.Errorf("incorrect p and q values")

	_, err := NewVSS(p, q)

	if !reflect.DeepEqual(err, expectedError) {
		t.Fatalf("actual error: %v\nexpected error: %v", err, expectedError)
	}
}
