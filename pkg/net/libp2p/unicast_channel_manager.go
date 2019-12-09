package libp2p

import (
	"context"
	"fmt"
	"sync"

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

	p2phost.SetStreamHandler(protocolID, manager.handleIncomingStream)

	notifyBundle := &network.NotifyBundle{}
	notifyBundle.ClosedStreamF = func(_ network.Network, stream network.Stream) {
		manager.handleClosedStream(stream)
	}
	p2phost.Network().Notify(notifyBundle)

	return manager
}

func (ucm *unicastChannelManager) handleIncomingStream(stream network.Stream) {
	logger.Warningf("[%v] new incoming stream from [%v]", ucm.identity.id, stream.Conn().RemotePeer())
	ucm.streamsMutex.Lock()
	ucm.streams[stream.Conn().RemotePeer()] = stream
	logger.Warningf("stream registered!")
	ucm.streamsMutex.Unlock()
}

func (ucm *unicastChannelManager) handleClosedStream(stream network.Stream) {
	ucm.streamsMutex.Lock()
	delete(ucm.streams, stream.Conn().RemotePeer())
	ucm.streamsMutex.Unlock()
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
		logger.Warningf("[%v] creating new stream for remote peer [%v]", ucm.identity.id, peerID)
		stream, err = ucm.p2phost.NewStream(ucm.ctx, remotePeer, protocolID)
		if err != nil {
			return nil, fmt.Errorf("could not create stream: [%v]", err)
		}
	} else {
		logger.Warningf("[%v] using existing stream for remote peer [%v]", ucm.identity.id, peerID)
	}

	channel := newUnicastChannel(ucm.ctx, ucm.identity, stream)
	go channel.handleMessages(ucm.ctx)

	return channel, nil
}
