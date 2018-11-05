package libp2p

import (
	"context"
	"fmt"
	"time"

	"github.com/keep-network/keep-core/pkg/net/gen/pb"
	inet "github.com/libp2p/go-libp2p-net"
	ma "github.com/multiformats/go-multiaddr"
)

var _ inet.Notifiee = (*channelNotifiee)(nil)

type channelNotifiee = channel

func (cn *channelNotifiee) Listen(n inet.Network, a ma.Multiaddr) {
}
func (cn *channelNotifiee) ListenClose(n inet.Network, a ma.Multiaddr) {
}
func (cn *channelNotifiee) Connected(n inet.Network, v inet.Conn) {
	fmt.Println("initiate handshake")
	// create a new entry in the message buffer for this peer
	// initiate a new nonce service
	cn.messageCache.nonceServiceLock.Lock()
	_, ok := cn.messageCache.nonceService[v.RemotePeer()]
	if ok {
		fmt.Println("already connected; have ref in nonce service")
		cn.messageCache.nonceServiceLock.Unlock()
		return
	} else {
		cn.messageCache.nonceService[v.RemotePeer()] = NewNonceService(
			cn.messageCache.identity,
		)
	}
	cn.messageCache.nonceServiceLock.Unlock()

	fmt.Println("made it to handashke")
	handshakeTimeoutCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cn.InitiateRequestForNonceHandler(
		handshakeTimeoutCtx, v.RemotePeer(),
	)

	cn.messageCache.messageBufferLock.Lock()
	_, ok = cn.messageCache.messageBuffer[v.RemotePeer()]
	if ok {
		fmt.Println("This is a really bad state, shouldn't be here")
		cn.messageCache.messageBufferLock.Unlock()
		return
	} else {
		cn.messageCache.messageBuffer[v.RemotePeer()] = make(
			[]pb.NetworkMessage, 100,
		)
	}
	cn.messageCache.messageBufferLock.Unlock()
}

func (cn *channelNotifiee) Disconnected(n inet.Network, v inet.Conn) {
	// clean up references to peer in nonce cache
	cn.messageCache.nonceServiceLock.Lock()
	_, ok := cn.messageCache.nonceService[v.RemotePeer()]
	cn.messageCache.nonceServiceLock.Unlock()
	if ok {
		delete(cn.messageCache.nonceService, v.RemotePeer())
	}

	cn.messageCache.messageBufferLock.Lock()
	_, ok = cn.messageCache.messageBuffer[v.RemotePeer()]
	cn.messageCache.messageBufferLock.Unlock()
	if ok {
		delete(cn.messageCache.messageBuffer, v.RemotePeer())
	}
}
func (cn *channelNotifiee) OpenedStream(n inet.Network, v inet.Stream) {
}
func (cn *channelNotifiee) ClosedStream(n inet.Network, v inet.Stream) {
}
