// node is the base relay client initalized on startup
package relayclient

import (
	"context"
	crand "crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	mrand "math/rand"
	"time"

	dstore "github.com/ipfs/go-datastore"
	dssync "github.com/ipfs/go-datastore/sync"
	addrutil "github.com/libp2p/go-addr-util"
	ci "github.com/libp2p/go-libp2p-crypto"
	host "github.com/libp2p/go-libp2p-host"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	peer "github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	protocol "github.com/libp2p/go-libp2p-protocol"
	routing "github.com/libp2p/go-libp2p-routing"
	swarm "github.com/libp2p/go-libp2p-swarm"
	bhost "github.com/libp2p/go-libp2p/p2p/host/basic"
	rhost "github.com/libp2p/go-libp2p/p2p/host/routed"
	smux "github.com/libp2p/go-stream-muxer"
	ma "github.com/multiformats/go-multiaddr"
	msmux "github.com/whyrusleeping/go-smux-multistream"
	yamux "github.com/whyrusleeping/go-smux-yamux"
)

// Protocol is an identifier to write protocol headers in streams.
// TODO: actually enforce this
var Protocol = protocol.ID("/keep/relay_client/0.0.1")

// hardcoded peers that will make this prototype work
// taken from go-experiments
var bootstrapPeers = []string{"/ip4/127.0.0.1/tcp/2701/ipfs/QmexAnfpHrhMmAC5UNQVS8iBuUUgDrMbMY17Cck2gKrqeX", "/ip4/127.0.0.1/tcp/2702/ipfs/Qmd3wzD2HWA95ZAs214VxnckwkwM4GHJyC6whKUCNQhNvW"}

// A node is the initialized relay client waiting to join a group
type Node struct {
	// Self
	Identity *Identity

	PeerStore pstore.Peerstore
	PeerHost  host.Host

	Routing routing.IpfsRouting

	Groups *GroupManager

	// groupDKG   chan *pb.DKGMessage
	// groupRelay chan *pb.RelayMessage

	// Use to detect node shutdowns
	ctx context.Context
}

// An Identity contains all libp2p and secret identifying information
type Identity struct {
	PeerID  peer.ID
	PubKey  ci.PubKey
	privKey ci.PrivKey
}

// NewNode should only be called once, on init
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
	n := &Node{Identity: &Identity{PeerID: pid, privKey: priv, PubKey: pub}}

	// The context governs the lifetime of the libp2p node
	n.ctx = ctx

	err = n.discoverAndConnect(n.ctx, port, n.Identity)
	if err != nil {
		return nil, err
	}

	// https: //github.com/libp2p/go-floodsub/issues/65#issuecomment-365680860
	dht := dht.NewDHT(n.ctx, n.PeerHost, dssync.MutexWrap(dstore.NewMapDatastore()))

	// TODO: add comments
	n.Routing = dht

	// TODO: add comments
	n.PeerHost = rhost.Wrap(n.PeerHost, n.Routing)

	if err := n.bootstrap(ctx); err != nil {
		return nil, fmt.Errorf("Failed to bootstrap nodes with err: %v", err)
	}

	n.Groups, err = NewGroupManager(ctx, n.Identity, n.PeerHost, dht)
	if err != nil {
		return nil, err
	}

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

func (n *Node) discoverAndConnect(ctx context.Context, port int, id *Identity) error {
	var err error

	// Ensure that other members in our broadcast channel can identify us
	// TODO: just pass in the Identity struct - maybe
	n.PeerStore, err = addToPeerStore(id)
	if err != nil {
		return err
	}

	// Convert available network ifaces to listen on into multiaddrs
	addrs, err := getListenAdresses(port)
	if err != nil {
		return err
	}

	// TODO: flesh out how we connect to libp2p
	n.PeerHost, err = buildPeerHost(ctx, addrs, id.PeerID, n.PeerStore)
	if err != nil {
		return err
	}

	// Ok, now we're ready to listen
	if err := n.PeerHost.Network().Listen(addrs...); err != nil {
		return err
	}
	// TODO: implement a standard and functional logger
	log.Printf("Listening at: %#v\n", addrs)

	return nil
}

