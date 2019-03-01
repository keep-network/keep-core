package local

import (
	"context"
	"fmt"
	"math/big"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"

	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
)

func TestSubmitTicketAndGetSelectedParticipants(t *testing.T) {
	groupSize := 4

	generateTicket := func(index int64) *relaychain.Ticket {
		return &relaychain.Ticket{
			Value: big.NewInt(10 * index),
			Proof: &relaychain.TicketProof{
				StakerValue:        big.NewInt(100 * index),
				VirtualStakerIndex: big.NewInt(index),
			},
		}
	}

	ticket1 := generateTicket(1)
	ticket2 := generateTicket(2)
	ticket3 := generateTicket(3)
	ticket4 := generateTicket(4)
	ticket5 := generateTicket(5)
	ticket6 := generateTicket(6)

	var tests = map[string]struct {
		submitTickets           func(chain relaychain.Interface)
		expectedSelectedTickets []*relaychain.Ticket
	}{
		"number of tickets is less than group size": {
			submitTickets: func(chain relaychain.Interface) {
				chain.SubmitTicket(ticket3)
				chain.SubmitTicket(ticket1)
				chain.SubmitTicket(ticket2)
			},
			expectedSelectedTickets: []*relaychain.Ticket{
				ticket1, ticket2, ticket3,
			},
		},
		"number of tickets is same as group size": {
			submitTickets: func(chain relaychain.Interface) {
				chain.SubmitTicket(ticket3)
				chain.SubmitTicket(ticket1)
				chain.SubmitTicket(ticket4)
				chain.SubmitTicket(ticket2)
			},
			expectedSelectedTickets: []*relaychain.Ticket{
				ticket1, ticket2, ticket3, ticket4,
			},
		},
		"number of tickets is greater than group size": {
			submitTickets: func(chain relaychain.Interface) {
				chain.SubmitTicket(ticket3)
				chain.SubmitTicket(ticket1)
				chain.SubmitTicket(ticket4)
				chain.SubmitTicket(ticket6)
				chain.SubmitTicket(ticket5)
				chain.SubmitTicket(ticket2)
			},
			expectedSelectedTickets: []*relaychain.Ticket{
				ticket1, ticket2, ticket3, ticket4,
			},
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			c := Connect(groupSize, 4, big.NewInt(200))
			chain := c.ThresholdRelay()

			test.submitTickets(chain)

			actualSelectedParticipants, err := chain.GetSelectedParticipants()
			if err != nil {
				t.Fatal(err)
			}

			expectedSelectedParticipants := make(
				[]relaychain.StakerAddress,
				len(test.expectedSelectedTickets),
			)
			for i, ticket := range test.expectedSelectedTickets {
				expectedSelectedParticipants[i] = ticket.Proof.StakerValue.Bytes()
			}

			if !reflect.DeepEqual(expectedSelectedParticipants, actualSelectedParticipants) {
				t.Fatalf(
					"\nexpected: %v\nactual:   %v\n",
					expectedSelectedParticipants,
					actualSelectedParticipants,
				)
			}
		})
	}
}

func TestLocalSubmitRelayEntry(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	chainHandle := Connect(10, 4, big.NewInt(200)).ThresholdRelay()
	relayEntryPromise := chainHandle.SubmitRelayEntry(
		&event.Entry{
			RequestID:   big.NewInt(int64(19)),
			GroupPubKey: []byte("1"),
		},
	)

	done := make(chan *event.Entry)
	relayEntryPromise.OnSuccess(func(entry *event.Entry) {
		done <- entry
	}).OnFailure(func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	})

	select {
	case entry := <-done:
		expected := int64(19)
		if entry.RequestID.Int64() != expected {
			t.Fatalf(
				"expected [%v], got [%v]",
				expected,
				entry.RequestID.Int64(),
			)
		}
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}

}

