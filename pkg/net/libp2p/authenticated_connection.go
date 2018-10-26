package libp2p

import (
	"net"

	libp2pcrypto "github.com/libp2p/go-libp2p-crypto"
	peer "github.com/libp2p/go-libp2p-peer"
)

// authenticatedConnection turns inbound and outbound unauthenticated,
// plain-text connections into authenticated, plain-text connections. Noticeably,
// it does not guarantee confidentiality as it does not encrypt connections.
type authenticatedConnection struct {
	net.Conn

	localPeerID         peer.ID
	localPeerPrivateKey libp2pcrypto.PrivKey

	remotePeerID        peer.ID
	remotePeerPublicKey libp2pcrypto.PubKey
}

func newAuthenticatedConnection(
	localPeerID peer.ID,
	privateKey libp2pcrypto.PrivKey,
	unauthenticatedConn net.Conn,
	remotePeerID peer.ID,
) (*authenticatedConnection, error) {
	var (
		remotePublicKey libp2pcrypto.PubKey
		err             error
	)

	if remotePeerID == "" {
		// SecureInbound case; if we don't have a remote peer.id, we
		// can't have their public key!
		remotePublicKey = nil
	} else {
		remotePublicKey, err = remotePeerID.ExtractPublicKey()
		if err != nil {
			return nil, err
		}
	}

	return &authenticatedConnection{
		Conn:                unauthenticatedConn,
		localPeerID:         localPeerID,
		localPeerPrivateKey: privateKey,
		remotePeerID:        remotePeerID,
		remotePeerPublicKey: remotePublicKey,
	}, nil
}
