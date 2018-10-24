package handshake

import (
	"context"
	"net"

	libp2pcrypto "github.com/libp2p/go-libp2p-crypto"
	peer "github.com/libp2p/go-libp2p-peer"
)

type authenticatedSession struct {
	net.Conn

	localPrivateKey libp2pcrypto.PrivKey
	localPeer       peer.ID

	remotePeer      peer.ID
	remotePublicKey libp2pcrypto.PubKey
}

func newAuthenticatedSession(
	ctx context.Context,
	local peer.ID,
	privateKey libp2pcrypto.PrivKey,
	unauthenticatedConn net.Conn,
	remotePeer peer.ID,
) (*authenticatedSession, error) {
	remotePublicKey, err := remotePeer.ExtractPublicKey()
	if err != nil {
		return nil, err
	}

	return &authenticatedSession{
		Conn:            unauthenticatedConn,
		localPeer:       local,
		localPrivateKey: privateKey,
		remotePeer:      remotePeer,
		remotePublicKey: remotePublicKey,
	}, nil
}
