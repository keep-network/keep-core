package operator

import (
	"crypto/ecdsa"
	"io"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/pborman/uuid"
)

// PrivateKey represents peer's static key associated with an on-chain
// stake. It is used to authenticate the peer and for attributability (signing).
type PrivateKey = ecdsa.PrivateKey

// PublicKey represents peer's static key associated with an on-chain
// stake. It is used to authenticate the peer and for attributability
// (verification).
type PublicKey = ecdsa.PublicKey

// GenerateKeyPair generates a new, random static key based on
// secp256k1 ethereum curve.
func GenerateKeyPair(rand io.Reader) (*PrivateKey, *PublicKey, error) {
	id := uuid.NewRandom()

	ecdsaKey, err := ecdsa.GenerateKey(secp256k1.S256(), rand)
	if err != nil {
		return nil, nil, err
	}

	ethereumKey := &keystore.Key{
		Id:         id,
		Address:    crypto.PubkeyToAddress(ecdsaKey.PublicKey),
		PrivateKey: ecdsaKey,
	}

	privateKey, publicKey := EthereumKeyToStaticKey(ethereumKey)
	return privateKey, publicKey, nil
}

// EthereumKeyToStaticKey transforms a `go-ethereum`-based ECDSA key into the
// format supported by all packages used in keep-core.
func EthereumKeyToStaticKey(ethereumKey *keystore.Key) (*PrivateKey, *PublicKey) {
	privKey := ethereumKey.PrivateKey
	return (*PrivateKey)(privKey), (*PublicKey)(&privKey.PublicKey)
}
