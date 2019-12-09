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

func TestReceiveRetransmissions(t *testing.T) {
	retransmitter := newRetransmitter(1, 1)

	var received []pb.NetworkMessage

	message := &pb.NetworkMessage{
		Sender:         []byte("this is sender"),
		Payload:        []byte("this is payload"),
		Retransmission: 0,
	}
	retransmitter.receive(message, func() error {
		received = append(received, *message)
		return nil
	})

	message.Retransmission = 1
	retransmitter.receive(message, func() error {
		received = append(received, *message)
		return nil
	})

	if len(received) != 1 {
		t.Fatalf("only one message should be accepted")
	}

	if received[0].Retransmission != 0 {
		t.Errorf("only the original message should be accepted")
	}
}

func TestPassReceivedUniqueMessages(t *testing.T) {
	retransmitter := newRetransmitter(1, 1)

	var received []*pb.NetworkMessage

	message1 := &pb.NetworkMessage{
		Sender:         []byte("donald duck"),
		Payload:        []byte("conan the barbarian"),
		Retransmission: 0,
	}
	message2 := &pb.NetworkMessage{
		Sender:         []byte("donald duck"),
		Payload:        []byte("conan the conqueror"),
		Retransmission: 0,
	}
	retransmitter.receive(message1, func() error {
		received = append(received, message1)
		return nil
	})

	retransmitter.receive(message2, func() error {
		received = append(received, message2)
		return nil
	})

	if len(received) != 2 {
		t.Fatalf("both unique messages should be accepted")
	}

	if received[0].Retransmission != 0 {
		t.Errorf("message is not a retransmission")
	}

	if received[1].Retransmission != 0 {
		t.Errorf("message is not a retransmission")
	}
}

func TestPassIdenticalMessages(t *testing.T) {
	retransmitter := newRetransmitter(1, 1)

	var received []*pb.NetworkMessage

	message1 := &pb.NetworkMessage{
		Sender:         []byte("masha"),
		Payload:        []byte("bear"),
		Retransmission: 0,
	}
	message2 := &pb.NetworkMessage{
		Sender:         []byte("masha"),
		Payload:        []byte("bear"),
		Retransmission: 0,
	}
	retransmitter.receive(message1, func() error {
		received = append(received, message1)
		return nil
	})

	retransmitter.receive(message2, func() error {
		received = append(received, message2)
		return nil
	})

	if len(received) != 2 {
		t.Fatalf("both identical messages should be accepted")
	}

	if received[0].Retransmission != 0 {
		t.Errorf("message is not a retransmission")
	}

	if received[1].Retransmission != 0 {
		t.Errorf("message is not a retransmission")
	}
}

func TestCalculateFingerprint(t *testing.T) {
	tests := map[string]struct {
		message1                *pb.NetworkMessage
		message2                *pb.NetworkMessage
		sameFingerprintExpected bool
	}{
		"different sender": {
			message1:                &pb.NetworkMessage{Sender: []byte("donald duck")},
			message2:                &pb.NetworkMessage{Sender: []byte("mickey mouse")},
			sameFingerprintExpected: false,
		},
		"same sender": {
			message1:                &pb.NetworkMessage{Sender: []byte("mickey mouse")},
			message2:                &pb.NetworkMessage{Sender: []byte("mickey mouse")},
			sameFingerprintExpected: true,
		},
		"different encrypted flag": {
			message1:                &pb.NetworkMessage{Encrypted: true},
			message2:                &pb.NetworkMessage{Encrypted: false},
			sameFingerprintExpected: false,
		},
		"same encrypted flag": {
			message1:                &pb.NetworkMessage{Encrypted: true},
			message2:                &pb.NetworkMessage{Encrypted: true},
			sameFingerprintExpected: true,
		},
		"different payload": {
			message1:                &pb.NetworkMessage{Payload: []byte("rambo")},
			message2:                &pb.NetworkMessage{Payload: []byte("terminator")},
			sameFingerprintExpected: false,
		},
		"same payload": {
			message1:                &pb.NetworkMessage{Payload: []byte("rambo")},
			message2:                &pb.NetworkMessage{Payload: []byte("rambo")},
			sameFingerprintExpected: true,
		},
		"different type": {
			message1:                &pb.NetworkMessage{Type: []byte("Om")},
			message2:                &pb.NetworkMessage{Type: []byte("Nom")},
			sameFingerprintExpected: false,
		},
		"same type": {
			message1:                &pb.NetworkMessage{Type: []byte("Om")},
			message2:                &pb.NetworkMessage{Type: []byte("Om")},
			sameFingerprintExpected: true,
		},
		"different channel": {
			message1:                &pb.NetworkMessage{Type: []byte("US")},
			message2:                &pb.NetworkMessage{Type: []byte("EU")},
			sameFingerprintExpected: false,
		},
		"same channel": {
			message1:                &pb.NetworkMessage{Type: []byte("US")},
			message2:                &pb.NetworkMessage{Type: []byte("US")},
			sameFingerprintExpected: true,
		},
		"different retransmission": {
			message1: &pb.NetworkMessage{
				Type:           []byte("US"),
				Retransmission: 0,
			},
			message2: &pb.NetworkMessage{
				Type:           []byte("US"),
				Retransmission: 1,
			},
			sameFingerprintExpected: true,
		},
		"same retransmission": {
			message1: &pb.NetworkMessage{
				Type:           []byte("US"),
				Retransmission: 0,
			},
			message2: &pb.NetworkMessage{
				Type:           []byte("US"),
				Retransmission: 0,
			},
			sameFingerprintExpected: true,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			fingerprint1, err := calculateFingerprint(test.message1)
			if err != nil {
				t.Fatal(err)
			}
			fingerprint2, err := calculateFingerprint(test.message2)
			if err != nil {
				t.Fatal(err)
			}

			if fingerprint1 == fingerprint2 && !test.sameFingerprintExpected {
				t.Errorf("same fingerprints")
			}

			if fingerprint1 != fingerprint2 && test.sameFingerprintExpected {
				t.Errorf("different fingerprints")
			}
		})
	}
}
