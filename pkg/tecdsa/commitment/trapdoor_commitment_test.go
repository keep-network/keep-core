package commitment

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/crypto/bn256"
)

func TestGenerateAndValidateCommitment(t *testing.T) {
	var tests = map[string]struct {
		modifySecret     func(secret *TrapdoorSecret)
		modifyCommitment func(commitment *TrapdoorCommitment)
		expectedResult   bool
	}{
		"positive validation - pass values used for generation": {
			modifySecret:     func(secret *TrapdoorSecret) {},
			modifyCommitment: func(commitment *TrapdoorCommitment) {},
			expectedResult:   true,
		},
		"negative validation - incorrect `secret`": {
			modifySecret: func(secret *TrapdoorSecret) {
				msg := []byte("top secret message2")
				secret.secret = &msg
			},
			modifyCommitment: func(commitment *TrapdoorCommitment) {},
			expectedResult:   false,
		},
		"negative validation - incorrect `r`": {
			modifySecret: func(secret *TrapdoorSecret) {
				secret.r = big.NewInt(3)
			},
			modifyCommitment: func(commitment *TrapdoorCommitment) {},
			expectedResult:   false,
		},
		"negative validation - incorrect `commitment`": {
			modifySecret: func(secret *TrapdoorSecret) {},
			modifyCommitment: func(commitment *TrapdoorCommitment) {
				commitment.commitment = new(bn256.G2).ScalarBaseMult(big.NewInt(3))
			},
			expectedResult: false,
		},
		"negative validation - incorrect `pubKey`": {
			modifySecret:     func(secret *TrapdoorSecret) {},
			modifyCommitment: func(commitment *TrapdoorCommitment) { commitment.pubKey = big.NewInt(3) },
			expectedResult:   false,
		},
		"negative validation - incorrect `h`": {
			modifySecret: func(secret *TrapdoorSecret) {},
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
			commitment, secret, err := GenerateCommitment(&msg)
			if err != nil {
				t.Fatalf("generation error [%v]", err)
			}

			// Invoke modification functions defined in test
			test.modifyCommitment(commitment)
			test.modifySecret(secret)

			// Validate Commitment
			result := ValidateCommitment(commitment, secret)

			// Check result
			if result != test.expectedResult {
				t.Fatalf("\nexpected: %v\nactual:   %v", test.expectedResult, result)
			}
		})
	}
}

func TestGenerateTwoCommitmentsCheckUniqueResults(t *testing.T) {
	msg := []byte("top secret message")

	// Generate Commitment 1
	commitment1, secret1, err := GenerateCommitment(&msg)
	if err != nil {
		t.Fatalf("generation error [%v]", err)
	}

	// Generate Commitment 2
	commitment2, secret2, err := GenerateCommitment(&msg)
	if err != nil {
		t.Fatalf("generation error [%v]", err)
	}

	// Check decommitments are unique
	if secret1.r.Cmp(secret2.r) == 0 {
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
