package relay

import (
	"fmt"
	"math/big"

	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/groupselection"
	"github.com/keep-network/keep-core/pkg/chain"
)

const ticketInitialTimeout = 5

var minimumStake = big.NewInt(1)
var naturalThreshold = big.NewInt((2 ^ 256) - 1)

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
			minimumStake,
			availableStake,
			[]byte(n.Staker.ID()),
			entryValue,
		)
	if err != nil {
		return err
	}

	errCh := make(chan error, len(tickets))
	for _, ticket := range tickets {
		if ticket.Value.Int().Cmp(naturalThreshold) < 0 {
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
