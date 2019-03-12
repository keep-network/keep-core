package static

import (
	"crypto/ecdsa"
	"io"

	"github.com/btcsuite/btcd/btcec"
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

// GenerateStaticKeyPair generates a new, random static key based on
// secp256k1 ethereum curve.
func GenerateStaticKeyPair(rand io.Reader) (*PrivateKey, *PublicKey, error) {
	id := uuid.NewRandom()

	ecdsaKey, err := ecdsa.GenerateKey(secp256k1.S256(), rand)
	if err != nil {
		return nil, nil, err
	}

	key := &keystore.Key{
		Id:         id,
		Address:    crypto.PubkeyToAddress(ecdsaKey.PublicKey),
		PrivateKey: ecdsaKey,
	}

	privKey, pubKey := EthereumKeyToStaticKey(key)
	return privKey, pubKey, nil
}

// EthereumKeyToStaticKey transforms `go-ethereum`-based ECDSA key into the
// format supported by the network layer. Because all curve parameters of
// secp256k1 curve defined by `go-ethereum` and all curve parameters of
// secp256k1 curve defined by `btcsuite` used by `lipb2b` under the hood are
// identical, we can simply rewrite the private key.
func EthereumKeyToStaticKey(ethereumKey *keystore.Key) (*PrivateKey, *PublicKey) {
	privKey, pubKey := btcec.PrivKeyFromBytes(
		btcec.S256(), ethereumKey.PrivateKey.D.Bytes(),
	)

	return (*PrivateKey)(privKey), (*PublicKey)(pubKey)
}

// PubKeyToEthAddress transforms a static PubKey into an Ethereum account
// address, in a string format.
func PubKeyToEthAddress(publicKey *PublicKey) string {
	return crypto.PubkeyToAddress(*publicKey).String()
}
