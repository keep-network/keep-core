package libp2p

import (
	"crypto/rand"
	"fmt"

	libp2pcrypto "github.com/libp2p/go-libp2p-crypto"
	peer "github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"
)

// identity represents a group member's network level identity. It
// implements the net.TransportIdentifier interface. A valid group member will
// generate or provide a keypair, which will correspond to a network ID.
//
// Consumers of the net package require an ID to register with protocol level
// IDs, as well as a public key for authentication.
type identity struct {
	id      peer.ID
	pubKey  libp2pcrypto.PubKey
	privKey libp2pcrypto.PrivKey
}

func (i *identity) ProviderName() string {
	return "libp2p"
}

func pubKeyToIdentifier(pub libp2pcrypto.PubKey) *identity {
	pid, err := peer.IDFromPublicKey(pub)
	if err != nil {
		panic(fmt.Sprintf("Failed to generate valid libp2p identity with err: %v", err))
	}
	return &identity{id: pid}
}

// AddIdentityToStore takes an identity and notifies the addressbook of the
// existance of a new client joining the network.
func addIdentityToStore(i *identity) (pstore.Peerstore, error) {
	// TODO: investigate a generic store interface that gives us a unified interface
	// to our address book (peerstore in libp2p) from secure storage (dht)
	peerstore := pstore.NewPeerstore()

	if err := peerstore.AddPrivKey(i.id, i.privKey); err != nil {
		return nil, fmt.Errorf("failed to add PrivateKey to store with error %s", err)
	}
	if err := peerstore.AddPubKey(i.id, i.pubKey); err != nil {
		return nil, fmt.Errorf("failed to add PubKey to store with error %s", err)
	}
	return peerstore, nil
}

// generateIdentity generates a public/private-key pair (using the libp2p/crypto
// wrapper for golang/crypto).
func generateIdentity() (*identity, error) {
	privKey, pubKey, err := libp2pcrypto.GenerateKeyPairWithReader(libp2pcrypto.Ed25519, 2048, rand.Reader)
	if err != nil {
		return nil, err
	}

	return &identity{privKey: privKey, pubKey: pubKey}, nil
}
