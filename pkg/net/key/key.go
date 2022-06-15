package key

import (
	"crypto/ecdsa"
	"crypto/elliptic"

	btcec_v2 "github.com/btcsuite/btcd/btcec/v2"
	"github.com/keep-network/keep-core/pkg/operator"
	libp2pcrypto "github.com/libp2p/go-libp2p-core/crypto"
)

// NetworkPrivate represents peer's static key associated with an on-chain
// stake. It is used to authenticate the peer and for message attributability
// - each message leaving the peer is signed with its private network key.
type NetworkPrivate = libp2pcrypto.Secp256k1PrivateKey

// NetworkPublic represents peer's static key associated with an on-chain
// stake. It is used to authenticate the peer and for message attributability
// - each received message is validated against sender's public network key.
type NetworkPublic = libp2pcrypto.Secp256k1PublicKey

// GenerateStaticNetworkKey generates a new, random static key based on
// secp256k1 curve.
func GenerateStaticNetworkKey() (*NetworkPrivate, *NetworkPublic, error) {
	privKey, pubKey, err := operator.GenerateKeyPair()
	if err != nil {
		return nil, nil, err
	}
	networkPrivateKey, networkPublicKey := OperatorKeyToNetworkKey(privKey, pubKey)
	return networkPrivateKey, networkPublicKey, nil
}

// OperatorKeyToNetworkKey transforms a static ECDSA key into the format
// supported by the network layer. Because all curve parameters of the secp256k1
// curve defined by `go-ethereum`-based modules and all curve parameters of
// secp256k1 curve defined by `btcsuite` used by `lipb2b` under the hood are
// identical, we can simply rewrite the private key.
//
// `libp2p` does not recognize `go-ethereum`-based curves and when it comes to
// creating peer's ID or deserializing the key, operation fails with
// unrecognized curve error. This is no longer a problem if we transform the
// key using this function.
func OperatorKeyToNetworkKey(
	operatorPrivateKey *operator.PrivateKey,
	operatorPublicKey *operator.PublicKey,
) (*NetworkPrivate, *NetworkPublic) {
	privKey, pubKey := btcec_v2.PrivKeyFromBytes(operatorPrivateKey.D.Bytes())
	return (*NetworkPrivate)(privKey), (*NetworkPublic)(pubKey)
}

// NetworkPubKeyToChainAddress transforms the network public key into a chain
// account address, in a string format.
func NetworkPubKeyToChainAddress(publicKey *NetworkPublic) string {
	ecdsaKey := (*btcec_v2.PublicKey)(publicKey).ToECDSA()
	return operator.PubkeyToAddress(*ecdsaKey).String()
}

// Marshal takes a network public key, converts it into an ecdsa
// public key, and uses go's standard library elliptic marshal method to
// convert the public key into a slice of bytes in the correct format for the key
// type. This allows external consumers of this key to verify integrity of the
// key without having to understand the internals of the net pkg.
func Marshal(publicKey *NetworkPublic) []byte {
	ecdsaKey := (*btcec_v2.PublicKey)(publicKey).ToECDSA()
	return elliptic.Marshal(ecdsaKey.Curve, ecdsaKey.X, ecdsaKey.Y)
}

// Libp2pKeyToNetworkKey takes an interface type, libp2pcrypto.PubKey, and
// returns the concrete type specific to this package. If it fails to do so, it
// returns nil.
func Libp2pKeyToNetworkKey(publicKey libp2pcrypto.PubKey) *NetworkPublic {
	switch networkKey := publicKey.(type) {
	case *libp2pcrypto.Secp256k1PublicKey:
		return (*NetworkPublic)(networkKey)
	}
	return nil
}

// ECDSAKeyToNetworkKey takes ecdsa.PublicKey from Go standard library and turns
// it into public NetworkKey.
func ECDSAKeyToNetworkKey(publicKey *ecdsa.PublicKey) *NetworkPublic {
	var x, y btcec_v2.FieldVal
	x.SetByteSlice(publicKey.X.Bytes())
	y.SetByteSlice(publicKey.Y.Bytes())

	return (*NetworkPublic)(btcec_v2.NewPublicKey(&x, &y))
}

// NetworkKeyToECDSAKey takes the public NetworkKey and turns it into
// ecdsa.PublicKey from Go standard library.
func NetworkKeyToECDSAKey(publicKey *NetworkPublic) *ecdsa.PublicKey {
	return (*btcec_v2.PublicKey)(publicKey).ToECDSA()
}
