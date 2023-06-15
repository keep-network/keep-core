package bls

import (
	"crypto/rand"
	"math/big"
	"testing"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"

	"github.com/keep-network/keep-core/internal/testutils"
)

func TestSignAndVerify(t *testing.T) {
	pi, _ := new(big.Int).SetString("31415926535897932384626433832795028841971693993751058209749445923078164062862", 10)
	message := pi.Bytes()

	secretKey := big.NewInt(123)
	publicKey := new(bn256.G2).ScalarBaseMult(secretKey)

	signature := Sign(secretKey, message)

	if !Verify(publicKey, message, signature) {
		t.Fatalf("Expected signature to be valid")
	}
}

func TestSignAndVerifyG1(t *testing.T) {
	pi, _ := new(big.Int).SetString("31415926535897932384626433832795028841971693993751058209749445923078164062862", 10)
	message := new(bn256.G1).ScalarBaseMult(pi)

	secretKey := big.NewInt(123)
	publicKey := new(bn256.G2).ScalarBaseMult(secretKey)

	signature := SignG1(secretKey, message)

	if !VerifyG1(publicKey, message, signature) {
		t.Fatalf("Expected signature to be valid")
	}
}

// Test verifying BLS aggregated signature.
func TestAggregateBLS(t *testing.T) {
	// Public keys and signatures to aggregate.
	var signatures []*bn256.G1
	var publicKeys []*bn256.G2

	// Message to sign.
	pi, _ := new(big.Int).SetString("31415926535897932384626433832795028841971693993751058209749445923078164062862", 10)
	message := new(bn256.G1).ScalarBaseMult(pi)

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
		sig := SignG1(k, message)
		signatures = append(signatures, sig)
	}

	aggSig := AggregateG1Points(signatures)
	aggPub := AggregateG2Points(publicKeys)

	result := VerifyG1(aggPub, message, aggSig)

	if !result {
		t.Errorf("Error verifying BLS multi signature.")
	}
}

// Test verifying BLS threshold signature.
func TestThresholdBLS(t *testing.T) {
	pi, _ := new(big.Int).SetString("31415926535897932384626433832795028841971693993751058209749445923078164062862", 10)
	message := new(bn256.G1).ScalarBaseMult(pi)

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
	//
	// NOTE: The loop begins from 1, not 0, as shares must be 1-indexed.
	for i := 1; i <= numOfPlayers; i++ {
		secretKeyShare := GetSecretKeyShare(masterSecretKey, int(i))
		publicKeyShares = append(publicKeyShares, secretKeyShare.PublicKeyShare())
		signatureShare := SignG1(secretKeyShare.V, message)
		signatureShares = append(signatureShares, &SignatureShare{
			I: i,
			V: signatureShare,
		})
	}

	var tests = map[string]struct {
		signatureShares func() []*SignatureShare
	}{
		"all shares in order": {
			signatureShares: func() []*SignatureShare {
				return signatureShares
			},
		},
		"all shares shuffled": {
			signatureShares: func() []*SignatureShare {
				return []*SignatureShare{
					signatureShares[1],
					signatureShares[0],
					signatureShares[3],
					signatureShares[2],
					signatureShares[4],
				}
			},
		},
		"shares [1,2,3,4]": {
			signatureShares: func() []*SignatureShare {
				return signatureShares[0:4]
			},
		},
		"shares [1,2,3]": {
			signatureShares: func() []*SignatureShare {
				return signatureShares[0:3]
			},
		},
		"shares [2,3,4,5]": {
			signatureShares: func() []*SignatureShare {
				return signatureShares[1:5]
			},
		},
		"shares [3,4,5]": {
			signatureShares: func() []*SignatureShare {
				return signatureShares[2:5]
			},
		},
		"shares [1,5,3]": {
			signatureShares: func() []*SignatureShare {
				return []*SignatureShare{
					signatureShares[0],
					signatureShares[4],
					signatureShares[2],
				}
			},
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			// Get full BLS signature. Only threshold amount  of valid shared will be
			// used to reconstruct the signature. The resulting signature is the same
			// as if it was produced using master secret key.
			signature, _ := RecoverSignature(test.signatureShares(), threshold)

			// Recovered public key should be the same as the main public key.
			publicKey, _ := RecoverPublicKey(publicKeyShares, threshold)
			testutils.AssertBytesEqual(t, publicKey.Marshal(), masterPublicKey[0].Marshal())

			if !VerifyG1(publicKey, message, signature) {
				t.Errorf("Error verifying BLS threshold signature.")
			}
		})
	}

}
