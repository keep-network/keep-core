package node

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	mrand "math/rand"
	"time"

	dstore "github.com/ipfs/go-datastore"
	dssync "github.com/ipfs/go-datastore/sync"
	floodsub "github.com/libp2p/go-floodsub"
	ci "github.com/libp2p/go-libp2p-crypto"
	host "github.com/libp2p/go-libp2p-host"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	peer "github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	routing "github.com/libp2p/go-libp2p-routing"
	swarm "github.com/libp2p/go-libp2p-swarm"
	bhost "github.com/libp2p/go-libp2p/p2p/host/basic"
	rhost "github.com/libp2p/go-libp2p/p2p/host/routed"
	smux "github.com/libp2p/go-stream-muxer"
	ma "github.com/multiformats/go-multiaddr"
	msmux "github.com/whyrusleeping/go-smux-multistream"
	yamux "github.com/whyrusleeping/go-smux-yamux"
)

var (
	bootstrapPeers = []string{"/ip4/127.0.0.1/tcp/2701/ipfs/QmexAnfpHrhMmAC5UNQVS8iBuUUgDrMbMY17Cck2gKrqeX", "/ip4/127.0.0.1/tcp/2702/ipfs/Qmd3wzD2HWA95ZAs214VxnckwkwM4GHJyC6whKUCNQhNvW"}
)

// A node is the initialized relay client waiting to join a group
type Node struct {
	// Self
	Identity *Identity

	PeerHost host.Host

	PeerStore pstore.Peerstore

	Floodsub *floodsub.PubSub
	Routing  routing.IpfsRouting

	// Use to detect node shutdowns
	ctx context.Context
}

type Identity struct {
	PeerID  peer.ID
	PrivKey ci.PrivKey
}

// Only call once on init
func NewNode(ctx context.Context, port int, randseed int64) (*Node, error) {
	n := &Node{
		Identity: &Identity{},
	}

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

	n.Identity.PeerID, n.Identity.PrivKey = pid, priv
	// Ensure that other members in our broadcast channel can identify us
	n.PeerStore = addToPeerStore(pid, priv, pub)
	// The context governs the lifetime of the libp2p node
	n.ctx = ctx

	if err := n.start(port); err != nil {
		return nil, fmt.Errorf("Failed to start Node process with err: %v", err)
	}

	dhtRouting := dht.NewDHTClient(n.ctx, n.PeerHost, dssync.MutexWrap(dstore.NewMapDatastore()))
	n.Routing = dhtRouting
	n.PeerHost = rhost.Wrap(n.PeerHost, dhtRouting)

	if err := n.bootstrap(); err != nil {
		return nil, fmt.Errorf("Failed to bootstrap nodes with err: %v", err)
	}

	return n, nil
}

func addToPeerStore(pid peer.ID, priv ci.PrivKey, pub ci.PubKey) pstore.Peerstore {
	ps := pstore.NewPeerstore()
	ps.AddPrivKey(pid, priv)
	ps.AddPubKey(pid, pub)
	return ps
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

func (n *Node) start(port int) error {
	// TODO: flesh out how we connect to libp2p
	if n.PeerHost != nil {
		return fmt.Errorf("already online")
	}
	// TODO: attach a muxer to a connection
	// TODO: figure out go-libp2p-interface-pnet.Protector and go-libp2p-pnet.NewProtector - later
	listen, err := ma.NewMultiaddr(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", port))
	if err != nil {
		return err
	}
	peerhost, err := buildPeerHost(n.ctx, []ma.Multiaddr{listen}, n.Identity.PeerID, n.PeerStore)
	if err != nil {
		return err
	}
	n.PeerHost = peerhost

	// Ok, now we're ready to listen
	// TODO: listen to more addresses, flesh this out
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

func makeSmuxTransport() smux.Transport {
	mstpt := msmux.NewBlankTransport()

	ymxtpt := &yamux.Transport{
		AcceptBacklog:          512,
		ConnectionWriteTimeout: time.Second * 10,
		KeepAliveInterval:      time.Second * 30,
		EnableKeepAlive:        true,
		MaxStreamWindowSize:    uint32(1024 * 512),
		LogOutput:              ioutil.Discard,
	}

	mstpt.AddTransport("/yamux/1.0.0", ymxtpt)
	return mstpt
}

func buildPeerHost(ctx context.Context, listenAddrs []ma.Multiaddr, pid peer.ID, ps pstore.Peerstore) (host.Host, error) {
	// Set up stream multiplexer
	tpt := makeSmuxTransport()

	// TODO: use NewSwarmWithProtector
	swrm, err := swarm.NewSwarmWithProtector(ctx, listenAddrs, pid, ps, nil, tpt, nil)
	if err != nil {
		return nil, err
	}

	network := (*swarm.Network)(swrm)
	// TODO: use our own host, I'm unsure about the utility of basic
	opts := &bhost.HostOpts{NATManager: bhost.NewNATManager(network)}
	h, err := bhost.NewHost(ctx, network, opts)
	if err != nil {
		h.Close()
		return nil, err
	}
	// TODO: do we need to enable the circuit relay? if so, do it here
	return h, nil

}

// lifted from github.com/keep-network/go-experiments
func (n *Node) bootstrap() error {
	log.Println("Bootstrapping peers...")
	for _, p := range bootstrapPeers {
		// The following code extracts target's the peer ID from the
		// given multiaddress
		ipfsaddr, err := ma.NewMultiaddr(p)
		if err != nil {
			log.Fatalln(err)
		}

		pid, err := ipfsaddr.ValueForProtocol(ma.P_IPFS)
		if err != nil {
			log.Fatalln(err)
		}

		peerid, err := peer.IDB58Decode(pid)
		if err != nil {
			log.Fatalln(err)
		}

		// Decapsulate the /ipfs/<peerID> part from the target
		// /ip4/<a.b.c.d>/ipfs/<peer> becomes /ip4/<a.b.c.d>
		targetPeerAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", peer.IDB58Encode(peerid)))
		targetAddr := ipfsaddr.Decapsulate(targetPeerAddr)
		if n.PeerHost.ID().String() != peerid.String() {
			// We have a peer ID and a targetAddr so we add it to the peerstore
			// so LibP2P knows how to contact it
			n.PeerHost.Peerstore().AddAddr(peerid, targetAddr, pstore.PermanentAddrTTL)
			n.PeerHost.Connect(n.ctx, pstore.PeerInfo{ID: peerid})
		}
	}
	return nil
}
