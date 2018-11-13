package ephemeral

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/x509"
	"fmt"
	"io"
)

func GenerateEphemeralKeypair(rand io.Reader) (PrivateKey, PublicKey, error) {
	ethCurve := elliptic.P256()

	ecdsaKey, err := ecdsa.GenerateKey(ethCurve, rand)
	if err != nil {
		return nil, nil, err
	}

	privKey := &ecdsaPrivateEphemeralKey{ecdsaKey}
	pubKey := &ecdsaPublicEphemeralKey{&ecdsaKey.PublicKey}

	return privKey, pubKey, nil
}

func UnmarshalPublicKey(bytes []byte) (PublicKey, error) {
	pubKey, err := x509.ParsePKIXPublicKey(bytes)
	if err != nil {
		return nil, fmt.Errorf("could not parse ephemeral public key [%v]", err)
	}

	ecdsaPubKey, ok := pubKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("unexpected type of ephemeral public key")
	}

	return &ecdsaPublicEphemeralKey{ecdsaPubKey}, nil
}

func UnmarshalPrivateKey(bytes []byte) (PrivateKey, error) {
	ecdsaPrivKey, err := x509.ParseECPrivateKey(bytes)
	if err != nil {
		return nil, fmt.Errorf("could not parse ephemeral private key [%v]", err)
	}

	return &ecdsaPrivateEphemeralKey{ecdsaPrivKey}, nil
}

type ecdsaPublicEphemeralKey struct {
	ecdsaKey *ecdsa.PublicKey
}

func (epek *ecdsaPublicEphemeralKey) Marshal() ([]byte, error) {
	return x509.MarshalPKIXPublicKey(epek.ecdsaKey)
}

type ecdsaPrivateEphemeralKey struct {
	ecdsaKey *ecdsa.PrivateKey
}

func (epek *ecdsaPrivateEphemeralKey) Marshal() ([]byte, error) {
	return x509.MarshalECPrivateKey(epek.ecdsaKey)
}

func (epek *ecdsaPrivateEphemeralKey) Decrypt(message []byte) ([]byte, error) {
	return nil, nil
}

func (epek *ecdsaPrivateEphemeralKey) Ecdh(pubKey PublicKey) SymmetricKey {
	return nil
}
