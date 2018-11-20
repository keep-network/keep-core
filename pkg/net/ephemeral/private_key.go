package ephemeral

import (
	"fmt"

	"github.com/btcsuite/btcd/btcec"
)

// PrivateKey is an ephemeral private elliptic curve key.
type PrivateKey btcec.PrivateKey

// PublicKey is an ephemeral public elliptic curve key.
type PublicKey btcec.PublicKey

func curve() *btcec.KoblitzCurve {
	return btcec.S256()
}

// GenerateKeypair generates a pair of public and private ephemeral keys
// that can be used as an input for ECDH.
func GenerateKeypair() (*PrivateKey, *PublicKey, error) {
	ecdsaKey, err := btcec.NewPrivateKey(curve())
	if err != nil {
		return nil, nil, fmt.Errorf(
			"could not generate new ephemeral keypair [%v]",
			err,
		)
	}

	return (*PrivateKey)(ecdsaKey), (*PublicKey)(&ecdsaKey.PublicKey), nil
}

// UnmarshalPrivateKey turns a slice of bytes into a `PrivateKey`.
func UnmarshalPrivateKey(bytes []byte) *PrivateKey {
	priv, _ := btcec.PrivKeyFromBytes(curve(), bytes)
	return (*PrivateKey)(priv)
}

// UnmarshalPublicKey turns a slice of bytes into a `PublicKey`.
func UnmarshalPublicKey(bytes []byte) (*PublicKey, error) {
	pubKey, err := btcec.ParsePubKey(bytes, curve())
	if err != nil {
		return nil, fmt.Errorf("could not parse ephemeral public key [%v]", err)
	}

	return (*PublicKey)(pubKey), nil
}

// Marshal turns a `PrivateKey` into a slice of bytes.
func (pk *PrivateKey) Marshal() []byte {
	return (*btcec.PrivateKey)(pk).Serialize()
}

// Marshal turns a `PublicKey` into a slice of bytes.
func (pk *PublicKey) Marshal() []byte {
	return (*btcec.PublicKey)(pk).SerializeCompressed()
}
