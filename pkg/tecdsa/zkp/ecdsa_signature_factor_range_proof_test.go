package zkp

import (
	"crypto/rand"
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/keep-network/keep-core/pkg/internal/testutils"
	"github.com/keep-network/keep-core/pkg/tecdsa/curve"
	"github.com/keep-network/paillier"
)

// TestCommitEcdsaSignatureFactorRangeProofValues creates a commitment and checks
// all commitment values against expected ones.
//
// This is not a full roundtrip test. We test the private commitment phase
// interface to make sure if anything goes wrong in future (e.g. curve
// implementation changes), we can isolate the problem easily.
// All expected values has been manually calculated based on the [GGN16] paper.
func TestCommitEcdsaSignatureFactorRangeProofValues(t *testing.T) {
	// GIVEN
	mockRandom := testutils.NewMockRandReader(big.NewInt(10))
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

	params := generateTestPublicParams()

	signatureFactorSecret := big.NewInt(3)
	signatureFactorMask := big.NewInt(5)

	paillierR := big.NewInt(7)

	signatureFactorPublic := curve.NewPoint(params.curve.ScalarBaseMult(signatureFactorSecret.Bytes()))
	signatureUnmask := &paillier.Cypher{C: big.NewInt(9)}
	secretKeyFactor := &paillier.Cypher{C: big.NewInt(7)}

	// WHEN
	zkp, err := CommitEcdsaSignatureFactorRangeProof(
		signatureFactorPublic,
		signatureUnmask,
		secretKeyFactor,
		signatureFactorSecret,
		signatureFactorMask,
		paillierR,
		params,
		mockRandom,
	)
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
	expectedU1 := curve.NewPoint(params.curve.ScalarBaseMult(big.NewInt(10).Bytes()))
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
	// 7^10 * 1082^(q*16) * 14^1081 mod 1168561 = 583302
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

	// hash(9, 7, 19261, 17194, g^10, 289613, 19547, 583302, 5738, 21032) = expectedHash
	expectedHash, _ := new(big.Int).SetString("25499118291731768019248475116926786227164997290767802448110811852901126765839", 10)
	if zkp.e.Cmp(expectedHash) != 0 {
		t.Errorf("Unexpected e\nActual: %v\nExpected: %v", zkp.e, expectedHash)
	}

	// 3 * e + 10 = expectedS1
	expectedS1, _ := new(big.Int).SetString("76497354875195304057745425350780358681494991872303407344332435558703380297527", 10)
	if zkp.s1.Cmp(expectedS1) != 0 {
		t.Errorf("Unexpected s1\nActual: %v\nExpected: %v", zkp.s1, expectedS1)
	}
	// 18 * e + 12 = expectedS2
	expectedS2, _ := new(big.Int).SetString("458984129251171824346472552104682152088969951233820444065994613352220281785114", 10)
	if zkp.s2.Cmp(expectedS2) != 0 {
		t.Errorf("Unexpected s2\nActual: %v\nExpected: %v", zkp.s2, expectedS2)
	}
	// 7^e * 14 mod 1081 = 811
	if zkp.t1.Cmp(big.NewInt(811)) != 0 {
		t.Errorf("Unexpected t1\nActual: %v\nExpected: 811", zkp.t1)
	}
	// 5 * e + 16 = expectedT2
	expectedT2, _ := new(big.Int).SetString("127495591458658840096242375584633931135824986453839012240554059264505633829211", 10)
	if zkp.t2.Cmp(expectedT2) != 0 {
		t.Errorf("Unexpected t2\nActual: %v\nExpected: %v", zkp.t2, expectedT2)
	}
	// 19 * e + 17 = expectedT3
	expectedT3, _ := new(big.Int).SetString("484483247542903592365721027221608938316134948524588246514105425205121408550958", 10)
	if zkp.t3.Cmp(expectedT3) != 0 {
		t.Errorf("Unexpected t3\nActual: %v\nExpected: %v", zkp.t3, expectedT3)
	}
}

