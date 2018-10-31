package gjkr

import (
	"math/big"
	"testing"
)

func TestReconstructPrivateKeyShares(t *testing.T) {
	threshold := 2
	groupSize := 5

	group, _ := initializeSharingMembersGroup(threshold, groupSize)

	disqualifiedIDs := []int{3, 5}

	expectedPrivateKeyShare1 := group[2].secretCoefficients[0] // for ID = 3
	expectedPrivateKeyShare2 := group[4].secretCoefficients[0] // for ID = 5

	// Prepare disqualified shares map for test run
	allDisqualifiedShares := make(map[int]map[int]*big.Int, groupSize-len(disqualifiedIDs))
	for _, disqualifiedID := range disqualifiedIDs {
		disqualifiedShares := make(map[int]*big.Int, groupSize-len(disqualifiedIDs))
		for _, m := range group {
			if !contains(disqualifiedIDs, m.ID) {
				for peerID, share := range m.receivedSharesS {
					if peerID == disqualifiedID {
						disqualifiedShares[m.ID] = share
					}
				}
			}
		}
		allDisqualifiedShares[disqualifiedID] = disqualifiedShares
	}
	for _, rm := range group {
		if !contains(disqualifiedIDs, rm.ID) {
			rm.ReconstructPrivateKeyShares(allDisqualifiedShares)

			if rm.reconstructedPrivateKeyShares[disqualifiedIDs[0]].Cmp(expectedPrivateKeyShare1) != 0 {
				t.Fatalf("\nexpected: %s\nactual:   %s\n",
					expectedPrivateKeyShare1,
					rm.reconstructedPrivateKeyShares[disqualifiedIDs[0]],
				)
			}

			if rm.reconstructedPrivateKeyShares[disqualifiedIDs[1]].Cmp(expectedPrivateKeyShare2) != 0 {
				t.Fatalf("\nexpected: %s\nactual:   %s\n",
					expectedPrivateKeyShare1,
					rm.reconstructedPrivateKeyShares[disqualifiedIDs[1]],
				)
			}
		}
	}
}

func contains(slice []int, value int) bool {
	for _, i := range slice {
		if i == value {
			return true
		}
	}
	return false
}
