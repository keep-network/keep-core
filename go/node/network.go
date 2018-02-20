// network contains our bridge to libp2p
package node

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	dstore "github.com/ipfs/go-datastore"
	dssync "github.com/ipfs/go-datastore/sync"
	addrutil "github.com/libp2p/go-addr-util"
	floodsub "github.com/libp2p/go-floodsub"
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

// Identifier to write protocol headers in streams.
var Protocol = protocol.ID("/keep/relay_client/0.0.1")

// hardcoded peers that will make this prototype work
var bootstrapPeers = []string{"/ip4/127.0.0.1/tcp/2701/ipfs/QmexAnfpHrhMmAC5UNQVS8iBuUUgDrMbMY17Cck2gKrqeX", "/ip4/127.0.0.1/tcp/2702/ipfs/Qmd3wzD2HWA95ZAs214VxnckwkwM4GHJyC6whKUCNQhNvW"}

type NetworkManager struct {
	PeerStore pstore.Peerstore
	PeerHost  host.Host
	Sub       *floodsub.PubSub
	Routing   routing.IpfsRouting
}

func NewNetworkManager(ctx context.Context, port int, pid peer.ID, priv ci.PrivKey,
	pub ci.PubKey) (*NetworkManager, error) {
	var err error
	n := &NetworkManager{}

	// Ensure that other members in our broadcast channel can identify us
	n.PeerStore, err = addToPeerStore(pid, priv, pub)
	if err != nil {
		return nil, err
	}

	// Convert available network ifaces to listen on into multiaddrs
	addrs, err := getListenAdresses(port)
	if err != nil {
		return nil, err
	}

	// TODO: flesh out how we connect to libp2p
	n.PeerHost, err = buildPeerHost(ctx, addrs, pid, n.PeerStore)
	if err != nil {
		return nil, err
	}

	// Ok, now we're ready to listen
	if err := n.PeerHost.Network().Listen(addrs...); err != nil {
		return nil, err
	}
	// TODO: implement a standard and functional logger
	log.Printf("Listening at: %#v\n", addrs)

	// n.Sub, err = floodsub.NewFloodSub(ctx, n.PeerHost)
	n.Sub, err = floodsub.NewGossipSub(ctx, n.PeerHost)
	if err != nil {
		return nil, err
	}
	// https: //github.com/libp2p/go-floodsub/issues/65#issuecomment-365680860
	n.Routing = dht.NewDHT(ctx, n.PeerHost, dssync.MutexWrap(dstore.NewMapDatastore()))
	n.PeerHost = rhost.Wrap(n.PeerHost, n.Routing)

	if err := n.bootstrap(ctx); err != nil {
		return nil, fmt.Errorf("Failed to bootstrap nodes with err: %v", err)
	}

	return n, nil
}

func addToPeerStore(pid peer.ID, priv ci.PrivKey, pub ci.PubKey) (pstore.Peerstore, error) {
	ps := pstore.NewPeerstore()
	// FIXME: I made a hack to go-libp2p-peer to get this work
	if err := ps.AddPrivKey(pid, priv); err != nil {
		fmt.Println("private key mishap")
		return nil, err
	}
	if err := ps.AddPubKey(pid, pub); err != nil {
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
func (n *NetworkManager) bootstrap(ctx context.Context) error {
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
			n.PeerHost.Connect(ctx, pstore.PeerInfo{ID: peerid})
		}
	}
	return nil
}
