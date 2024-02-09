package libp2p

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/keep-network/keep-core/pkg/operator"
	libp2pcrypto "github.com/libp2p/go-libp2p/core/crypto"
)

func TestOperatorPrivateKeyToNetworkKeyPair(t *testing.T) {
	operatorPrivateKey, operatorPublicKey, err := operator.GenerateKeyPair(DefaultCurve)
	if err != nil {
		t.Fatal(err)
	}

	networkPrivateKey, networkPublicKey, err := operatorPrivateKeyToNetworkKeyPair(
		operatorPrivateKey,
	)
	if err != nil {
		t.Fatal(err)
	}

	// To access X and Y coordinates `networkPublicKey` needs to be converted
	// to the underlying `type btcec.PublicKey` as the type
	// `libp2pcrypto.Secp256k1PublicKey` does not expose them
	btcecPublicKey := (*btcec.PublicKey)(networkPublicKey)

	if btcecPublicKey.X().Cmp(operatorPublicKey.X) != 0 {
		t.Errorf("network public key has a wrong X coordinate")
	}
	if btcecPublicKey.Y().Cmp(operatorPublicKey.Y) != 0 {
		t.Errorf("network public key has a wrong Y coordinate")
	}

	// To access the D value `networkPrivateKey`needs to be converted to the
	// underlying `type btcec.PrivateKey` as the type
	// `libp2pcrypto.Secp256k1PrivateKey` does not expose it.
	btcecPrivateKey := (*btcec.PrivateKey)(networkPrivateKey)

	if !reflect.DeepEqual(
		btcecPrivateKey.PubKey(),
		btcecPublicKey,
	) {
		t.Errorf("network private key contains wrong network public key")
	}

	if extractD(btcecPrivateKey).Cmp(operatorPrivateKey.D) != 0 {
		t.Errorf("network private key has a wrong D parameter")
	}
}

