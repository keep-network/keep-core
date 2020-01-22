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

	streamsMutex sync.Mutex
	streams      map[peer.ID]network.Stream
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
		streams:  make(map[peer.ID]network.Stream),
	}

	p2phost.SetStreamHandlerMatch(
		protocolID,
		func(protocol string) bool { return protocol == protocolID },
		manager.handleIncomingStream,
	)

	notifyBundle := &network.NotifyBundle{}
	notifyBundle.ClosedStreamF = func(_ network.Network, stream network.Stream) {
		if stream.Protocol() == protocolID {
			manager.handleClosedStream(stream)
		}
	}
	p2phost.Network().Notify(notifyBundle)

	return manager
}

func (ucm *unicastChannelManager) handleIncomingStream(stream network.Stream) {
	logger.Debugf(
		"[%v] new incoming stream [%v] from peer [%v]",
		ucm.identity.id,
		stream.Protocol(),
		stream.Conn().RemotePeer(),
	)

	ucm.streamsMutex.Lock()
	defer ucm.streamsMutex.Unlock()

	ucm.streams[stream.Conn().RemotePeer()] = stream

	logger.Debugf(
		"[%v] stream [%v] with peer [%v] registered successfully",
		ucm.identity.id,
		stream.Protocol(),
		stream.Conn().RemotePeer(),
	)
}

func (ucm *unicastChannelManager) handleClosedStream(stream network.Stream) {
	logger.Debugf(
		"[%v] detected closed stream [%v] with peer [%v]",
		ucm.identity.id,
		stream.Protocol(),
		stream.Conn().RemotePeer(),
	)

	ucm.streamsMutex.Lock()
	defer ucm.streamsMutex.Unlock()

	delete(ucm.streams, stream.Conn().RemotePeer())

	logger.Debugf(
		"[%v] stream [%v] with peer [%v] unregistered successfully",
		ucm.identity.id,
		stream.Protocol(),
		stream.Conn().RemotePeer(),
	)
}

func (ucm *unicastChannelManager) getUnicastChannel(
	peerID string,
) (*unicastChannel, error) {
	ucm.streamsMutex.Lock()
	defer ucm.streamsMutex.Unlock()

	remotePeer, err := peer.IDB58Decode(peerID)
	if err != nil {
		return nil, fmt.Errorf("invalid peer ID: [%v]", err)
	}

	stream, exists := ucm.streams[remotePeer]
	if !exists {
		logger.Debugf(
			"[%v] creating stream [%v] with peer [%v]",
			ucm.identity.id,
			protocolID,
			peerID,
		)

		stream, err = ucm.p2phost.NewStream(ucm.ctx, remotePeer, protocolID)
		if err != nil {
			return nil, fmt.Errorf(
				"could not create stream [%v] with peer [%v]: [%v]",
				protocolID,
				peerID,
				err,
			)
		}
	} else {
		logger.Debugf(
			"[%v] using existing stream [%v] with peer [%v]",
			ucm.identity.id,
			stream.Protocol(),
			peerID,
		)
	}

	channel := &unicastChannel{
		clientIdentity:     ucm.identity,
		stream:             stream,
		messageHandlers:    make([]*messageHandler, 0),
		unmarshalersByType: make(map[string]func() net.TaggedUnmarshaler),
	}

	go channel.handleMessages(ucm.ctx)

	return channel, nil
}
