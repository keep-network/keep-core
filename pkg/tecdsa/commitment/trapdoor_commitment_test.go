package commitment

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
)

func TestGenerateAndValidateCommitment(t *testing.T) {
	var tests = map[string]struct {
		modifySecret          func(secret []byte)
		modifyDecommitmentKey func(key *DecommitmentKey)
		modifyCommitment      func(commitment *TrapdoorCommitment)
		expectedResult        bool
	}{
		"positive validation": {
			expectedResult: true,
		},
		"negative validation - incorrect secret value": {
			modifySecret: func(secret []byte) {
				secret[0] = 0xff
			},
			expectedResult: false,
		},
		"negative validation - incorrect decommitment key `r`": {
			modifyDecommitmentKey: func(key *DecommitmentKey) {
				key.r = big.NewInt(3)
			},
			expectedResult: false,
		},
		"negative validation - incorrect `commitment`": {
			modifyCommitment: func(commitment *TrapdoorCommitment) {
				commitment.commitment = new(bn256.G2).ScalarBaseMult(big.NewInt(3))
			},
			expectedResult: false,
		},
		"negative validation - incorrect `pubKey`": {
			modifyCommitment: func(commitment *TrapdoorCommitment) {
				commitment.pubKey = big.NewInt(3)
			},
			expectedResult: false,
		},
		"negative validation - incorrect `h`": {
			modifyCommitment: func(commitment *TrapdoorCommitment) {
				commitment.h = new(bn256.G2).ScalarBaseMult(big.NewInt(6))
			},
			expectedResult: false,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			secret := []byte("top secret message")

			commitment, decommitmentKey, err := Generate(secret)
			if err != nil {
				t.Fatalf("generation error [%v]", err)
			}

			if test.modifyCommitment != nil {
				test.modifyCommitment(commitment)
			}

			if test.modifyDecommitmentKey != nil {
				test.modifyDecommitmentKey(decommitmentKey)
			}

			if test.modifySecret != nil {
				test.modifySecret(secret)
			}

			result := commitment.Verify(decommitmentKey, secret)

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

func TestGenerateAndValidateCommitmentForMultipleValues(t *testing.T) {
	var tests = map[string]struct {
		committedValue1    string
		committedValue2    string
		verificationValue1 string
		verificationValue2 string
		expectedResult     bool
	}{
		"positive validation": {
			committedValue1:    "eeyore",
			committedValue2:    "rabbit",
			verificationValue1: "eeyore",
			verificationValue2: "rabbit",
			expectedResult:     true,
		},
		"negative validation - first verification value is not the committed one": {
			committedValue1:    "eeyore",
			committedValue2:    "rabbit",
			verificationValue1: "pooh",
			verificationValue2: "rabbit",
			expectedResult:     false,
		},
		"negative validation - second verification value is not the committed one": {
			committedValue1:    "eeyore",
			committedValue2:    "rabbit",
			verificationValue1: "eeyore",
			verificationValue2: "pooh",
			expectedResult:     false,
		},
		"negative validation - verification values in different order": {
			committedValue1:    "eeyore",
			committedValue2:    "rabbit",
			verificationValue1: "rabbit",
			verificationValue2: "eeyore",
			expectedResult:     false,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			commitment, decommitmentKey, err := Generate(
				[]byte(test.committedValue1),
				[]byte(test.committedValue2),
			)
			if err != nil {
				t.Fatal(err)
			}

			result := commitment.Verify(
				decommitmentKey,
				[]byte(test.verificationValue1),
				[]byte(test.verificationValue2),
			)

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
