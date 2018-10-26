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
// this security transport.
const ID = "/keep/handshake/1.0.0"

// Compile time assertions of custom types
var _ secure.Transport = (*Transport)(nil)
var _ secure.Conn = (*authenticatedConnection)(nil)

// Transport constructs an authenticated communication connection for a peer.
type Transport struct {
	LocalPeerID peer.ID
	PrivateKey  libp2pcrypto.PrivKey
}

func New(pk libp2pcrypto.PrivKey) (*Transport, error) {
	id, err := peer.IDFromPrivateKey(pk)
	if err != nil {
		return nil, err
	}
	return &Transport{
		LocalPeerID: id,
		PrivateKey:  pk,
	}, nil
}

// SecureInbound secures an inbound connection.
func (t *Transport) SecureInbound(
	ctx context.Context,
	unauthenticatedConn net.Conn,
) (secure.Conn, error) {
	return newAuthenticatedSession(
		t.LocalPeerID,
		t.PrivateKey,
		unauthenticatedConn,
		"",
	)
}

// SecureOutbound secures an outbound connection.
func (t *Transport) SecureOutbound(
	ctx context.Context,
	unauthenticatedConn net.Conn,
	remotePeerID peer.ID,
) (secure.Conn, error) {
	return newAuthenticatedSession(
		t.LocalPeerID,
		t.PrivateKey,
		unauthenticatedConn,
		remotePeerID,
	)
}

// LocalPeer retrieves the local peer.
func (ac *authenticatedConnection) LocalPeer() peer.ID {
	return ac.localPeerID
}

// LocalPrivateKey retrieves the local peer's PrivateKey
func (ac *authenticatedConnection) LocalPrivateKey() libp2pcrypto.PrivKey {
	return ac.localPeerPrivateKey
}

// RemotePeer returns the remote peer ID if we initiated the dial. Otherwise, it
// returns "" (because this connection isn't actually secure).
func (ac *authenticatedConnection) RemotePeer() peer.ID {
	return ac.remotePeerID
}

// RemotePublicKey retrieves the remote public key.
func (ac *authenticatedConnection) RemotePublicKey() libp2pcrypto.PubKey {
	return ac.remotePeerPublicKey
}
