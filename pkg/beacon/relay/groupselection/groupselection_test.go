package groupselection

import (
	"encoding/binary"
	"math/big"
	"reflect"
	"sort"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/config"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/chain/local"
	"github.com/keep-network/keep-core/pkg/gen/async"
	"github.com/keep-network/keep-core/pkg/subscription"
)

func TestSubmission(t *testing.T) {
	var tests = map[string]struct {
		groupSize                 int
		initialSubmissionTickets  []*ticket
		reactiveSubmissionTickets []*ticket
		expectedSubmittedTickets  []uint64
	}{
		// Client has the same number of tickets below the natural threshold
		// (initial submission tickets) as the group size.
		// All initial submission tickets should be submitted to the chain.
		// Reactive ticket submission should not be executed.
		"only initial submission - the same number of tickets as group size": {
			groupSize: 4,
			initialSubmissionTickets: []*ticket{
				newTestTicket(1, 1001),
				newTestTicket(2, 1002),
				newTestTicket(3, 1003),
				newTestTicket(4, 1004),
			},
			reactiveSubmissionTickets: []*ticket{
				newTestTicket(5, 1005),
				newTestTicket(6, 1006),
			},
			expectedSubmittedTickets: []uint64{1001, 1002, 1003, 1004},
		},
		// Client has more tickets below the natural threshold (initial
		// submission tickets) than the group size. Only #group_size of initial
		// submission tickets should be submitted to the chain.
		// Reactive ticket submission should not be executed.
		"only initial submission - more tickets than group size": {
			groupSize: 2,
			initialSubmissionTickets: []*ticket{
				newTestTicket(1, 1001),
				newTestTicket(2, 1002),
				newTestTicket(3, 1003),
				newTestTicket(4, 1004),
			},
			reactiveSubmissionTickets: []*ticket{
				newTestTicket(5, 1005),
				newTestTicket(6, 1006),
			},
			expectedSubmittedTickets: []uint64{1001, 1002},
		},
		// Client has less tickets than the group size. Two tickets have value
		// below the natural threshold, and there are no tickets with value
		// above the natural threshold. All tickets should be submitted to
		// the chain.
		"only initial submission - less tickets than the group size": {
			groupSize: 5,
			initialSubmissionTickets: []*ticket{
				newTestTicket(1, 1001),
				newTestTicket(2, 1002),
			},
			reactiveSubmissionTickets: []*ticket{},
			expectedSubmittedTickets:  []uint64{1001, 1002},
		},
		// Client has less tickets below the natural threshold (initial
		// submission tickets) than the group size. Since no one else submitted
		// their tickets, reactive ticket submission should be executed.
		"with reactive submission phase - the same number of tickets as group size": {
			groupSize: 6,
			initialSubmissionTickets: []*ticket{
				newTestTicket(1, 1001),
				newTestTicket(2, 1002),
				newTestTicket(3, 1003),
				newTestTicket(4, 1004),
			},
			reactiveSubmissionTickets: []*ticket{
				newTestTicket(5, 1005),
				newTestTicket(6, 1006),
			},
			expectedSubmittedTickets: []uint64{1001, 1002, 1003, 1004, 1005, 1006},
		},
		// Client has less tickets below the natural threshold (initial
		// submission tickets) than the group size. Since no one else submitted
		// their tickets, reactive ticket submission should be executed.
		// No more tickets should be submitted by the client at overall than the
		// #group_size, though.
		"with reactive submission phase - more tickets than group size": {
			groupSize: 5,
			initialSubmissionTickets: []*ticket{
				newTestTicket(1, 1001),
				newTestTicket(2, 1002),
				newTestTicket(3, 1003),
				newTestTicket(4, 1004),
			},
			reactiveSubmissionTickets: []*ticket{
				newTestTicket(5, 1005),
				newTestTicket(6, 1006),
				newTestTicket(7, 1007),
				newTestTicket(8, 1008),
			},
			expectedSubmittedTickets: []uint64{1001, 1002, 1003, 1004, 1005},
		},
		// Client has no tickets below the natural threshold (initial
		// submission tickets). Since no one else submitted their tickets,
		// reactive ticket submission should be executed.
		// No more tickets should be submitted by the client at overall than the
		// #group_size, though.
		"with reactive submission phase - no initial submission tickets": {
			groupSize:                3,
			initialSubmissionTickets: []*ticket{},
			reactiveSubmissionTickets: []*ticket{
				newTestTicket(5, 1005),
				newTestTicket(6, 1006),
				newTestTicket(7, 1007),
				newTestTicket(8, 1008),
			},
			expectedSubmittedTickets: []uint64{1005, 1006, 1007},
		},
		// Client has less tickets than the group size. One ticket has value
		// below the natural threshold, the other ticket has value above
		// the natural threshold. Both those tickets should be submitted to
		// the chain.
		"with reactive submission phase - less tickets than the group size": {
			groupSize: 5,
			initialSubmissionTickets: []*ticket{
				newTestTicket(1, 1001),
			},
			reactiveSubmissionTickets: []*ticket{
				newTestTicket(2, 1002),
			},
			expectedSubmittedTickets: []uint64{1001, 1002},
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			chainConfig := &config.Chain{
				GroupSize: test.groupSize,
			}

			chain := &stubGroupInterface{
				groupSize: test.groupSize,
			}

			blockCounter, err := local.BlockCounter()
			if err != nil {
				t.Fatal(err)
			}

			onGroupSelected := func(result *Result) {}

			err = startTicketSubmission(
				test.initialSubmissionTickets,
				test.reactiveSubmissionTickets,
				chain,
				blockCounter,
				chainConfig,
				0, // start block height
				onGroupSelected,
				2,
			)
			if err != nil {
				t.Fatal(err)
			}

			ticketSubmissionTimeout, err := chain.TicketSubmissionTimeout()
			if err != nil {
				t.Fatal(err)
			}

			err = blockCounter.WaitForBlockHeight(
				ticketSubmissionTimeout.Uint64(),
			)
			if err != nil {
				t.Fatal(err)
			}

			submittedTickets, err := chain.GetSubmittedTickets()
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(test.expectedSubmittedTickets, submittedTickets) {
				t.Fatalf(
					"unexpected submitted tickets\nexpected: [%v]\nactual:   [%v]",
					test.expectedSubmittedTickets,
					submittedTickets,
				)
			}
		})
	}
}

