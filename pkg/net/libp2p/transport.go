package libp2p

import (
	"context"
	"net"

	"github.com/keep-network/keep-core/pkg/chain"
	secure "github.com/libp2p/go-conn-security"
	libp2pcrypto "github.com/libp2p/go-libp2p-crypto"
	peer "github.com/libp2p/go-libp2p-peer"
)

// ID is the multistream-select protocol ID that should be used when identifying
// this security transport.
const handshakeID = "/keep/handshake/1.0.0"

// Compile time assertions of custom types
var _ secure.Transport = (*transport)(nil)
var _ secure.Conn = (*authenticatedConnection)(nil)

// transport constructs an authenticated communication connection for a peer.
type transport struct {
	localPeerID  peer.ID
	privateKey   libp2pcrypto.PrivKey
	stakeMonitor chain.StakeMonitor
}

func newAuthenticatedTransport(
	pk libp2pcrypto.PrivKey,
	stakeMonitor chain.StakeMonitor,
) (*transport, error) {
	id, err := peer.IDFromPrivateKey(pk)
	if err != nil {
		return nil, err
	}
	return &transport{
		localPeerID:  id,
		privateKey:   pk,
		stakeMonitor: stakeMonitor,
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
		t.stakeMonitor,
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
		t.stakeMonitor,
	)
}

// LocalPeer retrieves the local peer.
func (ac *authenticatedConnection) LocalPeer() peer.ID {
	return ac.localPeerID
}

// LocalPrivateKey retrieves the local peer's privateKey
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
