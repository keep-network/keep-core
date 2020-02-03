package libp2p

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ipfs/go-log"

	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/net/key"
	"github.com/keep-network/keep-core/pkg/net/retransmission"
	"github.com/keep-network/keep-core/pkg/net/watchtower"

	dstore "github.com/ipfs/go-datastore"
	dssync "github.com/ipfs/go-datastore/sync"
	addrutil "github.com/libp2p/go-addr-util"
	libp2p "github.com/libp2p/go-libp2p"
	connmgr "github.com/libp2p/go-libp2p-connmgr"
	host "github.com/libp2p/go-libp2p-core/host"
	libp2pnet "github.com/libp2p/go-libp2p-core/network"
	peer "github.com/libp2p/go-libp2p-core/peer"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	rhost "github.com/libp2p/go-libp2p/p2p/host/routed"

	bootstrap "github.com/keep-network/go-libp2p-bootstrap"
	ma "github.com/multiformats/go-multiaddr"
)

var logger = log.Logger("keep-net-libp2p")

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

// watchtower constants
const (
	// StakeCheckTick is the amount of time between periodic checks for
	// minimum stake for all peers connected to this one.
	StakeCheckTick = time.Minute * 1
	// BootstrapCheckPeriod is the amount of time between periodic checks
	// for ensuring we are connected to an appropriate number of bootstrap
	// peers.
	BootstrapCheckPeriod = 10 * time.Second
)

// Config defines the configuration for the libp2p network provider.
type Config struct {
	Peers              []string
	Port               int
	AnnouncedAddresses []string
}

type provider struct {
	channelManagerMutex sync.Mutex
	channelManagr       *channelManager

	identity *identity
	host     host.Host
	routing  *dht.IpfsDHT
	addrs    []ma.Multiaddr

	connectionManager *connectionManager
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
		multiaddrStrings = append(
			multiaddrStrings,
			multiaddressWithIdentity(multiaddr, p.identity.id),
		)
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
		peers = append(peers, peer.String())
	}
	return peers
}

func (p *provider) ConnectionManager() net.ConnectionManager {
	return p.connectionManager
}

type connectionManager struct {
	host.Host
}

func (cm *connectionManager) ConnectedPeers() []string {
	var peers []string
	for _, connectedPeer := range cm.Network().Peers() {
		peers = append(peers, connectedPeer.String())
	}
	return peers
}

func (cm *connectionManager) GetPeerPublicKey(connectedPeer string) (*key.NetworkPublic, error) {
	peerID, err := peer.IDB58Decode(connectedPeer)
	if err != nil {
		return nil, fmt.Errorf(
			"Failed to decode peer ID from [%s] with error: [%v]",
			connectedPeer,
			err,
		)
	}

	peerPublicKey, err := peerID.ExtractPublicKey()
	if err != nil {
		return nil, fmt.Errorf(
			"Failed to extract peer [%s] public key with error: [%v]",
			connectedPeer,
			err,
		)
	}

	return key.Libp2pKeyToNetworkKey(peerPublicKey), nil
}

func (cm *connectionManager) DisconnectPeer(peerHash string) {
	peerID, err := peer.IDB58Decode(peerHash)
	if err != nil {
		logger.Errorf("failed to decode peer hash: [%v] [%v]", peerHash, err)
		return
	}

	connections := cm.Network().ConnsToPeer(peerID)
	for _, connection := range connections {
		if err := connection.Close(); err != nil {
			logger.Errorf("failed to disconnect: [%v]", err)
		}
	}
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
	staticKey *key.NetworkPrivate,
	stakeMonitor chain.StakeMonitor,
	ticker *retransmission.Ticker,
) (net.Provider, error) {
	identity, err := createIdentity(staticKey)
	if err != nil {
		return nil, err
	}

	host, err := discoverAndListen(
		ctx,
		identity,
		config.Port,
		config.AnnouncedAddresses,
		stakeMonitor,
	)
	if err != nil {
		return nil, err
	}

	host.Network().Notify(buildNotifiee())

	cm, err := newChannelManager(ctx, identity, host, ticker)
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

	if len(config.Peers) == 0 {
		logger.Infof("node's peers list is empty")
	}

	if err := provider.bootstrap(ctx, config.Peers); err != nil {
		return nil, fmt.Errorf("Failed to bootstrap nodes with err: %v", err)
	}

	provider.connectionManager = &connectionManager{provider.host}

	// Instantiates and starts the connection management background process
	watchtower.NewGuard(
		ctx, StakeCheckTick, stakeMonitor, provider.connectionManager,
	)

	return provider, nil
}

