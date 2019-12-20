package interception

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/keep-network/keep-core/pkg/net"
	netLocal "github.com/keep-network/keep-core/pkg/net/local"
)

func TestPassThruNetworkMessage(t *testing.T) {
	network := NewNetwork(
		netLocal.Connect(),
		func(msg net.TaggedMarshaler) net.TaggedMarshaler {
			// return message unchanged
			return msg
		},
	)

	channel, err := network.ChannelFor("badger288")
	if err != nil {
		t.Fatal(err)
	}

	inputMessage := &testMessage{"hello world"}
	onMessage := func(message *testMessage) {
		if !reflect.DeepEqual(message, inputMessage) {
			t.Errorf(
				"Unexpected message\nExpected: [%v]\nActual:   [%v]",
				inputMessage,
				message,
			)
		}
	}
	onTimeout := func() {
		t.Errorf("No message received")
	}

	testMessageRoundtrip(channel, inputMessage, onMessage, onTimeout)
}

func TestModifyNetworkMessage(t *testing.T) {
	modifiedMessage := &testMessage{"modified"}
	network := NewNetwork(
		netLocal.Connect(),
		func(msg net.TaggedMarshaler) net.TaggedMarshaler {
			// alter the message
			return modifiedMessage
		},
	)

	channel, err := network.ChannelFor("badger288")
	if err != nil {
		t.Fatal(err)
	}

	inputMessage := &testMessage{"hello world"}
	onMessage := func(message *testMessage) {
		if !reflect.DeepEqual(message, modifiedMessage) {
			t.Errorf(
				"Unexpected message\nExpected: [%v]\nActual:   [%v]",
				inputMessage,
				message,
			)
		}
	}
	onTimeout := func() {
		t.Errorf("No message received")
	}

	testMessageRoundtrip(channel, inputMessage, onMessage, onTimeout)
}

func TestDropNetworkMessage(t *testing.T) {
	network := NewNetwork(
		netLocal.Connect(),
		func(msg net.TaggedMarshaler) net.TaggedMarshaler {
			// return message unchanged
			return nil
		},
	)

	channel, err := network.ChannelFor("badger288")
	if err != nil {
		t.Fatal(err)
	}

	inputMessage := &testMessage{"hello world"}
	onMessage := func(message *testMessage) {
		t.Errorf("No message expected. Received [%v]", message)
	}
	onTimeout := func() {
		// ok, expected
	}

	testMessageRoundtrip(channel, inputMessage, onMessage, onTimeout)
}

func testMessageRoundtrip(
	channel net.BroadcastChannel,
	message *testMessage,
	onMessageReceived func(message *testMessage),
	onTimeout func(),
) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	channel.RegisterUnmarshaler(func() net.TaggedUnmarshaler {
		return &testMessage{}
	})

	handlerFiredChan := make(chan *testMessage)
	handler := net.HandleMessageFunc{
		Type: "test_message",
		Handler: func(msg net.Message) error {
			handlerFiredChan <- msg.Payload().(*testMessage)
			return nil
		},
	}

	channel.Recv(handler)
	channel.Send(ctx, message)

	select {
	case msg := <-handlerFiredChan:
		onMessageReceived(msg)
	case <-ctx.Done():
		onTimeout()
	}
}

type testMessage struct {
	payload string
}

func (tm *testMessage) Type() string {
	return "test_message"
}

func (tm *testMessage) Marshal() ([]byte, error) {
	return []byte(tm.payload), nil
}

func (tm *testMessage) Unmarshal(bytes []byte) error {
	tm.payload = string(bytes)
	return nil
}
