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

// Implementation of the TransportIdentifier interface
type networkID peer.ID

func (ni networkID) ProviderName() string {
	return "libp2p"
}

// peerIdentity represents a group member's network level identity. A valid group
// member will generate or provide a keypair, which will correspond to a network
// ID. Consumers of the net package require an ID to register with protocol level
// ID's, as well as a public key for authentication.
type peerIdentity struct {
	privKey ci.PrivKey
}

func (pi *peerIdentity) ID() TransportIdentifier {
	return networkID(pubKeyToID(pi.privKey.GetPublic()))
}

func pubKeyToID(pub ci.PubKey) peer.ID {
	pid, err := peer.IDFromPublicKey(pub)
	if err != nil {
		panic(fmt.Sprintf("Failed to generate valid libp2p identity with err: %v", err))
	}
	return pid
}

func (pi *peerIdentity) PubKeyFromID(id TransportIdentifier) (ci.PubKey, error) {
	pid := peer.ID(id.(networkID))
	return pid.ExtractPublicKey()
}

// AddIdentityToStore takes a peerIdentity and notifies the addressbook of the
// existance of a new client joining the network.
func addIdentityToStore(pi *peerIdentity) (pstore.Peerstore, error) {
	// TODO: investigate a generic store interface that gives us a unified interface
	// to our address book (peerstore in libp2p) from secure storage (dht)
	ps := pstore.NewPeerstore()
	id := peer.ID(pi.ID().(networkID))

	if err := ps.AddPrivKey(id, pi.privKey); err != nil {
		return nil, fmt.Errorf("failed to add PrivateKey with error %s", err)
	}
	if err := ps.AddPubKey(id, pi.privKey.GetPublic()); err != nil {
		return nil, fmt.Errorf("failed to add PubKey with error %s", err)
	}
	return ps, nil
}

// loadOrGenerateIdentity allows a client to provide or generate an Identity that
// will be used to reference the client in the peer-to-peer network.
func loadOrGenerateIdentity(randseed int64, filePath string) (*peerIdentity, error) {
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
func generateIdentity() (*peerIdentity, error) {
	var r io.Reader
	r = crand.Reader

	priv, _, err := ci.GenerateKeyPairWithReader(ci.Ed25519, 2048, r)
	if err != nil {
		return nil, err
	}

	return &peerIdentity{privKey: priv}, nil
}

func generateDeterministicIdentity(randseed int64) (*peerIdentity, error) {
	var r io.Reader
	r = mrand.New(mrand.NewSource(randseed))

	priv, _, err := ci.GenerateKeyPairWithReader(ci.Ed25519, 2048, r)
	if err != nil {
		return nil, err
	}

	return &peerIdentity{privKey: priv}, nil
}
