package retry

import (
	"fmt"
	"testing"

	"github.com/keep-network/keep-core/pkg/chain"
)

func isSubset(groupMembers []chain.Address, subset []chain.Address, t *testing.T) {
	memberMap := make(map[chain.Address]struct{})
	for _, operator := range groupMembers {
		memberMap[operator] = struct{}{}
	}
	for _, operator := range subset {
		if _, ok := memberMap[operator]; !ok {
			t.Errorf("Subset member [%s] is not in the operator group.", operator)
		}
	}
}

func testEq(a, b []chain.Address) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func isStable(groupMembers []chain.Address, subset []chain.Address, quantity uint, t *testing.T) {
	for i := 0; i < 30; i++ {
		newSubset, _ := EvaluateRetryParticipantsForSigning(groupMembers, int64(123), 0, quantity)
		if ok := testEq(subset, newSubset); !ok {
			t.Errorf(
				"The subsets changed\nexpected: [%v]\nactual:   [%v]",
				subset,
				newSubset,
			)
		}
	}
}

func isLargeEnough(subset []chain.Address, quantity uint, t *testing.T) {
	if len(subset) < int(quantity) {
		t.Errorf(
			"Subset isn't large enough\nexpected: [%d+]\nactual:   [%d]",
			quantity,
			len(subset),
		)
	}
}

// They don't all have to be different, but they shouldn't all be the same!
func affectedBySeed(groupMembers []chain.Address, subset []chain.Address, quantity uint, t *testing.T) {
	allTheSame := true
	for seed := int64(0); seed < 30 && allTheSame; seed++ {
		newSubset, _ := EvaluateRetryParticipantsForSigning(groupMembers, seed, 0, quantity)
		allTheSame = allTheSame && testEq(subset, newSubset)
	}
	if allTheSame {
		t.Error("The seed did not affect the subset generation. All subsets were the same.")
	}
}

// They don't all have to be different, but they shouldn't all be the same!
func affectedByRetryCount(groupMembers []chain.Address, quantity uint, t *testing.T) {
	allTheSame := true
	subset, _ := EvaluateRetryParticipantsForSigning(groupMembers, int64(72312), uint(0), quantity)
	for retryCount := uint(1); retryCount < 30 && allTheSame; retryCount++ {
		newSubset, _ := EvaluateRetryParticipantsForSigning(groupMembers, int64(72312), retryCount, quantity)
		allTheSame = allTheSame && testEq(subset, newSubset)
	}
	if allTheSame {
		t.Error("The seed did not affect the subset generation. All subsets were the same.")
	}
}

func invariants(groupMembers []chain.Address, subset []chain.Address, quantity uint, t *testing.T) {
	isSubset(groupMembers, subset, t)
	isStable(groupMembers, subset, quantity, t)
	isLargeEnough(subset, quantity, t)
	affectedBySeed(groupMembers, subset, quantity, t)
	affectedByRetryCount(groupMembers, quantity, t)
}

func TestEvaluateRetryParticipantsForSigning_100DifferentOperators(t *testing.T) {
	groupMembers := make([]chain.Address, 100)
	for i := 0; i < 100; i++ {
		groupMembers[i] = chain.Address(fmt.Sprintf("Operator-%d", i))
	}
	subset, _ := EvaluateRetryParticipantsForSigning(groupMembers, int64(123), 0, 51)
	invariants(groupMembers, subset, 51, t)
}

func TestEvaluateRetryParticipantsForSigning_FewOperators(t *testing.T) {
	groupMembers := make([]chain.Address, 100)
	for i := 0; i < 100; i++ {
		groupMembers[i] = chain.Address(fmt.Sprintf("Operator-%d", i%3))
	}
	subset, _ := EvaluateRetryParticipantsForSigning(groupMembers, int64(456), 0, 51)
	invariants(groupMembers, subset, 51, t)
}
