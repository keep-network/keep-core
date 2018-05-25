package libp2p

import (
	"context"
	"sync"

	"github.com/keep-network/keep-core/pkg/net"
	floodsub "github.com/libp2p/go-floodsub"
	host "github.com/libp2p/go-libp2p-host"
)

type channelManager struct {
	channelsMutex sync.Mutex
	channels      map[string]*channel

	pubsub *floodsub.PubSub
}

func newChannelManager(
	ctx context.Context,
	p2phost host.Host,
) (*channelManager, error) {
	gossipsub, err := floodsub.NewGossipSub(ctx, p2phost)
	if err != nil {
		return nil, err
	}
	return &channelManager{
		channels: make(map[string]*channel),
		pubsub:   gossipsub,
	}, nil
}

func (cm *channelManager) getChannel(name string) *channel {
	cm.channelsMutex.Lock()
	defer cm.channelsMutex.Unlock()

	channel, exists := cm.channels[name]
	if !exists {
		// TODO: no topic exists; create the broadcast channel
		// TODO: return something informative ie. return cm.JoinChannel(name)
		return nil
	}
	return channel
}

func (cm *channelManager) newChannel(name string) (*channel, error) {
	sub, err := cm.pubsub.Subscribe(name)
	if err != nil {
		return nil, err
	}

	channel := &channel{
		name:                        name,
		sub:                         sub,
		unmarshalersByType:          make(map[string]func() net.TaggedUnmarshaler, 0),
		transportToProtoIdentifiers: make(map[net.TransportIdentifier]net.ProtocolIdentifier),
		protoToTransportIdentifiers: make(map[net.ProtocolIdentifier]net.TransportIdentifier),
	}

	return channel, cm.joinChannel(name)
}

func (cm *channelManager) joinChannel(name string) error {
	return nil
}
