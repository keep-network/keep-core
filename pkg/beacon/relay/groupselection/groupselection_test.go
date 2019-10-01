package groupselection

import (
	"bytes"
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

	tickets, err := GenerateTickets(
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
		if expectedIndex != ticket.Proof.VirtualStakerIndex.Int64() {
			t.Fatalf("Got index [%d], want index [%d]",
				ticket.Proof.VirtualStakerIndex,
				expectedIndex,
			)
		}

		if ticket.Proof.VirtualStakerIndex == big.NewInt(0) {
			t.Fatal("Virutal stakers should be 1-indexed, not 0-indexed")
		}
	}

}

func TestValidateProofs(t *testing.T) {
	beaconOutput := []byte("test beacon output")
	beaconOutputPadded, _ := byteutils.LeftPadTo32Bytes(beaconOutput)

	stakingPublicKey, err := newTestPublicKey()
	if err != nil {
		t.Fatal(err)
	}
	stakingPublicKeyECDSA := stakingPublicKey.ToECDSA()
	stakingAddress := crypto.PubkeyToAddress(*stakingPublicKeyECDSA)
	stakerValuePadded, _ := byteutils.LeftPadTo32Bytes(stakingAddress.Bytes())

	minimumStake := big.NewInt(1)
	availableStake := big.NewInt(1)
	virtualStakers := big.NewInt(0).Quo(availableStake, minimumStake) // 1
	virtualStakerIndexPadded, _ := byteutils.LeftPadTo32Bytes(virtualStakers.Bytes())

	var valueBytes []byte

	valueBytes = append(valueBytes, beaconOutputPadded...) // V_i
	valueBytes = append(valueBytes, stakerValuePadded...)  // Q_j
	// only 1 virtual staker, which corresponds to the index, vs
	valueBytes = append(valueBytes, virtualStakerIndexPadded...)

	expectedValue := crypto.Keccak256(valueBytes[:])

	tickets, err := GenerateTickets(
		beaconOutput,
		stakingAddress.Bytes(),
		availableStake,
		minimumStake,
	)
	if err != nil {
		t.Fatal(err)
	}

	// we should have virtualStaker number of tickets
	if len(tickets) != int(virtualStakers.Int64()) {
		t.Fatalf(
			"expected [%d] tickets, received [%d] tickets",
			virtualStakers,
			len(tickets),
		)
	}

	if bytes.Compare(
		tickets[0].Value.Bytes(),
		expectedValue,
	) != 0 {
		t.Fatalf(
			"hashed value (%v) doesn't match ticket value (%v)",
			tickets[0].Value,
			expectedValue,
		)
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

	tickets := make([]*Ticket, 0)
	for i := 1; i <= 4; i++ {
		ticket, _ := NewTicket(beaconOutput, stakerValue, big.NewInt(int64(i)))
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

func fromChainTicket(ticket *chain.Ticket, t *testing.T) *Ticket {
	paddedTicketValue, err := byteutils.LeftPadTo32Bytes((ticket.Value.Bytes()))
	if err != nil {
		t.Errorf("could not pad ticket value [%v]", err)
	}

	value, err := SHAValue{}.SetBytes(paddedTicketValue)
	if err != nil {
		t.Errorf(
			"could not transform ticket from chain representation [%v]",
			err,
		)
	}

	return &Ticket{
		Value: value,
		Proof: &Proof{
			StakerValue:        ticket.Proof.StakerValue.Bytes(),
			VirtualStakerIndex: ticket.Proof.VirtualStakerIndex,
		},
	}
}

func TestCancelTicketSubmissionAfterATimeout(t *testing.T) {
	beaconOutput := big.NewInt(10).Bytes()
	stakerValue := []byte("StakerValue1001")

	tickets := make([]*Ticket, 0)
	for i := 1; i <= 6; i++ {
		ticket, _ := NewTicket(beaconOutput, stakerValue, big.NewInt(int64(i)))
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
