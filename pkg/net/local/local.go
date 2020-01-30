// Package local provides a local, non-networked implementation of the
// interfaces defined by the net package. It should largely be considered a
// sample implementation, and is not meant to be used at scale in any way.
package local

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/ipfs/go-log"

	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/net/gen/pb"
	"github.com/keep-network/keep-core/pkg/net/internal"
	"github.com/keep-network/keep-core/pkg/net/key"
	"github.com/keep-network/keep-core/pkg/net/retransmission"
)

var logger = log.Logger("keep-net-local")

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

func (lp *localProvider) ChannelWith(peerID net.TransportIdentifier) (net.UnicastChannel, error) {
	return nil, fmt.Errorf("not implemented")
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
		messageHandlers:      make([]*messageHandler, 0),
		unmarshalersMutex:    sync.Mutex{},
		unmarshalersByType:   make(map[string]func() net.TaggedUnmarshaler, 0),
		retransmissionTicker: retransmission.NewTimeTicker(
			context.Background(), 50*time.Millisecond,
		),
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

type messageHandler struct {
	ctx     context.Context
	channel chan retransmission.NetworkMessage
}

type localChannel struct {
	name                 string
	identifier           net.TransportIdentifier
	staticKey            *key.NetworkPublic
	messageHandlersMutex sync.Mutex
	messageHandlers      []*messageHandler
	unmarshalersMutex    sync.Mutex
	unmarshalersByType   map[string]func() net.TaggedUnmarshaler
	retransmissionTicker *retransmission.Ticker
}

func (lc *localChannel) Name() string {
	return lc.name
}

func doSend(channel *localChannel, message *pb.NetworkMessage) error {
	channelsMutex.Lock()
	targetChannels := channels[channel.name]
	channelsMutex.Unlock()

	unmarshaler, found := channel.unmarshalersByType[string(message.Type)]
	if !found {
		return fmt.Errorf("couldn't find unmarshaler for type %s", string(message.Type))
	}

	unmarshaled := unmarshaler()
	err := unmarshaled.Unmarshal(message.Payload)
	if err != nil {
		return err
	}

	fingerprint := retransmission.CalculateFingerprint(
		channel.identifier,
		message.Payload,
	)

	protocolMessage := retransmission.NewNetworkMessage(
		internal.BasicMessage(
			channel.identifier,
			unmarshaled,
			"local",
			key.Marshal(channel.staticKey),
		),
		fingerprint,
		message.Retransmission,
	)

	for _, targetChannel := range targetChannels {
		targetChannel.deliver(protocolMessage)
	}

	return nil
}

func (lc *localChannel) deliver(message retransmission.NetworkMessage) {
	lc.messageHandlersMutex.Lock()
	snapshot := make([]*messageHandler, len(lc.messageHandlers))
	copy(snapshot, lc.messageHandlers)
	lc.messageHandlersMutex.Unlock()

	for _, handler := range snapshot {
		go func(message retransmission.NetworkMessage, handler *messageHandler) {
			select {
			case handler.channel <- message:
			// Nothing to do here; we block until the message is handled
			// or until the context gets closed.
			// This way we don't lose any message but also don't stay
			// with any dangling goroutines if there is no longer anyone
			// to receive messages.
			case <-handler.ctx.Done():
				return
			}
		}(message, handler)
	}
}

func (lc *localChannel) Send(ctx context.Context, message net.TaggedMarshaler) error {
	messageProto, err := lc.messageProto(message)
	if err != nil {
		return err
	}

	retransmission.ScheduleRetransmissions(
		ctx,
		lc.retransmissionTicker,
		messageProto,
		func(msg *pb.NetworkMessage) error {
			return doSend(lc, msg)
		},
	)

	return doSend(lc, messageProto)
}

func (lc *localChannel) messageProto(
	message net.TaggedMarshaler,
) (*pb.NetworkMessage, error) {
	payloadBytes, err := message.Marshal()
	if err != nil {
		return nil, err
	}

	return &pb.NetworkMessage{
		Payload: payloadBytes,
		Sender:  []byte(lc.identifier.String()),
		Type:    []byte(message.Type()),
	}, nil
}

func (lc *localChannel) Recv(ctx context.Context, handler func(m net.Message)) {
	messageHandler := &messageHandler{
		ctx:     ctx,
		channel: make(chan retransmission.NetworkMessage),
	}

	lc.messageHandlersMutex.Lock()
	lc.messageHandlers = append(lc.messageHandlers, messageHandler)
	lc.messageHandlersMutex.Unlock()

	handleWithRetransmissions := retransmission.WithRetransmissionSupport(handler)

	go func() {
		for {
			select {
			case <-ctx.Done():
				logger.Debug("context is done, removing handler")
				lc.removeHandler(messageHandler)
				return

			case msg := <-messageHandler.channel:
				// Go language specification says that if one or more of the
				// communications in the select statement can proceed, a single
				// one that will proceed is chosen via a uniform pseudo-random
				// selection.
				// Thus, it can happen this communication is called when ctx is
				// already done. Since we guarantee in the network channel API
				// that handler is not called after ctx is done (client code
				// could e.g. perform come cleanup), we need to double-check
				// the context state here.
				if messageHandler.ctx.Err() != nil {
					continue
				}

				handleWithRetransmissions(msg)
			}
		}
	}()
}

func (lc *localChannel) removeHandler(handler *messageHandler) {
	lc.messageHandlersMutex.Lock()
	defer lc.messageHandlersMutex.Unlock()

	for i, h := range lc.messageHandlers {
		if h.channel == handler.channel {
			lc.messageHandlers[i] = lc.messageHandlers[len(lc.messageHandlers)-1]
			lc.messageHandlers = lc.messageHandlers[:len(lc.messageHandlers)-1]
		}
	}
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

func (lc *localChannel) SetFilter(filter net.BroadcastChannelFilter) error {
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
