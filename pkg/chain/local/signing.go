package local

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/asn1"
	"math/big"
)

type localSigning struct {
	operatorKey *ecdsa.PrivateKey
}

type ecdsaSignature struct {
	R, S *big.Int
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
	hash := sha256.Sum256(message)

	sig := &ecdsaSignature{}
	_, err := asn1.Unmarshal(signature, sig)
	if err != nil {
		return false, err
	}

	return ecdsa.Verify(&ls.operatorKey.PublicKey, hash[:], sig.R, sig.S), nil
}
