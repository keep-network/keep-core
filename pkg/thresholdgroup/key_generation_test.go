package thresholdgroup

import (
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"testing"

	"github.com/dfinity/go-dfinity-crypto/bls"
)

func TestLocalMemberCreation(t *testing.T) {
	id := fmt.Sprintf("%x", rand.Int31())

	member, err := NewMember(id, defaultDishonestThreshold, defaultGroupSize)
	if err != nil {
		t.Fatalf("unexpected error [%v]", err)
	}

	if member == nil {
		t.Fatal("\nexpected: non-nil\nactual: nil")
	}

	var propertyTests = map[string]struct {
		propertyFunc func(*LocalMember) string
		expected     string
	}{
		"id": {
			propertyFunc: func(lm *LocalMember) string { return lm.ID },
			expected:     fmt.Sprintf("0x%010s", id),
		},
		"BLS ID": {
			propertyFunc: func(lm *LocalMember) string { return lm.BlsID.GetHexString() },
			expected:     id,
		},
		"threshold": {
			propertyFunc: func(lm *LocalMember) string { return fmt.Sprintf("%v", lm.dishonestThreshold) },
			expected:     fmt.Sprintf("%v", defaultDishonestThreshold),
		},
		"group size": {
			propertyFunc: func(lm *LocalMember) string { return fmt.Sprintf("%v", lm.groupSize) },
			expected:     fmt.Sprintf("%v", defaultGroupSize),
		},
	}

	for testName, test := range propertyTests {
		t.Run(testName, func(t *testing.T) {
			property := test.propertyFunc(member)
			if test.expected != property {
				t.Errorf("\nexpected: %s\nactual:   %s", test.expected, property)
			}
		})
	}
}

func TestLocalMemberFailsForHighThreshold(t *testing.T) {
	member, err := NewMember(defaultID, defaultGroupSize/2, defaultGroupSize)
	if err == nil {
		t.Fatal("\nexpected error but got nil")
	}
	if member != nil {
		t.Fatalf("\nexpected nil member for errored instantiation\ngot [%v]", member)
	}
}

func TestLocalMemberCommitments(t *testing.T) {
	member, _ := NewMember(defaultID, defaultDishonestThreshold, defaultGroupSize)

	expectedShareCommitmentsCount := defaultDishonestThreshold + 1

	if len(member.Commitments()) != expectedShareCommitmentsCount {
		t.Errorf(
			"\nexpected: %v commitments\nactual:   %v commitments",
			expectedShareCommitmentsCount,
			len(member.Commitments()),
		)
	}

	// Smoke tests: all commitments are initialized, no two commitments in a row
	// are the same.
	uninitializedCommitment := bls.PublicKey{}
	previousCommitment := uninitializedCommitment
	for i, commitment := range member.Commitments() {
		if commitment.IsEqual(&uninitializedCommitment) {
			t.Fatalf(
				"at index %v\nexpected: initialized commitment\nactual:   uninitialized commitment",
				i,
			)
		} else if commitment.IsEqual(&previousCommitment) {
			t.Errorf(
				"at index %v\nexpected: different commitment\nactual:   same as previous commitment",
				i,
			)
		}

		previousCommitment = commitment
	}
}

func TestLocalMemberRegistration(t *testing.T) {
	member, _ := NewMember(defaultID, defaultDishonestThreshold, defaultGroupSize)
	member.RegisterMemberID(&member.BlsID)

	otherMemberCount := defaultGroupSize - 1
	for i := 0; i < otherMemberCount; i++ {
		if member.MemberListComplete() {
			t.Fatalf(
				"\nmember list complete after %v instead of %v additional members",
				i+1,
				otherMemberCount,
			)
		}

		id := bls.ID{}
		id.SetDecString(fmt.Sprintf("%v", i+1))
		member.RegisterMemberID(&id)
	}

	if !member.MemberListComplete() {
		t.Errorf(
			"\nexpected: member list complete after %v additional members\n"+
				"actual:   member list incomplete",
			otherMemberCount,
		)
	}
}

func TestSharingOtherMemberIDs(t *testing.T) {
	sharingMember := buildSharingMember("")

	otherMemberIDs := sharingMember.OtherMemberIDs()
	if len(otherMemberIDs) != defaultGroupSize-1 {
		t.Errorf(
			"\nexpected: %v other member ids\nactual:   %v other member ids",
			defaultGroupSize-1,
			len(otherMemberIDs),
		)
	}
	for i, id := range otherMemberIDs {
		if id.GetDecString() != fmt.Sprintf("%v", i+1) {
			t.Errorf(
				"at index %v\nexpected id: %v\nactual id:   %v",
				i,
				fmt.Sprintf("%v", i+1),
				id.GetDecString(),
			)
		}
	}
}

