package libp2p

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/keep-network/keep-core/pkg/net/retransmission"

	"github.com/keep-network/keep-core/pkg/chain/local"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/net/key"
	"github.com/multiformats/go-multiaddr"
)

func TestCreateUnicastChannel(t *testing.T) {
	ctx := context.Background()

	withNetwork(ctx, t, func(
		identity1 *identity,
		identity2 *identity,
		identity3 *identity,
		provider1 net.Provider,
		provider2 net.Provider,
		provider3 net.Provider,
	) {
		_, err := provider1.ChannelWith(identity2.id.String())
		if err != nil {
			t.Fatal(err)
		}

		_, err = provider2.ChannelWith(identity1.id.String())
		if err != nil {
			t.Fatal(err)
		}
	})
}

// Checks if two direct peers can exchange messages through an unicast channel.
// In this scenario, peer 2 has peer's 1 address in their config and open
// a direct connection during their cyclic `core bootstrap` round.
func TestSendUnicastChannel_MessagesBetweenDirectPeers(t *testing.T) {
	ctx := context.Background()

	withNetwork(ctx, t, func(
		identity1 *identity,
		identity2 *identity,
		identity3 *identity,
		provider1 net.Provider,
		provider2 net.Provider,
		provider3 net.Provider,
	) {
		// Channel instance of peer 1.
		channel1, err := provider1.ChannelWith(identity2.id.String())
		if err != nil {
			t.Fatal(err)
		}

		// Channel instance of peer 2.
		channel2, err := provider2.ChannelWith(identity1.id.String())
		if err != nil {
			t.Fatal(err)
		}

		if err := channel1.RegisterUnmarshaler(
			func() net.TaggedUnmarshaler { return &testMessage{} },
		); err != nil {
			t.Fatal(err)
		}

		if err := channel2.RegisterUnmarshaler(
			func() net.TaggedUnmarshaler { return &testMessage{} },
		); err != nil {
			t.Fatal(err)
		}

		// Register first handler for channel 1.
		channel1Receiver1 := newMessageReceiver("channel1Receiver1")
		channel1.Recv(ctx, channel1Receiver1.receive)

		// Register second handler for channel 1.
		channel1Receiver2 := newMessageReceiver("channel1Receiver2")
		channel1.Recv(ctx, channel1Receiver2.receive)

		// Register first handler for channel 2.
		channel2Receiver1 := newMessageReceiver("channel2Receiver1")
		channel2.Recv(ctx, channel2Receiver1.receive)

		// Prepare and send messages to channel 1.
		messagesToChannel1 := []testMessage{
			{Sender: identity1, Recipient: identity2, Payload: "one"},
			{Sender: identity1, Recipient: identity2, Payload: "two"},
			{Sender: identity1, Recipient: identity2, Payload: "three"},
		}
		go func() {
			for _, message := range messagesToChannel1 {
				if err := channel1.Send(ctx, &message); err != nil {
					t.Fatal(err)
				}
			}
		}()

		// Prepare and send messages to channel 2.
		messagesToChannel2 := []testMessage{
			{Sender: identity2, Recipient: identity1, Payload: "four"},
			{Sender: identity2, Recipient: identity1, Payload: "five"},
		}
		go func() {
			for _, message := range messagesToChannel2 {
				if err := channel2.Send(ctx, &message); err != nil {
					t.Fatal(err)
				}
			}
		}()

		// Wait a bit, messages must be sent and received.
		time.Sleep(10 * time.Second)

		assertReceivedMessages(t, channel1Receiver1, messagesToChannel2)
		assertReceivedMessages(t, channel1Receiver2, messagesToChannel2)
		assertReceivedMessages(t, channel2Receiver1, messagesToChannel1)
	})
}

// Checks if two discovered peers can exchange messages through an unicast channel.
// In this scenario, peers 1 and 3 discover themselves by performing a query to the
// intermediary peer 2, during their cyclic `DHT bootstrap` rounds. As a result
// a direct connection between peers 1 and 3 should be opened.
func TestSendUnicastChannel_MessagesBetweenDiscoveredPeers(t *testing.T) {
	ctx := context.Background()

	withNetwork(ctx, t, func(
		identity1 *identity,
		identity2 *identity,
		identity3 *identity,
		provider1 net.Provider,
		provider2 net.Provider,
		provider3 net.Provider,
	) {
		// Channel instance of peer 1.
		channel1, err := provider1.ChannelWith(identity3.id.String())
		if err != nil {
			t.Fatal(err)
		}

		// Channel instance of peer 3.
		channel3, err := provider3.ChannelWith(identity1.id.String())
		if err != nil {
			t.Fatal(err)
		}

		if err := channel3.RegisterUnmarshaler(
			func() net.TaggedUnmarshaler { return &testMessage{} },
		); err != nil {
			t.Fatal(err)
		}

		// Register handler for channel 3.
		channel3Receiver := newMessageReceiver("channel3Receiver")
		channel3.Recv(ctx, channel3Receiver.receive)

		// Prepare and send messages to channel 1.
		messagesToChannel1 := []testMessage{
			{Sender: identity1, Recipient: identity2, Payload: "one"},
			{Sender: identity1, Recipient: identity2, Payload: "two"},
		}
		go func() {
			for _, message := range messagesToChannel1 {
				if err := channel1.Send(ctx, &message); err != nil {
					t.Fatal(err)
				}
			}
		}()

		// Wait a bit, messages must be sent and received.
		time.Sleep(10 * time.Second)

		assertReceivedMessages(t, channel3Receiver, messagesToChannel1)
	})
}

