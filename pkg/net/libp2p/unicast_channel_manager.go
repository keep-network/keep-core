package libp2p

import (
	"context"
	"fmt"
	"sync"

	"github.com/keep-network/keep-core/pkg/net"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
)

const protocolID = "/keep/unicast/1.0.0"

type ChannelInitDirection int

const (
	Inbound ChannelInitDirection = iota
	Outbound
)

type unicastChannelManager struct {
	ctx context.Context

	identity *identity
	p2phost  host.Host

	channelsMutex sync.Mutex
	channels      map[net.TransportIdentifier]*unicastChannel

	channelOpenedHandler func(channel net.UnicastChannel)
}

func newUnicastChannelManager(
	ctx context.Context,
	identity *identity,
	p2phost host.Host,
) *unicastChannelManager {
	manager := &unicastChannelManager{
		ctx:      ctx,
		identity: identity,
		p2phost:  p2phost,
		channels: make(map[net.TransportIdentifier]*unicastChannel),
	}

	p2phost.SetStreamHandlerMatch(
		protocolID,
		func(protocol string) bool { return protocol == protocolID },
		manager.handleIncomingStream,
	)

	return manager
}

func (ucm *unicastChannelManager) onChannelOpened(
	handler func(channel net.UnicastChannel),
) {
	ucm.channelOpenedHandler = handler
}

func (ucm *unicastChannelManager) handleIncomingStream(stream network.Stream) {
	logger.Debugf(
		"[%v] processing incoming stream [%v] from peer [%v]",
		ucm.identity.id,
		stream.Protocol(),
		stream.Conn().RemotePeer(),
	)

	channel, err := ucm.getUnicastChannel(stream.Conn().RemotePeer(), Inbound)
	if err != nil {
		logger.Errorf(
			"[%v] incoming stream [%v] from peer [%v] dropped: [%v]",
			ucm.identity.id,
			stream.Protocol(),
			stream.Conn().RemotePeer(),
			err,
		)
		return
	}

	channel.handleStream(stream)
}

func (ucm *unicastChannelManager) getUnicastChannel(
	peerID net.TransportIdentifier,
	initDirection ChannelInitDirection,
) (
	*unicastChannel,
	error,
) {
	var (
		channel *unicastChannel
		exists  bool
	)

	ucm.channelsMutex.Lock()
	channel, exists = ucm.channels[peerID]
	ucm.channelsMutex.Unlock()

	if !exists {
		triggerHandshake := initDirection == Outbound
		newChannel, err := ucm.newUnicastChannel(peerID, triggerHandshake)
		if err != nil {
			return nil, err
		}

		// Creating a new channel can take some time. One should double-check
		// if some other channel wasn't created and cached in the same time.
		ucm.channelsMutex.Lock()
		channel, exists = ucm.channels[peerID]
		if !exists {
			channel = newChannel
			ucm.channels[peerID] = newChannel

			// One should invoke the channel opened handler only in case
			// when the new channel was initiated by the remote peer.
			if initDirection == Inbound && ucm.channelOpenedHandler != nil {
				ucm.channelOpenedHandler(newChannel)
			}
		}
		ucm.channelsMutex.Unlock()
	}

	return channel, nil
}

func (ucm *unicastChannelManager) newUnicastChannel(
	peerID net.TransportIdentifier,
	triggerHandshake bool,
) (*unicastChannel, error) {
	remotePeer, err := peer.IDB58Decode(peerID.String())
	if err != nil {
		return nil, fmt.Errorf("invalid peer ID: [%v]", err)
	}

	streamFactory := func(ctx context.Context, peerID peer.ID) (network.Stream, error) {
		return ucm.p2phost.NewStream(ctx, peerID, protocolID)
	}

	channel := &unicastChannel{
		clientIdentity:     ucm.identity,
		remotePeerID:       remotePeer,
		streamFactory:      streamFactory,
		messageHandlers:    make([]*unicastMessageHandler, 0),
		unmarshalersByType: make(map[string]func() net.TaggedUnmarshaler),
	}

	if triggerHandshake {
		// Trigger a handshake in order to check if peer is reachable.
		err = channel.handshake()
		if err != nil {
			return nil, fmt.Errorf("handshake failed: [%v]", err)
		}
	}

	return channel, nil
}
