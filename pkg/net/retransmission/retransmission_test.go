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
