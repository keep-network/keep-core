package state

import (
	"context"
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
	"github.com/keep-network/keep-core/pkg/chain"
	chainLocal "github.com/keep-network/keep-core/pkg/chain/local"
	"github.com/keep-network/keep-core/pkg/net"
	netLocal "github.com/keep-network/keep-core/pkg/net/local"
)

var testLog map[uint64][]string
var blockCounter chain.BlockCounter

func TestExecute(t *testing.T) {
	testLog = make(map[uint64][]string)

	localChain := chainLocal.Connect(10, 5, big.NewInt(200))
	blockCounter, _ = localChain.BlockCounter()
	provider := netLocal.Connect()
	channel, err := provider.BroadcastChannelFor("transitions_test")
	if err != nil {
		t.Fatal(err)
	}

	go func(blockCounter chain.BlockCounter) {
		blockCounter.WaitForBlockHeight(1)
		ctx, cancel := context.WithCancel(context.Background())
		channel.Send(ctx, &TestMessage{"message_1"})
		cancel()

		blockCounter.WaitForBlockHeight(4)
		ctx, cancel = context.WithCancel(context.Background())
		channel.Send(ctx, &TestMessage{"message_2"})
		cancel()

		blockCounter.WaitForBlockHeight(7)
		ctx, cancel = context.WithCancel(context.Background())
		channel.Send(ctx, &TestMessage{"message_3"})
		cancel()
	}(blockCounter)

	channel.RegisterUnmarshaler(func() net.TaggedUnmarshaler {
		return &TestMessage{}
	})

	initialState := testState1{
		memberIndex: group.MemberIndex(1),
		channel:     channel,
	}

	stateMachine := NewMachine(channel, blockCounter, initialState)

	finalState, endBlockHeight, err := stateMachine.Execute(1)
	if err != nil {
		t.Errorf("unexpected error [%v]", err)
	}

	if _, ok := finalState.(*testState5); !ok {
		t.Errorf("state is not final [%v]", finalState)
	}

	if endBlockHeight != 8 {
		t.Errorf("unexpected end block [%v]", endBlockHeight)
	}

	expectedTestLog := map[uint64][]string{
		1: {
			"1-state.testState1-initiate",
			"1-state.testState1-receive-message_1",
		},
		3: {"1-state.testState2-initiate"},
		4: {"1-state.testState2-receive-message_2"},
		6: {
			"1-state.testState3-initiate",
			"1-state.testState4-initiate",
		},
		7: {
			"1-state.testState4-receive-message_3",
		},
		8: {
			"1-state.testState5-initiate",
		},
	}

	if !reflect.DeepEqual(expectedTestLog, testLog) {
		t.Errorf("\nexpected: %v\nactual:   %v\n", expectedTestLog, testLog)
	}
}

func addToTestLog(testState State, functionName string) {
	currentBlock, _ := blockCounter.CurrentBlock()
	testLog[currentBlock] = append(
		testLog[currentBlock],
		fmt.Sprintf(
			"%v-%v-%v",
			testState.MemberIndex(),
			reflect.TypeOf(testState),
			functionName,
		),
	)
}

type testState1 struct {
	memberIndex group.MemberIndex
	channel     net.BroadcastChannel
}

func (ts testState1) DelayBlocks() uint64  { return 0 }
func (ts testState1) ActiveBlocks() uint64 { return 2 }
func (ts testState1) Initiate(ctx context.Context) error {
	addToTestLog(ts, "initiate")
	return nil
}
func (ts testState1) Receive(msg net.Message) error {
	addToTestLog(
		ts,
		fmt.Sprintf("receive-%v", msg.Payload().(*TestMessage).content),
	)
	return nil
}
func (ts testState1) Next() State                    { return &testState2{ts} }
func (ts testState1) MemberIndex() group.MemberIndex { return ts.memberIndex }

type testState2 struct {
	testState1
}

func (ts testState2) DelayBlocks() uint64  { return 0 }
func (ts testState2) ActiveBlocks() uint64 { return 2 }
func (ts testState2) Initiate(ctx context.Context) error {
	addToTestLog(ts, "initiate")
	return nil
}
func (ts testState2) Receive(msg net.Message) error {
	addToTestLog(
		ts,
		fmt.Sprintf("receive-%v", msg.Payload().(*TestMessage).content),
	)
	return nil
}
func (ts testState2) Next() State                    { return &testState3{ts} }
func (ts testState2) MemberIndex() group.MemberIndex { return ts.memberIndex }

type testState3 struct {
	testState2
}

func (ts testState3) DelayBlocks() uint64  { return 1 }
func (ts testState3) ActiveBlocks() uint64 { return 0 }

func (ts testState3) Initiate(ctx context.Context) error {
	addToTestLog(ts, "initiate")
	return nil
}

func (ts testState3) Receive(msg net.Message) error {
	addToTestLog(
		ts,
		fmt.Sprintf("receive-%v", msg.Payload().(*TestMessage).content),
	)
	return nil
}

func (ts testState3) Next() State                    { return &testState4{ts} }
func (ts testState3) MemberIndex() group.MemberIndex { return ts.memberIndex }

type testState4 struct {
	testState3
}

func (ts testState4) DelayBlocks() uint64  { return 0 }
func (ts testState4) ActiveBlocks() uint64 { return 2 }
func (ts testState4) Initiate(ctx context.Context) error {
	addToTestLog(ts, "initiate")
	return nil
}
func (ts testState4) Receive(msg net.Message) error {
	addToTestLog(
		ts,
		fmt.Sprintf("receive-%v", msg.Payload().(*TestMessage).content),
	)
	return nil
}
func (ts testState4) Next() State                    { return &testState5{ts} }
func (ts testState4) MemberIndex() group.MemberIndex { return ts.memberIndex }

type testState5 struct {
	testState4
}

func (ts testState5) DelayBlocks() uint64  { return 0 }
func (ts testState5) ActiveBlocks() uint64 { return 0 }
func (ts testState5) Initiate(ctx context.Context) error {
	addToTestLog(ts, "initiate")
	return nil
}
func (ts testState5) Receive(msg net.Message) error {
	addToTestLog(
		ts,
		fmt.Sprintf("receive-%v", msg.Payload().(*TestMessage).content),
	)
	return nil
}
func (ts testState5) Next() State                    { return nil }
func (ts testState5) MemberIndex() group.MemberIndex { return ts.memberIndex }

type TestMessage struct {
	content string
}

func (tm *TestMessage) Marshal() ([]byte, error) {
	return []byte(tm.content), nil
}

func (tm *TestMessage) Unmarshal(bytes []byte) error {
	tm.content = string(bytes)
	return nil
}

func (tm *TestMessage) Type() string {
	return "test_message"
}
