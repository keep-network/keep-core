package inactivity

import (
	"reflect"
	"testing"

	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/protocol/state"
)

func TestReceivedMessages(t *testing.T) {
	state := state.NewBaseAsyncState()

	message1 := &claimSignatureMessage{senderID: 1}
	message2 := &claimSignatureMessage{senderID: 2}
	message3 := &claimSignatureMessage{senderID: 1}
	message4 := &claimSignatureMessage{senderID: 3}
	message5 := &claimSignatureMessage{senderID: 3}

	state.ReceiveToHistory(newMockNetMessage(message1))
	state.ReceiveToHistory(newMockNetMessage(message2))
	state.ReceiveToHistory(newMockNetMessage(message3))
	state.ReceiveToHistory(newMockNetMessage(message4))
	state.ReceiveToHistory(newMockNetMessage(message5))

	expectedMessages := []*claimSignatureMessage{message1, message2, message4}
	actualType1Messages := receivedMessages[*claimSignatureMessage](state)
	if !reflect.DeepEqual(expectedMessages, actualType1Messages) {
		t.Errorf(
			"unexpected messages\n"+
				"expected: [%v]\n"+
				"actual:   [%v]",
			expectedMessages,
			actualType1Messages,
		)
	}
}

type mockNetMessage struct {
	payload interface{}
}

func newMockNetMessage(payload interface{}) *mockNetMessage {
	return &mockNetMessage{payload}
}

func (mnm *mockNetMessage) TransportSenderID() net.TransportIdentifier {
	panic("not implemented")
}

func (mnm *mockNetMessage) SenderPublicKey() []byte {
	panic("not implemented")
}

func (mnm *mockNetMessage) Payload() interface{} {
	return mnm.payload
}

func (mnm *mockNetMessage) Type() string {
	payload, ok := mnm.payload.(message)
	if !ok {
		panic("wrong payload type")
	}

	return payload.Type()
}

func (mnm *mockNetMessage) Seqno() uint64 {
	panic("not implemented")
}
