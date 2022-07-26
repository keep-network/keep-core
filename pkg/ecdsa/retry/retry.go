package retry

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

func calculateSeatCount(groupMembers []chain.Address) map[chain.Address]uint {
	operatorToSeatCount := make(map[chain.Address]uint)
	for _, operator := range groupMembers {
		operatorToSeatCount[operator]++
	}
	return operatorToSeatCount
}

// EvaluateRetryParticipantsForSigning takes in a slice of `groupMembers` and
// returns a subslice of those same members of length >= `quantity` randomly
// according to the provided `seed` and `retryCount`.
//
// This function is intended to be called during a signing protocol after a
// signing event has failed but *not* due to inactivity. Assuming that some of
// the `groupMembers` are sending corrupted information, either on purpose or
// accidentally, we keep trying to find a subset of `groupMembers` that is as
// small as possible, yet still larger than `quantity`.
//
// The `seed` param needs to vary on a per-message basis but must be the same
// seed between all operators for each invocation. This can be the hash of the
// message since cryptographically secure randomness isn't important.
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
	operatorToSeatCount := calculateSeatCount(groupMembers)
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
