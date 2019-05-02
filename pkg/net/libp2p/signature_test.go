package libp2p

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/keep-network/keep-core/pkg/chain/local"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/net/gen/pb"
	"github.com/keep-network/keep-core/pkg/net/key"
	host "github.com/libp2p/go-libp2p-host"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

// Integration test simulating malicious adversary tampering the network message
// put into the channel. There are two messages sent:
// - one with a valid signature evaluated with sender's key
// - one with a valid signature evaluated with key other than sender's key
// The first message should be properly delivered, the second message should get
// rejected.
func TestRejectMessageWithUnexpectedSignature(t *testing.T) {
	var (
		ctx     = context.Background()
		payload = "I did know once, only I've sort of forgotten."
	)

	connect := func(t *testing.T, a, b host.Host) {
		pinfo := a.Peerstore().PeerInfo(a.ID())
		err := b.Connect(ctx, pinfo)
		if err != nil {
			t.Fatal(err)
		}
	}

	peer1, err := createTestChannel(ctx, 8080)
	if err != nil {
		t.Fatal(err)
	}

	peer2, err := createTestChannel(ctx, 8081)
	if err != nil {
		t.Fatal(err)
	}

	connect(t, peer1.Host(), peer2.Host())
	time.Sleep(time.Millisecond * 100)

	// Create and publish message with a correct signature...
	if err := peer1.Send(&testMessage{Payload: payload}); err != nil {
		t.Fatal(err)
	}

	ensureNonMaliciousMessage := func(t *testing.T, msg *pubsub.Message) error {
		var envelope pb.NetworkEnvelope
		if err := proto.Unmarshal(msg.Data, &envelope); err != nil {
			t.Fatal(err)
		}

		var protoMessage pb.NetworkMessage
		if err := proto.Unmarshal(envelope.Message, &protoMessage); err != nil {
			t.Fatal(err)
		}

		unmarshaled, err := peer2.getUnmarshalingContainerByType(string(protoMessage.Type))
		if err != nil {
			t.Fatal(err)
		}

		if err := unmarshaled.Unmarshal(protoMessage.GetPayload()); err != nil {
			t.Fatal(err)
		}

		testPayload, ok := unmarshaled.(*testMessage)
		if !ok {
			return fmt.Errorf(
				"expected: payload type string\ngot:   payload type [%v]",
				testPayload,
			)
		}

		if payload != testPayload.Payload {
			return fmt.Errorf(
				"expected: message payload [%s]\ngot:   payload [%s]",
				payload,
				testPayload.Payload,
			)
		}
		return nil
	}
	// Check if the message with correct signature (from peer1) has been
	// properly delivered to peer2, and the message with the incorrect
	// signature has been dropped.
	for {
		select {
		// Ensure all messages are flushed before exiting
		case <-time.After(2 * time.Second):
			return
		case <-ctx.Done():
			t.Fatal(ctx.Err())
		default:
			msg, err := peer2.subscription.Next(ctx)
			if err != nil {
				t.Fatal(err)
			}
			if msg.GetSignature() == nil {
				t.Fatalf("expected signature in message")
			}
			if err := ensureNonMaliciousMessage(t, msg); err != nil {
				t.Fatal(err)
			}
		}
	}
}

// createTestChannel creates and initializes `BroadcastChannel` with all
// underlying libp2p setup steps. Created instance is then casted to
// `lib2p.channel` type so the private interface is available and can be
// tested.
func createTestChannel(
	ctx context.Context,
	port int,
) (*channel, error) {
	networkConfig := Config{Port: port}

	staticKey, _, err := key.GenerateStaticNetworkKey()
	if err != nil {
		return nil, err
	}

	provider, err := Connect(
		ctx,
		networkConfig,
		staticKey,
		local.NewStakeMonitor(big.NewInt(200)),
	)
	if err != nil {
		return nil, err
	}

	broadcastChannel, err := provider.ChannelFor("testchannel")
	if err != nil {
		return nil, err
	}

	if err := broadcastChannel.RegisterUnmarshaler(
		func() net.TaggedUnmarshaler { return &testMessage{} },
	); err != nil {
		return nil, err
	}

	ch, ok := broadcastChannel.(*channel)
	if !ok {
		return nil, fmt.Errorf("unexpected channel type")
	}

	return ch, nil
}
