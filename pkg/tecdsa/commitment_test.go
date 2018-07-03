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
	r, pubKey, h, commitment, err := GenerateCommitment(&secret)

	// Validate Commitment
	result, err := ValidateCommitment(&secret, r, commitment, pubKey, h)

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
	r1, pubKey1, h1, commitment1, err := GenerateCommitment(&secret)
	if err != nil {
		t.Fatalf("validation error [%v]", err)
	}

	// Generate Commitment 2
	r2, pubKey2, h2, commitment2, err := GenerateCommitment(&secret)
	if err != nil {
		t.Fatalf("validation error [%v]", err)
	}

	// Check decommitments are unique
	if r1.Cmp(r2) == 0 {
		t.Fatal("both decommitment keys `r` are equal")
	}

	// Check public keys are unique
	if pubKey1.String() == pubKey2.String() {
		t.Fatal("both public keys `pubKey` are equal")
	}

	// Check random points are unique
	if h1.String() == h2.String() {
		t.Fatal("both random points `h` are equal")
	}

	// Check commitments are unique
	if commitment1.String() == commitment2.String() {
		t.Fatal("both commitments are equal")
	}
}

func TestValidate(t *testing.T) {
	var tests = map[string]struct {
		secret         []byte
		r              *big.Int
		commitment     *bn256.G2
		pubKey         *bn256.G2
		h              *bn256.G2
		expectedResult bool
		expectedError  error
	}{
		"positive validation - pass values used for generation": {
			secret:         []byte("top secret message"),
			r:              big.NewInt(2),
			commitment:     convertBase64StringToCommitment("DZv5bcjEHmVfO+ukZUq7uLUKQhXTbE8mBkG4CTpt6bAiKnOXf2v6E3reMnW9qV0B5hjT9NPxaux6OQW5/vnw3iHJ2Gu8/icq3T7sKvoD+CXFDi8uYZv7pln7zyAI32OiGhjPPvBtZ/jtaIx0WTF58f1HVrQj1vWeQO43zFJCi38="),
			pubKey:         new(bn256.G2).ScalarBaseMult(big.NewInt(4)),
			h:              new(bn256.G2).ScalarBaseMult(big.NewInt(5)),
			expectedResult: true,
			expectedError:  nil,
		},
		"negative validation - incorrect `secret`": {
			secret:         []byte("top secret message2"),
			r:              big.NewInt(2),
			commitment:     convertBase64StringToCommitment("DZv5bcjEHmVfO+ukZUq7uLUKQhXTbE8mBkG4CTpt6bAiKnOXf2v6E3reMnW9qV0B5hjT9NPxaux6OQW5/vnw3iHJ2Gu8/icq3T7sKvoD+CXFDi8uYZv7pln7zyAI32OiGhjPPvBtZ/jtaIx0WTF58f1HVrQj1vWeQO43zFJCi38="),
			pubKey:         new(bn256.G2).ScalarBaseMult(big.NewInt(4)),
			h:              new(bn256.G2).ScalarBaseMult(big.NewInt(5)),
			expectedResult: false,
			expectedError:  fmt.Errorf("pairings doesn't match"),
		},
		"negative validation - incorrect `r`": {
			secret:         []byte("top secret message2"),
			r:              big.NewInt(3),
			commitment:     convertBase64StringToCommitment("DZv5bcjEHmVfO+ukZUq7uLUKQhXTbE8mBkG4CTpt6bAiKnOXf2v6E3reMnW9qV0B5hjT9NPxaux6OQW5/vnw3iHJ2Gu8/icq3T7sKvoD+CXFDi8uYZv7pln7zyAI32OiGhjPPvBtZ/jtaIx0WTF58f1HVrQj1vWeQO43zFJCi38="),
			pubKey:         new(bn256.G2).ScalarBaseMult(big.NewInt(4)),
			h:              new(bn256.G2).ScalarBaseMult(big.NewInt(5)),
			expectedResult: false,
			expectedError:  fmt.Errorf("pairings doesn't match"),
		},
		"negative validation - incorrect `commitmentBase64`": {
			secret:         []byte("top secret message"),
			r:              big.NewInt(2),
			commitment:     new(bn256.G2).ScalarBaseMult(big.NewInt(3)),
			pubKey:         new(bn256.G2).ScalarBaseMult(big.NewInt(4)),
			h:              new(bn256.G2).ScalarBaseMult(big.NewInt(5)),
			expectedResult: false,
			expectedError:  fmt.Errorf("pairings doesn't match"),
		},
		"negative validation - incorrect `pubKey`": {
			secret:         []byte("top secret message"),
			r:              big.NewInt(2),
			commitment:     convertBase64StringToCommitment("DZv5bcjEHmVfO+ukZUq7uLUKQhXTbE8mBkG4CTpt6bAiKnOXf2v6E3reMnW9qV0B5hjT9NPxaux6OQW5/vnw3iHJ2Gu8/icq3T7sKvoD+CXFDi8uYZv7pln7zyAI32OiGhjPPvBtZ/jtaIx0WTF58f1HVrQj1vWeQO43zFJCi38="),
			pubKey:         new(bn256.G2).ScalarBaseMult(big.NewInt(3)),
			h:              new(bn256.G2).ScalarBaseMult(big.NewInt(5)),
			expectedResult: false,
			expectedError:  fmt.Errorf("pairings doesn't match"),
		},
		"negative validation - incorrect `h`": {
			secret:         []byte("top secret message"),
			r:              big.NewInt(2),
			commitment:     convertBase64StringToCommitment("DZv5bcjEHmVfO+ukZUq7uLUKQhXTbE8mBkG4CTpt6bAiKnOXf2v6E3reMnW9qV0B5hjT9NPxaux6OQW5/vnw3iHJ2Gu8/icq3T7sKvoD+CXFDi8uYZv7pln7zyAI32OiGhjPPvBtZ/jtaIx0WTF58f1HVrQj1vWeQO43zFJCi38="),
			pubKey:         new(bn256.G2).ScalarBaseMult(big.NewInt(4)),
			h:              new(bn256.G2).ScalarBaseMult(big.NewInt(6)),
			expectedResult: false,
			expectedError:  fmt.Errorf("pairings doesn't match"),
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			result, err := ValidateCommitment(&test.secret, test.r, test.commitment, test.pubKey, test.h)

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
