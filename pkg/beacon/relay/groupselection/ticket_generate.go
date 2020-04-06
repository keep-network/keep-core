package groupselection

import (
	"math/big"
	"sort"
)

// generateTickets generates a set of tickets for the given staker and relay
// entry value given the specified stake parameters and natural threshold.
//
// Tickets are returned sorted in ascending order by their value.
func generateTickets(
	beaconValue []byte, // V_i
	stakerValue []byte, // Q_j
	availableStake *big.Int, // S_j
	minimumStake *big.Int,
) ([]*ticket, error) {
	stakingWeight := new(big.Int).Quo(availableStake, minimumStake) // W_j

	tickets := make([]*ticket, 0)
	for virtualStaker := int64(1); virtualStaker <= stakingWeight.Int64(); virtualStaker++ {
		ticket, err := newTicket(beaconValue, stakerValue, big.NewInt(virtualStaker))
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, ticket)
	}

	sort.Stable(byValue(tickets))

	return tickets, nil
}
