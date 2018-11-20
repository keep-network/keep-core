package ephemeral

import (
	"fmt"

	"github.com/btcsuite/btcd/btcec"
)

// EphemeralPrivateKey is an ephemeral private elliptic curve key.
type EphemeralPrivateKey btcec.PrivateKey

// EphemeralPublicKey is an ephemeral public elliptic curve key.
type EphemeralPublicKey btcec.PublicKey

func curve() *btcec.KoblitzCurve {
	return btcec.S256()
}

// GenerateEphemeralKeypair generates a pair of public and private ECDSA keys
// that can be used as an input for ECDH.
func GenerateEphemeralKeypair() (*EphemeralPrivateKey, *EphemeralPublicKey, error) {
	ecdsaKey, err := btcec.NewPrivateKey(curve())
	if err != nil {
		return nil, nil, fmt.Errorf(
			"could not generate new ephemeral keypair [%v]",
			err,
		)
	}

	return (*EphemeralPrivateKey)(ecdsaKey), (*EphemeralPublicKey)(&ecdsaKey.PublicKey), nil
}

// UnmarshalPrivateKey turns a slice of bytes into a `PrivateEcdsaKey`.
func UnmarshalPrivateKey(bytes []byte) *EphemeralPrivateKey {
	priv, _ := btcec.PrivKeyFromBytes(curve(), bytes)
	return (*EphemeralPrivateKey)(priv)
}

// UnmarshalPublicKey turns a slice of bytes into a `PublicEcdsaKey`.
func UnmarshalPublicKey(bytes []byte) (*EphemeralPublicKey, error) {
	pubKey, err := btcec.ParsePubKey(bytes, curve())
	if err != nil {
		return nil, fmt.Errorf("could not parse ephemeral public key [%v]", err)
	}

	return (*EphemeralPublicKey)(pubKey), nil
}

// Marshal turns a `PrivateEcdsaKey` into a slice of bytes.
func (pk *EphemeralPrivateKey) Marshal() []byte {
	return pk.toBtcec().Serialize()
}

func (pk *EphemeralPrivateKey) toBtcec() *btcec.PrivateKey {
	return (*btcec.PrivateKey)(pk)
}

// Marshal turns a `PublicEcdsaKey` into a slice of bytes.
func (pk *EphemeralPublicKey) Marshal() []byte {
	return pk.toBtcec().SerializeCompressed()
}

func (pk *EphemeralPublicKey) toBtcec() *btcec.PublicKey {
	return (*btcec.PublicKey)(pk)
}
