package operator

import (
	"crypto/rand"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
)

const testMessage = "Safe And Secure"

func TestOperatorKeySignAndverify(t *testing.T) {
	operatorPrivateKey, operatorPublicKey, err := GenerateKeyPair(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	hashedMessage := crypto.Keccak256([]byte(testMessage))
	sig, err := Sign(hashedMessage, operatorPrivateKey)
	if err != nil {
		t.Fatal(err)
	}

	if err := VerifySignature(operatorPublicKey, hashedMessage, sig); err != nil {
		t.Fatal(err)
	}
}
