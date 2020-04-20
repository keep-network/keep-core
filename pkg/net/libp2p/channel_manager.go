package libp2p

import (
	"context"
	"sync"

	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/net/retransmission"
	"github.com/libp2p/go-libp2p-core/host"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

const (
	libp2pMessageSigning              = true
	libp2pStrictSignatureVerification = true
	libp2pPeerOutboundQueueSize       = 256
	libp2pValidationQueueSize         = 4096
)

type channelManager struct {
	ctx context.Context

	identity  *identity
	peerStore peerstore.Peerstore

	channelsMutex sync.Mutex
	channels      map[string]*channel

	pubsub *pubsub.PubSub

	retransmissionTicker *retransmission.Ticker

	relaySubscriptionsMutex sync.Mutex
	relaySubscription       map[string]*pubsub.Subscription
}

func newChannelManager(
	ctx context.Context,
	identity *identity,
	p2phost host.Host,
	retransmissionTicker *retransmission.Ticker,
) (*channelManager, error) {
	floodsub, err := pubsub.NewFloodSub(
		ctx,
		p2phost,
		pubsub.WithMessageAuthor(identity.id),
		pubsub.WithMessageSigning(libp2pMessageSigning),
		pubsub.WithStrictSignatureVerification(libp2pStrictSignatureVerification),
		pubsub.WithPeerOutboundQueueSize(libp2pPeerOutboundQueueSize),
		pubsub.WithValidateQueueSize(libp2pValidationQueueSize),
	)
	if err != nil {
		return nil, err
	}
	return &channelManager{
		channels:             make(map[string]*channel),
		pubsub:               floodsub,
		peerStore:            p2phost.Peerstore(),
		identity:             identity,
		ctx:                  ctx,
		retransmissionTicker: retransmissionTicker,
		relaySubscription:    make(map[string]*pubsub.Subscription),
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
		// Ensure we update our cache of known channels
		cm.channelsMutex.Lock()
		defer cm.channelsMutex.Unlock()

		channel, exists = cm.channels[name]
		if exists {
			return channel, nil
		}

		channel, err = cm.newChannel(name)
		if err != nil {
			return nil, err
		}

		cm.channels[name] = channel
	}

	return channel, nil
}

func (cm *channelManager) newChannel(name string) (*channel, error) {
	sub, err := cm.pubsub.Subscribe(name)
	if err != nil {
		return nil, err
	}

	channel := &channel{
		name:                 name,
		clientIdentity:       cm.identity,
		peerStore:            cm.peerStore,
		pubsub:               cm.pubsub,
		subscription:         sub,
		incomingMessageQueue: make(chan *pubsub.Message, incomingMessageThrottle),
		messageHandlers:      make([]*messageHandler, 0),
		unmarshalersByType:   make(map[string]func() net.TaggedUnmarshaler),
		retransmissionTicker: cm.retransmissionTicker,
	}

	go channel.handleMessages(cm.ctx)

	return channel, nil
}

func (cm *channelManager) newRelay(name string) error {
	cm.relaySubscriptionsMutex.Lock()
	defer cm.relaySubscriptionsMutex.Unlock()

	if _, ok := cm.relaySubscription[name]; !ok {
		relaySubscription, err := cm.pubsub.Subscribe(name)
		if err != nil {
			return err
		}

		// TODO: invoke relaySubscription.Next() in a loop to avoid libp2p
		//  errors and make them context aware.

		cm.relaySubscription[name] = relaySubscription
	}

	return nil
}
