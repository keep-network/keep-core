package tecdsa

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/keep-network/paillier"
)

var publicParameters = &PublicParameters{
	groupSize: 10,
	threshold: 6,
	curve:     secp256k1.S256(),
}

func TestLocalSignerGenerateDsaKeyShare(t *testing.T) {
	group, err := newGroup(publicParameters)
	if err != nil {
		t.Fatal(err)
	}

	signer := group[0]

	dsaKeyShare, err := signer.generateDsaKeyShare()
	if err != nil {
		t.Fatal(err)
	}

	curveCardinality := publicParameters.curve.Params().N
	if curveCardinality.Cmp(dsaKeyShare.secretKeyShare) != 1 {
		t.Errorf("DSA secret key share must be less than Curve's cardinality")
	}

	if !publicParameters.curve.IsOnCurve(dsaKeyShare.publicKeyShare.x, dsaKeyShare.publicKeyShare.y) {
		t.Errorf("DSA public key share must be a point on Curve")
	}
}

func TestInitializeAndCombineDsaKey(t *testing.T) {
	group, err := newGroup(publicParameters)
	if err != nil {
		t.Fatal(err)
	}

	// Let each signer initialize the DSA key share and create InitMessage.
	// Each signer picks randomly secretKeyShare from Z_q and computes
	// publicKeyShare = g^secretKeyShare.
	//
	// E(secretKeyShare) and publicKeyShare are published by signer in a
	// broadcast InitMessage.
	// E is an additively homomorphic encryption scheme. For our implementation
	// we use Paillier.
	initMessages := make([]*InitMessage, publicParameters.groupSize)
	for i, signer := range group {
		initMessages[i], err = signer.InitializeDsaKeyGen()
		if err != nil {
			t.Fatal(err)
		}
	}

	// Combine all InitMessages from signers in order to create ThresholdDsaKey.
	dsaKey, err := group[0].CombineDsaKeyShares(initMessages)
	if err != nil {
		t.Fatal(err)
	}

	// Now we have ThresholdDsaKey with E(secretKey) and public key where
	// E(secretKey) is a threshold sharing of secretKey.
	//
	// We may check the correctness of E(secretKey) and publicKey:
	// 1. publicKey should be a point on the elliptic curve
	// 2. secretKey can be decrypted by a group of Signers (just for test)
	// 3. publicKey = g^secretKey should hold, according to how DSA key is
	//    constructed

	// 1. Check if publicKey is a point on curve
	if !publicParameters.curve.IsOnCurve(dsaKey.publicKey.x, dsaKey.publicKey.y) {
		t.Fatal("ThresholdDsaKey.y must be a point on Curve")
	}

	// 2. Decrypt secretKey from E(secretKey)
	xShares := make([]*paillier.PartialDecryption, publicParameters.groupSize)
	for i, signer := range group {
		xShares[i] = signer.paillierKey.Decrypt(dsaKey.secretKey.C)
	}
	secretKey, err := group[0].paillierKey.CombinePartialDecryptions(xShares)
	if err != nil {
		t.Fatal(err)
	}

	// Since xi is from Z_q when partial keys are generated, after we add all
	// xi shares we may produce number greater than q and exceed curve's
	// cardinality. specs256k1 can't handle scalars > 256 bits, so we need to
	// mod N here to stay in the curve's field.
	secretKey = new(big.Int).Mod(secretKey, publicParameters.curve.Params().N)

	// 3. Having secretKey, we can evaluate publicKey from
	//    publicKey = g^secretKey and compare with the actual
	//    value stored in ThresholdDsaKey.
	publicKeyX, publicKeyY := publicParameters.curve.ScalarBaseMult(secretKey.Bytes())

	if !reflect.DeepEqual(publicKeyX, dsaKey.publicKey.x) {
		t.Errorf(
			"Unexpected publicKey.x decoded\nActual %v\nExpected %v",
			publicKeyX,
			dsaKey.publicKey.x,
		)
	}
	if !reflect.DeepEqual(publicKeyY, dsaKey.publicKey.y) {
		t.Errorf(
			"Unexpected publicKey.y decoded\nActual %v\nExpected %v",
			publicKeyY,
			dsaKey.publicKey.y,
		)
	}
}

func TestCombineNotEnoughInitMessages(t *testing.T) {
	group, err := newGroup(publicParameters)
	if err != nil {
		t.Fatal(err)
	}

	message, err := group[1].InitializeDsaKeyGen()
	if err != nil {
		t.Fatal(err)
	}

	expectedError := fmt.Errorf(
		"InitMessages required from all group members; Got 1, expected 10",
	)

	shares := []*InitMessage{message}
	_, err = group[0].CombineDsaKeyShares(shares)
	if err == nil {
		t.Fatal("Error was expected")
	}
	if !reflect.DeepEqual(expectedError, err) {
		t.Errorf("Unexpected error\nActual %v\nExpected %v", expectedError, err)
	}

}
