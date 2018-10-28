package libp2p

import (
	"context"
	"net"

	"github.com/keep-network/keep-core/pkg/net/gen/pb"
	"github.com/keep-network/keep-core/pkg/net/security/handshake"
	libp2pcrypto "github.com/libp2p/go-libp2p-crypto"
	peer "github.com/libp2p/go-libp2p-peer"

	protoio "github.com/gogo/protobuf/io"
)

const maxFrameSize = 1 << 20

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
	ctx context.Context,
	unauthenticatedConn net.Conn,
	localPeerID peer.ID,
	privateKey libp2pcrypto.PrivKey,
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

	if err := ac.runHandshake(ctx); err != nil {
		// close the conn before returning otherwise we leak
		ac.Close()
		return nil, err
	}

	return ac, nil
}

func (ac *authenticatedConnection) runHandshakeAsInitiator(ctx context.Context) error {
	// initiator station

	//
	// Act 1
	//

	initiatorConnectionWriter := protoio.NewDelimitedWriter(ac.Conn)

	initiatorAct1, err := handshake.InitiateHandshake()
	if err != nil {
		return err
	}

	act1WireMessage := initiatorAct1.Message().Proto()
	if err := initiatorConnectionWriter.WriteMsg(act1WireMessage); err != nil {
		return err
	}

	initiatorAct2 := initiatorAct1.Next()

	//
	// Act 2
	//

	initiatorConnectionReader := protoio.NewDelimitedReader(ac.Conn, maxFrameSize)

	var act2WireResponseMessage pb.Act2Message
	if err := initiatorConnectionReader.ReadMsg(&act2WireResponseMessage); err != nil {
		return err
	}

	act2Message := handshake.Act2MessageFromProto(act2WireResponseMessage)
	initiatorAct3, err := initiatorAct2.Next(act2Message)
	if err != nil {
		return err
	}

	//
	// Act 3
	//

	act3WireMessage := initiatorAct3.Message().Proto()
	if err := initiatorConnectionWriter.WriteMsg(act3WireMessage); err != nil {
		return err
	}

	return nil
}

	// responder station
	err = responderAct3.FinalizeHandshake(act3Message)
	if err != nil {
		return err
	}
	return nil
}
