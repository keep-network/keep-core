package tecdsa

// TODO: rename to signer_smoke_test.go

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"math/big"

	"testing"

	"github.com/btcsuite/btcd/btcec"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/keep-network/keep-core/pkg/tecdsa/curve"
	"github.com/keep-network/keep-core/pkg/tecdsa/zkp"
	"github.com/keep-network/paillier"
	"github.com/keep-network/keep-core/pkg/internal/byteutils"
)

func TestFullInitAndSignPath(t *testing.T) {
	messageHash := make([]byte, 32)

	_, err := rand.Read(messageHash)
	if err != nil {
		t.Fatal(err)
	}

	localSigners, parameters, err := generateNewLocalGroup()
	if err != nil {
		t.Fatal(err)
	}

	masterPublicKeyShareMessagesKeyGeneration := make([]*MasterPublicKeyShareMessage, len(localSigners))
	publicKeyCommitmentMessages := make([]*PublicKeyShareCommitmentMessage, len(localSigners))
	keyShareRevealMessages := make([]*KeyShareRevealMessage, len(localSigners))

	masterPublicKeyShareMessagesSigning := make([]*MasterPublicKeyShareMessage, len(localSigners))

	round1Messages := make([]*SignRound1Message, len(localSigners))
	round2Messages := make([]*SignRound2Message, len(localSigners))
	round3Messages := make([]*SignRound3Message, len(localSigners))
	round4Messages := make([]*SignRound4Message, len(localSigners))
	round5Messages := make([]*SignRound5Message, len(localSigners))
	round6Messages := make([]*SignRound6Message, len(localSigners))

	round1Signers := make([]*Round1Signer, len(localSigners))
	round2Signers := make([]*Round2Signer, len(localSigners))
	round3Signers := make([]*Round3Signer, len(localSigners))
	round4Signers := make([]*Round4Signer, len(localSigners))
	round5Signers := make([]*Round5Signer, len(localSigners))

	//
	// Initialize master public key for multi-trapdoor commitment scheme for key
	// generation process
	//
	for i, signer := range localSigners {
		masterPublicKeyShareMessagesKeyGeneration[i], err = signer.GenerateMasterPublicKeyShare()
		if err != nil {
			t.Fatal(err)
		}
	}

	masterPublicKeyKeyGeneration, err := localSigners[0].CombineMasterPublicKeyShares(masterPublicKeyShareMessagesKeyGeneration)
	if err != nil {
		t.Fatal(err)
	}

	for _, signer := range localSigners {
		signer.commitmentMasterPublicKey = masterPublicKeyKeyGeneration
		if err != nil {
			t.Fatal(err)
		}
	}

	//
	// Execute the 1st key-gen round
	//
	for i, signer := range localSigners {
		publicKeyCommitmentMessages[i], err = signer.InitializeDsaKeyShares()
		if err != nil {
			t.Fatal(err)
		}
	}

	//
	// Execute the 2nd key-gen round
	//
	for i, signer := range localSigners {
		keyShareRevealMessages[i], err = signer.RevealDsaKeyShares()
		if err != nil {
			t.Fatal(err)
		}
	}

	dsaKey, err := localSigners[0].CombineDsaKeyShares(
		publicKeyCommitmentMessages,
		keyShareRevealMessages,
	)

	signers := make([]*Signer, len(localSigners))
	for i, localSigner := range localSigners {
		signers[i] = localSigner.WithDsaKey(dsaKey)
	}

	//
	// Initialize master public key for multi-trapdoor commitment scheme for
	// signing process
	//
	for i, signer := range signers {
		masterPublicKeyShareMessagesSigning[i], err = signer.GenerateMasterPublicKeyShare()
		if err != nil {
			t.Fatal(err)
		}
	}

	masterpublicKeySigning, err := signers[0].CombineMasterPublicKeyShares(masterPublicKeyShareMessagesSigning)
	if err != nil {
		t.Fatal(err)
	}

	for _, signer := range signers {
		signer.commitmentMasterPublicKey = masterpublicKeySigning
		if err != nil {
			t.Fatal(err)
		}
	}

	//
	// Execute the 1st signing round
	//
	for i, signer := range signers {
		round1Signers[i], round1Messages[i], err = signer.SignRound1()
		if err != nil {
			t.Fatal(err)
		}
	}

	//
	// Execute the 2nd signing round
	//
	for i, signer := range round1Signers {
		round2Signers[i], round2Messages[i], err = signer.SignRound2()
		if err != nil {
			t.Fatal(err)
		}
	}

	secretKeyRandomFactor, secretKeyMultiple, err :=
		round2Signers[0].CombineRound2Messages(
			round1Messages, round2Messages,
		)
	if err != nil {
		t.Fatal(err)
	}

	//
	// Execute the 3rd signing round
	//
	for i, signer := range round2Signers {
		round3Signers[i], round3Messages[i], err = signer.SignRound3(
			secretKeyRandomFactor, secretKeyMultiple,
		)
		if err != nil {
			t.Fatal(err)
		}
	}

	//
	// Execute the 4th signing round
	//
	for i, signer := range round3Signers {
		round4Signers[i], round4Messages[i], err = signer.SignRound4()
		if err != nil {
			t.Fatal(err)
		}
	}

	signatureUnmask, signatureRandomMultiplePublic, err :=
		round4Signers[0].CombineRound4Messages(
			round3Messages, round4Messages,
		)
	if err != nil {
		t.Fatal(err)
	}

	//
	// Execute the 5th signing round
	//
	for i, signer := range round4Signers {
		round5Signers[i], round5Messages[i], err = signer.SignRound5(
			signatureUnmask, signatureRandomMultiplePublic,
		)
		if err != nil {
			t.Fatal(err)
		}
	}

	signatureUnmaskDecrypted, err := round5Signers[0].CombineRound5Messages(
		round5Messages,
	)
	if err != nil {
		t.Fatal(err)
	}

	//
	// Execute the 6th signing round
	//
	for i, signer := range round5Signers {
		round6Messages[i], err = signer.SignRound6(
			signatureUnmaskDecrypted, messageHash,
		)
		if err != nil {
			t.Fatal(err)
		}
	}

	signature, err := round5Signers[0].CombineRound6Messages(round6Messages)
	if err != nil {
		t.Fatal(err)
	}

	//
	// Verify produced signature
	//
	err = verifySignatureInBitcoin(
		parameters.Curve,
		messageHash,
		round5Signers[0].dsaKey.PublicKey,
		signature,
	)
	if err != nil {
		t.Logf("H: %v\n", new(big.Int).SetBytes(messageHash))
		t.Logf("R: %v\n", signature.R)
		t.Logf("S: %v\n", signature.S)
		t.Logf("X: %v\n", round5Signers[0].dsaKey.PublicKey.X)
		t.Logf("Y: %v\n", round5Signers[0].dsaKey.PublicKey.Y)
		t.Fatalf("signature verification in bitcoin failed [%v]", err)
	}

	err = verifySignatureInEthereum(
		parameters.Curve,
		messageHash,
		round5Signers[0].dsaKey.PublicKey,
		signature,
	)
	if err != nil {
		fmt.Printf("H: %v\n", new(big.Int).SetBytes(messageHash))
		fmt.Printf("R: %v\n", signature.R)
		fmt.Printf("S: %v\n", signature.S)
		fmt.Printf("X: %v\n", round5Signers[0].dsaKey.PublicKey.X)
		fmt.Printf("Y: %v\n", round5Signers[0].dsaKey.PublicKey.Y)
		t.Fatalf("signature verification in ethereum failed [%v]", err)
	}
}

