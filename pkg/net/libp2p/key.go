package libp2p

import (
	"fmt"

	"github.com/btcsuite/btcd/btcec"
	"github.com/keep-network/keep-core/pkg/operator"
	libp2pcrypto "github.com/libp2p/go-libp2p-core/crypto"
)

// OperatorPrivateKeyToLibp2pKeyPair converts an operator private key to
// the libp2p key pair that uses the libp2p-specific curve implementation.
func OperatorPrivateKeyToLibp2pKeyPair(operatorPrivateKey *operator.PrivateKey) (
	*libp2pcrypto.Secp256k1PrivateKey,
	*libp2pcrypto.Secp256k1PublicKey,
	error,
) {
	// Make sure that libp2p package receives only secp256k1 operator keys.
	if operatorPrivateKey.Curve != operator.Secp256k1 {
		return nil, nil, fmt.Errorf("libp2p supports only secp256k1 operator keys")
	}

	// Libp2p keys are actually btcec keys under the hood.
	btcecPrivateKey, btcecPublicKey := btcec.PrivKeyFromBytes(
		btcec.S256(), operatorPrivateKey.D.Bytes(),
	)

	libp2pPrivateKey := libp2pcrypto.Secp256k1PrivateKey(*btcecPrivateKey)
	libp2pPublicKey := libp2pcrypto.Secp256k1PublicKey(*btcecPublicKey)

	return &libp2pPrivateKey, &libp2pPublicKey, nil
}

// OperatorPublicKeyToLibp2pPublicKey converts an operator public key to
// the libp2p public key that uses the libp2p-specific curve implementation.
func OperatorPublicKeyToLibp2pPublicKey(
	operatorPublicKey *operator.PublicKey,
) (*libp2pcrypto.Secp256k1PublicKey, error) {
	// Make sure that libp2p package receives only secp256k1 operator keys.
	if operatorPublicKey.Curve != operator.Secp256k1 {
		return nil, fmt.Errorf("libp2p supports only secp256k1 operator keys")
	}

	return &libp2pcrypto.Secp256k1PublicKey{
		Curve: btcec.S256(),
		X:     operatorPublicKey.X,
		Y:     operatorPublicKey.Y,
	}, nil
}

// Libp2pPublicKeyToOperatorPublicKey converts a libp2p public key to
// the operator public key that is agnostic regarding the curve implementation.
func Libp2pPublicKeyToOperatorPublicKey(libp2pPublicKey libp2pcrypto.PubKey) (*operator.PublicKey, error) {
	switch publicKey := libp2pPublicKey.(type) {
	case *libp2pcrypto.Secp256k1PublicKey:
		return &operator.PublicKey{
			Curve: operator.Secp256k1,
			X:     publicKey.X,
			Y:     publicKey.Y,
		}, nil
	}
	return nil, fmt.Errorf("unrecognized libp2p public key type")
}
