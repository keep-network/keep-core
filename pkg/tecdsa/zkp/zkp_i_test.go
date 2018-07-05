package zkp

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/keep-network/keep-core/pkg/tecdsa"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

func TestPIiCommitValues(t *testing.T) {
	mockRandom := &mockRandReader{
		counter: big.NewInt(10),
	}
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

	y := big.NewInt(11)
	w := big.NewInt(12)
	eta := big.NewInt(13)
	r := big.NewInt(14)

	commitment, err := CommitPIi(y, w, eta, r, parameters, mockRandom)
	if err != nil {
		t.Fatal(err)
	}

	// 49^13 * 22^12 mod 253 = 55
	if !reflect.DeepEqual(commitment.z, big.NewInt(55)) {
		t.Errorf("Unexpected z\nActual: %v\nExpected: 55", commitment.z)
	}

	// g^10
	expectedU1x, expectedU1y := secp256k1.S256().ScalarBaseMult(
		big.NewInt(10).Bytes(),
	)
	if !reflect.DeepEqual(commitment.u1.X, expectedU1x) {
		t.Errorf("Unexpected u1 x coordinate")
	}
	if !reflect.DeepEqual(commitment.u1.Y, expectedU1y) {
		t.Errorf("Unexpected u1 y coordinate")
	}

	// (1081+1)^10 * 11^1081 mod 1081^2 = 289613
	if !reflect.DeepEqual(commitment.u2, big.NewInt(289613)) {
		t.Errorf("Unexpected u2\nActual: %v\nExpected: 289613", commitment.u2)
	}

	// 49^10 * 22^13 mod 253 = 176
	if !reflect.DeepEqual(commitment.u3, big.NewInt(176)) {
		t.Errorf("Unexpected u3\nActual: %v\nExpected: 176", commitment.u3)
	}

	// e = hash(g, y, w, z, u1, u2, u3) =
	//     hash(1082, 11, 12, 55, g^10.x, g^10.y, 289613, 176)
	expectedHash := new(big.Int)
	expectedHash.SetString(
		"59167403082436634448058111708361841704129646999348477739812569953930856100700",
		10,
	)
	if !reflect.DeepEqual(commitment.e, expectedHash) {
		t.Errorf("Unexpected e\nActual: %v\nExpected: %v",
			commitment.e,
			expectedHash,
		)
	}

	// e*13 + 10
	expectedS1 := new(big.Int)
	expectedS1.SetString(
		"769176240071676247824755452208703942153685410991530210617563409401101129309110",
		10,
	)
	if !reflect.DeepEqual(commitment.s1, expectedS1) {
		t.Errorf("Unexpected s1\nActual: %v\nExpected: %v",
			commitment.s1,
			expectedS1,
		)
	}

	// 14^e * 11 mod 1081 = 267
	if !reflect.DeepEqual(commitment.s2, big.NewInt(287)) {
		t.Errorf("Unexpected s2\nActual: %v\nExpected: 287", commitment.s2)
	}

	// e*12 + 13
	expectedS3 := new(big.Int)
	expectedS3.SetString(
		"710008836989239613376697340500342100449555763992181732877750839447170273208413",
		10,
	)
	if !reflect.DeepEqual(commitment.s3, expectedS3) {
		t.Errorf("Unexpected s3\nActual: %v\nExpected: %v",
			commitment.s3,
			expectedS3,
		)
	}
}

func TestPIiVerificationValues(t *testing.T) {
	zkp := &PIi{
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
	}

	w := big.NewInt(674)

	y := tecdsa.NewCurvePoint(
		secp256k1.S256().ScalarBaseMult(big.NewInt(10).Bytes()),
	)

	// u1 = g^s1 * y^-e = g^22 * (g^10)^-881 = g^-8832
	expectedU1 := tecdsa.NewCurvePoint(
		secp256k1.S256().ScalarBaseMult(big.NewInt(-8832).Bytes()),
	)
	actualU1 := zkp.u1Verification(y, params)
	if !reflect.DeepEqual(expectedU1, actualU1) {
		t.Errorf(
			"Unexpected u1\nActual: %v\nExpected: %v",
			actualU1,
			expectedU1,
		)
	}

	// u2 = ((1081+1)^22 * 17^1081 * 674^-881) mod 1081^2
	expectedU2 := big.NewInt(227035)
	actualU2 := zkp.u2Verification(w, params)
	if !reflect.DeepEqual(expectedU2, actualU2) {
		t.Errorf(
			"Unexpected u2\nActual: %v\nExpected: %v",
			actualU2,
			expectedU2,
		)
	}

	// u3 = (11^22 * 27^63 * 191^-881) mod 253
	expectedU3 := big.NewInt(44)
	actualU3 := zkp.u3Verification(params)
	if !reflect.DeepEqual(expectedU3, actualU3) {
		t.Errorf(
			"Unexpected u3\nActual: %v\nExpected: %v",
			actualU3,
			expectedU3,
		)
	}
}
