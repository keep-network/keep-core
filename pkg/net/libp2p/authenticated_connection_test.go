package libp2p

import (
	"context"
	crand "crypto/rand"
	"io"
	"net"
	"testing"
	"time"

	"github.com/keep-network/keep-core/pkg/net/key"
	libp2pcrypto "github.com/libp2p/go-libp2p-crypto"
	peer "github.com/libp2p/go-libp2p-peer"
)

func TestHandshakeRoundTrip(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Connect the initiator and responder sessions
	authnInboundConn, authnOutboundConn := connectInitiatorAndResponderFull(t, ctx)

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

func connectInitiatorAndResponderFull(t *testing.T, ctx context.Context) (*authenticatedConnection, *authenticatedConnection) {
	initiatorStaticKey, initiatorPeerID := testStaticKeyAndID(t)
	responderStaticKey, responderPeerID := testStaticKeyAndID(t)
	initiatorConn, responderConn := newConnPair()

	var (
		done              = make(chan struct{})
		initiatorErr      error
		authnOutboundConn *authenticatedConnection
	)
	go func(
		ctx context.Context,
		initiatorConn net.Conn,
		initiatorPeerID peer.ID,
		initiatorStaticKey libp2pcrypto.PrivKey,
		responderPeerID peer.ID,
	) {
		authnOutboundConn, initiatorErr = newAuthenticatedOutboundConnection(
			ctx,
			initiatorConn,
			initiatorPeerID,
			initiatorStaticKey,
			responderPeerID,
		)
		done <- struct{}{}
	}(ctx, initiatorConn, initiatorPeerID, initiatorStaticKey, responderPeerID)

	authnInboundConn, err := newAuthenticatedInboundConnection(
		ctx,
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

func newConnPair() (net.Conn, net.Conn) {
	return net.Pipe()
}
