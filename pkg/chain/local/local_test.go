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

	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
)

func TestSubmitTicketAndGetSelectedParticipants(t *testing.T) {
	groupSize := 4

	generateTicket := func(index int64) *relaychain.Ticket {
		var value [8]byte
		copy(value[:], common.LeftPadBytes(big.NewInt(10*index).Bytes(), 8))

		return &relaychain.Ticket{
			Value: value,
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

	relayEntryPromise := chainHandle.SubmitRelayEntry(big.NewInt(19).Bytes())

	done := make(chan *event.EntrySubmitted)
	relayEntryPromise.OnSuccess(func(entry *event.EntrySubmitted) {
		done <- entry
	}).OnFailure(func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	})

	select {
	case <-done:
		// expected
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}

}

func TestLocalOnEntrySubmitted(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	chainHandle := Connect(10, 4, big.NewInt(200)).ThresholdRelay()

	eventFired := make(chan *event.EntrySubmitted)

	subscription, err := chainHandle.OnRelayEntrySubmitted(
		func(entry *event.EntrySubmitted) {
			eventFired <- entry
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	defer subscription.Unsubscribe()

	chainHandle.SubmitRelayEntry(big.NewInt(20).Bytes())

	select {
	case <-eventFired:
		// expected
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}
}

func TestLocalOnEntrySubmittedUnsubscribed(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	chainHandle := Connect(10, 4, big.NewInt(200)).ThresholdRelay()

	eventFired := make(chan *event.EntrySubmitted)

	subscription, err := chainHandle.OnRelayEntrySubmitted(
		func(entry *event.EntrySubmitted) {
			eventFired <- entry
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	subscription.Unsubscribe()

	chainHandle.SubmitRelayEntry(big.NewInt(1).Bytes())

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
	memberIndex := relaychain.GroupMemberIndex(1)
	dkgResult := &relaychain.DKGResult{GroupPublicKey: groupPublicKey}
	signatures := map[relaychain.GroupMemberIndex][]byte{
		1: {101},
		2: {102},
		3: {103},
		4: {104},
	}

	chainHandle.SubmitDKGResult(memberIndex, dkgResult, signatures)

	expectedGroupRegistrationEvent := &event.GroupRegistration{
		GroupPublicKey: groupPublicKey,
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
	memberIndex := relaychain.GroupMemberIndex(1)
	dkgResult := &relaychain.DKGResult{GroupPublicKey: groupPublicKey}
	signatures := map[relaychain.GroupMemberIndex][]byte{
		1: {101},
		2: {102},
		3: {103},
		4: {104},
	}

	chainHandle.SubmitDKGResult(memberIndex, dkgResult, signatures)

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
	memberIndex := relaychain.GroupMemberIndex(1)
	dkgResult := &relaychain.DKGResult{GroupPublicKey: groupPublicKey}
	signatures := map[relaychain.GroupMemberIndex][]byte{
		1: {101},
		2: {102},
		3: {103},
		4: {104},
	}

	chainHandle.SubmitDKGResult(memberIndex, dkgResult, signatures)

	expectedResultSubmissionEvent := &event.DKGResultSubmission{
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
	memberIndex := relaychain.GroupMemberIndex(1)
	dkgResult := &relaychain.DKGResult{GroupPublicKey: groupPublicKey}
	signatures := map[relaychain.GroupMemberIndex][]byte{
		1: {101},
		2: {102},
		3: {103},
		4: {104},
	}

	chainHandle.SubmitDKGResult(memberIndex, dkgResult, signatures)

	select {
	case event := <-eventFired:
		t.Fatalf("Event should have not been received due to the cancelled subscription: [%v]", event)
	case <-ctx.Done():
		// expected execution of goroutine
	}
}

func TestWatchBlocks(t *testing.T) {
	c := Connect(10, 4, big.NewInt(100))
	blockCounter, err := c.BlockCounter()
	if err != nil {
		t.Fatal(err)
	}

	ctx1, cancel1 := context.WithCancel(context.Background())
	defer cancel1()

	ctx2, cancel2 := context.WithCancel(context.Background())
	defer cancel2()

	watcher1 := blockCounter.WatchBlocks(ctx1)
	watcher2 := blockCounter.WatchBlocks(ctx2)

	watcher1ReceivedCount := 0
	watcher2ReceivedCount := 0
	go func() {
		for range watcher1 {
			watcher1ReceivedCount++
		}
	}()
	go func() {
		for range watcher2 {
			watcher2ReceivedCount++
		}
	}()

	time.Sleep(600 * time.Millisecond)
	cancel1()
	time.Sleep(600 * time.Millisecond)
	cancel2()

	if watcher1ReceivedCount != 1 {
		t.Errorf("watcher 1 should receive [1] block, has [%v]", watcher1ReceivedCount)
	}
	if watcher2ReceivedCount != 2 {
		t.Errorf("watcher 2 should receive [2] block, has [%v]", watcher2ReceivedCount)
	}
}

func TestWatchBlocksNonBlocking(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1100*time.Millisecond)
	defer cancel()

	c := Connect(10, 4, big.NewInt(100))
	blockCounter, err := c.BlockCounter()
	if err != nil {
		t.Fatal(err)
	}

	_ = blockCounter.WatchBlocks(ctx)        // does not read blocks
	watcher := blockCounter.WatchBlocks(ctx) // does read blocks

	var receivedCount uint64
	go func() {
		for range watcher {
			receivedCount++
		}
	}()

	<-ctx.Done()

	if receivedCount != 2 {
		t.Errorf("watcher should receive [2] blocks, has [%v]", receivedCount)
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

func TestLocalSubmitDKGResult(t *testing.T) {
	localChain := Connect(10, 4, big.NewInt(200))

	chainHandle := localChain.ThresholdRelay()

	memberIndex := relaychain.GroupMemberIndex(1)
	result := &relaychain.DKGResult{
		GroupPublicKey: []byte{11},
	}

	signatures := map[relaychain.GroupMemberIndex][]byte{
		1: {101},
		2: {102},
		3: {103},
		4: {104},
	}

	chainHandle.SubmitDKGResult(memberIndex, result, signatures)
	groupRegistered, err := chainHandle.IsGroupRegistered(result.GroupPublicKey)
	if err != nil {
		t.Fatal(err)
	}

	if !groupRegistered {
		t.Fatalf("Group not registered")
	}
}

func TestLocalSubmitDKGResultWithSignatures(t *testing.T) {
	groupSize := 5
	honestThreshold := 3

	localChain := Connect(groupSize, honestThreshold, big.NewInt(200)).(*localChain)
	chainHandle := localChain.ThresholdRelay()

	var tests = map[string]struct {
		signatures    map[relaychain.GroupMemberIndex][]byte
		expectedError error
	}{
		"no signatures": {
			signatures:    map[relaychain.GroupMemberIndex][]byte{},
			expectedError: fmt.Errorf("failed to submit result with [0] signatures for honest threshold [%v]", honestThreshold),
		},
		"one signature": {
			signatures: map[relaychain.GroupMemberIndex][]byte{
				1: {101},
			},
			expectedError: fmt.Errorf("failed to submit result with [1] signatures for honest threshold [%v]", honestThreshold),
		},
		"one less signature than threshold": {
			signatures: map[relaychain.GroupMemberIndex][]byte{
				1: {101},
				2: {102},
			},
			expectedError: fmt.Errorf("failed to submit result with [2] signatures for honest threshold [%v]", honestThreshold),
		},
		"threshold signatures": {
			signatures: map[relaychain.GroupMemberIndex][]byte{
				1: {101},
				2: {102},
				3: {103},
			},
			expectedError: nil,
		},
		"one more signature than threshold": {
			signatures: map[relaychain.GroupMemberIndex][]byte{
				1: {101},
				2: {102},
				3: {103},
				4: {104},
			},
			expectedError: nil,
		},
		"signatures from all group members": {
			signatures: map[relaychain.GroupMemberIndex][]byte{
				1: {101},
				2: {102},
				3: {103},
				4: {104},
				5: {105},
			},
			expectedError: nil,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			ctx, cancel := newTestContext()
			defer cancel()

			errorChan := make(chan error)

			memberIndex := uint8(1)
			result := &relaychain.DKGResult{
				GroupPublicKey: []byte{11},
			}

			promise := chainHandle.SubmitDKGResult(
				memberIndex,
				result,
				test.signatures,
			)
			promise.OnComplete(func(event *event.DKGResultSubmission, err error) {
				errorChan <- err
			})

			select {
			case err := <-errorChan:
				if !reflect.DeepEqual(test.expectedError, err) {
					t.Fatalf(
						"Unexpected error\nExpected: [%v]\nActual:   [%v]\n",
						test.expectedError,
						err,
					)
				}

			case <-ctx.Done():
				t.Fatalf("promise timed out")
			}
		})
	}
}

func TestCalculateDKGResultHash(t *testing.T) {
	localChain := &localChain{}

	dkgResult := &relaychain.DKGResult{
		GroupPublicKey: []byte{3, 40, 200},
		Misbehaved:     []byte{1, 2, 8, 14},
	}
	expectedHashString := "97a94a3b11a0f780c9510df852ac7f77072085d5bda4b07e5d198396dd4f68e5"

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
