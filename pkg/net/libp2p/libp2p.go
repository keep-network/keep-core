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

	"github.com/libp2p/go-libp2p/p2p/host/basic"
	smux "github.com/libp2p/go-stream-muxer"
	ma "github.com/multiformats/go-multiaddr"
	msmux "github.com/whyrusleeping/go-smux-multistream"
	yamux "github.com/whyrusleeping/go-smux-yamux"
)

type proxy struct {
	cm                  *channelManager
	channelManagerMutex sync.Mutex

	host host.Host
}

func (p *proxy) ChannelFor(name string) net.BroadcastChannel {
	p.channelManagerMutex.Lock()
	defer p.channelManagerMutex.Unlock()
	return p.cm.getChannel(name)
}

func (p *proxy) Type() string {
	return "libp2p"
}

func Connect(ctx context.Context, c *net.Config) (net.Provider, error) {
	cm, err := newChannelManager(ctx, nil)
	if err != nil {
		return nil, err
	}

	p := &proxy{cm: cm}
	host, addrs, err := discover(ctx, c.Port)
	if err != nil {
		return nil, err
	}
	p.host = host

	if err := p.host.Network().Listen(addrs...); err != nil {
		return nil, err
	}

	return p, nil
}

func discover(
	ctx context.Context,
	port int,
) (host.Host, []ma.Multiaddr, error) {
	// Convert available network ifaces to listen on into multiaddrs
	addrs, err := getListenAdresses(port)
	if err != nil {
		return nil, nil, err
	}

	nonLocalAddrs := make([]ma.Multiaddr, 0)
	for _, addr := range addrs {
		stringily := fmt.Sprintf("%v", addr)
		if strings.Contains(stringily, "10.240") {
			nonLocalAddrs = append(nonLocalAddrs, addr)
		}
	}

	h, err := buildPeerHost(ctx, nonLocalAddrs, peer.ID(""), nil)
	if err != nil {
		return nil, nil, err
	}
	return h, addrs, nil
}

// TODO: Allow for user-scoped listeners to either override this or union with this.
func getListenAdresses(port int) ([]ma.Multiaddr, error) {
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
	ps pstore.Peerstore,
) (host.Host, error) {
	// Set up stream multiplexer
	tpt := makeSmuxTransport()

	// TODO: Pass in protec and metrics reporter
	swrm, err := swarm.NewSwarmWithProtector(ctx, listenAddrs, pid, ps, nil, tpt, nil)
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
	mstpt := msmux.NewBlankTransport()
	ymxtpt := yamux.DefaultTransport
	mstpt.AddTransport("/yamux/1.0.0", ymxtpt)
	// TODO: compile error, return correct type
	return nil
}
