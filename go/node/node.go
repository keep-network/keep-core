package Node

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"

	floodsub "github.com/libp2p/go-floodsub"
	ci "github.com/libp2p/go-libp2p-crypto"
	host "github.com/libp2p/go-libp2p-host"
	peer "github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	swarm "github.com/libp2p/go-libp2p-swarm"
	bhost "github.com/libp2p/go-libp2p/p2p/host/basic"
	ma "github.com/multiformats/go-multiaddr"
)

// A node is the initialized Keep client waiting to join a group
type Node struct {
	// Self
	Identity *Identity

	PeerHost  host.Host
	Bootstrap []string // bootstrap peer addrs

	PeerStore pstore.Peerstore

	Floodsub *floodsub.PubSub

	// Use to detect node shutdowns
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
	// var n *Node
	n := &Node{
		Identity: &Identity{},
	}
	//TODO: allow the user to supply
	priv, pub, err := generatePKI()
	if err != nil {
		panic(fmt.Sprintf("Failed to generate valid key material with err: %v", err))
	}

	// From go-libp2p-peer: PKI-based identities for libp2p
	pid, err := peer.IDFromEd25519PublicKey(pub)
	if err != nil {
		panic(fmt.Sprintf("Failed to generate valid libp2p identity with err: %v", err))
	}

	n.Identity.PeerID, n.Identity.PrivKey = pid, priv
	// Ensure that other members in our broadcast channel can identify us
	n.PeerStore = addToPeerStore(pid, priv, pub)
	// The context governs the lifetime of the libp2p node
	n.ctx = ctx

	if err := n.Start(); err != nil {
		panic(fmt.Sprintf("Failed to start Node process with err: %v", err))
	}

	return n
}

func (n *Node) Start() error {
	// TODO: flesh out how we connect to libp2p
	if n.PeerHost != nil {
		return fmt.Errorf("already online")
	}
	// TODO: attach a muxer to a connection
	// TODO: figure out go-libp2p-interface-pnet.Protector and go-libp2p-pnet.NewProtector - later
	peerhost, err := buildPeerHost(n.ctx, n.Identity.PeerID, n.PeerStore)
	if err != nil {
		return err
	}
	n.PeerHost = peerhost

	// Ok, now we're ready to listen
	// TODO: listen to more addresses, flesh this out
	listen, _ := ma.NewMultiaddr(fmt.Sprint("/ip4/127.0.0.1/tcp/8080"))
	if err := n.PeerHost.Network().Listen([]ma.Multiaddr{listen}...); err != nil {
		return err
	}
	// TODO: implement a standard and functional logger
	log.Printf("Listening at: %s\n", listen)

	ps, err := floodsub.NewFloodSub(n.ctx, n.PeerHost)
	if err != nil {
		return err
	}
	n.Floodsub = ps

	return nil
}

func buildPeerHost(ctx context.Context, pid peer.ID, ps pstore.Peerstore) (host.Host, error) {
	// TODO: use NewSwarmWithProtector
	// TODO: customize transport with config, for now use default in go-libp2p-swarm
	// Start without any addresses...
	swrm, err := swarm.NewSwarm(ctx, nil, pid, ps, nil)
	if err != nil {
		return nil, err
	}
	network := (*swarm.Network)(swrm)
	// TODO: use our own host, basic is used in projects and examples, but outdated
	opts := &bhost.HostOpts{}
	h, err := bhost.NewHost(ctx, network, opts)
	if err != nil {
		h.Close()
		return nil, err
	}
	// TODO: do we need to enable the circuit relay? if so, do it here
	return h, nil

}

// func buildRoutingService(ctx context.Context, h host.Host) error {
// }
