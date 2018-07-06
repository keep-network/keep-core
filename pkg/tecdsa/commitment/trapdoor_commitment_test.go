package commitment

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
)

func TestGenerateAndValidateCommitment(t *testing.T) {
	var tests = map[string]struct {
		modifySecret     func(secret *Secret)
		modifyCommitment func(commitment *TrapdoorCommitment)
		expectedResult   bool
	}{
		"positive validation": {
			modifySecret:     nil,
			modifyCommitment: nil,
			expectedResult:   true,
		},
		"negative validation - incorrect `secret`": {
			modifySecret: func(secret *Secret) {
				msg := []byte("top secret message2")
				secret.secret = &msg
			},
			modifyCommitment: nil,
			expectedResult:   false,
		},
		"negative validation - incorrect `r`": {
			modifySecret: func(secret *Secret) {
				secret.r = big.NewInt(3)
			},
			modifyCommitment: nil,
			expectedResult:   false,
		},
		"negative validation - incorrect `commitment`": {
			modifySecret: nil,
			modifyCommitment: func(commitment *TrapdoorCommitment) {
				commitment.commitment = new(bn256.G2).ScalarBaseMult(big.NewInt(3))
			},
			expectedResult: false,
		},
		"negative validation - incorrect `pubKey`": {
			modifySecret:     nil,
			modifyCommitment: func(commitment *TrapdoorCommitment) { commitment.pubKey = big.NewInt(3) },
			expectedResult:   false,
		},
		"negative validation - incorrect `h`": {
			modifySecret: nil,
			modifyCommitment: func(commitment *TrapdoorCommitment) {
				commitment.h = new(bn256.G2).ScalarBaseMult(big.NewInt(6))
			},
			expectedResult: false,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			msg := []byte("top secret message")

			// Generate Commitment
			commitment, err := GenerateCommitment(&msg)
			if err != nil {
				t.Fatalf("generation error [%v]", err)
			}

			// Invoke modification functions defined in test
			if test.modifyCommitment != nil {
				test.modifyCommitment(commitment)
			}

			newSecret := commitment.secret
			if test.modifySecret != nil {
				test.modifySecret(&newSecret)
			}

			// Validate Commitment
			result := commitment.ValidateCommitment(&newSecret)

			// Check result
			if result != test.expectedResult {
				t.Fatalf("\nexpected: %v\nactual:   %v", test.expectedResult, result)
			}
		})
	}
}

func TestCommitmentRandomness(t *testing.T) {
	msg := []byte("top secret message")

	// Generate Commitment 1
	commitment1, err := GenerateCommitment(&msg)
	if err != nil {
		t.Fatalf("generation error [%v]", err)
	}

	// Generate Commitment 2
	commitment2, err := GenerateCommitment(&msg)
	if err != nil {
		t.Fatalf("generation error [%v]", err)
	}

	// Check decommitments are unique
	if commitment1.secret.r.Cmp(commitment2.secret.r) == 0 {
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
