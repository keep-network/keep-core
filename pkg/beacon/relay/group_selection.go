package relay

import "math/big"

var naturalThreshold = big.NewInt((2 ^ 256) - 1)

func (n *Node) SubmitTicketsForGroupSelection(entryValue []byte) error {
	tickets, err := n.Staker.GenerateTickets(entryValue)
	if err != nil {
		return err
	}

	for _, ticket := range tickets {
		// submit the tickets that fall under the natural threshold
		if ticket.Value.Int().Cmp(naturalThreshold) < 0 {
			// publish the result
		}
	}

	// Ensure this state terminates at TICKET_INITAL_TIMEOUT
	return nil
}
