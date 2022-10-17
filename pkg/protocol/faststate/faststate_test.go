package faststate

import (
	"testing"

	"github.com/keep-network/keep-core/pkg/internal/testutils"
	"github.com/keep-network/keep-core/pkg/net"
)

func TestBaseStateReceive(t *testing.T) {
	const type1 = "faststate/test_type1"
	const type2 = "faststate/test_type2"

	state := NewBaseState()

	messages := state.GetAllReceivedMessages(type1)
	if len(messages) != 0 {
		t.Fatalf("expected no messages of the first type yet")
	}
	messages = state.GetAllReceivedMessages(type2)
	if len(messages) != 0 {
		t.Fatalf("expected no messages of the second type yet")
	}

	state.receive(&mockMessage{_type: type1, payload: "a"})
	state.receive(&mockMessage{_type: type1, payload: "b"})
	state.receive(&mockMessage{_type: type2, payload: "a"})
	state.receive(&mockMessage{_type: type1, payload: "c"})

	messages = state.GetAllReceivedMessages(type1)
	testutils.AssertIntsEqual(
		t,
		"number of messages of the first type",
		3,
		len(messages),
	)
	messages = state.GetAllReceivedMessages(type2)
	testutils.AssertIntsEqual(
		t,
		"number of messages of the second type",
		1,
		len(messages),
	)
}

type mockMessage struct {
	payload string
	_type   string
}

func (mm *mockMessage) TransportSenderID() net.TransportIdentifier {
	panic("unsupported in mock")
}
func (mm *mockMessage) SenderPublicKey() []byte {
	panic("unsupported in mock")
}
func (mm *mockMessage) Payload() interface{} {
	return mm.payload
}
func (mm *mockMessage) Type() string {
	return mm._type
}
func (mm *mockMessage) Seqno() uint64 {
	panic("unsupported in mock")
}
