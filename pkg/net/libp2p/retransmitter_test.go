package libp2p

import (
	"testing"
	"time"

	"github.com/keep-network/keep-core/pkg/internal/testutils"
	"github.com/keep-network/keep-core/pkg/net/gen/pb"
)

func TestRetransmitExpectedNumberOfTimes(t *testing.T) {
	cycles := 5
	interval := 10

	message := &pb.NetworkMessage{}
	retransmitter := newRetransmitter(cycles, interval)

	retransmissions := make(chan *pb.NetworkMessage, cycles)

	retransmitter.scheduleRetransmission(
		message, func(msg *pb.NetworkMessage) error {
			retransmissions <- msg
			return nil
		},
	)

	time.Sleep(60 * time.Millisecond)

	if len(retransmissions) != cycles {
		t.Errorf(
			"unexpected number of retransmissions\nactual: [%v]\nexpected:   [%v]",
			cycles,
			len(retransmissions),
		)
	}
}

func TestUpdateRetransmissionCounter(t *testing.T) {
	cycles := 3
	interval := 10

	message := &pb.NetworkMessage{}
	retransmitter := newRetransmitter(cycles, interval)

	retransmissions := make(chan *pb.NetworkMessage, cycles)

	retransmitter.scheduleRetransmission(
		message, func(msg *pb.NetworkMessage) error {
			retransmissions <- msg
			return nil
		},
	)

	time.Sleep(40 * time.Millisecond)

	for i := 1; i <= cycles; i++ {
		message := <-retransmissions
		if uint32(i) != message.Retransmission {
			t.Errorf(
				"unexpected retransmission counter\nactual: [%v]\nexpected:   [%v]",
				i,
				message.Retransmission,
			)
		}
	}
}

func TestRetransmitOriginalContent(t *testing.T) {
	cycles := 3
	interval := 10

	sender := []byte("this is sender")
	encrypted := true
	payload := []byte("this is payload")
	messageType := []byte("this is type")
	channel := []byte("this is channel")

	message := &pb.NetworkMessage{
		Sender:    sender,
		Encrypted: encrypted,
		Payload:   payload,
		Type:      messageType,
		Channel:   channel,
	}

	retransmitter := newRetransmitter(cycles, interval)

	retransmissions := make(chan *pb.NetworkMessage, cycles)

	retransmitter.scheduleRetransmission(
		message, func(msg *pb.NetworkMessage) error {
			retransmissions <- msg
			return nil
		},
	)

	time.Sleep(40 * time.Millisecond)

	for i := 1; i <= cycles; i++ {
		message := <-retransmissions

		testutils.AssertBytesEqual(t, sender, message.Sender)
		testutils.AssertBytesEqual(t, payload, message.Payload)
		testutils.AssertBytesEqual(t, messageType, message.Type)
		testutils.AssertBytesEqual(t, channel, message.Channel)
		if encrypted != message.Encrypted {
			t.Errorf("unexpected 'Encrypted' field value")
		}
	}
}

func TestDoNotModifyOriginalMessage(t *testing.T) {
	cycles := 1
	interval := 10

	message := &pb.NetworkMessage{}
	retransmitter := newRetransmitter(cycles, interval)

	retransmissions := make(chan *pb.NetworkMessage, cycles)

	retransmitter.scheduleRetransmission(
		message, func(msg *pb.NetworkMessage) error {
			retransmissions <- msg
			return nil
		},
	)

	time.Sleep(10 * time.Millisecond)

	if message.Retransmission != 0 {
		t.Errorf("modified retransmission field of the original message")
	}
}
