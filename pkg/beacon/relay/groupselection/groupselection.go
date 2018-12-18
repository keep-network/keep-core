package groupselection

import (
	"fmt"
	"math/big"
	"sort"

	"github.com/keep-network/keep-core/pkg/staker"
)

var one = (&big.Int{}).SetUint64(uint64(1))

// GenerateTickets generates a set of tickets for the given staker and relay
// entry value given the specified minimum stake. Returns the resulting
// tickets in sorted order, or an error if there were issues computing the
// tickets.
func GenerateTickets(
	minimumStake *big.Int,
	staker staker.Staker, // S_j
	entryValue []byte, // V_i
) ([]*Ticket, error) {
	availableStake, err := staker.Stake()
	if err != nil {
		return nil, err
	}

	stakingWeight := &big.Int{} // W_j
	stakingWeight = stakingWeight.Quo(availableStake, minimumStake)

	stakerValue := &big.Int{} // Q_j
	stakerValue, ok := stakerValue.SetString(staker.ID(), 16)
	if !ok {
		return nil, fmt.Errorf(
			"staker ID [%v] failed to parse as hex string",
			staker.ID(),
		)
	}

	tickets := make(tickets, 0)
	for virtualStaker := (&big.Int{}).Set(one); virtualStaker.Cmp(stakingWeight) <= 0; virtualStaker.Add(virtualStaker, one) {
		tickets = append(
			tickets,
			calculateTicket(entryValue, stakerValue, virtualStaker), // prf
		)
	}
	sort.Stable(tickets)

	return tickets, nil
}
