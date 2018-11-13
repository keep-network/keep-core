package libp2p

import (
	"context"
	crand "crypto/rand"
	"io"
	"net"
	"testing"
	"time"

	protoio "github.com/gogo/protobuf/io"
	"github.com/keep-network/keep-core/pkg/net/gen/pb"
	"github.com/keep-network/keep-core/pkg/net/key"
	"github.com/keep-network/keep-core/pkg/net/security/handshake"
	libp2pcrypto "github.com/libp2p/go-libp2p-crypto"
	peer "github.com/libp2p/go-libp2p-peer"
)

func TestPinnedAndMessageKeyMismatch(t *testing.T) {
	initiatorStaticKey, initiatorPeerID := testStaticKeyAndID(t)
	responderStaticKey, responderPeerID := testStaticKeyAndID(t)
	initiatorConn, responderConn := newConnPair()

	go func(
		initiatorConn net.Conn,
		initiatorPeerID peer.ID,
		initiatorStaticKey libp2pcrypto.PrivKey,
		responderPeerID peer.ID,
		responderStaticKey libp2pcrypto.PrivKey,
	) {
		ac := &authenticatedConnection{
			Conn:                initiatorConn,
			localPeerID:         initiatorPeerID,
			localPeerPrivateKey: initiatorStaticKey,
			remotePeerID:        responderPeerID,
			remotePeerPublicKey: responderStaticKey.GetPublic(),
		}

		maliciousInitiatorHijacksHonestRun(t, ac)
		return
	}(initiatorConn, initiatorPeerID, initiatorStaticKey, responderPeerID, responderStaticKey)

	_, err := newAuthenticatedInboundConnection(
		responderConn,
		responderPeerID,
		responderStaticKey,
		"",
	)
	if err == nil {
		t.Fatal("should not have successfully completed handshake")
	}
}

// maliciousInitiatorHijacksHonestRun simulates an honest Acts 1 and 2 as an
// initiator, and then drops in a malicious peer for Act 3. Properly implemented
// peer-pinning should ensure that a malicious peer can't hijack a connection
// after the first act and sign subsequent messages.
func maliciousInitiatorHijacksHonestRun(t *testing.T, ac *authenticatedConnection) {
	initiatorConnectionReader := protoio.NewDelimitedReader(ac.Conn, maxFrameSize)
	initiatorConnectionWriter := protoio.NewDelimitedWriter(ac.Conn)

	initiatorAct1, err := handshake.InitiateHandshake()
	if err != nil {
		t.Fatal(err)
	}

	act1WireMessage, err := initiatorAct1.Message().Marshal()
	if err != nil {
		t.Fatal(err)
	}

	if err := ac.initiatorSendAct1(act1WireMessage, initiatorConnectionWriter); err != nil {
		t.Fatal(err)
	}

	initiatorAct2 := initiatorAct1.Next()

	act2Message, err := ac.initiatorReceiveAct2(initiatorConnectionReader)
	if err != nil {
		t.Fatal(err)
	}

	initiatorAct3, err := initiatorAct2.Next(act2Message)
	if err != nil {
		t.Fatal(err)
	}

	act3WireMessage, err := initiatorAct3.Message().Marshal()
	if err != nil {
		t.Fatal(err)
	}

	maliciousInitiatorStaticKey, maliciousInitiatorPeerID := testStaticKeyAndID(t)
	signedAct3Message, err := maliciousInitiatorStaticKey.Sign(act3WireMessage)
	if err != nil {
		t.Fatal(err)
	}

	act3Envelope := &pb.HandshakeEnvelope{
		Message:   act3WireMessage,
		PeerID:    []byte(maliciousInitiatorPeerID),
		Signature: signedAct3Message,
	}

	if err := initiatorConnectionWriter.WriteMsg(act3Envelope); err != nil {
		t.Fatal(err)
	}
}

func TestHandshakeRoundTrip(t *testing.T) {
	_, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Connect the initiator and responder sessions
	authnInboundConn, authnOutboundConn := connectInitiatorAndResponderFull(t)

	msg := []byte("brown fox blue tail")
	go func(authnOutboundConn *authenticatedConnection, msg []byte) {
		if _, err := authnOutboundConn.Write(msg); err != nil {
			t.Fatal(err)
		}
	}(authnOutboundConn, msg)

	msgContainer := make([]byte, len(msg))
	if _, err := io.ReadFull(authnInboundConn.Conn, msgContainer); err != nil {
		t.Fatal(err)
	}

	if string(msgContainer) != string(msg) {
		t.Fatalf("message mismatch got %v, want %v", string(msgContainer), string(msg))
	}
}

func connectInitiatorAndResponderFull(t *testing.T) (*authenticatedConnection, *authenticatedConnection) {
	initiatorStaticKey, initiatorPeerID := testStaticKeyAndID(t)
	responderStaticKey, responderPeerID := testStaticKeyAndID(t)
	initiatorConn, responderConn := newConnPair()

	var (
		done              = make(chan struct{})
		initiatorErr      error
		authnOutboundConn *authenticatedConnection
	)
	go func(
		initiatorConn net.Conn,
		initiatorPeerID peer.ID,
		initiatorStaticKey libp2pcrypto.PrivKey,
		responderPeerID peer.ID,
	) {
		authnOutboundConn, initiatorErr = newAuthenticatedOutboundConnection(
			initiatorConn,
			initiatorPeerID,
			initiatorStaticKey,
			responderPeerID,
		)
		done <- struct{}{}
	}(initiatorConn, initiatorPeerID, initiatorStaticKey, responderPeerID)

	authnInboundConn, err := newAuthenticatedInboundConnection(
		responderConn,
		responderPeerID,
		responderStaticKey,
		"",
	)
	if err != nil {
		t.Fatalf("failed to connect initiator with responder [%v]", err)
	}

	// handshake is done, and we'll know if the outbound failed
	<-done

	if initiatorErr != nil {
		t.Fatal(initiatorErr)
	}

	return authnInboundConn, authnOutboundConn
}

func testStaticKeyAndID(t *testing.T) (libp2pcrypto.PrivKey, peer.ID) {
	staticKey, err := key.GenerateStaticNetworkKey(crand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	peerID, err := peer.IDFromPrivateKey(staticKey)
	if err != nil {
		t.Fatal(err)
	}
	return staticKey, peerID
}

// Connect an initiator and responder via a full duplex network connection (reads
// on one end should be matched with writes on the other).
func newConnPair() (net.Conn, net.Conn) {
	return net.Pipe()
}