func discoverAndListen(
	ctx context.Context,
	identity *identity,
	port int,
	announcedAddresses []string,
	stakeMonitor chain.StakeMonitor,
) (host.Host, error) {
	var err error

	// Get available network ifaces, for a specific port, as multiaddrs
	addrs, err := getListenAddrs(port)
	if err != nil {
		return nil, err
	}

	transport, err := newEncryptedAuthenticatedTransport(
		identity.privKey,
		stakeMonitor,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"could not create authenticated transport [%v]",
			err,
		)
	}

	options := []libp2p.Option{
		libp2p.ListenAddrs(addrs...),
		libp2p.Identity(identity.privKey),
		libp2p.Security(handshakeID, transport),
		libp2p.ConnectionManager(
			connmgr.NewConnManager(
				DefaultConnMgrLowWater,
				DefaultConnMgrHighWater,
				DefaultConnMgrGracePeriod,
			),
		),
	}

	if addresses := parseMultiaddresses(announcedAddresses); len(addresses) > 0 {
		addressFactory := func(addrs []ma.Multiaddr) []ma.Multiaddr {
			logger.Debugf(
				"replacing default announced addresses [%v] with [%v]",
				addrs,
				addresses,
			)
			return addresses
		}
		options = append(options, libp2p.AddrsFactory(addressFactory))
	}

	return libp2p.New(ctx, options...)
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

func parseMultiaddresses(addresses []string) []ma.Multiaddr {
	multiaddresses := make([]ma.Multiaddr, 0)
	for _, address := range addresses {
		multiaddress, err := ma.NewMultiaddr(address)
		if err != nil {
			logger.Warningf(
				"could not parse address string [%v]: [%v]",
				address,
				err,
			)
			continue
		}
		multiaddresses = append(multiaddresses, multiaddress)
	}

	return multiaddresses
}

func (p *provider) bootstrap(ctx context.Context, bootstrapPeers []string) error {
	peerInfos, err := extractMultiAddrFromPeers(bootstrapPeers)
	if err != nil {
		return err
	}

	bootstraConfig := bootstrap.BootstrapConfigWithPeers(peerInfos)

	// TODO: allow this to be a configurable value
	bootstraConfig.Period = BootstrapCheckPeriod

	// TODO: use the io.Closer to shutdown the bootstrapper when we build out
	// a shutdown process.
	_, err = bootstrap.Bootstrap(
		p.identity.id,
		p.host,
		p.routing,
		bootstraConfig,
	)
	return err
}

func extractMultiAddrFromPeers(peers []string) ([]peerstore.PeerInfo, error) {
	var peerInfos []peerstore.PeerInfo
	for _, peer := range peers {
		ipfsaddr, err := ma.NewMultiaddr(peer)
		if err != nil {
			return nil, err
		}

		peerInfo, err := peerstore.InfoFromP2pAddr(ipfsaddr)
		if err != nil {
			return nil, err
		}

		peerInfos = append(peerInfos, *peerInfo)
	}
	return peerInfos, nil
}

func buildNotifiee() libp2pnet.Notifiee {
	notifyBundle := &libp2pnet.NotifyBundle{}

	notifyBundle.ConnectedF = func(_ libp2pnet.Network, connection libp2pnet.Conn) {
		logger.Infof(
			"established connection to [%v]",
			multiaddressWithIdentity(
				connection.RemoteMultiaddr(),
				connection.RemotePeer(),
			),
		)
	}
	notifyBundle.DisconnectedF = func(_ libp2pnet.Network, connection libp2pnet.Conn) {
		logger.Infof(
			"disconnected from [%v]",
			multiaddressWithIdentity(
				connection.RemoteMultiaddr(),
				connection.RemotePeer(),
			),
		)
	}

	return notifyBundle
}

func multiaddressWithIdentity(
	multiaddress ma.Multiaddr,
	peerID peer.ID,
) string {
	return fmt.Sprintf("%s/ipfs/%s", multiaddress.String(), peerID.String())
}
