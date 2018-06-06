package thresholdgroup

import (
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/dfinity/go-dfinity-crypto/bls"
)

func TestMain(m *testing.M) {
	bls.Init(bls.CurveSNARK1)

	os.Exit(m.Run())
}

var (
	defaultID        = "12345"
	defaultThreshold = 4
	defaultGroupSize = 12
)

func TestLocalMemberCreation(t *testing.T) {
	id := fmt.Sprintf("%x", rand.Int31())

	member, err := NewMember(id, defaultThreshold, defaultGroupSize)
	if err != nil {
		t.Fatalf("unexpected error [%v]", err)
	}

	if member == nil {
		t.Fatal("expected: non-nil\nactual: nil")
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
			propertyFunc: func(lm *LocalMember) string { return fmt.Sprintf("%v", lm.threshold) },
			expected:     fmt.Sprintf("%v", defaultThreshold),
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
	member, _ := NewMember(defaultID, defaultThreshold, defaultGroupSize)

	if len(member.Commitments()) != defaultThreshold {
		t.Errorf(
			"\nexpected: %v commitments\nactual:   %v commitments",
			defaultThreshold,
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
	member, _ := NewMember(defaultID, defaultThreshold, defaultGroupSize)
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

func buildSharingMember(id string) *SharingMember {
	if id == "" {
		id = defaultID
	}
	member, _ := NewMember(id, defaultThreshold, defaultGroupSize)

	defaultBlsID := &bls.ID{}
	defaultBlsID.SetHexString(defaultID)
	member.RegisterMemberID(defaultBlsID)
	for i := 1; i < defaultGroupSize; i++ {
		id := bls.ID{}
		id.SetDecString(fmt.Sprintf("%v", i))
		member.RegisterMemberID(&id)
	}

	return member.InitializeSharing()
}

func randomShares() []bls.SecretKey {
	secretKeys := make([]bls.SecretKey, 0)
	for i := 0; i < defaultThreshold; i++ {
		sk := bls.SecretKey{}
		sk.SetByCSPRNG()
		secretKeys = append(secretKeys, sk)
	}

	return secretKeys
}

func commitmentsFromShares(shares []bls.SecretKey) []bls.PublicKey {
	commitments := make([]bls.PublicKey, 0)
	for _, share := range shares {
		commitments = append(commitments, *share.GetPublicKey())
	}

	return commitments
}

func randomCommitments() []bls.PublicKey {
	return commitmentsFromShares(randomShares())
}

func buildCommittedSharingMember(id string) *SharingMember {
	sharingMember := buildSharingMember(id)

	for _, memberID := range sharingMember.OtherMemberIDs() {
		commitments := randomCommitments()
		sharingMember.AddCommitmentsFromID(memberID, commitments)
	}

	return sharingMember
}

func buildJustifyingMember(
	id string,
	accusationCount int,
) (*JustifyingMember, []*SharingMember) {
	sharingMember := buildCommittedSharingMember(id)
	otherMembers := make([]*SharingMember, 0)

	for i, memberID := range sharingMember.OtherMemberIDs() {
		// Until we get to accusationCount, add invalid shares.
		if i < accusationCount {
			sharingMember.AddShareFromID(memberID, &bls.SecretKey{})
		} else {
			otherMember := buildCommittedSharingMember(memberID.GetHexString())
			otherMembers = append(otherMembers, otherMember)
			sharingMember.AddCommitmentsFromID(memberID, otherMember.Commitments())
			memberShare := otherMember.SecretShareForID(&sharingMember.BlsID)
			sharingMember.AddShareFromID(memberID, memberShare)
		}
	}

	return sharingMember.InitializeJustification(), otherMembers
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
	var tests = map[string]struct {
		accuseFunc    func(*JustifyingMember, []*SharingMember, []*bls.ID)
		justifyFunc   func(*JustifyingMember, []*SharingMember, []*bls.ID)
		expectedError error
	}{
		"accused without justification": {
			accuseFunc: func(
				justifyingMember *JustifyingMember,
				otherMembers []*SharingMember,
				randomMemberIDs []*bls.ID,
			) {
				for i, otherMember := range otherMembers {
					justifyingMember.AddAccusationFromID(
						randomMemberIDs[i],
						&otherMember.BlsID,
					)
				}
			},
			justifyFunc:   func(*JustifyingMember, []*SharingMember, []*bls.ID) {},
			expectedError: fmt.Errorf("required 4 qualified members but only had 1"),
		},
		"all accused with less than threshold justifications": {
			accuseFunc: func(
				justifyingMember *JustifyingMember,
				otherMembers []*SharingMember,
				randomMemberIDs []*bls.ID,
			) {
				for i, otherMember := range otherMembers {
					justifyingMember.AddAccusationFromID(
						randomMemberIDs[i],
						&otherMember.BlsID,
					)
				}
			},
			justifyFunc: func(
				justifyingMember *JustifyingMember,
				otherMembers []*SharingMember,
				randomMemberIDs []*bls.ID,
			) {
				for i, otherMember := range otherMembers[0:2] {
					justifyingMember.RecordJustificationFromID(
						&otherMember.BlsID,
						randomMemberIDs[i],
						otherMember.SecretShareForID(randomMemberIDs[i]),
					)
				}
			},
			expectedError: fmt.Errorf("required 4 qualified members but only had 3"),
		},
		"threshold accused with no justifications": {
			accuseFunc: func(
				justifyingMember *JustifyingMember,
				otherMembers []*SharingMember,
				randomMemberIDs []*bls.ID,
			) {
				for i, otherMember := range otherMembers[0:9] {
					justifyingMember.AddAccusationFromID(
						randomMemberIDs[i],
						&otherMember.BlsID,
					)
				}
			},
			justifyFunc:   func(*JustifyingMember, []*SharingMember, []*bls.ID) {},
			expectedError: fmt.Errorf("required 4 qualified members but only had 3"),
		},
		"all accused with threshold justifications": {
			accuseFunc: func(
				justifyingMember *JustifyingMember,
				otherMembers []*SharingMember,
				randomMemberIDs []*bls.ID,
			) {
				for i, otherMember := range otherMembers {
					justifyingMember.AddAccusationFromID(
						randomMemberIDs[i],
						&otherMember.BlsID,
					)
				}
			},
			justifyFunc: func(
				justifyingMember *JustifyingMember,
				otherMembers []*SharingMember,
				randomMemberIDs []*bls.ID,
			) {
				for i, otherMember := range otherMembers[0:3] {
					justifyingMember.RecordJustificationFromID(
						&otherMember.BlsID,
						randomMemberIDs[i],
						otherMember.SecretShareForID(randomMemberIDs[i]),
					)
				}
			},
			expectedError: nil,
		},
		"all accused with all justifications": {
			accuseFunc: func(
				justifyingMember *JustifyingMember,
				otherMembers []*SharingMember,
				randomMemberIDs []*bls.ID,
			) {
				for i, otherMember := range otherMembers {
					justifyingMember.AddAccusationFromID(
						randomMemberIDs[i],
						&otherMember.BlsID,
					)
				}
			},
			justifyFunc: func(
				justifyingMember *JustifyingMember,
				otherMembers []*SharingMember,
				randomMemberIDs []*bls.ID,
			) {
				for i, otherMember := range otherMembers {
					justifyingMember.RecordJustificationFromID(
						&otherMember.BlsID,
						randomMemberIDs[i],
						otherMember.SecretShareForID(randomMemberIDs[i]),
					)
				}
			},
			expectedError: nil,
		},
		"threshold accused with threshold justifications": {
			accuseFunc: func(
				justifyingMember *JustifyingMember,
				otherMembers []*SharingMember,
				randomMemberIDs []*bls.ID,
			) {
				for i, otherMember := range otherMembers {
					justifyingMember.AddAccusationFromID(
						randomMemberIDs[i],
						&otherMember.BlsID,
					)
				}
			},
			justifyFunc: func(
				justifyingMember *JustifyingMember,
				otherMembers []*SharingMember,
				randomMemberIDs []*bls.ID,
			) {
				for i, otherMember := range otherMembers {
					justifyingMember.RecordJustificationFromID(
						&otherMember.BlsID,
						randomMemberIDs[i],
						otherMember.SecretShareForID(randomMemberIDs[i]),
					)
				}
			},
			expectedError: nil,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			justifyingMember, otherMembers := buildJustifyingMember("", 0)

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
