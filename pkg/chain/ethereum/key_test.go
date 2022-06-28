package ethereum

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/keep-network/keep-core/pkg/operator"
	"math/big"
	"reflect"
	"testing"
)

func TestChainPrivateKeyToOperatorKeyPair(t *testing.T) {
	chainPrivateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}

	operatorPrivateKey, operatorPublicKey, err := ChainPrivateKeyToOperatorKeyPair(
		chainPrivateKey,
	)
	if err != nil {
		t.Fatal(err)
	}

	if operatorPublicKey.Curve != operator.Secp256k1 {
		t.Errorf("operator public key uses the wrong curve")
	}

	if operatorPublicKey.X.Cmp(chainPrivateKey.X) != 0 {
		t.Errorf("operator public key has a wrong X coordinate")
	}
	if operatorPublicKey.Y.Cmp(chainPrivateKey.Y) != 0 {
		t.Errorf("operator public key has a wrong Y coordinate")
	}

	if !reflect.DeepEqual(operatorPrivateKey.PublicKey, *operatorPublicKey) {
		t.Errorf("operator private key contains wrong operator public key")
	}

	if operatorPrivateKey.D.Cmp(chainPrivateKey.D) != 0 {
		t.Errorf("network private key has a wrong D parameter")
	}
}

func TestChainPrivateKeyToOperatorKeyPair_NotSecp256k1(t *testing.T) {
	chainPrivateKey := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: elliptic.P256(),
		},
	}

	_, _, err := ChainPrivateKeyToOperatorKeyPair(
		chainPrivateKey,
	)

	expectedError := fmt.Errorf(
		"ethereum chain key does not use secp256k1 curve",
	)
	if !reflect.DeepEqual(expectedError, err) {
		t.Errorf(
			"unexpected error\nexpected: %v\nactual:   %v\n",
			expectedError,
			err,
		)
	}
}

func TestOperatorPublicKeyToChainPublicKey(t *testing.T) {
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

	chainPublicKey, err := operatorPublicKeyToChainPublicKey(operatorPublicKey)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(chainPublicKey.Curve.Params(), DefaultCurve.Params()) {
		t.Errorf("chain public key uses the wrong curve")
	}

	if chainPublicKey.X.Cmp(operatorPublicKey.X) != 0 {
		t.Errorf("chain public key has a wrong X coordinate")
	}
	if chainPublicKey.Y.Cmp(operatorPublicKey.Y) != 0 {
		t.Errorf("chain public key has a wrong Y coordinate")
	}
}

func TestOperatorPublicKeyToChainPublicKey_NotSecp256k1(t *testing.T) {
	operatorPublicKey := &operator.PublicKey{
		Curve: -1,
	}

	_, err := operatorPublicKeyToChainPublicKey(operatorPublicKey)

	expectedError := fmt.Errorf(
		"ethereum supports only secp256k1 operator keys",
	)
	if !reflect.DeepEqual(expectedError, err) {
		t.Errorf(
			"unexpected error\nexpected: %v\nactual:   %v\n",
			expectedError,
			err,
		)
	}
}

func TestOperatorPublicKeyToChainAddress(t *testing.T) {
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

	chainAddress, err := operatorPublicKeyToChainAddress(operatorPublicKey)
	if err != nil {
		t.Fatal(err)
	}

	expectedAddress := "0x09e303E34F5aC4350caF327aD92d752602f3B061"
	actualAddress := chainAddress.Hex()
	if expectedAddress != actualAddress {
		t.Errorf(
			"unexpected address\nexpected: %v\nactual:   %v\n",
			expectedAddress,
			actualAddress,
		)
	}
}

func TestOperatorPublicKeyToChainAddress_NotSecp256k1(t *testing.T) {
	operatorPublicKey := &operator.PublicKey{
		Curve: -1,
	}

	_, err := operatorPublicKeyToChainAddress(operatorPublicKey)

	expectedError := fmt.Errorf(
		"cannot convert from operator to chain key: " +
			"[ethereum supports only secp256k1 operator keys]",
	)
	if !reflect.DeepEqual(expectedError, err) {
		t.Errorf(
			"unexpected error\nexpected: %v\nactual:   %v\n",
			expectedError,
			err,
		)
	}
}
