package ecdsa

import (
	"fmt"
	"math/rand"
	"sort"

	"github.com/keep-network/keep-core/pkg/chain"
)

type byAddress []chain.Address

func (ba byAddress) Len() int           { return len(ba) }
func (ba byAddress) Swap(i, j int)      { ba[i], ba[j] = ba[j], ba[i] }
func (ba byAddress) Less(i, j int) bool { return ba[i] < ba[j] }

func EvaluateRetryParticipantsForSigning(
	groupMembers []chain.Address,
	seed int64,
	retryCount uint,
	quantity uint,
) ([]chain.Address, error) {
	if int(quantity) > len(groupMembers) {
		return nil, fmt.Errorf(
			"Asked for too many seats. %d seats were requested, but there are only %d available.",
			quantity,
			len(groupMembers),
		)
	}
	operatorToSeatCount := make(map[chain.Address]uint)
	for _, operator := range groupMembers {
		operatorToSeatCount[operator]++
	}
	rand.Seed(seed + int64(retryCount))

	operators := make(byAddress, len(operatorToSeatCount))
	i := 0
	for operator := range operatorToSeatCount {
		operators[i] = operator
		i++
	}
	sort.Sort(operators)
	rand.Shuffle(len(operators), func(i, j int) {
		operators[i], operators[j] = operators[j], operators[i]
	})

	seatCount := uint(0)
	acceptedOperators := make(map[chain.Address]bool)
	for j := 0; seatCount < quantity; j++ {
		operator := operators[j]
		seatCount += operatorToSeatCount[operator]
		acceptedOperators[operator] = true
	}

	var seats []chain.Address
	for _, operator := range groupMembers {
		if acceptedOperators[operator] {
			seats = append(seats, operator)
		}
	}
	return seats, nil
}
