package libp2p

import (
	"crypto/rand"
	"fmt"
	"io"

	ci "github.com/libp2p/go-libp2p-crypto"
	peer "github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"
)

// peerIdentifier represents a group member's network level identity. It
// implements the net.TransportIdentifier interface. A valid group member will
// generate or provide a keypair, which will correspond to a network ID.
// Consumers of the net package require an ID to register with protocol level IDs, as well as a public key for authentication.
type peerIdentifier struct {
	id peer.ID
	sk ci.PrivKey
}

func (p peerIdentifier) ProviderName() string {
	return "libp2p"
}

func pubKeyFromIdentifier(pi peerIdentifier) (ci.PubKey, error) {
	pid := peer.ID(pi.id)
	return pid.ExtractPublicKey()
}

func pubKeyToIdentifier(pub ci.PubKey) peerIdentifier {
	pid, err := peer.IDFromPublicKey(pub)
	if err != nil {
		panic(fmt.Sprintf("Failed to generate valid libp2p identity with err: %v", err))
	}
	return peerIdentifier{id: pid}
}

// AddIdentityToStore takes a peerIdentifier and notifies the addressbook of the
// existance of a new client joining the network.
func addIdentityToStore(pi peerIdentifier) (pstore.Peerstore, error) {
	// TODO: investigate a generic store interface that gives us a unified interface
	// to our address book (peerstore in libp2p) from secure storage (dht)
	ps := pstore.NewPeerstore()
	pid := peer.ID(pi.id)

	if err := ps.AddPrivKey(pid, pi.sk); err != nil {
		return nil, fmt.Errorf("failed to add PrivateKey with error %s", err)
	}
	if err := ps.AddPubKey(pid, pi.sk.GetPublic()); err != nil {
		return nil, fmt.Errorf("failed to add PubKey with error %s", err)
	}
	return ps, nil
}

// generateIdentity generates a public/private-key pair
// (using the libp2p/crypto wrapper for golang/crypto) provided a reader.
// Use randseed for deterministic IDs, otherwise we'll use cryptographically secure psuedorandomness.
func generateIdentity() (*peerIdentifier, error) {
	var r io.Reader
	r = rand.Reader

	priv, _, err := ci.GenerateKeyPairWithReader(ci.Ed25519, 2048, r)
	if err != nil {
		return nil, err
	}

	return &peerIdentifier{sk: priv}, nil
}
