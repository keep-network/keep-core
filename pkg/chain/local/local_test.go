package local

import (
	"context"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/keep-network/keep-core/pkg/beacon/relay/event"

	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
)

func TestSubmitTicketAndGetSelectedTickets(t *testing.T) {
	c := Connect(10, 4, big.NewInt(200))
	chain := c.ThresholdRelay()

	ticket1 := &relaychain.Ticket{Value: big.NewInt(1)}
	ticket2 := &relaychain.Ticket{Value: big.NewInt(2)}
	ticket3 := &relaychain.Ticket{Value: big.NewInt(3)}
	ticket4 := &relaychain.Ticket{Value: big.NewInt(4)}

	chain.SubmitTicket(ticket3)
	chain.SubmitTicket(ticket1)
	chain.SubmitTicket(ticket4)
	chain.SubmitTicket(ticket2)

	expectedResult := []*relaychain.Ticket{
		ticket1, ticket2, ticket3, ticket4,
	}

	actualResult, err := chain.GetSelectedTickets()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expectedResult, actualResult) {
		t.Fatalf(
			"\nexpected: %v\nactual:   %v\n",
			expectedResult,
			actualResult,
		)
	}
}
func TestLocalSubmitRelayEntry(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	chainHandle := Connect(10, 4, big.NewInt(200)).ThresholdRelay()
	relayEntryPromise := chainHandle.SubmitRelayEntry(
		&event.Entry{
			RequestID: big.NewInt(int64(19)),
			GroupID:   big.NewInt(int64(1)),
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

			// Block waiter should wait for test.expectedWaitTime at minimum.
			if elapsed < test.expectedWaitTime {
				t.Errorf(
					"waited less than expected; expected [%v] at min, waited [%v]",
					test.expectedWaitTime,
					elapsed,
				)
			}

			// Block waiter should wait for test.expectedWaitTime plus some
			// margin at maximum; the margin is the time needed for the return
			// instructions to execute, setting it to 25ms for this test.
			margin := time.Duration(25) * time.Millisecond
			if elapsed > (test.expectedWaitTime + margin) {
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
	submittedResults := make(map[string][]*relaychain.DKGResult)

	submittedRequestID := big.NewInt(1)
	submittedResult := &relaychain.DKGResult{
		GroupPublicKey: []byte{11},
	}

	submittedResults[submittedRequestID.String()] = append(
		submittedResults[submittedRequestID.String()],
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

func TestLocalSubmitDKGResult(t *testing.T) {
	ctx, cancel := newTestContext()
	defer cancel()

	// Initialize local chain.
	submittedResults := make(map[string][]*relaychain.DKGResult)
	localChain := &localChain{
		submittedResults:             submittedResults,
		dkgResultPublicationHandlers: make(map[int]func(dkgResultPublication *event.DKGResultPublication)),
	}
	chainHandle := localChain.ThresholdRelay()

	// Channel for DKGResultPublication events.
	dkgResultPublicationChan := make(chan *event.DKGResultPublication)
	localChain.OnDKGResultPublished(
		func(dkgResultPublication *event.DKGResultPublication) {
			dkgResultPublicationChan <- dkgResultPublication
		},
	)

	if len(localChain.submittedResults) > 0 {
		t.Fatalf("initial submitted results map is not empty")
	}

	// Submit new result for request ID 1
	requestID1 := big.NewInt(1)
	submittedResult1 := &relaychain.DKGResult{
		GroupPublicKey: []byte{11},
	}
	expectedEvent1 := &event.DKGResultPublication{
		RequestID:      requestID1,
		GroupPublicKey: submittedResult1.GroupPublicKey[:],
	}

	chainHandle.SubmitDKGResult(requestID1, submittedResult1)
	if !reflect.DeepEqual(
		localChain.submittedResults[requestID1.String()],
		[]*relaychain.DKGResult{submittedResult1},
	) {
		t.Fatalf("invalid submitted results for request ID %v\nexpected: %v\nactual:   %v\n",
			requestID1,
			[]*relaychain.DKGResult{submittedResult1},
			localChain.submittedResults[requestID1.String()],
		)
	}
	select {
	case dkgResultPublicationEvent := <-dkgResultPublicationChan:
		if !reflect.DeepEqual(expectedEvent1, dkgResultPublicationEvent) {
			t.Fatalf("\nexpected: %v\nactual:   %v\n",
				expectedEvent1,
				dkgResultPublicationEvent,
			)
		}
	case <-ctx.Done():
		t.Fatalf("expected event was not emitted")
	}

	// Submit the same result for request ID 2
	requestID2 := big.NewInt(2)
	expectedEvent2 := &event.DKGResultPublication{
		RequestID:      requestID2,
		GroupPublicKey: submittedResult1.GroupPublicKey[:],
	}

	chainHandle.SubmitDKGResult(requestID2, submittedResult1)
	if !reflect.DeepEqual(
		localChain.submittedResults[requestID2.String()],
		[]*relaychain.DKGResult{submittedResult1},
	) {
		t.Fatalf("invalid submitted results for request ID %v\nexpected: %v\nactual:   %v\n",
			requestID2,
			[]*relaychain.DKGResult{submittedResult1},
			localChain.submittedResults[requestID2.String()],
		)
	}
	select {
	case dkgResultPublicationEvent := <-dkgResultPublicationChan:
		if !reflect.DeepEqual(expectedEvent2, dkgResultPublicationEvent) {
			t.Fatalf("\nexpected: %v\nactual:   %v\n",
				expectedEvent2,
				dkgResultPublicationEvent,
			)
		}
	case <-ctx.Done():
		t.Fatalf("expected event was not emitted")
	}

	// Submit already submitted result for request ID 1
	chainHandle.SubmitDKGResult(requestID1, submittedResult1)
	if !reflect.DeepEqual(
		localChain.submittedResults[requestID1.String()],
		[]*relaychain.DKGResult{submittedResult1},
	) {
		t.Fatalf("invalid submitted results for request ID %v\nexpected: %v\nactual:   %v\n",
			requestID1,
			[]*relaychain.DKGResult{submittedResult1},
			localChain.submittedResults[requestID1.String()],
		)
	}
	select {
	case dkgResultPublicationEvent := <-dkgResultPublicationChan:
		t.Fatalf("unexpected event was emitted: %v", dkgResultPublicationEvent)
	case <-ctx.Done():
		t.Logf("DKG result publication event not generated")
	}
}

func TestLocalOnDKGResultPublishedUnsubscribe(t *testing.T) {
	ctx, cancel := newTestContext()
	defer cancel()

	localChain := &localChain{
		submittedResults:             make(map[string][]*relaychain.DKGResult),
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

func newTestContext(timeout ...time.Duration) (context.Context, context.CancelFunc) {
	defaultTimeout := 3 * time.Second
	if len(timeout) > 0 {
		defaultTimeout = timeout[0]
	}
	return context.WithTimeout(context.Background(), defaultTimeout)
}

func TestDKGResultVote(t *testing.T) {
	requestID := big.NewInt(12)

	var tests = map[string]struct {
		expectedVotes  int
		callVoteNtimes int
		requestID      *big.Int
	}{
		"after submission and 1 vote": {
			expectedVotes:  2,
			callVoteNtimes: 1,
			requestID:      requestID,
		},
		"after submission and 2 votes": {
			expectedVotes:  3,
			callVoteNtimes: 2,
			requestID:      requestID,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			localChain := &localChain{
				submittedResults:             make(map[string][]*relaychain.DKGResult),
				dkgResultPublicationHandlers: make(map[int]func(dkgResultPublication *event.DKGResultPublication)),
			}
			chainHandle := localChain.ThresholdRelay()
			dkgResult := &relaychain.DKGResult{
				Success:        true,
				GroupPublicKey: []byte{11},
			}
			chainHandle.SubmitDKGResult(test.requestID, dkgResult)
			for i := 0; i < test.callVoteNtimes; i++ {
				chainHandle.DKGResultVote(test.requestID, dkgResult.Hash())
			}
			submissions := chainHandle.GetDKGSubmissions(test.requestID)
			if submissions.DKGSubmissions[0].Votes != test.expectedVotes {
				t.Fatalf("\nexpected: %v\nactual:   %v\n",
					test.expectedVotes,
					submissions.DKGSubmissions[0].Votes,
				)
			}
		})
	}
}

func TestGetDKGSubmissions(t *testing.T) {
	requestID := big.NewInt(12)
	localChain := &localChain{
		submittedResults:             make(map[string][]*relaychain.DKGResult),
		dkgResultPublicationHandlers: make(map[int]func(dkgResultPublication *event.DKGResultPublication)),
	}
	chainHandle := localChain.ThresholdRelay()
	dkgResult := &relaychain.DKGResult{
		Success:        true,
		GroupPublicKey: []byte{11},
	}
	chainHandle.SubmitDKGResult(requestID, dkgResult)
	submissions := chainHandle.GetDKGSubmissions(requestID)

	if len(submissions.DKGSubmissions) != 1 {
		t.Fatalf("\nexpected: %v\nactual:   %v\n",
			1,
			len(submissions.DKGSubmissions),
		)
	}

	if submissions.DKGSubmissions[0].Votes != 1 {
		t.Fatalf("\nexpected: %v\nactual:   %v\n",
			1,
			submissions.DKGSubmissions[0].Votes,
		)
	}

	if !submissions.DKGSubmissions[0].DKGResult.Equals(dkgResult) {
		t.Fatalf("\nexpected: %v\nactual:   %v\n",
			dkgResult,
			submissions.DKGSubmissions[0].DKGResult,
		)
	}
}

func TestOnDKGResultVote(t *testing.T) {
	requestID := big.NewInt(12)
	requestID2 := big.NewInt(22)

	var tests = map[string]struct {
		callVoteNtimes int
		requestID      *big.Int
		expectMessage  bool
		groupPublicKey []byte
	}{
		"after submission and 1 vote": {
			callVoteNtimes: 1,
			requestID:      requestID,
			expectMessage:  true,
			groupPublicKey: []byte{11},
		},
		"verify multiple votes": {
			callVoteNtimes: 3,
			requestID:      requestID2,
			expectMessage:  false,
			groupPublicKey: []byte{21},
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			localChain := &localChain{
				submittedResults:             make(map[string][]*relaychain.DKGResult),
				dkgResultPublicationHandlers: make(map[int]func(dkgResultPublication *event.DKGResultPublication)),
			}
			chainHandle := localChain.ThresholdRelay()
			dkgResult := &relaychain.DKGResult{
				Success:        true,
				GroupPublicKey: test.groupPublicKey,
			}
			chainHandle.SubmitDKGResult(test.requestID, dkgResult)

			messages := make(chan string, 3) // big enough
			defer close(messages)
			localChain.OnDKGResultVote(func(dkgResultVote *event.DKGResultVote) {
				messages <- "got vote"
			})
			for i := 0; i < test.callVoteNtimes; i++ {
				chainHandle.DKGResultVote(test.requestID, dkgResult.Hash())
			}

			nMessageReceived := test.callVoteNtimes
			for i := 0; i < 30; i++ {
				select {
				case msg, ok := <-messages:
					if ok {
						nMessageReceived--
						if msg != "got vote" {
							t.Fatalf("\nexpected: %v\nactual:   %v\n",
								"got vote",
								msg,
							)
						}
						if nMessageReceived <= 0 {
							goto done
						}
					} else {
						t.Fatalf("\nexpected: %v\nactual:   %v\n",
							test.callVoteNtimes,
							"failed : channel closed",
						)
					}
				default:
				}
				time.Sleep(100 * time.Millisecond)
			}
			t.Fatalf("\nexpected: %v\nactual:   %v\n",
				test.callVoteNtimes,
				"failed to get correct number of calls",
			)
		done:
		})
	}
}

func TestOnDKGResultVoteBadRequestID(t *testing.T) {
	requestID := big.NewInt(12)
	requestID2 := big.NewInt(22)

	groupPublicKey := []byte{11}

	localChain := &localChain{
		submittedResults:             make(map[string][]*relaychain.DKGResult),
		dkgResultPublicationHandlers: make(map[int]func(dkgResultPublication *event.DKGResultPublication)),
	}
	chainHandle := localChain.ThresholdRelay()
	dkgResult := &relaychain.DKGResult{
		Success:        true,
		GroupPublicKey: groupPublicKey,
	}
	chainHandle.SubmitDKGResult(requestID, dkgResult)

	chainHandle.DKGResultVote(requestID2, dkgResult.Hash())
	submissions := chainHandle.GetDKGSubmissions(requestID)

	// check Votes == 1
	if submissions.DKGSubmissions[0].Votes != 1 {
		t.Fatalf("\nexpected: %v\nactual:   %v\n",
			1,
			submissions.DKGSubmissions[0].Votes,
		)
	}

	chainHandle.DKGResultVote(requestID, dkgResult.Hash())
	submissions = chainHandle.GetDKGSubmissions(requestID)

	// check Votes == 2
	if submissions.DKGSubmissions[0].Votes != 2 {
		t.Fatalf("\nexpected: %v\nactual:   %v\n",
			2,
			submissions.DKGSubmissions[0].Votes,
		)
	}

}

/*
# github.com/keep-network/keep-core/pkg/chain/local
./local_test.go:229:42: undefined: submittedResult1
./local_test.go:232:27: undefined: submittedResult1
./local_test.go:236:28: undefined: submittedResult1
./local_test.go:259:42: undefined: submittedResult1
./local_test.go:262:27: undefined: submittedResult1
./local_test.go:266:28: undefined: submittedResult1
./local_test.go:266:28: too many errors
*/
