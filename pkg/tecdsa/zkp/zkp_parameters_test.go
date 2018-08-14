package zkp

import (
	"errors"
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

func TestGeneratePublicParameters(t *testing.T) {
	paillierModulus := big.NewInt(13)
	curve := secp256k1.S256()

	params, err := GeneratePublicParameters(paillierModulus, curve)
	if err != nil {
		t.Fatal(err)
	}

	if params.NTilde.Cmp(big.NewInt(0)) <= 0 {
		t.Errorf("NTilde should be greater than 0. Actual: %v", params.NTilde)
	}
	if params.h1.Cmp(big.NewInt(0)) <= 0 {
		t.Errorf("h1 should be greater than 0. Actual: %v", params.h1)
	}
	if params.h2.Cmp(big.NewInt(0)) <= 0 {
		t.Errorf("h2 should be greater than 0. Actual: %v", params.h2)
	}
	if params.q.Cmp(secp256k1.S256().Params().N) != 0 {
		t.Errorf(
			"Unexpected q\nExpected: %v\nActual: %v",
			secp256k1.S256().Params().N,
			params.q,
		)
	}
	if params.N != paillierModulus {
		t.Errorf("Unexpected N\nExpected: %v\nActual: %v",
			paillierModulus,
			params.N,
		)
	}
	if params.curve != curve {
		t.Errorf("Unexpected Elliptic Curve used")
	}
}

func TestGeneratePublicParametersFromSafePrimes(t *testing.T) {
	paillierModulus := big.NewInt(13)
	curve := secp256k1.S256()

	var tests = map[string]struct {
		pTilde        *big.Int
		qTilde        *big.Int
		expectedError error
	}{
		"ZKP parameters successfully created": {
			pTilde: big.NewInt(64319),
			qTilde: big.NewInt(59063),
		},
		"can't create ZKP parameters for equal safe primes": {
			pTilde:        big.NewInt(64319),
			qTilde:        big.NewInt(64319),
			expectedError: errors.New("safe primes must not be equal, got 64319 and 64319"),
		},
		"can't create ZKP parameters for safe primes not having the same bit length": {
			pTilde:        big.NewInt(59063),
			qTilde:        big.NewInt(1283),
			expectedError: errors.New("safe primes must have the same bit length, got 16 and 11 bits"),
		},
		"can't create ZKP if pTilde is not prime": {
			pTilde:        big.NewInt(64318),
			qTilde:        big.NewInt(59063),
			expectedError: errors.New("both numbers must be safe primes"),
		},
		"can't crate ZKP if qTilde is not prime": {
			pTilde:        big.NewInt(64319),
			qTilde:        big.NewInt(59064),
			expectedError: errors.New("both numbers must be safe primes"),
		},
		"can't create ZKP if pTilde is not a safe prime": {
			pTilde:        big.NewInt(109),
			qTilde:        big.NewInt(107),
			expectedError: errors.New("both numbers must be safe primes"),
		},
		"can't create ZKP if qTilde is not a safe prime": {
			pTilde:        big.NewInt(107),
			qTilde:        big.NewInt(109),
			expectedError: errors.New("both numbers must be safe primes"),
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			params, err := GeneratePublicParametersFromSafePrimes(
				paillierModulus, test.pTilde, test.qTilde, curve,
			)

			if !reflect.DeepEqual(err, test.expectedError) {
				t.Fatalf(
					"Unexpected error\nExpected: %v\nActual: %v",
					test.expectedError,
					err,
				)
			}

			if test.expectedError == nil {
				if params == nil {
					t.Fatal("Generated parameters can't be nil")
				}

				expectedNTilde := new(big.Int).Mul(test.pTilde, test.qTilde)
				if params.NTilde.Cmp(expectedNTilde) != 0 {
					t.Fatalf(
						"Unexpected NTilde\nExpected: %v\nActual: %v",
						expectedNTilde,
						params.NTilde,
					)
				}
			}

		})
	}
}
