package libp2p

import (
	"fmt"
	"net"

	"github.com/keep-network/keep-core/pkg/net/gen/pb"
	"github.com/keep-network/keep-core/pkg/net/security/handshake"
	secure "github.com/libp2p/go-conn-security"
	libp2pcrypto "github.com/libp2p/go-libp2p-crypto"
	peer "github.com/libp2p/go-libp2p-peer"

	protoio "github.com/gogo/protobuf/io"
)

// Enough space for a proto-encoded envelope with a message, peer.ID, and sig.
const maxFrameSize = 1024

// Compile time assertions of the libp2p interfaces we implement
var _ secure.Conn = (*authenticatedConnection)(nil)

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

// newAuthenticatedInboundConnection is the connection that's formed by
// transport.SecureInbound. This function is executed by the receiver of a new
// connection, who acts as the responder. This side of the connection lacks
// knowledge of the remotePeerID (passed in as empty string). On success running
// the responder side of the handshake, it returns a fully-authenticated
// connection, which grants access to the network.
func newAuthenticatedInboundConnection(
	unauthenticatedConn net.Conn,
	localPeerID peer.ID,
	privateKey libp2pcrypto.PrivKey,
	remotePeerID peer.ID,
) (*authenticatedConnection, error) {
	ac := &authenticatedConnection{
		Conn:                unauthenticatedConn,
		localPeerID:         localPeerID,
		localPeerPrivateKey: privateKey,
	}

	if err := ac.runHandshakeAsResponder(); err != nil {
		// close the conn before returning (if it hasn't already)
		// otherwise we leak.
		ac.Close()
		return nil, fmt.Errorf("connection handshake failed [%v]", err)
	}

	return ac, nil
}

// newAuthenticatedOutboundConnection is the connection that's formed by
// transport.SecureOutbound. This function is executed by the initiator of a new
// connection. This side of the connection has knowledge of the remote peer
// identity. On success running the initiator side of the handshake, it returns a
// fully-authenticated connection, which grants access to the network.
func newAuthenticatedOutboundConnection(
	unauthenticatedConn net.Conn,
	localPeerID peer.ID,
	privateKey libp2pcrypto.PrivKey,
	remotePeerID peer.ID,
) (*authenticatedConnection, error) {
	remotePublicKey, err := remotePeerID.ExtractPublicKey()
	if err != nil {
		return nil, fmt.Errorf(
			"could not create new authenticated outbound connection [%v]",
			err,
		)
	}

	ac := &authenticatedConnection{
		Conn:                unauthenticatedConn,
		localPeerID:         localPeerID,
		localPeerPrivateKey: privateKey,
		remotePeerID:        remotePeerID,
		remotePeerPublicKey: remotePublicKey,
	}

	if err := ac.runHandshakeAsInitiator(); err != nil {
		ac.Close()
		return nil, fmt.Errorf("connection handshake failed [%v]", err)
	}

	return ac, nil
}

func (ac *authenticatedConnection) runHandshakeAsInitiator() error {
	// initiator station

	initiatorConnectionReader := protoio.NewDelimitedReader(ac.Conn, maxFrameSize)
	initiatorConnectionWriter := protoio.NewDelimitedWriter(ac.Conn)

	//
	// Act 1
	//

	initiatorAct1, err := handshake.InitiateHandshake()
	if err != nil {
		return err
	}

	act1WireMessage, err := initiatorAct1.Message().Marshal()
	if err != nil {
		return err
	}

	if err := initiatorSendAct1(
		act1WireMessage,
		initiatorConnectionWriter,
		ac.localPeerPrivateKey,
		ac.localPeerID,
	); err != nil {
		return err
	}

	initiatorAct2 := initiatorAct1.Next()

	//
	// Act 2
	//

	act2Message, err := initiatorReceiveAct2(initiatorConnectionReader)
	if err != nil {
		return err
	}

	initiatorAct3, err := initiatorAct2.Next(act2Message)
	if err != nil {
		return err
	}

	//
	// Act 3
	//

	act3WireMessage, err := initiatorAct3.Message().Marshal()
	if err != nil {
		return err
	}

	if err := initiatorSendAct3(
		act3WireMessage,
		initiatorConnectionWriter,
		ac.localPeerPrivateKey,
		ac.localPeerID,
	); err != nil {
		return err
	}

	return nil
}

