package faststate

import (
	"context"
	"encoding/json"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/keep-network/keep-core/pkg/internal/testutils"
	"github.com/keep-network/keep-core/pkg/net"
	netLocal "github.com/keep-network/keep-core/pkg/net/local"
	"github.com/keep-network/keep-core/pkg/protocol/group"
)

const (
	partyCount        = 3
	round2MessageType = "test/round_2_message"
	round3MessageType = "test/round_3_message"
)

func TestExecute(t *testing.T) {
	provider := netLocal.Connect()
	channel, err := provider.BroadcastChannelFor("test")
	if err != nil {
		t.Fatal(err)
	}

	channel.SetUnmarshaler(func() net.TaggedUnmarshaler {
		return &round2Message{}
	})
	channel.SetUnmarshaler(func() net.TaggedUnmarshaler {
		return &round3Message{}
	})

	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	var logger = &testutils.MockLogger{}

	initialState1 := &testState1{
		BaseState:   NewBaseState(),
		memberIndex: group.MemberIndex(1),
		channel:     channel,
	}
	initialState2 := &testState1{
		BaseState:   NewBaseState(),
		memberIndex: group.MemberIndex(2),
		channel:     channel,
	}
	initialState3 := &testState1{
		BaseState:   NewBaseState(),
		memberIndex: group.MemberIndex(3),
		channel:     channel,
	}

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		NewMachine(logger, ctx, channel, initialState1).Execute()
		wg.Done()
	}()
	go func() {
		NewMachine(logger, ctx, channel, initialState2).Execute()
		wg.Done()
	}()
	go func() {
		NewMachine(logger, ctx, channel, initialState3).Execute()
		wg.Done()
	}()

	wg.Wait()

	expectedExecutionLog := "1-initiate 1-done 2-initiate 2-done 3-initiate 3-done"
	testutils.AssertStringsEqual(
		t,
		"member 1 execution log",
		expectedExecutionLog,
		strings.Join(initialState1.testLog, " "),
	)
	testutils.AssertStringsEqual(
		t,
		"member 2 execution log",
		expectedExecutionLog,
		strings.Join(initialState2.testLog, " "),
	)
	testutils.AssertStringsEqual(
		t,
		"member 3 execution log",
		expectedExecutionLog,
		strings.Join(initialState3.testLog, " "),
	)
}

//
// testState1 can transition to the next state immediately;
// it is not receiving or sending any messages.
//
type testState1 struct {
	*BaseState

	memberIndex group.MemberIndex
	channel     net.BroadcastChannel
	testLog     []string
}

func (ts *testState1) addToTestLog(log string) {
	ts.testLog = append(ts.testLog, log)
}

func (ts *testState1) CanTransition() bool {
	return true
}
func (ts *testState1) Initiate(ctx context.Context) error {
	ts.addToTestLog("1-initiate")
	return nil
}
func (ts *testState1) Receive(msg net.Message) error {
	ts.receive(msg)
	return nil
}
func (ts *testState1) Next() (State, error) {
	ts.addToTestLog("1-done")
	return &testState2{testState1: ts}, nil
}
func (ts *testState1) MemberIndex() group.MemberIndex { return ts.memberIndex }

//
// testState2 waits for `partyCount` messages to be received;
// it is sending one message immediately upon the initiation.
//
type testState2 struct {
	*testState1
}

func (ts *testState2) CanTransition() bool {
	return len(ts.GetAllReceivedMessages(round2MessageType)) == partyCount
}
func (ts *testState2) Initiate(ctx context.Context) error {
	ts.addToTestLog("2-initiate")
	ts.channel.Send(ctx, newRound2Message(strconv.Itoa(int(ts.memberIndex))))
	return nil
}
func (ts *testState2) Receive(msg net.Message) error {
	ts.receive(msg)
	return nil
}
func (ts *testState2) Next() (State, error) {
	ts.addToTestLog("2-done")
	return &testState3{testState2: ts}, nil
}
func (ts *testState2) MemberIndex() group.MemberIndex { return ts.memberIndex }

//
// testState3 waits for `partyCount` messages to be received.
// It is sending one message after some random delay upon the initiation.
//
type testState3 struct {
	*testState2
}

func (ts *testState3) CanTransition() bool {
	return len(ts.GetAllReceivedMessages(round3MessageType)) == partyCount
}
func (ts *testState3) Initiate(ctx context.Context) error {
	ts.addToTestLog("3-initiate")
	rand.Seed(time.Now().UnixNano())
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	ts.channel.Send(ctx, newRound3Message(strconv.Itoa(int(ts.memberIndex))))
	return nil
}
func (ts *testState3) Receive(msg net.Message) error {
	ts.receive(msg)
	return nil
}
func (ts *testState3) Next() (State, error) {
	ts.addToTestLog("3-done")
	return nil, nil
}
func (ts *testState3) MemberIndex() group.MemberIndex { return ts.memberIndex }

//
// testState2 message
//
type round2Message struct {
	Type_   string `json:"type"`
	Payload string `json:"payload"`
}

func (r2m *round2Message) Type() string {
	return "test/round_2_message"
}
func (r2m *round2Message) Marshal() ([]byte, error) {
	return json.Marshal(r2m)
}
func (r2m *round2Message) Unmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, r2m)
}
func newRound2Message(payload string) *round2Message {
	return &round2Message{
		Type_:   round2MessageType,
		Payload: payload,
	}
}

//
// testState3 message
//
type round3Message struct {
	Type_   string `json:"type"`
	Payload string `json:"payload"`
}

func (r3m *round3Message) Type() string {
	return "test/round_3_message"
}
func (r3m *round3Message) Marshal() ([]byte, error) {
	return json.Marshal(r3m)
}
func (r3m *round3Message) Unmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, r3m)
}
func newRound3Message(payload string) *round3Message {
	return &round3Message{
		Type_:   round3MessageType,
		Payload: payload,
	}
}
