package commitment

import (
	"math/big"
	"testing"
)

func TestGenerateAndValidateCommitment(t *testing.T) {
	committedValue := "eeyore"
	parameters, err := GenerateParameters()
	if err != nil {
		t.Fatalf("parameters generation error [%v]", err)
	}

	var tests = map[string]struct {
		verificationValue     string
		modifyDecommitmentKey func(key *DecommitmentKey)
		modifyCommitment      func(commitment *Commitment)
		modifyParameters      func(parameters *Parameters)
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
		"negative validation - incorrect `g`": {
			verificationValue: committedValue,
			modifyParameters: func(parameters *Parameters) {
				parameters.g = big.NewInt(23)
			},
			expectedResult: false,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {

			commitment, decommitmentKey, err := Generate(
				parameters, []byte(committedValue),
			)
			if err != nil {
				t.Fatalf("generation error [%v]", err)
			}

			if test.modifyCommitment != nil {
				test.modifyCommitment(commitment)
			}

			if test.modifyDecommitmentKey != nil {
				test.modifyDecommitmentKey(decommitmentKey)
			}

			if test.modifyParameters != nil {
				test.modifyParameters(parameters)
			}

			result := commitment.Verify(parameters, decommitmentKey, []byte(test.verificationValue))

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
