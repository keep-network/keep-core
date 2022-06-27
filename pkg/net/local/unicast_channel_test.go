package local

import (
	"context"
	"github.com/keep-network/keep-core/pkg/operator"
	"testing"
	"time"

	"github.com/keep-network/keep-core/pkg/net"
)

func TestNewChannelNotification(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	peer1Provider, _ := initTestProvider(t)
	peer2Provider, peer2OperatorPublicKey := initTestProvider(t)

	peer1NewChannelNotificationCount := 0
	peer1Provider.OnUnicastChannelOpened(func(channel net.UnicastChannel) {
		peer1NewChannelNotificationCount++
	})

	peer2NewChannelNotificationCount := 0
	peer2Provider.OnUnicastChannelOpened(func(channel net.UnicastChannel) {
		peer2NewChannelNotificationCount++
	})

	remotePeerID, err := createLocalIdentifier(peer2OperatorPublicKey)
	if err != nil {
		t.Fatal(err)
	}

	_, err = peer1Provider.UnicastChannelWith(remotePeerID)
	if err != nil {
		t.Fatal(err)
	}

	<-ctx.Done() // give some time for notifications...

	if peer1NewChannelNotificationCount != 0 {
		t.Errorf(
			"expected no notifications, has [%v]",
			peer1NewChannelNotificationCount,
		)
	}
	if peer2NewChannelNotificationCount != 1 {
		t.Errorf(
			"expected [1] notification, has [%v]",
			peer2NewChannelNotificationCount,
		)
	}
}

