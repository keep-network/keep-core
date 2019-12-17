// Package local provides a local, non-networked implementation of the
// interfaces defined by the net package. It should largely be considered a
// sample implementation, and is not meant to be used at scale in any way.
package local

import (
	"fmt"
	"math/rand"
	"sync"

	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/net/internal"
	"github.com/keep-network/keep-core/pkg/net/key"
)

type localIdentifier string

func (li localIdentifier) String() string {
	return string(li)
}

var channelsMutex sync.Mutex
var channels map[string][]*localChannel

// Provider is an extension of net.Provider. This interface exposes additional
// functions useful for testing.
type Provider interface {
	net.Provider

	// AddPeer allows the simulation of adding a peer to the client's local
	// registry of peers.
	AddPeer(peerID string, pubKey *key.NetworkPublic)
}

type localProvider struct {
	id        localIdentifier
	staticKey *key.NetworkPublic
	cm        *localConnectionManager
}

func (lp *localProvider) ID() net.TransportIdentifier {
	return lp.id
}

func (lp *localProvider) ChannelFor(name string) (net.BroadcastChannel, error) {
	return channel(name, lp.staticKey), nil
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
	lp.cm.peers[peerID] = pubKey
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
		id:        localIdentifier(randomIdentifier()),
		staticKey: staticKey,
		cm:        &localConnectionManager{peers: make(map[string]*key.NetworkPublic)},
	}
}

func (lp *localProvider) ConnectionManager() net.ConnectionManager {
	return lp.cm
}

// channel returns a BroadcastChannel designed to mediate between local
// participants. It delivers all messages sent to the channel through its
// receive channels. RecvChan on a LocalChannel creates a new receive channel
// that is returned to the caller, so that all receive channels can receive
// the message.
func channel(name string, staticKey *key.NetworkPublic) net.BroadcastChannel {
	channelsMutex.Lock()
	defer channelsMutex.Unlock()
	if channels == nil {
		channels = make(map[string][]*localChannel)
	}

	localChannels, exists := channels[name]
	if !exists {
		localChannels = make([]*localChannel, 0)
		channels[name] = localChannels
	}

	identifier := localIdentifier(randomIdentifier())
	channel := &localChannel{
		name:                 name,
		identifier:           &identifier,
		staticKey:            staticKey,
		messageHandlersMutex: sync.Mutex{},
		messageHandlers:      make([]net.HandleMessageFunc, 0),
		unmarshalersMutex:    sync.Mutex{},
		unmarshalersByType:   make(map[string]func() net.TaggedUnmarshaler, 0),
	}
	channels[name] = append(channels[name], channel)

	return channel
}

var letterRunes = [52]rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j',
	'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y',
	'z', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N',
	'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}

func randomIdentifier() string {
	runes := make([]rune, 32)
	for i := range runes {
		runes[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(runes)
}

type localChannel struct {
	name                 string
	identifier           net.TransportIdentifier
	staticKey            *key.NetworkPublic
	messageHandlersMutex sync.Mutex
	messageHandlers      []net.HandleMessageFunc
	unmarshalersMutex    sync.Mutex
	unmarshalersByType   map[string]func() net.TaggedUnmarshaler
}

func (lc *localChannel) Name() string {
	return lc.name
}

func doSend(channel *localChannel, payload net.TaggedMarshaler) error {
	channelsMutex.Lock()
	targetChannels := channels[channel.name]
	channelsMutex.Unlock()

	bytes, err := payload.Marshal()
	if err != nil {
		return err
	}

	unmarshaler, found := channel.unmarshalersByType[payload.Type()]
	if !found {
		return fmt.Errorf("Couldn't find unmarshaler for type %s", payload.Type())
	}

	unmarshaled := unmarshaler()
	err = unmarshaled.Unmarshal(bytes)
	if err != nil {
		return err
	}

	for _, targetChannel := range targetChannels {
		targetChannel.deliver(channel.identifier, channel.staticKey, unmarshaled) // TODO error handling?
	}

	return nil
}

func (lc *localChannel) deliver(
	senderTransportIdentifier net.TransportIdentifier,
	senderPublicKey *key.NetworkPublic,
	payload interface{},
) {
	lc.messageHandlersMutex.Lock()
	snapshot := make([]net.HandleMessageFunc, len(lc.messageHandlers))
	copy(snapshot, lc.messageHandlers)
	lc.messageHandlersMutex.Unlock()

	message :=
		internal.BasicMessage(
			senderTransportIdentifier,
			payload,
			"local",
			key.Marshal(senderPublicKey),
		)

	for _, handler := range snapshot {
		// Invoking each handler in a separate goroutine in order
		// to avoid situation when one of the handlers blocks the
		// whole loop execution.
		go handler.Handler(message)
	}
}

func (lc *localChannel) Send(
	message net.TaggedMarshaler,
	retransmission ...net.RetransmissionOptions,
) error {
	return doSend(lc, message)
}

func (lc *localChannel) Recv(handler net.HandleMessageFunc) error {
	lc.messageHandlersMutex.Lock()
	lc.messageHandlers = append(lc.messageHandlers, handler)
	lc.messageHandlersMutex.Unlock()

	return nil
}

func (lc *localChannel) UnregisterRecv(handlerType string) error {
	lc.messageHandlersMutex.Lock()
	defer lc.messageHandlersMutex.Unlock()

	removedCount := 0

	// updated slice shares the same backing array and capacity as the original,
	// so the storage is reused for the filtered slice.
	updated := lc.messageHandlers[:0]

	for _, mh := range lc.messageHandlers {
		if mh.Type != handlerType {
			updated = append(updated, mh)
		} else {
			removedCount++
		}
	}

	lc.messageHandlers = updated[:len(lc.messageHandlers)-removedCount]

	return nil
}

func (lc *localChannel) RegisterUnmarshaler(
	unmarshaler func() net.TaggedUnmarshaler,
) (err error) {
	tpe := unmarshaler().Type()

	lc.unmarshalersMutex.Lock()
	_, exists := lc.unmarshalersByType[tpe]
	if exists {
		err = fmt.Errorf("type %s already has an associated unmarshaler", tpe)
	} else {
		lc.unmarshalersByType[tpe] = unmarshaler
	}
	lc.unmarshalersMutex.Unlock()
	return
}

func (lc *localChannel) AddFilter(filter net.BroadcastChannelFilter) error {
	return nil // no-op
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
