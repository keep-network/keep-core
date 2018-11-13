package ephemeral

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/x509"
	"fmt"
	"io"
)

// SymmetricKey is a session-scoped ECDH symmetric key.
type SymmetricKey interface {
	Encrypt([]byte) ([]byte, error)
}

// PrivateKey is a session-scoped private elliptic curve key.
type PrivateKey struct {
	ecdsa.PrivateKey
}

// PublicKey is a session-scoped public elliptic curve key.
type PublicKey struct {
	ecdsa.PublicKey
}

func GenerateEphemeralKeypair(rand io.Reader) (*PrivateKey, *PublicKey, error) {
	ethCurve := elliptic.P256()

	ecdsaKey, err := ecdsa.GenerateKey(ethCurve, rand)
	if err != nil {
		return nil, nil, err
	}

	privKey := &PrivateKey{*ecdsaKey}
	pubKey := &PublicKey{ecdsaKey.PublicKey}

	return privKey, pubKey, nil
}

func (pk *PrivateKey) Decrypt(message []byte) ([]byte, error) {
	return nil, nil
}

func (pk *PrivateKey) Ecdh(publicKey *PublicKey) SymmetricKey {
	return nil
}

func (pk *PrivateKey) Marshal() ([]byte, error) {
	return x509.MarshalECPrivateKey(&pk.PrivateKey)
}

func (pk *PrivateKey) Unmarshal(bytes []byte) error {
	ecdsaPrivKey, err := x509.ParseECPrivateKey(bytes)
	if err != nil {
		return fmt.Errorf("could not parse ephemeral private key [%v]", err)
	}

	pk.PrivateKey = *ecdsaPrivKey
	return nil
}

func (pk *PublicKey) Marshal() ([]byte, error) {
	return x509.MarshalPKIXPublicKey(&pk.PublicKey)
}

func (pk *PublicKey) Unmarshal(bytes []byte) error {
	pubKey, err := x509.ParsePKIXPublicKey(bytes)
	if err != nil {
		return fmt.Errorf("could not parse ephemeral public key [%v]", err)
	}

	ecdsaPubKey, ok := pubKey.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("unexpected type of ephemeral public key")
	}

	pk.PublicKey = *ecdsaPubKey
	return nil
}
