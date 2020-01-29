package dkgtest

import (
	"testing"

	"github.com/keep-network/keep-core/pkg/altbn128"
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
	"github.com/keep-network/keep-core/pkg/internal/testutils"
)

// AssertDkgResultPublished checks if DKG result has been published to the
// chain. It does not inspect the result.
func AssertDkgResultPublished(t *testing.T, testResult *Result) {
	if testResult.dkgResult == nil {
		t.Fatal("dkg result is nil")
	}
}

// AssertSuccessfulSignersCount checks the number of successful signers. It does
// not check which particular signers were successful.
func AssertSuccessfulSignersCount(
	t *testing.T,
	testResult *Result,
	expectedCount int,
) {
	if len(testResult.signers) != expectedCount {
		t.Errorf(
			"unexpected number of successful signers\nexpected: [%v]\nactual:   [%v]",
			expectedCount,
			len(testResult.signers),
		)
	}
}

// AssertSuccessfulSigners checks which particular signers were successful.
func AssertSuccessfulSigners(
	t *testing.T,
	testResult *Result,
	expectedSuccessfulMembers ...group.MemberIndex,
) {
	actualSuccessfulMembers := make([]group.MemberIndex, len(testResult.signers))
	for _, signer := range testResult.signers {
		memberIndex := signer.MemberID()
		actualSuccessfulMembers = append(actualSuccessfulMembers, memberIndex)

		isSuccessfulExpected := containsMemberIndex(
			memberIndex,
			expectedSuccessfulMembers,
		)

		if !isSuccessfulExpected {
			t.Errorf(
				"member [%v] should not be a successful signer",
				memberIndex,
			)
		}
	}

	for _, memberIndex := range expectedSuccessfulMembers {
		isSuccessful := containsMemberIndex(
			memberIndex,
			actualSuccessfulMembers,
		)

		if !isSuccessful {
			t.Errorf(
				"member [%v] should be a successful signer",
				memberIndex,
			)
		}
	}
}

// AssertMemberFailuresCount checks the number of members who failed the
// protocol execution. It does not check which particular members failed.
func AssertMemberFailuresCount(
	t *testing.T,
	testResult *Result,
	expectedCount int,
) {
	if len(testResult.memberFailures) != expectedCount {
		t.Errorf(
			"unexpected number of member failures\nexpected: [%v]\nactual:   [%v]",
			expectedCount,
			len(testResult.memberFailures),
		)
	}
}

// AssertNoDisqualifiedMembers checks there were no disqualified members during
// the protocol execution.
func AssertNoDisqualifiedMembers(t *testing.T, testResult *Result) {
	AssertDisqualifiedMembers(t, testResult)
}

// AssertDisqualifiedMembers checks which members were disqualified
// during the protocol execution and compares them against expected ones.
func AssertDisqualifiedMembers(
	t *testing.T,
	testResult *Result,
	expectedDisqualifiedMembers ...group.MemberIndex,
) {
	actualDisqualifiedMembers := make(
		[]group.MemberIndex,
		len(testResult.dkgResult.Disqualified),
	)

	for _, dq := range testResult.dkgResult.Disqualified {
		memberIndex := group.MemberIndex(uint8(dq))
		actualDisqualifiedMembers = append(actualDisqualifiedMembers, memberIndex)

		disqualificationExpected := containsMemberIndex(
			memberIndex,
			expectedDisqualifiedMembers,
		)

		if !disqualificationExpected {
			t.Errorf(
				"member [%v] should not be marked as disqualified",
				memberIndex,
			)
		}
	}

	for _, memberIndex := range expectedDisqualifiedMembers {
		isDisqualified := containsMemberIndex(
			memberIndex,
			actualDisqualifiedMembers,
		)

		if !isDisqualified {
			t.Errorf(
				"member [%v] should be marked as disqualified",
				memberIndex,
			)
		}
	}
}

func containsMemberIndex(
	index group.MemberIndex,
	indexes []group.MemberIndex,
) bool {
	for _, i := range indexes {
		if i == index {
			return true
		}
	}

	return false
}

// AssertNoInactiveMembers checks there were no inactive members during the
// protocol execution.
func AssertNoInactiveMembers(t *testing.T, testResult *Result) {
	AssertInactiveMembers(t, testResult)
}

// AssertInactiveMembers checks which members were inactive during the protocol
// execution and compares them against expected ones.
func AssertInactiveMembers(
	t *testing.T,
	testResult *Result,
	expectedInactiveMembers ...group.MemberIndex,
) {
	actualInactiveMembers := make(
		[]group.MemberIndex,
		len(testResult.dkgResult.Inactive),
	)

	for _, ia := range testResult.dkgResult.Inactive {
		memberIndex := group.MemberIndex(uint8(ia))
		actualInactiveMembers = append(actualInactiveMembers, memberIndex)

		inactivityExpected := containsMemberIndex(
			group.MemberIndex(memberIndex),
			expectedInactiveMembers,
		)

		if !inactivityExpected {
			t.Errorf(
				"member [%v] should not be marked as inactive",
				memberIndex,
			)
		}
	}

	for _, memberIndex := range expectedInactiveMembers {
		isInactive := containsMemberIndex(
			memberIndex,
			actualInactiveMembers,
		)

		if !isInactive {
			t.Errorf(
				"member [%v] should be marked as inactive",
				memberIndex,
			)
		}
	}
}

// AssertSamePublicKey checks if all members of the group generated the same
// group public key during DKG.
func AssertSamePublicKey(t *testing.T, testResult *Result) {
	for _, signer := range testResult.signers {
		testutils.AssertBytesEqual(
			t,
			testResult.dkgResult.GroupPublicKey,
			signer.GroupPublicKeyBytes(),
		)
	}
}

// AssertValidGroupPublicKey checks if the generated group public key is valid.
func AssertValidGroupPublicKey(t *testing.T, testResult *Result) {
	_, err := altbn128.DecompressToG2(testResult.dkgResult.GroupPublicKey)
	if err != nil {
		t.Errorf("invalid group public key: [%v]", err)
	}
}

// AssertResultSupportingMembers checks which particular members
// actually support the final result with their signature.
func AssertResultSupportingMembers(
	t *testing.T,
	testResult *Result,
	expectedSupportingMembers ...group.MemberIndex,
) {
	actualSupportingMembers := make(
		[]group.MemberIndex,
		len(testResult.dkgResultSignatures),
	)
	for memberIndex := range testResult.dkgResultSignatures {
		actualSupportingMembers = append(actualSupportingMembers, memberIndex)

		isSupportingExpected := containsMemberIndex(
			memberIndex,
			expectedSupportingMembers,
		)

		if !isSupportingExpected {
			t.Errorf(
				"member [%v] should not support the result",
				memberIndex,
			)
		}
	}

	for _, memberIndex := range expectedSupportingMembers {
		isSupporting := containsMemberIndex(
			memberIndex,
			actualSupportingMembers,
		)

		if !isSupporting {
			t.Errorf(
				"member [%v] should support the result",
				memberIndex,
			)
		}
	}
}