func TestRoundCandidateTickets(t *testing.T) {
	groupSize := 9
	reactiveSubmissionRounds := uint64(7)

	tickets := []*ticket{
		newTestTicket(1, 36028797018963968),
		newTestTicket(2, 72057594037927936),
		newTestTicket(3, 144115188075855872),
		newTestTicket(4, 288230376151711744),
		newTestTicket(5, 576460752303423488),
		newTestTicket(6, 1152921504606846976),
		newTestTicket(7, 2305843009213693952),
		newTestTicket(8, 4611686018427387904),
		newTestTicket(9, 9223372036854775808),
	}

	var tests = map[string]struct {
		existingChainTickets             []uint64
		expectedCandidateTicketsPerRound map[uint64][]*ticket
		expectedChainTickets             []uint64
	}{
		"no existing chain tickets - all tickets should be submitted": {
			existingChainTickets: []uint64{},
			expectedCandidateTicketsPerRound: map[uint64][]*ticket{
				0: {tickets[0], tickets[1]},
				1: {tickets[2]},
				2: {tickets[3]},
				3: {tickets[4]},
				4: {tickets[5]},
				5: {tickets[6]},
				6: {tickets[7]},
				7: {tickets[8]},
			},
			expectedChainTickets: []uint64{
				36028797018963968,
				72057594037927936,
				144115188075855872,
				288230376151711744,
				576460752303423488,
				1152921504606846976,
				2305843009213693952,
				4611686018427387904,
				9223372036854775808,
			},
		},
		"better chain tickets exists and their count is below the group size - " +
			"only best tickets should be submitted": {
			existingChainTickets: []uint64{1000, 1001, 1002, 1003},
			expectedCandidateTicketsPerRound: map[uint64][]*ticket{
				0: {tickets[0], tickets[1]},
				1: {tickets[2]},
				2: {tickets[3]},
				3: {tickets[4]},
				4: {},
				5: {},
				6: {},
				7: {},
			},
			expectedChainTickets: []uint64{
				1000,
				1001,
				1002,
				1003,
				36028797018963968,
				72057594037927936,
				144115188075855872,
				288230376151711744,
				576460752303423488,
			},
		},
		"better chain tickets exists and their count is equal the group size - " +
			"no tickets should be submitted": {
			existingChainTickets: []uint64{
				1000, 1001, 1002, 1003, 1004, 1005, 1006, 1007, 1008,
			},
			expectedCandidateTicketsPerRound: map[uint64][]*ticket{
				0: {}, 1: {}, 2: {}, 3: {}, 4: {}, 5: {}, 6: {}, 7: {},
			},
			expectedChainTickets: []uint64{
				1000, 1001, 1002, 1003, 1004, 1005, 1006, 1007, 1008,
			},
		},
		"worse chain tickets exists and their count is below the group size - " +
			"all tickets should be submitted": {
			existingChainTickets: []uint64{
				9223372036854775809,
				9223372036854775810,
				9223372036854775811,
				9223372036854775812,
			},
			expectedCandidateTicketsPerRound: map[uint64][]*ticket{
				0: {tickets[0], tickets[1]},
				1: {tickets[2]},
				2: {tickets[3]},
				3: {tickets[4]},
				4: {tickets[5]},
				5: {tickets[6]},
				6: {tickets[7]},
				7: {tickets[8]},
			},
			expectedChainTickets: []uint64{
				36028797018963968,
				72057594037927936,
				144115188075855872,
				288230376151711744,
				576460752303423488,
				1152921504606846976,
				2305843009213693952,
				4611686018427387904,
				9223372036854775808,
			},
		},
		"worse chain tickets exists and their count is equal the group size - " +
			"all tickets should be submitted": {
			existingChainTickets: []uint64{
				9223372036854775809,
				9223372036854775810,
				9223372036854775811,
				9223372036854775812,
				9223372036854775813,
				9223372036854775814,
				9223372036854775815,
				9223372036854775816,
				9223372036854775817,
			},
			expectedCandidateTicketsPerRound: map[uint64][]*ticket{
				0: {tickets[0], tickets[1]},
				1: {tickets[2]},
				2: {tickets[3]},
				3: {tickets[4]},
				4: {tickets[5]},
				5: {tickets[6]},
				6: {tickets[7]},
				7: {tickets[8]},
			},
			expectedChainTickets: []uint64{
				36028797018963968,
				72057594037927936,
				144115188075855872,
				288230376151711744,
				576460752303423488,
				1152921504606846976,
				2305843009213693952,
				4611686018427387904,
				9223372036854775808,
			},
		},
		"better and worse chain tickets exists and their count is below the group size - " +
			"only best tickets should be submitted": {
			existingChainTickets: []uint64{
				1000,
				1001,
				1002,
				9223372036854775809,
				9223372036854775810,
				9223372036854775811,
			},
			expectedCandidateTicketsPerRound: map[uint64][]*ticket{
				0: {tickets[0], tickets[1]},
				1: {tickets[2]},
				2: {tickets[3]},
				3: {tickets[4]},
				4: {tickets[5]},
				5: {},
				6: {},
				7: {},
			},
			expectedChainTickets: []uint64{
				1000,
				1001,
				1002,
				36028797018963968,
				72057594037927936,
				144115188075855872,
				288230376151711744,
				576460752303423488,
				1152921504606846976,
			},
		},
		"better and worse chain tickets exists and their count is equal the group size - " +
			"only best tickets should be submitted": {
			existingChainTickets: []uint64{
				1000,
				1001,
				1002,
				1003,
				9223372036854775809,
				9223372036854775810,
				9223372036854775811,
				9223372036854775812,
				9223372036854775813,
			},
			expectedCandidateTicketsPerRound: map[uint64][]*ticket{
				0: {tickets[0], tickets[1]},
				1: {tickets[2]},
				2: {tickets[3]},
				3: {tickets[4]},
				4: {},
				5: {},
				6: {},
				7: {},
			},
			expectedChainTickets: []uint64{
				1000,
				1001,
				1002,
				1003,
				36028797018963968,
				72057594037927936,
				144115188075855872,
				288230376151711744,
				576460752303423488,
			},
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			existingChainTickets := make([]*chain.Ticket, 0)
			for _, existingChainTicket := range test.existingChainTickets {
				chainTicket, err := toChainTicket(
					newTestTicket(0, existingChainTicket),
				)
				if err != nil {
					t.Fatal(err)
				}

				existingChainTickets = append(existingChainTickets, chainTicket)
			}

			relayChain := &stubGroupInterface{
				groupSize:        groupSize,
				submittedTickets: existingChainTickets,
			}

			for roundIndex := uint64(0); roundIndex <= reactiveSubmissionRounds; roundIndex++ {
				roundLeadingZeros := reactiveSubmissionRounds - roundIndex

				candidateTickets, err := roundCandidateTickets(
					relayChain,
					tickets,
					roundIndex,
					roundLeadingZeros,
					groupSize,
				)
				if err != nil {
					t.Fatal(err)
				}

				if !reflect.DeepEqual(
					test.expectedCandidateTicketsPerRound[roundIndex],
					candidateTickets,
				) {
					t.Fatalf(
						"unexpected candidate tickets for round [%v]\n"+
							"expected: [%v]\nactual:   [%v]",
						roundIndex,
						test.expectedCandidateTicketsPerRound[roundIndex],
						candidateTickets,
					)
				}

				for _, ticket := range candidateTickets {
					chainTicket, err := toChainTicket(ticket)
					if err != nil {
						t.Fatal(err)
					}

					relayChain.SubmitTicket(chainTicket)
				}
			}

			chainTickets, err := relayChain.GetSubmittedTickets()
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(test.expectedChainTickets, chainTickets) {
				t.Fatalf(
					"unexpected chain tickets\n"+
						"expected: [%v]\nactual:   [%v]",
					test.expectedChainTickets,
					chainTickets,
				)
			}
		})
	}
}

