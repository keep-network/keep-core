package zkp

import (
	"crypto/rand"
	"math/big"
	"reflect"
	"testing"

	"github.com/keep-network/keep-core/pkg/tecdsa"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/keep-network/paillier"
)

// TestZKP2CommitValues creates a commitment and checks all
// commitment values against expected ones.
//
// This is not a full roundtrip test. We test the private commitment phase
// interface to make sure if anything goes wrong in future (e.g. curve
// implementation changes), we can isolate the problem easily.
// All expected values has been manually calculated based on the [GGN16] paper.
func TestZKP2CommitValues(t *testing.T) {
	// GIVEN
	mockRandom := &mockRandReader{
		counter: big.NewInt(10),
	}
	// Following values are assigned to ZKP parameters as a result of
	// calling mockRandom:
	//
	// alpha=10
	// beta=11
	// gamma=12
	// delta=13
	// mu=14
	// nu=15
	// theta=16
	// tau=17
	// rho1=18
	// rho2=19

	params := &PublicParameters{
		N:      big.NewInt(1081),  // 23 * 47
		NTilde: big.NewInt(25651), // 23 * 11

		h1: big.NewInt(20535),
		h2: big.NewInt(20919),

		q:     secp256k1.S256().Params().N,
		curve: secp256k1.S256(),
	}

	eta1 := big.NewInt(3)
	eta2 := big.NewInt(5)

	rc := big.NewInt(7)

	g := params.CurveBasePoint()
	r := tecdsa.NewCurvePoint(params.curve.ScalarBaseMult(eta1.Bytes()))
	u := big.NewInt(7)
	w := big.NewInt(9)

	// WHEN
	zkp, err := CommitZkpPi2(g, r, w, u, eta1, eta2, rc, params, mockRandom)
	if err != nil {
		t.Fatal(err)
	}

	// THEN

	// 20535^3 * 20919^18 mod 25651 = 19261
	if zkp.z1.Cmp(big.NewInt(19261)) != 0 {
		t.Errorf("Unexpected z1\nActual: %v\nExpected: 19261", zkp.z1)
	}
	// 20535^5 * 20919^19 mod 25651 = 17194
	if zkp.z2.Cmp(big.NewInt(17194)) != 0 {
		t.Errorf("Unexpected z2\nActual: %v\nExpected: 17194", zkp.z1)
	}
	// g^10
	expectedU1 := tecdsa.NewCurvePoint(params.curve.ScalarBaseMult(big.NewInt(10).Bytes()))
	if !reflect.DeepEqual(expectedU1, zkp.u1) {
		t.Errorf("Unexpected u1\nActual: %v\nExpected: %v", zkp.u1, expectedU1)
	}
	// 1082^10 * 11^1081 mod 1168561 = 289613
	if zkp.u2.Cmp(big.NewInt(289613)) != 0 {
		t.Errorf("Unexpected u2\nActual: %v\nExpected: 289613", zkp.u2)
	}
	// 20535^10 * 20919^12 mod 25651 = 19547
	if zkp.u3.Cmp(big.NewInt(19547)) != 0 {
		t.Errorf("Unexpected u3\nActual: %v\nExpected: 19547", zkp.u3)
	}
	// 7^10 * 1082^(115792089237316195423570985008687907852837564279074904382605163141518161494337*16) * 14^1081 mod 1168561 = 583302
	if zkp.v1.Cmp(big.NewInt(583302)) != 0 {
		t.Errorf("Unexpected v1\nActual: %v\nExpected: 583302", zkp.v1)
	}
	// 20535^13 * 20919^15 mod 25651 = 5738
	if zkp.v2.Cmp(big.NewInt(5738)) != 0 {
		t.Errorf("Unexpected v2\nActual: %v\nExpected: 5738", zkp.v2)
	}
	// 20535^16 * 20919^17 mod 25651 = 21032
	if zkp.v3.Cmp(big.NewInt(21032)) != 0 {
		t.Errorf("Unexpected v3\nActual: %v\nExpected: 21032", zkp.v3)
	}

	// hash(g, 9, 7, 19261, 17194, g^10, 289613, 19547, 583302, 5738, 21032) = expectedHash
	expectedHash, _ := new(big.Int).SetString("5409986892274499674842741357332043231522465088980131051275882178073046897869", 10)
	if zkp.e.Cmp(expectedHash) != 0 {
		t.Errorf("Unexpected e\nActual: %v\nExpected: %v", zkp.e, expectedHash)
	}

	// expectedHash * 3 + 10 = expectedS1
	expectedS1, _ := new(big.Int).SetString("16229960676823499024528224071996129694567395266940393153827646534219140693617", 10)
	if zkp.s1.Cmp(expectedS1) != 0 {
		t.Errorf("Unexpected s1\nActual: %v\nExpected: %v", zkp.s1, expectedS1)
	}
	// expectedHash * 18 + 12 = expectedS2
	expectedS2, _ := new(big.Int).SetString("97379764060940994147169344431976778167404371601642358922965879205314844161654", 10)
	if zkp.s2.Cmp(expectedS2) != 0 {
		t.Errorf("Unexpected s2\nActual: %v\nExpected: %v", zkp.s2, expectedS2)
	}
	// 7^expectedHash * 14 mod 1081 = 213
	if zkp.t1.Cmp(big.NewInt(213)) != 0 {
		t.Errorf("Unexpected t1\nActual: %v\nExpected: 213", zkp.t1)
	}
	// expectedHash * 5 + 16 = expectedT2
	expectedT2, _ := new(big.Int).SetString("27049934461372498374213706786660216157612325444900655256379410890365234489361", 10)
	if zkp.t2.Cmp(expectedT2) != 0 {
		t.Errorf("Unexpected t2\nActual: %v\nExpected: %v", zkp.t2, expectedT2)
	}
	// expectedHash * 19 + 17 = 213
	expectedT3, _ := new(big.Int).SetString("102789750953215493822012085789308821398926836690622489974241761383387891059528", 10)
	if zkp.t3.Cmp(expectedT3) != 0 {
		t.Errorf("Unexpected t3\nActual: %v\nExpected: %v", zkp.t3, expectedT3)
	}
}
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
