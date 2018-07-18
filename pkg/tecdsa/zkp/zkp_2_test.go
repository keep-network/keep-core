package zkp

import (
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/keep-network/keep-core/pkg/tecdsa"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/keep-network/paillier"
)

func TestRoundTripZKP2(t *testing.T) {
	// GIVEN
	privateKey := paillier.CreatePrivateKey(big.NewInt(23), big.NewInt(47))

	params, err := GeneratePublicParameters(privateKey.N, secp256k1.S256())
	if err != nil {
		t.Fatal(err)
	}

	eta1, err := rand.Int(rand.Reader, params.q) // factor from ZKP PI2,1
	if err != nil {
		t.Fatal(err)
	}
	eta2, err := rand.Int(rand.Reader, params.QSix()) // factor from ZKP PI2,1
	if err != nil {
		t.Fatal(err)
	}

	rc, err := paillier.GetRandomNumberInMultiplicativeGroup(params.N, rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	g := params.CurveBasePoint()                                         // curveBasePoint
	r := tecdsa.NewCurvePoint(params.curve.ScalarBaseMult(eta1.Bytes())) // eta1CurvePoint
	u, _ := privateKey.EncryptWithR(big.NewInt(5), rc)                   // encryptedFactor1 (rho from round 1), u = E(ρ)

	// D(w) = η1*D(u) + q*η2 = η1*ρ + q*η2
	// w = E(D(w)) = E(η1*ρ + q*η2) = E(η1*D(u) + E(q*η2) = E(η1*u) + E(q*η2)
	// E(η1*u)
	encryptedUEta1 := privateKey.PublicKey.Mul(u, eta1)

	// E(q*η2)
	qEta2 := new(big.Int).Mod(new(big.Int).Mul(params.q, eta2), params.N)
	encryptedQEta2, err := privateKey.EncryptWithR(qEta2, rc)
	if err != nil {
		t.Fatal(err)
	}

	w := privateKey.PublicKey.Add(encryptedUEta1, encryptedQEta2) // encryptedMaskedFactor

	// WHEN
	zkp, err := CommitZkpPi2(g, r, w.C, u.C, eta1, eta2, rc, params, rand.Reader)

	if err != nil {
		t.Fatal(err)
	}

	var tests = map[string]struct {
		verify         func() bool
		expectedResult bool
	}{
		"positive validation": {
			verify: func() bool {
				return zkp.Verify(g, r, w.C, u.C, params)
			},
			expectedResult: true,
		},
		"negative validation - wrong g": {
			verify: func() bool {
				wrongG := tecdsa.NewCurvePoint(big.NewInt(1), big.NewInt(2))
				return zkp.Verify(wrongG, r, w.C, u.C, params)
			},
			expectedResult: false,
		},
		"negative validation - wrong r": {
			verify: func() bool {
				wrongR := tecdsa.NewCurvePoint(big.NewInt(3), big.NewInt(4))
				return zkp.Verify(g, wrongR, w.C, u.C, params)
			},
			expectedResult: false,
		},
		"negative validation - wrong w": {
			verify: func() bool {
				wrongW := new(big.Int).Add(w.C, big.NewInt(1))
				return zkp.Verify(g, r, wrongW, u.C, params)
			},
			expectedResult: false,
		},
		"negative validation - wrong u": {
			verify: func() bool {
				wrongU := new(big.Int).Add(w.C, big.NewInt(1))
				return zkp.Verify(g, r, w.C, wrongU, params)
			},
			expectedResult: false,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			expectedResult := test.expectedResult
			actualResult := test.verify()
			if actualResult != expectedResult {
				t.Fatalf(
					"Expected %v from commitment validation. Got %v",
					expectedResult,
					actualResult,
				)
			}

		})
	}
}
