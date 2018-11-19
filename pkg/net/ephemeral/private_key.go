package ephemeral

import (
	"fmt"

	"github.com/btcsuite/btcd/btcec"
)

// PrivateEcdsaKey is an ephemeral private elliptic curve key.
type PrivateEcdsaKey btcec.PrivateKey

// PublicEcdsaKey is an ephemeral public elliptic curve key.
type PublicEcdsaKey btcec.PublicKey

func curve() *btcec.KoblitzCurve {
	return btcec.S256()
}

// GenerateEphemeralKeypair generates a pair of public and private ECDSA keys
// that can be used as an input for ECDH.
func GenerateEphemeralKeypair() (*PrivateEcdsaKey, *PublicEcdsaKey, error) {
	ecdsaKey, err := btcec.NewPrivateKey(curve())
	if err != nil {
		return nil, nil, fmt.Errorf(
			"could not generate new ephemeral keypair [%v]",
			err,
		)
	}

	return (*PrivateEcdsaKey)(ecdsaKey), (*PublicEcdsaKey)(&ecdsaKey.PublicKey), nil
}

// UnmarshalPrivateKey turns a slice of bytes into a `PrivateEcdsaKey`.
func UnmarshalPrivateKey(bytes []byte) *PrivateEcdsaKey {
	priv, _ := btcec.PrivKeyFromBytes(curve(), bytes)
	return (*PrivateEcdsaKey)(priv)
}

// UnmarshalPublicKey turns a slice of bytes into a `PublicEcdsaKey`.
func UnmarshalPublicKey(bytes []byte) (*PublicEcdsaKey, error) {
	pubKey, err := btcec.ParsePubKey(bytes, curve())
	if err != nil {
		return nil, fmt.Errorf("could not parse ephemeral public key [%v]", err)
	}

	return (*PublicEcdsaKey)(pubKey), nil
}

// Marshal turns a `PrivateEcdsaKey` into a slice of bytes.
func (privk *PrivateEcdsaKey) Marshal() []byte {
	return (*btcec.PrivateKey)(privk).Serialize()

}

// Marshal turns a `PublicEcdsaKey` into a slice of bytes.
func (pubk *PublicEcdsaKey) Marshal() []byte {
	return (*btcec.PublicKey)(pubk).SerializeCompressed()

}
