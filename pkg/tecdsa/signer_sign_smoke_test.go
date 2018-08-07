package tecdsa

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/btcsuite/btcd/btcec"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/keep-network/keep-core/pkg/tecdsa/curve"
	"github.com/keep-network/paillier"
)

func TestGenerateValidKey(t *testing.T) {
	messageHash := make([]byte, 32)

	_, err := rand.Read(messageHash)
	if err != nil {
		t.Fatal(err)
	}

	group, dsaKey, err := initializeNewLocalGroupWithFullKey()
	if err != nil {
		t.Fatal(err)
	}

	// Decrypt secretKey from E(secretKey)
	dShares := make([]*paillier.PartialDecryption, publicParameters.groupSize)
	for i, signer := range group {
		dShares[i] = signer.paillierKey.Decrypt(dsaKey.secretKey.C)
	}
	D, err := group[0].paillierKey.CombinePartialDecryptions(dShares)
	if err != nil {
		t.Fatal(err)
	}

	key := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: publicParameters.curve,
			X:     dsaKey.publicKey.X,
			Y:     dsaKey.publicKey.Y,
		},
		D: D,
	}

	r, s, err := ecdsa.Sign(rand.Reader, key, messageHash)
	if err != nil {
		t.Fatal(err)
	}

	if !ecdsa.Verify(&key.PublicKey, messageHash, r, s) {
		t.Fatal("Signature verification failed")
	}
}

func TestFullSignPath(t *testing.T) {
	messageHash := make([]byte, 32)

	_, err := rand.Read(messageHash)
	if err != nil {
		t.Fatal(err)
	}

	signers, err := initializeNewSignerGroup()
	if err != nil {
		t.Fatal(err)
	}

	round1Messages := make([]*SignRound1Message, len(signers))
	round2Messages := make([]*SignRound2Message, len(signers))
	round3Messages := make([]*SignRound3Message, len(signers))
	round4Messages := make([]*SignRound4Message, len(signers))
	round5Messages := make([]*SignRound5Message, len(signers))
	round6Messages := make([]*SignRound6Message, len(signers))

	round1Signers := make([]*Round1Signer, len(signers))
	round2Signers := make([]*Round2Signer, len(signers))
	round3Signers := make([]*Round3Signer, len(signers))
	round4Signers := make([]*Round4Signer, len(signers))
	round5Signers := make([]*Round5Signer, len(signers))

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
		round6Messages[i] = signer.SignRound6(
			signatureUnmaskDecrypted, messageHash,
		)
	}

	signature, err := round5Signers[0].CombineRound6Messages(round6Messages)
	if err != nil {
		t.Fatal(err)
	}

	//
	// Verify produced signature
	//
	err = verifySignatureInBitcoin(
		publicParameters.curve,
		messageHash,
		round5Signers[0].dsaKey.publicKey,
		signature,
	)
	if err != nil {
		t.Fatal(err)
	}

	err = verifySignatureInEthereum(
		publicParameters.curve,
		messageHash,
		round5Signers[0].dsaKey.publicKey,
		signature,
	)
	if err != nil {
		t.Fatal(err)
	}
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
