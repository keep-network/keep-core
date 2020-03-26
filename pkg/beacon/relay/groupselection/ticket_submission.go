package groupselection

import (
	"math/big"

	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
)

// submitTickets submits tickets to the chain. It checks to see if the submission
// period is over in between ticket submits.
func submitTickets(
	tickets []*ticket,
	relayChain relaychain.GroupSelectionInterface,
	quit <-chan struct{},
) {
	for _, ticket := range tickets {
		select {
		case <-quit:
			return
		default:
			chainTicket, err := toChainTicket(ticket)
			if err != nil {
				logger.Errorf(
					"could not transform ticket to chain format: [%v]",
					err,
				)
				continue
			}

			relayChain.SubmitTicket(chainTicket).OnFailure(
				func(err error) {
					logger.Errorf(
						"ticket submission failed: [%v]",
						err,
					)
				},
			)
		}
	}
}

func toChainTicket(ticket *ticket) (*relaychain.Ticket, error) {
	return &relaychain.Ticket{
		Value: new(big.Int).SetBytes(ticket.value[:]),
		Proof: &relaychain.TicketProof{
			StakerValue:        new(big.Int).SetBytes(ticket.proof.stakerValue),
			VirtualStakerIndex: ticket.proof.virtualStakerIndex,
		},
	}, nil
}
