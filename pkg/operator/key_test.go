package operator

import (
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
)

const testMessage = "Safe And Secure"

func TestOperatorKeySignAndVerify(t *testing.T) {
	operatorPrivateKey, operatorPublicKey, err := GenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	var tests = map[string]struct {
		sign              func(h []byte, p *PrivateKey) (Signature, error)
		verificationError func(sig []byte) string
	}{
		"signature is equal to 65 bytes": {
			sign: func(hash []byte, priv *PrivateKey) (Signature, error) {
				return Sign(hash, priv)
			},
			verificationError: nil,
		},
		"signature is greater than 65 bytes": {
			sign: func(hash []byte, priv *PrivateKey) (Signature, error) {
				sig, err := Sign(hash, priv)
				if err != nil {
					return nil, err
				}
				sig = append(sig, byte('0'))
				return Signature(sig), nil
			},
			verificationError: func(sig []byte) string {
				return fmt.Sprintf(
					"malformed signature %+v with length %d",
					sig,
					len(sig),
				)
			},
		},
		"signature is equal to 64 bytes": {
			sign: func(hash []byte, priv *PrivateKey) (Signature, error) {
				sig, err := Sign(hash, priv)
				if err != nil {
					return nil, err
				}
				sig = sig[:len(sig)-1]
				return Signature(sig), nil
			},
			verificationError: nil,
		},
		"signature is less than 64 bytes": {
			sign: func(hash []byte, priv *PrivateKey) (Signature, error) {
				sig, err := Sign(hash, priv)
				if err != nil {
					return nil, err
				}
				return Signature(sig[:32]), nil
			},
			verificationError: func(signature []byte) string {
				return fmt.Sprintf(
					"malformed signature %+v with length %d",
					signature,
					len(signature),
				)
			},
		},
		"incorrect signature algorithm": {
			sign: func(hash []byte, priv *PrivateKey) (Signature, error) {
				// Use the go crypto library ecdsa sign method
				sig, err := priv.Sign(rand.Reader, hash, nil)
				if err != nil {
					return nil, err
				}
				return Signature(sig), nil
			},
			// verificationError: nil,
			verificationError: func(signature []byte) string {
				return fmt.Sprintf(
					"malformed signature %+v with length %d",
					signature,
					len(signature),
				)
			},
		},
		"invalid signature": {
			sign: func(hash []byte, priv *PrivateKey) (Signature, error) {
				// hash a different message
				fakeHash := crypto.Keccak256([]byte("fake!"))
				sig, err := Sign(fakeHash, priv)
				if err != nil {
					return nil, err
				}
				return Signature(sig), nil
			},
			// verificationError: nil,
			verificationError: func(signature []byte) string {
				return fmt.Sprint("failed to verify signature")
			},
		},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			hashedMessage := crypto.Keccak256([]byte(testMessage))
			sig, err := test.sign(hashedMessage, operatorPrivateKey)
			if err != nil {
				t.Fatal(err)
			}

			err = VerifySignature(operatorPublicKey, hashedMessage, sig)
			if err != nil && err.Error() != test.verificationError(sig) {
				t.Fatal(err)
			}
		})
	}
}
