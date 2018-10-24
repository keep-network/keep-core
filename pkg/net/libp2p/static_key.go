package libp2p

import (
	"github.com/btcsuite/btcd/btcec"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	libp2pcrypto "github.com/libp2p/go-libp2p-crypto"
)

// toLibp2pKey transforms `go-ethereum`-based ECDSA key into format supported
// by `libp2p`. Because all curve parameters of secp256k1 curve defined by
// `go-ethereum` and all curve parameters of secp256k1 curve defined
// by `btcsuite` used by `lipb2b` under the hood are identical, we can simply
// rewrite the private key.
//
// `libp2p` does not recognize `go-ethereum` curves and when it comes to
// creating peer's ID or deserializing the key, operation fails with
// unrecognized curve error.
func toLibp2pKey(ethereumKey *keystore.Key) *libp2pcrypto.Secp256k1PrivateKey {
	privKey, _ := btcec.PrivKeyFromBytes(
		btcec.S256(), ethereumKey.PrivateKey.D.Bytes(),
	)

	return (*libp2pcrypto.Secp256k1PrivateKey)(privKey)
}
