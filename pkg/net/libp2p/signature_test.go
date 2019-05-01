package libp2p

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/keep-network/keep-core/pkg/chain/local"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/net/key"
	host "github.com/libp2p/go-libp2p-host"
)

// Integration test simulating malicious adversary tampering the network message
// put into the channel. There are two messages sent:
// - one with a valid signature evaluated with sender's key
// - one with a valid signature evaluated with key other than sender's key
// The first message should be properly delivered, the second message should get
// rejected.
func TestRejectMessageWithUnexpectedSignature(t *testing.T) {
	var (
		ctx           = context.Background()
		honestPayload = "I did know once, only I've sort of forgotten."
		// maliciousPayload = "You never can tell with bees."
	)

	connect := func(t *testing.T, a, b host.Host) {
		pinfo := a.Peerstore().PeerInfo(a.ID())
		err := b.Connect(context.Background(), pinfo)
		if err != nil {
			t.Fatal(err)
		}
	}

	honestPeer1, err := createTestChannel(ctx, 8080)
	if err != nil {
		t.Fatal(err)
	}

	honestPeer2, err := createTestChannel(ctx, 3030)
	if err != nil {
		t.Fatal(err)
	}

	adversarialPeer, err := createTestChannel(ctx, 2020)
	if err != nil {
		t.Fatal(err)
	}

	// Create the following network topology:
	// honestPeer1 <-> adversarialPeer <-> honestPeer2
	connect(t, honestPeer1.Host(), adversarialPeer.Host())
	connect(t, honestPeer2.Host(), adversarialPeer.Host())
	time.Sleep(time.Millisecond * 100)

	// Create and publish message with a correct signature...
	if err := honestPeer1.Send(&testMessage{Payload: honestPayload}); err != nil {
		t.Fatal(err)
	}

	// Create a function which in a loop calls next, grabs the honest message,
	// attempts to perturb it and then send it back on the broadcast channel.
	// No one should receive that message. And they should check they don't get it

	// for {
	// 	msg, err := adversarialPeer.subscription.Next(context.Background())
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		continue
	// 	}
	// 	var envelope pb.NetworkEnvelope
	// 	if err := proto.Unmarshal(msg.Data, &envelope); err != nil {
	// 		fmt.Println(err)
	// 		continue
	// 	}
	// 	// somehow tamper with the message and broadcast this?
	// }

	// Check if the message with correct signature (from honestPeer1) has been
	// properly delivered to honestPeer2, and the message with the incorrect
	// signature has been dropped.
	honestPeer2RecvChan := make(chan net.Message)
	if err := honestPeer2.Recv(net.HandleMessageFunc{
		Type: "test",
		Handler: func(msg net.Message) error {
			honestPeer2RecvChan <- msg
			return nil
		},
	}); err != nil {
		t.Fatal(err)
	}

	// Ensure that honestPeer1 drops the malicious message
	honestPeer1RecvChan := make(chan net.Message)
	if err := honestPeer1.Recv(net.HandleMessageFunc{
		Type: "test",
		Handler: func(msg net.Message) error {
			honestPeer1RecvChan <- msg
			return nil
		},
	}); err != nil {
		t.Fatal(err)
	}

	ensureNonMaliciousMessage := func(t *testing.T, msg net.Message) error {
		testPayload, ok := msg.Payload().(*testMessage)
		if !ok {
			return fmt.Errorf(
				"expected: payload type string\ngot:   payload type [%v]",
				testPayload,
			)
		}

		if honestPayload != testPayload.Payload {
			return fmt.Errorf(
				"expected: message payload [%s]\ngot:   payload [%s]",
				honestPayload,
				testPayload.Payload,
			)
		}
		return nil
	}

	for {
		select {
		case msg := <-honestPeer2RecvChan:
			if err := ensureNonMaliciousMessage(t, msg); err != nil {
				t.Fatal(err)
			}

		case msg := <-honestPeer1RecvChan:
			if err := ensureNonMaliciousMessage(t, msg); err != nil {
				t.Fatal(err)
			}

		// Ensure all messages are flushed before exiting
		case <-time.After(2 * time.Second):
			return
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
