// Package local provides a local, non-networked implementation of the
// interfaces defined by the net package. It should largely be considered a
// sample implementation, and is not meant to be used at scale in any way.
package local

import (
	"context"
	"crypto/ecdsa"
	"sync"

	"github.com/ipfs/go-log"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/net/key"
)

var logger = log.Logger("keep-net-local")

// Provider is an extension of net.Provider. This interface exposes additional
// functions useful for testing.
type Provider interface {
	net.Provider

	// AddPeer allows the simulation of adding a peer to the client's local
	// registry of peers.
	AddPeer(peerID string, pubKey *key.NetworkPublic)
}

type localProvider struct {
	id                    localIdentifier
	staticKey             *key.NetworkPublic
	connectionManager     *localConnectionManager
	unicastChannelManager *unicastChannelManager
}

func (lp *localProvider) ID() net.TransportIdentifier {
	return lp.id
}

func (lp *localProvider) UnicastChannelWith(peerID net.TransportIdentifier) (
	net.UnicastChannel,
	error,
) {
	return lp.unicastChannelManager.UnicastChannelWith(peerID)
}

func (lp *localProvider) OnUnicastChannelOpened(
	handler func(channel net.UnicastChannel),
) {
	lp.unicastChannelManager.OnUnicastChannelOpened(context.Background(), handler)
}

func (lp *localProvider) BroadcastChannelFor(name string) (net.BroadcastChannel, error) {
	return getBroadcastChannel(name, lp.staticKey), nil
}

func (lp *localProvider) Type() string {
	return "local"
}

func (lp *localProvider) AddPeer(peerID string, pubKey *key.NetworkPublic) {
	lp.connectionManager.peers[peerID] = pubKey
}

func (lp *localProvider) CreateTransportIdentifier(publicKey ecdsa.PublicKey) (
	net.TransportIdentifier,
	error,
) {
	networkPublicKey := key.ECDSAKeyToNetworkKey(&publicKey)
	return createLocalIdentifier(networkPublicKey), nil
}

func (lp *localProvider) BroadcastChannelForwarderFor(name string) {
	//no-op
}

// Connect returns a local instance of a net provider that does not go over the
// network.
func Connect() Provider {
	_, staticKey, err := key.GenerateStaticNetworkKey()
	if err != nil {
		panic(err)
	}

	return ConnectWithKey(staticKey)
}

// ConnectWithKey returns a local instance of net provider that does not go
// over the network. The returned instance uses the provided network key to
// identify network messages.
func ConnectWithKey(staticKey *key.NetworkPublic) Provider {
	return &localProvider{
		id:                    randomLocalIdentifier(),
		staticKey:             staticKey,
		connectionManager:     &localConnectionManager{peers: make(map[string]*key.NetworkPublic)},
		unicastChannelManager: newUnicastChannelManager(staticKey),
	}
}

func (lp *localProvider) ConnectionManager() net.ConnectionManager {
	return lp.connectionManager
}

type localConnectionManager struct {
	mutex sync.Mutex

	peers map[string]*key.NetworkPublic
}

func (lcm *localConnectionManager) ConnectedPeers() []string {
	lcm.mutex.Lock()
	defer lcm.mutex.Unlock()
	var connectedPeers []string
	for peer := range lcm.peers {
		connectedPeers = append(connectedPeers, peer)
	}
	return connectedPeers
}

func (lcm *localConnectionManager) GetPeerPublicKey(
	connectedPeer string,
) (*key.NetworkPublic, error) {
	lcm.mutex.Lock()
	defer lcm.mutex.Unlock()

	return lcm.peers[connectedPeer], nil
}

func (lcm *localConnectionManager) DisconnectPeer(connectedPeer string) {
	lcm.mutex.Lock()
	defer lcm.mutex.Unlock()

	delete(lcm.peers, connectedPeer)
}

func (lcm *localConnectionManager) AddrStrings() []string {
	return make([]string, 0)
}

func (lcm *localConnectionManager) IsConnected(address string) bool {
	panic("not implemented")
}
