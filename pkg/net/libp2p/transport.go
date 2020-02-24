package libp2p

import (
	"context"
	"net"

	secio "github.com/libp2p/go-libp2p-secio"

	"github.com/keep-network/keep-core/pkg/chain"
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
	stakeMonitor    chain.StakeMonitor
	encryptionLayer sec.SecureTransport
}

func newEncryptedAuthenticatedTransport(
	pk libp2pcrypto.PrivKey,
	stakeMonitor chain.StakeMonitor,
) (*transport, error) {
	id, err := peer.IDFromPrivateKey(pk)
	if err != nil {
		return nil, err
	}

	encryptionLayer, err := secio.New(pk)
	if err != nil {
		return nil, err
	}

	return &transport{
		localPeerID:     id,
		privateKey:      pk,
		stakeMonitor:    stakeMonitor,
		encryptionLayer: encryptionLayer,
	}, nil
}

// SecureInbound secures an inbound connection.
func (t *transport) SecureInbound(
	ctx context.Context,
	connection net.Conn,
) (sec.SecureConn, error) {
	encryptedConnection, err := t.encryptionLayer.SecureInbound(ctx, connection)
	if err != nil {
		return nil, err
	}

	return newAuthenticatedInboundConnection(
		encryptedConnection,
		t.localPeerID,
		t.privateKey,
		t.stakeMonitor,
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
		t.stakeMonitor,
	)
}