func TestSharingMemberSecretShares(t *testing.T) {
	sharingMember := buildSharingMember("")

	memberIDs := sharingMember.OtherMemberIDs()
	memberIDs = append(memberIDs, &sharingMember.BlsID)
	uninitializedShare := &bls.SecretKey{}
	previousShare := uninitializedShare
	for _, id := range memberIDs {
		share := sharingMember.SecretShareForID(id)
		if share.IsEqual(uninitializedShare) {
			t.Fatalf(
				"for id %v\nexpected: initialized share\nactual:   uninitialized share",
				id.GetHexString(),
			)
		} else if share.IsEqual(previousShare) {
			t.Errorf(
				"for id %v\nexpected: different share\nactual:   same as previous share",
				id.GetHexString(),
			)
		}

		previousShare = share
	}
}

func TestSharingMemberCommitmentReception(t *testing.T) {
	sharingMember := buildSharingMember("")
	completeCommitmentCount := defaultGroupSize - 1

	for i, memberID := range sharingMember.OtherMemberIDs() {
		if sharingMember.CommitmentsComplete() {
			t.Fatalf(
				"\ncommitments complete after %v instead of %v members",
				i+1,
				completeCommitmentCount,
			)
		}

		commitments := randomCommitments()
		sharingMember.AddCommitmentsFromID(memberID, commitments)
	}

	if !sharingMember.CommitmentsComplete() {
		t.Errorf(
			"\nexpected: commitments complete after %v members\nactual:   commitments incomplete",
			completeCommitmentCount,
		)
	}
}

func TestSharingMemberPrivateShareReception(t *testing.T) {
	completeShareCount := defaultGroupSize - 1

	var tests = map[string]struct {
		memberShareFunc     func(*SharingMember, *bls.ID) *bls.SecretKey
		expectedAccusations int
	}{
		"empty member commitments": {
			memberShareFunc: func(sharingMember *SharingMember, otherMemberID *bls.ID) *bls.SecretKey {
				sharingMember.AddCommitmentsFromID(
					otherMemberID,
					make([]bls.PublicKey, 0),
				)
				return &bls.SecretKey{}
			},
			expectedAccusations: 11,
		},
		"valid shares": {
			memberShareFunc: func(sharingMember *SharingMember, otherMemberID *bls.ID) *bls.SecretKey {
				otherMember := buildCommittedSharingMember(otherMemberID.GetHexString())
				sharingMember.AddCommitmentsFromID(otherMemberID, otherMember.Commitments())
				memberShare := otherMember.SecretShareForID(&sharingMember.BlsID)
				return memberShare
			},
			expectedAccusations: 0,
		},
		"invalid shares": {
			memberShareFunc: func(_ *SharingMember, _ *bls.ID) *bls.SecretKey {
				return &bls.SecretKey{}
			},
			expectedAccusations: 11,
		},
		"invalid and valid shares": {
			memberShareFunc: func(sharingMember *SharingMember, otherMemberID *bls.ID) *bls.SecretKey {
				// Return an invalid share for 5 shares.
				if id, _ := strconv.Atoi(otherMemberID.GetDecString()); id < 6 {
					return &bls.SecretKey{}
				}

				otherMember := buildCommittedSharingMember(otherMemberID.GetHexString())
				sharingMember.AddCommitmentsFromID(otherMemberID, otherMember.Commitments())
				memberShare := otherMember.SecretShareForID(&sharingMember.BlsID)
				return memberShare
			},
			expectedAccusations: 5,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			sharingMember := buildCommittedSharingMember("")
			for i, memberID := range sharingMember.OtherMemberIDs() {
				if sharingMember.SharesComplete() {
					t.Fatalf(
						"\nshares complete after %v instead of %v members",
						i+1,
						completeShareCount,
					)
				}

				memberShare := test.memberShareFunc(sharingMember, memberID)
				sharingMember.AddShareFromID(memberID, memberShare)
			}

			if !sharingMember.SharesComplete() {
				t.Errorf(
					"\nexpected: shares complete after %v members\nactual:   shares incomplete",
					completeShareCount,
				)
			}

			accusedIDs := sharingMember.AccusedIDs()
			if len(accusedIDs) != test.expectedAccusations {
				t.Errorf(
					"\nexpected: %v accusations\nactual:   %v accusations",
					test.expectedAccusations,
					len(accusedIDs),
				)
			}
		})
	}
}

