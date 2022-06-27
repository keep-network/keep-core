package libp2p

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"reflect"
	"testing"

	"github.com/keep-network/keep-core/pkg/operator"
	libp2pcrypto "github.com/libp2p/go-libp2p-core/crypto"
)

func TestOperatorPrivateKeyToNetworkKeyPair(t *testing.T) {
	operatorPrivateKey, operatorPublicKey, err := operator.GenerateKeyPair(DefaultCurve)
	if err != nil {
		t.Fatal(err)
	}

	networkPrivateKey, networkPublicKey, err := OperatorPrivateKeyToNetworkKeyPair(
		operatorPrivateKey,
	)
	if err != nil {
		t.Fatal(err)
	}

	rawNetworkPublicKey, err := networkPublicKey.Raw()
	if err != nil {
		t.Fatal(err)
	}
	x, y := elliptic.Unmarshal(DefaultCurve, rawNetworkPublicKey)

	if x.Cmp(operatorPublicKey.X) != 0 {
		t.Errorf("network public key has a wrong X coordinate")
	}
	if y.Cmp(operatorPublicKey.Y) != 0 {
		t.Errorf("network public key has a wrong Y coordinate")
	}

	if !reflect.DeepEqual(
		networkPrivateKey.PublicKey,
		(ecdsa.PublicKey)(*networkPublicKey),
	) {
		t.Errorf("network private key contains wrong network public key")
	}

	if networkPrivateKey.D.Cmp(operatorPrivateKey.D) != 0 {
		t.Errorf("network private key has a wrong D parameter")
	}
}

func TestOperatorPrivateKeyToNetworkKeyPair_NotSecp256k1(t *testing.T) {
	operatorPrivateKey, _, err := operator.GenerateKeyPair(DefaultCurve)
	if err != nil {
		t.Fatal(err)
	}

	// Alter the curve information to be undefined.
	operatorPrivateKey.Curve = operator.Undefined

	_, _, err = OperatorPrivateKeyToNetworkKeyPair(
		operatorPrivateKey,
	)

	expectedError := fmt.Errorf(
		"libp2p supports only secp256k1 operator keys",
	)
	if !reflect.DeepEqual(expectedError, err) {
		t.Errorf(
			"unexpected error\nexpected: %v\nactual:   %v\n",
			expectedError,
			err,
		)
	}
}

func TestOperatorPublicKeyToNetworkPublicKey(t *testing.T) {
	_, operatorPublicKey, err := operator.GenerateKeyPair(DefaultCurve)
	if err != nil {
		t.Fatal(err)
	}

	networkPublicKey, err := OperatorPublicKeyToNetworkPublicKey(
		operatorPublicKey,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(networkPublicKey.Curve.Params(), DefaultCurve.Params()) {
		t.Errorf("network public key uses the wrong curve")
	}

	if networkPublicKey.X.Cmp(operatorPublicKey.X) != 0 {
		t.Errorf("network public key has a wrong X coordinate")
	}
	if networkPublicKey.Y.Cmp(operatorPublicKey.Y) != 0 {
		t.Errorf("network public key has a wrong Y coordinate")
	}
}

func TestOperatorPublicKeyToNetworkPublicKey_NotSecp256k1(t *testing.T) {
	_, operatorPublicKey, err := operator.GenerateKeyPair(DefaultCurve)
	if err != nil {
		t.Fatal(err)
	}

	// Alter the curve information to be undefined.
	operatorPublicKey.Curve = operator.Undefined

	_, err = OperatorPublicKeyToNetworkPublicKey(
		operatorPublicKey,
	)

	expectedError := fmt.Errorf(
		"libp2p supports only secp256k1 operator keys",
	)
	if !reflect.DeepEqual(expectedError, err) {
		t.Errorf(
			"unexpected error\nexpected: %v\nactual:   %v\n",
			expectedError,
			err,
		)
	}
}

func TestNetworkPublicKeyToOperatorPublicKey(t *testing.T) {
	_, networkPublicKey, err := libp2pcrypto.GenerateSecp256k1Key(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	operatorPublicKey, err := NetworkPublicKeyToOperatorPublicKey(networkPublicKey)
	if err != nil {
		t.Fatal(err)
	}

	if operatorPublicKey.Curve != operator.Secp256k1 {
		t.Errorf("operator public key uses the wrong curve")
	}

	if operatorPublicKey.X.Cmp(
		networkPublicKey.(*libp2pcrypto.Secp256k1PublicKey).X,
	) != 0 {
		t.Errorf("operator public key has a wrong X coordinate")
	}
	if operatorPublicKey.Y.Cmp(
		networkPublicKey.(*libp2pcrypto.Secp256k1PublicKey).Y,
	) != 0 {
		t.Errorf("operator public key has a wrong Y coordinate")
	}
}

func TestNetworkPublicKeyToOperatorPublicKey_NotSecp256k1(t *testing.T) {
	_, networkPublicKey, err := libp2pcrypto.GenerateEd25519Key(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	_, err = NetworkPublicKeyToOperatorPublicKey(networkPublicKey)

	expectedError := fmt.Errorf(
		"unrecognized libp2p public key type",
	)
	if !reflect.DeepEqual(expectedError, err) {
		t.Errorf(
			"unexpected error\nexpected: %v\nactual:   %v\n",
			expectedError,
			err,
		)
	}
}
