package dkg2

import (
	"fmt"
	"math/big"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/keep-network/keep-core/pkg/beacon/relay/gjkr"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/net/local"
)

func TestFullStateTransitions(t *testing.T) {
	threshold := 2
	groupSize := 5

	channels := make([]net.BroadcastChannel, groupSize)
	states := make([]keyGenerationState, groupSize)

	provider := local.Connect()

	dkg, err := gjkr.GenerateDKG()
	if err != nil {
		t.Fatal(err)
	}

	// Initialize one state and one channel per one member.
	// Each state gets a separate channel to have a different transport
	// identifier.
	for i := 0; i < groupSize; i++ {
		channel, err := provider.ChannelFor("transitions_test_" + string(i))
		if err != nil {
			t.Fatal(err)
		}

		member := gjkr.NewMember(
			gjkr.MemberID(i+1),
			make([]gjkr.MemberID, 0),
			threshold,
			dkg,
		)

		Init(channel)

		channels[i] = channel
		states[i] = &initializationState{channel, member}
	}

	// Perform all possible state transitions.
	for {
		nextStates, err := doStateTransition(states, channels)
		if err != nil {
			t.Fatal(err)
		}

		if nextStates == nil {
			break
		}

		states = nextStates
	}

	// Check whether all states are final and extract generated group
	// public keys.
	groupPublicKeys := make([]*big.Int, groupSize)
	for i, state := range states {
		finalState, ok := state.(*combiningState)
		if !ok {
			t.Fatalf("not a final state: %#v", state)
		}

		groupPublicKeys[i] = finalState.member.GroupPublicKey()
	}

	// Check whether all group public keys are the same.
	for i := 1; i < len(groupPublicKeys); i++ {
		if !reflect.DeepEqual(groupPublicKeys[i], groupPublicKeys[0]) {
			t.Fatalf(
				"unexpected group public key\nexpected: %v\nactual:   %v",
				groupPublicKeys[i],
				groupPublicKeys[0],
			)
		}
	}
}

func doStateTransition(
	states []keyGenerationState,
	channels []net.BroadcastChannel,
) ([]keyGenerationState, error) {
	var waitGroup sync.WaitGroup
	waitGroup.Add(expectedMessagesCount(states))

	// All messages in the protocol are broadcast so if any message arrive
	// during the current stage, it is going to be stored in `phaseMessages`
	// slice.
	var phaseMessagesMutex = &sync.Mutex{}
	var phaseMessages []net.Message

	for _, channel := range channels {
		if err := channel.Recv(net.HandleMessageFunc{
			Type: "test",
			Handler: func(msg net.Message) error {
				phaseMessagesMutex.Lock()
				phaseMessages = append(phaseMessages, msg)
				phaseMessagesMutex.Unlock()
				waitGroup.Done()
				return nil
			},
		}); err != nil {
			return nil, fmt.Errorf("message handler failed [%v]", err)
		}

		defer channel.UnregisterRecv("test")
	}

	// Once we have the message handler installed, we let all members to init
	// the phase and send their messages if they want to.
	for _, state := range states {
		fmt.Printf("[member:%v, state:%T] Executing\n", state.memberID(), state)

		if err := state.initiate(); err != nil {
			return nil, fmt.Errorf("initiate failed [%v]", err)
		}
	}

	// Now we need to wait until the expected amount of messages arrive.
	// We can't wait forever, so we set the timeout to 5 seconds.
	if waitWithTimeout(&waitGroup, 5*time.Second) {
		return nil, fmt.Errorf("timed out waiting for messages")
	}

	// All members sent their messages, now it's part to process received
	// messages and proceed to the next phase.
	nextStates := make([]keyGenerationState, len(states))
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

	// When there is no next phase to be exected, `nextState()` in
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

// expectedMessagesCount returns number of messages that are expected to be
// produced by all members of the group at the given stage of DKG protocol.
func expectedMessagesCount(states []keyGenerationState) int {
	statesCount := len(states)
	if statesCount == 0 {
		return 0
	}

	switch states[0].(type) {
	case *joinState:
		return statesCount
	case *ephemeralKeyPairGeneratingState:
		return statesCount
	case *committingState:
		return statesCount * 2 // shares + commitments
	case *commitmentsVerificationState:
		return statesCount
	case *pointsSharingState:
		return statesCount
	case *pointsValidationState:
		return statesCount
	default:
		return 0
	}
}

// waitWithTimeout waits for the waitgroup for the specified max timeout.
// Returns true if waiting timed out.
func waitWithTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}
