package tecdsa

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/keep-network/keep-core/pkg/tecdsa/commitment"

	"github.com/keep-network/keep-core/pkg/tecdsa/curve"
	"github.com/keep-network/keep-core/pkg/tecdsa/zkp"
	"github.com/keep-network/paillier"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
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

	if !publicParameters.curve.IsOnCurve(dsaKeyShare.publicKeyShare.X, dsaKeyShare.publicKeyShare.Y) {
		t.Errorf("DSA public key share must be a point on Curve")
	}
}

func TestInitializeAndCombineDsaKey(t *testing.T) {
	group, commitmentMessages, revealMessages, err := initializeNewGroup()
	if err != nil {
		t.Fatal(err)
	}

	// Combine all PublicKeyShareCommitmentMessages and KeyShareRevealMessages
	// from signers in order to create a ThresholdDsaKey.
	dsaKey, err := group[0].CombineDsaKeyShares(
		commitmentMessages,
		revealMessages,
	)
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
	if !publicParameters.curve.IsOnCurve(dsaKey.publicKey.X, dsaKey.publicKey.Y) {
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

	if !reflect.DeepEqual(publicKeyX, dsaKey.publicKey.X) {
		t.Errorf(
			"Unexpected publicKey.x decoded\nActual %v\nExpected %v",
			publicKeyX,
			dsaKey.publicKey.X,
		)
	}
	if !reflect.DeepEqual(publicKeyY, dsaKey.publicKey.Y) {
		t.Errorf(
			"Unexpected publicKey.y decoded\nActual %v\nExpected %v",
			publicKeyY,
			dsaKey.publicKey.Y,
		)
	}
}

func TestCombineWithNotEnoughCommitMessages(t *testing.T) {
	group, commitmentMessages, revealMessages, err := initializeNewGroup()
	if err != nil {
		t.Fatal(err)
	}

	expectedError := fmt.Errorf(
		"commitments required from all group members; got 1, expected 10",
	)

	_, err = group[0].CombineDsaKeyShares(
		[]*PublicKeyShareCommitmentMessage{commitmentMessages[0]},
		revealMessages,
	)
	if err == nil {
		t.Fatal("Error was expected")
	}
	if !reflect.DeepEqual(expectedError, err) {
		t.Errorf("Unexpected error\nActual %v\nExpected %v", expectedError, err)
	}
}

func TestCombineWithNotEnoughRevealMessages(t *testing.T) {
	group, commitmentMessages, revealMessages, err := initializeNewGroup()
	if err != nil {
		t.Fatal(err)
	}

	expectedError := fmt.Errorf(
		"all group members should reveal shares; Got 1, expected 10",
	)

	_, err = group[0].CombineDsaKeyShares(
		commitmentMessages,
		[]*KeyShareRevealMessage{revealMessages[0]},
	)
	if err == nil {
		t.Fatal("Error was expected")
	}
	if !reflect.DeepEqual(expectedError, err) {
		t.Errorf("Unexpected error\nActual %v\nExpected %v", expectedError, err)
	}
}

func TestCombineWithInvalidCommitment(t *testing.T) {
	group, commitmentMessages, revealMessages, err := initializeNewGroup()
	if err != nil {
		t.Fatal(err)
	}

	invalidCommitment, _, err := commitment.GenerateCommitment(
		big.NewInt(1).Bytes(),
	)
	if err != nil {
		t.Fatal(err)
	}
	commitmentMessages[len(commitmentMessages)-1].commitment = invalidCommitment

	expectedError := fmt.Errorf("KeyShareRevealMessage rejected")

	_, err = group[0].CombineDsaKeyShares(commitmentMessages, revealMessages)
	if err == nil {
		t.Fatal("Error was expected")
	}
	if !reflect.DeepEqual(expectedError, err) {
		t.Errorf("Unexpected error\nActual %v\nExpected %v", expectedError, err)
	}
}

func TestCombineWithInvalidZKP(t *testing.T) {
	group, commitmentMessages, revealMessages, err := initializeNewGroup()
	if err != nil {
		t.Fatal(err)
	}

	// Let's modify one of reveal message ZKPs to make it fail
	invalidProof, err := zkp.CommitDsaPaillierKeyRange(
		big.NewInt(1),
		&curve.Point{X: big.NewInt(1), Y: big.NewInt(2)},
		&paillier.Cypher{C: big.NewInt(3)},
		big.NewInt(1),
		group[0].zkpParameters,
		rand.Reader,
	)
	if err != nil {
		t.Fatal(err)
	}
	revealMessages[len(revealMessages)-1].secretKeyProof = invalidProof

	expectedError := fmt.Errorf("KeyShareRevealMessage rejected")

	_, err = group[0].CombineDsaKeyShares(commitmentMessages, revealMessages)
	if err == nil {
		t.Fatal("Error was expected")
	}
	if !reflect.DeepEqual(expectedError, err) {
		t.Errorf("Unexpected error\nActual %v\nExpected %v", expectedError, err)
	}
}

func initializeNewGroup() (
	[]*LocalSigner,
	[]*PublicKeyShareCommitmentMessage,
	[]*KeyShareRevealMessage,
	error,
) {
	group, err := newGroup(publicParameters)
	if err != nil {
		return nil, nil, nil, err
	}

	// Let each signer initialize the DSA key share and create
	// PublicKeyShareCommitmentMessage. Each signer picks randomly
	// secretKeyShare from Z_q and computes publicKeyShare = g^secretKeyShare.
	// Generated key shares are saved internally by each Signer. Each Signer
	// generates commitment for the public key share. Commitment is broadcasted
	// in the PublicKeyShareCommitmentMessage.
	publicKeyCommitmentMessages := make(
		[]*PublicKeyShareCommitmentMessage,
		publicParameters.groupSize,
	)
	for i, signer := range group {
		publicKeyCommitmentMessages[i], err = signer.InitializeDsaKeyShares()
		if err != nil {
			return nil, nil, nil, err
		}
	}

	// In the next phase, each Signer publishes KeyShareRevealMessage with:
	// - E(secretKeyShare)
	// - decommitment for the public key share
	// - ZKP for the secret key
	//
	// E is an additively homomorphic encryption scheme. For our implementation
	// we use Paillier.
	keyShareRevealMessages := make(
		[]*KeyShareRevealMessage,
		publicParameters.groupSize,
	)
	for i, signer := range group {
		keyShareRevealMessages[i], err = signer.RevealDsaKeyShares()
		if err != nil {
			return nil, nil, nil, err
		}
	}

	return group, publicKeyCommitmentMessages, keyShareRevealMessages, nil
}
