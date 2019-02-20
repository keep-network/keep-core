package groupselection

import (
	"math/big"
	"sort"
)

var one = int64(1)

// GenerateTickets generates a set of tickets for the given staker and relay
// entry value given the specified minimum stake. Returns the resulting
// tickets in sorted order, or an error if there were issues computing the
// tickets.
func GenerateTickets(
	beaconValue []byte, // V_i
	stakerValue []byte, // Q_j
	availableStake *big.Int, // S_j
	minimumStake *big.Int,
) ([]*Ticket, error) {
	stakingWeight := (&big.Int{}).Quo(availableStake, minimumStake) // W_j

	tickets := make(tickets, 0)
	for virtualStaker := one; virtualStaker <= stakingWeight.Int64(); virtualStaker++ {
		tickets = append(
			tickets,
			NewTicket(beaconValue, stakerValue, big.NewInt(virtualStaker)), // prf
		)
	}
	sort.Stable(tickets)

	return tickets, nil
}

// Result represents ordered, selected tickets from those submitted to the chain.
type Result struct {
	SelectedStakers [][]byte
}
