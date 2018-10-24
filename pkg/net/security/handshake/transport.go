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
var _ secure.Conn = (*authenticatedSession)(nil)

// Transport constructs an authenticated communication connection for a peer.
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
func (t *Transport) SecureInbound(
	ctx context.Context,
	unauthenticatedConn net.Conn,
) (secure.Conn, error) {
	return newAuthenticatedSession(
		ctx,
		t.LocalID,
		t.PrivateKey,
		unauthenticatedConn,
		"",
	)
}

// SecureOutbound secures an outbound connection.
func (t *Transport) SecureOutbound(
	ctx context.Context,
	unauthenticatedConn net.Conn,
	remotePeer peer.ID,
) (secure.Conn, error) {
	return newAuthenticatedSession(
		ctx,
		t.LocalID,
		t.PrivateKey,
		unauthenticatedConn,
		remotePeer,
	)
}

// LocalPeer retrieves the local peer.
func (ss *authenticatedSession) LocalPeer() peer.ID {
	return ss.localPeerID
}

// LocalPrivateKey retrieves the local peer's PrivateKey
func (ss *authenticatedSession) LocalPrivateKey() libp2pcrypto.PrivKey {
	return ss.localPeerPrivateKey
}

// RemotePeer retrieves the remote peer.
func (ss *authenticatedSession) RemotePeer() peer.ID {
	return ss.remotePeerID
}

// RemotePublicKey retrieves the remote public key.
func (ss *authenticatedSession) RemotePublicKey() libp2pcrypto.PubKey {
	return ss.remotePeerPublicKey
}
