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
		provider1 net.Provider,
		provider2 net.Provider,
	) {
		_, err := provider1.UnicastChannelWith(identity2.id)
		if err != nil {
			t.Fatal(err)
		}

		_, err = provider2.UnicastChannelWith(identity1.id)
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestSendReceiveUnicastChannel(t *testing.T) {
	ctx := context.Background()

	withNetwork(ctx, t, func(
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
	testFn func(
		identity1 *identity,
		identity2 *identity,
		provider1 net.Provider,
		provider2 net.Provider,
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

	stakeMonitor := local.NewStakeMonitor(big.NewInt(0))

	provider1, err := Connect(
		ctx,
		Config{
			Peers: []string{
				multiaddressWithIdentity(
					multiaddr2,
					identity2.id,
				),
			},
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
				),
			},
			Port: 8082,
		},
		privKey2,
		stakeMonitor,
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
