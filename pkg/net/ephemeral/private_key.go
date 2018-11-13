package ephemeral

import (
	"fmt"

	"github.com/btcsuite/btcd/btcec"
)

// PrivateEcdsaKey is a session-scoped private elliptic curve key.
type PrivateEcdsaKey btcec.PrivateKey

// PublicEcdsaKey is a session-scoped public elliptic curve key.
type PublicEcdsaKey btcec.PublicKey

func curve() *btcec.KoblitzCurve {
	return btcec.S256()
}

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

func UnmarshalPrivateKey(bytes []byte) *PrivateEcdsaKey {
	priv, _ := btcec.PrivKeyFromBytes(curve(), bytes)
	return (*PrivateEcdsaKey)(priv)
}

func UnmarshalPublicKey(bytes []byte) (*PublicEcdsaKey, error) {
	pubKey, err := btcec.ParsePubKey(bytes, curve())
	if err != nil {
		return nil, fmt.Errorf("could not parse ephemeral public key [%v]", err)
	}

	return (*PublicEcdsaKey)(pubKey), nil
}

func (pk *PrivateEcdsaKey) Marshal() []byte {
	return pk.toBtcec().Serialize()
}

func (pk *PrivateEcdsaKey) toBtcec() *btcec.PrivateKey {
	return (*btcec.PrivateKey)(pk)
}

func (pk *PublicEcdsaKey) Marshal() []byte {
	return pk.toBtcec().SerializeCompressed()
}

func (pk *PublicEcdsaKey) toBtcec() *btcec.PublicKey {
	return (*btcec.PublicKey)(pk)
}
