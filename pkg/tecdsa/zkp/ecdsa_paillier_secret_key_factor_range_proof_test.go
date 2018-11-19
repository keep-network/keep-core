package zkp

import (
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/keep-network/keep-core/pkg/internal/testutils"
	"github.com/keep-network/paillier"
)

// TestEcdsaPaillierSecretKeyFactorRangeProofCommitValues creates a commitment and
// checks all commitment values against expected ones.
//
// This is not a full roundtrip test. We test the private commitment phase
// interface to make sure if anything goes wrong in future (e.g. curve
// implementation changes), we can isolate the problem easily.
// All expected values has been manually calculated based on the [GGN16] paper.
func TestEcdsaPaillierSecretKeyFactorRangeProofCommitValues(t *testing.T) {
	// GIVEN
	mockRandom := testutils.NewMockRandReader(big.NewInt(10))
	// Following values are assigned to ZKP parameters as a result of
	// calling mockRandom:
	//
	// alpha=10
	// beta=11
	// rho=12
	// gamma=13

	params := generateTestPublicParams()

	factor := big.NewInt(8)

	secretEcdsaKeyMultiple := &paillier.Cypher{C: big.NewInt(9)}
	secretEcdsaKey := &paillier.Cypher{C: big.NewInt(7)}
	secretEcdsaKeyFactor := &paillier.Cypher{C: big.NewInt(9)}

	r := big.NewInt(7)

	// WHEN
	zkp, err := CommitEcdsaPaillierSecretKeyFactorRange(secretEcdsaKeyMultiple,
		secretEcdsaKey, secretEcdsaKeyFactor, factor, r, params, mockRandom)
	if err != nil {
		t.Fatal(err)
	}

	// THEN

	// 1082^10 * 11^1081 mod 1168561 = 289613
	if zkp.z.Cmp(big.NewInt(289613)) != 0 {
		t.Errorf("Unexpected z\nActual: %v\nExpected: 289613", zkp.z)
	}

	// 7^10 mod 1168561 = 852048
	if zkp.v.Cmp(big.NewInt(852048)) != 0 {
		t.Errorf("Unexpected v\nActual: %v\nExpected: 852048", zkp.v)
	}

	// 20535^8 * 20919^12  mod 25651 = 7002
	if zkp.u1.Cmp(big.NewInt(7002)) != 0 {
		t.Errorf("Unexpected u1\nActual: %v\nExpected: 7002", zkp.u1)
	}

	// 20535^10 * 20919^13 mod 25651 = 1102
	if zkp.u2.Cmp(big.NewInt(1102)) != 0 {
		t.Errorf("Unexpected u2\nActual: %v\nExpected: 1102", zkp.u2)
	}

	// hash(9, 7, 8, 289613, 7002, 1101, 852048) = expectedHash
	expectedHash, _ := new(big.Int).SetString("5682122828670513769378897932249359273638734078677473798462399325423003506113", 10)
	if zkp.e.Cmp(expectedHash) != 0 {
		t.Errorf("Unexpected e\nActual: %v\nExpected: %v", zkp.e, expectedHash)
	}

	// expectedHash * 8 + 10 = expectedS1
	expectedS1, _ := new(big.Int).SetString("45456982629364110155031183457994874189109872629419790387699194603384028048914", 10)
	if zkp.s1.Cmp(expectedS1) != 0 {
		t.Errorf("Unexpected s1\nActual: %v\nExpected: %v", zkp.s1, expectedS1)
	}

	// 7^expectedHash * 11 mod 1081 = 869
	if zkp.s2.Cmp(big.NewInt(869)) != 0 {
		t.Errorf("Unexpected s2\nActual: %v\nExpected: 869", zkp.s2)
	}

	// expectedHash * 12 + 13
	expectedS3, _ := new(big.Int).SetString("68185473944046165232546775186992311283664808944129685581548791905076042073369", 10)
	if zkp.s3.Cmp(expectedS3) != 0 {
		t.Errorf("Unexpected s3\nActual: %v\nExpected: %v", zkp.s3, expectedS3)
	}
}

