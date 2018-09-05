package zkp

import (
	"crypto/elliptic"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/keep-network/paillier"
)

func TestGeneratePublicParameters(t *testing.T) {
	paillierModulus := big.NewInt(64319 * 59063)
	curve := secp256k1.S256()

	params, err := GeneratePublicParameters(paillierModulus, curve)
	if err != nil {
		t.Fatal(err)
	}

	if params.NTilde.BitLen() != paillierModulus.BitLen() {
		t.Errorf(
			"Unexpected NTilde bit length\nExpected: %v\nActual: %v",
			paillierModulus.BitLen(),
			params.NTilde.BitLen(),
		)
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

// createTestZkpParameters creates `paillier.PrivateKey` and
// `zkp.PublicParameters` from the predefined safe prime values.
//
// Since prime generation is expensive and we want our unit tests to be blazing
// fast we need to use predefined values. This approach is only valid for unit
// tests. In case of integration tests testing whole key-gen or signature
// processes, we should use a real code for generating paillier keys and ZKP
// parameters.
func createTestZkpParameters(curve elliptic.Curve) (
	*paillier.PrivateKey,
	*PublicParameters,
	error,
) {
	p, _ := new(big.Int).SetString(
		"IYZo9UHyYhqclSgETOQ9hgfVxdW480zaODUjjTlECgHsihK7mf5Ujf8MLbM0thyNlK9Ga"+
			"drGNYKqsUqMqEKusO52NUydpTgqVwykdNjJRtNle7SPP5CfKnVHBQi5JrVty1luws"+
			"WRAl1AYMv3mPlUumWbw7oopLJ3oHIKAkXlaP07",
		62,
	)

	q, _ := new(big.Int).SetString(
		"KPLO1v0WSSFFgtDWgxggB6ug8IIhVcG2CRRCveXouTqloltXOaN4xitDebsnXjktkgSUi"+
			"HLHrzzWnLnAr5Dt5w2UKSDs2wdnJeh1MKKXdnfpApiLDfv5TwbS54UzZSOW3Ajqky"+
			"iGperkiwsdzdhZQUG0kKSK177STEUSBZ9QUIKn",
		62,
	)

	paillierKey := paillier.CreatePrivateKey(p, q)

	pTilde, _ := new(big.Int).SetString(
		"T7J6D6kn5JGzrJpskwFuabU33cBMKgkfjucmEuSvCWH4PmgcghcdSBV9u8oPWikoKokF"+
			"XVVr8NriqFizRqbRiG7R0OABnxHl1lL2CtS4pv3Raoc8Ubv6uGzbat2dx0kgGkve"+
			"LfIQOVmLrFxYOvFXKHEKal9OYlZTdIcRe5vA9XfR", 62,
	)
	qTilde, _ := new(big.Int).SetString(
		"UOzoNzDzf4dFKuTyVfEXGnbaiLXTDBtWoWLR83wl34AgEOteIhDtVdngNPbcCKP1A0H9"+
			"CBRUE081rsad6ftTyiVeDhEHIToSf3LRenAJAxic4kNrzYtyKN1yFtzMKPH6ndcT"+
			"8NPTI7gEOt789ZAf3queB54mbZjWWsm1sV2plmhR", 62,
	)

	zkpParameters, err := GeneratePublicParametersFromSafePrimes(
		paillierKey.N, pTilde, qTilde, curve,
	)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"could not generate public ZKP parameters [%v]", err,
		)
	}

	return paillierKey, zkpParameters, nil
}
