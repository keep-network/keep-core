package bls

import (
	"crypto/rand"
	"math/big"
	mrand "math/rand"
	"testing"

	"github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/keep-network/keep-core/pkg/internal/testutils"
)

// Test verifying BLS aggregated signature.
func TestAggregateBLS(t *testing.T) {

	// Public keys and signatures to aggregate.
	var signatures []*bn256.G1
	var publicKeys []*bn256.G2

	// Message to sign.
	msg := []byte("Hello!")

	for i := 0; i < 100; i++ {
		// Get secret key.
		k, _, err := bn256.RandomG1(rand.Reader)

		if err != nil {
			t.Errorf("Error generating random point on G1")
		}

		// Get public key.
		pub := new(bn256.G2).ScalarBaseMult(k)
		publicKeys = append(publicKeys, pub)

		// Sign the message.
		sig := Sign(k, msg)
		signatures = append(signatures, sig)
	}

	aggSig := AggregateG1Points(signatures)
	aggPub := AggregateG2Points(publicKeys)

	result := Verify(aggPub, msg, aggSig)

	if !result {
		t.Errorf("Error verifying BLS multi signature.")
	}
}

// Test verifying BLS threshold signature.
func TestThresholdBLS(t *testing.T) {

	message := []byte("Hello!")

	numOfPlayers := 5
	threshold := 3

	var masterSecretKey []*big.Int
	var masterPublicKey []*bn256.G2
	var signatureShares []*SignatureShare
	var publicKeyShares []*PublicKeyShare

	// Set up master keys. Based on Shamir's Secret Sharing scheme these are
	// polynomial coefficients where the first one is the secret and the
	// rest (threshold - 1) are sufficient to reconstruct this secret.
	for i := 0; i < threshold; i++ {
		sk, pk, _ := bn256.RandomG2(rand.Reader)
		masterSecretKey = append(masterSecretKey, sk)
		masterPublicKey = append(masterPublicKey, pk)
	}

	// Each member of the group gets their secret share by evaluating their
	// participant index over polynomial based on the coefficients from
	// masterSecretKey array. Threshold amount of these secret shares are
	// sufficient to reconstruct the main secret.
	// The resulting shares are used to sign the same message creating
	// signature shares. Threshold amount of these signature shares are
	// sufficient to reconstruct the signature which is the same value as
	// if the message was sign with the main secret.
	for i := 0; i < numOfPlayers; i++ {
		secretKeyShare := GetSecretKeyShare(masterSecretKey, int(i))
		publicKeyShares = append(publicKeyShares, secretKeyShare.PublicKeyShare())
		signatureShare := Sign(secretKeyShare.V, message)
		signatureShares = append(signatureShares, &SignatureShare{
			I: i + 1, // always 1-indexed
			V: signatureShare,
		})
	}

	// Shuffle signatureShares array. It's irrelevant which signatures shares
	// are used and in what order as long as they carry the corresponding
	// participant index.
	for i := range signatureShares {
		j := mrand.Intn(1 + i)
		signatureShares[i], signatureShares[j] = signatureShares[j], signatureShares[i]
	}

	// Get full BLS signature. Only threshold amount  of valid shared will be
	// used to reconstruct the signature. The resulting signature is the same
	// as if it was produced using master secret key.
	signature, _ := RecoverSignature(signatureShares, threshold)

	// Recovered public key should be the same as the main public key.
	publicKey, _ := RecoverPublicKey(publicKeyShares, threshold)
	testutils.AssertBytesEqual(t, publicKey.Marshal(), masterPublicKey[0].Marshal())

	result := Verify(publicKey, message, signature)

	if !result {
		t.Errorf("Error verifying BLS threshold signature.")
	}
}
