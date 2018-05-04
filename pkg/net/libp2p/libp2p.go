package libp2p

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"
	"time"

	"github.com/keep-network/keep-core/pkg/net"
	addrutil "github.com/libp2p/go-addr-util"
	host "github.com/libp2p/go-libp2p-host"
	peer "github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	swarm "github.com/libp2p/go-libp2p-swarm"

	bhost "github.com/libp2p/go-libp2p/p2p/host/basic"
	smux "github.com/libp2p/go-stream-muxer"
	ma "github.com/multiformats/go-multiaddr"
	msmux "github.com/whyrusleeping/go-smux-multistream"
	yamux "github.com/whyrusleeping/go-smux-yamux"
)

type Peer struct {
	cm *channelManager
	mu sync.Mutex // guards channel manager

	host host.Host
}

func (p *Peer) ChannelFor(name string) net.BroadcastChannel {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.cm.getChannel(name)
}

func (p *Peer) Type() string {
	return "libp2p"
}

func Connect(ctx context.Context, port int) (net.Provider, error) {
	cm, err := newChannelManager(ctx, nil)
	if err != nil {
		return nil, err
	}

	p := &Peer{cm: cm}

	host, addrs, err := discv(ctx, port)
	if err != nil {
		return nil, err
	}
	p.host = host

	// Ok, now we're ready to listen
	if err := p.host.Network().Listen(addrs...); err != nil {
		return nil, err
	}

	return p, nil
}

func discv(
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
	opts := &bhost.HostOpts{NATManager: bhost.NewNATManager(network)}
	h, err := bhost.NewHost(ctx, network, opts)
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

	ymxtpt := &yamux.Transport{
		AcceptBacklog:          512,
		ConnectionWriteTimeout: time.Second * 10,
		KeepAliveInterval:      time.Second * 30,
		EnableKeepAlive:        true,
		MaxStreamWindowSize:    uint32(1024 * 512),
		LogOutput:              ioutil.Discard,
	}

	mstpt.AddTransport("/yamux/1.0.0", ymxtpt)
	// TODO: compile error, return correct type
	return nil
}
