package local

import (
	"context"
	"fmt"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/keep-network/keep-core/pkg/beacon/relay/chain"
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
	dkgSubmissions := make(map[string]*relaychain.DKGSubmissions)

	requestID := big.NewInt(1)
	dkgSubmission := &relaychain.DKGSubmission{
		DKGResult: &relaychain.DKGResult{
			GroupPublicKey: []byte{11},
		},
		Votes: 1,
	}

	dkgSubmissions[requestID.String()] = &relaychain.DKGSubmissions{
		DKGSubmissions: []*relaychain.DKGSubmission{
			dkgSubmission,
		},
	}

	localChain := &localChain{
		submissions: dkgSubmissions,
	}
	chainHandle := localChain.ThresholdRelay()

	var tests = map[string]struct {
		requestID      *big.Int
		expectedResult bool
	}{
		"matched": {
			requestID:      requestID,
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
	localChain := &localChain{
		submissions:        make(map[string]*relaychain.DKGSubmissions),
		submissionHandlers: make(map[int]func(dkgResultPublication *event.DKGResultPublication)),
	}
	chainHandle := localChain.ThresholdRelay()

	// Channel for DKGResultPublication events.
	dkgResultPublicationChan := make(chan *event.DKGResultPublication)
	defer close(dkgResultPublicationChan)
	subscription, err := localChain.OnDKGResultPublished(
		func(dkgResultPublication *event.DKGResultPublication) {
			dkgResultPublicationChan <- dkgResultPublication
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	defer subscription.Unsubscribe()

	if len(localChain.submissions) > 0 {
		t.Fatalf("initial submitted results map is not empty")
	}

	// Submit new result for request ID 1
	requestID1 := big.NewInt(1)
	dkgResult1 := &relaychain.DKGResult{
		GroupPublicKey: []byte{11},
	}
	expectedEvent1 := &event.DKGResultPublication{
		RequestID:      requestID1,
		GroupPublicKey: dkgResult1.GroupPublicKey,
	}
	expectedsubmissions1 := &relaychain.DKGSubmissions{
		DKGSubmissions: []*relaychain.DKGSubmission{
			&relaychain.DKGSubmission{
				DKGResult: dkgResult1,
				Votes:     1,
			},
		},
	}

	chainHandle.SubmitDKGResult(requestID1, dkgResult1).
		OnFailure(func(err error) {
			t.Fatal(err)
		})
	if !reflect.DeepEqual(
		expectedsubmissions1,
		localChain.submissions[requestID1.String()],
	) {
		t.Fatalf("\nexpected: %+v\nactual:   %+v\n",
			expectedsubmissions1,
			localChain.submissions[requestID1.String()],
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
		GroupPublicKey: dkgResult1.GroupPublicKey,
	}
	expectedSubmissions2 := &relaychain.DKGSubmissions{
		DKGSubmissions: []*relaychain.DKGSubmission{
			&relaychain.DKGSubmission{
				DKGResult: dkgResult1,
				Votes:     1,
			},
		},
	}

	chainHandle.SubmitDKGResult(requestID2, dkgResult1).
		OnFailure(func(err error) {
			t.Fatal(err)
		})
	if !reflect.DeepEqual(
		expectedSubmissions2,
		localChain.submissions[requestID2.String()],
	) {
		t.Fatalf("\nexpected: %v\nactual:   %v\n",
			expectedSubmissions2,
			localChain.submissions[requestID2.String()],
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
	expectedError := fmt.Errorf("result already submitted")

	chainHandle.SubmitDKGResult(requestID1, dkgResult1).
		OnFailure(func(err error) {
			if !reflect.DeepEqual(expectedError, err) {
				t.Fatalf("\nexpected: %v\nactual:   %v\n",
					expectedError,
					err,
				)
			}
		})

	if !reflect.DeepEqual(
		expectedsubmissions1,
		localChain.submissions[requestID1.String()],
	) {
		t.Fatalf("\nexpected: %v\nactual:   %v\n",
			expectedsubmissions1,
			localChain.submissions[requestID1.String()],
		)
	}
	select {
	case dkgResultPublicationEvent := <-dkgResultPublicationChan:
		t.Fatalf("unexpected event was emitted: %v", dkgResultPublicationEvent)
	case <-time.After(500 * time.Millisecond):
		break
	}

	// Submit new result for request ID 1
	dkgResult2 := &relaychain.DKGResult{
		GroupPublicKey: []byte{12},
	}
	expectedEvent3 := &event.DKGResultPublication{
		RequestID:      requestID1,
		GroupPublicKey: dkgResult2.GroupPublicKey,
	}
	expectedSubmissions3 := &relaychain.DKGSubmissions{
		DKGSubmissions: []*relaychain.DKGSubmission{
			&relaychain.DKGSubmission{
				DKGResult: dkgResult1,
				Votes:     1,
			},
			&relaychain.DKGSubmission{
				DKGResult: dkgResult2,
				Votes:     1,
			},
		},
	}

	chainHandle.SubmitDKGResult(requestID1, dkgResult2).
		OnFailure(func(err error) {
			t.Fatal(err)
		})
	if !reflect.DeepEqual(
		expectedSubmissions3,
		localChain.submissions[requestID1.String()],
	) {
		t.Fatalf("\nexpected: %+v\nactual:   %+v\n",
			expectedSubmissions3,
			localChain.submissions[requestID1.String()],
		)
	}

	select {
	case dkgResultPublicationEvent := <-dkgResultPublicationChan:
		if !reflect.DeepEqual(expectedEvent3, dkgResultPublicationEvent) {
			t.Fatalf("\nexpected: %v\nactual:   %v\n",
				expectedEvent3,
				dkgResultPublicationEvent,
			)
		}
	case <-ctx.Done():
		t.Fatalf("expected event was not emitted")
	}
}

func TestLocalOnDKGResultPublishedUnsubscribe(t *testing.T) {
	ctx, cancel := newTestContext()
	defer cancel()

	localChain := &localChain{
		submissions:        make(map[string]*relaychain.DKGSubmissions),
		submissionHandlers: make(map[int]func(dkgResultPublication *event.DKGResultPublication)),
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
		break
	}
}

func TestDKGResultVote(t *testing.T) {
	ctx, cancel := newTestContext()
	defer cancel()

	c := Connect(10, 4, big.NewInt(200)).(*localChain)
	chain := c.ThresholdRelay()

	voteChan := make(chan *event.DKGResultVote)
	defer close(voteChan)

	subscription, err := chain.OnDKGResultVote(
		func(dkgResultVote *event.DKGResultVote) {
			voteChan <- dkgResultVote
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	defer subscription.Unsubscribe()

	dkgResult1 := &relaychain.DKGResult{
		GroupPublicKey: []byte{100},
	}
	dkgResult2 := &relaychain.DKGResult{
		GroupPublicKey: []byte{200},
	}
	requestID0 := big.NewInt(0)
	requestID1 := big.NewInt(1)

	c.submissions = map[string]*relaychain.DKGSubmissions{
		requestID1.String(): &relaychain.DKGSubmissions{
			DKGSubmissions: []*relaychain.DKGSubmission{
				&relaychain.DKGSubmission{
					DKGResult: dkgResult1,
					Votes:     3,
				},
				&relaychain.DKGSubmission{
					DKGResult: dkgResult2,
					Votes:     1,
				},
			},
		},
	}

	tests := map[string]struct {
		requestID           *big.Int
		dkgResult           *relaychain.DKGResult
		expectedSubmissions func(chain *localChain) map[string]*relaychain.DKGSubmissions
		expectedEvent       *event.DKGResultVote
		expectedError       error
	}{
		"invalid vote when no submissions for given request id": {
			requestID: requestID0,
			dkgResult: dkgResult1,
			expectedSubmissions: func(chain *localChain) map[string]*relaychain.DKGSubmissions {
				return chain.submissions
			},
			expectedEvent: nil,
			expectedError: fmt.Errorf("no submissions for given request id"),
		},
		"invalid vote for dkg result not matching submitted one": {
			requestID: requestID1,
			dkgResult: &relaychain.DKGResult{GroupPublicKey: []byte{99}},
			expectedSubmissions: func(chain *localChain) map[string]*relaychain.DKGSubmissions {
				return chain.submissions
			},
			expectedError: fmt.Errorf("no submissions matching given dkg result hash"),
		},
		"valid vote for submitted dkg result": {
			requestID: requestID1,
			dkgResult: dkgResult2,
			expectedSubmissions: func(chain *localChain) map[string]*relaychain.DKGSubmissions {
				chain.submissions[requestID1.String()].DKGSubmissions[1].Votes = 2
				return chain.submissions
			},
			expectedEvent: &event.DKGResultVote{RequestID: requestID1},
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {

			promise := chain.DKGResultVote(test.requestID, test.dkgResult.Hash())

			promise.OnComplete(func(event *event.DKGResultVote, err error) {
				if !reflect.DeepEqual(test.expectedError, err) {
					t.Fatalf("\nexpected: %v\nactual:   %v\n",
						test.expectedError,
						err,
					)
				}
			})

			expectedSubmissions := test.expectedSubmissions(c)
			if !reflect.DeepEqual(
				expectedSubmissions,
				c.submissions,
			) {
				t.Fatalf("\nexpected: %v\nactual:   %v\n",
					expectedSubmissions,
					c.submissions,
				)
			}

			if test.expectedEvent == nil {
				select {
				case dkgResultVoteEvent := <-voteChan:
					t.Fatalf("unexpected event was emitted: %v", dkgResultVoteEvent)
				case <-time.After(100 * time.Millisecond):
					break
				}
			} else {
				select {
				case dkgResultVoteEvent := <-voteChan:
					if !reflect.DeepEqual(test.expectedEvent, dkgResultVoteEvent) {
						t.Fatalf("\nexpected: %v\nactual:   %v\n",
							test.expectedEvent,
							dkgResultVoteEvent,
						)
					}
				case <-ctx.Done():
					t.Fatalf("expected event was not emitted")
				}
			}
		})
	}
}

func TestSubmitAndGetDKGSubmissions(t *testing.T) {
	requestID1 := big.NewInt(1)
	requestID2 := big.NewInt(2)
	dkgResult1 := &chain.DKGResult{
		GroupPublicKey: []byte{100},
	}
	dkgResult2 := &chain.DKGResult{
		GroupPublicKey: []byte{200},
	}

	localChain := &localChain{
		submissionHandlers: make(map[int]func(dkgResultPublication *event.DKGResultPublication)),
		submissions:        make(map[string]*relaychain.DKGSubmissions),
	}

	var tests = map[string]struct {
		requestID           *big.Int
		resultsToSubmit     []*chain.DKGResult
		expectedSubmissions *chain.DKGSubmissions
	}{
		"no submissions": {
			requestID:           big.NewInt(0),
			resultsToSubmit:     []*chain.DKGResult{},
			expectedSubmissions: &chain.DKGSubmissions{},
		},
		"one submission for given request ID": {
			requestID:       requestID1,
			resultsToSubmit: []*chain.DKGResult{dkgResult1},
			expectedSubmissions: &chain.DKGSubmissions{
				DKGSubmissions: []*chain.DKGSubmission{
					&chain.DKGSubmission{
						DKGResult: dkgResult1,
						Votes:     1,
					},
				},
			},
		},
		"two submissions for given request ID": {
			requestID:       requestID2,
			resultsToSubmit: []*chain.DKGResult{dkgResult1, dkgResult2},
			expectedSubmissions: &chain.DKGSubmissions{
				DKGSubmissions: []*relaychain.DKGSubmission{
					&chain.DKGSubmission{
						DKGResult: dkgResult1,
						Votes:     1,
					},
					&chain.DKGSubmission{
						DKGResult: dkgResult2,
						Votes:     1,
					},
				},
			},
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			chainHandle := localChain.ThresholdRelay()
			for _, dkgResult := range test.resultsToSubmit {
				chainHandle.SubmitDKGResult(test.requestID, dkgResult)
			}

			submissions := chainHandle.GetDKGSubmissions(test.requestID)
			if !reflect.DeepEqual(test.expectedSubmissions, submissions) {
				t.Fatalf("\nexpected: %v\nactual:   %v\n",
					test.expectedSubmissions,
					submissions,
				)
			}
		})
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
				submissions:        make(map[string]*relaychain.DKGSubmissions),
				submissionHandlers: make(map[int]func(dkgResultPublication *event.DKGResultPublication)),
				voteHandler:        make(map[int]func(dkgResultVote *event.DKGResultVote)),
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

func newTestContext(timeout ...time.Duration) (context.Context, context.CancelFunc) {
	defaultTimeout := 3 * time.Second
	if len(timeout) > 0 {
		defaultTimeout = timeout[0]
	}
	return context.WithTimeout(context.Background(), defaultTimeout)
}
