package libp2p

import (
	"context"
	"sync"

	"github.com/keep-network/keep-core/pkg/net"
	"github.com/libp2p/go-libp2p-core/host"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

type channelManager struct {
	ctx context.Context

	identity  *identity
	peerStore peerstore.Peerstore

	channelsMutex sync.Mutex
	channels      map[string]*channel

	pubsub *pubsub.PubSub

	retransmissionOptions *retransmissionOptions
}

func newChannelManager(
	ctx context.Context,
	identity *identity,
	p2phost host.Host,
	retransmissionOptions *retransmissionOptions,
) (*channelManager, error) {
	floodsub, err := pubsub.NewFloodSub(
		ctx,
		p2phost,
		pubsub.WithMessageAuthor(identity.id),
		pubsub.WithMessageSigning(true),
		pubsub.WithStrictSignatureVerification(true),
	)
	if err != nil {
		return nil, err
	}
	return &channelManager{
		channels:              make(map[string]*channel),
		pubsub:                floodsub,
		peerStore:             p2phost.Peerstore(),
		identity:              identity,
		ctx:                   ctx,
		retransmissionOptions: retransmissionOptions,
	}, nil
}

func (cm *channelManager) getChannel(name string) (*channel, error) {
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

func (cm *channelManager) newChannel(name string) (*channel, error) {
	sub, err := cm.pubsub.Subscribe(name)
	if err != nil {
		return nil, err
	}

	channel := &channel{
		name:               name,
		clientIdentity:     cm.identity,
		peerStore:          cm.peerStore,
		pubsub:             cm.pubsub,
		subscription:       sub,
		messageHandlers:    make([]*messageHandler, 0),
		unmarshalersByType: make(map[string]func() net.TaggedUnmarshaler),
		retransmitter: newRetransmitter(
			cm.retransmissionOptions.cycles,
			cm.retransmissionOptions.intervalMilliseconds,
		),
	}

	go channel.handleMessages(cm.ctx)

	return channel, nil
}
