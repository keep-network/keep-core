package local

import (
	"fmt"
	"github.com/keep-network/keep-core/pkg/operator"
	"reflect"
	"testing"
)

func TestOperatorPublicKeyToChainPublicKey(t *testing.T) {
	_, operatorPublicKey, err := operator.GenerateKeyPair(DefaultCurve)
	if err != nil {
		t.Fatal(err)
	}

	chainPublicKey, err := OperatorPublicKeyToChainPublicKey(operatorPublicKey)
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
	_, operatorPublicKey, err := operator.GenerateKeyPair(DefaultCurve)
	if err != nil {
		t.Fatal(err)
	}

	// Alter the curve information to be undefined.
	operatorPublicKey.Curve = operator.Undefined

	_, err = OperatorPublicKeyToChainPublicKey(operatorPublicKey)

	expectedError := fmt.Errorf(
		"local chain supports only secp256k1 operator keys",
	)
	if !reflect.DeepEqual(expectedError, err) {
		t.Errorf(
			"unexpected error\nexpected: %v\nactual:   %v\n",
			expectedError,
			err,
		)
	}
}

func TestOperatorPrivateKeyToChainKeyPair(t *testing.T) {
	operatorPrivateKey, operatorPublicKey, err := operator.GenerateKeyPair(DefaultCurve)
	if err != nil {
		t.Fatal(err)
	}

	chainPrivateKey, chainPublicKey, err := OperatorPrivateKeyToChainKeyPair(
		operatorPrivateKey,
	)
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

	if !reflect.DeepEqual(chainPrivateKey.PublicKey, *chainPublicKey) {
		t.Errorf("chain private key contains wrong chain public key")
	}

	if chainPrivateKey.D.Cmp(operatorPrivateKey.D) != 0 {
		t.Errorf("chain private key has a wrong D parameter")
	}
}

func TestOperatorPrivateKeyToChainKeyPair_NotSecp256k1(t *testing.T) {
	operatorPrivateKey, _, err := operator.GenerateKeyPair(DefaultCurve)
	if err != nil {
		t.Fatal(err)
	}

	// Alter the curve information to be undefined.
	operatorPrivateKey.Curve = operator.Undefined

	_, _, err = OperatorPrivateKeyToChainKeyPair(
		operatorPrivateKey,
	)

	expectedError := fmt.Errorf(
		"cannot convert operator public key to chain public key: " +
			"[local chain supports only secp256k1 operator keys]",
	)
	if !reflect.DeepEqual(expectedError, err) {
		t.Errorf(
			"unexpected error\nexpected: %v\nactual:   %v\n",
			expectedError,
			err,
		)
	}
}
