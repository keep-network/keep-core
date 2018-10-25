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

// EthereumStaticKey is an implementation of StaticKey interface that supports
// ethereum keys based on secp256k1 curve.
type EthereumStaticKey struct {
	ethereumKey *keystore.Key
}

// NewEthereumStaticKey instantiates EhtereumStaticKey from the provided
// ethereum ECDSA key.
func NewEthereumStaticKey(ethereumKey *keystore.Key) *EthereumStaticKey {
	return &EthereumStaticKey{ethereumKey}
}

// GenerateEthereumStaticKey generates a new, random EthereumStaticKey based on
// secp256k1 curve.
func GenerateEthereumStaticKey(rand io.Reader) (*EthereumStaticKey, error) {
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

	return NewEthereumStaticKey(key), nil
}

// PrivateKey returns ethereum secp256k1 ECDSA key as a libp2p PrivKey instance.
func (esk *EthereumStaticKey) PrivateKey() libp2pcrypto.PrivKey {
	return esk.toLibp2pKey()
}

// toLibp2pKey transforms `go-ethereum`-based ECDSA key into format supported
// by `libp2p`. Because all curve parameters of secp256k1 curve defined by
// `go-ethereum` and all curve parameters of secp256k1 curve defined
// by `btcsuite` used by `lipb2b` under the hood are identical, we can simply
// rewrite the private key.
//
// `libp2p` does not recognize `go-ethereum` curves and when it comes to
// creating peer's ID or deserializing the key, operation fails with
// unrecognized curve error.
func (esk *EthereumStaticKey) toLibp2pKey() *libp2pcrypto.Secp256k1PrivateKey {
	privKey, _ := btcec.PrivKeyFromBytes(
		btcec.S256(), esk.ethereumKey.PrivateKey.D.Bytes(),
	)

	return (*libp2pcrypto.Secp256k1PrivateKey)(privKey)
}
