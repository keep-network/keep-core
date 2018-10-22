package libp2p

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"

	libp2pcrypto "github.com/libp2p/go-libp2p-crypto"
)

type ethereumPubKey struct {
	libp2pcrypto.PubKey

	delegate *ecdsa.PublicKey
}

func (epk *ethereumPubKey) Bytes() ([]byte, error) {
	pk := epk.delegate
	return elliptic.Marshal(pk.Curve, pk.X, pk.Y), nil
}

type ethereumPrivKey struct {
	libp2pcrypto.PrivKey
}

func identityKeyPair(ethereumKey *ecdsa.PrivateKey) (
	*ethereumPrivKey,
	*ethereumPubKey,
	error,
) {
	priv, pub, err := libp2pcrypto.ECDSAKeyPairFromKey(ethereumKey)
	if err != nil {
		return nil, nil, fmt.Errorf("could not create identity key [%v]", err)
	}

	return &ethereumPrivKey{priv}, &ethereumPubKey{pub, &ethereumKey.PublicKey}, nil
}
