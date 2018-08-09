package tecdsa

// TODO: rename to signer_smoke_test.go

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/btcsuite/btcd/btcec"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/keep-network/keep-core/pkg/tecdsa/curve"
	"github.com/keep-network/keep-core/pkg/tecdsa/zkp"
	"github.com/keep-network/paillier"
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

	publicKeyCommitmentMessages := make([]*PublicKeyShareCommitmentMessage, len(localSigners))
	keyShareRevealMessages := make([]*KeyShareRevealMessage, len(localSigners))

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
		signers[i] = &Signer{
			dsaKey:     dsaKey,
			signerCore: localSigner.signerCore,
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
		round5Signers[0].dsaKey.publicKey,
		signature,
	)
	if err != nil {
		t.Fatal(err)
	}

	err = verifySignatureInEthereum(
		parameters.Curve,
		messageHash,
		round5Signers[0].dsaKey.publicKey,
		signature,
	)
	if err != nil {
		t.Fatal(err)
	}
}

func generateNewLocalGroup() (
	[]*LocalSigner,
	*PublicParameters,
	error,
) {
	parameters := &PublicParameters{
		GroupSize:            20,
		Threshold:            12,
		Curve:                secp256k1.S256(),
		PaillierKeyBitLength: 2048,
	}

	paillierKeyGen, err := paillier.GetThresholdKeyGenerator(
		parameters.PaillierKeyBitLength,
		parameters.GroupSize,
		parameters.Threshold,
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
		parameters.Curve,
	)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"could not generate public ZKP parameters [%v]", err,
		)
	}

	members := make([]*signerCore, len(paillierKeys))
	for i := 0; i < len(members); i++ {
		members[i] = &signerCore{
			ID:              generateMemberID(),
			paillierKey:     paillierKeys[i],
			groupParameters: parameters,
			zkpParameters:   zkpParameters,
		}
	}

	localSigners := make([]*LocalSigner, len(members))
	for i := 0; i < len(localSigners); i++ {
		localSigners[i] = &LocalSigner{
			signerCore: *members[i],
		}
	}

	return localSigners, parameters, nil
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

	ethSignatureRS := append(signature.R.Bytes(), signature.S.Bytes()...)

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
