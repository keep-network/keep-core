package relay

import (
	"bytes"
	"fmt"
	"math/big"
	"time"

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
	quitTicketSubmission := make(chan struct{}, 0)
	quitTicketChallenge := make(chan struct{}, 0)
	groupCandidate := &groupCandidate{address: n.Staker.ID(), tickets: tickets}

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

	// Phase 2a: submit all tickets that fall under the natural threshold
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
			selectedTickets := relayChain.GetOrderedTickets()

			// Read the selected, ordered tickets from the chain,
			// determine if we're eligible for the next group.
			go n.JoinGroupIfEligible(
				relayChain,
				&groupselection.Result{SelectedTickets: selectedTickets},
				entryRequestID,
				entrySeed,
			)
			quitTicketChallenge <- struct{}{}
			return nil
		}
	}
}

// submitTickets checks to see if the submission period is over in between ticket
// submits.
func (gc *groupCandidate) submitTickets(
	relayChain relaychain.GroupInterface,
	naturalThreshold *big.Int,
	quit <-chan struct{},
	errCh chan<- error,
) {
	for _, ticket := range gc.tickets {
		relayChain.SubmitTicket(ticket).OnFailure(
			func(err error) { errCh <- err },
		)

		select {
		case <-quit:
			return
		}
	}
}

func (gc *groupCandidate) verifyTicket(
	relayChain relaychain.GroupInterface,
	beaconValue []byte,
	quit <-chan struct{},
) {
	t := time.NewTimer(1)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			for _, ticket := range relayChain.GetOrderedTickets() {
				if !costlyCheck(beaconValue, ticket) {
					challenge := &groupselection.TicketChallenge{
						Ticket:        ticket,
						SenderAddress: gc.address,
					}
					relayChain.SubmitChallenge(challenge).OnFailure(
						func(err error) {
							fmt.Printf(
								"Failed to submit challenge with err: [%v]",
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
