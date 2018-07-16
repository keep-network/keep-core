package libp2p

import (
	"context"
	"fmt"
	"sync"

	"github.com/keep-network/keep-core/pkg/net"

	dstore "github.com/ipfs/go-datastore"
	dssync "github.com/ipfs/go-datastore/sync"
	addrutil "github.com/libp2p/go-addr-util"
	host "github.com/libp2p/go-libp2p-host"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/libp2p/go-libp2p-peerstore"
	routing "github.com/libp2p/go-libp2p-routing"
	swarm "github.com/libp2p/go-libp2p-swarm"
	basichost "github.com/libp2p/go-libp2p/p2p/host/basic"
	rhost "github.com/libp2p/go-libp2p/p2p/host/routed"

	smux "github.com/libp2p/go-stream-muxer"
	ma "github.com/multiformats/go-multiaddr"
	msmux "github.com/whyrusleeping/go-smux-multistream"
	yamux "github.com/whyrusleeping/go-smux-yamux"
)

// Config defines the configuration for the libp2p network provider.
type Config struct {
	Peers []string
	Port  int
	Seed  int
}

type provider struct {
	channelManagerMutex sync.Mutex
	channelManagr       *channelManager

	identity *identity
	host     host.Host
	routing  routing.IpfsRouting
	addrs    []ma.Multiaddr
}

func (p *provider) ChannelFor(name string) (net.BroadcastChannel, error) {
	p.channelManagerMutex.Lock()
	defer p.channelManagerMutex.Unlock()
	return p.channelManagr.getChannel(name)
}

func (p *provider) Type() string {
	return "libp2p"
}

func (p *provider) ID() net.TransportIdentifier {
	return networkIdentity(p.identity.id)
}

func (p *provider) Addrs() []ma.Multiaddr {
	return p.addrs
}

// Connect connects to a libp2p network based on the provided config. The
// connection is managed in part by the passed context, and provides access to
// the functionality specified in the net.Provider interface.
//
// An error is returned if any part of the connection or bootstrap process
// fails.
func Connect(ctx context.Context, config Config) (net.Provider, error) {
	identity, err := generateIdentity(config.Seed)
	if err != nil {
		return nil, err
	}

	host, err := discoverAndListen(ctx, identity, config.Port)
	if err != nil {
		return nil, err
	}

	cm, err := newChannelManager(ctx, identity, host)
	if err != nil {
		return nil, err
	}

	router := dht.NewDHT(ctx, host, dssync.MutexWrap(dstore.NewMapDatastore()))

	provider := &provider{
		channelManagr: cm,
		identity:      identity,
		host:          rhost.Wrap(host, router),
		routing:       router,
		addrs:         host.Addrs(),
	}

	// FIXME: return an error if we don't provide bootstrap peers
	if len(config.Peers) == 0 {
		return provider, nil
	}

	if err := provider.bootstrap(ctx, config.Peers); err != nil {
		return nil, fmt.Errorf("Failed to bootstrap nodes with err: %v", err)
	}

	return provider, nil
}

func discoverAndListen(
	ctx context.Context,
	identity *identity,
	port int,
) (host.Host, error) {
	var err error

	// Get available network ifaces to listen on into multiaddrs
	addrs, err := getListenAddrs(port)
	if err != nil {
		return nil, err
	}

	peerStore, err := addIdentityToStore(identity)
	if err != nil {
		return nil, err
	}

	peerHost, err := buildPeerHost(ctx, addrs, peer.ID(identity.id), peerStore)
	if err != nil {
		return nil, err
	}

	if err := peerHost.Network().Listen(addrs...); err != nil {
		return nil, err
	}

	return peerHost, nil
}

func getListenAddrs(port int) ([]ma.Multiaddr, error) {
	ia, err := addrutil.InterfaceAddresses()
	if err != nil {
		return nil, err
	}
	addrs := make([]ma.Multiaddr, len(ia))
	for _, addr := range ia {
		portAddr, err := ma.NewMultiaddr(fmt.Sprintf("/tcp/%d", port))
		if err != nil {
			return nil, err
		}
		addrs = append(addrs, addr.Encapsulate(portAddr))
	}
	return addrs, nil
}

func buildPeerHost(
	ctx context.Context,
	listenAddrs []ma.Multiaddr,
	pid peer.ID,
	peerStore peerstore.Peerstore,
) (host.Host, error) {
	smuxTransport := makeSmuxTransport()

	swrm, err := swarm.NewSwarmWithProtector(ctx, listenAddrs, pid, peerStore, nil, smuxTransport, nil)
	if err != nil {
		return nil, err
	}

	network := (*swarm.Network)(swrm)
	opts := &basichost.HostOpts{NATManager: basichost.NewNATManager(network)}
	h, err := basichost.NewHost(ctx, network, opts)
	if err != nil {
		if cerr := h.Close(); cerr != nil {
			return nil, cerr
		}
		return nil, err
	}

	return h, nil
}

func makeSmuxTransport() smux.Transport {
	multiStreamTransport := msmux.NewBlankTransport()
	yamuxTransport := yamux.DefaultTransport

	multiStreamTransport.AddTransport("/yamux/1.0.0", yamuxTransport)
	return multiStreamTransport
}

func (p *provider) bootstrap(ctx context.Context, bootstrapPeers []string) error {
	var waitGroup sync.WaitGroup

	peerInfos, err := extractMultiAddrFromPeers(bootstrapPeers)
	if err != nil {
		return err
	}

	for _, peerInfo := range peerInfos {
		if p.host.ID() == peerInfo.ID {
			// We shouldn't bootstrap to ourself if we're the
			// bootstrap node.
			continue
		}
		waitGroup.Add(1)
		go func(pi *peerstore.PeerInfo) {
			defer waitGroup.Done()
			if err := p.host.Connect(ctx, *pi); err != nil {
				fmt.Println(err)
				return
			}
		}(peerInfo)
	}
	waitGroup.Wait()

	// Bootstrap the host
	return p.routing.Bootstrap(ctx)
}

func extractMultiAddrFromPeers(peers []string) ([]*peerstore.PeerInfo, error) {
	var peerInfos []*peerstore.PeerInfo
	for _, peer := range peers {
		ipfsaddr, err := ma.NewMultiaddr(peer)
		if err != nil {
			return nil, err
		}

		peerInfo, err := peerstore.InfoFromP2pAddr(ipfsaddr)
		if err != nil {
			return nil, err
		}

		peerInfos = append(peerInfos, peerInfo)
	}
	return peerInfos, nil
}
