// Package local provides a local, non-networked implementation of the
// interfaces defined by the net package. It should largely be considered a
// sample implementation, and is not meant to be used at scale in any way.
package local

import (
	"fmt"
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
	id                localIdentifier
	staticKey         *key.NetworkPublic
	connectionManager *localConnectionManager
}

func (lp *localProvider) ID() net.TransportIdentifier {
	return lp.id
}

func (lp *localProvider) UnicastChannelWith(peerID net.TransportIdentifier) (
	net.UnicastChannel,
	error,
) {
	return nil, fmt.Errorf("not implemented")
}

func (lp *localProvider) OnUnicastChannelOpened(
	handler func(channel net.UnicastChannel),
) {
	// no-op
}

func (lp *localProvider) BroadcastChannelFor(name string) (net.BroadcastChannel, error) {
	return getBroadcastChannel(name, lp.staticKey), nil
}

func (lp *localProvider) Type() string {
	return "local"
}

func (lp *localProvider) AddrStrings() []string {
	return make([]string, 0)
}

func (lp *localProvider) Peers() []string {
	return make([]string, 0)
}

func (lp *localProvider) AddPeer(peerID string, pubKey *key.NetworkPublic) {
	lp.connectionManager.peers[peerID] = pubKey
}

// Connect returns a local instance of a net provider that does not go over the
// network.
func Connect() Provider {
	_, public, err := key.GenerateStaticNetworkKey()
	if err != nil {
		panic(err)
	}

	return ConnectWithKey(public)
}

// ConnectWithKey returns a local instance of net provider that does not go
// over the network. The returned instance uses the provided network key to
// identify network messages.
func ConnectWithKey(staticKey *key.NetworkPublic) Provider {
	return &localProvider{
		id:                randomIdentifier(),
		staticKey:         staticKey,
		connectionManager: &localConnectionManager{peers: make(map[string]*key.NetworkPublic)},
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
	connectedPeers := make([]string, len(lcm.peers))
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
