package groupselection

import (
	"math/big"

	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
)

// submitTicketsOnChain submits tickets to the chain.
func submitTicketsOnChain(
	tickets []*ticket,
	relayChain relaychain.GroupSelectionInterface,
) {
	for _, ticket := range tickets {
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

func toChainTicket(ticket *ticket) (*relaychain.Ticket, error) {
	return &relaychain.Ticket{
		Value: new(big.Int).SetBytes(ticket.value[:]),
		Proof: &relaychain.TicketProof{
			StakerValue:        new(big.Int).SetBytes(ticket.proof.stakerValue),
			VirtualStakerIndex: ticket.proof.virtualStakerIndex,
		},
	}, nil
}
