package libp2p

import (
	"context"
	"sync"

	floodsub "github.com/libp2p/go-floodsub"
	host "github.com/libp2p/go-libp2p-host"
)

type channelManager struct {
	channelMutex sync.Mutex
	channels     map[string]*channel

	pubsub *floodsub.PubSub
	host   host.Host
}

func newChannelManager(
	ctx context.Context,
	h host.Host,
) (*channelManager, error) {
	gs, err := floodsub.NewGossipSub(ctx, h)
	if err != nil {
		return nil, err
	}
	cm := &channelManager{
		channels: make(map[string]*channel),
		pubsub:   gs,
		host:     h,
	}
	return cm, nil
}

func (cm *channelManager) getChannel(name string) *channel {
	cm.channelMutex.Lock()
	defer cm.channelMutex.Unlock()

	channel, exists := cm.channels[name]
	if !exists {
		// TODO: no topic exists; create the broadcast channel
		// TODO: return something informative ie. return cm.JoinChannel(name)
		return nil
	}
	return channel
}
