package relay

import (
	"context"
	"fmt"
	"math/big"

	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
	"github.com/keep-network/keep-core/pkg/chain"
)

const TicketInitialTimeout = 5

var NaturalThreshold = big.NewInt((2 ^ 256) - 1)

func (n *Node) SubmitTicketsForGroupSelection(
	entryValue []byte,
	relayChain relaychain.Interface,
	blockCounter chain.BlockCounter,
) error {
	tickets, err := n.Staker.GenerateTickets(entryValue)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, len(tickets))

	go func(
		cancel context.CancelFunc,
		tickets group.Tickets,
		relayChain relaychain.Interface,
		blockCounter chain.BlockCounter,
		errCh chan error,
	) {
		initialTimeout, err := blockCounter.BlockWaiter(
			TicketInitialTimeout,
		)
		if err != nil {
			errCh <- fmt.Errorf(
				"failed to initialize blockCounter with err %v",
				err,
			)
			return
		}
		for _, ticket := range tickets {
			// submit the tickets that fall under the natural threshold
			if ticket.Value.Int().Cmp(NaturalThreshold) < 0 {
				// publish the result
				relayChain.SubmitTicket(
					ticket,
				).OnFailure(func(err error) {
					errCh <- err
				})
			}
			select {
			// if TicketInitialTimeout blocks have passed, close the context
			case <-initialTimeout:
				cancel()
				return
			}
		}
	}(cancel, tickets, relayChain, blockCounter, errCh)

	for {
		select {
		case err := <-errCh:
			// TODO: log this error
			fmt.Println(err)
		case <-ctx.Done():
			close(errCh)
			return nil
		}
	}

	return nil
}