func TestOperatorPrivateKeyToNetworkKeyPair_NotSecp256k1(t *testing.T) {
	operatorPrivateKey, _, err := operator.GenerateKeyPair(DefaultCurve)
	if err != nil {
		t.Fatal(err)
	}

	// Alter the curve information to be non-supported.
	operatorPrivateKey.Curve = -1

	_, _, err = operatorPrivateKeyToNetworkKeyPair(
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

func TestOperatorPrivateKeyToNetworkKeyPair_OperatorKeyBytesLengthShorterThan32(t *testing.T) {
	operatorPublicKey := operator.PublicKey{
		Curve: operator.Secp256k1,
		X:     big.NewInt(0),
		Y:     big.NewInt(0),
	}
	operatorPublicKey.X.UnmarshalText([]byte(
		"97427329544416289364384965587373348295173301176370762853904909401438964894172",
	))
	operatorPublicKey.Y.UnmarshalText([]byte(
		"74629372855232926020228171103363510898329534375325069073334079268933673928899",
	))

	operatorPrivateKey := &operator.PrivateKey{
		PublicKey: operatorPublicKey,
		D:         big.NewInt(0),
	}
	// The D value will be 31-byte long.
	operatorPrivateKey.D.UnmarshalText([]byte(
		"93742012702679576205257802297208635061933563746198248090895937324950509169",
	))

	networkPrivateKey, networkPublicKey, err := operatorPrivateKeyToNetworkKeyPair(
		operatorPrivateKey,
	)
	if err != nil {
		t.Fatal(err)
	}

	// To access X and Y coordinates `networkPublicKey` needs to be converted
	// to the underlying `type btcec.PublicKey` as the type
	// `libp2pcrypto.Secp256k1PublicKey` does not expose them
	btcecPublicKey := (*btcec.PublicKey)(networkPublicKey)

	if btcecPublicKey.X().Cmp(operatorPublicKey.X) != 0 {
		t.Errorf("network public key has a wrong X coordinate")
	}
	if btcecPublicKey.Y().Cmp(operatorPublicKey.Y) != 0 {
		t.Errorf("network public key has a wrong Y coordinate")
	}

	// To access the D value `networkPrivateKey`needs to be converted to the
	// underlying `type btcec.PrivateKey` as the type
	// `libp2pcrypto.Secp256k1PrivateKey` does not expose it.
	btcecPrivateKey := (*btcec.PrivateKey)(networkPrivateKey)

	if !reflect.DeepEqual(
		btcecPrivateKey.PubKey(),
		btcecPublicKey,
	) {
		t.Errorf("network private key contains wrong network public key")
	}

	if extractD(btcecPrivateKey).Cmp(operatorPrivateKey.D) != 0 {
		t.Errorf("network private key has a wrong D parameter")
	}
}

func TestOperatorPrivateKeyToNetworkKeyPair_OperatorKeyBytesLengthEqualTo32(t *testing.T) {
	operatorPublicKey := operator.PublicKey{
		Curve: operator.Secp256k1,
		X:     big.NewInt(0),
		Y:     big.NewInt(0),
	}
	operatorPublicKey.X.UnmarshalText([]byte(
		"100773973786886702528067362833887681119299455607054007058082928411975884796343",
	))
	operatorPublicKey.Y.UnmarshalText([]byte(
		"114537311513994193611577651952591602316970335059766760441557249896390043285921",
	))

	operatorPrivateKey := &operator.PrivateKey{
		PublicKey: operatorPublicKey,
		D:         big.NewInt(0),
	}
	// The D value will be 32-byte long.
	operatorPrivateKey.D.UnmarshalText([]byte(
		"27821820627376159532519075045721078275565093530261898135224182555762310064846",
	))

	networkPrivateKey, networkPublicKey, err := operatorPrivateKeyToNetworkKeyPair(
		operatorPrivateKey,
	)
	if err != nil {
		t.Fatal(err)
	}

	// To access X and Y coordinates `networkPublicKey` needs to be converted
	// to the underlying `type btcec.PublicKey` as the type
	// `libp2pcrypto.Secp256k1PublicKey` does not expose them.
	btcecPublicKey := (*btcec.PublicKey)(networkPublicKey)

	if btcecPublicKey.X().Cmp(operatorPublicKey.X) != 0 {
		t.Errorf("network public key has a wrong X coordinate")
	}
	if btcecPublicKey.Y().Cmp(operatorPublicKey.Y) != 0 {
		t.Errorf("network public key has a wrong Y coordinate")
	}

	// To access the D value `networkPrivateKey`needs to be converted to the
	// underlying `type btcec.PrivateKey` as the type
	// `libp2pcrypto.Secp256k1PrivateKey` does not expose it.
	btcecPrivateKey := (*btcec.PrivateKey)(networkPrivateKey)

	if !reflect.DeepEqual(
		btcecPrivateKey.PubKey(),
		btcecPublicKey,
	) {
		t.Errorf("network private key contains wrong network public key")
	}

	if extractD(btcecPrivateKey).Cmp(operatorPrivateKey.D) != 0 {
		t.Errorf("network private key has a wrong D parameter")
	}
}

func TestOperatorPrivateKeyToNetworkKeyPair_OperatorKeyBytesLengthLongerThan32(t *testing.T) {
	operatorPrivateKey, _, err := operator.GenerateKeyPair(DefaultCurve)
	if err != nil {
		t.Fatal(err)
	}

	// Alter the D value so that its length exceeds 32 bytes
	operatorPrivateKeyBytes := [33]byte{0: 1}
	operatorPrivateKey.D.SetBytes(operatorPrivateKeyBytes[:])

	_, _, err = operatorPrivateKeyToNetworkKeyPair(
		operatorPrivateKey,
	)

	expectedError := fmt.Errorf(
		"operator private key byte length is greater than 32 bytes",
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

	networkPublicKey, err := operatorPublicKeyToNetworkPublicKey(
		operatorPublicKey,
	)
	if err != nil {
		t.Fatal(err)
	}

	// To access X and Y coordinates `networkPublicKey` needs to be converted
	// to the underlying `type btcec.PublicKey` as the type
	// `libp2pcrypto.Secp256k1PublicKey` does not expose them
	btcecPublicKey := (*btcec.PublicKey)(networkPublicKey)

	if btcecPublicKey.X().Cmp(operatorPublicKey.X) != 0 {
		t.Errorf("network public key has a wrong X coordinate")
	}
	if btcecPublicKey.Y().Cmp(operatorPublicKey.Y) != 0 {
		t.Errorf("network public key has a wrong Y coordinate")
	}
}

func TestOperatorPublicKeyToNetworkPublicKey_NotSecp256k1(t *testing.T) {
	_, operatorPublicKey, err := operator.GenerateKeyPair(DefaultCurve)
	if err != nil {
		t.Fatal(err)
	}

	// Alter the curve information to be non-supported.
	operatorPublicKey.Curve = -1

	_, err = operatorPublicKeyToNetworkPublicKey(
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

	operatorPublicKey, err := networkPublicKeyToOperatorPublicKey(networkPublicKey)
	if err != nil {
		t.Fatal(err)
	}

	if operatorPublicKey.Curve != operator.Secp256k1 {
		t.Errorf("operator public key uses the wrong curve")
	}

	libp2pPublicKey, ok := networkPublicKey.(*libp2pcrypto.Secp256k1PublicKey)
	if !ok {
		t.Fatal("network public key is not an instance of the libp2p secp256k1 public key")
	}

	// To access X and Y coordinates `libp2pPublicKey` needs to be converted
	// to the underlying `type btcec.PublicKey` as the type
	// `libp2pcrypto.Secp256k1PublicKey` does not expose them
	btcecPublicKey := (*btcec.PublicKey)(libp2pPublicKey)

	if operatorPublicKey.X.Cmp(btcecPublicKey.X()) != 0 {
		t.Errorf("operator public key has a wrong X coordinate")
	}
	if operatorPublicKey.Y.Cmp(btcecPublicKey.Y()) != 0 {
		t.Errorf("operator public key has a wrong Y coordinate")
	}
}

func TestNetworkPublicKeyToOperatorPublicKey_NotSecp256k1(t *testing.T) {
	_, networkPublicKey, err := libp2pcrypto.GenerateEd25519Key(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	_, err = networkPublicKeyToOperatorPublicKey(networkPublicKey)

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

func extractD(privateKey *btcec.PrivateKey) *big.Int {
	privateKeyBytes := privateKey.Key.Bytes()
	return new(big.Int).SetBytes(privateKeyBytes[:])
}
