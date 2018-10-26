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

// GenerateStaticKey generates a new, random static key based on secp256k1
// ethereum curve.
func GenerateStaticKey(rand io.Reader) (*libp2pcrypto.Secp256k1PrivateKey, error) {
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

	return EthereumKeyToLibp2pKey(key), nil
}

// EthereumKeyToLibp2pKey transforms `go-ethereum`-based ECDSA key into format
// supported by `libp2p`. Because all curve parameters of secp256k1 curve
// defined by `go-ethereum` and all curve parameters of secp256k1 curve defined
// by `btcsuite` used by `lipb2b` under the hood are identical, we can simply
// rewrite the private key.
//
// `libp2p` does not recognize `go-ethereum` curves and when it comes to
// creating peer's ID or deserializing the key, operation fails with
// unrecognized curve error. This is no longer a problem if we transform the
// key using this function.
func EthereumKeyToLibp2pKey(
	ethereumKey *keystore.Key,
) *libp2pcrypto.Secp256k1PrivateKey {
	privKey, _ := btcec.PrivKeyFromBytes(
		btcec.S256(), ethereumKey.PrivateKey.D.Bytes(),
	)

	return (*libp2pcrypto.Secp256k1PrivateKey)(privKey)
}
