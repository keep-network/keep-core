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

func TestLocalSignerGenerateDsaKeyShare(t *testing.T) {
	group, err := createNewLocalGroup()
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
	group, dsaKey, err := initializeNewLocalGroupWithFullKey()
	if err != nil {
		t.Fatal(err)
	}

	// We have ThresholdDsaKey with E(secretKey) and public key where
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
	group, commitmentMessages, revealMessages, err := initializeNewLocalGroupWithKeyShares()
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
	group, commitmentMessages, revealMessages, err := initializeNewLocalGroupWithKeyShares()
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
	group, commitmentMessages, revealMessages, err := initializeNewLocalGroupWithKeyShares()
	if err != nil {
		t.Fatal(err)
	}

	invalidCommitment, _, err := commitment.Generate(
		big.NewInt(1).Bytes(),
	)
	if err != nil {
		t.Fatal(err)
	}
	commitmentMessages[len(commitmentMessages)-1].publicKeyShareCommitment =
		invalidCommitment

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
	group, commitmentMessages, revealMessages, err := initializeNewLocalGroupWithKeyShares()
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

// Test group public parameters used to construct and validate groups for tests
var publicParameters = &PublicParameters{
	groupSize: 10,
	threshold: 6,
	curve:     secp256k1.S256(),
}

// createNewCoreGroup generates a group of `signerCore`s backed by a threshold
// Paillier key and ZKP public parameters built from the generated Paillier key.
// This approach works in an oracle mode - one party is responsible for
// generating Paillier keys and distributing them and should be used only for
// testing.
func createNewCoreGroup() ([]*signerCore, error) {
	paillierKeyGen := paillier.GetThresholdKeyGenerator(
		paillierModulusBitLength,
		publicParameters.groupSize,
		publicParameters.threshold,
		rand.Reader,
	)

	paillierKeys, err := paillierKeyGen.Generate()
	if err != nil {
		return nil, fmt.Errorf(
			"could not generate threshold Paillier keys [%v]", err,
		)
	}

	zkpParameters, err := zkp.GeneratePublicParameters(
		paillierKeys[0].N,
		publicParameters.curve,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"could not generate public ZKP parameters [%v]", err,
		)
	}

	members := make([]*signerCore, len(paillierKeys))
	for i := 0; i < len(members); i++ {
		members[i] = &signerCore{
			ID:              generateMemberID(),
			paillierKey:     paillierKeys[i],
			groupParameters: publicParameters,
			zkpParameters:   zkpParameters,
		}
	}

	return members, nil
}

// createNewLocalGroup creates a new group of `LocalSigner`s that did not
// started initialization process yet.
func createNewLocalGroup() ([]*LocalSigner, error) {
	coreGroup, err := createNewCoreGroup()
	if err != nil {
		return nil, err
	}

	localSigners := make([]*LocalSigner, len(coreGroup))
	for i := 0; i < len(localSigners); i++ {
		localSigners[i] = &LocalSigner{
			signerCore: *coreGroup[i],
		}
	}

	return localSigners, nil
}

// initializeNewLocalGroupWithKeyShares creates and initializes a new group of
// `LocalSigner`s. It simulates a real initialization process by first calling
// `InitializeDsaKeyShares` and then `RevealDsaKeyShares`. Messages produced by
// those functions are returned along with all `LocalSigner`s created.
// It's responsibility of code calling this function to execute
// `CombineDsaKeyShares`, in order to produce signers with a fully initialized
// threshold ECDSA key, if needed.
func initializeNewLocalGroupWithKeyShares() (
	[]*LocalSigner,
	[]*PublicKeyShareCommitmentMessage,
	[]*KeyShareRevealMessage,
	error,
) {
	group, err := createNewLocalGroup()
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

// initializeNewLocalGroupWithFullKey creates and initializes a new group of
// `LocalSigner`s simulating a real initialization process just like the
// `initializeNewLocalGroupWithKeyShares` except that it also calls
// `ConbineDsaKeyShares` in order to produce a full `ThresholdDsaKey`.
func initializeNewLocalGroupWithFullKey() (
	[]*LocalSigner, *ThresholdDsaKey, error,
) {
	group, commitmentMessages, revealMessages, err :=
		initializeNewLocalGroupWithKeyShares()
	if err != nil {
		return nil, nil, err
	}

	// Combine all PublicKeyShareCommitmentMessages and KeyShareRevealMessages
	// from signers in order to create a ThresholdDsaKey.
	dsaKey, err := group[0].CombineDsaKeyShares(
		commitmentMessages,
		revealMessages,
	)
	if err != nil {
		return nil, nil, err
	}

	return group, dsaKey, nil
}
