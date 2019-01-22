package relay

import (
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/keep-network/keep-core/pkg/beacon/relay/groupselection"
	"github.com/keep-network/keep-core/pkg/gen/async"
)

func TestSubmitAllTickets(t *testing.T) {
	// 2^257 is bigger than any SHA256 generated number. We want all tickets to
	// be accepted
	naturalThreshold := new(big.Int).Exp(big.NewInt(2), big.NewInt(257), nil)

	beaconOutput := big.NewInt(10).Bytes()

	tickets := []*groupselection.Ticket{
		groupselection.NewTicket(beaconOutput, big.NewInt(11).Bytes(), big.NewInt(1)),
		groupselection.NewTicket(beaconOutput, big.NewInt(12).Bytes(), big.NewInt(2)),
		groupselection.NewTicket(beaconOutput, big.NewInt(13).Bytes(), big.NewInt(3)),
		groupselection.NewTicket(beaconOutput, big.NewInt(14).Bytes(), big.NewInt(4)),
	}

	candidate := &groupCandidate{
		tickets: tickets,
	}

	errCh := make(chan error, len(tickets))
	quit := make(chan struct{}, 0)
	submittedTickets := make([]*groupselection.Ticket, 0)

	mockInterface := &mockGroupInterface{
		mockSubmitTicketFn: func(t *groupselection.Ticket) *async.GroupTicketPromise {
			submittedTickets = append(submittedTickets, t)
			promise := &async.GroupTicketPromise{}
			promise.Fulfill(t)
			return promise
		},
	}

	candidate.submitTickets(mockInterface, naturalThreshold, quit, errCh)

	if !reflect.DeepEqual(tickets, submittedTickets) {
		t.Errorf(
			"unexpected tickets submitted\n[%v]\n[%v]",
			tickets,
			submittedTickets,
		)
	}
}

func TestCancelTicketSubmissionAfterATimeout(t *testing.T) {
	// 2^257 is bigger than any SHA256 generated number. We want all tickets to
	// be accepted
	naturalThreshold := new(big.Int).Exp(big.NewInt(2), big.NewInt(257), nil)

	beaconOutput := big.NewInt(10).Bytes()

	tickets := []*groupselection.Ticket{
		groupselection.NewTicket(beaconOutput, big.NewInt(11).Bytes(), big.NewInt(1)),
		groupselection.NewTicket(beaconOutput, big.NewInt(12).Bytes(), big.NewInt(2)),
		groupselection.NewTicket(beaconOutput, big.NewInt(13).Bytes(), big.NewInt(3)),
		groupselection.NewTicket(beaconOutput, big.NewInt(14).Bytes(), big.NewInt(4)),
		groupselection.NewTicket(beaconOutput, big.NewInt(15).Bytes(), big.NewInt(5)),
		groupselection.NewTicket(beaconOutput, big.NewInt(16).Bytes(), big.NewInt(6)),
	}

	candidate := &groupCandidate{
		tickets: tickets,
	}

	errCh := make(chan error, len(tickets))
	quit := make(chan struct{}, 0)
	submittedTickets := make([]*groupselection.Ticket, 0)

	mockInterface := &mockGroupInterface{
		mockSubmitTicketFn: func(t *groupselection.Ticket) *async.GroupTicketPromise {
			submittedTickets = append(submittedTickets, t)
			promise := &async.GroupTicketPromise{}
			promise.Fulfill(t)

			time.Sleep(500 * time.Millisecond)

			return promise
		},
	}

	go func() {
		time.Sleep(1 * time.Second)
		quit <- struct{}{}
	}()

	candidate.submitTickets(mockInterface, naturalThreshold, quit, errCh)

	if len(tickets) == len(submittedTickets) {
		t.Errorf("ticket submission has not been cancelled")
	}
}

type mockGroupInterface struct {
	mockSubmitTicketFn func(t *groupselection.Ticket) *async.GroupTicketPromise
}

func (mgi *mockGroupInterface) SubmitTicket(
	ticket *groupselection.Ticket,
) *async.GroupTicketPromise {
	if mgi.mockSubmitTicketFn != nil {
		return mgi.mockSubmitTicketFn(ticket)
	}

	panic("unexpected")
}

func (mgi *mockGroupInterface) SubmitChallenge(
	ticket *groupselection.TicketChallenge,
) *async.GroupTicketChallengePromise {
	panic("unexpected")
}

func (mgi *mockGroupInterface) GetOrderedTickets() ([]*groupselection.Ticket, error) {
	panic("unexpected")
}