func TestEcdsaPaillierSecretKeyFactorRangeProofVerification(t *testing.T) {
	//GIVEN
	params := generateTestPublicParams()

	secretEcdsaKeyMultiple := &paillier.Cypher{C: big.NewInt(729674)}
	secretEcdsaKey := &paillier.Cypher{C: big.NewInt(133808)}
	secretEcdsaKeyFactor := &paillier.Cypher{C: big.NewInt(688361)}

	expectedZ := big.NewInt(289613)
	expectedV := big.NewInt(285526)
	expectedU2 := big.NewInt(1102)

	e, _ := new(big.Int).SetString("28665131959061509990138847422722847282246667596979352654045645230544684705784", 10)
	s1, _ := new(big.Int).SetString("315316451549676609891527321649951320104713343566772879194502097535991531763634", 10)
	s3, _ := new(big.Int).SetString("343981583508738119881666169072674167386960011163752231848547742766536216469421", 10)

	zkp := &EcdsaPaillierSecretKeyFactorRangeProof{
		z:  big.NewInt(289613),
		v:  big.NewInt(285526),
		u1: big.NewInt(10797),
		u2: big.NewInt(1102),

		e: e,

		s1: s1,
		s2: big.NewInt(986),
		s3: s3,
	}

	//WHEN
	actualZ := evaluateVerificationZ(secretEcdsaKeyFactor, zkp.s1, zkp.s2, zkp.e, params)
	actualV := evaluateVerificationV(secretEcdsaKeyMultiple, secretEcdsaKey,
		zkp.s1, zkp.e, params)
	actualU2 := evaluateVerificationU2(zkp.u1, zkp.s1, zkp.s3, zkp.e, params)
	verificationResult := zkp.Verify(secretEcdsaKeyMultiple, secretEcdsaKey, secretEcdsaKeyFactor, params)

	//THEN
	if expectedZ.Cmp(actualZ) != 0 {
		t.Errorf("Unexpected Z\nActual: %v\nExpected: %v", actualZ, expectedZ)
	}
	if expectedV.Cmp(actualV) != 0 {
		t.Errorf("Unexpected Z\nActual: %v\nExpected: %v", actualV, expectedV)
	}
	if expectedU2.Cmp(actualU2) != 0 {
		t.Errorf("Unexpected U2\nActual: %v\nExpected: %v", actualU2, expectedU2)
	}
	if !verificationResult {
		t.Errorf("Verification failed")
	}
}

func TestEcdsaPaillierSecretKeyFactorRangeProofCommitAndVerify(t *testing.T) {
	// GIVEN
	message := big.NewInt(430)

	paillierKey, params, err := createTestZkpParameters(secp256k1.S256())
	if err != nil {
		t.Fatal(err)
	}

	r, err := paillier.GetRandomNumberInMultiplicativeGroup(params.N, rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	secretDsaKey, err := paillierKey.EncryptWithR(message, r)
	if err != nil {
		t.Fatal(err)
	}

	factor, err := rand.Int(rand.Reader, params.q)
	if err != nil {
		t.Fatalf("could not generate eta [%v]", err)
	}

	secretEcdsaKeyMultiple := paillierKey.Mul(secretDsaKey, factor)
	secretEcdsaKeyFactor, err := paillierKey.EncryptWithR(factor, r)
	if err != nil {
		t.Fatal(err)
	}

	// WHEN
	zkp, err := CommitEcdsaPaillierSecretKeyFactorRange(secretEcdsaKeyMultiple,
		secretDsaKey, secretEcdsaKeyFactor, factor, r, params, rand.Reader)
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
				return zkp.Verify(
					secretEcdsaKeyMultiple,
					secretDsaKey,
					secretEcdsaKeyFactor,
					params,
				)
			},
			expectedResult: true,
		},
		"negative validation - wrong secret dsa key multiple": {
			verify: func() bool {
				wrongEncryptedSecretDsaKeyMultiple := &paillier.Cypher{C: big.NewInt(1411)}
				return zkp.Verify(
					wrongEncryptedSecretDsaKeyMultiple,
					secretDsaKey,
					secretEcdsaKeyFactor,
					params,
				)
			},
			expectedResult: false,
		},
		"negative validation - wrong secret dsa key": {
			verify: func() bool {
				wrongSecretDsaKey := &paillier.Cypher{C: big.NewInt(856)}
				return zkp.Verify(
					secretEcdsaKeyMultiple,
					wrongSecretDsaKey,
					secretEcdsaKeyFactor,
					params,
				)
			},
			expectedResult: false,
		},
		"negative validation - wrong secret dsa key factor": {
			verify: func() bool {
				wrongFactor := &paillier.Cypher{C: big.NewInt(798)}
				return zkp.Verify(
					secretEcdsaKeyMultiple,
					secretDsaKey,
					wrongFactor,
					params,
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

func generateTestPublicParams() *PublicParameters {
	return &PublicParameters{
		N:      big.NewInt(1081),  // 23 * 47
		NTilde: big.NewInt(25651), // 23 * 11

		h1: big.NewInt(20535),
		h2: big.NewInt(20919),

		q:     secp256k1.S256().Params().N,
		curve: secp256k1.S256(),
	}
}
