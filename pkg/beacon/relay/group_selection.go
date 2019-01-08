package relay

import (
	"fmt"
	"math/big"
	"time"

	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/groupselection"
	"github.com/keep-network/keep-core/pkg/chain"
)

func (n *Node) SubmitTicketsForGroupSelection(
	beaconValue []byte,
	relayChain relaychain.GroupInterface,
	blockCounter chain.BlockCounter,
) error {
	submissionTimeout, err := blockCounter.BlockWaiter(
		n.chainConfig.TicketSubmissionTimeout,
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
			n.chainConfig.MinimumStake,
			availableStake,
			[]byte(n.Staker.ID()),
			beaconValue,
		)
	if err != nil {
		return err
	}

	errCh := make(chan error, len(tickets))
	quitTicketSubmission := make(chan struct{}, 0)
	quitTicketChallenge := make(chan struct{}, 0)

	// Phase 2a: submit all tickets that fall under the natural threshold
	go submitTickets(
		relayChain,
		tickets,
		n.chainConfig.NaturalThreshold,
		quitTicketSubmission,
		errCh,
	)

	// kick off background loop to check submitted tickets
	go verifyTicket(relayChain, n.StakeID, quitTicketChallenge)

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
			quitTicketChallenge <- struct{}{}
			return nil
		}
	}
}

func submitTickets(
	relayChain relaychain.GroupInterface,
	tickets []*groupselection.Ticket,
	naturalThreshold *big.Int,
	quit <-chan struct{},
	errCh chan error,
) {
	for _, ticket := range tickets {
		if ticket.Value.Int().Cmp(naturalThreshold) < 0 {
			relayChain.SubmitTicket(ticket).OnFailure(func(err error) {
				errCh <- err
			})
		}

		select {
		case <-quit:
			return
		}
	}
}

func verifyTicket(
	relayChain relaychain.GroupInterface,
	self string,
	quit <-chan struct{},
) {
	t := time.NewTimer(1)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			for _, ticket := range relayChain.GetOrderedTickets() {
				if !costlyCheck(ticket) {
					challenge := &groupselection.Challenge{
						Ticket: ticket,
						Sender: self,
					}
					relayChain.SubmitChallenge(challenge).OnFailure(
						func(err error) {
							// TODO: implement
						},
					)
				}
			}
			t.Reset(5 * time.Second)
		case <-quit:
			// Exit this loop when we get a signal from quit.
			return
		}
	}
}

func costlyCheck(ticket *groupselection.Ticket) bool {
	// TODO: implement
	return true
}