// initiatorSendAct1 signs a marshaled *handshake.Act1Message, prepares
// the message in a pb.HandshakeEnvelope, and sends the message to the responder
// (over the open connection) from the initiator.
func initiatorSendAct1(
	act1WireMessage []byte,
	initiatorConnectionWriter protoio.WriteCloser,
	localPeerPrivateKey libp2pcrypto.PrivKey,
	localPeerID peer.ID,
) error {
	signedAct1Message, err := localPeerPrivateKey.Sign(act1WireMessage)
	if err != nil {
		return err
	}

	act1Envelope := &pb.HandshakeEnvelope{
		Message:   act1WireMessage,
		PeerID:    []byte(localPeerID),
		Signature: signedAct1Message,
	}

	if err := initiatorConnectionWriter.WriteMsg(act1Envelope); err != nil {
		return err
	}

	return nil
}

// initiatorReceiveAct2 unmarshals a pb.HandshakeEnvelope from a responder,
// verifies that the signed messages matches the expected peer.ID, and returns
// the handshake.Act2Message for processing by the initiator.
func initiatorReceiveAct2(
	initiatorConnectionReader protoio.ReadCloser,
) (*handshake.Act2Message, error) {
	var (
		act2Envelope pb.HandshakeEnvelope
		act2Message  = &handshake.Act2Message{}
	)
	if err := initiatorConnectionReader.ReadMsg(&act2Envelope); err != nil {
		return nil, err
	}

	if err := ac.verify(
		ac.remotePeerID,
		peer.ID(act2Envelope.GetPeerID()),
		act2Envelope.GetMessage(),
		act2Envelope.GetSignature(),
	); err != nil {
		return nil, err
	}

	if err := act2Message.Unmarshal(act2Envelope.Message); err != nil {
		return nil, err
	}

	return act2Message, nil
}

// initiatorSendAct3 signs a marshaled *handshake.Act3Message, prepares the
// message in a pb.HandshakeEnvelope, and sends the message to the responder
// (over the open connection) from the initiator.
func initiatorSendAct3(
	act3WireMessage []byte,
	initiatorConnectionWriter protoio.WriteCloser,
	localPeerPrivateKey libp2pcrypto.PrivKey,
	localPeerID peer.ID,
) error {
	signedAct3Message, err := localPeerPrivateKey.Sign(act3WireMessage)
	if err != nil {
		return err
	}

	act3Envelope := &pb.HandshakeEnvelope{
		Message:   act3WireMessage,
		PeerID:    []byte(localPeerID),
		Signature: signedAct3Message,
	}

	if err := initiatorConnectionWriter.WriteMsg(act3Envelope); err != nil {
		return err
	}

	return nil
}

func (ac *authenticatedConnection) runHandshakeAsResponder() error {
	// responder station

	responderConnectionReader := protoio.NewDelimitedReader(ac.Conn, maxFrameSize)
	responderConnectionWriter := protoio.NewDelimitedWriter(ac.Conn)

	//
	// Act 1
	//

	act1Message, err := responderReceiveAct1(responderConnectionReader)
	if err != nil {
		return err
	}

	responderAct2, err := handshake.AnswerHandshake(act1Message)
	if err != nil {
		return err
	}

	//
	// Act 2
	//

	act2WireMessage, err := responderAct2.Message().Marshal()
	if err != nil {
		return err
	}
	if err := responderSendAct2(
		act2WireMessage,
		responderConnectionWriter,
		ac.localPeerPrivateKey,
		ac.localPeerID,
	); err != nil {
		return err
	}

	responderAct3 := responderAct2.Next()

	//
	// Act 3
	//

	act3Message, err := responderReceiveAct3(responderConnectionReader)
	if err != nil {
		return err
	}

	if err := responderAct3.FinalizeHandshake(act3Message); err != nil {
		return err
	}

	return nil
}

