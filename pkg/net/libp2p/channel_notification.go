package libp2p

import (
	"context"
	"fmt"
	"time"

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
}

func (cn *channelNotifiee) Disconnected(n inet.Network, v inet.Conn) {
}
func (cn *channelNotifiee) OpenedStream(n inet.Network, v inet.Stream) {
	// initiate a new nonce service
	cn.messageCache.nonceServiceLock.Lock()
	_, ok := cn.messageCache.nonceService[v.Conn().RemotePeer()]
	if ok {
		fmt.Println("already connected; have ref in nonce service")
		cn.messageCache.nonceServiceLock.Unlock()
		return
	} else {
		cn.messageCache.nonceService[v.Conn().RemotePeer()] = NewNonceService(
			cn.messageCache.identity,
		)
	}
	cn.messageCache.nonceServiceLock.Unlock()

	fmt.Println("made it to shake")
	handshakeTimeoutCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cn.InitiateRequestForNonceHandler(
		handshakeTimeoutCtx, v.Conn().RemotePeer(),
	)
}
func (cn *channelNotifiee) ClosedStream(n inet.Network, v inet.Stream) {
	// clean up references to peer in nonce cache
	cn.messageCache.nonceServiceLock.Lock()
	_, ok := cn.messageCache.nonceService[v.Conn().RemotePeer()]
	cn.messageCache.nonceServiceLock.Unlock()
	if ok {
		delete(cn.messageCache.nonceService, v.Conn().RemotePeer())
	}
}
