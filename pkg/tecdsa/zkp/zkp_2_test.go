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
	message := big.NewInt(5)

	p, _ := new(big.Int).SetString("23", 10)
	q, _ := new(big.Int).SetString("47", 10)

	privateKey := paillier.CreatePrivateKey(p, q)

	params, err := GeneratePublicParameters(privateKey.N, secp256k1.S256())
	if err != nil {
		t.Fatal(err)
	}

	eta1, err := rand.Int(rand.Reader, params.q)
	if err != nil {
		t.Fatal(err)
	}
	eta2, err := rand.Int(rand.Reader, params.QSix())
	if err != nil {
		t.Fatal(err)
	}

	rc, err := paillier.GetRandomNumberInMultiplicativeGroup(params.N, rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	g := tecdsa.NewCurvePoint(params.curve.ScalarBaseMult(big.NewInt(1).Bytes()))

	u := message

	qEta2 := new(big.Int).Mod(new(big.Int).Mul(params.q, eta2), params.N)
	mask, err := privateKey.EncryptWithR(qEta2, rc)
	if err != nil {
		t.Fatal(err)
	}
	uEta1 := new(big.Int).Exp(u, eta1, params.NSquare())
	w := new(big.Int).Mod(new(big.Int).Mul(uEta1, mask.C), params.NSquare())

	r := tecdsa.NewCurvePoint(params.curve.ScalarBaseMult(eta1.Bytes()))

	// WHEN
	zkp, err := CommitZkpPi2(r, g, w, u, eta1, eta2, rc, params, rand.Reader)

	if err != nil {
		t.Fatal(err)
	}

	var tests = map[string]struct {
		verify         func() bool
		expectedResult bool
	}{
		"positive validation": {
			verify: func() bool {
				return zkp.Verify(params)
			},
			expectedResult: true,
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