type stubGroupInterface struct {
	groupSize        int
	submittedTickets []*chain.Ticket
}

func (stg *stubGroupInterface) SubmitTicket(ticket *chain.Ticket) *async.EventGroupTicketSubmissionPromise {
	promise := &async.EventGroupTicketSubmissionPromise{}

	stg.submittedTickets = append(stg.submittedTickets, ticket)

	sort.SliceStable(stg.submittedTickets, func(i, j int) bool {
		return stg.submittedTickets[i].Value.Cmp(stg.submittedTickets[j].Value) == -1
	})

	if len(stg.submittedTickets) > stg.groupSize {
		stg.submittedTickets = stg.submittedTickets[:stg.groupSize]
	}

	_ = promise.Fulfill(&event.GroupTicketSubmission{
		TicketValue: ticket.Value,
		BlockNumber: 222,
	})

	return promise
}

func (stg *stubGroupInterface) GetSubmittedTickets() ([]uint64, error) {
	tickets := make([]uint64, len(stg.submittedTickets))

	for i := range tickets {
		valueBytes := common.LeftPadBytes(stg.submittedTickets[i].Value.Bytes(), 32)
		tickets[i] = binary.BigEndian.Uint64(valueBytes)
	}

	return tickets, nil
}

func (stg *stubGroupInterface) TicketSubmissionTimeout() (*big.Int, error) {
	return big.NewInt(10), nil
}

func (stg *stubGroupInterface) GetSelectedParticipants() ([]chain.StakerAddress, error) {
	selected := make([]chain.StakerAddress, stg.groupSize)
	for i := 0; i < stg.groupSize; i++ {
		selected[i] = []byte("whatever")
	}

	return selected, nil
}

func (stg *stubGroupInterface) OnGroupSelectionStarted(
	func(groupSelectionStart *event.GroupSelectionStart),
) (subscription.EventSubscription, error) {
	panic("not implemented")
}
