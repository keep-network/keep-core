package libp2p

import (
	"net"

	"github.com/keep-network/keep-core/pkg/net/security/handshake"
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

	ac := &authenticatedConnection{
		Conn:                unauthenticatedConn,
		localPeerID:         localPeerID,
		localPeerPrivateKey: privateKey,
		remotePeerID:        remotePeerID,
		remotePeerPublicKey: remotePublicKey,
	}

	return ac.run()
}

func (ac *authenticatedConnection) run() (*authenticatedConnection, error) {
	// TODO: placeholder code
	//
	// Act 1
	//

	// initiator station
	initiatorAct1, err := handshake.InitiateHandshake()
	if err != nil {
		return nil, err
	}
	act1Message := initiatorAct1.Message()
	initiatorAct2 := initiatorAct1.Next()

	// responder station
	responderAct2, err := handshake.AnswerHandshake(act1Message)
	if err != nil {
		return nil, err
	}

	//
	// Act 2
	//

	// responder station
	act2Message := responderAct2.Message()
	responderAct3 := responderAct2.Next()

	// initiator station
	initiatorAct3, err := initiatorAct2.Next(act2Message)
	if err != nil {
		return nil, err
	}

	//
	// Act 3
	//

	// initiator station
	act3Message := initiatorAct3.Message()

	// responder station
	err = responderAct3.FinalizeHandshake(act3Message)
	if err != nil {
		return nil, err
	}
	return ac, nil
}
