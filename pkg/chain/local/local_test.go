package local

import (
	"context"
	"fmt"
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

func TestLocalSubmitRelayEntry(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	chainHandle := Connect(10, 4, big.NewInt(200)).ThresholdRelay()
	signingId := int64(19)
	relayEntryPromise := chainHandle.SubmitRelayEntry(
		&event.Entry{
			SigningId:   big.NewInt(signingId),
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
		if entry.SigningId.Int64() != signingId {
			t.Fatalf(
				"Unexpected relay entry request id\nExpected: [%v]\nActual:  [%v]",
				signingId,
				entry.SigningId.Int64(),
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
		SigningId:   big.NewInt(42),
		Value:       big.NewInt(19),
		GroupPubKey: []byte("1"),
		Seed:        big.NewInt(30),
		BlockNumber: 123,
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

func TestLocalOnGroupRegistered(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	chainHandle := Connect(10, 4, big.NewInt(200)).ThresholdRelay()

	eventFired := make(chan *event.GroupRegistration)

	subscription, err := chainHandle.OnGroupRegistered(
		func(entry *event.GroupRegistration) {
			eventFired <- entry
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	defer subscription.Unsubscribe()

	groupPublicKey := []byte("1")
	signingId := big.NewInt(42)
	memberIndex := group.MemberIndex(1)
	dkgResult := &relaychain.DKGResult{GroupPublicKey: groupPublicKey}
	signatures := make(map[group.MemberIndex]operator.Signature)

	chainHandle.SubmitDKGResult(signingId, memberIndex, dkgResult, signatures)

	expectedGroupRegistrationEvent := &event.GroupRegistration{
		GroupPublicKey: groupPublicKey,
		SigningId:      signingId,
	}

	select {
	case event := <-eventFired:
		if !reflect.DeepEqual(event, expectedGroupRegistrationEvent) {
			t.Fatalf(
				"Unexpected group registration entry\nExpected: [%v]\nActual:   [%v]",
				expectedGroupRegistrationEvent,
				event,
			)
		}
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}
}

func TestLocalOnGroupRegisteredUnsubscribed(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	chainHandle := Connect(10, 4, big.NewInt(200)).ThresholdRelay()

	eventFired := make(chan *event.GroupRegistration)

	subscription, err := chainHandle.OnGroupRegistered(
		func(entry *event.GroupRegistration) {
			eventFired <- entry
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	subscription.Unsubscribe()

	groupPublicKey := []byte("1")
	signingId := big.NewInt(42)
	memberIndex := group.MemberIndex(1)
	dkgResult := &relaychain.DKGResult{GroupPublicKey: groupPublicKey}
	signatures := make(map[group.MemberIndex]operator.Signature)

	chainHandle.SubmitDKGResult(signingId, memberIndex, dkgResult, signatures)

	select {
	case event := <-eventFired:
		t.Fatalf("Event should have not been received due to the cancelled subscription: [%v]", event)
	case <-ctx.Done():
		// expected execution of goroutine
	}
}

func TestLocalOnDKGResultSubmitted(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	chainHandle := Connect(10, 4, big.NewInt(200)).ThresholdRelay()

	eventFired := make(chan *event.DKGResultSubmission)

	subscription, err := chainHandle.OnDKGResultSubmitted(
		func(request *event.DKGResultSubmission) {
			eventFired <- request
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	defer subscription.Unsubscribe()

	groupPublicKey := []byte("1")
	signingId := big.NewInt(42)
	memberIndex := group.MemberIndex(1)
	dkgResult := &relaychain.DKGResult{GroupPublicKey: groupPublicKey}
	signatures := make(map[group.MemberIndex]operator.Signature)

	chainHandle.SubmitDKGResult(signingId, memberIndex, dkgResult, signatures)

	expectedResultSubmissionEvent := &event.DKGResultSubmission{
		SigningId:      signingId,
		MemberIndex:    uint32(memberIndex),
		GroupPublicKey: groupPublicKey,
	}

	select {
	case event := <-eventFired:
		if !reflect.DeepEqual(event, expectedResultSubmissionEvent) {
			t.Fatalf(
				"Unexpected DKG result submission event\nExpected: [%v]\nActual:   [%v]",
				expectedResultSubmissionEvent,
				event,
			)
		}
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}
}

func TestLocalOnDKGResultSubmittedUnsubscribed(t *testing.T) {
	ctx, cancel := newTestContext()
	defer cancel()

	chainHandle := Connect(10, 4, big.NewInt(200)).ThresholdRelay()

	eventFired := make(chan *event.DKGResultSubmission)

	subscription, err := chainHandle.OnDKGResultSubmitted(
		func(event *event.DKGResultSubmission) {
			eventFired <- event
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	subscription.Unsubscribe()

	groupPublicKey := []byte("1")
	signingId := big.NewInt(42)
	memberIndex := group.MemberIndex(1)
	dkgResult := &relaychain.DKGResult{GroupPublicKey: groupPublicKey}
	signatures := make(map[group.MemberIndex]operator.Signature)

	chainHandle.SubmitDKGResult(signingId, memberIndex, dkgResult, signatures)

	select {
	case event := <-eventFired:
		t.Fatalf("Event should have not been received due to the cancelled subscription: [%v]", event)
	case <-ctx.Done():
		// expected execution of goroutine
	}
}

func TestLocalBlockHeightWaiter(t *testing.T) {
	var tests = map[string]struct {
		blockHeight      uint64
		initialDelay     time.Duration
		expectedWaitTime time.Duration
	}{
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

func TestLocalIsGroupStale(t *testing.T) {
	group1 := localGroup{
		groupPublicKey:          []byte{'v'},
		registrationBlockHeight: 1,
	}

	group2 := localGroup{
		groupPublicKey:          []byte{'i'},
		registrationBlockHeight: 1,
	}

	group3 := localGroup{
		groupPublicKey:          []byte{'d'},
		registrationBlockHeight: 1,
	}

	availableGroups := []localGroup{group1, group2, group3}

	var tests = map[string]struct {
		group           localGroup
		expectedResult  bool
		simulatedHeight uint64
	}{
		"found a first group": {
			group: localGroup{
				groupPublicKey: group1.groupPublicKey,
			},
			simulatedHeight: group1.registrationBlockHeight + 2,
			expectedResult:  false,
		},
		"found a second group": {
			group: localGroup{
				groupPublicKey: group2.groupPublicKey,
			},
			simulatedHeight: group2.registrationBlockHeight + 3,
			expectedResult:  false,
		},
		"group was not found": {
			group: localGroup{
				groupPublicKey: []byte{'z'},
			},
			simulatedHeight: 1,
			expectedResult:  true,
		},
		"a third group was found and current block has passed the expiration and operation timeout": {
			group: localGroup{
				groupPublicKey: group3.groupPublicKey,
			},
			simulatedHeight: group3.registrationBlockHeight +
				groupActiveTime +
				relayRequestTimeout +
				1,
			expectedResult: true,
		},
		"a second group was found and current block is the same as an active time and operation timeout": {
			group: localGroup{
				groupPublicKey: group2.groupPublicKey,
			},
			simulatedHeight: group2.registrationBlockHeight +
				groupActiveTime +
				relayRequestTimeout,
			expectedResult: false,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			localChain := &localChain{
				groups:          availableGroups,
				simulatedHeight: test.simulatedHeight,
			}
			chainHandle := localChain.ThresholdRelay()
			actualResult, err := chainHandle.IsStaleGroup(test.group.groupPublicKey)
			if err != nil {
				t.Fatal(err)
			}

			if actualResult != test.expectedResult {
				t.Fatalf("\nCheck for a group removal eligibility failed. \nexpected: %v\nactual:   %v\n", test.expectedResult, actualResult)
			}
		})
	}
}

func TestLocalIsDKGResultSubmitted(t *testing.T) {
	submittedResults := make(map[*big.Int][]*relaychain.DKGResult)

	submittedSigningId := big.NewInt(1)
	submittedResult := &relaychain.DKGResult{
		GroupPublicKey: []byte{11},
	}

	submittedResults[submittedSigningId] = append(
		submittedResults[submittedSigningId],
		submittedResult,
	)

	chainHandle := Connect(10, 4, big.NewInt(100)).ThresholdRelay()

	chainHandle.SubmitDKGResult(
		submittedSigningId,
		group.MemberIndex(1),
		submittedResult,
		make(map[group.MemberIndex]operator.Signature),
	)

	var tests = map[string]struct {
		signingId      *big.Int
		expectedResult bool
	}{
		"result for the request ID submitted": {
			signingId:      submittedSigningId,
			expectedResult: true,
		},
		"result for the given request ID not yet submitted": {
			signingId:      big.NewInt(3),
			expectedResult: false,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			actualResult, err := chainHandle.IsDKGResultSubmitted(test.signingId)
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
	signingId1 := big.NewInt(1)
	memberIndex := uint32(1)
	submittedResult11 := &relaychain.DKGResult{
		GroupPublicKey: []byte{11},
	}
	expectedEvent1 := &event.DKGResultSubmission{
		SigningId:      signingId1,
		MemberIndex:    memberIndex,
		GroupPublicKey: submittedResult11.GroupPublicKey[:],
		BlockNumber:    0,
	}

	signatures := map[group.MemberIndex]operator.Signature{
		1: operator.Signature{101},
		2: operator.Signature{102},
		3: operator.Signature{103},
	}

	chainHandle.SubmitDKGResult(signingId1, 1, submittedResult11, signatures)
	if !reflect.DeepEqual(
		localChain.submittedResults[signingId1.String()],
		submittedResult11,
	) {
		t.Fatalf("invalid submitted results for request ID %v\nexpected: %v\nactual:   %v\n",
			signingId1,
			[]*relaychain.DKGResult{submittedResult11},
			localChain.submittedResults[signingId1.String()],
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
	signingId2 := big.NewInt(2)
	expectedEvent2 := &event.DKGResultSubmission{
		SigningId:      signingId2,
		MemberIndex:    memberIndex,
		GroupPublicKey: submittedResult11.GroupPublicKey[:],
	}

	chainHandle.SubmitDKGResult(signingId2, 1, submittedResult11, signatures)
	if !reflect.DeepEqual(
		localChain.submittedResults[signingId2.String()],
		submittedResult11,
	) {
		t.Fatalf("invalid submitted results for request ID %v\nexpected: %v\nactual:   %v\n",
			signingId2,
			[]*relaychain.DKGResult{submittedResult11},
			localChain.submittedResults[signingId2.String()],
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
	promise := chainHandle.SubmitDKGResult(signingId1, 1, submittedResult11, signatures)
	promise.OnSuccess(func(result *event.DKGResultSubmission) {
		t.Fatalf("Should not be able to submit result for the given ID more than once")
	})
	promise.OnFailure(func(err error) {
		expectedError := fmt.Errorf("result for request ID [1] already submitted")
		if !reflect.DeepEqual(err, expectedError) {
			t.Fatalf(
				"Unexpected error\nExpected: [%v]\nActual:   [%v]\n",
				expectedError,
				err,
			)
		}
	})
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

func TestNextGroupIndex(t *testing.T) {
	var tests = map[string]struct {
		previousEntry  int
		numberOfGroups int
		expectedIndex  int
	}{
		"zero groups": {
			previousEntry:  12,
			numberOfGroups: 0,
			expectedIndex:  0,
		},
		"fewer groups than the previous entry value": {
			previousEntry:  13,
			numberOfGroups: 4,
			expectedIndex:  1,
		},
		"more groups than the previous entry value": {
			previousEntry:  3,
			numberOfGroups: 12,
			expectedIndex:  3,
		},
	}

	for nextGroupIndexTest, test := range tests {
		t.Run(nextGroupIndexTest, func(t *testing.T) {
			bigPreviousEntry := big.NewInt(int64(test.previousEntry))
			bigNumberOfGroups := test.numberOfGroups
			expectedIndex := test.expectedIndex

			actualIndex := selectGroup(bigPreviousEntry, bigNumberOfGroups)

			if actualIndex != expectedIndex {
				t.Fatalf(
					"Unexpected group index selected\nexpected: [%v]\nactual:   [%v]\n",
					expectedIndex,
					actualIndex,
				)
			}
		})
	}
}
