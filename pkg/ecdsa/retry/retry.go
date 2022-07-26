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

func EvaluateRetryParticipantsForKeyGeneration(
	groupMembers []chain.Address,
	seed int64,
	retryCount uint,
	quantity uint,
) ([]chain.Address, error) {
	originalRetryCount := retryCount
	if int(quantity) > len(groupMembers) {
		return nil, fmt.Errorf(
			"Asked for too many seats. %d seats were requested, but there are only %d available.",
			quantity,
			len(groupMembers),
		)
	}
	operatorToSeatCount := calculateSeatCount(groupMembers)
	rand.Seed(seed)

	operators := make(byAddress, 0, len(operatorToSeatCount))
	for operator := range operatorToSeatCount {
		if len(groupMembers)-int(operatorToSeatCount[operator]) >= int(quantity) {
			operators = append(operators, operator)
		}
	}
	sort.Sort(operators)
	rand.Shuffle(len(operators), func(i, j int) {
		operators[i], operators[j] = operators[j], operators[i]
	})

	if int(retryCount) < len(operators) {
		removedOperator := operators[retryCount]
		usedOperators := make([]chain.Address, 0, len(groupMembers))
		for _, operator := range groupMembers {
			if operator != removedOperator {
				usedOperators = append(usedOperators, operator)
			}
		}
		return usedOperators, nil
	} else {
		retryCount -= uint(len(operators))
		for i := 0; i < len(operators)-1; i++ {
			for j := i + 1; j < len(operators); j++ {
				leftOperator := operators[i]
				rightOperator := operators[j]
				count := len(groupMembers) -
					int(operatorToSeatCount[leftOperator]) -
					int(operatorToSeatCount[rightOperator])

				if count >= int(quantity) {
					if retryCount == 0 {
						usedOperators := make([]chain.Address, 0, len(groupMembers))
						for _, operator := range groupMembers {
							if operator != leftOperator && operator != rightOperator {
								usedOperators = append(usedOperators, operator)
							}
							return usedOperators, nil
						}
					} else {
						retryCount -= 1
					}
				}
			}
		}

		for i := 0; i < len(operators)-2; i++ {
			for j := i + 1; i < len(operators)-1; j++ {
				for k := j + 1; j < len(operators); k++ {
					leftOperator := operators[i]
					middleOperator := operators[j]
					rightOperator := operators[j]
					count := len(groupMembers) -
						int(operatorToSeatCount[leftOperator]) -
						int(operatorToSeatCount[middleOperator]) -
						int(operatorToSeatCount[rightOperator])

					if count >= int(quantity) {
						if retryCount == 0 {
							usedOperators := make([]chain.Address, 0, len(groupMembers))
							for _, operator := range groupMembers {
								if operator != leftOperator && operator != middleOperator && operator != rightOperator {
									usedOperators = append(usedOperators, operator)
								}
								return usedOperators, nil
							}
						} else {
							retryCount -= 1
						}
					}
				}
			}
		}

		return nil, fmt.Errorf(
			"The retry count [%d] was too large to handle! Tried every single, pair, and triplet, but still needed [%d] more.",
			originalRetryCount,
			retryCount,
		)
	}
}
