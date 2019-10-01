package groupselection

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/btcsuite/btcd/btcec"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/gen/async"
	"github.com/keep-network/keep-core/pkg/internal/byteutils"
	"github.com/keep-network/keep-core/pkg/subscription"
)

func TestGenerateTickets(t *testing.T) {
	minimumStake := big.NewInt(1)
	availableStake := big.NewInt(10)
	virtualStakers := availableStake.Int64() / minimumStake.Int64()

	stakingPublicKey, err := newTestPublicKey()
	if err != nil {
		t.Fatal(err)
	}
	stakingPublicKeyECDSA := stakingPublicKey.ToECDSA()
	stakingAddress := crypto.PubkeyToAddress(*stakingPublicKeyECDSA)
	previousBeaconOutput := []byte("test beacon output")

	tickets, err := generateTickets(
		previousBeaconOutput,
		stakingAddress.Bytes(),
		availableStake,
		minimumStake,
	)
	if err != nil {
		t.Fatal(err)
	}

	// We should have 10 tickets
	if len(tickets) != int(virtualStakers) {
		t.Fatalf(
			"expected [%d] tickets, received [%d] tickets",
			virtualStakers,
			len(tickets),
		)
	}

	for i, ticket := range tickets {
		expectedIndex := int64(i + 1)
		// Tickets should be sorted in ascending order
		if expectedIndex != ticket.proof.virtualStakerIndex.Int64() {
			t.Fatalf(
				"got index [%d], want index [%d]",
				ticket.proof.virtualStakerIndex,
				expectedIndex,
			)
		}

		if ticket.proof.virtualStakerIndex == big.NewInt(0) {
			t.Fatal("virutal stakers should be 1-indexed, not 0-indexed")
		}
	}
}

func newTestPublicKey() (*btcec.PublicKey, error) {
	ecdsaPrivateKey, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return nil, fmt.Errorf(
			"could not generate new ephemeral keypair [%v]",
			err,
		)
	}

	return ecdsaPrivateKey.PubKey(), nil
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

func (mgi *mockGroupInterface) GetSelectedParticipants() ([]chain.StakerAddress, error) {
	panic("unexpected")
}

func (mgi *mockGroupInterface) OnGroupSelectionStarted(
	func(groupSelectionStart *event.GroupSelectionStart),
) (subscription.EventSubscription, error) {
	panic("not implemented")
}
