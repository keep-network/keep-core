package handshake

import (
	"context"
	"net"

	libp2pcrypto "github.com/libp2p/go-libp2p-crypto"
	peer "github.com/libp2p/go-libp2p-peer"
)

type secureSession struct {
	net.Conn

	localPrivateKey libp2pcrypto.PrivKey
	localPeer       peer.ID

	remotePeer      peer.ID
	remotePublicKey libp2pcrypto.PubKey
}

func newSecureSession(ctx context.Context, local peer.ID, privateKey libp2pcrypto.PrivKey, insecure net.Conn, remotePeer peer.ID) (*secureSession, error) {
	remotePublicKey, err := remotePeer.ExtractPublicKey()
	if err != nil {
		return nil, err
	}

	session := &secureSession{
		Conn:            insecure,
		localPeer:       local,
		localPrivateKey: privateKey,
		remotePeer:      remotePeer,
		remotePublicKey: remotePublicKey,
	}
	return session, nil
}
