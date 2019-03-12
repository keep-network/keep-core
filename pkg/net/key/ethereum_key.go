package key

import (
	"io"

	"github.com/btcsuite/btcd/btcec"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/keep-network/keep-core/pkg/static"
	libp2pcrypto "github.com/libp2p/go-libp2p-crypto"
)

// NetworkPrivateKey represents peer's static key associated with an on-chain
// stake. It is used to authenticate the peer and for message attributability
// - each message leaving the peer is signed with its private network key.
type NetworkPrivateKey = libp2pcrypto.Secp256k1PrivateKey

// NetworkPublicKey represents peer's static key associated with an on-chain
// stake. It is used to authenticate the peer and for message attributability
// - each received message is validated against sender's public network key.
type NetworkPublicKey = libp2pcrypto.Secp256k1PublicKey

// GenerateStaticNetworkKey generates a new, random static key based on
// secp256k1 ethereum curve.
func GenerateStaticNetworkKey(rand io.Reader) (*NetworkPrivateKey, *NetworkPublicKey, error) {
	privKey, pubKey, err := static.GenerateStaticKeyPair(rand)
	if err != nil {
		return nil, nil, err
	}
	networkPrivateKey, networkPublicKey := StaticKeyToNetworkKey(privKey, pubKey)
	return networkPrivateKey, networkPublicKey, nil
}

// StaticKeyToNetworkKey transforms a static ECDSA key into the format supported
// by the network layer. Because all curve parameters of the secp256k1 curve
// defined by `go-ethereum` and all curve parameters of secp256k1 curve defined
// by `btcsuite` used by `lipb2b` under the hood are identical, we can simply
// rewrite the private key.
//
// `libp2p` does not recognize `go-ethereum` curves and when it comes to
// creating peer's ID or deserializing the key, operation fails with
// unrecognized curve error. This is no longer a problem if we transform the
// key using this function.
func StaticKeyToNetworkKey(
	staticPrivateKey *static.PrivateKey,
	staticPublicKey *static.PublicKey,
) (*NetworkPrivateKey, *NetworkPublicKey) {
	privKey, pubKey := btcec.PrivKeyFromBytes(
		btcec.S256(), staticPrivateKey.D.Bytes(),
	)
	return (*NetworkPrivateKey)(privKey), (*NetworkPublicKey)(pubKey)
}

// NetworkPubKeyToEthAddress transforms NetworkPublicKey into Ethereum account
// address in a string format.
func NetworkPubKeyToEthAddress(publicKey *NetworkPublicKey) string {
	ecdsaKey := (*btcec.PublicKey)(publicKey).ToECDSA()
	return crypto.PubkeyToAddress(*ecdsaKey).String()
}
