package libp2p

import (
	"context"
	"net"

	secure "github.com/libp2p/go-conn-security"
	libp2pcrypto "github.com/libp2p/go-libp2p-crypto"
	peer "github.com/libp2p/go-libp2p-peer"
)

// ID is the multistream-select protocol ID that should be used when identifying
// this security transport.
const handshakeID = "/keep/handshake/1.0.0"

// Compile time assertions of the libp2p interfaces we implement
var _ secure.Transport = (*transport)(nil)

// transport constructs an authenticated communication connection for a peer.
type transport struct {
	localPeerID peer.ID
	privateKey  libp2pcrypto.PrivKey
}

func newAuthenticatedTransport(pk libp2pcrypto.PrivKey) (*transport, error) {
	id, err := peer.IDFromPrivateKey(pk)
	if err != nil {
		return nil, err
	}
	return &transport{
		localPeerID: id,
		privateKey:  pk,
	}, nil
}

// SecureInbound secures an inbound connection.
func (t *transport) SecureInbound(
	ctx context.Context,
	unauthenticatedConn net.Conn,
) (secure.Conn, error) {
	return newAuthenticatedInboundConnection(
		unauthenticatedConn,
		t.localPeerID,
		t.privateKey,
		"",
	)
}

// SecureOutbound secures an outbound connection.
func (t *transport) SecureOutbound(
	ctx context.Context,
	unauthenticatedConn net.Conn,
	remotePeerID peer.ID,
) (secure.Conn, error) {
	return newAuthenticatedOutboundConnection(
		unauthenticatedConn,
		t.localPeerID,
		t.privateKey,
		remotePeerID,
	)
}
