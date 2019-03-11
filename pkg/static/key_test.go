package static

import (
	"crypto/ecdsa"
	"crypto/rand"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/pborman/uuid"
)

func TestStaticPubKeyToAddress(t *testing.T) {
	ethereumKey, err := generateEthereumKey()
	if err != nil {
		t.Fatal(err)
	}

	_, staticPublicKey := EthereumKeyToStaticKey(ethereumKey)

	ethAddress := crypto.PubkeyToAddress(ethereumKey.PrivateKey.PublicKey).String()

	staticKeyAddress := PubKeyToEthAddress(staticPublicKey)

	if ethAddress != staticKeyAddress {
		t.Errorf(
			"unexpected address\nexpected: %v\nactual: %v",
			ethAddress,
			staticKeyAddress,
		)
	}
}

func generateEthereumKey() (*keystore.Key, error) {
	ethCurve := secp256k1.S256()

	ethereumKey, err := ecdsa.GenerateKey(ethCurve, rand.Reader)
	if err != nil {
		return nil, err
	}

	id := uuid.NewRandom()

	return &keystore.Key{
		Id:         id,
		Address:    crypto.PubkeyToAddress(ethereumKey.PublicKey),
		PrivateKey: ethereumKey,
	}, nil
}
