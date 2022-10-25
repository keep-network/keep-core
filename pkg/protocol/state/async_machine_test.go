package state

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
	netlocal "github.com/keep-network/keep-core/pkg/net/local"
	"github.com/keep-network/keep-core/pkg/protocol/group"
)

const (
	partyCount        = 3
	round2MessageType = "test/round_2_message"
	round3MessageType = "test/round_3_message"
)

func TestAsyncExecute(t *testing.T) {
	provider := netlocal.Connect()
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

	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	var logger = &testutils.MockLogger{}

	initialState1 := &testAsyncState1{
		BaseAsyncState: NewBaseAsyncState(),
		memberIndex:    group.MemberIndex(1),
		channel:        channel,
	}
	initialState2 := &testAsyncState1{
		BaseAsyncState: NewBaseAsyncState(),
		memberIndex:    group.MemberIndex(2),
		channel:        channel,
	}
	initialState3 := &testAsyncState1{
		BaseAsyncState: NewBaseAsyncState(),
		memberIndex:    group.MemberIndex(3),
		channel:        channel,
	}

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		NewAsyncMachine(logger, ctx, channel, initialState1).Execute()
		wg.Done()
	}()
	go func() {
		NewAsyncMachine(logger, ctx, channel, initialState2).Execute()
		wg.Done()
	}()
	go func() {
		NewAsyncMachine(logger, ctx, channel, initialState3).Execute()
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

func TestAsyncExecute_ContextCancelled(t *testing.T) {
	provider := netlocal.Connect()
	channel, err := provider.BroadcastChannelFor("test")
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	var logger = &testutils.MockLogger{}

	initialState := &testAsyncState1{
		BaseAsyncState: NewBaseAsyncState(),
		memberIndex:    group.MemberIndex(1),
		channel:        channel,
	}

	go func() {
		time.Sleep(1 * time.Second)
		cancelCtx()
	}()

	_, err = NewAsyncMachine(logger, ctx, channel, initialState).Execute()

	testutils.AssertErrorsSame(t, context.Canceled, err)
}

//
// testAsyncState1 can transition to the next state immediately;
// it is not receiving or sending any messages.
//
type testAsyncState1 struct {
	*BaseAsyncState

	memberIndex group.MemberIndex
	channel     net.BroadcastChannel
	testLog     []string
}

func (tas *testAsyncState1) addToTestLog(log string) {
	tas.testLog = append(tas.testLog, log)
}

func (tas *testAsyncState1) CanTransition() bool {
	return true
}
func (tas *testAsyncState1) Initiate(ctx context.Context) error {
	tas.addToTestLog("1-initiate")
	return nil
}
func (tas *testAsyncState1) Receive(msg net.Message) error {
	tas.ReceiveToHistory(msg)
	return nil
}
func (tas *testAsyncState1) Next() (AsyncState, error) {
	tas.addToTestLog("1-done")
	return &testAsyncState2{testAsyncState1: tas}, nil
}
func (tas *testAsyncState1) MemberIndex() group.MemberIndex { return tas.memberIndex }

//
// testAsyncState2 waits for `partyCount` messages to be received;
// it is sending one message immediately upon the initiation.
//
type testAsyncState2 struct {
	*testAsyncState1
}

func (tas *testAsyncState2) CanTransition() bool {
	return len(tas.GetAllReceivedMessages(round2MessageType)) == partyCount
}
func (tas *testAsyncState2) Initiate(ctx context.Context) error {
	tas.addToTestLog("2-initiate")
	tas.channel.Send(ctx, newRound2Message(strconv.Itoa(int(tas.memberIndex))))
	return nil
}
func (tas *testAsyncState2) Receive(msg net.Message) error {
	tas.ReceiveToHistory(msg)
	return nil
}
func (tas *testAsyncState2) Next() (AsyncState, error) {
	tas.addToTestLog("2-done")
	return &testAsyncState3{testAsyncState2: tas}, nil
}
func (tas *testAsyncState2) MemberIndex() group.MemberIndex { return tas.memberIndex }

//
// testAsyncState3 waits for `partyCount` messages to be received.
// It is sending one message after some random delay upon the initiation.
//
type testAsyncState3 struct {
	*testAsyncState2
}

func (tas *testAsyncState3) CanTransition() bool {
	return len(tas.GetAllReceivedMessages(round3MessageType)) == partyCount
}
func (tas *testAsyncState3) Initiate(ctx context.Context) error {
	tas.addToTestLog("3-initiate")
	rand.Seed(time.Now().UnixNano())
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	tas.channel.Send(ctx, newRound3Message(strconv.Itoa(int(tas.memberIndex))))
	return nil
}
func (tas *testAsyncState3) Receive(msg net.Message) error {
	tas.ReceiveToHistory(msg)
	return nil
}
func (tas *testAsyncState3) Next() (AsyncState, error) {
	tas.addToTestLog("3-done")
	return nil, nil
}
func (tas *testAsyncState3) MemberIndex() group.MemberIndex { return tas.memberIndex }

//
// testAsyncState2 message
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
// testAsyncState3 message
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
