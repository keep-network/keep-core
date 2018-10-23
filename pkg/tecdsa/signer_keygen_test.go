package tecdsa

import (
	"bytes"
	"crypto/rand"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"math/big"
	"reflect"
	"testing"

	"github.com/keep-network/keep-core/pkg/tecdsa/commitment"

	"github.com/keep-network/keep-core/pkg/tecdsa/curve"
	"github.com/keep-network/keep-core/pkg/tecdsa/zkp"
	"github.com/keep-network/paillier"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

func TestLocalSignerGenerateEcdsaKeyShare(t *testing.T) {
	group, parameters, err := createNewLocalGroup()
	if err != nil {
		t.Fatal(err)
	}

	signer := group[0]

	ecdsaKeyShare, err := signer.generateEcdsaKeyShare()
	if err != nil {
		t.Fatal(err)
	}

	if parameters.curveCardinality().Cmp(ecdsaKeyShare.secretKeyShare) != 1 {
		t.Errorf("DSA secret key share must be less than Curve's cardinality")
	}

	if !parameters.Curve.IsOnCurve(
		ecdsaKeyShare.publicKeyShare.X,
		ecdsaKeyShare.publicKeyShare.Y,
	) {
		t.Errorf("DSA public key share must be a point on Curve")
	}
}

func TestInitializeAndCombineDsaKey(t *testing.T) {
	group, parameters, dsaKey, err := initializeNewLocalGroupWithFullKey()
	if err != nil {
		t.Fatal(err)
	}

	// We have ThresholdEcdsaKey with E(secretKey) and public key where
	// E(secretKey) is a threshold sharing of secretKey.
	//
	// We may check the correctness of E(secretKey) and publicKey:
	// 1. publicKey should be a point on the elliptic curve
	// 2. secretKey can be decrypted by a group of Signers (just for test)
	// 3. publicKey = g^secretKey should hold, according to how DSA key is
	//    constructed

	// 1. Check if publicKey is a point on curve
	if !parameters.Curve.IsOnCurve(dsaKey.PublicKey.X, dsaKey.PublicKey.Y) {
		t.Fatal("ThresholdEcdsaKey.y must be a point on Curve")
	}

	// 2. Decrypt secretKey from E(secretKey)
	xShares := make([]*paillier.PartialDecryption, len(group))
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
	secretKey = new(big.Int).Mod(secretKey, parameters.Curve.Params().N)

	// 3. Having secretKey, we can evaluate publicKey from
	//    publicKey = g^secretKey and compare with the actual
	//    value stored in ThresholdEcdsaKey.
	publicKeyX, publicKeyY := parameters.Curve.ScalarBaseMult(secretKey.Bytes())

	if !reflect.DeepEqual(publicKeyX, dsaKey.PublicKey.X) {
		t.Errorf(
			"Unexpected publicKey.x decoded\nActual %v\nExpected %v",
			publicKeyX,
			dsaKey.PublicKey.X,
		)
	}
	if !reflect.DeepEqual(publicKeyY, dsaKey.PublicKey.Y) {
		t.Errorf(
			"Unexpected publicKey.y decoded\nActual %v\nExpected %v",
			publicKeyY,
			dsaKey.PublicKey.Y,
		)
	}
}

func TestCombineWithNotEnoughCommitMessages(t *testing.T) {
	group, _, commitmentMessages, revealMessages, err :=
		initializeNewLocalGroupWithKeyShares()
	if err != nil {
		t.Fatal(err)
	}

	expectedError := fmt.Errorf(
		"commitments required from all group peer members; got 1, expected 9",
	)

	receiver := group[0]
	receiverCommitmentMessages := publicKeyShareCommitmentMessagesForReceiver(
		commitmentMessages, receiver.ID,
	)
	receiverRevealMessages := publicKeyShareRevealMessagesForReceiver(
		revealMessages, receiver.ID,
	)
	_, err = receiver.CombineEcdsaKeyShares(
		[]*PublicEcdsaKeyShareCommitmentMessage{receiverCommitmentMessages[0]},
		receiverRevealMessages,
	)
	if err == nil {
		t.Fatal("Error was expected")
	}
	if !reflect.DeepEqual(expectedError, err) {
		t.Errorf("Unexpected error\nActual %v\nExpected %v", expectedError, err)
	}
}

func TestCombineWithNotEnoughRevealMessages(t *testing.T) {
	group, _, commitmentMessages, revealMessages, err :=
		initializeNewLocalGroupWithKeyShares()
	if err != nil {
		t.Fatal(err)
	}

	expectedError := fmt.Errorf(
		"all group peer members should reveal shares; Got 1, expected 9",
	)

	receiver := group[0]
	receiverCommitmentMessages := publicKeyShareCommitmentMessagesForReceiver(
		commitmentMessages, receiver.ID,
	)
	receiverRevealMessages := publicKeyShareRevealMessagesForReceiver(
		revealMessages, receiver.ID,
	)
	_, err = receiver.CombineEcdsaKeyShares(
		receiverCommitmentMessages,
		[]*KeyShareRevealMessage{receiverRevealMessages[0]},
	)
	if err == nil {
		t.Fatal("Error was expected")
	}
	if !reflect.DeepEqual(expectedError, err) {
		t.Errorf("Unexpected error\nActual %v\nExpected %v", expectedError, err)
	}
}

func TestCombineWithInvalidCommitment(t *testing.T) {
	group, _, commitmentMessages, revealMessages, err :=
		initializeNewLocalGroupWithKeyShares()
	if err != nil {
		t.Fatal(err)
	}
	receiver := group[0]

	invalidCommitment, _, err := commitment.Generate(
		receiver.protocolParameters[receiver.ID].commitmentMasterPublicKey,
		big.NewInt(1).Bytes(),
	)
	if err != nil {
		t.Fatal(err)
	}

	receiverCommitmentMessages := publicKeyShareCommitmentMessagesForReceiver(
		commitmentMessages, receiver.ID,
	)
	receiverRevealMessages := publicKeyShareRevealMessagesForReceiver(
		revealMessages, receiver.ID,
	)
	receiverCommitmentMessages[1].publicKeyShareCommitment =
		invalidCommitment

	expectedError := fmt.Errorf("KeyShareRevealMessage rejected")

	_, err = receiver.CombineEcdsaKeyShares(
		receiverCommitmentMessages,
		receiverRevealMessages,
	)
	if err == nil {
		t.Fatal("Error was expected")
	}
	if !reflect.DeepEqual(expectedError, err) {
		t.Errorf("Unexpected error\nActual %v\nExpected %v", expectedError, err)
	}
}

func TestCombineWithInvalidZKP(t *testing.T) {
	group, _, commitmentMessages, revealMessages, err :=
		initializeNewLocalGroupWithKeyShares()
	if err != nil {
		t.Fatal(err)
	}

	// Let's modify one of reveal message ZKPs to make it fail
	invalidProof, err := zkp.CommitEcdsaPaillierKeyRange(
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

	receiver := group[0]
	receiverCommitmentMessages := publicKeyShareCommitmentMessagesForReceiver(
		commitmentMessages, receiver.ID,
	)
	receiverRevealMessages := publicKeyShareRevealMessagesForReceiver(
		revealMessages, receiver.ID,
	)

	receiverRevealMessages[len(receiverRevealMessages)-1].secretKeyProof =
		invalidProof

	expectedError := fmt.Errorf("KeyShareRevealMessage rejected")

	_, err = group[0].CombineEcdsaKeyShares(
		receiverCommitmentMessages, receiverRevealMessages,
	)
	if err == nil {
		t.Fatal("Error was expected")
	}
	if !reflect.DeepEqual(expectedError, err) {
		t.Errorf("Unexpected error\nActual %v\nExpected %v", expectedError, err)
	}
}

// readTestParameters returns test `PublicParameters`, threshold Paillier
// key generated for those parameters and ZKP `PublicParameters` generated for
// the key.
//
// Since Paillier key initialization is quite expensive, for secp256k1 curve
// Paillier key must be at least 2048-bit long and we don't really want to test
// the initialization process here, we just load the existing key from file.
// Moreover, we create ZKP parameters from predefined 1024-bit safe prime
// numbers instead of generating them for the same reason as for Paillier key.
//
// Please bear in mind this is not a correct approach for all T-ECDSA tests!
func readTestParameters() (
	[]paillier.ThresholdPrivateKey,
	*PublicParameters,
	*zkp.PublicParameters,
	*signerGroup,
	error,
) {
	publicParameters := &PublicParameters{
		Curve:                secp256k1.S256(),
		PaillierKeyBitLength: 2048,
	}

	signerGroup := &signerGroup{
		InitialGroupSize: 10,
		Threshold:        6,
		signerIDs:        make([]string, 0),
	}

	var paillierKey []paillier.ThresholdPrivateKey
	data, err := ioutil.ReadFile("paillier.test.key")
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("Could not read test Paillier key [%v]", err)
	}
	buffer := bytes.NewBuffer(data)

	decoder := gob.NewDecoder(buffer)
	err = decoder.Decode(&paillierKey)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("Could not read test Paillier key [%v]", err)
	}

	pTilde, _ := new(big.Int).SetString(
		"T7J6D6kn5JGzrJpskwFuabU33cBMKgkfjucmEuSvCWH4PmgcghcdSBV9u8oPWikoKokF"+
			"XVVr8NriqFizRqbRiG7R0OABnxHl1lL2CtS4pv3Raoc8Ubv6uGzbat2dx0kgGkve"+
			"LfIQOVmLrFxYOvFXKHEKal9OYlZTdIcRe5vA9XfR", 62,
	)
	qTilde, _ := new(big.Int).SetString(
		"UOzoNzDzf4dFKuTyVfEXGnbaiLXTDBtWoWLR83wl34AgEOteIhDtVdngNPbcCKP1A0H9"+
			"CBRUE081rsad6ftTyiVeDhEHIToSf3LRenAJAxic4kNrzYtyKN1yFtzMKPH6ndcT"+
			"8NPTI7gEOt789ZAf3queB54mbZjWWsm1sV2plmhR", 62,
	)

	zkpParameters, err := zkp.GeneratePublicParametersFromSafePrimes(
		paillierKey[0].N, pTilde, qTilde, publicParameters.Curve,
	)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf(
			"could not generate public ZKP parameters [%v]", err,
		)
	}

	return paillierKey, publicParameters, zkpParameters, signerGroup, nil
}

