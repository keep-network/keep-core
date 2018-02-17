// node is the base relay client initalized on startup
package node

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	mrand "math/rand"

	ci "github.com/libp2p/go-libp2p-crypto"
	peer "github.com/libp2p/go-libp2p-peer"
)

// A node is the initialized relay client waiting to join a group
type Node struct {
	// Self
	Identity *Identity

	Network *NetworkManager
	Groups  *GroupManager
	// Use to detect node shutdowns
	ctx context.Context
}

type Identity struct {
	PeerID  peer.ID
	PubKey  ci.PubKey
	PrivKey ci.PrivKey
}

// Only call once on init
func NewNode(ctx context.Context, port int, randseed int64) (*Node, error) {
	//TODO: allow the user to supply
	priv, pub, err := generatePKI(randseed)
	if err != nil {
		return nil, fmt.Errorf("Failed to generate valid key material with err: %v", err)
	}

	// From go-libp2p-peer: PKI-based identities for libp2p
	pid, err := peer.IDFromEd25519PublicKey(pub)
	if err != nil {
		return nil, fmt.Errorf("Failed to generate valid libp2p identity with err: %v", err)
	}

	n := &Node{Identity: &Identity{PeerID: pid, PrivKey: priv, PubKey: pub}}
	// The context governs the lifetime of the libp2p node
	n.ctx = ctx
	n.Network, err = NewNetworkManager(n.ctx, port, n.Identity.PeerID, n.Identity.PrivKey, n.Identity.PubKey)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func generatePKI(randseed int64) (ci.PrivKey, ci.PubKey, error) {
	// If the seed is zero, use real cryptographic randomness. Otherwise, use a
	// deterministic randomness source to make generated keys stay the same
	// across multiple runs
	var r io.Reader
	if randseed == 0 {
		r = rand.Reader
	} else {
		r = mrand.New(mrand.NewSource(randseed))
	}

	priv, pub, err := ci.GenerateKeyPairWithReader(ci.Ed25519, 2048, r)
	if err != nil {
		return nil, nil, err
	}
	return priv, pub, nil
}
