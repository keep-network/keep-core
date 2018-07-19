package commitment

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
)

func TestGenerateAndValidateCommitment(t *testing.T) {
	var tests = map[string]struct {
		committedValues       []string
		verificationValues    []string
		modifyDecommitmentKey func(key *DecommitmentKey)
		modifyCommitment      func(commitment *TrapdoorCommitment)
		expectedResult        bool
	}{
		"positive validation": {
			committedValues:    []string{"eeyore", "rabbit"},
			verificationValues: []string{"eeyore", "rabbit"},
			expectedResult:     true,
		},
		"negative validation - first verification value is not the committed one": {
			committedValues:    []string{"eeyore", "rabbit"},
			verificationValues: []string{"pooh", "rabbit"},
			expectedResult:     false,
		},
		"negative validation - second verification value is not the committed one": {
			committedValues:    []string{"eeyore", "rabbit"},
			verificationValues: []string{"eeyore", "pooh"},
			expectedResult:     false,
		},
		"negative validation - verification values in wrong order": {
			committedValues:    []string{"eeyore", "rabbit"},
			verificationValues: []string{"rabbit", "eeyore"},
			expectedResult:     false,
		},
		"negative validation - incorrect decommitment key `r`": {
			committedValues:    []string{"eeyore", "rabbit"},
			verificationValues: []string{"eeyore", "rabbit"},
			modifyDecommitmentKey: func(key *DecommitmentKey) {
				key.r = big.NewInt(3)
			},
			expectedResult: false,
		},
		"negative validation - incorrect `commitment`": {
			committedValues:    []string{"eeyore", "rabbit"},
			verificationValues: []string{"eeyore", "rabbit"},
			modifyCommitment: func(commitment *TrapdoorCommitment) {
				commitment.commitment = new(bn256.G2).ScalarBaseMult(big.NewInt(3))
			},
			expectedResult: false,
		},
		"negative validation - incorrect `pubKey`": {
			committedValues:    []string{"eeyore", "rabbit"},
			verificationValues: []string{"eeyore", "rabbit"},
			modifyCommitment: func(commitment *TrapdoorCommitment) {
				commitment.pubKey = big.NewInt(3)
			},
			expectedResult: false,
		},
		"negative validation - incorrect `h`": {
			committedValues:    []string{"eeyore", "rabbit"},
			verificationValues: []string{"eeyore", "rabbit"},
			modifyCommitment: func(commitment *TrapdoorCommitment) {
				commitment.h = new(bn256.G2).ScalarBaseMult(big.NewInt(6))
			},
			expectedResult: false,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			commitmentBytes := make([][]byte, len(test.committedValues))
			for i, committedValue := range test.committedValues {
				commitmentBytes[i] = []byte(committedValue)
			}

			verificationBytes := make([][]byte, len(test.verificationValues))
			for i, verificationValue := range test.verificationValues {
				verificationBytes[i] = []byte(verificationValue)
			}

			commitment, decommitmentKey, err := Generate(commitmentBytes...)
			if err != nil {
				t.Fatalf("generation error [%v]", err)
			}

			if test.modifyCommitment != nil {
				test.modifyCommitment(commitment)
			}

			if test.modifyDecommitmentKey != nil {
				test.modifyDecommitmentKey(decommitmentKey)
			}

			result := commitment.Verify(decommitmentKey, verificationBytes...)

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

func TestCommitmentRandomness(t *testing.T) {
	secret := []byte("top secret message")

	// Generate Commitment 1
	commitment1, decommitmentKey1, err := Generate(secret)
	if err != nil {
		t.Fatalf("generation error [%v]", err)
	}

	// Generate Commitment 2
	commitment2, decommitmentKey2, err := Generate(secret)
	if err != nil {
		t.Fatalf("generation error [%v]", err)
	}

	// Check decommitment keys are unique
	if decommitmentKey1.r.Cmp(decommitmentKey2.r) == 0 {
		t.Fatal("both decommitment keys `r` are equal")
	}

	// Check public keys are unique
	if commitment1.pubKey.String() == commitment2.pubKey.String() {
		t.Fatal("both public keys `pubKey` are equal")
	}

	// Check master public keys are unique
	if commitment1.h.String() == commitment2.h.String() {
		t.Fatal("both master public keys `h` are equal")
	}

	// Check commitments are unique
	if commitment1.commitment.String() == commitment2.commitment.String() {
		t.Fatal("both commitments are equal")
	}
}
