package zkp

import (
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/keep-network/paillier"
)

// TestDsaPaillierSecretKeyFactorRangeProofCommitValues creates a commitment and checks all
// commitment values against expected ones.
//
// This is not a full roundtrip test. We test the private commitment phase
// interface to make sure if anything goes wrong in future (e.g. curve
// implementation changes), we can isolate the problem easily.
// All expected values has been manually calculated based on the [GGN16] paper.
func TestDsaPaillierSecretKeyFactorRangeProofCommitValues(t *testing.T) {
	// GIVEN
	mockRandom := &mockRandReader{
		counter: big.NewInt(10),
	}
	// Following values are assigned to ZKP parameters as a result of
	// calling mockRandom:
	//
	// alpha=10
	// beta=11
	// rho=12
	// gamma=13

	params := generateTestPublicParams()

	factor := big.NewInt(8)

	encryptedSecretDsaKeyMultiple := &paillier.Cypher{C: big.NewInt(9)}
	encryptedSecretDsaKey := &paillier.Cypher{C: big.NewInt(7)}
	encryptedFactor := &paillier.Cypher{C: big.NewInt(9)}

	r := big.NewInt(7)

	// WHEN
	zkp, err := CommitDsaPaillierSecretKeyFactorRange(encryptedSecretDsaKeyMultiple,
		encryptedSecretDsaKey, encryptedFactor, factor, r, params, mockRandom)
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

func TestDsaPaillierSecretKeyFactorRangeProofVerification(t *testing.T) {
	//GIVEN
	params := generateTestPublicParams()

	encryptedSecretDsaKeyMultiple := &paillier.Cypher{C: big.NewInt(729674)}
	encryptedSecretDsaKey := &paillier.Cypher{C: big.NewInt(133808)}
	encryptedFactor := &paillier.Cypher{C: big.NewInt(688361)}

	expectedZ := big.NewInt(289613)
	expectedV := big.NewInt(285526)
	expectedU2 := big.NewInt(1102)

	e, _ := new(big.Int).SetString("28665131959061509990138847422722847282246667596979352654045645230544684705784", 10)
	s1, _ := new(big.Int).SetString("315316451549676609891527321649951320104713343566772879194502097535991531763634", 10)
	s3, _ := new(big.Int).SetString("343981583508738119881666169072674167386960011163752231848547742766536216469421", 10)

	zkp := &DsaPaillierSecretKeyFactorRangeProof{
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
	actualZ := evaluateVerificationZ(encryptedFactor, zkp.s1, zkp.s2, zkp.e, params)
	actualV := evaluateVerificationV(encryptedSecretDsaKeyMultiple, encryptedSecretDsaKey,
		zkp.s1, zkp.e, params)
	actualU2 := evaluateVerificationU2(zkp.u1, zkp.s1, zkp.s3, zkp.e, params)
	verificationResult := zkp.Verify(encryptedSecretDsaKeyMultiple, encryptedSecretDsaKey, encryptedFactor, params)

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

func TestDsaPaillierSecretKeyFactorRangeProofRoundTrip(t *testing.T) {
	// GIVEN
	message := big.NewInt(430)

	p, _ := new(big.Int).SetString("104479735358598948369258156463683391052543755432914893102752306517616376250927", 10)
	q, _ := new(big.Int).SetString("110280671641689691092051226222060939019447720119674706500089479951904142152567", 10)
	paillierKey := paillier.CreatePrivateKey(p, q)

	params, err := GeneratePublicParameters(paillierKey.N, secp256k1.S256())
	if err != nil {
		t.Fatal(err)
	}

	r, err := paillier.GetRandomNumberInMultiplicativeGroup(params.N, rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	encryptedSecretDsaKey, err := paillierKey.EncryptWithR(message, r)
	if err != nil {
		t.Fatal(err)
	}

	factor, err := rand.Int(rand.Reader, params.q)
	if err != nil {
		t.Fatalf("could not generate eta [%v]", err)
	}

	// An encrypted plaintext raised to the power of another plaintext will
	// decrypt to the product of the two plaintexts:
	// (E(c2))^η = E(c2 * η)
	encryptedSecretDsaKeyMultiple := &paillier.Cypher{C: new(big.Int).Exp(encryptedSecretDsaKey.C, factor, params.NSquare())}

	encryptedFactor, err := paillierKey.EncryptWithR(factor, r)
	t.Logf("encryptedfactor: %s", encryptedFactor.C)
	if err != nil {
		t.Fatal(err)
	}

	// WHEN
	zkp, err := CommitDsaPaillierSecretKeyFactorRange(encryptedSecretDsaKeyMultiple,
		encryptedSecretDsaKey, encryptedFactor, factor, r, params, rand.Reader)
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
					encryptedSecretDsaKeyMultiple,
					encryptedSecretDsaKey,
					encryptedFactor,
					params,
				)
			},
			expectedResult: true,
		},
		"negative validation - wrong encrypted secret dsa key multiple": {
			verify: func() bool {
				wrongEncryptedSecretDsaKeyMultiple := &paillier.Cypher{C: big.NewInt(1411)}
				return zkp.Verify(
					wrongEncryptedSecretDsaKeyMultiple,
					encryptedSecretDsaKey,
					encryptedFactor,
					params,
				)
			},
			expectedResult: false,
		},
		"negative validation - wrong encrypted secret dsa key": {
			verify: func() bool {
				wrongEncryptedSecretDsaKey := &paillier.Cypher{C: big.NewInt(856)}
				return zkp.Verify(
					encryptedSecretDsaKeyMultiple,
					wrongEncryptedSecretDsaKey,
					encryptedFactor,
					params,
				)
			},
			expectedResult: false,
		},
		"negative validation - wrong encrypted factor": {
			verify: func() bool {
				wrongEncryptedFactor := &paillier.Cypher{C: big.NewInt(798)}
				return zkp.Verify(
					encryptedSecretDsaKeyMultiple,
					encryptedSecretDsaKey,
					wrongEncryptedFactor,
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
