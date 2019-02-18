package relay

import (
	"encoding/hex"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/config"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/beacon/relay/groupselection"
	"github.com/keep-network/keep-core/pkg/gen/async"
)

func TestSubmitAllTickets(t *testing.T) {
	// 2^257 is bigger than any SHA256 generated number. We want all tickets to
	// be accepted
	naturalThreshold := new(big.Int).Exp(big.NewInt(2), big.NewInt(257), nil)

	beaconOutput := big.NewInt(10).Bytes()
	stakerValue := []byte("StakerValue1001")

	tickets := []*groupselection.Ticket{
		groupselection.NewTicket(beaconOutput, stakerValue, big.NewInt(1)),
		groupselection.NewTicket(beaconOutput, stakerValue, big.NewInt(2)),
		groupselection.NewTicket(beaconOutput, stakerValue, big.NewInt(3)),
		groupselection.NewTicket(beaconOutput, stakerValue, big.NewInt(4)),
	}

	candidate := &Node{
		chainConfig: &config.Chain{
			NaturalThreshold: naturalThreshold,
		},
	}

	errCh := make(chan error, len(tickets))
	quit := make(chan struct{}, 0)
	submittedTickets := make([]*chain.Ticket, 0)

	mockInterface := &mockGroupInterface{
		mockSubmitTicketFn: func(t *chain.Ticket) *async.GroupTicketPromise {
			submittedTickets = append(submittedTickets, t)
			promise := &async.GroupTicketPromise{}
			promise.Fulfill(&event.GroupTicketSubmission{TicketValue: t.Value})
			return promise
		},
	}

	candidate.submitTickets(tickets, mockInterface, quit, errCh)

	if len(tickets) != len(submittedTickets) {
		t.Errorf(
			"unexpected number of tickets submitted\nexpected: [%v]\nactual: [%v]",
			len(tickets),
			len(submittedTickets),
		)
	}

	for i, ticket := range tickets {
		submitted, err := fromChainTicket(submittedTickets[i])
		if err != nil {
			t.Fatal(err)
		}

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

func TestCancelTicketSubmissionAfterATimeout(t *testing.T) {
	// 2^257 is bigger than any SHA256 generated number. We want all tickets to
	// be accepted
	naturalThreshold := new(big.Int).Exp(big.NewInt(2), big.NewInt(257), nil)

	beaconOutput := big.NewInt(10).Bytes()
	stakerValue := []byte("StakerValue1001")

	tickets := []*groupselection.Ticket{
		groupselection.NewTicket(beaconOutput, stakerValue, big.NewInt(1)),
		groupselection.NewTicket(beaconOutput, stakerValue, big.NewInt(2)),
		groupselection.NewTicket(beaconOutput, stakerValue, big.NewInt(3)),
		groupselection.NewTicket(beaconOutput, stakerValue, big.NewInt(4)),
		groupselection.NewTicket(beaconOutput, stakerValue, big.NewInt(5)),
		groupselection.NewTicket(beaconOutput, stakerValue, big.NewInt(6)),
	}

	candidate := &Node{
		chainConfig: &config.Chain{
			NaturalThreshold: naturalThreshold,
		},
	}

	errCh := make(chan error, len(tickets))
	quit := make(chan struct{}, 0)
	submittedTickets := make([]*chain.Ticket, 0)

	mockInterface := &mockGroupInterface{
		mockSubmitTicketFn: func(t *chain.Ticket) *async.GroupTicketPromise {
			submittedTickets = append(submittedTickets, t)
			promise := &async.GroupTicketPromise{}

			time.Sleep(500 * time.Millisecond)

			promise.Fulfill(&event.GroupTicketSubmission{TicketValue: t.Value})
			return promise
		},
	}

	go func() {
		time.Sleep(1 * time.Second)
		quit <- struct{}{}
	}()

	candidate.submitTickets(tickets, mockInterface, quit, errCh)

	if len(submittedTickets) == 0 {
		t.Errorf("no tickets submitted")
	}

	if len(tickets) == len(submittedTickets) {
		t.Errorf("ticket submission has not been cancelled")
	}
}

func TestToFromChainTicket(t *testing.T) {
	tests := map[string]struct {
		shaValue      string
		expectedError error
	}{
		"ticket sha value": {
			shaValue: "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
		},
		"ticket sha value starts with zeros": {
			shaValue: "00de2289dfca6b3f9034688598756c996eda3e29eb665240d137248610b4137e",
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			hash, err := hex.DecodeString(test.shaValue)
			if err != nil {
				t.Fatal(err)
			}
			var array [32]byte
			copy(array[:], hash[:])
			ticketValue := groupselection.SHAValue(array)

			ticket := &groupselection.Ticket{
				Value: ticketValue,
				Proof: &groupselection.Proof{
					StakerValue:        []byte("staker-value"),
					VirtualStakerIndex: big.NewInt(123),
				},
			}

			chainTicket, err := toChainTicket(ticket)
			if err != nil {
				t.Fatal(err)
			}

			actualTicket, err := fromChainTicket(chainTicket)

			if !reflect.DeepEqual(err, test.expectedError) {
				t.Fatalf("\nexpected: [%+v]\nactual: [%+v]", test.expectedError, err)
			}

			if !reflect.DeepEqual(actualTicket, ticket) {
				t.Fatalf(
					"\nexpected: [%+v]\nactual:   [%+v]",
					ticket,
					actualTicket,
				)
			}
		})
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

func (mgi *mockGroupInterface) GetSelectedParticipants() ([]chain.StakerAddress, error) {
	panic("unexpected")
}
