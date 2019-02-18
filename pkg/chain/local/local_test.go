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

func TestLocalSubmitDKGResult(t *testing.T) {
	ctx, cancel := newTestContext()
	defer cancel()

	// Initialize local chain.
	submittedResults := make(map[*big.Int][]*relaychain.DKGResult)
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
	submittedResult11 := &relaychain.DKGResult{
		GroupPublicKey: []byte{11},
	}
	expectedEvent1 := &event.DKGResultPublication{
		RequestID:      requestID1,
		GroupPublicKey: submittedResult11.GroupPublicKey[:],
	}

	chainHandle.SubmitDKGResult(requestID1, submittedResult11)
	if !reflect.DeepEqual(
		localChain.submittedResults[requestID1],
		[]*relaychain.DKGResult{submittedResult11},
	) {
		t.Fatalf("invalid submitted results for request ID %v\nexpected: %v\nactual:   %v\n",
			requestID1,
			[]*relaychain.DKGResult{submittedResult11},
			localChain.submittedResults[requestID1],
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
		GroupPublicKey: submittedResult11.GroupPublicKey[:],
	}

	chainHandle.SubmitDKGResult(requestID2, submittedResult11)
	if !reflect.DeepEqual(
		localChain.submittedResults[requestID2],
		[]*relaychain.DKGResult{submittedResult11},
	) {
		t.Fatalf("invalid submitted results for request ID %v\nexpected: %v\nactual:   %v\n",
			requestID2,
			[]*relaychain.DKGResult{submittedResult11},
			localChain.submittedResults[requestID2],
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
	chainHandle.SubmitDKGResult(requestID1, submittedResult11)
	if !reflect.DeepEqual(
		localChain.submittedResults[requestID1],
		[]*relaychain.DKGResult{submittedResult11},
	) {
		t.Fatalf("invalid submitted results for request ID %v\nexpected: %v\nactual:   %v\n",
			requestID1,
			[]*relaychain.DKGResult{submittedResult11},
			localChain.submittedResults[requestID1],
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

func newTestContext(timeout ...time.Duration) (context.Context, context.CancelFunc) {
	defaultTimeout := 3 * time.Second
	if len(timeout) > 0 {
		defaultTimeout = timeout[0]
	}
	return context.WithTimeout(context.Background(), defaultTimeout)
}
