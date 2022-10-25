package signing

import (
	"context"
	"fmt"
	"github.com/keep-network/keep-core/pkg/internal/testutils"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/protocol/faststate"
	"reflect"
	"testing"
	"time"
)

func TestStateAction(t *testing.T) {
	actionOne := &stateAction{}

	testutils.AssertBoolsEqual(
		t,
		"action one initial state",
		false,
		actionOne.isDone(),
	)

	testutils.AssertErrorsSame(
		t,
		nil,
		actionOne.error(),
	)

	funcOneCtx, funcOneCancelCtx := context.WithTimeout(
		context.Background(),
		1 * time.Second,
	)
	funcOne := func() error {
		defer funcOneCancelCtx()
		return nil
	}

	actionOne.run(funcOne)

	<-funcOneCtx.Done()

	// Context should be cancelled as function one was actually executed.
	testutils.AssertErrorsSame(
		t,
		funcOneCtx.Err(),
		context.Canceled,
	)

	testutils.AssertBoolsEqual(
		t,
		"action one state after running function one",
		true,
		actionOne.isDone(),
	)

	testutils.AssertErrorsSame(
		t,
		nil,
		actionOne.error(),
	)

	funcTwoCtx, funcTwoCancelCtx := context.WithTimeout(
		context.Background(),
		1 * time.Second,
	)
	funcTwo := func() error {
		defer funcTwoCancelCtx()
		return nil
	}

	// This should be a no-op.
	actionOne.run(funcTwo)

	<-funcTwoCtx.Done()

	// Context should be deadline exceeded as function two was not executed.
	testutils.AssertErrorsSame(
		t,
		funcTwoCtx.Err(),
		context.DeadlineExceeded,
	)

	actionTwo := &stateAction{}

	testutils.AssertBoolsEqual(
		t,
		"action two initial state",
		false,
		actionTwo.isDone(),
	)

	testutils.AssertErrorsSame(
		t,
		nil,
		actionTwo.error(),
	)

	funcThreeErr := fmt.Errorf("func three error")
	funcThreeCtx, funcThreeCancelCtx := context.WithTimeout(
		context.Background(),
		1 * time.Second,
	)
	funcThree := func() error {
		defer funcThreeCancelCtx()
		return funcThreeErr
	}

	actionTwo.run(funcThree)

	<-funcThreeCtx.Done()

	// Context should be cancelled as function one was actually executed.
	testutils.AssertErrorsSame(
		t,
		funcThreeCtx.Err(),
		context.Canceled,
	)

	testutils.AssertBoolsEqual(
		t,
		"action two state after running function three",
		true,
		actionTwo.isDone(),
	)

	testutils.AssertErrorsSame(
		t,
		funcThreeErr,
		actionTwo.error(),
	)
}

func TestReceivedMessages(t *testing.T) {
	state := faststate.NewBaseState()

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
	actualType1Messages :=  receivedMessages[*ephemeralPublicKeyMessage](state)
	if !reflect.DeepEqual(expectedType1Messages, actualType1Messages) {
		t.Errorf(
			"unexpected messages\n" +
				"expected: [%v]\n" +
				"actual:   [%v]",
			expectedType1Messages,
			actualType1Messages,
		)
	}

	expectedType2Messages := []*tssRoundOneMessage{message2, message3}
	actualType2Messages :=  receivedMessages[*tssRoundOneMessage](state)
	if !reflect.DeepEqual(expectedType2Messages, actualType2Messages) {
		t.Errorf(
			"unexpected messages\n" +
				"expected: [%v]\n" +
				"actual:   [%v]",
			expectedType2Messages,
			actualType2Messages,
		)
	}

	expectedType3Messages := make([]*tssRoundTwoMessage, 0)
	actualType3Messages :=  receivedMessages[*tssRoundTwoMessage](state)
	if !reflect.DeepEqual(expectedType3Messages, actualType3Messages) {
		t.Errorf(
			"unexpected messages\n" +
				"expected: [%v]\n" +
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
