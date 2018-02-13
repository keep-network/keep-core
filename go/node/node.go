package node

import (
	"bufio"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"

	iaddr "github.com/ipfs/go-ipfs-addr"
	floodsub "github.com/libp2p/go-floodsub"
	ci "github.com/libp2p/go-libp2p-crypto"
	host "github.com/libp2p/go-libp2p-host"
	peer "github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	routing "github.com/libp2p/go-libp2p-routing"
	swarm "github.com/libp2p/go-libp2p-swarm"
	bhost "github.com/libp2p/go-libp2p/p2p/host/basic"
	ma "github.com/multiformats/go-multiaddr"
)

// ErrNotEnoughBootstrapPeers signals that we do not have enough bootstrap
// peers to bootstrap correctly.
var ErrNotEnoughBootstrapPeers = errors.New("not enough bootstrap peers to bootstrap")

// for safety, publish this list to ipfs
var DefaultBootstrapAddresses = []string{
	"/ip4/127.0.0.1/tcp/2701/ipfs/QmexAnfpHrhMmAC5UNQVS8iBuUUgDrMbMY17Cck2gKrqeX",
	"/ip4/127.0.0.1/tcp/2702/ipfs/Qmd3wzD2HWA95ZAs214VxnckwkwM4GHJyC6whKUCNQhNvW",
}

// A node is the initialized Keep client waiting to join a group
type Node struct {
	// Self
	Identity *Identity

	PeerHost  host.Host
	Bootstrap []string // bootstrap peer addrs

	PeerStore pstore.Peerstore

	Floodsub *floodsub.PubSub
	Routing  routing.IpfsRouting // ugh does this have to be ipfsrouting?

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
	// Our temp routing option
	n.Routing = NewNilRouter()

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
	listen, err := ma.NewMultiaddr(fmt.Sprint("/ip4/127.0.0.1/tcp/8080"))
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

	if err := n.bootstrap(); err != nil {
		return err
	}

	return nil
}

func buildPeerHost(ctx context.Context, listenAddrs []ma.Multiaddr, pid peer.ID, ps pstore.Peerstore) (host.Host, error) {
	// TODO: use NewSwarmWithProtector
	// TODO: customize transport with config, for now use default in go-libp2p-swarm
	// Start without any addresses...
	swrm, err := swarm.NewSwarm(ctx, listenAddrs, pid, ps, nil)
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

func (n *Node) bootstrap() error {
	// first we get our list of bootstrap peers
	peers, err := getBootstrapPeers()
	if err != nil {
		return err
	}

	// next we connect to all known peers
	if err := bootstrapConnect(n.ctx, n.PeerHost, peers); err != nil {
		return err
	}

	// lastly kick off routing bootstrap
	if n.Routing != nil {
		if err := n.Routing.Bootstrap(n.ctx); err != nil {
			return err
		}
	}
	return nil
}

// copied (altered) from github.com/go-ipfs/core/bootstrap.go
func bootstrapConnect(ctx context.Context, ph host.Host, peers []pstore.PeerInfo) error {
	if len(peers) < 1 {
		return ErrNotEnoughBootstrapPeers
	}

	errs := make(chan error, len(peers))
	var wg sync.WaitGroup
	for _, p := range peers {

		// performed asynchronously because when performed synchronously, if
		// one `Connect` call hangs, subsequent calls are more likely to
		// fail/abort due to an expiring context.
		// Also, performed asynchronously for dial speed.

		wg.Add(1)
		go func(p pstore.PeerInfo) {
			defer wg.Done()
			log.Printf("%s bootstrapping to %s\n", ph.ID(), p.ID)

			ph.Peerstore().AddAddrs(p.ID, p.Addrs, pstore.PermanentAddrTTL)
			if err := ph.Connect(ctx, p); err != nil {
				log.Printf("failed to bootstrap with %v: %s\n", p.ID, err)
				errs <- err
				return
			}
			log.Printf("bootstrapDialSuccess %s", p.ID)
			log.Printf("bootstrapped with %v", p.ID)
		}(p)
	}
	wg.Wait()

	// our failure condition is when no connection attempt succeeded.
	// So drain the errs channel, counting the results.
	close(errs)
	count := 0
	var err error
	for err = range errs {
		if err != nil {
			count++
		}
	}
	if count == len(peers) {
		return fmt.Errorf("failed to bootstrap. %s", err)
	}
	return nil
}

func getBootstrapPeers() ([]pstore.PeerInfo, error) {
	peers := make([]iaddr.IPFSAddr, 0)
	for _, addr := range DefaultBootstrapAddresses {
		ia, err := iaddr.ParseString(addr)
		if err != nil {
			return nil, err
		}
		peers = append(peers, ia)
	}
	return toPeerInfos(peers), nil
}

// copied (altered) from github.com/go-ipfs/core/bootstrap.go
func toPeerInfos(bpeers []iaddr.IPFSAddr) []pstore.PeerInfo {
	pinfos := make(map[peer.ID]*pstore.PeerInfo)
	for _, bootstrap := range bpeers {
		pinfo, ok := pinfos[bootstrap.ID()]
		if !ok {
			pinfo = new(pstore.PeerInfo)
			pinfos[bootstrap.ID()] = pinfo
			pinfo.ID = bootstrap.ID()
		}

		pinfo.Addrs = append(pinfo.Addrs, bootstrap.Transport())
	}

	var peers []pstore.PeerInfo
	for _, pinfo := range pinfos {
		peers = append(peers, *pinfo)
	}

	return peers
}

func getPeers(ha host.Host) []peer.ID {
	return ha.Peerstore().Peers()
}

// Modified version from github.com/keep-network/go-experiments
func addPeers(ha host.Host, peers []peer.ID) {
	for _, p := range peers {
		if ha.ID().String() != p.String() {
			stream, err := ha.NewStream(context.Background(), p, "/add/1.0.0")
			if err != nil {
				continue
			}

			for _, addr := range ha.Addrs() {
				if addr.String() != "" {
					_, err = stream.Write([]byte(addr.String() + "/ipfs/" + ha.ID().Pretty() + "\n"))
					if err != nil {
						log.Println(err)
					}
				}
			}
			buf := bufio.NewReader(stream)
			str, err := buf.ReadString('\n')
			// The following code extracts target's the peer ID from the
			// given multiaddress
			t := strings.TrimSpace(str)
			addresses := strings.Split(t, ",")
			for _, address := range addresses {
				ipfsaddr, err := ma.NewMultiaddr(address)
				if err != nil {
					log.Println(err)
				}

				// TODO: do we need any of this ipfsaddr stuff?
				pid, err := ipfsaddr.ValueForProtocol(ma.P_IPFS)
				if err != nil {
					log.Println(err)
				}

				peerid, err := peer.IDB58Decode(pid)
				if err != nil {
					log.Println(err)
				}
				// Decapsulate the /ipfs/<peerID> part from the target
				// /ip4/<a.b.c.d>/ipfs/<peer> becomes /ip4/<a.b.c.d>
				targetAddr := ipfsaddr.Decapsulate(ipfsaddr)
				ha.Peerstore().AddAddr(peerid, targetAddr, pstore.PermanentAddrTTL)
			}
		}
	}
}

// func buildRoutingService(ctx context.Context, h host.Host) error {
// }
