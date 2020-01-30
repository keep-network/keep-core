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

func (ucm *unicastChannelManager) handleIncomingStream(stream network.Stream) {
	logger.Debugf(
		"[%v] processing incoming stream [%v] from peer [%v]",
		ucm.identity.id,
		stream.Protocol(),
		stream.Conn().RemotePeer(),
	)

	channel, err := ucm.getUnicastChannel(stream.Conn().RemotePeer())
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

func (ucm *unicastChannelManager) getUnicastChannel(peerID net.TransportIdentifier) (
	*unicastChannel,
	error,
) {
	var (
		channel *unicastChannel
		exists  bool
		err     error
	)

	ucm.channelsMutex.Lock()
	channel, exists = ucm.channels[peerID]
	ucm.channelsMutex.Unlock()

	if !exists {
		channel, err = ucm.newUnicastChannel(peerID)
		if err != nil {
			return nil, err
		}

		// Ensure we update our cache of known channels
		ucm.channelsMutex.Lock()
		ucm.channels[peerID] = channel
		ucm.channelsMutex.Unlock()
	}

	return channel, nil
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
