package groupselection

import (
	"fmt"
	"math/big"
	"sort"
)

var startingIndex = big.NewInt(1)

// GenerateTickets generates a set of tickets for the given staker and relay
// entry value given the specified minimum stake. Returns the resulting
// tickets in sorted order, or an error if there were issues computing the
// tickets.
func GenerateTickets(
	minimumStake *big.Int,
	stakerID string, // S_j id
	availableStake *big.Int, // S_j stake
	entryValue []byte, // V_i
) ([]*Ticket, error) {
	stakingWeight := &big.Int{} // W_j
	stakingWeight = stakingWeight.Quo(availableStake, minimumStake)

	stakerValue := &big.Int{} // Q_j
	stakerValue, ok := stakerValue.SetString(stakerID, 16)
	if !ok {
		return nil, fmt.Errorf(
			"staker ID [%v] failed to parse as hex string",
			stakerID,
		)
	}

	tickets := make(tickets, 0)
	for virtualStaker := startingIndex.Int64(); virtualStaker <= stakingWeight.Int64(); virtualStaker++ {
		tickets = append(
			tickets,
			calculateTicket(entryValue, stakerValue, big.NewInt(virtualStaker)), // prf
		)
	}
	sort.Stable(tickets)

	return tickets, nil
}
