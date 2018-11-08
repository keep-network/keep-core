package gjkr

import (
	"math/big"
	"testing"

	"github.com/keep-network/keep-core/pkg/beacon/relay/pedersen"
)

func TestReconstructIndividualPrivateKeys(t *testing.T) {
	threshold := 2
	groupSize := 5
	disqualifiedIDs := []int{3, 5}

	group, allDisqualifiedShares := initializeReconstructingMembersGroup(threshold, groupSize, disqualifiedIDs)

	expectedIndividualPrivateKey1 := group[2].secretCoefficients[0] // for ID = 3
	expectedIndividualPrivateKey2 := group[4].secretCoefficients[0] // for ID = 5

	for _, rm := range group {
		if !contains(disqualifiedIDs, rm.ID) {
			rm.ReconstructIndividualPrivateKeys(allDisqualifiedShares)

			if rm.reconstructedIndividualPrivateKeys[disqualifiedIDs[0]].Cmp(expectedIndividualPrivateKey1) != 0 {
				t.Fatalf("\nexpected: %s\nactual:   %s\n",
					expectedIndividualPrivateKey1,
					rm.reconstructedIndividualPrivateKeys[disqualifiedIDs[0]],
				)
			}

			if rm.reconstructedIndividualPrivateKeys[disqualifiedIDs[1]].Cmp(expectedIndividualPrivateKey2) != 0 {
				t.Fatalf("\nexpected: %s\nactual:   %s\n",
					expectedIndividualPrivateKey1,
					rm.reconstructedIndividualPrivateKeys[disqualifiedIDs[1]],
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

func TestCalculateReconstructedIndividualPublicKeys(t *testing.T) {
	groupSize := 3
	p := big.NewInt(179)
	g := big.NewInt(7)

	disqualifiedIDs := []int{4, 5} // m

	reconstructedIndividualPrivateKeys := make(map[int]*big.Int, len(disqualifiedIDs)) // z_m
	reconstructedIndividualPrivateKeys[4] = big.NewInt(14)                             // z_4
	reconstructedIndividualPrivateKeys[5] = big.NewInt(15)                             // z_5

	expectedIndividualPublicKeys := make(map[int]*big.Int, len(disqualifiedIDs)) // y_m = g^{z_m} mod p
	expectedIndividualPublicKeys[4] = big.NewInt(43)                             // 7^14 mod 179
	expectedIndividualPublicKeys[5] = big.NewInt(122)                            // 7^15 mod 179

	members := make([]*ReconstructingMember, groupSize)
	for i := range members {
		members[i] = &ReconstructingMember{
			SharingMember: &SharingMember{
				QualifiedMember: &QualifiedMember{
					SharesJustifyingMember: &SharesJustifyingMember{
						CommittingMember: &CommittingMember{
							memberCore: &memberCore{
								ID:             i,
								protocolConfig: &DKG{P: p},
							},
							vss: &pedersen.VSS{G: g},
						},
					},
				},
			},
			reconstructedIndividualPrivateKeys: reconstructedIndividualPrivateKeys,
		}
	}

	for _, rm := range members {
		rm.CalculateReconstructedIndividualPublicKeys()

		for m, expectedIndividualPublicKey := range expectedIndividualPublicKeys {
			if rm.reconstructedIndividualPublicKeys[m].Cmp(expectedIndividualPublicKey) != 0 {
				t.Fatalf("\nexpected: %s\nactual:   %s\n",
					expectedIndividualPublicKey,
					rm.reconstructedIndividualPublicKeys[m],
				)
			}
		}
	}
}

func initializeReconstructingMembersGroup(
	threshold, groupSize int,
	disqualifiedIDs []int,
) ([]*ReconstructingMember, []*DisqualifiedShares) {
	sharingMembers, _ := initializeSharingMembersGroup(threshold, groupSize)

	var reconstructingMembers []*ReconstructingMember
	for _, sm := range sharingMembers {
		reconstructingMembers = append(reconstructingMembers,
			&ReconstructingMember{
				SharingMember: sm,
			},
		)
	}

	// Disqualified shares for test run
	allDisqualifiedShares := make([]*DisqualifiedShares, len(disqualifiedIDs))
	for i, disqualifiedID := range disqualifiedIDs {
		shares := make(map[int]*big.Int, groupSize-len(disqualifiedIDs))
		for _, m := range sharingMembers {
			if !contains(disqualifiedIDs, m.ID) {
				for peerID, share := range m.receivedSharesS {
					if peerID == disqualifiedID {
						shares[m.ID] = share
						break
					}
				}
			}
		}
		allDisqualifiedShares[i] = &DisqualifiedShares{
			disqualifiedMemberID: disqualifiedID,
			peerSharesS:          shares,
		}
	}

	return reconstructingMembers, allDisqualifiedShares
}
