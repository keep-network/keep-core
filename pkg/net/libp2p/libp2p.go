package libp2p

import (
	"context"
	"fmt"
	"sync"

	"github.com/keep-network/keep-core/pkg/net"

	dstore "github.com/ipfs/go-datastore"
	dssync "github.com/ipfs/go-datastore/sync"
	"github.com/libp2p/go-addr-util"
	"github.com/libp2p/go-libp2p-host"
	"github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p-peer"
	"github.com/libp2p/go-libp2p-peerstore"
	"github.com/libp2p/go-libp2p-routing"
	"github.com/libp2p/go-libp2p-swarm"
	"github.com/libp2p/go-libp2p/p2p/host/basic"
	rhost "github.com/libp2p/go-libp2p/p2p/host/routed"

	"github.com/keep-network/keep-core/util"
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
	addrs   []ma.Multiaddr
}

// ListenAddrs for this host.
var ListenAddrs []ma.Multiaddr

func (p *provider) ChannelFor(name string) (net.BroadcastChannel, error) {
	p.channelManagerMutex.Lock()
	defer p.channelManagerMutex.Unlock()
	return p.channelManagr.getChannel(name)
}

func (p *provider) Type() string {
	return "libp2p"
}

func (p *provider) Addrs() []ma.Multiaddr {
	return p.addrs
}

const ipfsURLPattern = `.+\/.*`

var ifpsURLRegex = util.CompileRegex(ipfsURLPattern)

// NodeConfig contains the config values for this node.
type NodeConfig struct {
	Port  int
	Seed  int
	Peers []string
}

var nodeConfig = &NodeConfig{Port: 0, Seed: 0, Peers: nil}

// DefaultNodeConfig is a non-bootrap node.
var DefaultNodeConfig = NodeConfig{
	Port:  27001,
	Seed:  0,
	Peers: []string{"/ip4/127.0.0.1/tcp/27001/ipfs/12D3KooWKRyzVWW6ChFjQjK4miCty85Niy49tpPV95XdKu1BcvMA"},
}

// Config contains the data needed to configure lib2p resources for this provider.
type Config struct {
	NodeConfig

	listenAddrs []ma.Multiaddr
	identity    *identity
}

// ValidationError returns validation errors for all config values.
func (c *Config) ValidationError() error {
	var errMsgs []string
	if c.Port <= 0 {
		errMsgs = append(errMsgs,
			fmt.Sprintf("Node.Port (%d) invalid; see node section in config file or use --port flag",
				c.Port))
	}
	if len(c.Peers) == 0 && c.Seed <= 0 {
		errMsgs = append(errMsgs, fmt.Sprintf("either supply valid Node.Peers or a valid Bootstrap.Seed"))
	}
	if len(c.Peers) > 0 && c.Seed != 0 {
		errMsgs = append(errMsgs, fmt.Sprintf("non-bootstrap node should have Bootstrap.URL and a Bootstrap.Seed of 0"))
	}
	if len(c.Peers) > 0 {
		for _, ipfsURL := range c.Peers {
			if !util.MatchFound(ifpsURLRegex, ipfsURL) {
				errMsgs = append(errMsgs,
					fmt.Sprintf("Node.Peers (%s) invalid; format expected: %s",
						ipfsURL,
						ipfsURLPattern))
			}
		}
		if util.DuplicatesExist(c.Peers) {
			errMsgs = append(errMsgs,
				fmt.Sprintf("Node.Peers invalid; duplicates found: %s",
					util.Join(util.Duplicates(c.Peers), " ")))
		}
	}
	return util.Err(errMsgs)
}

// Connect returns the host provider with channel manager, router and listen addresses.
func Connect(ctx context.Context, config *Config) (net.Provider, error) {
	host, identity, err := discoverAndListen(ctx, config)
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
	config *Config,
) (host.Host, *identity, error) {
	var err error

	addrs := config.listenAddrs
	if addrs == nil {
		// Get available network ifaces to listen on into multiaddrs.
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
	var waitGroup sync.WaitGroup

	peerInfos, err := extractMultiAddrFromPeers(bootstrapPeers)
	if err != nil {
		return err
	}

	for _, peerInfo := range peerInfos {
		if p.host.ID() == peerInfo.ID {
			// We shouldn't bootstrap to ourself if we're the bootstrap node.
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

	// Bootstrap the host.
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
