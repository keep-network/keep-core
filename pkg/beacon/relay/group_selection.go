package relay

import (
	"fmt"

	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/groupselection"
	"github.com/keep-network/keep-core/pkg/chain"
)

func (n *Node) SubmitTicketsForGroupSelection(
	entryValue []byte,
	relayChain relaychain.GroupInterface,
	blockCounter chain.BlockCounter,
) error {
	initialTimeout, err := blockCounter.BlockWaiter(
		n.chainConfig.TicketTimeout,
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
			entryValue,
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
				entryValue,
				err,
			)
		case <-initialTimeout:
			return nil
		}
	}

	return nil
}
