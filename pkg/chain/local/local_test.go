package local

import (
	"context"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
	"github.com/keep-network/keep-core/pkg/operator"

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

func TestLocalRequestRelayEntry(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	chainHandle := Connect(10, 4, big.NewInt(200)).ThresholdRelay()
	seed := big.NewInt(42)
	relayRequestPromise := chainHandle.RequestRelayEntry(seed)

	done := make(chan *event.Request)
	relayRequestPromise.OnSuccess(func(entry *event.Request) {
		done <- entry
	}).OnFailure(func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	})

	select {
	case entry := <-done:
		if entry.Seed.Cmp(seed) != 0 {
			t.Fatalf(
				"Unexpected relay entry seed\nExpected: [%v]\nActual:  [%v]",
				seed,
				entry.Seed.Int64(),
			)
		}
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}

}

func TestLocalSubmitRelayEntry(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	chainHandle := Connect(10, 4, big.NewInt(200)).ThresholdRelay()
	requestID := int64(19)
	relayEntryPromise := chainHandle.SubmitRelayEntry(
		&event.Entry{
			RequestID:   big.NewInt(requestID),
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
		if entry.RequestID.Int64() != requestID {
			t.Fatalf(
				"Unexpected relay entry request id\nExpected: [%v]\nActual:  [%v]",
				requestID,
				entry.RequestID.Int64(),
			)
		}
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}

}

func TestLocalOnRelayEntryGenerated(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	chainHandle := Connect(10, 4, big.NewInt(200)).ThresholdRelay()

	eventFired := make(chan *event.Entry)

	subscription, err := chainHandle.OnRelayEntryGenerated(
		func(entry *event.Entry) {
			eventFired <- entry
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	defer subscription.Unsubscribe()

	expectedEntry := &event.Entry{
		RequestID:   big.NewInt(42),
		Value:       big.NewInt(19),
		GroupPubKey: []byte("1"),
		Seed:        big.NewInt(30),
		BlockNumber: uint64(123),
	}

	chainHandle.SubmitRelayEntry(expectedEntry)

	select {
	case event := <-eventFired:
		if !reflect.DeepEqual(event, expectedEntry) {
			t.Fatalf(
				"Unexpected relay entry\nExpected: [%v]\nActual:   [%v]",
				expectedEntry,
				event,
			)
		}
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}
}

func TestLocalOnRelayEntryGeneratedUnsubscribed(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	chainHandle := Connect(10, 4, big.NewInt(200)).ThresholdRelay()

	eventFired := make(chan *event.Entry)

	subscription, err := chainHandle.OnRelayEntryGenerated(
		func(entry *event.Entry) {
			eventFired <- entry
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	subscription.Unsubscribe()

	chainHandle.SubmitRelayEntry(
		&event.Entry{},
	)

	select {
	case event := <-eventFired:
		t.Fatalf("Event should have not been received due to the cancelled subscription: [%v]", event)
	case <-ctx.Done():
		// expected execution of goroutine
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
			actualResult, err := chainHandle.IsDKGResultSubmitted(test.requestID)
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
	localChain := Connect(10, 4, big.NewInt(200)).(*localChain)

	chainHandle := localChain.ThresholdRelay()

	// Channel for DKGResultSubmission events.
	DKGResultSubmissionChan := make(chan *event.DKGResultSubmission)
	chainHandle.OnDKGResultSubmitted(
		func(DKGResultSubmission *event.DKGResultSubmission) {
			DKGResultSubmissionChan <- DKGResultSubmission
		},
	)

	if len(localChain.submittedResults) > 0 {
		t.Fatalf("initial submitted results map is not empty")
	}

	// Submit new result for request ID 1
	requestID1 := big.NewInt(1)
	memberIndex := uint32(1)
	submittedResult11 := &relaychain.DKGResult{
		GroupPublicKey: []byte{11},
	}
	expectedEvent1 := &event.DKGResultSubmission{
		RequestID:      requestID1,
		MemberIndex:    memberIndex,
		GroupPublicKey: submittedResult11.GroupPublicKey[:],
		BlockNumber:    0,
	}

	signatures := map[group.MemberIndex]operator.Signature{
		1: operator.Signature{101},
		2: operator.Signature{102},
		3: operator.Signature{103},
	}

	chainHandle.SubmitDKGResult(requestID1, 1, submittedResult11, signatures)
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
	case DKGResultSubmissionEvent := <-DKGResultSubmissionChan:
		if !reflect.DeepEqual(expectedEvent1, DKGResultSubmissionEvent) {
			t.Fatalf("\nexpected: %+v\nactual:   %+v\n",
				expectedEvent1,
				DKGResultSubmissionEvent,
			)
		}
	case <-ctx.Done():
		t.Fatalf("expected event was not emitted")
	}

	// Submit the same result for request ID 2
	requestID2 := big.NewInt(2)
	expectedEvent2 := &event.DKGResultSubmission{
		RequestID:      requestID2,
		MemberIndex:    memberIndex,
		GroupPublicKey: submittedResult11.GroupPublicKey[:],
	}

	chainHandle.SubmitDKGResult(requestID2, 1, submittedResult11, signatures)
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
	case DKGResultSubmissionEvent := <-DKGResultSubmissionChan:
		if !reflect.DeepEqual(expectedEvent2, DKGResultSubmissionEvent) {
			t.Fatalf("\nexpected: %v\nactual:   %v\n",
				expectedEvent2,
				DKGResultSubmissionEvent,
			)
		}
	case <-ctx.Done():
		t.Fatalf("expected event was not emitted")
	}

	// Submit already submitted result for request ID 1
	chainHandle.SubmitDKGResult(requestID1, 1, submittedResult11, signatures)
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
	case DKGResultSubmissionEvent := <-DKGResultSubmissionChan:
		t.Fatalf("unexpected event was emitted: %v", DKGResultSubmissionEvent)
	case <-ctx.Done():
		t.Logf("DKG result publication event not generated")
	}
}

func TestLocalOnDKGResultPublishedUnsubscribe(t *testing.T) {
	ctx, cancel := newTestContext()
	defer cancel()

	localChain := Connect(10, 4, big.NewInt(200)).(*localChain)

	chainHandle := localChain.ThresholdRelay()

	DKGResultSubmissionChan := make(chan *event.DKGResultSubmission)
	subscription, err := localChain.OnDKGResultSubmitted(
		func(DKGResultSubmission *event.DKGResultSubmission) {
			DKGResultSubmissionChan <- DKGResultSubmission
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	// Unsubscribe from the event - from this point, callback should
	// never be called.
	subscription.Unsubscribe()

	chainHandle.SubmitDKGResult(
		big.NewInt(999),
		1,
		&relaychain.DKGResult{
			GroupPublicKey: []byte{88},
		},
		nil, // TODO: Update test to include signatures
	)

	select {
	case <-DKGResultSubmissionChan:
		t.Fatalf("event should not be emitted - I have unsubscribed!")
	case <-ctx.Done():
		// ok
	}
}

func TestCalculateDKGResultHash(t *testing.T) {
	localChain := &localChain{}

	dkgResult := &relaychain.DKGResult{
		GroupPublicKey: []byte{3, 40, 200},
		Disqualified:   []byte{1, 0, 1, 0},
		Inactive:       []byte{0, 1, 1, 0},
	}
	expectedHashString := "f65d6c5e938537224bbd2716d2f24895746a756978d29e1eaaf46fb97a555716"

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

func newTestContext(timeout ...time.Duration) (context.Context, context.CancelFunc) {
	defaultTimeout := 3 * time.Second
	if len(timeout) > 0 {
		defaultTimeout = timeout[0]
	}
	return context.WithTimeout(context.Background(), defaultTimeout)
}