func addToPeerStore(id *Identity) (pstore.Peerstore, error) {
	ps := pstore.NewPeerstore()
	// HACK: see github.com/rargulati/go-libp2p-crypto for fix
	if err := ps.AddPrivKey(id.PeerID, id.privKey); err != nil {
		fmt.Println("private key mishap")
		return nil, err
	}
	if err := ps.AddPubKey(id.PeerID, id.PubKey); err != nil {
		fmt.Println("pub key mishap")
		return nil, err
	}
	return ps, nil
}

// TODO: Allow for user-scoped listeners to either override this or union with this.
func getListenAdresses(port int) ([]ma.Multiaddr, error) {
	// TODO: attach a muxer to a connection
	// TODO: figure out go-libp2p-interface-pnet.Protector and go-libp2p-pnet.NewProtector - later
	ia, err := addrutil.InterfaceAddresses()
	if err != nil {
		return nil, err
	}
	addrs := make([]ma.Multiaddr, len(ia), len(ia))
	for _, addr := range ia {
		portAddr, err := ma.NewMultiaddr(fmt.Sprintf("/tcp/%d", port))
		if err != nil {
			return nil, err
		}
		addrs = append(addrs, addr.Encapsulate(portAddr))
	}
	return addrs, nil
}

func buildPeerHost(ctx context.Context, listenAddrs []ma.Multiaddr, pid peer.ID, ps pstore.Peerstore) (host.Host, error) {
	// Set up stream multiplexer
	tpt := makeSmuxTransport()

	// TODO: Pass in protec and metrics reporter
	swrm, err := swarm.NewSwarmWithProtector(ctx, listenAddrs, pid, ps, nil, tpt, nil)
	if err != nil {
		return nil, err
	}

	network := (*swarm.Network)(swrm)
	// TODO: use our own host, I'm unsure about the utility of basic
	opts := &bhost.HostOpts{NATManager: bhost.NewNATManager(network)}
	// TODO: does host leak?
	h, err := bhost.NewHost(ctx, network, opts)
	if err != nil {
		h.Close()
		return nil, err
	}
	// TODO: do we need to enable the circuit relay? if so, do it here
	return h, nil
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

// lifted and modified from github.com/keep-network/go-experiments
func (n *Node) bootstrap(ctx context.Context) error {
	log.Println("Bootstrapping peers...")
	for _, p := range bootstrapPeers {
		// The following code extracts target's the peer ID from the
		// given multiaddress
		ipfsaddr, err := ma.NewMultiaddr(p)
		if err != nil {
			log.Fatalln(err)
			return err
		}

		pid, err := ipfsaddr.ValueForProtocol(ma.P_IPFS)
		if err != nil {
			log.Fatalln(err)
			return err
		}

		peerid, err := peer.IDB58Decode(pid)
		if err != nil {
			log.Fatalln(err)
			return err
		}

		// Decapsulate the /ipfs/<peerID> part from the target
		// /ip4/<a.b.c.d>/ipfs/<peer> becomes /ip4/<a.b.c.d>
		targetPeerAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", peer.IDB58Encode(peerid)))
		targetAddr := ipfsaddr.Decapsulate(targetPeerAddr)
		if n.PeerHost.ID().String() != peerid.String() {
			// We have a peer ID and a targetAddr so we add it to the peerstore
			// so LibP2P knows how to contact it
			n.PeerHost.Peerstore().AddAddr(peerid, targetAddr, pstore.PermanentAddrTTL)
			n.PeerHost.Connect(ctx, pstore.PeerInfo{ID: peerid})
			// implement handling errors for connect after you have retries
			// if err != nil {
			// 	log.Fatalln(err)
			// 	// return err
			// }
		}
	}

	// Bootstrap the host
	err := n.Routing.Bootstrap(ctx)
	if err != nil {
		log.Fatalln(err)
		return err
	}

	return nil
}
