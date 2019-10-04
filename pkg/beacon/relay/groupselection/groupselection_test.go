package groupselection

import (
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/gen/async"
	"github.com/keep-network/keep-core/pkg/internal/byteutils"
	"github.com/keep-network/keep-core/pkg/subscription"
)

var stakingAddress = []byte("staking address")
var previousBeaconOutput = []byte("test beacon output")

func naturalThreshold() *big.Int { // 2^256 / 2
	return new(big.Int).Div(
		new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil),
		big.NewInt(2),
	)
}

func TestAllTicketsGenerated(t *testing.T) {
	minimumStake := big.NewInt(20)
	availableStake := big.NewInt(1000)
	virtualStakers := availableStake.Int64() / minimumStake.Int64()

	initialTickets, reactiveTickets, err := generateTickets(
		previousBeaconOutput,
		stakingAddress,
		availableStake,
		minimumStake,
		naturalThreshold(),
	)
	if err != nil {
		t.Fatal(err)
	}

	// We should have 1000/20 = 50 tickets
	allTicketsCount := len(initialTickets) + len(reactiveTickets)
	if allTicketsCount != int(virtualStakers) {
		t.Fatalf(
			"expected [%d] tickets, has [%d] tickets",
			virtualStakers,
			allTicketsCount,
		)
	}
}

func TestTicketsGeneratedInOrder(t *testing.T) {
	minimumStake := big.NewInt(1)
	availableStake := big.NewInt(100)

	initialTickets, reactiveTickets, err := generateTickets(
		previousBeaconOutput,
		stakingAddress,
		availableStake,
		minimumStake,
		naturalThreshold(),
	)
	if err != nil {
		t.Fatal(err)
	}

	allTickets := append(initialTickets, reactiveTickets...)

	// Tickets should be sorted in ascending order
	for i := 0; i < len(allTickets)-1; i++ {
		value := allTickets[i].intValue()
		nextValue := allTickets[i+1].intValue()

		if value.Cmp(nextValue) > 0 {
			t.Errorf("tickets not sorted in ascending order")
		}
	}
}

func TestInitialTicketsGeneatedBelowNaturalThreshold(t *testing.T) {
	minimumStake := big.NewInt(1)
	availableStake := big.NewInt(10000)

	initialTickets, _, err := generateTickets(
		previousBeaconOutput,
		stakingAddress,
		availableStake,
		minimumStake,
		naturalThreshold(),
	)
	if err != nil {
		t.Fatal(err)
	}

	// All initial submission tickets should have value below the natural
	// threshold
	for _, ticket := range initialTickets {
		if ticket.intValue().Cmp(naturalThreshold()) >= 0 {
			t.Errorf(
				"initial submission ticket value should be below natural "+
					"threshold\nvalue:     [%v]\nthreshold: [%v]",
				ticket.intValue(),
				naturalThreshold(),
			)
		}
	}
}

func TestReactiveTicketsGeneatedAboveNaturalThreshold(t *testing.T) {
	minimumStake := big.NewInt(1)
	availableStake := big.NewInt(10000)

	_, reactiveTickets, err := generateTickets(
		previousBeaconOutput,
		stakingAddress,
		availableStake,
		minimumStake,
		naturalThreshold(),
	)
	if err != nil {
		t.Fatal(err)
	}

	// All reactive submission tickets should have value equal or above
	// the natural threshold
	for _, ticket := range reactiveTickets {
		if ticket.intValue().Cmp(naturalThreshold()) <= 0 {
			t.Errorf(
				"reactive submission ticket value should not be below natural "+
					"threshold\nvalue:     [%v]\nthreshold: [%v]",
				ticket.intValue(),
				naturalThreshold(),
			)
		}
	}
}

