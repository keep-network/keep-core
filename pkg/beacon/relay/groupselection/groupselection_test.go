package groupselection

import (
	"math/big"
	"testing"

	"github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/config"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/chain/local"
	"github.com/keep-network/keep-core/pkg/gen/async"
	"github.com/keep-network/keep-core/pkg/subscription"
)

func TestSubmission(t *testing.T) {
	var tests = map[string]struct {
		groupSize                     int
		initialSubmissionTickets      []*ticket
		reactiveSubmissionTickets     []*ticket
		expectedSubmittedTicketsCount int
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
			expectedSubmittedTicketsCount: 4,
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
			expectedSubmittedTicketsCount: 2,
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
			reactiveSubmissionTickets:     []*ticket{},
			expectedSubmittedTicketsCount: 2,
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
			expectedSubmittedTicketsCount: 6,
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
			expectedSubmittedTicketsCount: 5,
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
			expectedSubmittedTicketsCount: 3,
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
			expectedSubmittedTicketsCount: 2,
		},
	}

	for _, test := range tests {
		chainConfig := &config.Chain{
			GroupSize:                     test.groupSize,
			TicketSubmissionRoundDuration: 3,
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
		)
		if err != nil {
			t.Fatal(err)
		}

		err = blockCounter.WaitForBlockHeight(
			chainConfig.TicketSubmissionRoundDuration,
		)
		if err != nil {
			t.Fatal(err)
		}

		submittedTickets, err := chain.GetSubmittedTickets()
		if err != nil {
			t.Fatal(err)
		}

		submittedTicketsCount := big.NewInt(int64(len(submittedTickets)))

		expectedCount := big.NewInt(int64(test.expectedSubmittedTicketsCount))
		if expectedCount.Cmp(submittedTicketsCount) != 0 {
			t.Fatalf(
				"unexpected number of submitted tickets\nexpected: [%v]\nactual:   [%v]",
				expectedCount,
				submittedTicketsCount,
			)
		}
	}
}

type stubGroupInterface struct {
	groupSize        int
	submittedTickets []*chain.Ticket
}

func (stg *stubGroupInterface) SubmitTicket(ticket *chain.Ticket) *async.EventGroupTicketSubmissionPromise {
	stg.submittedTickets = append(stg.submittedTickets, ticket)
	promise := &async.EventGroupTicketSubmissionPromise{}

	promise.Fulfill(&event.GroupTicketSubmission{
		TicketValue: ticket.Value,
		BlockNumber: 222,
	})

	return promise
}

func (stg *stubGroupInterface) GetSubmittedTickets() ([]uint64, error) {
	tickets := make([]uint64, len(stg.submittedTickets))

	for i := range tickets {
		tickets[i] = stg.submittedTickets[i].Value.Uint64()
	}

	return tickets, nil
}

func (stg *stubGroupInterface) TicketSubmissionTimeout() (*big.Int, error) {
	return big.NewInt(6), nil
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