// Test31ByteSignatureRS is used to confirm that our signature verification
// algorithms work as expected if R or S is not 32 bytes long.
func Test31ByteSignatureRS(t *testing.T) {
	curve256 := secp256k1.S256()
	hash, _ := new(big.Int).SetString("8212313713408286312196617183996305874840581803582507267077647863768629906917", 10)
	publicKeyX, _ := new(big.Int).SetString("37243867901665327053253589157822427909743265115168368728514491795447858153874", 10)
	publicKeyY, _ := new(big.Int).SetString("48390273199951608338554842648959247259879464398730289908850755020939488517653", 10)
	signatureR, _ := new(big.Int).SetString("364606010805150545511962786008183839616327659698238570520068502825199705412", 10)
	signatureS, _ := new(big.Int).SetString("13781549995437993932032462513201290378095678483995393941371114222574658241776", 10)

	publicKey := curve.NewPoint(publicKeyX, publicKeyY)
	signature := &Signature{R: signatureR, S: signatureS}

	err := verifySignatureInBitcoin(
		curve256,
		hash.Bytes(),
		publicKey,
		signature,
	)
	if err != nil {
		t.Fatalf("signature verification in bitcoin failed [%v]", err)
	}

	err = verifySignatureInEthereum(
		curve256,
		hash.Bytes(),
		publicKey,
		signature,
	)
	if err != nil {
		t.Fatalf("signature verification in ethereum failed [%v]", err)
	}
}

