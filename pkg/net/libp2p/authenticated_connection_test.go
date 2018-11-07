package libp2p

import (
	"context"
	crand "crypto/rand"
	"fmt"
	"io"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/chain/local"
	"github.com/keep-network/keep-core/pkg/net/key"
	libp2pcrypto "github.com/libp2p/go-libp2p-crypto"
	peer "github.com/libp2p/go-libp2p-peer"
)

func TestHandshake(t *testing.T) {
	_, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	initiator := createTestConnectionConfig(t)
	responder := createTestConnectionConfig(t)

	stakeMonitoring := local.NewStakeMonitoring()
	stakeMonitoring.StakeTokens(key.NetworkPubKeyToEthAddress(initiator.pubKey))
	stakeMonitoring.StakeTokens(key.NetworkPubKeyToEthAddress(responder.pubKey))

	authnInboundConn, authnOutboundConn, inboundError, outboundError :=
		connectInitiatorAndResponder(initiator, responder, stakeMonitoring, t)
	if inboundError != nil {
		t.Fatal(inboundError)
	}
	if outboundError != nil {
		t.Fatal(outboundError)
	}

	// send a test message over the established connection
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

func TestHandshakeNoInitiatorStake(t *testing.T) {
	_, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	initiator := createTestConnectionConfig(t)
	responder := createTestConnectionConfig(t)

	stakeMonitoring := local.NewStakeMonitoring()
	// only responder is staked
	stakeMonitoring.StakeTokens(key.NetworkPubKeyToEthAddress(responder.pubKey))

	_, _, inboundError, outboundError :=
		connectInitiatorAndResponder(initiator, responder, stakeMonitoring, t)

	if inboundError != nil {
		t.Fatal(inboundError)
	}

	expectedOutboundError := fmt.Errorf("connection handshake failed - remote peer has no minimum stake")
	if !reflect.DeepEqual(expectedOutboundError, outboundError) {
		t.Fatalf(
			"unexpected outbound connection error\nexpected: %v\nactual: %v",
			expectedOutboundError,
			outboundError,
		)
	}
}

func TestHanshakeNoResponderStake(t *testing.T) {
	_, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	initiator := createTestConnectionConfig(t)
	responder := createTestConnectionConfig(t)

	stakeMonitoring := local.NewStakeMonitoring()
	// only initiator is staked
	stakeMonitoring.StakeTokens(key.NetworkPubKeyToEthAddress(initiator.pubKey))

	_, _, inboundError, outboundError :=
		connectInitiatorAndResponder(initiator, responder, stakeMonitoring, t)

	expectedInboundError := fmt.Errorf("connection handshake failed - remote peer has no minimum stake")
	if !reflect.DeepEqual(expectedInboundError, inboundError) {
		t.Fatalf(
			"unexpected outbound connection error\nexpected: %v\nactual: %v",
			expectedInboundError,
			inboundError,
		)
	}

	if outboundError != nil {
		t.Fatal(outboundError)
	}
}

func connectInitiatorAndResponder(
	initiator *testConnectionConfig,
	responder *testConnectionConfig,
	stakeMonitoring chain.StakeMonitoring,
	t *testing.T,
) (
	authnInboundConn *authenticatedConnection,
	authnOutboundConn *authenticatedConnection,
	outboundError error,
	inboundError error,
) {

	initiatorConn, responderConn := newConnPair()

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
	}(initiatorConn, initiator.peerID, initiator.privKey, responder.peerID)

	authnInboundConn, inboundError = newAuthenticatedInboundConnection(
		responderConn,
		responder.peerID,
		responder.privKey,
		stakeMonitoring,
	)

	<-done // handshake is done

	return
}

type testConnectionConfig struct {
	privKey *key.NetworkPrivateKey
	pubKey  *key.NetworkPublicKey
	peerID  peer.ID
}

func createTestConnectionConfig(t *testing.T) *testConnectionConfig {
	privKey, pubKey, err := key.GenerateStaticNetworkKey(crand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	peerID, err := peer.IDFromPrivateKey(privKey)
	if err != nil {
		t.Fatal(err)
	}

	return &testConnectionConfig{privKey, pubKey, peerID}
}

// Connect an initiator and responder via a full duplex network connection (reads
// on one end should be matched with writes on the other).
func newConnPair() (net.Conn, net.Conn) {
	return net.Pipe()
}
