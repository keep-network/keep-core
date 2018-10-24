// Package handshake is used to turn inbound and outbound unauthenticated,
// plain-text connections into authenticated, plain-text connections. Noticeably,
// it does not guarantee confidentiality as it does not encrypt connections.
package handshake

import (
	"context"
	"net"

	secure "github.com/libp2p/go-conn-security"
	libp2pcrypto "github.com/libp2p/go-libp2p-crypto"
	peer "github.com/libp2p/go-libp2p-peer"
)

// ID is the multistream-select protocol ID that should be used when identifying
// this security transport. Unfortunately, listeners are configured to match on
// either /plaintext/1.0.0 or /secio/1.0.0. For now, we must lie until custom listeners.
const ID = "/secio/1.0.0"

// Compile time assertions of custom types
var _ secure.Transport = (*Transport)(nil)
var _ secure.Conn = (*secureSession)(nil)

// Transport constructs secure communication sessions for a peer.
type Transport struct {
	LocalID    peer.ID
	PrivateKey libp2pcrypto.PrivKey
}

func New(pk libp2pcrypto.PrivKey) (*Transport, error) {
	id, err := peer.IDFromPrivateKey(pk)
	if err != nil {
		return nil, err
	}
	return &Transport{
		LocalID:    id,
		PrivateKey: pk,
	}, nil
}

// SecureInbound secures an inbound connection.
func (t *Transport) SecureInbound(ctx context.Context, insecure net.Conn) (secure.Conn, error) {
	return newSecureSession(ctx, t.LocalID, t.PrivateKey, insecure, "")
}

// SecureOutbound secures an outbound connection.
func (t *Transport) SecureOutbound(ctx context.Context, insecure net.Conn, p peer.ID) (secure.Conn, error) {
	return newSecureSession(ctx, t.LocalID, t.PrivateKey, insecure, p)
}

// LocalPeer retrieves the local peer.
func (ss *secureSession) LocalPeer() peer.ID {
	return ss.localPeer
}

// LocalPrivateKey retrieves the local peer's PrivateKey
func (ss *secureSession) LocalPrivateKey() libp2pcrypto.PrivKey {
	return ss.localPrivateKey
}

// RemotePeer retrieves the remote peer.
func (ss *secureSession) RemotePeer() peer.ID {
	return ss.remotePeer
}

// RemotePublicKey retrieves the remote public key.
func (ss *secureSession) RemotePublicKey() libp2pcrypto.PubKey {
	return ss.remotePublicKey
}