func TestSubmitAllTickets(t *testing.T) {
	beaconOutput := big.NewInt(10).Bytes()
	stakerValue := []byte("StakerValue1001")

	tickets := make([]*ticket, 0)
	for i := 1; i <= 4; i++ {
		ticket, _ := newTicket(beaconOutput, stakerValue, big.NewInt(int64(i)))
		tickets = append(tickets, ticket)
	}

	errCh := make(chan error, len(tickets))
	quit := make(chan struct{}, 0)
	submittedTickets := make([]*chain.Ticket, 0)

	mockInterface := &mockGroupInterface{
		mockSubmitTicketFn: func(t *chain.Ticket) *async.GroupTicketPromise {
			submittedTickets = append(submittedTickets, t)
			promise := &async.GroupTicketPromise{}
			promise.Fulfill(&event.GroupTicketSubmission{
				TicketValue: t.Value,
				BlockNumber: 111,
			})
			return promise
		},
	}

	submitTickets(tickets, mockInterface, quit, errCh)

	if len(tickets) != len(submittedTickets) {
		t.Errorf(
			"unexpected number of tickets submitted\nexpected: [%v]\nactual: [%v]",
			len(tickets),
			len(submittedTickets),
		)
	}

	for i, ticket := range tickets {
		submitted := fromChainTicket(submittedTickets[i], t)

		if !reflect.DeepEqual(ticket, submitted) {
			t.Errorf(
				"unexpected ticket at index [%v]\nexpected: [%v]\nactual: [%v]",
				i,
				ticket,
				submitted,
			)
		}
	}
}

func fromChainTicket(chainTicket *chain.Ticket, t *testing.T) *ticket {
	paddedTicketValue, err := byteutils.LeftPadTo32Bytes((chainTicket.Value.Bytes()))
	if err != nil {
		t.Errorf("could not pad ticket value [%v]", err)
	}

	var value [32]byte
	copy(value[:], paddedTicketValue)

	return &ticket{
		value: value,
		proof: &proof{
			stakerValue:        chainTicket.Proof.StakerValue.Bytes(),
			virtualStakerIndex: chainTicket.Proof.VirtualStakerIndex,
		},
	}
}

func TestCancelTicketSubmissionAfterATimeout(t *testing.T) {
	beaconOutput := big.NewInt(10).Bytes()
	stakerValue := []byte("StakerValue1001")

	tickets := make([]*ticket, 0)
	for i := 1; i <= 6; i++ {
		ticket, _ := newTicket(beaconOutput, stakerValue, big.NewInt(int64(i)))
		tickets = append(tickets, ticket)
	}

	errCh := make(chan error, len(tickets))
	quit := make(chan struct{}, 0)
	submittedTickets := make([]*chain.Ticket, 0)

	mockInterface := &mockGroupInterface{
		mockSubmitTicketFn: func(t *chain.Ticket) *async.GroupTicketPromise {
			submittedTickets = append(submittedTickets, t)
			promise := &async.GroupTicketPromise{}

			time.Sleep(500 * time.Millisecond)

			promise.Fulfill(&event.GroupTicketSubmission{
				TicketValue: t.Value,
				BlockNumber: 222,
			})
			return promise
		},
	}

	go func() {
		time.Sleep(1 * time.Second)
		quit <- struct{}{}
	}()

	submitTickets(tickets, mockInterface, quit, errCh)

	if len(submittedTickets) == 0 {
		t.Errorf("no tickets submitted")
	}

	if len(tickets) == len(submittedTickets) {
		t.Errorf("ticket submission has not been cancelled")
	}
}

type mockGroupInterface struct {
	mockSubmitTicketFn func(t *chain.Ticket) *async.GroupTicketPromise
}

func (mgi *mockGroupInterface) SubmitTicket(
	ticket *chain.Ticket,
) *async.GroupTicketPromise {
	if mgi.mockSubmitTicketFn != nil {
		return mgi.mockSubmitTicketFn(ticket)
	}

	panic("unexpected")
}

func (mgi *mockGroupInterface) GetSubmittedTicketsCount() (*big.Int, error) {
	panic("not implemented")
}

func (mgi *mockGroupInterface) GetSelectedParticipants() ([]chain.StakerAddress, error) {
	panic("unexpected")
}

func (mgi *mockGroupInterface) OnGroupSelectionStarted(
	func(groupSelectionStart *event.GroupSelectionStart),
) (subscription.EventSubscription, error) {
	panic("not implemented")
}