func TestSharingMemberMissingShareAccusations(t *testing.T) {
	// Note: current calling code can't enter this state, but we're checking
	// that it works just in case.
	var tests = map[string]struct {
		memberShareFunc     func(*SharingMember, *bls.ID) *bls.SecretKey
		expectedAccusations int
	}{
		"missing shares": {
			memberShareFunc: func(_ *SharingMember, _ *bls.ID) *bls.SecretKey {
				return nil
			},
			expectedAccusations: 11,
		},
		"missing and valid shares": {
			memberShareFunc: func(sharingMember *SharingMember, otherMemberID *bls.ID) *bls.SecretKey {
				// Return no share for 5 shares.
				if id, _ := strconv.Atoi(otherMemberID.GetDecString()); id < 6 {
					return nil
				}

				otherMember := buildCommittedSharingMember(otherMemberID.GetHexString())
				sharingMember.AddCommitmentsFromID(otherMemberID, otherMember.Commitments())
				memberShare := otherMember.SecretShareForID(&sharingMember.BlsID)
				return memberShare
			},
			expectedAccusations: 5,
		},
		"missing, invalid, and valid shares": {
			memberShareFunc: func(sharingMember *SharingMember, otherMemberID *bls.ID) *bls.SecretKey {
				// Return no share for 5 shares, invalid share for the next 4.
				if id, _ := strconv.Atoi(otherMemberID.GetDecString()); id < 6 {
					return nil
				} else if id, _ := strconv.Atoi(otherMemberID.GetDecString()); id < 10 {
					return &bls.SecretKey{}
				}

				otherMember := buildCommittedSharingMember(otherMemberID.GetHexString())
				sharingMember.AddCommitmentsFromID(otherMemberID, otherMember.Commitments())
				memberShare := otherMember.SecretShareForID(&sharingMember.BlsID)
				return memberShare
			},
			expectedAccusations: 9,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			sharingMember := buildCommittedSharingMember("")
			for _, memberID := range sharingMember.OtherMemberIDs() {
				memberShare := test.memberShareFunc(sharingMember, memberID)
				if memberShare != nil {
					sharingMember.AddShareFromID(memberID, memberShare)
				}
			}

			accusedIDs := sharingMember.AccusedIDs()
			if len(accusedIDs) != test.expectedAccusations {
				t.Errorf(
					"\nexpected: %v accusations\nactual:   %v accusations",
					test.expectedAccusations,
					len(accusedIDs),
				)
			}
		})
	}
}

func TestJustifyingMemberSelfAccusationRecording(t *testing.T) {
	justifyingMember, _ := buildJustifyingMember("", 0)

	otherMemberIDs := justifyingMember.OtherMemberIDs()
	randomMemberIDs := make([]*bls.ID, 0)
	for _, i := range rand.Perm(len(otherMemberIDs)) {
		randomMemberIDs = append(randomMemberIDs, otherMemberIDs[i])
	}

	t.Run("accusations against self produce justifications", func(t *testing.T) {
		for _, id := range randomMemberIDs {
			justifyingMember.AddAccusationFromID(id, &justifyingMember.BlsID)
		}

		justifications := justifyingMember.Justifications()
		for _, id := range justifyingMember.OtherMemberIDs() {
			share, exists := justifications[*id]
			if !exists {
				t.Errorf(
					"\nexpected: justification for %v\nactual:   no justification",
					id.GetHexString(),
				)
				continue
			}

			if !share.IsEqual(justifyingMember.memberShares[*id]) {
				t.Errorf(
					"\nexpected: valid justification for %v\nactual:   bad justification",
					id.GetHexString(),
				)
			}
		}
	})
}

