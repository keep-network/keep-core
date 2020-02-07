package libp2p

import (
	"context"
	"fmt"
	"sync"

	"github.com/keep-network/keep-core/pkg/net"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	peer "github.com/libp2p/go-libp2p-core/peer"
)

const protocolID = "/keep/unicast/1.0.0"

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

	channel, isExistingChannel, err := ucm.getUnicastChannel(stream.Conn().RemotePeer())
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

	if !isExistingChannel && ucm.channelOpenedHandler != nil {
		ucm.channelOpenedHandler(channel)
	}

	channel.handleStream(stream)
}

func (ucm *unicastChannelManager) getUnicastChannel(
	peerID net.TransportIdentifier,
) (
	*unicastChannel,
	bool,
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
		newChannel, err := ucm.newUnicastChannel(peerID)
		if err != nil {
			return nil, exists, err
		}

		// Creating a new channel can take some time. One should double-check
		// if some other channel wasn't created and cached in the same time.
		ucm.channelsMutex.Lock()
		channel, exists = ucm.channels[peerID]
		if !exists {
			channel = newChannel
			ucm.channels[peerID] = newChannel
		}
		ucm.channelsMutex.Unlock()
	}

	return channel, exists, nil
}

func (ucm *unicastChannelManager) newUnicastChannel(
	peerID net.TransportIdentifier,
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

	return channel, nil
}
