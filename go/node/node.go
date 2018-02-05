package Node

import (
	"context"
	"crypto/rand"
	"fmt"

	crypto "github.com/libp2p/go-libp2p-crypto"
	host "github.com/libp2p/go-libp2p-host"
	peer "github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	swarm "github.com/libp2p/go-libp2p-swarm"
	bhost "github.com/libp2p/go-libp2p/p2p/basichost"
	ma "github.com/multiformats/go-multiaddr"
)

type Node struct {
	host.Host //libp2p identity
}

func addToPeerStore(pid peer.ID, priv crypto.PrivKey, pub crypto.PubKey) pstore.Peerstore {
	ps := pstore.NewPeerstore()
	ps.AddPrivKey(pid, priv)
	ps.AddPubKey(pid, pub)
	return ps
}

func generatePKI() (crypto.PrivKey, crypto.PubKey, error) {
	// TODO: deterministic randomness for tests
	r := rand.Reader

	priv, pub, err := crypto.GenerateKeyPairWithReader(crypto.Ed25519, 2048, r)
	if err != nil {
		return nil, nil, err
	}
	return priv, pub, nil
}

func NewNode() *Node {
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

	// Ensure that other members in our broadcast channel can identify us
	ps := addToPeerStore(pid, priv, pub)
	ctx := context.Background()

	// TODO: redo how we connect to libp2p
	listen, _ := ma.NewMultiaddr(fmt.Sprint("/ip4/127.0.0.1/tcp/80"))
	n, _ := swarm.NewNetwork(ctx, []ma.Multiaddr{listen}, pid, ps, nil)
	// FIXME: Easypath
	h := bhost.NewHost(ctx, n, nil)

	return &Node{Host: h}
}