func TestEcdsaSignatureFactorRangeProofVerification(t *testing.T) {
	//GIVEN
	params := generateTestPublicParams()

	signatureFactorSecret := big.NewInt(3)
	signatureFactorPublic := curve.NewPoint(params.curve.ScalarBaseMult(signatureFactorSecret.Bytes()))
	secretKeyFactor := &paillier.Cypher{C: big.NewInt(7)}
	// `w` is calculated from:
	// w = u^η1 * Γ^qη2 * rc^N mod N^2
	//
	// With following values:
	// u = 7
	// η1 = 3
	// Γ = params.G()
	// q = params.q
	// η2 = 5
	// rc = 7
	// N = params.N
	signatureUnmask := &paillier.Cypher{C: big.NewInt(1005955)}

	e, _ := new(big.Int).SetString("5409986892274499674842741357332043231522465088980131051275882178073046897869", 10)
	s1, _ := new(big.Int).SetString("16229960676823499024528224071996129694567395266940393153827646534219140693617", 10)
	s2, _ := new(big.Int).SetString("97379764060940994147169344431976778167404371601642358922965879205314844161654", 10)
	t2, _ := new(big.Int).SetString("27049934461372498374213706786660216157612325444900655256379410890365234489361", 10)
	t3, _ := new(big.Int).SetString("102789750953215493822012085789308821398926836690622489974241761383387891059528", 10)

	zkp := &EcdsaSignatureFactorRangeProof{
		z1: big.NewInt(19261),
		z2: big.NewInt(17194),

		u1: curve.NewPoint(params.curve.ScalarBaseMult(big.NewInt(10).Bytes())),
		u2: big.NewInt(289613),
		u3: big.NewInt(19547),

		v1: big.NewInt(583302),
		v2: big.NewInt(5738),
		v3: big.NewInt(21032),

		e:  e,
		s1: s1,
		s2: s2,

		t1: big.NewInt(213),
		t2: t2,
		t3: t3,
	}

	expectedU1 := curve.NewPoint(params.curve.ScalarBaseMult(big.NewInt(10).Bytes()))
	expectedU3 := big.NewInt(19547)
	expectedV1 := big.NewInt(583302)
	expectedV3 := big.NewInt(21032)

	// WHEN
	actualU1 := zkp.evaluateVerificationU1(signatureFactorPublic, params)
	actualU3 := zkp.evaluateVerificationU3(params)
	actualV1 := zkp.evaluateVerificationV1(signatureUnmask, secretKeyFactor, params)
	actualV3 := zkp.evaluateVerificationV3(params)

	// THEN
	if !reflect.DeepEqual(expectedU1, actualU1) {
		t.Errorf("Unexpected u1\nActual: %v\nExpected: %v", actualU1, expectedU1)
	}
	if actualU3.Cmp(expectedU3) != 0 {
		t.Errorf("Unexpected u3\nActual: %v\nExpected: %v", actualU3, expectedU3)
	}
	if actualV1.Cmp(expectedV1) != 0 {
		t.Errorf("Unexpected v1\nActual: %v\nExpected: %v", actualV1, expectedV1)
	}
	if actualV3.Cmp(expectedV3) != 0 {
		t.Errorf("Unexpected v3\nActual: %v\nExpected: %v", actualV3, expectedV3)
	}
}

func TestEcdsaSignatureFactorRangeProofCommitAndVerify(t *testing.T) {
	// GIVEN
	privateKey, params, err := createTestZkpParameters(secp256k1.S256())
	if err != nil {
		t.Fatal(err)
	}

	signatureFactorSecret, err := rand.Int(rand.Reader, params.q) // factor from ZKP PI2,1
	if err != nil {
		t.Fatal(err)
	}
	signatureFactorMask, err := rand.Int(rand.Reader, params.QPow6()) // factor from ZKP PI2,1
	if err != nil {
		t.Fatal(err)
	}

	paillierR, err := paillier.GetRandomNumberInMultiplicativeGroup(params.N, rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	signatureFactorPublic := curve.NewPoint(params.curve.ScalarBaseMult(signatureFactorSecret.Bytes())) // eta1CurvePoint
	secretKeyFactor, _ := privateKey.EncryptWithR(big.NewInt(5), paillierR)                             // encryptedFactor1 (rho from round 1), u = E(ρ)

	// D(w) = η1*D(u) + q*η2 = η1*ρ + q*η2
	// w = E(D(w)) = E(η1*D(u) + q*η2) = E(η1*D(u) + E(q*η2) = E(η1*u) + E(q*η2)
	// E(η1*u)
	encryptedUEta1 := privateKey.PublicKey.Mul(secretKeyFactor, signatureFactorSecret)

	// E(q*η2)
	qEta2 := new(big.Int).Mod(new(big.Int).Mul(params.q, signatureFactorMask), params.N)
	encryptedQEta2, err := privateKey.EncryptWithR(qEta2, paillierR)
	if err != nil {
		t.Fatal(err)
	}

	signatureUnmask := privateKey.PublicKey.Add(encryptedUEta1, encryptedQEta2)

	// WHEN
	zkp, err := CommitEcdsaSignatureFactorRangeProof(
		signatureFactorPublic,
		signatureUnmask,
		secretKeyFactor,
		signatureFactorSecret,
		signatureFactorMask,
		paillierR,
		params,
		rand.Reader,
	)

	if err != nil {
		t.Fatal(err)
	}

	// THEN
	var tests = map[string]struct {
		verify         func() bool
		expectedResult bool
	}{
		"positive validation": {
			verify: func() bool {
				return zkp.Verify(signatureFactorPublic, signatureUnmask, secretKeyFactor, params)
			},
			expectedResult: true,
		},
		"negative validation - wrong signatureFactorPublic": {
			verify: func() bool {
				wrongSignatureFactorPublic := curve.NewPoint(big.NewInt(3), big.NewInt(4))
				return zkp.Verify(wrongSignatureFactorPublic, signatureUnmask, secretKeyFactor, params)
			},
			expectedResult: false,
		},
		"negative validation - wrong signatureUnmask": {
			verify: func() bool {
				wrongSignatureUnmask := privateKey.PublicKey.Mul(signatureUnmask, big.NewInt(2))
				return zkp.Verify(signatureFactorPublic, wrongSignatureUnmask, secretKeyFactor, params)
			},
			expectedResult: false,
		},
		"negative validation - wrong secretKeyFactor": {
			verify: func() bool {
				wrongSecretKeyFactor := privateKey.PublicKey.Mul(secretKeyFactor, big.NewInt(2))
				return zkp.Verify(signatureFactorPublic, signatureUnmask, wrongSecretKeyFactor, params)
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
