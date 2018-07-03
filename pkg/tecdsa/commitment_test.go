package tecdsa

import (
	"encoding/base64"
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/crypto/bn256"
)

func TestGenerateAndValidateCommitment(t *testing.T) {
	secret := []byte("top secret message")

	// Generate Commitment
	params, err := GenerateCommitment(&secret)

	// Validate Commitment
	result, err := ValidateCommitment(params)

	if err != nil {
		t.Fatalf("validation error [%v]", err)
	}

	if !result {
		t.Fatalf("validation result: [%v]", result)
	}
}

func TestGenerateTwoCommitmentsCheckUniqueResults(t *testing.T) {
	secret := []byte("top secret message")

	// Generate Commitment 1
	params1, err := GenerateCommitment(&secret)
	if err != nil {
		t.Fatalf("validation error [%v]", err)
	}

	// Generate Commitment 2
	params2, err := GenerateCommitment(&secret)
	if err != nil {
		t.Fatalf("validation error [%v]", err)
	}

	// Check decommitments are unique
	if params1.r.Cmp(params2.r) == 0 {
		t.Fatal("both decommitment keys `r` are equal")
	}

	// Check public keys are unique
	if params1.pubKey.String() == params2.pubKey.String() {
		t.Fatal("both public keys `pubKey` are equal")
	}

	// Check random points are unique
	if params1.h.String() == params2.h.String() {
		t.Fatal("both random points `h` are equal")
	}

	// Check commitments are unique
	if params1.commitment.String() == params2.commitment.String() {
		t.Fatal("both commitments are equal")
	}
}

func TestValidate(t *testing.T) {
	secret := []byte("top secret message")
	h := new(bn256.G2).ScalarBaseMult(big.NewInt(5))
	r := big.NewInt(2)
	pubKey := new(bn256.G2).ScalarBaseMult(big.NewInt(4))
	commitment := convertBase64StringToCommitment("DZv5bcjEHmVfO+ukZUq7uLUKQhXTbE8mBkG4CTpt6bAiKnOXf2v6E3reMnW9qV0B5hjT9NPxaux6OQW5/vnw3iHJ2Gu8/icq3T7sKvoD+CXFDi8uYZv7pln7zyAI32OiGhjPPvBtZ/jtaIx0WTF58f1HVrQj1vWeQO43zFJCi38=")

	secretInvalid := []byte("top secret message2")

	var tests = map[string]struct {
		params         *commitmentParams
		expectedResult bool
		expectedError  error
	}{
		"positive validation - pass values used for generation": {
			params: &commitmentParams{
				publicParams: publicParams{
					commitment: commitment,
					pubKey:     pubKey,
					h:          h,
				},
				secret: &secret,
				r:      r,
			},
			expectedResult: true,
			expectedError:  nil,
		},
		"negative validation - incorrect `secret`": {
			params: &commitmentParams{
				publicParams: publicParams{
					commitment: commitment,
					pubKey:     pubKey,
					h:          new(bn256.G2).ScalarBaseMult(big.NewInt(5)),
				},
				secret: &secretInvalid,
				r:      r,
			},
			expectedResult: false,
			expectedError:  fmt.Errorf("pairings doesn't match"),
		},
		"negative validation - incorrect `r`": {
			params: &commitmentParams{
				publicParams: publicParams{
					commitment: commitment,
					pubKey:     pubKey,
					h:          new(bn256.G2).ScalarBaseMult(big.NewInt(5)),
				},
				secret: &secret,
				r:      big.NewInt(3),
			},
			expectedResult: false,
			expectedError:  fmt.Errorf("pairings doesn't match"),
		},
		"negative validation - incorrect `commitment`": {
			params: &commitmentParams{
				publicParams: publicParams{
					commitment: new(bn256.G2).ScalarBaseMult(big.NewInt(3)),
					pubKey:     pubKey,
					h:          h,
				},
				secret: &secret,
				r:      r,
			},
			expectedResult: false,
			expectedError:  fmt.Errorf("pairings doesn't match"),
		},
		"negative validation - incorrect `pubKey`": {
			params: &commitmentParams{
				publicParams: publicParams{
					commitment: commitment,
					pubKey:     new(bn256.G2).ScalarBaseMult(big.NewInt(3)),
					h:          h,
				},
				secret: &secret,
				r:      r,
			},
			expectedResult: false,
			expectedError:  fmt.Errorf("pairings doesn't match"),
		},
		"negative validation - incorrect `h`": {
			params: &commitmentParams{
				publicParams: publicParams{
					commitment: commitment,
					pubKey:     pubKey,
					h:          new(bn256.G2).ScalarBaseMult(big.NewInt(6)),
				},
				secret: &secret,
				r:      r,
			},
			expectedResult: false,
			expectedError:  fmt.Errorf("pairings doesn't match"),
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			result, err := ValidateCommitment(test.params)

			if result != test.expectedResult {
				t.Fatalf("\nexpected: %v\nactual:   %v", test.expectedResult, result)
			}

			if !reflect.DeepEqual(err, test.expectedError) {
				t.Fatalf("\nexpected: %v\nactual:   %v", test.expectedError, err)
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