// createNewLocalGroup creates a new group of `LocalSigner`s that did not
// started initialization process yet.
func createNewLocalGroup() ([]*LocalSigner, *PublicParameters, error) {
	paillierKeys, signerParameters, zkpParameters, signerGroup, err := readTestParameters()
	if err != nil {
		return nil, nil, err
	}

	localSigners := make([]*LocalSigner, len(paillierKeys))
	for i := 0; i < len(localSigners); i++ {
		signer := NewLocalSigner(
			&paillierKeys[i], signerParameters, zkpParameters, signerGroup,
		)

		signerGroup.signerIDs = append(signerGroup.signerIDs, signer.ID)
		localSigners[i] = signer
	}

	return localSigners, signerParameters, nil
}

// initializeNewLocalGroupWithKeyShares creates and initializes a new group of
// `LocalSigner`s.
// It simulates a real initialization process by first generating master public key
// for multi-trapdoor commitment scbheme and then calling `InitializeEcdsaKeyShares`
//  and then `RevealEcdsaKeyShares`. Messages produced by those functions are returned
// along with all `LocalSigner`s created.
// It's responsibility of code calling this function to execute
// `CombineEcdsaKeyShares`, in order to produce signers with a fully initialized
// threshold ECDSA key, if needed.
func initializeNewLocalGroupWithKeyShares() (
	[]*LocalSigner,
	*PublicParameters,
	[]*PublicEcdsaKeyShareCommitmentMessage,
	[]*KeyShareRevealMessage,
	error,
) {
	group, parameters, err := createNewLocalGroup()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	err = setupGroup(group)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// Let each signer initialize the DSA key share and create
	// PublicKeyShareCommitmentMessage. Each signer picks randomly
	// secretKeyShare from Z_q and computes publicKeyShare = g^secretKeyShare.
	// Generated key shares are saved internally by each Signer. Each Signer
	// generates commitment for the public key share. Commitment is broadcasted
	// in the PublicKeyShareCommitmentMessage.
	publicKeyCommitmentMessages := make([]*PublicEcdsaKeyShareCommitmentMessage, 0)
	for _, signer := range group {
		messages, err := signer.InitializeEcdsaKeyShares()
		if err != nil {
			return nil, nil, nil, nil, err
		}
		publicKeyCommitmentMessages = append(
			publicKeyCommitmentMessages, messages...,
		)
	}

	// In the next phase, each Signer publishes KeyShareRevealMessage with:
	// - E(secretKeyShare)
	// - decommitment for the public key share
	// - ZKP for the secret key
	//
	// E is an additively homomorphic encryption scheme. For our implementation
	// we use Paillier.
	keyShareRevealMessages := make([]*KeyShareRevealMessage, 0)
	for _, signer := range group {
		messages, err := signer.RevealEcdsaKeyShares()
		if err != nil {
			return nil, nil, nil, nil, err
		}
		keyShareRevealMessages = append(keyShareRevealMessages, messages...)
	}

	return group, parameters, publicKeyCommitmentMessages, keyShareRevealMessages, nil
}

