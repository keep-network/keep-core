package libp2p

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/keep-network/keep-core/pkg/net"

	addrutil "github.com/libp2p/go-addr-util"
	host "github.com/libp2p/go-libp2p-host"
	peer "github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	swarm "github.com/libp2p/go-libp2p-swarm"
	basichost "github.com/libp2p/go-libp2p/p2p/host/basic"

	smux "github.com/libp2p/go-stream-muxer"
	ma "github.com/multiformats/go-multiaddr"
	msmux "github.com/whyrusleeping/go-smux-multistream"
	yamux "github.com/whyrusleeping/go-smux-yamux"
)

type provider struct {
	cm                  *channelManager
	channelManagerMutex sync.Mutex

	host host.Host
}

func (p *provider) ChannelFor(name string) (net.BroadcastChannel, error) {
	p.channelManagerMutex.Lock()
	defer p.channelManagerMutex.Unlock()
	return p.cm.getChannel(name)
}

func (p *provider) Type() string {
	return "libp2p"
}

type Config struct {
	port        int
	listenAddrs []ma.Multiaddr
	identity    *identity
}

func Connect(ctx context.Context, config *Config) (net.Provider, error) {
	host, err := discoverAndListen(ctx, config)
	if err != nil {
		return nil, err
	}

	cm, err := newChannelManager(ctx, host)
	if err != nil {
		return nil, err
	}

	return &provider{cm: cm, host: host}, nil
}

func discoverAndListen(
	ctx context.Context,
	config *Config,
) (host.Host, error) {
	var err error

	addrs := config.listenAddrs
	if addrs == nil {
		// Get available network ifaces to listen on into multiaddrs
		addrs, err = getListenAddrs(config.port)
		if err != nil {
			return nil, err
		}
	}

	nonLocalAddrs := make([]ma.Multiaddr, 0)
	for _, addr := range addrs {
		stringily := fmt.Sprintf("%v", addr)
		if strings.Contains(stringily, "10.240") {
			nonLocalAddrs = append(nonLocalAddrs, addr)
		}
	}

	peerIdentity := config.identity
	if peerIdentity == nil {
		// FIXME: revisit this fallback decision. We run into the case
		// where the user's config isn't right and then they're in the
		// network as an identity they aren't familiar with.
		peerIdentity, err = generateIdentity()
		if err != nil {
			return nil, err
		}
	}

	peerStore, err := addIdentityToStore(peerIdentity)
	if err != nil {
		return nil, err
	}

	peerHost, err := buildPeerHost(ctx, nonLocalAddrs, peerIdentity.id, peerStore)
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

// TODO: encrypted connections
func buildPeerHost(
	ctx context.Context,
	listenAddrs []ma.Multiaddr,
	pid peer.ID,
	peerStore pstore.Peerstore,
) (host.Host, error) {
	smuxTransport := makeSmuxTransport()

	// TODO: Pass in protec and metrics reporter
	swrm, err := swarm.NewSwarmWithProtector(ctx, listenAddrs, pid, peerStore, nil, smuxTransport, nil)
	if err != nil {
		return nil, err
	}

	network := (*swarm.Network)(swrm)
	// TODO: use our own host
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