func TestLocalBlockHeightWaiter(t *testing.T) {
	var tests = map[string]struct {
		blockHeight      int
		initialDelay     time.Duration
		expectedWaitTime time.Duration
	}{
		"does not wait for negative block height": {
			blockHeight:      -1,
			expectedWaitTime: 0,
		},
		"returns immediately for genesis block": {
			blockHeight:      0,
			expectedWaitTime: 0,
		},
		"returns immediately for block height already reached": {
			blockHeight:      2,
			initialDelay:     3 * blockTime,
			expectedWaitTime: 0,
		},
		"waits for block height not yet reached": {
			blockHeight:      5,
			initialDelay:     2 * blockTime,
			expectedWaitTime: 3 * blockTime,
		},
	}

	for testName, test := range tests {
		test := test
		t.Run(testName, func(t *testing.T) {
			t.Parallel()
			c := Connect(10, 4, big.NewInt(100))

			blockCounter, err := c.BlockCounter()
			if err != nil {
				t.Fatalf("failed to set up block counter: [%v]", err)
			}

			time.Sleep(test.initialDelay)

			start := time.Now().UTC()
			blockCounter.WaitForBlockHeight(test.blockHeight)
			end := time.Now().UTC()

			elapsed := end.Sub(start)

			// Block waiter should wait for test.expectedWaitTime minus some
			// margin at minimum; the margin is needed because clock is not
			// always that precise. Setting it to 5ms for this test.
			minMargin := time.Duration(5) * time.Millisecond
			if elapsed < (test.expectedWaitTime - minMargin) {
				t.Errorf(
					"waited less than expected; expected [%v] at min, waited [%v]",
					test.expectedWaitTime,
					elapsed,
				)
			}

			// Block waiter should wait for test.expectedWaitTime plus some
			// margin at maximum; the margin is the time needed for the return
			// instructions to execute, setting it to 25ms for this test.
			maxMargin := time.Duration(25) * time.Millisecond
			if elapsed > (test.expectedWaitTime + maxMargin) {
				t.Errorf(
					"waited longer than expected; expected %v at max, waited %v",
					test.expectedWaitTime,
					elapsed,
				)
			}
		})
	}
}

func TestLocalBlockWaiter(t *testing.T) {
	var tests = map[string]struct {
		blockWait        int
		expectedWaitTime time.Duration
	}{
		"does wait for a block": {
			blockWait:        1,
			expectedWaitTime: blockTime,
		},
		"does wait for two blocks": {
			blockWait:        2,
			expectedWaitTime: 2 * blockTime,
		},
		"does wait for three blocks": {
			blockWait:        3,
			expectedWaitTime: 3 * blockTime,
		},
		"does not wait for 0 blocks": {
			blockWait:        0,
			expectedWaitTime: 0,
		},
		"does not wait for negative number of blocks": {
			blockWait:        -1,
			expectedWaitTime: 0,
		},
	}

	for testName, test := range tests {
		test := test
		t.Run(testName, func(t *testing.T) {
			t.Parallel()
			c := Connect(10, 4, big.NewInt(200))
			countWait, err := c.BlockCounter()
			if err != nil {
				t.Fatalf("failed to set up block counter: [%v]", err)
			}

			start := time.Now().UTC()
			countWait.WaitForBlocks(test.blockWait)
			end := time.Now().UTC()

			elapsed := end.Sub(start)

			// Block waiter should wait for test.expectedWaitTime minus some
			// margin at minimum; the margin is needed because clock is not
			// always that precise. Setting it to 5ms for this test.
			minMargin := time.Duration(5) * time.Millisecond
			if elapsed < (test.expectedWaitTime - minMargin) {
				t.Errorf(
					"waited less than expected; expected [%v] at min, waited [%v]",
					test.expectedWaitTime+minMargin,
					elapsed,
				)
			}

			// Block waiter should wait for test.expectedWaitTime plus some
			// margin at maximum; the margin is the time needed for the return
			// instructions to execute, setting it to 25ms for this test.
			maxMargin := time.Duration(25) * time.Millisecond
			if elapsed > (test.expectedWaitTime + maxMargin) {
				t.Errorf(
					"waited longer than expected; expected %v at max, waited %v",
					test.expectedWaitTime,
					elapsed,
				)
			}
		})
	}
}

