package net

import (
	"fmt"
	"io"

	crand "crypto/rand"
	mrand "math/rand"

	ci "github.com/libp2p/go-libp2p-crypto"
	peer "github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"
)

type Identity interface {
	ID() peer.ID
	AddIdentityToStore() (pstore.Peerstore, error)
	PubKey() ci.PubKey
	PubKeyFromID(peer.ID) (ci.PubKey, error)
}

type PeerIdentity struct {
	privKey ci.PrivKey
}

func (pi *PeerIdentity) ID() peer.ID {
	return pubKeyToID(pi.privKey.GetPublic())
}

func pubKeyToID(pub ci.PubKey) peer.ID {
	// From go-libp2p-peer: PKI-based identities for libp2p
	pid, err := peer.IDFromPublicKey(pub)
	if err != nil {
		panic(fmt.Sprintf("Failed to generate valid libp2p identity with err: %v", err))
	}
	return pid
}

func (pi *PeerIdentity) KeyPair() (ci.PrivKey, ci.PubKey) {
	return pi.privKey, pi.privKey.GetPublic()
}

func (pi *PeerIdentity) AddIdentityToStore() (pstore.Peerstore, error) {
	ps := pstore.NewPeerstore()
	// HACK: see github.com/rargulati/go-libp2p-crypto for fix
	if err := ps.AddPrivKey(pi.ID(), pi.privKey); err != nil {
		return nil, fmt.Errorf("failed to add PrivateKey with error %s", err)
	}
	if err := ps.AddPubKey(pi.ID(), pi.PubKey()); err != nil {
		return nil, fmt.Errorf("failed to add PubKey with error %s", err)
	}
	return ps, nil
}

func (pi *PeerIdentity) PubKey() ci.PubKey {
	return pi.privKey.GetPublic()
}

func (pi *PeerIdentity) PubKeyFromID(peer.ID) (ci.PubKey, error) {
	return pi.ID().ExtractPublicKey()
}

func LoadOrGenerateIdentity(randseed int64, filePath string) (Identity, error) {
	if filePath != "" {
		// TODO: unmarshal and build out PKI
		// TODO: ensure this is associated with some staking address
	}
	if randseed != 0 {
		return generateDeterministicIdentity(randseed)
	}
	return generateIdentity()
}

// generateIdentity generates a public/private-key pair
// (using the libp2p/crypto wrapper for golang/crypto) provided a reader.
// Use randseed for deterministic IDs, otherwise we'll use cryptographically secure psuedorandomness.
func generateIdentity() (Identity, error) {
	var r io.Reader
	r = crand.Reader

	priv, _, err := ci.GenerateKeyPairWithReader(ci.Ed25519, 2048, r)
	if err != nil {
		return nil, err
	}

	return &PeerIdentity{privKey: priv}, nil
}

func generateDeterministicIdentity(randseed int64) (Identity, error) {
	var r io.Reader
	r = mrand.New(mrand.NewSource(randseed))

	priv, _, err := ci.GenerateKeyPairWithReader(ci.Ed25519, 2048, r)
	if err != nil {
		return nil, err
	}

	return &PeerIdentity{privKey: priv}, nil
}
