package gjkr

import (
	"math/big"
	"testing"

	"github.com/keep-network/keep-core/pkg/beacon/relay/pedersen"
)

func TestReconstructIndividualPrivateKeys(t *testing.T) {
	threshold := 2
	groupSize := 5

	disqualifiedMembersIDs := []int{3, 5}

	group := initializeReconstructingMembersGroup(threshold, groupSize)

	disqualifiedMember1 := group[2] // for ID = 3
	disqualifiedMember2 := group[4] // for ID = 5

	expectedIndividualPrivateKey1 := disqualifiedMember1.secretCoefficients[0]
	expectedIndividualPrivateKey2 := disqualifiedMember2.secretCoefficients[0]

	allDisqualifiedShares := disqualifyMembers(group, disqualifiedMembersIDs)

	for _, m := range group {
		if !contains(disqualifiedMembersIDs, m.ID) {
			m.ReconstructIndividualPrivateKeys(allDisqualifiedShares)

			if m.reconstructedIndividualPrivateKeys[disqualifiedMember1.ID].Cmp(expectedIndividualPrivateKey1) != 0 {
				t.Fatalf("invalid reconstructed private key 1\nexpected: %s\nactual:   %s\n",
					expectedIndividualPrivateKey1,
					m.reconstructedIndividualPrivateKeys[disqualifiedMember1.ID],
				)
			}

			if m.reconstructedIndividualPrivateKeys[disqualifiedMember2.ID].Cmp(expectedIndividualPrivateKey2) != 0 {
				t.Fatalf("invalid reconstructed private key 2\nexpected: %s\nactual:   %s\n",
					expectedIndividualPrivateKey2,
					m.reconstructedIndividualPrivateKeys[disqualifiedMember2.ID],
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

	disqualifiedMembersIDs := []int{4, 5} // m

	reconstructedIndividualPrivateKeys := make(map[int]*big.Int, len(disqualifiedMembersIDs)) // z_m
	reconstructedIndividualPrivateKeys[4] = big.NewInt(14)                                    // z_4
	reconstructedIndividualPrivateKeys[5] = big.NewInt(15)                                    // z_5

	expectedIndividualPublicKeys := make(map[int]*big.Int, len(disqualifiedMembersIDs)) // y_m = g^{z_m} mod p
	expectedIndividualPublicKeys[4] = big.NewInt(43)                                    // 7^14 mod 179
	expectedIndividualPublicKeys[5] = big.NewInt(122)                                   // 7^15 mod 179

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

	for _, reconstructingMember := range members {
		reconstructingMember.ReconstructIndividualPublicKeys()

		for disqualifiedMemberID, expectedIndividualPublicKey := range expectedIndividualPublicKeys {
			if reconstructingMember.reconstructedIndividualPublicKeys[disqualifiedMemberID].
				Cmp(expectedIndividualPublicKey) != 0 {
				t.Fatalf("\nexpected: %s\nactual:   %s\n",
					expectedIndividualPublicKey,
					reconstructingMember.reconstructedIndividualPublicKeys[disqualifiedMemberID],
				)
			}
		}
	}
}

func initializeReconstructingMembersGroup(threshold, groupSize int) []*ReconstructingMember {
	// TODO When whole protocol is implemented check if SharingMember type is really
	// the one expected here (should be the member from Phase 10)
	sharingMembers, _ := initializeSharingMembersGroup(threshold, groupSize)

	var reconstructingMembers []*ReconstructingMember
	// TODO Should be handled by the `.Next()`` function
	for _, sm := range sharingMembers {
		reconstructingMembers = append(reconstructingMembers,
			&ReconstructingMember{
				SharingMember: sm,
			},
		)
	}

	return reconstructingMembers
}

// disqualifyMembers disqualifies specific members for a test run. It collects
// shares calculated by disqualified members for their peers and reveals them.
func disqualifyMembers(
	members []*ReconstructingMember,
	disqualifiedMembersIDs []int,
) []*DisqualifiedShares {
	allDisqualifiedShares := make([]*DisqualifiedShares, len(disqualifiedMembersIDs))
	for i, disqualifiedMemberID := range disqualifiedMembersIDs {
		sharesReceivedFromDisqualifiedMember := make(map[int]*big.Int, len(members)-len(disqualifiedMembersIDs))
		// for each group member
		for _, m := range members {
			// if the member has not been disqualified
			if !contains(disqualifiedMembersIDs, m.ID) {
				// collect all shares which this member received from disqualified
				// member and store them in sharesReceivedFromDisqualifiedMember
				for peerID, receivedShare := range m.receivedValidSharesS {
					if peerID == disqualifiedMemberID {
						sharesReceivedFromDisqualifiedMember[m.ID] = receivedShare
						break
					}
				}
			}
		}
		allDisqualifiedShares[i] = &DisqualifiedShares{
			disqualifiedMemberID: disqualifiedMemberID,
			peerSharesS:          sharesReceivedFromDisqualifiedMember,
		}
	}

	return allDisqualifiedShares
}
