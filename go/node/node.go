package Node

import (
	"context"
	"crypto/rand"
	"fmt"

	floodsub "github.com/libp2p/go-floodsub"
	ci "github.com/libp2p/go-libp2p-crypto"
	host "github.com/libp2p/go-libp2p-host"
	peer "github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	swarm "github.com/libp2p/go-libp2p-swarm"
	bhost "github.com/libp2p/go-libp2p/p2p/host/basic"
)

// A node is the initialized Keep client waiting to join a group
type Node struct {
	// Self
	Identity Identity

	PeerHost  host.Host
	Bootstrap []string // bootstrap peer addrs

	PeerStore pstore.Peerstore

	Floodsub *floodsub.PubSub

	// Need a method of detecting shutdowns
	// Maybe use ctx.Cancel() and ctx.Done() until then...
	ctx context.Context
}

type Identity struct {
	PeerID  peer.ID
	PrivKey ci.PrivKey
}

func addToPeerStore(pid peer.ID, priv ci.PrivKey, pub ci.PubKey) pstore.Peerstore {
	ps := pstore.NewPeerstore()
	ps.AddPrivKey(pid, priv)
	ps.AddPubKey(pid, pub)
	return ps
}

func generatePKI() (ci.PrivKey, ci.PubKey, error) {
	// TODO: deterministic randomness for tests
	r := rand.Reader

	priv, pub, err := ci.GenerateKeyPairWithReader(ci.Ed25519, 2048, r)
	if err != nil {
		return nil, nil, err
	}
	return priv, pub, nil
}

// Only call once on init
func NewNode(ctx context.Context) *Node {
	var n *Node
	//TODO: allow the user to supply
	priv, pub, err := generatePKI()
	if err != nil {
		panic("Failed to generate valid key material")
	}

	// From go-libp2p-peer: PKI-based identities for libp2p
	pid, err := peer.IDFromEd25519PublicKey(pub)
	if err != nil {
		panic("Failed to generate valid libp2p identity")
	}

	n.Identity.PeerID, n.Identity.PrivKey = pid, priv
	// Ensure that other members in our broadcast channel can identify us
	n.PeerStore = addToPeerStore(pid, priv, pub)
	// The context governs the lifetime of the libp2p node
	n.ctx = context.Background()

	if err := n.Start(); err != nil {
		panic("Failed to start Node process")
	}

	return n
}

func (n *Node) Start(ctx context.Context) error {
	// TODO: flesh out how we connect to libp2p
	// listen, _ := ma.NewMultiaddr(fmt.Sprint("/ip4/127.0.0.1/tcp/80"))
	if n.PeerHost != nil {
		return fmt.Errorf("already online")
	}
	// TODO: init a new transport - asap
	// TODO: attach a muxer to a connection asap
	// TODO: figure out go-libp2p-interface-pnet.Protector and go-libp2p-pnet.NewProtector - later

	return nil
}

func buildPeerHost(ctx context.Context, pid peer.ID, ps pstore.Peerstore) (host.Host, error) {
	// TODO: use NewSwarmWithProtector
	// Start without any addresses...
	swrm, err := swarm.NewNetwork(ctx, nil, pid, ps, nil)
	if err != nil {
		return nil, err
	}
	// network := (*swarm.Network)(swrm)
	h, err := bhost.NewHost(ctx, swrm, nil)
	if err != nil {
		h.Close()
		return nil, err
	}
	// TODO: do I need a circuit relay?
	return h, nil

}
