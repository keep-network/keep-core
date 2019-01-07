package relay

import (
	"fmt"

	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/groupselection"
	"github.com/keep-network/keep-core/pkg/chain"
)

func (n *Node) SubmitTicketsForGroupSelection(
	beaconValue []byte,
	relayChain relaychain.GroupInterface,
	blockCounter chain.BlockCounter,
) error {
	// Timeout for initial ticket submission, Phase 2a
	initialTimeout, err := blockCounter.BlockWaiter(
		n.chainConfig.TicketInitialTimeout,
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
	for _, ticket := range tickets {
		if ticket.Value.Int().Cmp(n.chainConfig.NaturalThreshold) < 0 {
			relayChain.SubmitTicket(ticket).OnFailure(func(err error) {
				errCh <- err
			})
		}
	}

	for {
		select {
		case err := <-errCh:
			fmt.Printf(
				"Error during ticket submission for entry [%v]: [%v].",
				beaconValue,
				err,
			)
		case <-initialTimeout:
			// Phase 2b: reactive ticket submission. Submit all
			// tickets that are above the natural threshold.
			// get the current best threshold

			return nil
		}
	}

	return nil
}
