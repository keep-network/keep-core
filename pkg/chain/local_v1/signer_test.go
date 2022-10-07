package local_v1

import (
	"github.com/keep-network/keep-core/pkg/operator"
	"math/big"
	"testing"
)

func TestSigner_PublicKeyToAddress(t *testing.T) {
	x, ok := new(big.Int).SetString(
		"3f89dfad9a9ace8437a2c752448b6de75aac78613ce97e0469f13c92006c54cb",
		16,
	)
	if !ok {
		t.Fatal("cannot set X coordinate")
	}

	y, ok := new(big.Int).SetString(
		"96bd09fc1b36e316a369a82f5d5e11c3225352deafca2772f8e0f62813cfccb3",
		16,
	)
	if !ok {
		t.Fatal("cannot set Y coordinate")
	}

	operatorPublicKey := &operator.PublicKey{
		Curve: operator.Secp256k1,
		X:     x,
		Y:     y,
	}

	// The operator private key is not relevant in this scenario.
	signer := NewSigner(&operator.PrivateKey{
		PublicKey: *operatorPublicKey,
	})

	address, err := signer.PublicKeyToAddress(operatorPublicKey)
	if err != nil {
		t.Fatal(err)
	}

	expectedAddress :=
		"043f89dfad9a9ace8437a2c752448b6de75aac78613ce97e0469f13c92006c54cb9" +
			"6bd09fc1b36e316a369a82f5d5e11c3225352deafca2772f8e0f62813cfccb3"
	actualAddress := address.String()
	if expectedAddress != actualAddress {
		t.Errorf(
			"unexpected address\nexpected: %v\nactual:   %v\n",
			expectedAddress,
			actualAddress,
		)
	}
}