// responderReceiveAct1 unmarshals a pb.HandshakeEnvelope from an initiator,
// verifies that the signed messages matches the expected peer.ID, and returns
// the handshake.Act1Message for processing by the responder.
func responderReceiveAct1(
	responderConnectionReader protoio.ReadCloser,
) (*handshake.Act1Message, error) {
	var (
		act1Envelope pb.HandshakeEnvelope
		act1Message  = &handshake.Act1Message{}
	)
	if err := responderConnectionReader.ReadMsg(&act1Envelope); err != nil {
		return nil, err
	}

	// Libp2p specific step: the responder has no knowledge of the initiator
	// until after the handshake has succeeded, the connection has been
	// upgraded, and identity information is exchanged. This provides an
	// element of identity hiding for the initiator. To help prevent
	// malicious interference, we want to pin this identity for the duration
	// of the connection.
	ac.remotePeerID = peer.ID(act1Envelope.GetPeerID())

	if err := ac.verify(
		ac.remotePeerID,
		peer.ID(act1Envelope.GetPeerID()),
		act1Envelope.GetMessage(),
		act1Envelope.GetSignature(),
	); err != nil {
		return nil, err
	}

	if err := act1Message.Unmarshal(act1Envelope.Message); err != nil {
		return nil, err
	}

	return act1Message, nil
}

// responderSendAct2 signs a marshaled *handshake.Act2Message, prepares the
// message in a pb.HandshakeEnvelope, and sends the message to the initiator
// (over the open connection) from the responder.
func responderSendAct2(
	act2WireMessage []byte,
	responderConnectionWriter protoio.WriteCloser,
	localPeerPrivateKey libp2pcrypto.PrivKey,
	localPeerID peer.ID,
) error {
	signedAct2Message, err := localPeerPrivateKey.Sign(act2WireMessage)
	if err != nil {
		return err
	}

	act2Envelope := &pb.HandshakeEnvelope{
		Message:   act2WireMessage,
		PeerID:    []byte(localPeerID),
		Signature: signedAct2Message,
	}

	if err := responderConnectionWriter.WriteMsg(act2Envelope); err != nil {
		return err
	}

	return nil
}

// responderReceiveAct3 unmarshals a pb.HandshakeEnvelope from an initiator,
// verifies that the signed messages matches the expected peer.ID, and returns
// the handshake.Act3Message for processing by the responder.
func responderReceiveAct3(
	responderConnectionReader protoio.ReadCloser,
) (*handshake.Act3Message, error) {
	var (
		act3Envelope pb.HandshakeEnvelope
		act3Message  = &handshake.Act3Message{}
	)
	if err := responderConnectionReader.ReadMsg(&act3Envelope); err != nil {
		return nil, err
	}

	if err := ac.verify(
		ac.remotePeerID,
		peer.ID(act3Envelope.GetPeerID()),
		act3Envelope.GetMessage(),
		act3Envelope.GetSignature(),
	); err != nil {
		return nil, err
	}

	if err := act3Message.Unmarshal(act3Envelope.Message); err != nil {
		return nil, err
	}

	return act3Message, nil
}

// verify checks to see if the pinned (expected) identity matches the message
// sender's identity before running through the signature verification check.
func (ac *authenticatedConnection) verify(
	expectedSender, actualSender peer.ID,
	messageBytes, signatureBytes []byte,
) error {
	if expectedSender != actualSender {
		return fmt.Errorf(
			"pinned identity [%v] does not match sender identity [%v]",
			expectedSender,
			actualSender,
		)
	}
	return verifyEnvelope(actualSender, messageBytes, signatureBytes)
}
