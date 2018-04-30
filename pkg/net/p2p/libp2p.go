package p2p

import (
	"fmt"
	"io"

	crand "crypto/rand"
	mrand "math/rand"

	ci "github.com/libp2p/go-libp2p-crypto"
	peer "github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"

	"github.com/keep-network/keep-core/pkg/net"
)

type peerIdentity struct {
	privKey ci.PrivKey
}

func (pi *peerIdentity) ProviderName() string {
	return "libp2p"
}

func (pi *peerIdentity) ID() string {
	return pubKeyToID(pi.privKey.GetPublic())
}

func pubKeyToID(pub ci.PubKey) string {
	pid, err := peer.IDFromPublicKey(pub)
	if err != nil {
		panic(fmt.Sprintf("Failed to generate valid libp2p identity with err: %v", err))
	}
	return string(pid)
}

func (pi *peerIdentity) PubKeyFromID(id string) ([]byte, error) {
	key, err := peer.ID(id).ExtractPublicKey()
	if err != nil {
		return nil, err
	}
	return key.Bytes()
}

// AddIdentityToStore takes a peerIdentity and notifies the addressbook of the
// existance of a new client joining the network.
func AddIdentityToStore(pi *peerIdentity) (pstore.Peerstore, error) {
	// TODO: investigate a generic store interface that gives us a unified interface
	// to our address book (peerstore in libp2p) from secure storage (dht)
	ps := pstore.NewPeerstore()

	id := peer.ID(pi.ID())

	if err := ps.AddPrivKey(id, pi.privKey); err != nil {
		return nil, fmt.Errorf("failed to add PrivateKey with error %s", err)
	}
	if err := ps.AddPubKey(id, pi.privKey.GetPublic()); err != nil {
		return nil, fmt.Errorf("failed to add PubKey with error %s", err)
	}
	return ps, nil
}

// LoadOrGenerateIdentity allows a client to provide or generate an Identity that
// will be used to reference the client in the peer-to-peer network.
func LoadOrGenerateIdentity(randseed int64, filePath string) (net.Identity, error) {
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
func generateIdentity() (net.Identity, error) {
	var r io.Reader
	r = crand.Reader

	priv, _, err := ci.GenerateKeyPairWithReader(ci.Ed25519, 2048, r)
	if err != nil {
		return nil, err
	}

	return &peerIdentity{privKey: priv}, nil
}

func generateDeterministicIdentity(randseed int64) (net.Identity, error) {
	var r io.Reader
	r = mrand.New(mrand.NewSource(randseed))

	priv, _, err := ci.GenerateKeyPairWithReader(ci.Ed25519, 2048, r)
	if err != nil {
		return nil, err
	}

	return &peerIdentity{privKey: priv}, nil
}
