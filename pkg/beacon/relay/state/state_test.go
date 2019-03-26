package state

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/keep-network/keep-core/pkg/beacon/relay/member"
	"github.com/keep-network/keep-core/pkg/chain"
	chainLocal "github.com/keep-network/keep-core/pkg/chain/local"
	"github.com/keep-network/keep-core/pkg/net"
	netLocal "github.com/keep-network/keep-core/pkg/net/local"
)

var testLog map[int][]string
var blockCounter chain.BlockCounter

func TestExecute(t *testing.T) {
	testLog = make(map[int][]string)

	localChain := chainLocal.Connect(10, 5, big.NewInt(200))
	blockCounter, _ = localChain.BlockCounter()
	provider := netLocal.Connect()
	channel, err := provider.ChannelFor("transitions_test")
	if err != nil {
		t.Fatal(err)
	}

	go func(blockCounter chain.BlockCounter) {
		blockCounter.WaitForBlockHeight(1)
		channel.Send(&TestMessage{"message_1"})

		blockCounter.WaitForBlockHeight(5)
		channel.Send(&TestMessage{"message_2"})

		blockCounter.WaitForBlockHeight(9)
		channel.Send(&TestMessage{"message_3"})
	}(blockCounter)

	initChannel := func(channel net.BroadcastChannel) {
		currentBlock, _ := blockCounter.CurrentBlock()
		testLog[currentBlock] = append(
			testLog[currentBlock],
			"init_channel-"+channel.Name(),
		)

		channel.RegisterUnmarshaler(func() net.TaggedUnmarshaler {
			return &TestMessage{}
		})
	}

	initialState := testState1{
		memberIndex: member.Index(1),
		channel:     channel,
	}

	stateMachine := NewMachine(channel, blockCounter, initialState)

	finalState, err := stateMachine.Execute(initChannel)
	if err != nil {
		t.Errorf("unexpected error [%v]", err)
	}

	if _, ok := finalState.(*testState4); !ok {
		t.Errorf("state is not final [%v]", finalState)
	}

	expectedTestLog := map[int][]string{
		0: []string{"init_channel-transitions_test"},
		1: []string{
			"1-state.testState1-initiate",
			"1-state.testState1-receive-message_1",
		},
		4: []string{"1-state.testState2-initiate"},
		5: []string{"1-state.testState2-receive-message_2"},
		7: []string{"1-state.testState3-initiate"},
		10: []string{
			"1-state.testState4-initiate",
			"1-state.testState4-receive-message_3",
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
	memberIndex member.Index
	channel     net.BroadcastChannel
}

func (ts testState1) ActiveBlocks() int { return 2 }
func (ts testState1) Initiate() error {
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
func (ts testState1) NextState() State          { return &testState2{ts} }
func (ts testState1) MemberIndex() member.Index { return ts.memberIndex }

type testState2 struct {
	testState1
}

func (ts testState2) ActiveBlocks() int { return 2 }
func (ts testState2) Initiate() error {
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
func (ts testState2) NextState() State          { return &testState3{ts} }
func (ts testState2) MemberIndex() member.Index { return ts.memberIndex }

type testState3 struct {
	testState2
}

func (ts testState3) ActiveBlocks() int { return 2 }

func (ts testState3) Initiate() error {
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

func (ts testState3) NextState() State          { return &testState4{ts} }
func (ts testState3) MemberIndex() member.Index { return ts.memberIndex }

type testState4 struct {
	testState3
}

func (ts testState4) ActiveBlocks() int { return 2 }
func (ts testState4) Initiate() error {
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
func (ts testState4) NextState() State          { return nil }
func (ts testState4) MemberIndex() member.Index { return ts.memberIndex }

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
