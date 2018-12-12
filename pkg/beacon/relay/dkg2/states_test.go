package dkg2

import (
	"fmt"
	"testing"
	"time"

	"github.com/keep-network/keep-core/pkg/beacon/relay/gjkr"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/net/local"
)

func TestFullStateTransitions(t *testing.T) {
	threshold := 1
	groupSize := 2

	channels := make([]net.BroadcastChannel, groupSize)
	states := make([]keyGenerationState, groupSize)

	provider := local.Connect()

	for i := 0; i < groupSize; i++ {
		channel, err := provider.ChannelFor("transitions_test_" + string(i))
		if err != nil {
			t.Fatal(err)
		}

		Init(channel)

		channels[i] = channel

		member, err := gjkr.NewMember(
			gjkr.MemberID(i),
			make([]gjkr.MemberID, 0),
			threshold,
		)
		if err != nil {
			t.Fatal(err)
		}

		states[i] = &initializationState{channel, member}
	}

	for states != nil {
		nextStates, err := doStateTransition(states, channels)
		if err != nil {
			t.Fatal(err)
		}

		states = nextStates
	}
}

func doStateTransition(
	states []keyGenerationState,
	channels []net.BroadcastChannel,
) ([]keyGenerationState, error) {
	if len(states) == 0 {
		return nil, fmt.Errorf("at least one state required")
	}

	// phase duration time, assuming 1 block = 2 sec
	duration := 2 * time.Second * time.Duration(states[0].activeBlocks())

	nextStates := make([]keyGenerationState, len(states))

	var phaseMessages []net.Message

	for _, channel := range channels {
		if err := channel.Recv(net.HandleMessageFunc{
			Type: "test",
			Handler: func(msg net.Message) error {
				phaseMessages = append(phaseMessages, msg)
				return nil
			},
		}); err != nil {
			return nil, fmt.Errorf("message handler failed [%v]", err)
		}

		defer channel.UnregisterRecv("test")
	}

	// let all members to init and send their messages
	for _, state := range states {
		fmt.Printf(
			"[member:%v, state:%T] Executing\n",
			state.memberID(),
			state,
		)

		if err := state.initiate(); err != nil {
			return nil, fmt.Errorf("initiate failed %v", err)
		}
	}

	time.Sleep(duration)
	fmt.Printf("Exchanged %v messages in this phase\n", len(phaseMessages))

	// all members sent their messages, now it's part to receive them and
	// proceed to the next state
	for i, state := range states {
		for _, message := range phaseMessages {
			if err := state.receive(message); err != nil {
				return nil, fmt.Errorf("receive failed [%v]", err)
			}
		}

		next, err := state.nextState()
		if err != nil {
			return nil, fmt.Errorf("nextState failed [%v]", err)
		}

		fmt.Printf(
			"[member:%v, state:%T] Successfully transitioned to the next state\n",
			state.memberID(),
			state,
		)

		nextStates[i] = next
	}

	// When there is no next state to be exected, `nextState()` in
	// `keyGenerationState` returns `nil`. To say that all states have
	// executed completely, all `states` must be `nil`.
	allCompleted := true
	for _, state := range nextStates {
		if state != nil {
			allCompleted = false
			break
		}
	}

	if allCompleted {
		return nil, nil
	}

	return nextStates, nil
}
