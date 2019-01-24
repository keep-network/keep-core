package relay

import (
	"bytes"
	"fmt"
	"log"
	"math/big"
	"sync"
	"time"

	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/groupselection"
	"github.com/keep-network/keep-core/pkg/chain"
)

// getTicketListInterval is the number of seconds we wait before requesting the
// ordered ticket list (to run ticket verification)from the chain.
const getTicketListInterval = 5 * time.Second

type groupCandidate struct {
	address []byte
	tickets []*groupselection.Ticket

	selectedTickets     []*groupselection.Ticket
	selectedTicketsLock sync.Mutex
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
	groupSize int,
) error {
	availableStake, err := n.Staker.Stake()
	if err != nil {
		return err
	}
	availableStake = big.NewInt(200000000)

	fmt.Printf(
		"input value for generate tickets: beaconValue [%+v], stakerID: [%+v], availableStake [%+v], min stake [%+v]\n",
		beaconValue, string(n.Staker.ID()), availableStake, n.chainConfig.MinimumStake,
	)
	tickets, err :=
		groupselection.GenerateTickets(
			beaconValue,
			n.Staker.ID(),
			availableStake,
			n.chainConfig.MinimumStake,
		)
	if err != nil {
		return err
	}

	fmt.Printf("Generated [%d] tickets [%+v]\n", len(tickets), tickets)

	var (
		errCh                = make(chan error, len(tickets))
		quitTicketSubmission = make(chan struct{}, 0)
		quitTicketChallenge  = make(chan struct{}, 0)

		groupCandidate = &groupCandidate{
			address: n.Staker.ID(),
			tickets: tickets,
		}
	)

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

	fmt.Println("attempting to submit tickets")
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
				"Error during ticket submission: [%v].\n",
				err,
			)
		case <-submissionTimeout:
			quitTicketSubmission <- struct{}{}
			fmt.Println("submission timeout end")
		case <-challengeTimeout:
			quitTicketChallenge <- struct{}{}
			fmt.Println("challenge timeout end")

			selectedTickets, err := relayChain.GetOrderedTickets()
			if err != nil {
				fmt.Printf(
					"error getting submitted tickets [%v].\n",
					err,
				)
			}

			if len(selectedTickets) == 0 {
				groupCandidate.selectedTicketsLock.Lock()
				selectedTickets = groupCandidate.selectedTickets
				groupCandidate.selectedTicketsLock.Unlock()
			}

			if len(selectedTickets) == 0 {
				fmt.Println("no tickets selected to the group")
				return nil
			}

			groupSelectedTickets := selectedTickets[0:groupSize]

			// Read the selected, ordered tickets from the chain,
			// determine if we're eligible for the next group.
			go n.JoinGroupIfEligible(
				relayChain,
				&groupselection.Result{SelectedTickets: groupSelectedTickets},
				entryRequestID,
				entrySeed,
			)

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
	log.Printf("in ticket submit loop, naturalThreshold [%+v]\n", naturalThreshold)
	for _, ticket := range gc.tickets {
		log.Printf("submitting ticket [%+v]\n", ticket)
		relayChain.SubmitTicket(ticket).OnFailure(
			func(err error) { errCh <- err },
		)
	}

	select {
	case <-quit:
		return
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
			selectedTickets, err := relayChain.GetOrderedTickets()
			if err != nil {
				fmt.Printf(
					"error getting submitted tickets [%v].\n",
					err,
				)
			}

			gc.selectedTicketsLock.Lock()
			gc.selectedTickets = selectedTickets
			gc.selectedTicketsLock.Unlock()

			for _, ticket := range selectedTickets {
				if !costlyCheck(beaconValue, ticket) {
					challenge := &groupselection.TicketChallenge{
						Ticket:        ticket,
						SenderAddress: gc.address,
					}
					relayChain.SubmitChallenge(challenge).OnFailure(
						func(err error) {
							fmt.Printf(
								"Failed to submit challenge with err: [%v]\n",
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

	fmt.Printf("incorrect ticket: [%v] vs [%v]\n", computedValue, ticket.Value)
	return false
}
