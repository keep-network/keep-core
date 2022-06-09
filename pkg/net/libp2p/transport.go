package libp2p

import (
	"context"
	"net"

	libp2ptls "github.com/libp2p/go-libp2p-tls"

	keepNet "github.com/keep-network/keep-core/pkg/net"
	libp2pcrypto "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/sec"
)

// ID is the multistream-select protocol ID that should be used when identifying
// this security transport.
const handshakeID = "/keep/handshake/1.0.0"

// Compile time assertions of custom types
var _ sec.SecureTransport = (*transport)(nil)
var _ sec.SecureConn = (*authenticatedConnection)(nil)

// transport constructs an encrypted and authenticated connection for a peer.
type transport struct {
	localPeerID     peer.ID
	privateKey      libp2pcrypto.PrivKey
	protocol        string
	firewall        keepNet.Firewall
	encryptionLayer sec.SecureTransport
}

func newEncryptedAuthenticatedTransport(
	pk libp2pcrypto.PrivKey,
	protocol string,
	firewall keepNet.Firewall,
) (*transport, error) {
	id, err := peer.IDFromPrivateKey(pk)
	if err != nil {
		return nil, err
	}

	encryptionLayer, err := libp2ptls.New(pk)
	if err != nil {
		return nil, err
	}

	return &transport{
		localPeerID:     id,
		privateKey:      pk,
		firewall:        firewall,
		encryptionLayer: encryptionLayer,
		protocol:        protocol,
	}, nil
}

// SecureInbound secures an inbound connection.
func (t *transport) SecureInbound(
	ctx context.Context,
	connection net.Conn,
	peerID peer.ID,
) (sec.SecureConn, error) {
	encryptedConnection, err := t.encryptionLayer.SecureInbound(ctx, connection, peerID)
	if err != nil {
		return nil, err
	}

	return newAuthenticatedInboundConnection(
		encryptedConnection,
		t.localPeerID,
		t.privateKey,
		t.firewall,
		t.protocol,
	)
}

// SecureOutbound secures an outbound connection.
func (t *transport) SecureOutbound(
	ctx context.Context,
	connection net.Conn,
	remotePeerID peer.ID,
) (sec.SecureConn, error) {
	encryptedConnection, err := t.encryptionLayer.SecureOutbound(
		ctx,
		connection,
		remotePeerID,
	)
	if err != nil {
		return nil, err
	}

	return newAuthenticatedOutboundConnection(
		encryptedConnection,
		t.localPeerID,
		t.privateKey,
		remotePeerID,
		t.firewall,
		t.protocol,
	)
}
