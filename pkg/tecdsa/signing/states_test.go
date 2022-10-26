package signing

import (
	"reflect"
	"testing"

	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/protocol/state"
)

func TestReceivedMessages(t *testing.T) {
	state := state.NewBaseAsyncState()

	message1 := &ephemeralPublicKeyMessage{senderID: 1}
	message2 := &tssRoundOneMessage{senderID: 1}
	message3 := &tssRoundOneMessage{senderID: 2}
	message4 := &ephemeralPublicKeyMessage{senderID: 3}
	message5 := &ephemeralPublicKeyMessage{senderID: 2}
	message6 := &ephemeralPublicKeyMessage{senderID: 3}
	message7 := &tssRoundOneMessage{senderID: 2}

	state.ReceiveToHistory(newMockNetMessage(message1))
	state.ReceiveToHistory(newMockNetMessage(message2))
	state.ReceiveToHistory(newMockNetMessage(message3))
	state.ReceiveToHistory(newMockNetMessage(message4))
	state.ReceiveToHistory(newMockNetMessage(message5))
	state.ReceiveToHistory(newMockNetMessage(message6))
	state.ReceiveToHistory(newMockNetMessage(message7))

	expectedType1Messages := []*ephemeralPublicKeyMessage{message1, message4, message5}
	actualType1Messages := receivedMessages[*ephemeralPublicKeyMessage](state)
	if !reflect.DeepEqual(expectedType1Messages, actualType1Messages) {
		t.Errorf(
			"unexpected messages\n"+
				"expected: [%v]\n"+
				"actual:   [%v]",
			expectedType1Messages,
			actualType1Messages,
		)
	}

	expectedType2Messages := []*tssRoundOneMessage{message2, message3}
	actualType2Messages := receivedMessages[*tssRoundOneMessage](state)
	if !reflect.DeepEqual(expectedType2Messages, actualType2Messages) {
		t.Errorf(
			"unexpected messages\n"+
				"expected: [%v]\n"+
				"actual:   [%v]",
			expectedType2Messages,
			actualType2Messages,
		)
	}

	expectedType3Messages := make([]*tssRoundTwoMessage, 0)
	actualType3Messages := receivedMessages[*tssRoundTwoMessage](state)
	if !reflect.DeepEqual(expectedType3Messages, actualType3Messages) {
		t.Errorf(
			"unexpected messages\n"+
				"expected: [%v]\n"+
				"actual:   [%v]",
			expectedType3Messages,
			actualType3Messages,
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
