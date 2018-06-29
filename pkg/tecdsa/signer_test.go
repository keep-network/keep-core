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
	if curveCardinality.Cmp(dsaKeyShare.xi) != 1 {
		t.Errorf("xi DSA key share must be less than Curve's cardinality")
	}

	if !publicParameters.curve.IsOnCurve(dsaKeyShare.yi.x, dsaKeyShare.yi.y) {
		t.Errorf("yi DSA key share must be a point on Curve")
	}
}

func TestInitializeAndCombineDsaKey(t *testing.T) {
	group, err := newGroup(publicParameters)
	if err != nil {
		t.Fatal(err)
	}

	// Let each signer initialize the DSA key share and create InitMessage.
	// Each signer picks randomly xi from Z_q and computes yi = g^xi.
	// xi and yi are shares of a DSA key (x, y).
	//
	// E(xi) and yi are published by signer in a broadcast InitMessage.
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

	// Now we have ThresholdDsaKey with E(x) and y where E(x) is a threshold
	// sharing of x.
	//
	// We may check the correctness of E(x) and y:
	// 1. y should be a point on the elliptic curve
	// 2. x can be decomposed by a group of Signers (just for test)
	// 3. y = g^x should hold, according to how DSA key is constructed

	// 1. Check if y is a point on curve
	if !publicParameters.curve.IsOnCurve(dsaKey.y.x, dsaKey.y.y) {
		t.Fatal("ThresholdDsaKey.y must be a point on Curve")
	}

	// 2. Deconstruct x from E(x)
	xShares := make([]*paillier.PartialDecryption, publicParameters.groupSize)
	for i, signer := range group {
		xShares[i] = signer.paillerKey.Decrypt(dsaKey.x.C)
	}
	x, err := group[0].paillerKey.CombinePartialDecryptions(xShares)
	if err != nil {
		t.Fatal(err)
	}

	// Since xi is from Z_q when partial keys are generated, after we add all
	// xi shares we may produce number greater than q and exceed curve's
	// cardinality. specs256k1 can't handle scalars > 256 bits, so we need to
	// mod N here to stay in the curve's field.
	x = new(big.Int).Mod(x, publicParameters.curve.Params().N)

	// 3. Having x, we can evaluate y from y = g^x and compare with the actual
	//    value stored in ThresholdDsaKey.
	yx, yy := publicParameters.curve.ScalarBaseMult(x.Bytes())
	if !reflect.DeepEqual(yx, dsaKey.y.x) {
		t.Errorf("Unexpected y.x decoded\nActual %v\nExpected %v", yx, dsaKey.y.x)
	}
	if !reflect.DeepEqual(yy, dsaKey.y.y) {
		t.Errorf("Unexpected y.y decoded\nActual %v\nExpected %v", yy, dsaKey.y.y)
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
