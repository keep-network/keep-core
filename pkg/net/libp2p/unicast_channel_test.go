package libp2p

import (
	"context"
	"github.com/keep-network/keep-core/pkg/operator"
	"strconv"
	"testing"
	"time"

	"github.com/keep-network/keep-core/pkg/firewall"
	"github.com/keep-network/keep-core/pkg/net/retransmission"

	"github.com/keep-network/keep-core/pkg/net"
	"github.com/multiformats/go-multiaddr"
)

func TestSendReceiveUnicastChannel(t *testing.T) {
	ctx := context.Background()

	withNetwork(ctx, t, 9000, func(
		identity1 *identity,
		identity2 *identity,
		provider1 net.Provider,
		provider2 net.Provider,
	) {
		// Peer 1 initiates the channel.
		peer1Channel, err := provider1.UnicastChannelWith(identity2.id)
		if err != nil {
			t.Fatal(err)
		}

		// Peer 1 registers a message receiver and unmarshaller on the channel instance.
		peer1Receiver := newMessageReceiver("peer1Receiver")
		peer1Channel.Recv(ctx, peer1Receiver.receive)
		peer1Channel.SetUnmarshaler(func() net.TaggedUnmarshaler {
			return &testMessage{}
		})

		// Peer 2 registers a message receiver and unmarshaller using `OnUnicastChannelOpened`.
		onUnicastChannelOpenedInvocations := 0
		peer2Receiver := newMessageReceiver("peer2Receiver")
		provider2.OnUnicastChannelOpened(func(channel net.UnicastChannel) {
			onUnicastChannelOpenedInvocations++
			channel.Recv(ctx, peer2Receiver.receive)
			channel.SetUnmarshaler(func() net.TaggedUnmarshaler {
				return &testMessage{}
			})
		})

		// Peer 1 prepares and sends messages.
		peer1Messages := []testMessage{
			{Sender: identity1, Recipient: identity2, Payload: "one"},
			{Sender: identity1, Recipient: identity2, Payload: "two"},
			{Sender: identity1, Recipient: identity2, Payload: "three"},
		}
		for _, message := range peer1Messages {
			send := withRetry(func() error {
				return peer1Channel.Send(&message)
			}, 3, 10*time.Second)
			if err := send; err != nil {
				t.Fatal(err)
			}
		}

		// Peer 2 get the channel from cache.
		peer2Channel, err := provider2.UnicastChannelWith(identity1.id)
		if err != nil {
			t.Fatal(err)
		}

		// Peer 2 prepares and sends messages.
		peer2Messages := []testMessage{
			{Sender: identity1, Recipient: identity2, Payload: "four"},
			{Sender: identity1, Recipient: identity2, Payload: "five"},
		}
		for _, message := range peer2Messages {
			send := withRetry(func() error {
				return peer2Channel.Send(&message)
			}, 3, 10*time.Second)
			if err := send; err != nil {
				t.Fatal(err)
			}
		}

		// Wait a bit, messages must be sent and received.
		time.Sleep(2 * time.Second)

		assertReceivedMessages(t, peer1Receiver, peer2Messages)
		assertReceivedMessages(t, peer2Receiver, peer1Messages)

		expectedOnUnicastChannelOpenedInvocations := 1
		if onUnicastChannelOpenedInvocations != expectedOnUnicastChannelOpenedInvocations {
			t.Errorf(
				"unexpected number of `OnUnicastChannelOpened` invocations\nactual:   %v\nexpected: %v\n",
				onUnicastChannelOpenedInvocations,
				expectedOnUnicastChannelOpenedInvocations,
			)
		}
	})
}

func withRetry(function func() error, retryCount int, waitTime time.Duration) error {
	var err error

	for i := 0; i < retryCount+1; i++ {
		err = function()
		if err == nil {
			return nil
		}
		time.Sleep(waitTime)
	}

	return err
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

func withNetwork(
	ctx context.Context,
	t *testing.T,
	startingPort int,
	testFn func(
		identity1 *identity,
		identity2 *identity,
		provider1 net.Provider,
		provider2 net.Provider,
	),
) {
	operatorPrivateKey1, _, err := operator.GenerateKeyPair(DefaultCurve)
	if err != nil {
		t.Fatal(err)
	}

	networkPrivateKey1, _, err := OperatorPrivateKeyToNetworkKeyPair(operatorPrivateKey1)
	if err != nil {
		t.Fatal(err)
	}

	identity1, err := createIdentity(networkPrivateKey1)
	if err != nil {
		t.Fatal(err)
	}

	provider1Port := startingPort

	multiaddr1, err := multiaddr.NewMultiaddr(
		"/ip4/127.0.0.1/tcp/" + strconv.Itoa(provider1Port),
	)
	if err != nil {
		t.Fatal(err)
	}

	operatorPrivateKey2, _, err := operator.GenerateKeyPair(DefaultCurve)
	if err != nil {
		t.Fatal(err)
	}

	networkPrivateKey2, _, err := OperatorPrivateKeyToNetworkKeyPair(operatorPrivateKey2)
	if err != nil {
		t.Fatal(err)
	}

	identity2, err := createIdentity(networkPrivateKey2)
	if err != nil {
		t.Fatal(err)
	}

	provider2Port := startingPort + 1

	multiaddr2, err := multiaddr.NewMultiaddr(
		"/ip4/127.0.0.1/tcp/" + strconv.Itoa(provider2Port),
	)
	if err != nil {
		t.Fatal(err)
	}

	provider1, err := Connect(
		ctx,
		Config{
			Peers: []string{
				multiaddressWithIdentity(
					multiaddr2,
					identity2.id,
				),
			},
			Port: provider1Port,
		},
		operatorPrivateKey1,
		ProtocolBeacon,
		firewall.Disabled,
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
				),
			},
			Port: provider2Port,
		},
		operatorPrivateKey2,
		ProtocolBeacon,
		firewall.Disabled,
		retransmission.NewTicker(make(chan uint64)),
	)
	if err != nil {
		t.Fatal(err)
	}

	testFn(
		identity1,
		identity2,
		provider1,
		provider2,
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
