package libp2p

import (
	"context"
	"fmt"
	"sync"

	"github.com/keep-network/keep-core/pkg/net"
	floodsub "github.com/libp2p/go-floodsub"
	host "github.com/libp2p/go-libp2p-host"
)

type channelManager struct {
	channelsMutex sync.Mutex
	channels      map[string]*channel

	pubsubMutex sync.Mutex
	pubsub      *floodsub.PubSub
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

func (cm *channelManager) getChannel(name string) (*channel, error) {
	cm.channelsMutex.Lock()
	channel, exists := cm.channels[name]
	cm.channelsMutex.Unlock()

	if !exists {
		return cm.newChannel(name)
	}

	return channel, nil
}

func verifyGroupName(name string) error {
	if name == "" {
		return fmt.Errorf("invalid channel name")
	}
	// TODO: some other conditions
	return nil
}

func (cm *channelManager) newChannel(name string) (*channel, error) {
	if err := verifyGroupName(name); err != nil {
		return nil, err
	}

	cm.pubsubMutex.Lock()
	sub, err := cm.pubsub.Subscribe(name)
	if err != nil {
		return nil, err
	}
	cm.pubsubMutex.Unlock()

	channel := &channel{
		name:                        name,
		sub:                         sub,
		unmarshalersByType:          make(map[string]func() net.TaggedUnmarshaler),
		transportToProtoIdentifiers: make(map[net.TransportIdentifier]net.ProtocolIdentifier),
		protoToTransportIdentifiers: make(map[net.ProtocolIdentifier]net.TransportIdentifier),
	}

	// Ensure we update our cache of known channels
	cm.channelsMutex.Lock()
	cm.channels[name] = channel
	cm.channelsMutex.Unlock()

	return channel, cm.joinChannel(name)
}

func (cm *channelManager) joinChannel(name string) error {
	return nil
}
