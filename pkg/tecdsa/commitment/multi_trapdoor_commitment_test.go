package commitment

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
)

func TestGenerateAndValidateCommitment(t *testing.T) {
	committedValues := []string{"eeyore", "rabbit"}

	var tests = map[string]struct {
		verificationValues    []string
		modifyDecommitmentKey func(key *DecommitmentKey)
		modifyCommitment      func(commitment *MultiTrapdoorCommitment)
		expectedResult        bool
	}{
		"positive validation": {
			verificationValues: committedValues,
			expectedResult:     true,
		},
		"negative validation - first verification value is not the committed one": {
			verificationValues: []string{"pooh", "rabbit"},
			expectedResult:     false,
		},
		"negative validation - second verification value is not the committed one": {
			verificationValues: []string{"eeyore", "pooh"},
			expectedResult:     false,
		},
		"negative validation - verification values in wrong order": {
			verificationValues: []string{"rabbit", "eeyore"},
			expectedResult:     false,
		},
		"negative validation - incorrect decommitment key `r`": {
			verificationValues: committedValues,
			modifyDecommitmentKey: func(key *DecommitmentKey) {
				key.r = big.NewInt(3)
			},
			expectedResult: false,
		},
		"negative validation - incorrect `signature`": {
			verificationValues: committedValues,
			modifyDecommitmentKey: func(key *DecommitmentKey) {
				key.signature.s = big.NewInt(4)
			},
			expectedResult: false,
		},
		"negative validation - incorrect `commitment`": {
			verificationValues: committedValues,
			modifyCommitment: func(commitment *MultiTrapdoorCommitment) {
				commitment.commitment = new(bn256.G2).ScalarBaseMult(big.NewInt(3))
			},
			expectedResult: false,
		},
		"negative validation - incorrect `verification key`": {
			verificationValues: committedValues,
			modifyCommitment: func(commitment *MultiTrapdoorCommitment) {
				commitment.verificationKey.X = big.NewInt(3)
			},
			expectedResult: false,
		},
		"negative validation - incorrect `h`": {
			verificationValues: committedValues,
			modifyCommitment: func(commitment *MultiTrapdoorCommitment) {
				commitment.h = new(bn256.G2).ScalarBaseMult(big.NewInt(6))
			},
			expectedResult: false,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			commitmentBytes := toBytes(committedValues)
			verificationBytes := toBytes(test.verificationValues)

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

func toBytes(strings []string) [][]byte {
	bytes := make([][]byte, len(strings))
	for i, s := range strings {
		bytes[i] = []byte(s)
	}
	return bytes
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

	// Check signatures are unique
	if reflect.DeepEqual(decommitmentKey1.signature, decommitmentKey2.signature) {
		t.Fatal("both signatures are equal")
	}

	// Check verification keys are unique
	if reflect.DeepEqual(commitment1.verificationKey, commitment2.verificationKey) {
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
