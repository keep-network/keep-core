package groupselection

import (
	"math/big"
	"sort"
)

// generateTickets generates a set of tickets for the given staker and relay
// entry value given the specified stake parameters and natural threshold.
// It returns the tickets in two slices:
// - initialSubmissionTickets contains tickets with values below the natural
//   threshold. These tickets should be submitted first as they have the highest
//   chance of being selected to the group.
// - reactiveSubmissionTickets contains tickets with values equal to or above
//   the natural threshold. These tickets should be submitted in the reactive
//   submission phase if there is still a chance to become a group member.
//
// Tickets are returned sorted in ascending order by their value.
func generateTickets(
	beaconValue []byte, // V_i
	stakerValue []byte, // Q_j
	availableStake *big.Int, // S_j
	minimumStake *big.Int,
	naturalThreshold *big.Int,
) (
	initialSubmissionTickets []*ticket,
	reactiveSubmissionTickets []*ticket,
	err error,
) {
	stakingWeight := (&big.Int{}).Quo(availableStake, minimumStake) // W_j

	tickets := make([]*ticket, 0)
	for virtualStaker := int64(1); virtualStaker <= stakingWeight.Int64(); virtualStaker++ {
		ticket, err := newTicket(beaconValue, stakerValue, big.NewInt(virtualStaker))
		if err != nil {
			return nil, nil, err
		}
		tickets = append(tickets, ticket)
	}

	sort.Stable(byValue(tickets))

	for _, ticket := range tickets {
		if ticket.intValue().Cmp(naturalThreshold) < 0 {
			initialSubmissionTickets = append(initialSubmissionTickets, ticket)
		} else {
			reactiveSubmissionTickets = append(reactiveSubmissionTickets, ticket)
		}
	}

	return initialSubmissionTickets, reactiveSubmissionTickets, nil
}
