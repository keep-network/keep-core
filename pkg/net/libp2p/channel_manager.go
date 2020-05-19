package libp2p

import (
	"context"
	"sync"
	"time"

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

	forwarderSubscriptionsMutex sync.Mutex
	forwarderSubscriptions      map[string]*pubsub.Subscription
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
		channels:               make(map[string]*channel),
		pubsub:                 floodsub,
		peerStore:              p2phost.Peerstore(),
		identity:               identity,
		ctx:                    ctx,
		retransmissionTicker:   retransmissionTicker,
		forwarderSubscriptions: make(map[string]*pubsub.Subscription),
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

func (cm *channelManager) newForwarder(name string, ttl time.Duration) error {
	cm.forwarderSubscriptionsMutex.Lock()
	defer cm.forwarderSubscriptionsMutex.Unlock()

	if _, ok := cm.forwarderSubscriptions[name]; !ok {
		forwarderSubscription, err := cm.pubsub.Subscribe(name)
		if err != nil {
			return err
		}

		go func() {
			ctx, cancelCtx := context.WithTimeout(cm.ctx, ttl)
			defer cancelCtx()

			for {
				select {
				case <-ctx.Done():
					cm.shutdownForwarder(name)
					return
				default:
					// Just pull the message from subscription to unblock
					// the channel and avoid warnings from libp2p. We
					// are not interested with their content.
					_, _ = forwarderSubscription.Next(ctx)
				}
			}
		}()

		cm.forwarderSubscriptions[name] = forwarderSubscription
	}

	return nil
}

func (cm *channelManager) shutdownForwarder(name string) {
	cm.forwarderSubscriptionsMutex.Lock()
	defer cm.forwarderSubscriptionsMutex.Unlock()

	logger.Infof("shutting down message forwarder for channel: [%v]", name)

	forwarderSubscription, ok := cm.forwarderSubscriptions[name]

	if !ok {
		return
	}

	forwarderSubscription.Cancel()
	delete(cm.forwarderSubscriptions, name)
}
