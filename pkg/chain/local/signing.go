package local

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/asn1"
	"fmt"
	"math/big"
)

type localSigning struct {
	operatorKey *ecdsa.PrivateKey
}

type ecdsaSignature struct {
	R, S *big.Int
}

func (ls *localSigning) PublicKey() []byte {
	publicKey := ls.operatorKey.PublicKey
	return elliptic.Marshal(publicKey.Curve, publicKey.X, publicKey.Y)
}

func (ls *localSigning) Sign(message []byte) ([]byte, error) {
	hash := sha256.Sum256(message)

	r, s, err := ecdsa.Sign(rand.Reader, ls.operatorKey, hash[:])
	if err != nil {
		return nil, err
	}

	return asn1.Marshal(ecdsaSignature{r, s})
}

func (ls *localSigning) Verify(message []byte, signature []byte) (bool, error) {
	return verifySignature(message, signature, &ls.operatorKey.PublicKey)
}

func (ls *localSigning) VerifyWithPublicKey(
	message []byte,
	signature []byte,
	publicKey []byte,
) (bool, error) {
	unmarshalledPubKey, err := unmarshalPublicKey(
		publicKey,
		ls.operatorKey.Curve,
	)
	if err != nil {
		return false, err
	}

	return verifySignature(message, signature, unmarshalledPubKey)
}

func verifySignature(
	message []byte,
	signature []byte,
	publicKey *ecdsa.PublicKey,
) (bool, error) {
	hash := sha256.Sum256(message)

	sig := &ecdsaSignature{}
	_, err := asn1.Unmarshal(signature, sig)
	if err != nil {
		return false, err
	}

	return ecdsa.Verify(publicKey, hash[:], sig.R, sig.S), nil
}

func unmarshalPublicKey(
	bytes []byte,
	curve elliptic.Curve,
) (*ecdsa.PublicKey, error) {
	x, y := elliptic.Unmarshal(curve, bytes)
	if x == nil {
		return nil, fmt.Errorf(
			"invalid public key bytes",
		)
	}
	ecdsaPublicKey := &ecdsa.PublicKey{Curve: curve, X: x, Y: y}
	return (*ecdsa.PublicKey)(ecdsaPublicKey), nil
}

func (ls *localSigning) PublicKeyToAddress(publicKey ecdsa.PublicKey) []byte {
	return elliptic.Marshal(publicKey.Curve, publicKey.X, publicKey.Y)
}

func (ls *localSigning) PublicKeyBytesToAddress(publicKey []byte) []byte {
	return publicKey
}
