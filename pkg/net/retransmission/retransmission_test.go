package retransmission

import (
	"context"
	"testing"
	"time"

	"github.com/keep-network/keep-core/pkg/internal/testutils"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/net/gen/pb"
)

func TestRetransmitExpectedNumberOfTimes(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 510*time.Millisecond)
	defer cancel()

	timeTicker := time.NewTicker(50 * time.Millisecond)
	defer timeTicker.Stop()

	retransmissionsCount := 0
	ScheduleRetransmissions(
		ctx,
		NewTimeTicker(timeTicker),
		&pb.NetworkMessage{},
		func(msg *pb.NetworkMessage) error {
			retransmissionsCount++
			return nil
		},
	)

	<-ctx.Done()

	if retransmissionsCount != 10 {
		t.Errorf("expected [10] retransmissions, has [%v]", retransmissionsCount)
	}
}

func TestUpdateRetransmissionCounter(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 260*time.Millisecond)
	defer cancel()

	timeTicker := time.NewTicker(50 * time.Millisecond)
	defer timeTicker.Stop()

	var retransmissions []*pb.NetworkMessage
	ScheduleRetransmissions(
		ctx,
		NewTimeTicker(timeTicker),
		&pb.NetworkMessage{},
		func(msg *pb.NetworkMessage) error {
			retransmissions = append(retransmissions, msg)
			return nil
		},
	)

	<-ctx.Done()

	for i := 1; i <= 5; i++ {
		message := retransmissions[i-1]
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

	ctx, cancel := context.WithTimeout(context.Background(), 110*time.Millisecond)
	defer cancel()

	timeTicker := time.NewTicker(50 * time.Millisecond)
	defer timeTicker.Stop()

	var retransmissions []*pb.NetworkMessage
	ScheduleRetransmissions(
		ctx,
		NewTimeTicker(timeTicker),
		message,
		func(msg *pb.NetworkMessage) error {
			retransmissions = append(retransmissions, msg)
			return nil
		},
	)

	<-ctx.Done()

	for i := 0; i < 2; i++ {
		message := retransmissions[i]

		testutils.AssertBytesEqual(t, sender, message.Sender)
		testutils.AssertBytesEqual(t, payload, message.Payload)
		testutils.AssertBytesEqual(t, messageType, message.Type)
		testutils.AssertBytesEqual(t, channel, message.Channel)
		if encrypted != message.Encrypted {
			t.Errorf("unexpected 'Encrypted' field value")
		}
	}
}

func TestHandlerReceiveMessage(t *testing.T) {
	var received []net.Message

	handler := WithRetransmissionSupport(func(message net.Message) {
		received = append(received, message)
	})

	handler(&mockNetworkMessage{
		fingerprint:    "ABC",
		retransmission: 0,
	})

	if len(received) != 1 {
		t.Fatalf(
			"unexpected number of accepted messages\nactual:   [%v]\nexpected: [1]",
			len(received),
		)
	}
}

func TestHandlerReceiveRetransmission(t *testing.T) {
	var received []net.Message

	handler := WithRetransmissionSupport(func(message net.Message) {
		received = append(received, message)
	})

	handler(&mockNetworkMessage{
		fingerprint:    "ABC",
		retransmission: 1,
	})

	if len(received) != 1 {
		t.Fatalf(
			"unexpected number of accepted messages\nactual:   [%v]\nexpected: [1]",
			len(received),
		)
	}
}

func TestHandlerReceiveMessageAndRetransmission(t *testing.T) {
	var received []net.Message

	handler := WithRetransmissionSupport(func(message net.Message) {
		received = append(received, message)
	})

	handler(&mockNetworkMessage{
		fingerprint:    "ABC",
		retransmission: 0,
	})

	handler(&mockNetworkMessage{
		fingerprint:    "ABC",
		retransmission: 1,
	})

	if len(received) != 1 {
		t.Fatalf(
			"unexpected number of accepted messages\nactual:   [%v]\nexpected: [1]",
			len(received),
		)
	}
}

func TestHandlerReceiveUniqueMessages(t *testing.T) {
	var received []net.Message

	handler := WithRetransmissionSupport(func(message net.Message) {
		received = append(received, message)
	})

	handler(&mockNetworkMessage{
		fingerprint:    "ABC",
		retransmission: 0,
	})

	handler(&mockNetworkMessage{
		fingerprint:    "DEF",
		retransmission: 0,
	})

	if len(received) != 2 {
		t.Fatalf(
			"unexpected number of accepted messages\nactual:   [%v]\nexpected: [2]",
			len(received),
		)
	}
}

func TestHandlerIdenticalMessages(t *testing.T) {
	var received []net.Message

	handler := WithRetransmissionSupport(func(message net.Message) {
		received = append(received, message)
	})

	handler(&mockNetworkMessage{
		fingerprint:    "ABC",
		retransmission: 0,
	})

	handler(&mockNetworkMessage{
		fingerprint:    "ABC",
		retransmission: 0,
	})

	if len(received) != 2 {
		t.Fatalf(
			"unexpected number of accepted messages\nactual:   [%v]\nexpected: [2]",
			len(received),
		)
	}
}

type mockNetworkMessage struct {
	fingerprint    string
	retransmission uint32
}

func (mnm *mockNetworkMessage) Fingerprint() string {
	return mnm.fingerprint
}

func (mnm *mockNetworkMessage) Retransmission() uint32 {
	return mnm.retransmission
}

func (mnm *mockNetworkMessage) TransportSenderID() net.TransportIdentifier {
	panic("not implemented")
}

func (mnm *mockNetworkMessage) Payload() interface{} {
	panic("not implemented")
}

func (mnm *mockNetworkMessage) Type() string {
	panic("not implemented")
}

func (mnm *mockNetworkMessage) SenderPublicKey() []byte {
	panic("not implemented")
}

func TestCalculateFingerprint(t *testing.T) {
	tests := map[string]struct {
		sender1                 string
		sender2                 string
		payload1                string
		payload2                string
		sameFingerprintExpected bool
	}{
		"different sender, same payload": {
			sender1:                 "bob",
			sender2:                 "alice",
			payload1:                "blockchain",
			payload2:                "blockchain",
			sameFingerprintExpected: false,
		},
		"same sender, different payload": {
			sender1:                 "bob",
			sender2:                 "bob",
			payload1:                "btc",
			payload2:                "eth",
			sameFingerprintExpected: false,
		},
		"same sender, same payload": {
			sender1:                 "alice",
			sender2:                 "alice",
			payload1:                "blockchain",
			payload2:                "blockchain",
			sameFingerprintExpected: true,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			fingerprint1 := CalculateFingerprint(
				localIdentifier(test.sender1),
				[]byte(test.payload1),
			)

			fingerprint2 := CalculateFingerprint(
				localIdentifier(test.sender2),
				[]byte(test.payload2),
			)

			if fingerprint1 == fingerprint2 && !test.sameFingerprintExpected {
				t.Errorf("same fingerprints")
			}

			if fingerprint1 != fingerprint2 && test.sameFingerprintExpected {
				t.Errorf("different fingerprints")
			}
		})
	}
}

type localIdentifier string

func (li localIdentifier) String() string {
	return string(li)
}