func assertReceivedMessages(
	t *testing.T,
	receiver *messageReceiver,
	expectedMessages []testMessage,
) {
	if len(receiver.messages) != len(expectedMessages) {
		t.Errorf(
			"[%v] unexpected number of messages\nactual:   %v\nexpected: %v\n",
			receiver.name,
			len(receiver.messages),
			len(expectedMessages),
		)
	}

	for _, expectedMessage := range expectedMessages {
		isReceived := false
		for _, message := range receiver.messages {
			if message.Payload == expectedMessage.Payload {
				isReceived = true
				break
			}
		}

		if !isReceived {
			t.Errorf(
				"[%v] expected message [%v] not received",
				receiver.name,
				expectedMessage.Payload,
			)
		}
	}
}

// Builds a network consisting of three peers. At the beginning:
// - Peer 1 waits for incoming connections passively
// - Peer 2 tries to connect to peer 1
// - Peer 3 tries to connect to peer 2
// Eventually, all peers should be connected together due to libp2p peer discovery.
func withNetwork(
	ctx context.Context,
	t *testing.T,
	testFn func(
		identity1 *identity,
		identity2 *identity,
		identity3 *identity,
		provider1 net.Provider,
		provider2 net.Provider,
		provider3 net.Provider,
	),
) {
	privKey1, _, err := key.GenerateStaticNetworkKey()
	if err != nil {
		t.Fatal(err)
	}

	identity1, err := createIdentity(privKey1)
	if err != nil {
		t.Fatal(err)
	}

	multiaddr1, err := multiaddr.NewMultiaddr("/ip4/127.0.0.1/tcp/8081")
	if err != nil {
		t.Fatal(err)
	}

	privKey2, _, err := key.GenerateStaticNetworkKey()
	if err != nil {
		t.Fatal(err)
	}

	identity2, err := createIdentity(privKey2)
	if err != nil {
		t.Fatal(err)
	}

	multiaddr2, err := multiaddr.NewMultiaddr("/ip4/127.0.0.1/tcp/8082")
	if err != nil {
		t.Fatal(err)
	}

	privKey3, _, err := key.GenerateStaticNetworkKey()
	if err != nil {
		t.Fatal(err)
	}

	identity3, err := createIdentity(privKey3)
	if err != nil {
		t.Fatal(err)
	}

	stakeMonitor := local.NewStakeMonitor(big.NewInt(0))

	provider1, err := Connect(
		ctx,
		Config{
			Port: 8081,
		},
		privKey1,
		stakeMonitor,
		retransmission.NewTicker(make(chan uint64)),
	)
	if err != nil {
		t.Fatal(err)
	}

	provider2, err := Connect(
		ctx,
		Config{
			Peers: []string{
				multiaddressWithIdentity(
					multiaddr1,
					identity1.id,
				)},
			Port: 8082,
		},
		privKey2,
		stakeMonitor,
		retransmission.NewTicker(make(chan uint64)),
	)
	if err != nil {
		t.Fatal(err)
	}

	provider3, err := Connect(
		ctx,
		Config{
			Peers: []string{
				multiaddressWithIdentity(
					multiaddr2,
					identity2.id,
				)},
			Port: 8083,
		},
		privKey3,
		stakeMonitor,
		retransmission.NewTicker(make(chan uint64)),
	)
	if err != nil {
		t.Fatal(err)
	}

	// Wait peer discovery end for one minute.
	peerDiscoveryCompleted := false
	for i := 0; i < 60; i++ {
		if len(provider1.Peers()) == 2 &&
			len(provider2.Peers()) == 2 &&
			len(provider3.Peers()) == 2 {
			peerDiscoveryCompleted = true
			break
		}
		time.Sleep(1 * time.Second)
	}

	if !peerDiscoveryCompleted {
		t.Fatal("peer discovery timeout")
	}

	testFn(
		identity1,
		identity2,
		identity3,
		provider1,
		provider2,
		provider3,
	)
}

type messageReceiver struct {
	name     string
	messages []testMessage
}

func newMessageReceiver(name string) *messageReceiver {
	return &messageReceiver{
		name:     name,
		messages: make([]testMessage, 0),
	}
}

func (mr *messageReceiver) receive(message net.Message) {
	mr.messages = append(mr.messages, *message.Payload().(*testMessage))
}
