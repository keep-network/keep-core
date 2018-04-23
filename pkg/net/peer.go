package net

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/keep-network/keep-core/pkg/net/identity"
	addrutil "github.com/libp2p/go-addr-util"
	host "github.com/libp2p/go-libp2p-host"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	swarm "github.com/libp2p/go-libp2p-swarm"
	smux "github.com/libp2p/go-stream-muxer"
	ma "github.com/multiformats/go-multiaddr"
	msmux "github.com/whyrusleeping/go-smux-multistream"
	yamux "github.com/whyrusleeping/go-smux-yamux"
)

type Peer struct {
	host.Host
	Authenticator

	ID    identity.Identity
	Store pstore.Peerstore
}

func NewPeer(randseed int64, filepath string) *Peer {
	pi, err := identity.LoadOrGenerateIdentity(randseed, filepath)
	if err != nil {
		panic(fmt.Sprintf("Failed to generate Identity with error %s", err))
	}
	ps, err := pi.AddIdentityToStore()
	if err != nil {
		panic(fmt.Sprintf("Failed to add Identity to PeerStore with error %s", err))
	}
	return &Peer{ID: pi, Store: ps}
}

func (p *Peer) Connect(ctx context.Context, port int) error {
	// Convert available network ifaces to listen on into multiaddrs
	addrs, err := getListenAdresses(port)
	if err != nil {
		return err
	}

	// TODO: should we limit which interfaces we attempt to listen to?
	nonLocalAddrs := make([]ma.Multiaddr, 0)
	for _, addr := range addrs {
		stringily := fmt.Sprintf("%v", addr)
		// TODO: limited to this cidr in testing
		if strings.Contains(stringily, "10.240") {
			nonLocalAddrs = append(nonLocalAddrs, addr)
		}
	}

	p.Host, err = p.buildPeerHost(nonLocalAddrs)
	if err != nil {
		return err
	}

	// TODO: bootstrap client

	// Ok, now we're ready to listen
	if err := p.Network().Listen(addrs...); err != nil {
		return err
	}

	return nil
}

func (p *Peer) buildPeerHost(listenAddrs []ma.Multiaddr) (host.Host, error) {
	// Set up stream multiplexer
	tpt := makeSmuxTransport()

	// TODO: Pass in protec and metrics reporter
	swrm, err := swarm.NewSwarmWithProtector(ctx, listenAddrs, p.ID.ID(), p.Store, nil, tpt, nil)
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

// TODO: Allow for user-scoped listeners to either override this or union with this.
func getListenAdresses(port int) ([]ma.Multiaddr, error) {
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