func TestLocalIsDKGResultPublished(t *testing.T) {
	submittedResults := make(map[*big.Int][]*relaychain.DKGResult)

	submittedRequestID := big.NewInt(1)
	submittedResult := &relaychain.DKGResult{
		GroupPublicKey: []byte{11},
	}

	submittedResults[submittedRequestID] = append(
		submittedResults[submittedRequestID],
		submittedResult,
	)

	localChain := &localChain{
		submittedResults: submittedResults,
	}
	chainHandle := localChain.ThresholdRelay()

	var tests = map[string]struct {
		requestID      *big.Int
		expectedResult bool
	}{
		"matched": {
			requestID:      submittedRequestID,
			expectedResult: true,
		},
		"not matched - different request ID": {
			requestID:      big.NewInt(3),
			expectedResult: false,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			actualResult, err := chainHandle.IsDKGResultPublished(test.requestID)
			if err != nil {
				t.Fatal(err)
			}

			if actualResult != test.expectedResult {
				t.Fatalf("\nexpected: %v\nactual:   %v\n", test.expectedResult, actualResult)
			}
		})
	}
}

func TestSubmitDKGResult(t *testing.T) {
	localChain := Connect(10, 4, big.NewInt(200))
	chainHandle := localChain.ThresholdRelay()

	// Channel for callback on DKG result submission.
	onResultSubmissionCallbackChan := make(chan *event.DKGResultPublication)
	subscription, err := localChain.OnDKGResultPublished(
		func(dkgResultPublication *event.DKGResultPublication) {
			onResultSubmissionCallbackChan <- dkgResultPublication
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	defer subscription.Unsubscribe()

	// Test data.
	requestID0 := big.NewInt(0)
	requestID1 := big.NewInt(1)
	requestID2 := big.NewInt(2)

	dkgResult0 := &relaychain.DKGResult{GroupPublicKey: []byte{00}}
	dkgResult0Hash, _ := chainHandle.CalculateDKGResultHash(dkgResult0)

	dkgResult1 := &relaychain.DKGResult{GroupPublicKey: []byte{11}}
	dkgResult1Hash, _ := chainHandle.CalculateDKGResultHash(dkgResult1)

	// Register a result in the chain as initial state.
	localChain.dkgResults = map[string][]*relaychain.DKGResult{
		requestID0.String(): []*relaychain.DKGResult{dkgResult0},
	}
	localChain.submissions = map[string]relaychain.DKGResultsVotes{
		requestID0.String(): relaychain.DKGResultsVotes{
			dkgResult0Hash: 1,
		},
	}

	var tests = map[string]struct {
		requestID      *big.Int
		resultToSubmit *relaychain.DKGResult

		expectedResultsUpdate      func(initialResults map[string][]*relaychain.DKGResult)
		expectedResultsVotesUpdate func(initialResultsVotes map[string]relaychain.DKGResultsVotes)
		expectedEvent              *event.DKGResultPublication
		expectedError              error
	}{
		"submit a new result for a request ID with no previous submissions": {
			requestID:      requestID1,
			resultToSubmit: dkgResult1,

			expectedResultsUpdate: func(initialResults map[string][]*relaychain.DKGResult) {
				initialResults[requestID1.String()] = []*relaychain.DKGResult{dkgResult1}
			},
			expectedResultsVotesUpdate: func(initialResultsVotes map[string]relaychain.DKGResultsVotes) {
				initialResultsVotes[requestID1.String()] =
					map[relaychain.DKGResultHash]int{
						dkgResult1Hash: 1,
					}
			},
			expectedEvent: &event.DKGResultPublication{
				RequestID:      requestID1,
				GroupPublicKey: dkgResult1.GroupPublicKey[:],
			},
		},
		"submit a result which was previously submitted but for different request ID": {
			requestID:      requestID2,
			resultToSubmit: dkgResult0,
			expectedResultsUpdate: func(initialResults map[string][]*relaychain.DKGResult) {
				initialResults[requestID2.String()] =
					[]*relaychain.DKGResult{dkgResult0}
			},
			expectedResultsVotesUpdate: func(
				initialResultsVotes map[string]relaychain.DKGResultsVotes,
			) {
				initialResultsVotes[requestID2.String()] =
					map[relaychain.DKGResultHash]int{
						dkgResult0Hash: 1,
					}
			},
			expectedEvent: &event.DKGResultPublication{
				RequestID:      requestID2,
				GroupPublicKey: dkgResult0.GroupPublicKey[:],
			},
		},
		"submit a new result for a request ID with previous submissions": {
			requestID:      requestID0,
			resultToSubmit: dkgResult1,

			expectedResultsUpdate: func(initialResults map[string][]*relaychain.DKGResult) {
				initialResults[requestID0.String()] = []*relaychain.DKGResult{dkgResult0, dkgResult1}
			},
			expectedResultsVotesUpdate: func(initialResultsVotes map[string]relaychain.DKGResultsVotes) {
				initialResultsVotes[requestID0.String()] =
					map[relaychain.DKGResultHash]int{
						dkgResult0Hash: 1,
						dkgResult1Hash: 1,
					}
			},
			expectedEvent: &event.DKGResultPublication{
				RequestID:      requestID0,
				GroupPublicKey: dkgResult1.GroupPublicKey[:],
			},
		},
		"submit a result for a request ID which already has this result registered": {
			requestID:      requestID0,
			resultToSubmit: dkgResult0,

			expectedError: fmt.Errorf("result already submitted"),
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			ctx, cancel := newTestContext()
			defer cancel()

			expectedResults := make(map[string][]*relaychain.DKGResult)
			expectedResultsVotes := make(map[string]relaychain.DKGResultsVotes)

			for k, v := range localChain.dkgResults {
				expectedResults[k] = v
			}
			for k, v := range localChain.submissions {
				expectedResultsVotes[k] = v
			}

			if test.expectedResultsUpdate != nil {
				test.expectedResultsUpdate(expectedResults)
			}
			if test.expectedResultsVotesUpdate != nil {
				test.expectedResultsVotesUpdate(expectedResultsVotes)
			}

			if test.expectedEvent != nil {
				currentBlock, _ := localChain.blockCounter.CurrentBlock()
				test.expectedEvent.BlockNumber = uint64(currentBlock)
			}

			waitForCompleted := sync.WaitGroup{}
			waitForCompleted.Add(1)

			chainHandle.SubmitDKGResult(test.requestID, test.resultToSubmit).
				// Validate the promise.
				OnComplete(func(event *event.DKGResultPublication, err error) {
					waitForCompleted.Done()

					if !reflect.DeepEqual(test.expectedError, err) {
						t.Errorf("\nexpected: %v\nactual:   %v\n", test.expectedError, err)
					}
					if !reflect.DeepEqual(test.expectedEvent, event) {
						t.Errorf("\nexpected: %+v\nactual:   %+v\n", test.expectedEvent, event)
					}
				})
			waitForCompleted.Wait()

			// Validate registered results and votes.
			if !reflect.DeepEqual(expectedResults, localChain.dkgResults) {
				t.Errorf("\nexpected: %+v\nactual:   %+v\n",
					expectedResults,
					localChain.dkgResults,
				)
			}
			if !reflect.DeepEqual(expectedResultsVotes, localChain.submissions) {
				t.Errorf("\nexpected: %+v\nactual:   %+v\n",
					expectedResultsVotes,
					localChain.submissions,
				)
			}

			// Validate event in callback.
			select {
			case dkgResultPublicationEvent := <-onResultSubmissionCallbackChan:
				if !reflect.DeepEqual(test.expectedEvent, dkgResultPublicationEvent) {
					t.Errorf("\nexpected: %v\nactual:   %v\n",
						test.expectedEvent,
						dkgResultPublicationEvent,
					)
				}
			case <-ctx.Done():
				if test.expectedError == nil {
					t.Errorf("expected event was not emitted")
				}
			}

		})
	}
}

func TestLocalOnDKGResultPublishedUnsubscribe(t *testing.T) {
	ctx, cancel := newTestContext()
	defer cancel()

	localChain := &localChain{
		submittedResults:             make(map[*big.Int][]*relaychain.DKGResult),
		dkgResultPublicationHandlers: make(map[int]func(dkgResultPublication *event.DKGResultPublication)),
	}
	relay := localChain.ThresholdRelay()

	dkgResultPublicationChan := make(chan *event.DKGResultPublication)
	subscription, err := localChain.OnDKGResultPublished(
		func(dkgResultPublication *event.DKGResultPublication) {
			dkgResultPublicationChan <- dkgResultPublication
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	// Unsubscribe from the event - from this point, callback should
	// never be called.
	subscription.Unsubscribe()

	relay.SubmitDKGResult(big.NewInt(999), &relaychain.DKGResult{
		GroupPublicKey: []byte{88},
	})

	select {
	case <-dkgResultPublicationChan:
		t.Fatalf("event should not be emitted - I have unsubscribed!")
	case <-ctx.Done():
		// ok
	}
}

func TestCalculateDKGResultHash(t *testing.T) {
	localChain := &localChain{}

	dkgResult := &relaychain.DKGResult{
		Success:        true,
		GroupPublicKey: []byte{3, 40, 200},
		Disqualified:   []byte{1, 0, 1, 0},
		Inactive:       []byte{0, 1, 1, 0},
	}
	expectedHashString := "135a0a776b24afbdb70a3548d2b01d197f67972b7482df68703caeae4453134e"

	actualHash, err := localChain.CalculateDKGResultHash(dkgResult)
	if err != nil {
		t.Fatal(err)
	}

	expectedHash := relaychain.DKGResultHash{}
	copy(
		expectedHash[:],
		common.Hex2Bytes(expectedHashString)[:32],
	)

	if expectedHash != actualHash {
		t.Fatalf("\nexpected: %x\nactual:   %x\n",
			expectedHash,
			actualHash,
		)
	}
}

func TestVoteOnDKGResult(t *testing.T) {
	chain := Connect(10, 4, big.NewInt(200)).(*localChain)
	chainHandle := chain.ThresholdRelay()

	// Channel for callback on DKG result submission.
	onResultVoteCallbackChan := make(chan *event.DKGResultVote)
	subscription, err := chain.OnDKGResultVote(
		func(event *event.DKGResultVote) {
			onResultVoteCallbackChan <- event
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	defer subscription.Unsubscribe()

	// Test data.
	requestID0 := big.NewInt(0)
	requestID1 := big.NewInt(1)

	dkgResult0Hash := relaychain.DKGResultHash{00}
	dkgResult1Hash := relaychain.DKGResultHash{11}

	// Register a result in the chain as initial state.
	chain.dkgResultsVotes = map[string]relaychain.DKGResultsVotes{
		requestID0.String(): relaychain.DKGResultsVotes{
			dkgResult0Hash: 1,
		},
	}
	chain.alreadySubmittedOrVoted = map[string][]int{
		requestID0.String(): []int{0},
	}

	var tests = map[string]struct {
		requestID     *big.Int
		memberIndex   int
		dkgResultHash relaychain.DKGResultHash

		expectedResultsVotesUpdate            func(initialResultsVotes map[string]relaychain.DKGResultsVotes)
		expectedAlreadySubmittedOrVotedUpdate func(alreadySubmittedOrVoted map[string][]int)
		expectedEvent                         *event.DKGResultVote
		expectedError                         error
	}{
		"submit a vote for a request ID with no previous submissions": {
			requestID:     requestID1,
			memberIndex:   1,
			dkgResultHash: dkgResult1Hash,

			expectedError: fmt.Errorf("no registered submissions or votes for given request id"),
		},
		"submit a vote with request ID with previous submissions but not existing DKG result hash": {
			requestID:     requestID0,
			memberIndex:   1,
			dkgResultHash: dkgResult1Hash,

			expectedError: fmt.Errorf("result hash is not registered in dkg results votes map"),
		},
		"submit again a vote for the same request ID": {
			requestID:     requestID0,
			memberIndex:   0,
			dkgResultHash: dkgResult1Hash,

			expectedError: fmt.Errorf("this member already submitted or voted on DKG result for given request id"),
		},
		"submit a vote for a DKG result submitted before": {
			requestID:     requestID0,
			memberIndex:   1,
			dkgResultHash: dkgResult0Hash,

			expectedResultsVotesUpdate: func(initialResultsVotes map[string]relaychain.DKGResultsVotes) {
				initialResultsVotes[requestID0.String()] =
					map[relaychain.DKGResultHash]int{
						dkgResult0Hash: 2,
					}
			},
			expectedAlreadySubmittedOrVotedUpdate: func(alreadySubmittedOrVoted map[string][]int) {
				alreadySubmittedOrVoted[requestID0.String()] = append(alreadySubmittedOrVoted[requestID0.String()],
					1,
				)
			},
			expectedEvent: &event.DKGResultVote{
				RequestID:     requestID0,
				MemberIndex:   1,
				DKGResultHash: dkgResult0Hash,
			},
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			ctx, cancel := newTestContext()
			defer cancel()

			expectedResultsVotes := make(map[string]relaychain.DKGResultsVotes)
			for k, v := range chain.dkgResultsVotes {
				expectedResultsVotes[k] = v
			}

			if test.expectedResultsVotesUpdate != nil {
				test.expectedResultsVotesUpdate(expectedResultsVotes)
			}

			expectedAlreadySubmittedOrVoted := make(map[string][]int)
			for k, v := range chain.alreadySubmittedOrVoted {
				expectedAlreadySubmittedOrVoted[k] = v
			}
			if test.expectedAlreadySubmittedOrVotedUpdate != nil {
				test.expectedAlreadySubmittedOrVotedUpdate(
					expectedAlreadySubmittedOrVoted,
				)
			}

			if test.expectedEvent != nil {
				currentBlock, _ := chain.blockCounter.CurrentBlock()
				test.expectedEvent.BlockNumber = uint64(currentBlock)
			}

			waitForCompleted := sync.WaitGroup{}
			waitForCompleted.Add(1)

			chainHandle.VoteOnDKGResult(
				test.requestID,
				test.memberIndex,
				test.dkgResultHash,
			).
				// Validate the promise.
				OnComplete(func(event *event.DKGResultVote, err error) {
					waitForCompleted.Done()

					if !reflect.DeepEqual(test.expectedError, err) {
						t.Errorf("\nexpected: %v\nactual:   %v\n", test.expectedError, err)
					}
					if !reflect.DeepEqual(test.expectedEvent, event) {
						t.Errorf("\nexpected: %+v\nactual:   %+v\n", test.expectedEvent, event)
					}
				})
			waitForCompleted.Wait()

			// Validate registered votes.
			if !reflect.DeepEqual(expectedResultsVotes, chain.dkgResultsVotes) {
				t.Errorf("\nexpected: %+v\nactual:   %+v\n",
					expectedResultsVotes,
					chain.dkgResultsVotes,
				)
			}
			if !reflect.DeepEqual(
				expectedAlreadySubmittedOrVoted,
				chain.alreadySubmittedOrVoted,
			) {
				t.Errorf("\nexpected: %+v\nactual:   %+v\n",
					expectedAlreadySubmittedOrVoted,
					chain.alreadySubmittedOrVoted,
				)
			}

			// Validate event in callback.
			select {
			case event := <-onResultVoteCallbackChan:
				if !reflect.DeepEqual(test.expectedEvent, event) {
					t.Errorf("\nexpected: %v\nactual:   %v\n",
						test.expectedEvent,
						event,
					)
				}
			case <-ctx.Done():
				if test.expectedError == nil {
					t.Errorf("expected event was not emitted")
				}
			}
		})
	}
}

func newTestContext(timeout ...time.Duration) (context.Context, context.CancelFunc) {
	defaultTimeout := 3 * time.Second
	if len(timeout) > 0 {
		defaultTimeout = timeout[0]
	}
	return context.WithTimeout(context.Background(), defaultTimeout)
}
