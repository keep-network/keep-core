package libp2p

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/net/key"

	dstore "github.com/ipfs/go-datastore"
	dssync "github.com/ipfs/go-datastore/sync"
	addrutil "github.com/libp2p/go-addr-util"
	libp2p "github.com/libp2p/go-libp2p"
	connmgr "github.com/libp2p/go-libp2p-connmgr"
	host "github.com/libp2p/go-libp2p-host"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p-peerstore"
	rhost "github.com/libp2p/go-libp2p/p2p/host/routed"

	ma "github.com/multiformats/go-multiaddr"
)

// Defaults from ipfs
const (
	// DefaultConnMgrHighWater is the default value for the connection managers
	// 'high water' mark
	DefaultConnMgrHighWater = 900

	// DefaultConnMgrLowWater is the default value for the connection managers 'low
	// water' mark
	DefaultConnMgrLowWater = 600

	// DefaultConnMgrGracePeriod is the default value for the connection managers
	// grace period
	DefaultConnMgrGracePeriod = time.Second * 20
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
	routing  *dht.IpfsDHT
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

func (p *provider) AddrStrings() []string {
	multiaddrStrings := make([]string, 0, len(p.addrs))
	for _, multiaddr := range p.addrs {
		addrWithIdentity := fmt.Sprintf("%s/ipfs/%s", multiaddr.String(), p.identity.id.Pretty())
		multiaddrStrings = append(multiaddrStrings, addrWithIdentity)
	}

	return multiaddrStrings
}

func (p *provider) Peers() []string {
	var peers []string
	peersIDSlice := p.host.Peerstore().Peers()
	for _, peer := range peersIDSlice {
		// filter out our own node
		if peer == p.identity.id {
			continue
		}
		peers = append(peers, peer.Pretty())
	}
	return peers
}

// Connect connects to a libp2p network based on the provided config. The
// connection is managed in part by the passed context, and provides access to
// the functionality specified in the net.Provider interface.
//
// An error is returned if any part of the connection or bootstrap process
// fails.
func Connect(
	ctx context.Context,
	config Config,
	staticKey *key.StaticNetworkKey,
) (net.Provider, error) {
	identity, err := createIdentity(staticKey)
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

	// Get available network ifaces, for a specific port, as multiaddrs
	addrs, err := getListenAddrs(port)
	if err != nil {
		return nil, err
	}

	return libp2p.New(ctx,
		libp2p.ListenAddrs(addrs...),
		libp2p.Identity(identity.privKey),
		libp2p.Security(handshakeID, newAuthenticatedTransport),
		libp2p.ConnectionManager(
			connmgr.NewConnManager(
				DefaultConnMgrLowWater,
				DefaultConnMgrHighWater,
				DefaultConnMgrGracePeriod,
			),
		),
	)
}

func getListenAddrs(port int) ([]ma.Multiaddr, error) {
	ia, err := addrutil.InterfaceAddresses()
	if err != nil {
		return nil, err
	}
	addrs := make([]ma.Multiaddr, 0)
	for _, addr := range ia {
		portAddr, err := ma.NewMultiaddr(fmt.Sprintf("/tcp/%d", port))
		if err != nil {
			return nil, err
		}
		addrs = append(addrs, addr.Encapsulate(portAddr))
	}
	return addrs, nil
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