// initializeNewLocalGroupWithFullKey creates and initializes a new group of
// `LocalSigner`s simulating a real initialization process just like the
// `initializeNewLocalGroupWithKeyShares` except that it also calls
// `CombineEcdsaKeyShares` in order to produce a full `ThresholdEcdsaKey`.
func initializeNewLocalGroupWithFullKey() (
	[]*LocalSigner, *PublicParameters, *ThresholdEcdsaKey, error,
) {
	group, parameters, commitmentMessages, revealMessages, err :=
		initializeNewLocalGroupWithKeyShares()
	if err != nil {
		return nil, nil, nil, err
	}

	// Combine all PublicKeyShareCommitmentMessages and KeyShareRevealMessages
	// from signers in order to create a ThresholdDsaKey.
	receiver := group[0]
	ecdsaKey, err := receiver.CombineEcdsaKeyShares(
		publicKeyShareCommitmentMessagesForReceiver(commitmentMessages, receiver.ID),
		publicKeyShareRevealMessagesForReceiver(revealMessages, receiver.ID),
	)
	if err != nil {
		return nil, nil, nil, err
	}

	return group, parameters, ecdsaKey, nil
}

func publicKeyShareCommitmentMessagesForReceiver(
	messages []*PublicEcdsaKeyShareCommitmentMessage,
	receiverID string,
) []*PublicEcdsaKeyShareCommitmentMessage {
	filtered := make([]*PublicEcdsaKeyShareCommitmentMessage, 0)
	for _, message := range messages {
		if message.receiverID == receiverID {
			filtered = append(filtered, message)
		}
	}
	return filtered
}

func publicKeyShareRevealMessagesForReceiver(
	messages []*KeyShareRevealMessage,
	receiverID string,
) []*KeyShareRevealMessage {
	filtered := make([]*KeyShareRevealMessage, 0)
	for _, message := range messages {
		if message.receiverID == receiverID {
			filtered = append(filtered, message)
		}
	}
	return filtered
}
