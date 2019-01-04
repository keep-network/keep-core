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
	minimumStake *big.Int,
	availableStake *big.Int, // S_j stake
	stakerValue []byte, // Q_j
	entryValue []byte, // V_i
) ([]*Ticket, error) {
	stakingWeight := &big.Int{} // W_j
	stakingWeight = stakingWeight.Quo(availableStake, minimumStake)

	tickets := make(tickets, 0)
	for virtualStaker := one; virtualStaker <= stakingWeight.Int64(); virtualStaker++ {
		tickets = append(
			tickets,
			calculateTicket(entryValue, stakerValue, big.NewInt(virtualStaker)), // prf
		)
	}
	sort.Stable(tickets)

	return tickets, nil
}
