package tecdsa

import (
	"encoding/base64"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/crypto/bn256"
)

func TestGenerateAndValidateCommitment(t *testing.T) {
	msg := []byte("top secret message")

	// Generate Commitment
	commitment, secret, err := GenerateCommitment(&msg)
	if err != nil {
		t.Fatalf("generation error [%v]", err)
	}

	// Validate Commitment
	result := ValidateCommitment(commitment, secret)

	if !result {
		t.Fatalf("validation result: [%v]", result)
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

func TestValidate(t *testing.T) {
	msg := []byte("top secret message")
	h := new(bn256.G2).ScalarBaseMult(big.NewInt(5))
	r := big.NewInt(2)
	pubKey := big.NewInt(4)
	commitment := convertBase64StringToCommitment("DZv5bcjEHmVfO+ukZUq7uLUKQhXTbE8mBkG4CTpt6bAiKnOXf2v6E3reMnW9qV0B5hjT9NPxaux6OQW5/vnw3iHJ2Gu8/icq3T7sKvoD+CXFDi8uYZv7pln7zyAI32OiGhjPPvBtZ/jtaIx0WTF58f1HVrQj1vWeQO43zFJCi38=")

	msgInvalid := []byte("top secret message2")

	var tests = map[string]struct {
		secret         *Secret
		commitment     *Commitment
		expectedResult bool
		expectedError  error
	}{
		"positive validation - pass values used for generation": {
			secret: &Secret{
				secret: &msg,
				r:      r,
			},
			commitment: &Commitment{
				commitment: commitment,
				pubKey:     pubKey,
				h:          h,
			},
			expectedResult: true,
		},
		"negative validation - incorrect `secret`": {
			secret: &Secret{
				secret: &msgInvalid,
				r:      r,
			},
			commitment: &Commitment{
				commitment: commitment,
				pubKey:     pubKey,
				h:          new(bn256.G2).ScalarBaseMult(big.NewInt(5)),
			},
			expectedResult: false,
		},
		"negative validation - incorrect `r`": {
			secret: &Secret{

				secret: &msg,
				r:      big.NewInt(3),
			},
			commitment: &Commitment{
				commitment: commitment,
				pubKey:     pubKey,
				h:          new(bn256.G2).ScalarBaseMult(big.NewInt(5)),
			},
			expectedResult: false,
		},
		"negative validation - incorrect `commitment`": {
			secret: &Secret{

				secret: &msg,
				r:      r,
			},
			commitment: &Commitment{
				commitment: new(bn256.G2).ScalarBaseMult(big.NewInt(3)),
				pubKey:     pubKey,
				h:          h,
			},
			expectedResult: false,
		},
		"negative validation - incorrect `pubKey`": {
			secret: &Secret{

				secret: &msg,
				r:      r,
			},
			commitment: &Commitment{
				commitment: commitment,
				pubKey:     big.NewInt(3),
				h:          h,
			},
			expectedResult: false,
		},
		"negative validation - incorrect `h`": {
			secret: &Secret{

				secret: &msg,
				r:      r,
			},
			commitment: &Commitment{
				commitment: commitment,
				pubKey:     pubKey,
				h:          new(bn256.G2).ScalarBaseMult(big.NewInt(6)),
			},
			expectedResult: false,
			expectedError:  fmt.Errorf("pairings doesn't match"),
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			result := ValidateCommitment(test.commitment, test.secret)

			if result != test.expectedResult {
				t.Fatalf("\nexpected: %v\nactual:   %v", test.expectedResult, result)
			}
		})
	}
}

func convertBase64StringToCommitment(str string) *bn256.G2 {
	commitmentBytes, _ := base64.StdEncoding.DecodeString(str)

	commitment := new(bn256.G2)
	commitment.Unmarshal(commitmentBytes)

	return commitment
}
