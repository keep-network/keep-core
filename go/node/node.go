// node is the base relay client initalized on startup
package node

import (
	"context"
	crand "crypto/rand"
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

	// onChain    chan *pb.ChainMessage
	// groupDKG   chan *pb.DKGMessage
	// groupRelay chan *pb.RelayMessage

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

	n.Groups = NewGroupManager(n.Network.Sub, n.Network.PeerHost, n.Network.Routing, n.Network.DHT)

	return n, nil
}

// generatePKI generates a public/private-key pair
// (using the libp2p/crypto wrapper for golang/crypto) provided a reader.
// Use randseed for deterministic IDs, otherwise we'll use cryptographically secure psuedorandomness.
func generatePKI(randseed int64) (ci.PrivKey, ci.PubKey, error) {
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

// EventLoop handles all inputs arriving on channels
// func EventLoop(ctx context.Context) {
// 	for {
// 		select {
// 		case cMsg <- onChain:
// 			return
// 		case dkgMsg <- groupDKGChannel:
// 			return
// 		case relayMsg <- groupRelayChain:
// 			return
// 		case <-ctx.Done():
// 			return
// 		}
// 	}
// }
