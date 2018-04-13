package net

import (
	"fmt"
	"io"

	ci "github.com/libp2p/go-libp2p-crypto"
	peer "github.com/libp2p/go-libp2p-peer"
)

type Identity interface {
	ID() peer.ID
	PubKey() ci.PubKey
	PubKey(peer.ID) (ci.PubKey, error)
}

type PeerIdentity struct {
	PrivKey ci.PrivKey
}

func (pi *PeerIdentity) ID() peer.ID {
	return pubKeyToID(pi.PrivKey.GetPublic())
}

func pubKeyToID(pk ci.PubKey) peer.ID {
	// From go-libp2p-peer: PKI-based identities for libp2p
	pid, err := peer.IDFromEd25519PublicKey(pub)
	if err != nil {
		return nil, fmt.Errorf("Failed to generate valid libp2p identity with err: %v", err)
	}
}

func (pi *PeerIdentity) PubKey() ci.PubKey {
	return pi.PrivKey.GetPublic()
}

func (pi *PeerIdentity) PubKey(peer.ID) (ci.PubKey, error) {
	return pi.ID().ExtractEd25519PublicKey()
}

func LoadOrGeneratePKI(filePath string) (*PeerIdentity, error) {
	if filePath != "" {
		// TODO: unmarshal and build out PKI
	}
	return generatePKI(0, "") // TODO: ensure this is associated with some staking address
}

// generatePKI generates a public/private-key pair
// (using the libp2p/crypto wrapper for golang/crypto) provided a reader.
// Use randseed for deterministic IDs, otherwise we'll use cryptographically secure psuedorandomness.
func generatePKI(randseed int64, keyType string) (ci.PrivKey, ci.PubKey, error) {
	var r io.Reader
	if randseed == 0 {
		r = crand.Reader
	} else {
		r = mrand.New(mrand.NewSource(randseed))
	}
	// TODO: explore if we use PublicKeyToCurve25519 (converts an Ed25519 public key into the curve25519)
	priv, pub, err := ci.GenerateKeyPairWithReader(ci.Ed25519, 2048, r)
	if err != nil {
		return nil, nil, err
	}
	return priv, pub, nil
}