func TestJustifyingMemberFinalization(t *testing.T) {
	accuse := func(count int) func(*JustifyingMember, []*SharingMember, []*bls.ID) {
		return func(
			justifyingMember *JustifyingMember,
			otherMembers []*SharingMember,
			randomMemberIDs []*bls.ID,
		) {
			for i, otherMember := range otherMembers[0:count] {
				justifyingMember.AddAccusationFromID(
					randomMemberIDs[i],
					&otherMember.BlsID,
				)
			}
		}
	}
	accuseAll := accuse(defaultGroupSize - 1)

	justify := func(goodCount int, badCount int) func(*JustifyingMember, []*SharingMember, []*bls.ID) {
		return func(
			justifyingMember *JustifyingMember,
			otherMembers []*SharingMember,
			randomMemberIDs []*bls.ID,
		) {
			if goodCount > 0 {
				for i, otherMember := range otherMembers[0:goodCount] {
					justifyingMember.RecordJustificationFromID(
						&otherMember.BlsID,
						randomMemberIDs[i],
						otherMember.SecretShareForID(randomMemberIDs[i]),
					)
				}
			}

			if badCount > 0 {
				for i, otherMember := range otherMembers[goodCount:badCount] {
					justifyingMember.RecordJustificationFromID(
						&otherMember.BlsID,
						randomMemberIDs[i+2],
						&bls.SecretKey{},
					)
				}
			}
		}
	}

	var tests = map[string]struct {
		badShares     int
		accuseFunc    func(*JustifyingMember, []*SharingMember, []*bls.ID)
		justifyFunc   func(*JustifyingMember, []*SharingMember, []*bls.ID)
		expectedError error
	}{
		"all accused without justification": {
			accuseFunc:    accuseAll,
			justifyFunc:   justify(0, 0),
			expectedError: fmt.Errorf("required 5 qualified members but only had 1"),
		},
		"all accused with less than `defaultDishonestThreshold` justifications": {
			accuseFunc:    accuseAll,
			justifyFunc:   justify(defaultDishonestThreshold-1, 0),
			expectedError: fmt.Errorf("required 5 qualified members but only had 4"),
		},
		"all accused with bad justifications and less than `defaultDishonestThreshold` good justifications": {
			accuseFunc:    accuseAll,
			justifyFunc:   justify(defaultDishonestThreshold-1, defaultGroupSize-2),
			expectedError: fmt.Errorf("required 5 qualified members but only had 4"),
		},
		"all accused with `defaultDishonestThreshold` justifications": {
			accuseFunc:  accuseAll,
			justifyFunc: justify(defaultDishonestThreshold, 0),
		},
		"all accused with `defaultDishonestThreshold+1` justifications": {
			accuseFunc:  accuseAll,
			justifyFunc: justify(defaultDishonestThreshold+1, 0),
		},
		"all accused with all justifications": {
			accuseFunc:  accuseAll,
			justifyFunc: justify(defaultGroupSize-1, 0),
		},
		"`defaultDishonestThreshold` honest with no justifications": {
			accuseFunc:    accuse(defaultGroupSize - defaultDishonestThreshold),
			justifyFunc:   justify(0, 0),
			expectedError: fmt.Errorf("required 5 qualified members but only had 4"),
		},
		"`defaultDishonestThreshold+1` accused with `defaultDishonestThreshold` justifications": {
			accuseFunc:  accuse(defaultDishonestThreshold + 1),
			justifyFunc: justify(defaultDishonestThreshold, 0),
		},
		"bad shares without justification": {
			badShares:  1,
			accuseFunc: accuseAll,
			justifyFunc: func(
				justifyingMember *JustifyingMember,
				otherMembers []*SharingMember,
				randomMemberIDs []*bls.ID,
			) {
				for _, otherMember := range otherMembers[0:1] {
					justifyingMember.RecordJustificationFromID(
						&otherMember.BlsID,
						&justifyingMember.BlsID,
						&bls.SecretKey{},
					)
				}
			},
			expectedError: fmt.Errorf("required 5 qualified members but only had 1"),
		},
		"bad shares with justification": {
			badShares:  5,
			accuseFunc: accuse(0),
			justifyFunc: func(
				justifyingMember *JustifyingMember,
				otherMembers []*SharingMember,
				randomMemberIDs []*bls.ID,
			) {
				for _, otherMember := range otherMembers[0:5] {
					// FIXME Can we make this happen automatically?
					justifyingMember.AddCommitmentsFromID(&otherMember.BlsID, otherMember.Commitments())
					justifyingMember.RecordJustificationFromID(
						&otherMember.BlsID,
						&justifyingMember.BlsID,
						otherMember.SecretShareForID(&justifyingMember.BlsID),
					)
				}
			},
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			justifyingMember, otherMembers := buildJustifyingMember("", test.badShares)

			randomMemberIDs := make([]*bls.ID, 0)
			for _, i := range rand.Perm(len(otherMembers)) {
				randomMemberIDs = append(randomMemberIDs, &otherMembers[i].BlsID)
			}

			test.accuseFunc(justifyingMember, otherMembers, randomMemberIDs)
			test.justifyFunc(justifyingMember, otherMembers, randomMemberIDs)

			_, err := justifyingMember.FinalizeMember()
			if !reflect.DeepEqual(err, test.expectedError) {
				t.Errorf("\nexpected: %v\nactual:   %v", test.expectedError, err)
			}
		})
	}
}
