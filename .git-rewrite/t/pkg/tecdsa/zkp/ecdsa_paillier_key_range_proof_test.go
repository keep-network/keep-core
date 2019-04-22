package zkp

import (
	"crypto/rand"
	"math/big"
	"reflect"
	"testing"

	"github.com/keep-network/keep-core/pkg/internal/testutils"
	"github.com/keep-network/keep-core/pkg/tecdsa/curve"
	"github.com/keep-network/paillier"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

// TestEcdsaPaillierKeyRangeProofCommitValues creates a commitment and checks all
// commitment values against expected ones. This is not a full roundtrip test.
// We test the private commitment phase interface to make sure if anything goes
// wrong in future (e.g. curve implementation changes), we can isolate the
// problem easily. All expected values has been manually calculated basis on
// the [GGN16] paper.
func TestEcdsaPaillierKeyRangeProofCommitValues(t *testing.T) {
	mockRandom := testutils.NewMockRandReader(big.NewInt(10))
	// Following values are assigned to ZKP parameters as a result of
	// calling mockRandom:
	//
	// alpha=10
	// beta=11
	// rho=12
	// gamma=13

	parameters := &PublicParameters{
		N:      big.NewInt(1081), // 23 * 47
		NTilde: big.NewInt(253),  // 23 * 11
		h1:     big.NewInt(49),
		h2:     big.NewInt(22),
		q:      secp256k1.S256().Params().N,
		curve:  secp256k1.S256(),
	}

	secretEcdsaKeyShare := big.NewInt(13)
	publicEcdsaKeyShare := curve.NewPoint(
		secp256k1.S256().ScalarBaseMult(big.NewInt(11).Bytes()),
	)

	encryptedSecretEcdsaKeyShare := &paillier.Cypher{C: big.NewInt(12)}
	r := big.NewInt(14)

	commitment, err := CommitEcdsaPaillierKeyRange(
		secretEcdsaKeyShare,
		publicEcdsaKeyShare,
		encryptedSecretEcdsaKeyShare,
		r,
		parameters,
		mockRandom,
	)
	if err != nil {
		t.Fatal(err)
	}

	// 49^13 * 22^12 mod 253 = 55
	if commitment.z.Cmp(big.NewInt(55)) != 0 {
		t.Errorf("Unexpected z\nActual: %v\nExpected: 55", commitment.z)
	}

	// g^10
	expectedU1x, expectedU1y := secp256k1.S256().ScalarBaseMult(
		big.NewInt(10).Bytes(),
	)
	if commitment.u1.X.Cmp(expectedU1x) != 0 {
		t.Errorf("Unexpected u1 x coordinate")
	}
	if commitment.u1.Y.Cmp(expectedU1y) != 0 {
		t.Errorf("Unexpected u1 y coordinate")
	}

	// (1081+1)^10 * 11^1081 mod 1081^2 = 289613
	if commitment.u2.Cmp(big.NewInt(289613)) != 0 {
		t.Errorf("Unexpected u2\nActual: %v\nExpected: 289613", commitment.u2)
	}

	// 49^10 * 22^13 mod 253 = 176
	if commitment.u3.Cmp(big.NewInt(176)) != 0 {
		t.Errorf("Unexpected u3\nActual: %v\nExpected: 176", commitment.u3)
	}

	// e = hash(y, w, z, u1, u2, u3) =
	//     hash(g^11.X, g^11.Y, 12, 55, g^10.X, g^10.Y, 289613, 176)
	expectedHash := new(big.Int)
	expectedHash.SetString(
		"51984478426836913558864603258469889500681512521977850701426158002380794165890",
		10,
	)
	if commitment.e.Cmp(expectedHash) != 0 {
		t.Errorf("Unexpected e\nActual: %v\nExpected: %v",
			commitment.e,
			expectedHash,
		)
	}

	// e*13 + 10
	expectedS1 := new(big.Int)
	expectedS1.SetString(
		"675798219548879876265239842360108563508859662785712059118540054030950324156580",
		10,
	)
	if commitment.s1.Cmp(expectedS1) != 0 {
		t.Errorf("Unexpected s1\nActual: %v\nExpected: %v",
			commitment.s1,
			expectedS1,
		)
	}

	// 14^e * 11 mod 1081 = 605
	if commitment.s2.Cmp(big.NewInt(91)) != 0 {
		t.Errorf("Unexpected s2\nActual: %v\nExpected: 91", commitment.s2)
	}

	// e*12 + 13
	expectedS3 := new(big.Int)
	expectedS3.SetString(
		"623813741122042962706375239101638674008178150263734208417113896028569529990693",
		10,
	)
	if commitment.s3.Cmp(expectedS3) != 0 {
		t.Errorf("Unexpected s3\nActual: %v\nExpected: %v",
			commitment.s3,
			expectedS3,
		)
	}
}

// TestEcdsaPaillierKeyRangeProofVerification runs over the verification phase
// using supplied, hardcoded commitment values. This is not a full roundtrip
// test. We test the private verification phase interface to make sure if
// anything goes wrong in future (e.g. curve implementation changes), we can
// isolate the problem easily. All expected values has been manually calculated
// basis on the [GGN16] paper.
func TestEcdsaPaillierKeyRangeProofVerification(t *testing.T) {
	zkp := &EcdsaPaillierKeyRangeProof{
		s1: big.NewInt(22),
		s2: big.NewInt(17),
		s3: big.NewInt(63),
		e:  big.NewInt(881),
		z:  big.NewInt(191),
	}

	params := &PublicParameters{
		N:      big.NewInt(1081), // 23 * 47
		NTilde: big.NewInt(253),  // 23 * 11
		h1:     big.NewInt(11),
		h2:     big.NewInt(27),
		curve:  secp256k1.S256(),
		q:      secp256k1.S256().Params().N,
	}

	encryptedSecretEcdsaKeyShare := big.NewInt(674)
	publicEcdsaKeyShare := curve.NewPoint(
		secp256k1.S256().ScalarBaseMult(big.NewInt(10).Bytes()),
	)

	// u1 = g^s1 * y^-e = g^22 * (g^10)^-881 = g^{q-8810}
	expectedU1 := curve.NewPoint(secp256k1.S256().ScalarBaseMult(
		new(big.Int).Sub(params.q, big.NewInt(8788)).Bytes(),
	))
	actualU1 := zkp.evaluateU1Verification(publicEcdsaKeyShare, params)
	if !reflect.DeepEqual(expectedU1, actualU1) {
		t.Errorf(
			"Unexpected u1\nActual: %v\nExpected: %v",
			actualU1,
			expectedU1,
		)
	}

	// u2 = ((1081+1)^22 * 17^1081 * 674^-881) mod 1081^2
	expectedU2 := big.NewInt(227035)
	actualU2 := zkp.evaluateU2Verification(encryptedSecretEcdsaKeyShare, params)
	if expectedU2.Cmp(actualU2) != 0 {
		t.Errorf(
			"Unexpected u2\nActual: %v\nExpected: %v",
			actualU2,
			expectedU2,
		)
	}

	// u3 = (11^22 * 27^63 * 191^-881) mod 253
	expectedU3 := big.NewInt(44)
	actualU3 := zkp.evaluateU3Verification(params)
	if expectedU3.Cmp(actualU3) != 0 {
		t.Errorf(
			"Unexpected u3\nActual: %v\nExpected: %v",
			actualU3,
			expectedU3,
		)
	}
}

// TestEcdsaPaillierKeyRangeProofCommitAndVerify is a full roundtrip
// test of the ZKP, including generating public parameters, positive
// and negative validation scenarios.
func TestEcdsaPaillierKeyRangeProofCommitAndVerify(t *testing.T) {
	ellipticCurve := secp256k1.S256()
	paillierKey, parameters, err := createTestZkpParameters(ellipticCurve)
	if err != nil {
		t.Fatal(err)
	}

	r, err := paillier.GetRandomNumberInMultiplicativeGroup(parameters.N, rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	secretEcdsaKeyShare := big.NewInt(1410)
	encryptedSecretEcdsaKeyShare, err := paillierKey.EncryptWithR(
		secretEcdsaKeyShare,
		r,
	)
	if err != nil {
		t.Fatal(err)
	}
	publicEcdsaKeyShare := curve.NewPoint(
		ellipticCurve.ScalarBaseMult(secretEcdsaKeyShare.Bytes()),
	)

	commitment, err := CommitEcdsaPaillierKeyRange(
		secretEcdsaKeyShare,
		publicEcdsaKeyShare,
		encryptedSecretEcdsaKeyShare,
		r,
		parameters,
		rand.Reader,
	)
	if err != nil {
		t.Fatal(err)
	}

	var tests = map[string]struct {
		verify         func() bool
		expectedResult bool
	}{
		"positive validation": {
			verify: func() bool {
				return commitment.Verify(
					encryptedSecretEcdsaKeyShare,
					publicEcdsaKeyShare,
					parameters,
				)
			},
			expectedResult: true,
		},
		"negative validation - wrong encrypted secret DSA key share": {
			verify: func() bool {
				wrongEncryptedSecretEcdsaKeyShare := &paillier.Cypher{
					C: big.NewInt(1411),
				}
				return commitment.Verify(
					wrongEncryptedSecretEcdsaKeyShare,
					publicEcdsaKeyShare,
					parameters,
				)
			},
			expectedResult: false,
		},
		"negative validation - wrong public DSA key share": {
			verify: func() bool {
				wrongPublicEcdsaKeyShare := &curve.Point{
					X: big.NewInt(997),
					Y: big.NewInt(998),
				}
				return commitment.Verify(
					encryptedSecretEcdsaKeyShare,
					wrongPublicEcdsaKeyShare,
					parameters,
				)
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

// TestEcdsaPaillierKeyRangeProofParamsInRange runs a test of preliminary
// commitment validation parameters check. The check is a preliminary step to
// test if commitment is not corrupted (MiM attack).
func TestEcdsaPaillierKeyRangeProofParamsInRange(t *testing.T) {
	params := &PublicParameters{
		N:      big.NewInt(1081), // 23 * 47
		NTilde: big.NewInt(253),  // 23 * 11
		h1:     big.NewInt(11),
		h2:     big.NewInt(27),
		curve:  secp256k1.S256(),
		q:      secp256k1.S256().Params().N,
	}

	var tests = map[string]struct {
		modifyProof    func(zkp *EcdsaPaillierKeyRangeProof)
		expectedResult bool
	}{
		"positive parameters range validation": {
			modifyProof:    nil,
			expectedResult: true,
		},
		"negative validation - z less than 0": {
			modifyProof: func(zkp *EcdsaPaillierKeyRangeProof) {
				zkp.z = big.NewInt(-1)
			},
			expectedResult: false,
		},
		"negative validation - z equal NTilde": {
			modifyProof: func(zkp *EcdsaPaillierKeyRangeProof) {
				zkp.z = params.NTilde
			},
			expectedResult: false,
		},
		"negative validation - z greater than NTilde": {
			modifyProof: func(zkp *EcdsaPaillierKeyRangeProof) {
				zkp.z = new(big.Int).Add(params.NTilde, big.NewInt(1))
			},
			expectedResult: false,
		},
		"negative validation - u2 less than 0": {
			modifyProof: func(zkp *EcdsaPaillierKeyRangeProof) {
				zkp.u2 = big.NewInt(-1)
			},
			expectedResult: false,
		},
		"negative validation - u2 equal NSquare": {
			modifyProof: func(zkp *EcdsaPaillierKeyRangeProof) {
				zkp.u2 = params.NSquare()
			},
			expectedResult: false,
		},
		"negative validation - u2 greater than NSquare": {
			modifyProof: func(zkp *EcdsaPaillierKeyRangeProof) {
				zkp.u2 = new(big.Int).Add(params.NSquare(), big.NewInt(1))
			},
			expectedResult: false,
		},
		"negative validation - u3 less than 0": {
			modifyProof: func(zkp *EcdsaPaillierKeyRangeProof) {
				zkp.u3 = big.NewInt(-1)
			},
			expectedResult: false,
		},
		"negative validation - u3 equal NTilde": {
			modifyProof: func(zkp *EcdsaPaillierKeyRangeProof) {
				zkp.u3 = params.NTilde
			},
			expectedResult: false,
		},
		"negative validation - u3 greater than NTilde": {
			modifyProof: func(zkp *EcdsaPaillierKeyRangeProof) {
				zkp.u3 = new(big.Int).Add(params.NTilde, big.NewInt(1))
			},
			expectedResult: false,
		},
		"negative validation - s2 less than 0": {
			modifyProof: func(zkp *EcdsaPaillierKeyRangeProof) {
				zkp.u3 = big.NewInt(-1)
			},
			expectedResult: false,
		},
		"negative validation - s2 equal N": {
			modifyProof: func(zkp *EcdsaPaillierKeyRangeProof) {
				zkp.u3 = params.N
			},
			expectedResult: false,
		},
		"negative validation - s2 greater than N": {
			modifyProof: func(zkp *EcdsaPaillierKeyRangeProof) {
				zkp.u3 = new(big.Int).Add(params.N, big.NewInt(1))
			},
			expectedResult: false,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			zkp := &EcdsaPaillierKeyRangeProof{
				z:  big.NewInt(250),
				u2: big.NewInt(224),
				u3: big.NewInt(123),
				s2: big.NewInt(17),
			}

			if test.modifyProof != nil {
				test.modifyProof(zkp)
			}

			actualResult := zkp.allParametersInRange(params)

			if actualResult != test.expectedResult {
				t.Fatalf(
					"expected: %v\nactual: %v",
					test.expectedResult,
					actualResult,
				)
			}
		})
	}
}
