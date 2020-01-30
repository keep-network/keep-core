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
	libp2pPeerOutboundQueueSize       = 128
	libp2pValidationQueueSize         = 128
)

type channelManager struct {
	ctx context.Context

	identity  *identity
	peerStore peerstore.Peerstore

	channelsMutex sync.Mutex
	channels      map[string]*channel

	pubsub *pubsub.PubSub

	retransmissionTicker *retransmission.Ticker
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
	}, nil
}

func (cm *channelManager) getChannel(name string) (*channel, error) {
	var (
		channel *channel
		exists  bool
	)

	cm.channelsMutex.Lock()
	channel, exists = cm.channels[name]
	cm.channelsMutex.Unlock()

	if !exists {
		newChannel, err := cm.newChannel(name)
		if err != nil {
			return nil, err
		}

		// Creating a new channel can take some time. One should double-check
		// if some other channel wasn't created and cached in the same time.
		cm.channelsMutex.Lock()
		channel, exists = cm.channels[name]
		if !exists {
			channel = newChannel
			cm.channels[name] = newChannel
		}
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
		name:                 name,
		clientIdentity:       cm.identity,
		peerStore:            cm.peerStore,
		pubsub:               cm.pubsub,
		subscription:         sub,
		messageHandlers:      make([]*messageHandler, 0),
		unmarshalersByType:   make(map[string]func() net.TaggedUnmarshaler),
		retransmissionTicker: cm.retransmissionTicker,
	}

	go channel.handleMessages(cm.ctx)

	return channel, nil
}
