package libp2p

import (
	"context"
	crand "crypto/rand"
	"io"
	"net"
	"testing"
	"time"

	"github.com/keep-network/keep-core/pkg/chain/local"
	"github.com/keep-network/keep-core/pkg/net/key"
	libp2pcrypto "github.com/libp2p/go-libp2p-crypto"
	peer "github.com/libp2p/go-libp2p-peer"
)

func TestHandshakeRoundTrip(t *testing.T) {
	_, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Connect the initiator and responder sessions
	authnInboundConn, authnOutboundConn, inboundError, outboundError :=
		connectInitiatorAndResponderFull(t)
	if inboundError != nil {
		t.Fatal(inboundError)
	}
	if outboundError != nil {
		t.Fatal(outboundError)
	}

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

func connectInitiatorAndResponderFull(t *testing.T) (
	authnInboundConn *authenticatedConnection,
	authnOutboundConn *authenticatedConnection,
	outboundError error,
	inboundError error,
) {
	initiatorPrivKey, initiatorPubKey, initiatorPeerID := testStaticKeyAndID(t)
	responderPrivKey, responderPubKey, responderPeerID := testStaticKeyAndID(t)
	initiatorConn, responderConn := newConnPair()

	stakeMonitoring := local.NewStakeMonitoring()
	stakeMonitoring.StakeTokens(key.NetworkPubKeyToEthAddress(initiatorPubKey))
	stakeMonitoring.StakeTokens(key.NetworkPubKeyToEthAddress(responderPubKey))

	done := make(chan struct{})

	go func(
		initiatorConn net.Conn,
		initiatorPeerID peer.ID,
		initiatorPrivKey libp2pcrypto.PrivKey,
		responderPeerID peer.ID,
	) {
		authnOutboundConn, outboundError = newAuthenticatedOutboundConnection(
			initiatorConn,
			initiatorPeerID,
			initiatorPrivKey,
			responderPeerID,
			stakeMonitoring,
		)
		done <- struct{}{}
	}(initiatorConn, initiatorPeerID, initiatorPrivKey, responderPeerID)

	authnInboundConn, inboundError = newAuthenticatedInboundConnection(
		responderConn,
		responderPeerID,
		responderPrivKey,
		stakeMonitoring,
	)

	<-done // handshake is done

	return
}

func testStaticKeyAndID(t *testing.T) (
	*key.NetworkPrivateKey,
	*key.NetworkPublicKey,
	peer.ID,
) {
	privKey, pubKey, err := key.GenerateStaticNetworkKey(crand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	peerID, err := peer.IDFromPrivateKey(privKey)
	if err != nil {
		t.Fatal(err)
	}
	return privKey, pubKey, peerID
}

// Connect an initiator and responder via a full duplex network connection (reads
// on one end should be matched with writes on the other).
func newConnPair() (net.Conn, net.Conn) {
	return net.Pipe()
}
