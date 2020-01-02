package retransmission

import (
	"testing"

	"github.com/keep-network/keep-core/pkg/net"
)

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
