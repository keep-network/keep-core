package libp2p

import (
	"context"
	"sync"

	floodsub "github.com/libp2p/go-floodsub"
	host "github.com/libp2p/go-libp2p-host"
)

type channelManager struct {
	channeslMutex sync.Mutex
	channels      map[string]*channel

	pubsub *floodsub.PubSub
	host   host.Host
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
		host:     p2phost,
	}, nil
}

func (cm *channelManager) getChannel(name string) *channel {
	cm.channeslMutex.Lock()
	defer cm.channeslMutex.Unlock()

	channel, exists := cm.channels[name]
	if !exists {
		// TODO: no topic exists; create the broadcast channel
		// TODO: return something informative ie. return cm.JoinChannel(name)
		return nil
	}
	return channel
}
