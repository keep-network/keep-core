package relay

import (
	"bytes"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/groupselection"
	"github.com/keep-network/keep-core/pkg/chain"
)

// getTicketListInterval is the number of seconds we wait before requesting the
// ordered ticket list (to run ticket verification)from the chain.
const getTicketListInterval = 5 * time.Second

type groupCandidate struct {
	address string
	tickets []*groupselection.Ticket
}

// SubmitTicketsForGroupSelection takes the previous beacon value and attempts to
// generate the appropriate number of tickets for the staker. After ticket
// generation begins an interactive process, where the staker submits tickets
// that fall under the natural threshold, while challenging tickets on chain
// that fail verification. Submission ends at the end of the submission period,
// and the staker can only contest incorrect tickets up to the challenge period.
//
// See the group selection protocol specification for more information.
func (n *Node) SubmitTicketsForGroupSelection(
	relayChain relaychain.Interface,
	blockCounter chain.BlockCounter,
	beaconValue []byte,
	entryRequestID *big.Int,
	entrySeed *big.Int,
) error {
	// we use reactive submission timeout temporarily, until we properly
	// implement phases 2a and 2b.
	submissionTimeout, err := blockCounter.BlockWaiter(
		n.chainConfig.TicketReactiveSubmissionTimeout,
	)
	if err != nil {
		return err
	}

	challengeTimeout, err := blockCounter.BlockWaiter(
		n.chainConfig.TicketChallengeTimeout,
	)
	if err != nil {
		return err
	}

	availableStake, err := n.Staker.Stake()
	if err != nil {
		return err
	}

	tickets, err :=
		groupselection.GenerateTickets(
			beaconValue,
			[]byte(n.Staker.ID()),
			availableStake,
			n.chainConfig.MinimumStake,
		)
	if err != nil {
		return err
	}

	errCh := make(chan error, len(tickets))
	quitTicketSubmission := make(chan struct{}, 1)
	quitTicketChallenge := make(chan struct{}, 0)
	groupCandidate := &groupCandidate{address: n.Staker.ID(), tickets: tickets}

	go groupCandidate.submitTickets(
		relayChain,
		n.chainConfig.NaturalThreshold,
		quitTicketSubmission,
		errCh,
	)

	// kick off background loop to check submitted tickets
	go groupCandidate.verifyTicket(
		relayChain,
		beaconValue,
		quitTicketChallenge,
	)

	for {
		select {
		case err := <-errCh:
			fmt.Printf(
				"Error during ticket submission for entry [%v]: [%v].",
				beaconValue,
				err,
			)
		case <-submissionTimeout:
			quitTicketSubmission <- struct{}{}
		case <-challengeTimeout:
			selectedTickets, err := relayChain.GetOrderedTickets()
			if err != nil {
				quitTicketChallenge <- struct{}{}
				return fmt.Errorf(
					"could not fetch ordered tickets after challenge timeout: [%v]",
					err,
				)
			}

			var tickets []*groupselection.Ticket
			for _, chainTicket := range selectedTickets {
				ticket, err := fromChainTicket(chainTicket)
				if err != nil {
					fmt.Fprintf(
						os.Stderr,
						"incorrect ticket format: [%v]",
						err,
					)

					continue // ignore incorrect ticket
				}

				tickets = append(tickets, ticket)
			}

			// Read the selected, ordered tickets from the chain,
			// determine if we're eligible for the next group.
			go n.JoinGroupIfEligible(
				relayChain,
				&groupselection.Result{SelectedTickets: tickets},
				entryRequestID,
				entrySeed,
			)
			quitTicketChallenge <- struct{}{}
			return nil
		}
	}
}

// submitTickets submits tickets to the chain. It checks to see if the submission
// period is over in between ticket submits.
func (gc *groupCandidate) submitTickets(
	relayChain relaychain.GroupSelectionInterface,
	naturalThreshold *big.Int,
	quit <-chan struct{},
	errCh chan<- error,
) {
	for _, ticket := range gc.tickets {
		select {
		case <-quit:
			// Exit this loop when we get a signal from quit.
			return
		default:
			chainTicket, err := toChainTicket(ticket)
			if err != nil {
				errCh <- err
				continue
			}

			relayChain.SubmitTicket(chainTicket).OnFailure(
				func(err error) { errCh <- err },
			)
		}
	}
}

func (gc *groupCandidate) verifyTicket(
	relayChain relaychain.GroupSelectionInterface,
	beaconValue []byte,
	quit <-chan struct{},
) {
	t := time.NewTimer(1)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			selectedTickets, err := relayChain.GetOrderedTickets()
			if err != nil {
				fmt.Fprintf(
					os.Stderr,
					"error getting submitted tickets: [%v]",
					err,
				)
			}

			for _, selectedTicket := range selectedTickets {
				ticket, err := fromChainTicket(selectedTicket)
				if err != nil {
					fmt.Fprintf(
						os.Stderr,
						"incorrect ticket format: [%v]",
						err,
					)

					continue // ignore incorrect ticket
				}

				if !costlyCheck(beaconValue, ticket) {
					relayChain.SubmitChallenge(ticket.Value.Int()).OnFailure(
						func(err error) {
							fmt.Fprintf(
								os.Stderr,
								"Failed to submit challenge: [%v]",
								err,
							)
						},
					)
				}
			}
			t.Reset(getTicketListInterval)
		case <-quit:
			// Exit this loop when we get a signal from quit.
			return
		}
	}
}

// costlyCheck takes the on-chain Proof, computes the sha256 hash from the Proof,
// and then uses a constant time compare to determine if the on-chain value
// matches the value the client computes for them.
func costlyCheck(beaconValue []byte, ticket *groupselection.Ticket) bool {
	// cheapCheck is done on chain
	computedValue := groupselection.CalculateTicketValue(
		beaconValue,
		ticket.Proof.StakerValue,
		ticket.Proof.VirtualStakerIndex,
	)
	switch bytes.Compare(computedValue[:], ticket.Value[:]) {
	case 0:
		return true
	}
	return false
}

func toChainTicket(ticket *groupselection.Ticket) (*relaychain.Ticket, error) {
	stakerValueInt, err := hexutil.DecodeBig(string(ticket.Proof.StakerValue))
	if err != nil {
		return nil, fmt.Errorf(
			"could not transform ticket to chain representation: [%v]",
			err,
		)
	}

	return &relaychain.Ticket{
		Value: ticket.Value.Int(),
		Proof: &relaychain.TicketProof{
			StakerValue:        stakerValueInt,
			VirtualStakerIndex: ticket.Proof.VirtualStakerIndex,
		},
	}, nil

}

func fromChainTicket(ticket *relaychain.Ticket) (*groupselection.Ticket, error) {
	value, err := groupselection.SHAValue{}.SetBytes(ticket.Value.Bytes())
	if err != nil {
		return nil, fmt.Errorf("incorrect ticket value [%v]", err)
	}

	return &groupselection.Ticket{
		Value: value,
		Proof: &groupselection.Proof{
			StakerValue: []byte(
				hexutil.EncodeBig(ticket.Proof.StakerValue),
			),
			VirtualStakerIndex: ticket.Proof.VirtualStakerIndex,
		},
	}, nil
}
