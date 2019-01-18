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

	n.tickets = tickets

	initialTimeout, err := blockCounter.BlockWaiter(
		n.chainConfig.TicketInitialSubmissionTimeout,
	)
	if err != nil {
		return err
	}

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

	var (
		initialSubmitErrorChannel    = make(chan error, len(tickets))
		reactiveSubmitErrorChannel   = make(chan error, len(tickets))
		quitTicketInitialSubmission  = make(chan struct{}, 1)
		quitTicketReactiveSubmission = make(chan struct{}, 1)
		quitTicketChallenge          = make(chan struct{}, 1)
	)
	// Phase 2a: submit all tickets that fall under the natural threshold
	go n.submitTickets(
		relayChain,
		quitTicketInitialSubmission,
		initialSubmitErrorChannel,
	)

	// kick off background loop to check submitted tickets
	go n.verifyTicket(
		relayChain,
		beaconValue,
		quitTicketChallenge,
	)

	for {
		select {
		case err := <-initialSubmitErrorChannel:
			fmt.Printf(
				"Error during initial ticket submission for entry [%v]: [%v].",
				beaconValue,
				err,
			)
		case <-initialTimeout:
			quitTicketInitialSubmission <- struct{}{}

			// Phase 2b: submit all tickets, even those above the
			// natural threshold.
			go n.submitTicketsReactive(
				relayChain,
				quitTicketReactiveSubmission,
				reactiveSubmitErrorChannel,
			)
		case err := <-reactiveSubmitErrorChannel:
			fmt.Printf(
				"Error during reactive ticket submission for entry [%v]: [%v].",
				beaconValue,
				err,
			)
		case <-submissionTimeout:
			quitTicketReactiveSubmission <- struct{}{}
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

// submitTickets submits tickets to the chain. It checks to see if the submission
// period is over in between ticket submits.
func (n *Node) submitTickets(
	relayChain relaychain.GroupSelectionInterface,
	quit <-chan struct{},
	errCh chan<- error,
) {
	for _, ticket := range n.tickets {
		select {
		case <-quit:
			// Exit this loop when we get a signal from quit.
			return
		default:
			if ticket.Value.Int().Cmp(n.chainConfig.NaturalThreshold) < 0 {
				relayChain.SubmitTicket(ticket).OnSuccess(
					func(submittedTicket *groupselection.Ticket) {
						if submittedTicket == nil {
							return
						}

						n.submittedTicketsMutex.Lock()
						n.submittedTickets[submittedTicket.Value] = true
						n.submittedTicketsMutex.Unlock()
					},
				).OnFailure(
					func(err error) { errCh <- err },
				)
			}
		}
	}
}

// submitTicketsReactive submits tickets to the chain. It checks to see if the submission
// period is over in between ticket submits.
func (n *Node) submitTicketsReactive(
	relayChain relaychain.GroupSelectionInterface,
	quit <-chan struct{},
	errCh chan<- error,
) {
	for _, ticket := range n.tickets {
		select {
		case <-quit:
			// Exit this loop when we get a signal from quit.
			return
		default:
			n.submittedTicketsMutex.Lock()
			defer n.submittedTicketsMutex.Unlock()

			if !n.submittedTickets[ticket.Value] {
				relayChain.SubmitTicket(ticket).OnSuccess(
					func(submittedTicket *groupselection.Ticket) {
						n.submittedTickets[submittedTicket.Value] = true
					},
				).OnFailure(
					func(err error) { errCh <- err },
				)
			}
		}
	}
}

func (n *Node) verifyTicket(
	relayChain relaychain.GroupSelectionInterface,
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
						SenderAddress: n.Staker.ID(),
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
