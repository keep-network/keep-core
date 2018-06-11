package thresholdgroup

import (
	"fmt"
	"os"
	"testing"

	"github.com/dfinity/go-dfinity-crypto/bls"
)

func TestMain(m *testing.M) {
	bls.Init(bls.CurveSNARK1)

	os.Exit(m.Run())
}

const (
	defaultID                 = "12345"
	defaultDishonestThreshold = 4
	defaultGroupSize          = 12
)

func buildSharingMember(id string) *SharingMember {
	if id == "" {
		id = defaultID
	}
	member, _ := NewMember(id, defaultDishonestThreshold, defaultGroupSize)

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
	for i := 0; i < defaultDishonestThreshold; i++ {
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
		otherMember := buildCommittedSharingMember(memberID.GetHexString())
		otherMembers = append(otherMembers, otherMember)

		// Until we get to accusationCount, add invalid shares.
		if i < accusationCount {
			sharingMember.AddShareFromID(memberID, &bls.SecretKey{})
		} else {
			sharingMember.AddCommitmentsFromID(memberID, otherMember.Commitments())
			memberShare := otherMember.SecretShareForID(&sharingMember.BlsID)
			sharingMember.AddShareFromID(memberID, memberShare)
		}
	}

	return sharingMember.InitializeJustification(), otherMembers
}

func buildMembers(id string) (*Member, []*Member) {
	justifyingMember, otherMembers := buildJustifyingMember(id, 0)

	finalOtherMembers := make([]*Member, 0)
	for i, otherMember := range otherMembers {
		otherMember.AddCommitmentsFromID(
			&justifyingMember.BlsID,
			justifyingMember.Commitments(),
		)
		otherMember.AddShareFromID(
			&justifyingMember.BlsID,
			justifyingMember.SecretShareForID(&otherMember.BlsID),
		)
		for j, otherOtherMember := range otherMembers {
			if i == j {
				continue
			}

			otherMember.AddCommitmentsFromID(
				&otherOtherMember.BlsID,
				otherOtherMember.Commitments(),
			)
			otherMember.AddShareFromID(
				&otherOtherMember.BlsID,
				otherOtherMember.SecretShareForID(&otherMember.BlsID),
			)
		}

		finalMember, err := otherMember.InitializeJustification().FinalizeMember()
		if err != nil {
			panic(err)
		}

		finalOtherMembers = append(finalOtherMembers, finalMember)
	}

	finalMember, err := justifyingMember.FinalizeMember()
	if err != nil {
		panic(err)
	}

	return finalMember, finalOtherMembers
}