func TestExistingChannelNotification(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	peer1Provider, _ := initTestProvider(t)
	peer2Provider, peer2OperatorPublicKey := initTestProvider(t)

	newChannelNotificationCount := 0
	peer2Provider.OnUnicastChannelOpened(func(channel net.UnicastChannel) {
		newChannelNotificationCount++
	})

	remotePeerID, err := createLocalIdentifier(peer2OperatorPublicKey)
	if err != nil {
		t.Fatal(err)
	}

	_, err = peer1Provider.UnicastChannelWith(remotePeerID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = peer1Provider.UnicastChannelWith(remotePeerID)
	if err != nil {
		t.Fatal(err)
	}

	<-ctx.Done() // give some time for notifications...

	if newChannelNotificationCount != 1 {
		t.Errorf(
			"expected [1] notification, has [%v]",
			newChannelNotificationCount,
		)
	}
}

func TestSendAndReceive(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	//
	// Prepare communication channel between peer1 and peer2
	//
	peer1Provider, peer1OperatorPublicKey := initTestProvider(t)
	peer2Provider, peer2OperatorPublicKey := initTestProvider(t)

	remotePeer1ID, err := createLocalIdentifier(peer1OperatorPublicKey)
	if err != nil {
		t.Fatal(err)
	}

	remotePeer2ID, err := createLocalIdentifier(peer2OperatorPublicKey)
	if err != nil {
		t.Fatal(err)
	}

	channel1, err := peer1Provider.UnicastChannelWith(remotePeer2ID)
	if err != nil {
		t.Fatal(err)
	}
	channel2, err := peer2Provider.UnicastChannelWith(remotePeer1ID)
	if err != nil {
		t.Fatal(err)
	}

	channel1.SetUnmarshaler(func() net.TaggedUnmarshaler {
		return &mockMessage{}
	})
	channel2.SetUnmarshaler(func() net.TaggedUnmarshaler {
		return &mockMessage{}
	})

	peer1Received := make(chan net.Message)
	peer2Received := make(chan net.Message)

	channel1.Recv(ctx, func(msg net.Message) {
		peer1Received <- msg
	})
	channel2.Recv(ctx, func(msg net.Message) {
		peer2Received <- msg
	})

	//
	// peer1 sends a message to peer2
	// make sure peer2 receives it
	//

	channel1Message := &mockMessage{"yolo1"}
	err = channel1.Send(channel1Message)
	if err != nil {
		t.Fatal(err)
	}

	select {
	case msg := <-peer2Received:
		switch message := msg.Payload().(type) {
		case *mockMessage:
			if message.content != channel1Message.content {
				t.Fatalf(
					"unexpected message content\nactual:   [%v]\nexpected: [%v]",
					message.content,
					channel1Message.content,
				)
			}
		default:
			t.Fatal("unexpected message type")
		}

	case <-peer1Received:
		t.Fatal("peer 1 should not receive this message")
	case <-ctx.Done():
		t.Fatal("expected message not arrived to peer 2")
	}

	//
	// peer2 sends a message to peer1
	// make sure peer1 receives it
	//

	channel2Message := &mockMessage{"yolo2"}
	err = channel2.Send(channel2Message)
	if err != nil {
		t.Fatal(err)
	}

	select {
	case msg := <-peer1Received:
		switch message := msg.Payload().(type) {
		case *mockMessage:
			if message.content != channel2Message.content {
				t.Fatalf(
					"unexpected message content\nactual:   [%v]\nexpected: [%v]",
					message.content,
					channel2Message.content,
				)
			}
		default:
			t.Fatal("unexpected message type")
		}
	case <-peer2Received:
		t.Fatal("peer 2 should not receive this message")
	case <-ctx.Done():
		t.Fatal("expected message not arrived")
	}
}

func TestTalkToSelf(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	//
	// Prepare self-communication channel (e.g. two goroutines)
	//
	peerProvider, peerOperatorPublicKey := initTestProvider(t)

	peerTransportID, err := createLocalIdentifier(peerOperatorPublicKey)
	if err != nil {
		t.Fatal(err)
	}

	channel1, err := peerProvider.UnicastChannelWith(peerTransportID)
	if err != nil {
		t.Fatal(err)
	}
	channel2, err := peerProvider.UnicastChannelWith(peerTransportID)
	if err != nil {
		t.Fatal(err)
	}

	channel1.SetUnmarshaler(func() net.TaggedUnmarshaler {
		return &mockMessage{}
	})

	chan1Received := make(chan net.Message)
	chan2Received := make(chan net.Message)

	channel1.Recv(ctx, func(msg net.Message) {
		chan1Received <- msg
	})
	channel2.Recv(ctx, func(msg net.Message) {
		chan2Received <- msg
	})

	//
	// send message to self via the first channel
	// both handlers receive it
	//

	err = channel1.Send(&mockMessage{"yolo1"})
	if err != nil {
		t.Fatal(err)
	}

	select {
	case <-chan1Received: // ok
	case <-ctx.Done():
		t.Fatal("expected message not arrived")
	}

	select {
	case <-chan2Received: // ok
	case <-ctx.Done():
		t.Fatal("expected message not arrived")
	}

	//
	// send message to self via the second channel
	// again, both handlers should receive it
	//

	err = channel2.Send(&mockMessage{"yolo2"})
	if err != nil {
		t.Fatal(err)
	}

	select {
	case <-chan1Received: // ok
	case <-ctx.Done():
		t.Fatal("expected message not arrived")
	}

	select {
	case <-chan2Received: // ok
	case <-ctx.Done():
		t.Fatal("expected message not arrived")
	}
}

func TestTimedOutHandlerNotReceiveUnicastMessage(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx2, cancel2 := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel2()

	peer1ID := localIdentifier("peer-0xAAEF12")
	_, peer1OperatorPublicKey, _ := operator.GenerateKeyPair(DefaultCurve)

	peer2ID := localIdentifier("peer-0x121211")

	unicastChannel := newUnicastChannel(peer1ID, peer1OperatorPublicKey, peer2ID)
	unicastChannel.SetUnmarshaler(func() net.TaggedUnmarshaler {
		return &mockMessage{}
	})

	received := make(chan net.Message)
	unicastChannel.Recv(ctx, func(msg net.Message) {
		received <- msg
	})

	received2 := make(chan net.Message)
	unicastChannel.Recv(ctx2, func(msg net.Message) {
		received2 <- msg
	})

	cancel() // cancel the first context

	message := &mockMessage{"hello"}
	marshaled, err := message.Marshal()
	if err != nil {
		t.Fatal(err)
	}

	unicastChannel.receiveMessage(marshaled, message.Type())

	select {
	case <-received:
		t.Fatal("receiver should not be called")
	default:
		// ok, should not receive
	}

	select {
	case <-ctx2.Done():
		t.Fatal("expected message not received")
	case <-received2:
		// ok, should receive
	}
}

func initTestProvider(t *testing.T) (net.Provider, *operator.PublicKey) {
	_, operatorPublicKey, err := operator.GenerateKeyPair(DefaultCurve)
	if err != nil {
		t.Fatal(err)
	}

	provider := ConnectWithKey(operatorPublicKey)

	return provider, operatorPublicKey
}

type mockMessage struct {
	content string
}

func (mm *mockMessage) Type() string {
	return "mock_message"
}

func (mm *mockMessage) Marshal() ([]byte, error) {
	return []byte(mm.content), nil
}

func (mm *mockMessage) Unmarshal(bytes []byte) error {
	mm.content = string(bytes)
	return nil
}
