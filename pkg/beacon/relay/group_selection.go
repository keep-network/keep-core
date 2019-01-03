package relay

import (
	"fmt"
	"math/big"

	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/groupselection"
	"github.com/keep-network/keep-core/pkg/chain"
)

const ticketInitialTimeout = 5

func (n *Node) SubmitTicketsForGroupSelection(
	entryValue []byte,
	relayChain relaychain.GroupInterface,
	blockCounter chain.BlockCounter,
) error {
	initialTimeout, err := blockCounter.BlockWaiter(ticketInitialTimeout)
	if err != nil {
		return err
	}

	availableStake, err := n.Staker.Stake()
	if err != nil {
		return err
	}

	tickets, err :=
		groupselection.GenerateTickets(
			groupselection.MinimumStake,
			availableStake,
			[]byte(n.Staker.ID()),
			entryValue,
		)
	if err != nil {
		return err
	}

	errCh := make(chan error, len(tickets))
	virtualStakers := big.NewInt(int64(len(tickets)))
	for _, ticket := range tickets {
		if ticket.Value.Int().Cmp(groupselection.NaturalThreshold(virtualStakers)) < 0 {
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
