package libp2p

import (
	"fmt"

	"github.com/keep-network/keep-core/pkg/net/gen/pb"
	"github.com/keep-network/keep-core/pkg/net/security/handshake"

	libp2pcrypto "github.com/libp2p/go-libp2p-crypto"
	peer "github.com/libp2p/go-libp2p-peer"

	protoio "github.com/gogo/protobuf/io"
)

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
	remotePeerID peer.ID,
) (*handshake.Act2Message, error) {
	var (
		act2Envelope pb.HandshakeEnvelope
		act2Message  = &handshake.Act2Message{}
	)
	if err := initiatorConnectionReader.ReadMsg(&act2Envelope); err != nil {
		return nil, err
	}

	if err := verifyHandshakeMessage(
		remotePeerID,
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

// responderReceiveAct1 unmarshals a pb.HandshakeEnvelope from an initiator,
// verifies that the signed messages matches the expected peer.ID, and returns
// the handshake.Act1Message for processing by the responder.
func responderReceiveAct1(
	responderConnectionReader protoio.ReadCloser,
) (*handshake.Act1Message, peer.ID, error) {
	var (
		act1Envelope pb.HandshakeEnvelope
		act1Message  = &handshake.Act1Message{}
		remotePeerID peer.ID
	)
	if err := responderConnectionReader.ReadMsg(&act1Envelope); err != nil {
		return nil, remotePeerID, err
	}

	// New remote pinned identity
	remotePeerID = peer.ID(act1Envelope.GetPeerID())

	if err := verifyHandshakeMessage(
		remotePeerID,
		peer.ID(act1Envelope.GetPeerID()),
		act1Envelope.GetMessage(),
		act1Envelope.GetSignature(),
	); err != nil {
		return nil, remotePeerID, err
	}

	if err := act1Message.Unmarshal(act1Envelope.Message); err != nil {
		return nil, remotePeerID, err
	}

	return act1Message, remotePeerID, nil
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
	remotePeerID peer.ID,
) (*handshake.Act3Message, error) {
	var (
		act3Envelope pb.HandshakeEnvelope
		act3Message  = &handshake.Act3Message{}
	)
	if err := responderConnectionReader.ReadMsg(&act3Envelope); err != nil {
		return nil, err
	}

	if err := verifyHandshakeMessage(
		remotePeerID,
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

// verifyHandshakeMessage checks to see if the pinned (expected) identity
// matches the message sender's identity before running through the signature
// verification check.
func verifyHandshakeMessage(
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
	return verify(actualSender, messageBytes, signatureBytes)
}
