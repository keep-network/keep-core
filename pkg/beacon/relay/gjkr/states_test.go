package gjkr

import (
	"fmt"
	"math/big"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/net/local"
)

func TestFullStateTransitions(t *testing.T) {
	dishonestThreshold := 2
	groupSize := 5

	seed := big.NewInt(18293712839)

	channels := make([]net.BroadcastChannel, groupSize)
	states := make([]keyGenerationState, groupSize)

	provider := local.Connect()

	// Initialize one state and one channel per one member.
	// Each state gets a separate channel to have a different transport
	// identifier.
	for i := 0; i < groupSize; i++ {
		channel, err := provider.ChannelFor("transitions_test_" + string(i))
		if err != nil {
			t.Fatal(err)
		}

		member, err := NewMember(
			group.MemberIndex(i+1),
			groupSize,
			dishonestThreshold,
			seed,
		)
		if err != nil {
			t.Fatal(err)
		}

		InitializeChannel(channel)

		channels[i] = channel
		states[i] = &joinState{channel, member}
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

	// Check whether all states are final and produced the same result.
	results := make([]*Result, groupSize)
	for i, state := range states {
		finalState, ok := state.(*finalizationState)
		if !ok {
			t.Fatalf("not a final state: %#v", state)
		}

		results[i] = finalState.member.Result()
	}

	// Check whether all group public keys are the same, and they have no
	// disqualified or inactive  members.
	for _, result := range results {
		if len(result.Group.InactiveMemberIDs()) != 0 {
			t.Errorf("expected no IA members\n[%v]", result)
		}
		if len(result.Group.DisqualifiedMemberIDs()) != 0 {
			t.Errorf("expected no DQ members\n[%v]", result)
		}

		if result.GroupPublicKey.String() != results[0].GroupPublicKey.String() {
			t.Fatalf(
				"Unexpected group public key\nExpected: [%v]\nActual:   [%v]\n",
				result.GroupPublicKey,
				results[0].GroupPublicKey,
			)
		}
		if !reflect.DeepEqual(result.Group, results[0].Group) {
			t.Fatalf(
				"Unexpected group information\nExpected: [%v]\nActual:   [%v]\n",
				result.Group,
				results[0].Group,
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
		fmt.Printf("[member:%v, state:%T] Executing\n", state.MemberIndex(), state)

		if err := state.Initiate(); err != nil {
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
			if err := state.Receive(message); err != nil {
				return nil, fmt.Errorf("receive failed [%v]", err)
			}
		}

		next := state.Next()

		if next != nil {
			fmt.Printf(
				"[member:%v, state:%T] Successfully transitioned to the next state\n",
				state.MemberIndex(),
				state,
			)
		}

		nextStates[i] = next
	}

	// When there is no next phase to be exected, `Next()` in
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
	case *ephemeralKeyPairGenerationState:
		return statesCount
	case *commitmentState:
		return statesCount * 2 // shares + commitments
	case *commitmentsVerificationState:
		return statesCount
	case *pointsShareState:
		return statesCount
	case *pointsValidationState:
		return statesCount
	case *keyRevealState:
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
