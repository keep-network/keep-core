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
// returns a subslice of those same members of length >=
// `retryParticipantsCount` randomly according to the provided `seed` and
// `retryCount`.
//
// This function is intended to be called during a signing protocol after a
// signing event has failed but *not* due to inactivity. Assuming that some of
// the `groupMembers` are sending corrupted information, either on purpose or
// accidentally, we keep trying to find a subset of `groupMembers` that is as
// small as possible, yet still larger than `retryParticipantsCount`.
//
// The `seed` param needs to vary on a per-message basis but must be the same
// seed between all operators for each invocation. This can be the hash of the
// message since cryptographically secure randomness isn't important.
//
// The `retryCount` denotes the number of the given retry, so that should be
// incremented after each attempt while the `seed` stays consistent on a
// per-message basis.
func EvaluateRetryParticipantsForSigning(
	groupMembers []chain.Address,
	seed int64,
	retryCount uint,
	retryParticipantsCount uint,
) ([]chain.Address, error) {
	if int(retryParticipantsCount) > len(groupMembers) {
		return nil, fmt.Errorf(
			"asked for too many seats. [%d] seats were requested, but there are only [%d] available.",
			retryParticipantsCount,
			len(groupMembers),
		)
	}
	operatorToSeatCount := calculateSeatCount(groupMembers)
	rand.Seed(seed + int64(retryCount))

	operators := make([]chain.Address, len(operatorToSeatCount))
	i := 0
	for operator := range operatorToSeatCount {
		operators[i] = operator
		i++
	}
	sort.Sort(byAddress(operators))
	rand.Shuffle(len(operators), func(i, j int) {
		operators[i], operators[j] = operators[j], operators[i]
	})

	seatCount := uint(0)
	acceptedOperators := make(map[chain.Address]bool)
	for j := 0; seatCount < retryParticipantsCount; j++ {
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

func excludeSingleOperator(
	rng *rand.Rand,
	groupMembers []chain.Address,
	index int,
	operatorToSeatCount map[chain.Address]uint,
	operators []chain.Address,
) ([]chain.Address, int, bool) {
	if index < len(operators) {
		rng.Shuffle(len(operators), func(i, j int) {
			operators[i], operators[j] = operators[j], operators[i]
		})
		removedOperator := operators[index]
		usedOperators := make([]chain.Address, 0, len(groupMembers))
		for _, operator := range groupMembers {
			if operator != removedOperator {
				usedOperators = append(usedOperators, operator)
			}
		}
		return usedOperators, 0, true
	} else {
		return nil, len(operators), false
	}
}

func excludeOperatorPairs(
	rng *rand.Rand,
	groupMembers []chain.Address,
	index int,
	operatorToSeatCount map[chain.Address]uint,
	operators []chain.Address,
	retryParticipantsCount int,
) ([]chain.Address, int, bool) {
	pairIndexes := make([][2]int, 0, len(operators)*len(operators))
	for i := 0; i < len(operators)-1; i++ {
		for j := i + 1; j < len(operators); j++ {
			leftOperator := operators[i]
			rightOperator := operators[j]
			count := len(groupMembers) -
				int(operatorToSeatCount[leftOperator]) -
				int(operatorToSeatCount[rightOperator])

			if count >= int(retryParticipantsCount) {
				pairIndexes = append(pairIndexes, [2]int{i, j})
			}
		}
	}
	if index < len(pairIndexes) {
		rng.Shuffle(len(pairIndexes), func(i, j int) {
			pairIndexes[i], pairIndexes[j] = pairIndexes[j], pairIndexes[i]
		})
		pair := pairIndexes[index]
		leftOperator := operators[pair[0]]
		rightOperator := operators[pair[1]]
		usedOperators := make([]chain.Address, 0, len(groupMembers))
		for _, operator := range groupMembers {
			if operator != leftOperator && operator != rightOperator {
				usedOperators = append(usedOperators, operator)
			}
		}
		return usedOperators, 0, true
	} else {
		return nil, len(pairIndexes), false
	}
}

func excludeOperatorTriplets(
	rng *rand.Rand,
	groupMembers []chain.Address,
	index int,
	operatorToSeatCount map[chain.Address]uint,
	operators []chain.Address,
	retryParticipantsCount int,
) ([]chain.Address, int, bool) {
	tripletIndexes := make([][3]int, 0, len(operators)*len(operators)*len(operators))
	for i := 0; i < len(operators)-2; i++ {
		for j := i + 1; j < len(operators)-1; j++ {
			for k := j + 1; k < len(operators); k++ {
				leftOperator := operators[i]
				middleOperator := operators[j]
				rightOperator := operators[j]
				count := len(groupMembers) -
					int(operatorToSeatCount[leftOperator]) -
					int(operatorToSeatCount[middleOperator]) -
					int(operatorToSeatCount[rightOperator])

				if count >= int(retryParticipantsCount) {
					tripletIndexes = append(tripletIndexes, [3]int{i, j, k})
				}
			}
		}
	}
	if index < len(tripletIndexes) {
		rng.Shuffle(len(tripletIndexes), func(i, j int) {
			tripletIndexes[i], tripletIndexes[j] = tripletIndexes[j], tripletIndexes[i]
		})
		triplet := tripletIndexes[index]
		leftOperator := operators[triplet[0]]
		middleOperator := operators[triplet[1]]
		rightOperator := operators[triplet[2]]
		usedOperators := make([]chain.Address, 0, len(groupMembers))
		for _, operator := range groupMembers {
			if operator != leftOperator && operator != middleOperator && operator != rightOperator {
				usedOperators = append(usedOperators, operator)
			}
		}
		return usedOperators, 0, true
	} else {
		return nil, len(tripletIndexes), false
	}
}

// EvaluateRetryParticipantsForKeyGeneration takes in a slice of `groupMembers`
// and returns a subslice of those same members of length >=
// `retryParticipantsCount` randomly according to the provided `seed` and
// `retryCount`.
//
// This function is intended to be called during key generation after a failure
// *not* due to inactivity. Assuming that some of the `groupMembers` are
// sending corrupted information, either on purpose or accidentally, we keep
// trying to find a subset of `groupMembers` that is as large as possible by
// first excluding single operators, then pairs of operators, then triplets of
// operators.
//
// The `seed` param needs to vary on a per-message basis but must be the same
// seed between all operators for each invocation. This can be the hash of the
// message since cryptographically secure randomness isn't important.
//
// The `retryCount` denotes the number of the given retry, so that should be
// incremented after each attempt while the `seed` stays consistent on a
// per-message basis.
func EvaluateRetryParticipantsForKeyGeneration(
	groupMembers []chain.Address,
	seed int64,
	retryCount uint,
	retryParticipantsCount uint,
) ([]chain.Address, error) {
	originalRetryCount := retryCount
	if int(retryParticipantsCount) > len(groupMembers) {
		return nil, fmt.Errorf(
			"asked for too many seats. [%d] seats were requested, "+
				"but there are only [%d] available.",
			retryParticipantsCount,
			len(groupMembers),
		)
	}
	operatorToSeatCount := calculateSeatCount(groupMembers)
	source := rand.NewSource(seed)
	rng := rand.New(source)

	operators := make([]chain.Address, 0, len(operatorToSeatCount))
	for operator := range operatorToSeatCount {
		if len(groupMembers)-int(operatorToSeatCount[operator]) >= int(retryParticipantsCount) {
			operators = append(operators, operator)
		}
	}
	sort.Sort(byAddress(operators))

	usedOperators, tries, ok := excludeSingleOperator(
		rng,
		groupMembers,
		int(retryCount),
		operatorToSeatCount,
		operators,
	)
	if ok {
		return usedOperators, nil
	} else {
		retryCount -= uint(tries)
	}

	usedOperators, tries, ok = excludeOperatorPairs(
		rng,
		groupMembers,
		int(retryCount),
		operatorToSeatCount,
		operators,
		int(retryParticipantsCount),
	)
	if ok {
		return usedOperators, nil
	} else {
		retryCount -= uint(tries)
	}

	usedOperators, tries, ok = excludeOperatorTriplets(
		rng,
		groupMembers,
		int(retryCount),
		operatorToSeatCount,
		operators,
		int(retryParticipantsCount),
	)
	if ok {
		return usedOperators, nil
	} else {
		retryCount -= uint(tries)
		return nil, fmt.Errorf(
			"the retry count [%d] was too large to handle! "+
				"Tried every single, pair, and triplet, but still needed [%d] more.",
			originalRetryCount,
			retryCount,
		)
	}
}
