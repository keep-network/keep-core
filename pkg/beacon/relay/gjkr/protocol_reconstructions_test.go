package gjkr

import (
	"math/big"
	"testing"
)

func TestReconstructPrivateKeyShares(t *testing.T) {
	threshold := 2
	groupSize := 5
	disqualifiedIDs := []int{3, 5}

	group, allDisqualifiedShares, _ := initializeReconstructingMembersGroup(threshold, groupSize, disqualifiedIDs)

	expectedPrivateKeyShare1 := group[2].secretCoefficients[0] // for ID = 3
	expectedPrivateKeyShare2 := group[4].secretCoefficients[0] // for ID = 5

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

func initializeReconstructingMembersGroup(
	threshold, groupSize int,
	disqualifiedIDs []int,
) ([]*ReconstructingMember, map[int]map[int]*big.Int, error) {
	sharingMembers, _ := initializeSharingMembersGroup(threshold, groupSize)

	var reconstructingMembers []*ReconstructingMember
	for _, sm := range sharingMembers {
		reconstructingMembers = append(reconstructingMembers,
			&ReconstructingMember{
				SharingMember: sm,
			},
		)
	}

	// Disqualified shares map for test run:
	// <disqualifiedID, <peerID, shareS>>
	allDisqualifiedShares := make(map[int]map[int]*big.Int, groupSize-len(disqualifiedIDs))
	for _, disqualifiedID := range disqualifiedIDs {
		disqualifiedShares := make(map[int]*big.Int, groupSize-len(disqualifiedIDs))
		for _, m := range sharingMembers {
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

	return reconstructingMembers, allDisqualifiedShares, nil
}
