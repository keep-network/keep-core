package bls

import (
	"crypto/rand"
	"math/big"
	mrand "math/rand"
	"testing"

	"github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
)

// Test verifying BLS multi signature.
func TestMultisigBLS(t *testing.T) {

	// Public keys and signatures to aggregate.
	var signatures []*bn256.G1
	var publicKeys []*bn256.G2

	// Message to sign.
	msg := []byte("Hello!")

	for i := 0; i < 100; i++ {
		// Get private key.
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

	msg := []byte("Hello!")

	numOfPlayers := 5
	threshold := 3

	var masterSecretKey []*big.Int
	var masterPublicKey []*bn256.G2
	var signatureShares []*bn256.G1

	// Set up master keys.
	for i := 0; i < threshold; i++ {
		sk, pk, _ := bn256.RandomG2(rand.Reader)
		masterSecretKey = append(masterSecretKey, sk)
		masterPublicKey = append(masterPublicKey, pk)
	}

	// Each member of the group signs the same message creating signature share.
	for i := 0; i < numOfPlayers; i++ {
		sk := SecretKeyShare(masterSecretKey, int64(i))
		share := Sign(sk, msg)
		signatureShares = append(signatureShares, share)
	}

	// Shuffle signatureShares array.
	for i := range signatureShares {
		j := mrand.Intn(1 + i)
		signatureShares[i], signatureShares[j] = signatureShares[j], signatureShares[i]
	}

	// Get full BLS signature. Use only threshold amount of signatures to
	// demonstrate it's a functional threshold construction.
	sig := Recover(signatureShares[:threshold], threshold)

	result := Verify(masterPublicKey[0], msg, sig)

	if !result {
		t.Errorf("Error verifying BLS threshold signature.")
	}
}
