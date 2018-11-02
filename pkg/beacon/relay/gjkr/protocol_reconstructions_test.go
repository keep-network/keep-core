package gjkr

import (
	"math/big"
	"testing"

	"github.com/keep-network/keep-core/pkg/beacon/relay/pedersen"
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

func TestCalculateReconstructedPublicKeyShares(t *testing.T) {
	groupSize := 3
	p := big.NewInt(179)
	g := big.NewInt(7)

	disqualifiedIDs := []int{4, 5} // m

	reconstructedPrivateKeyShares := make(map[int]*big.Int, len(disqualifiedIDs)) // z_m
	reconstructedPrivateKeyShares[4] = big.NewInt(14)                             // z_4
	reconstructedPrivateKeyShares[5] = big.NewInt(15)                             // z_5

	expectedPublicKeyShares := make(map[int]*big.Int, len(disqualifiedIDs)) // y_m = g^{z_m} mod p
	expectedPublicKeyShares[4] = big.NewInt(43)                             // 7^14 mod 179
	expectedPublicKeyShares[5] = big.NewInt(122)                            // 7^15 mod 179

	members := make([]*ReconstructingMember, groupSize)
	for i := range members {
		members[i] = &ReconstructingMember{
			SharingMember: &SharingMember{
				CommittingMember: &CommittingMember{
					memberCore: &memberCore{
						ID:             i,
						protocolConfig: &DKG{P: p},
					},
					vss: &pedersen.VSS{G: g},
				},
			},
			reconstructedPrivateKeyShares: reconstructedPrivateKeyShares,
		}
	}

	for _, rm := range members {
		rm.CalculateReconstructedPublicKeyShares()

		for m, expectedPublicKeyShare := range expectedPublicKeyShares {
			if rm.reconstructedPublicKeyShares[m].Cmp(expectedPublicKeyShare) != 0 {
				t.Fatalf("\nexpected: %s\nactual:   %s\n",
					expectedPublicKeyShare,
					rm.reconstructedPublicKeyShares[m],
				)
			}
		}
	}
}

func TestCombineGroupPublicKey(t *testing.T) {
	threshold := 3
	groupSize := 5
	p := big.NewInt(1907)

	var tests = map[string]struct {
		disqualifiedIDs        []int
		expectedError          error
		expectedGroupPublicKey *big.Int
	}{
		"no disqualified members - no reconstructed individual public key": {
			expectedError:          nil,
			expectedGroupPublicKey: big.NewInt(1156),
		},
		"2 disqualified members - 2 reconstructed individual public keys": {
			disqualifiedIDs:        []int{6, 7},
			expectedError:          nil,
			expectedGroupPublicKey: big.NewInt(1037),
		},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			// Prepare members group.
			members := make([]*ReconstructingMember, groupSize)
			for i := range members {
				members[i] = &ReconstructingMember{
					SharingMember: &SharingMember{
						CommittingMember: &CommittingMember{
							memberCore: &memberCore{
								ID: i + 1,
								protocolConfig: &DKG{
									P: p,
								},
							},
						},
					},
				}
			}

			// Generate member's public coefficients.
			for _, member := range members {
				member.publicCoefficients = make([]*big.Int, threshold+1)
				for k := range member.publicCoefficients {
					member.publicCoefficients[k] = big.NewInt(int64(member.ID*10 + k))
				}
			}

			// Configure public coefficients received from peer members.
			for _, member := range members {
				member.receivedGroupPublicKeyShares = make(map[int]*big.Int, groupSize-1)
				for _, peer := range members {
					if member.ID != peer.ID {
						member.receivedGroupPublicKeyShares[peer.ID] =
							peer.publicCoefficients[0]
					}
				}
			}

			// Configure reconstructed individual public key of disqualified members.
			for _, member := range members {
				member.reconstructedPublicKeyShares = make(map[int]*big.Int, len(test.disqualifiedIDs))
				for _, disqualifiedID := range test.disqualifiedIDs {
					member.reconstructedPublicKeyShares[disqualifiedID] =
						big.NewInt(int64(20 + disqualifiedID))
				}
			}

			for _, member := range members {
				member.CombineGroupPublicKey()

				if member.groupPublicKey.Cmp(test.expectedGroupPublicKey) != 0 {
					t.Fatalf(
						"incorrect group public key for member %d\nexpected: %v\nactual:   %v\n",
						member.ID,
						test.expectedGroupPublicKey,
						member.groupPublicKey,
					)
				}
			}
		})
	}
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
