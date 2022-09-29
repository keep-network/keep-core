// Package local provides a local, non-networked implementation of the
// interfaces defined by the net package. It should largely be considered a
// sample implementation, and is not meant to be used at scale in any way.
package local

import (
	"sync"

	"github.com/keep-network/keep-core/pkg/operator"

	"github.com/ipfs/go-log"
	"github.com/keep-network/keep-core/pkg/net"
)

var logger = log.Logger("keep-netlocal")

// Provider is an extension of net.Provider. This interface exposes additional
// functions useful for testing.
type Provider interface {
	net.Provider

	// AddPeer allows the simulation of adding a peer to the client's local
	// registry of peers.
	AddPeer(peerID string, publicKey *operator.PublicKey)
}

type localProvider struct {
	id                localIdentifier
	operatorPublicKey *operator.PublicKey
	connectionManager *localConnectionManager
}

func (lp *localProvider) ID() net.TransportIdentifier {
	return lp.id
}

func (lp *localProvider) BroadcastChannelFor(name string) (net.BroadcastChannel, error) {
	return getBroadcastChannel(name, lp.operatorPublicKey), nil
}

func (lp *localProvider) Type() string {
	return "local"
}

func (lp *localProvider) AddPeer(peerID string, publicKey *operator.PublicKey) {
	lp.connectionManager.peers[peerID] = publicKey
}

func (lp *localProvider) CreateTransportIdentifier(
	operatorPublicKey *operator.PublicKey,
) (
	net.TransportIdentifier,
	error,
) {
	return createLocalIdentifier(operatorPublicKey)
}

func (lp *localProvider) BroadcastChannelForwarderFor(name string) {
	//no-op
}

// Connect returns a local instance of a net provider that does not go over the
// network.
func Connect() Provider {
	_, operatorPublicKey, err := operator.GenerateKeyPair(DefaultCurve)
	if err != nil {
		panic(err)
	}

	return ConnectWithKey(operatorPublicKey)
}

// ConnectWithKey returns a local instance of net provider that does not go
// over the network. The returned instance uses the provided network key to
// identify network messages.
func ConnectWithKey(operatorPublicKey *operator.PublicKey) Provider {
	return &localProvider{
		id:                randomLocalIdentifier(),
		operatorPublicKey: operatorPublicKey,
		connectionManager: &localConnectionManager{peers: make(map[string]*operator.PublicKey)},
	}
}

func (lp *localProvider) ConnectionManager() net.ConnectionManager {
	return lp.connectionManager
}

type localConnectionManager struct {
	mutex sync.Mutex

	peers map[string]*operator.PublicKey
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

func (lcm *localConnectionManager) ConnectedPeersAddrInfo() map[string][]string {
	lcm.mutex.Lock()
	defer lcm.mutex.Unlock()
	var peersAddrInfo map[string][]string
	addresses := []string{"/ip4/localhost/"}
	for peer := range lcm.peers {
		peersAddrInfo[peer] = addresses
	}
	return peersAddrInfo
}

func (lcm *localConnectionManager) GetPeerPublicKey(
	connectedPeer string,
) (*operator.PublicKey, error) {
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
