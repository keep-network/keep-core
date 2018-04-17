package net

import (
	"fmt"
	"io"

	ci "github.com/libp2p/go-libp2p-crypto"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	peer "github.com/rargulati/go-libp2p-peer"
)

type Identity interface {
	ID() peer.ID
	AddIdentityToStore() (pstore.Peerstore, error)
	PubKey() (ci.PubKey, error)
	PubKey(peer.ID) (ci.PubKey, error)
}

type PeerIdentity struct {
	privKey ci.PrivKey
}

func (pi *PeerIdentity) ID() peer.ID {
	return pubKeyToID(pi.privKey.GetPublic())
}

func pubKeyToID(pk ci.PubKey) peer.ID {
	// From go-libp2p-peer: PKI-based identities for libp2p
	pid, err := peer.IDFromEd25519PublicKey(pub)
	if err != nil {
		return nil, fmt.Errorf("Failed to generate valid libp2p identity with err: %v", err)
	}
}

func (pi *PeerIdentity) KeyPair() (ci.PrivKey, ci.PubKey) {
	return pi.privkey, pi.privKey.GetPublic()
}

func (pi *PeerIdentity) AddIdentityToStore() (*pstore.Peerstore, error) {
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

func (pi *PeerIdentity) PubKey(peer.ID) (ci.PubKey, error) {
	return pi.privKey.GetPublic()
}

func (pi *PeerIdentity) PubKey(peer.ID) (ci.PubKey, error) {
	return pi.ID().ExtractEd25519PublicKey()
}

func LoadOrGenerateIdentity(randseed int64, filePath string) (*Identity, error) {
	if filePath != "" {
		// TODO: unmarshal and build out PKI
		// TODO: ensure this is associated with some staking address
	}
	if randseed != 0 {
		return generateDeterministicKeypair(randseed)
	}
	return generateKeypair()
}

// generateIdentity generates a public/private-key pair
// (using the libp2p/crypto wrapper for golang/crypto) provided a reader.
// Use randseed for deterministic IDs, otherwise we'll use cryptographically secure psuedorandomness.
func generateIdentity() (*PeerIdentity, error) {
	var r io.Reader
	r = crand.Reader

	priv, _, err := ci.GenerateKeyPairWithReader(ci.Ed25519, 2048, r)
	if err != nil {
		return nil, err
	}

	return &PeerIdentity{privKey: priv}, nil
}

func generateDeterministicIdentity(randseed int64) (*PeerIdentity, error) {
	var r io.Reader
	r = mrand.New(mrand.NewSource(randseed))

	priv, _, err := ci.GenerateKeyPairWithReader(ci.Ed25519, 2048, r)
	if err != nil {
		return nil, err
	}

	return &PeerIdentity{privKey: priv}, nil
}
