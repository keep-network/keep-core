package retransmission

import (
	"context"
	"github.com/keep-network/keep-core/pkg/internal/testutils"
	"sync/atomic"
	"testing"
	"time"

	"github.com/keep-network/keep-core/pkg/net"
)

func TestRetransmitExpectedNumberOfTimes(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 510*time.Millisecond)
	defer cancel()

	var retransmissionsCount uint64
	ScheduleRetransmissions(
		ctx,
		&testutils.MockLogger{},
		NewTimeTicker(ctx, 50*time.Millisecond),
		func() error {
			atomic.AddUint64(&retransmissionsCount, 1)
			return nil
		},
	)

	<-ctx.Done()

	if retransmissionsCount != 10 {
		t.Errorf("expected [10] retransmissions, has [%v]", retransmissionsCount)
	}
}

func TestHandlerReceiveUniqueMessages(t *testing.T) {
	var received []net.Message

	handler := WithRetransmissionSupport(func(message net.Message) {
		received = append(received, message)
	})

	handler(&mockNetworkMessage{senderID: "a", seqno: 1})
	handler(&mockNetworkMessage{senderID: "a", seqno: 2})
	handler(&mockNetworkMessage{senderID: "a", seqno: 4})
	handler(&mockNetworkMessage{senderID: "b", seqno: 1})
	handler(&mockNetworkMessage{senderID: "b", seqno: 2})

	if len(received) != 5 {
		t.Fatalf(
			"unexpected number of accepted messages\nactual:   [%v]\nexpected: [5]",
			len(received),
		)
	}
}

func TestHandlerReceiveRetransmissions(t *testing.T) {
	var received []net.Message

	handler := WithRetransmissionSupport(func(message net.Message) {
		received = append(received, message)
	})

	handler(&mockNetworkMessage{senderID: "a", seqno: 1})
	handler(&mockNetworkMessage{senderID: "a", seqno: 2})
	handler(&mockNetworkMessage{senderID: "a", seqno: 2})
	handler(&mockNetworkMessage{senderID: "a", seqno: 1})
	handler(&mockNetworkMessage{senderID: "b", seqno: 2})
	handler(&mockNetworkMessage{senderID: "b", seqno: 1})
	handler(&mockNetworkMessage{senderID: "b", seqno: 1})

	if len(received) != 4 {
		t.Fatalf(
			"unexpected number of accepted messages\nactual:   [%v]\nexpected: [4]",
			len(received),
		)
	}
}

type mockNetworkMessage struct {
	senderID string
	seqno    uint64
}

func (mnm *mockNetworkMessage) TransportSenderID() net.TransportIdentifier {
	return &mockTransportIdentifier{mnm.senderID}
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

func (mnm *mockNetworkMessage) Seqno() uint64 {
	return mnm.seqno
}

type mockTransportIdentifier struct {
	senderID string
}

func (mti *mockTransportIdentifier) String() string {
	return mti.senderID
}
