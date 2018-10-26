package key

import (
	"crypto/ecdsa"
	"io"

	"github.com/btcsuite/btcd/btcec"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	libp2pcrypto "github.com/libp2p/go-libp2p-crypto"
	"github.com/pborman/uuid"
)

// StaticNetworkKey represents static peer's key equal to the delegate key
// associated with an on-chain stake. It is used to authenticate the peer and
// for message attributability - each message leaving the peer is signed with
// its static key.
type StaticNetworkKey = libp2pcrypto.Secp256k1PrivateKey

// GenerateStaticNetworkKey generates a new, random static key based on
// secp256k1 ethereum curve.
func GenerateStaticNetworkKey(rand io.Reader) (*StaticNetworkKey, error) {
	id := uuid.NewRandom()

	ecdsaKey, err := ecdsa.GenerateKey(secp256k1.S256(), rand)
	if err != nil {
		return nil, err
	}

	key := &keystore.Key{
		Id:         id,
		Address:    crypto.PubkeyToAddress(ecdsaKey.PublicKey),
		PrivateKey: ecdsaKey,
	}

	return EthereumKeyToNetworkKey(key), nil
}

// EthereumKeyToNetworkKey transforms `go-ethereum`-based ECDSA key into format
// supported by `libp2p`. Because all curve parameters of secp256k1 curve
// defined by `go-ethereum` and all curve parameters of secp256k1 curve defined
// by `btcsuite` used by `lipb2b` under the hood are identical, we can simply
// rewrite the private key.
//
// `libp2p` does not recognize `go-ethereum` curves and when it comes to
// creating peer's ID or deserializing the key, operation fails with
// unrecognized curve error. This is no longer a problem if we transform the
// key using this function.
func EthereumKeyToNetworkKey(ethereumKey *keystore.Key) *StaticNetworkKey {
	privKey, _ := btcec.PrivKeyFromBytes(
		btcec.S256(), ethereumKey.PrivateKey.D.Bytes(),
	)

	return (*StaticNetworkKey)(privKey)
}
