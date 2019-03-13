package operator

import (
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
		sign              func(hash []byte, priv *PrivateKey) ([]byte, error)
		signingError      error
		verificationError error
	}{
		"signature is equal to 65 bytes": {
			sign: func(hash []byte, priv *PrivateKey) ([]byte, error) {
				return Sign(hash, priv)
			},
			signingError:      nil,
			verificationError: nil,
		},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			hashedMessage := crypto.Keccak256([]byte(testMessage))
			sig, err := test.sign(hashedMessage, operatorPrivateKey)
			if err != nil && err != test.signingError {
				t.Fatal(err)
			}

			err = VerifySignature(operatorPublicKey, hashedMessage, sig)
			if err != nil && err != test.verificationError {
				t.Fatal(err)
			}
		})
	}
}
