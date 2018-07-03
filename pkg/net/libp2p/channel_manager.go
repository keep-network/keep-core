package libp2p

import (
	"context"
	"sync"

	"github.com/keep-network/keep-core/pkg/net"
	floodsub "github.com/libp2p/go-floodsub"
	host "github.com/libp2p/go-libp2p-host"
	"github.com/libp2p/go-libp2p-peerstore"
)

// ChannelManager providing peer identity, peerStore, and messaging.
type ChannelManager struct {
	ctx context.Context

	identity  *Identity
	peerStore peerstore.Peerstore

	channelsMutex sync.Mutex
	channels      map[string]*channel

	pubsub *floodsub.PubSub
}

// NewChannelManager returns a new ChannelManager
func NewChannelManager(
	ctx context.Context,
	identity *Identity,
	p2phost host.Host,
) (*ChannelManager, error) {
	gossipsub, err := floodsub.NewGossipSub(ctx, p2phost)
	if err != nil {
		return nil, err
	}
	return &ChannelManager{
		channels:  make(map[string]*channel),
		pubsub:    gossipsub,
		peerStore: p2phost.Peerstore(),
		identity:  identity,
		ctx:       ctx,
	}, nil
}

func (cm *ChannelManager) getChannel(name string) (*channel, error) {
	var (
		channel *channel
		exists  bool
		err     error
	)

	cm.channelsMutex.Lock()
	channel, exists = cm.channels[name]
	cm.channelsMutex.Unlock()

	if !exists {
		channel, err = cm.newChannel(name)
		if err != nil {
			return nil, err
		}

		// Ensure we update our cache of known channels
		cm.channelsMutex.Lock()
		cm.channels[name] = channel
		cm.channelsMutex.Unlock()
	}

	return channel, nil
}

func (cm *ChannelManager) newChannel(name string) (*channel, error) {
	sub, err := cm.pubsub.Subscribe(name)
	if err != nil {
		return nil, err
	}

	channel := &channel{
		name:                        name,
		clientIdentity:              cm.identity,
		peerStore:                   cm.peerStore,
		pubsub:                      cm.pubsub,
		subscription:                sub,
		messageHandlers:             make([]net.HandleMessageFunc, 0),
		unmarshalersByType:          make(map[string]func() net.TaggedUnmarshaler),
		transportToProtoIdentifiers: make(map[net.TransportIdentifier]net.ProtocolIdentifier),
		protoToTransportIdentifiers: make(map[net.ProtocolIdentifier]net.TransportIdentifier),
	}

	go channel.handleMessages(cm.ctx)

	return channel, nil
}