func generateNewLocalGroup() (
	[]*LocalSigner,
	*PublicParameters,
	error,
) {
	publicParameters := &PublicParameters{

		Curve:                secp256k1.S256(),
		PaillierKeyBitLength: 2048,
	}

	signerGroup := &signerGroup{
		InitialGroupSize: 20,
		Threshold:        12,
	}

	paillierKeyGen, err := paillier.GetThresholdKeyGenerator(
		publicParameters.PaillierKeyBitLength,
		signerGroup.InitialGroupSize,
		signerGroup.Threshold,
		rand.Reader,
	)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"could not generate threshold Paillier keys [%v]", err,
		)
	}

	paillierKeys, err := paillierKeyGen.Generate()
	if err != nil {
		return nil, nil, fmt.Errorf(
			"could not generate threshold Paillier keys [%v]", err,
		)
	}

	zkpParameters, err := zkp.GeneratePublicParameters(
		paillierKeys[0].N,
		publicParameters.Curve,
	)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"could not generate public ZKP parameters [%v]", err,
		)
	}

	localSigners := make([]*LocalSigner, len(paillierKeys))
	for i := 0; i < len(localSigners); i++ {
		localSigners[i] = NewLocalSigner(
			paillierKeys[i], publicParameters, zkpParameters, signerGroup,
		)
	}

	return localSigners, publicParameters, nil
}

func verifySignatureInBitcoin(
	curve elliptic.Curve,
	hash []byte,
	publicKey *curve.Point,
	signature *Signature,
) error {
	btcPublicKey := &btcec.PublicKey{
		Curve: curve,
		X:     publicKey.X,
		Y:     publicKey.Y,
	}

	btcSignature := &btcec.Signature{R: signature.R, S: signature.S}

	// Verify if signature is valid for the given hash and public key
	if !btcSignature.Verify(hash, btcPublicKey) {
		return fmt.Errorf("Signature verification failed")
	}

	// Serialize type Signature {R,S} to DER format supported by Bitcoin:
	// 0x30 <length> 0x02 <length r> r 0x02 <length s> s
	btcSigSerialized := btcSignature.Serialize()

	// Deserialize signature in DER format to a Signature type {R,S}
	btcSigDeserialized, err := btcec.ParseDERSignature(btcSigSerialized, curve)
	if err != nil {
		return err
	}

	// Validate if ∂deserialized signature matches original signature
	if !btcSigDeserialized.IsEqual(btcSignature) {
		return fmt.Errorf("Signatures are not equal")
	}

	// All is fine
	return nil
}

func verifySignatureInEthereum(
	curve elliptic.Curve,
	hash []byte,
	publicKey *curve.Point,
	signature *Signature,
) error {
	ethPublicKey := &ecdsa.PublicKey{
		Curve: curve,
		X:     publicKey.X,
		Y:     publicKey.Y,
	}

	// We need to add padding to the R and S.
	// Ethereum requires that both values are 32 bytes long each.
	paddedR, err := byteutils.LeftPadTo32Bytes(signature.R.Bytes())
	if err != nil {
		return err
	}
	paddedS, err := byteutils.LeftPadTo32Bytes(signature.S.Bytes())
	if err != nil {
		return err
	}

	ethSignatureRS := append(paddedR, paddedS...)

	// Verify if signature is valid for the given hash and public key
	if !crypto.VerifySignature(
		crypto.CompressPubkey(ethPublicKey),
		hash,
		ethSignatureRS,
	) {
		return fmt.Errorf("Signature verification failed")
	}

	// All is fine
	return nil
}
