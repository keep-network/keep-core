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

type provider struct {
	channelManagerMutex sync.Mutex
	channelManagr       *channelManager

	host    host.Host
	routing routing.IpfsRouting
}

func (p *provider) ChannelFor(name string) (net.BroadcastChannel, error) {
	p.channelManagerMutex.Lock()
	defer p.channelManagerMutex.Unlock()
	return p.channelManagr.getChannel(name)
}

func (p *provider) Type() string {
	return "libp2p"
}

type Config struct {
	Peers []string
	Port  int
	Seed  int

	listenAddrs []ma.Multiaddr
	identity    *identity
}

func Connect(ctx context.Context, config *Config) (net.Provider, error) {
	host, identity, err := discoverAndListen(ctx, config)
	if err != nil {
		return nil, err
	}

	cm, err := newChannelManager(ctx, identity, host)
	if err != nil {
		return nil, err
	}

	provider := &provider{channelManagr: cm, host: host}

	dht := dht.NewDHT(ctx, provider.host, dssync.MutexWrap(dstore.NewMapDatastore()))

	provider.routing = dht

	// Wrap our host and router together into the routed host.
	// This helps us find addresses for identities we encounter in the network
	provider.host = rhost.Wrap(provider.host, provider.routing)

	// TODO: panic if we don't provide bootstrap peers
	if len(config.Peers) > 0 {
		if err := provider.bootstrap(ctx, config.Peers); err != nil {
			return nil, fmt.Errorf("Failed to bootstrap nodes with err: %v", err)
		}
	}

	return provider, nil
}

func discoverAndListen(
	ctx context.Context,
	config *Config,
) (host.Host, *identity, error) {
	var err error

	addrs := config.listenAddrs
	if addrs == nil {
		// Get available network ifaces to listen on into multiaddrs
		addrs, err = getListenAddrs(config.Port)
		if err != nil {
			return nil, nil, err
		}
	}

	peerIdentity := config.identity
	if peerIdentity == nil {
		// FIXME: revisit this fallback decision. We run into the case
		// where the user's config isn't right and then they're in the
		// network as an identity they aren't familiar with.
		peerIdentity, err = generateIdentity(config.Seed)
		if err != nil {
			return nil, nil, err
		}
	}

	peerStore, err := addIdentityToStore(peerIdentity)
	if err != nil {
		return nil, nil, err
	}

	peerHost, err := buildPeerHost(ctx, addrs, peer.ID(peerIdentity.id), peerStore)
	if err != nil {
		return nil, nil, err
	}

	if err := peerHost.Network().Listen(addrs...); err != nil {
		return nil, nil, err
	}

	return peerHost, peerIdentity, nil
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
	var (
		peers         []*peerstore.PeerInfo
		waitGroup     sync.WaitGroup
		internalError error
	)
	for _, bp := range bootstrapPeers {
		// The following code extracts target's peer ID from the
		// given multiaddress
		ipfsaddr, err := ma.NewMultiaddr(bp)
		if err != nil {
			return err
		}

		peerInfo, err := peerstore.InfoFromP2pAddr(ipfsaddr)
		if err != nil {
			return err
		}

		peers = append(peers, peerInfo)
	}

	for _, pi := range peers {
		if p.host.ID() == pi.ID {
			// We shouldn't bootstrap to ourself if we're the bootstrap node
			continue
		}
		waitGroup.Add(1)
		go func(peerInfo *peerstore.PeerInfo) {
			defer waitGroup.Done()
			if err := p.host.Connect(ctx, *peerInfo); err != nil {
				internalError = err
				return
			}
		}(pi)
	}

	waitGroup.Wait()
	return internalError
}
