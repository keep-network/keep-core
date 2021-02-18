//+build celo

package operator

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"

	"github.com/celo-org/celo-blockchain/accounts/keystore"
	"github.com/celo-org/celo-blockchain/crypto"
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
func GenerateKeyPair() (*PrivateKey, *PublicKey, error) {
	ecdsaKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, nil, err
	}

	return (*PrivateKey)(ecdsaKey), (*PublicKey)(&ecdsaKey.PublicKey), nil
}

// ChainKeyToOperatorKey transforms a `go-ethereum`-based ECDSA key into the
// format supported by all packages used in keep-core.
func ChainKeyToOperatorKey(chainKey *keystore.Key) (*PrivateKey, *PublicKey) {
	privKey := chainKey.PrivateKey
	return (*PrivateKey)(privKey), (*PublicKey)(&privKey.PublicKey)
}

// Marshal takes an operator's PublicKey and produces an uncompressed public key
// as a slice of bytes (as specified in ANSI X9.62).
func Marshal(publicKey *PublicKey) []byte {
	return elliptic.Marshal(publicKey.Curve, publicKey.X, publicKey.Y)
}

// Unmarshal takes raw bytes and produces an uncompressed, operator's PublicKey.
// Unmarshal assumes the PublicKey's curve is of type S256 as defined in geth.
func Unmarshal(data []byte) (*PublicKey, error) {
	x, y := elliptic.Unmarshal(crypto.S256(), data)
	if x == nil {
		return nil, fmt.Errorf(
			"incorrect public key bytes",
		)
	}
	ecdsaPublicKey := &ecdsa.PublicKey{Curve: crypto.S256(), X: x, Y: y}
	return (*PublicKey)(ecdsaPublicKey), nil
}

// PubkeyToAddress converts operator's PublicKey to an address in string format.
func PubkeyToAddress(publicKey PublicKey) string {
	return crypto.PubkeyToAddress(publicKey).String()
}
