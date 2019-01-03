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

func TestLocalSubmitRelayEntry(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	chainHandle := Connect(10, 4).ThresholdRelay()
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
		blockWait    int
		expectation  time.Duration
		errorMessage string
	}{
		"does wait for a block": {
			blockWait:    1,
			expectation:  time.Duration(525) * time.Millisecond,
			errorMessage: "Failed to wait for a single block; expected %s but took %s.",
		},
		"waited for a longer time": {
			blockWait:    2,
			expectation:  time.Duration(525*2) * time.Millisecond,
			errorMessage: "Failed to wait for 2 blocks; expected %s but took %s.",
		},
		"doesn't wait if 0 blocks": {
			blockWait:    0,
			expectation:  time.Duration(20) * time.Millisecond,
			errorMessage: "Failed for a 0 block wait; expected %s but took %s.",
		},
		"invalid value": {
			blockWait:    -1,
			expectation:  time.Duration(20) * time.Millisecond,
			errorMessage: "Waiting for a time when it should have errored; expected %s but took %s.",
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			t.Parallel()

			c := Connect(10, 4)
			countWait, err := c.BlockCounter()
			if err != nil {
				t.Fatalf("failed to set up block counter: [%v]", err)
			}

			start := time.Now().UTC()
			countWait.WaitForBlocks(test.blockWait)
			end := time.Now().UTC()

			elapsed := end.Sub(start)
			if test.expectation < elapsed {
				t.Errorf(test.errorMessage, test.expectation, elapsed)
			}
		})
	}
}

func TestLocalIsDKGResultPublished(t *testing.T) {
	submittedResults := make(map[*big.Int][]*relaychain.DKGResult)

	submittedRequestID1 := big.NewInt(1)
	submittedResult11 := &relaychain.DKGResult{
		GroupPublicKey: big.NewInt(11),
	}

	submittedRequestID2 := big.NewInt(2)
	submittedResult21 := &relaychain.DKGResult{
		GroupPublicKey: big.NewInt(21),
	}

	submittedResults[submittedRequestID1] = append(
		submittedResults[submittedRequestID1],
		submittedResult11,
	)

	submittedResults[submittedRequestID2] = append(
		submittedResults[submittedRequestID2],
		submittedResult21,
	)

	localChain := &localChain{
		submittedResults: submittedResults,
	}
	chainHandle := localChain.ThresholdRelay()

	var tests = map[string]struct {
		requestID      *big.Int
		dkgResult      *relaychain.DKGResult
		expectedResult bool
	}{
		"matched": {
			requestID:      submittedRequestID2,
			dkgResult:      submittedResult21,
			expectedResult: true,
		},
		"not matched - different request ID": {
			requestID:      submittedRequestID2,
			dkgResult:      submittedResult11,
			expectedResult: false,
		},
		"not matched - different request ID and DKG Result": {
			requestID:      big.NewInt(3),
			dkgResult:      &relaychain.DKGResult{GroupPublicKey: big.NewInt(31)},
			expectedResult: false,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			actualResult := chainHandle.IsDKGResultPublished(test.requestID, test.dkgResult)

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
		GroupPublicKey: big.NewInt(11),
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
		if dkgResultPublicationEvent.RequestID.Cmp(requestID1) != 0 {
			t.Fatalf("\nexpected: %v\nactual:   %v\n",
				requestID1,
				dkgResultPublicationEvent.RequestID,
			)
		}
	case <-ctx.Done():
		t.Fatalf("expected event was not emitted")
	}

	// Submit the same result for request ID 2
	requestID2 := big.NewInt(2)

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
		if dkgResultPublicationEvent.RequestID.Cmp(requestID2) != 0 {
			t.Fatalf("\nexpected: %v\nactual:   %v\n",
				requestID2,
				dkgResultPublicationEvent.RequestID,
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
	subscription := localChain.OnDKGResultPublished(
		func(dkgResultPublication *event.DKGResultPublication) {
			dkgResultPublicationChan <- dkgResultPublication
		},
	)

	// Unsubscribe from the event - from this point, callback should
	// never be called.
	subscription.Unsubscribe()

	relay.SubmitDKGResult(big.NewInt(999), &relaychain.DKGResult{
		GroupPublicKey: big.NewInt(888),
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
