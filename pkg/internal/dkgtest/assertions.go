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
	disqualifiedMemberByte := byte(0x01)

	for i, dq := range testResult.dkgResult.Disqualified {
		if dq == disqualifiedMemberByte {
			t.Errorf("member [%v] has been unexpectedly disqualified", i)
		}
	}
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
	inactiveMemberByte := byte(0x01)
	activeMemberByte := byte(0x00)

	containsMemberIndex := func(
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

	for i, ia := range testResult.dkgResult.Inactive {
		memberIndex := i + 1 // member indexes starts from 1
		inactiveExpected := containsMemberIndex(
			group.MemberIndex(memberIndex),
			expectedInactiveMembers,
		)

		if ia == inactiveMemberByte && !inactiveExpected {
			t.Errorf(
				"member [%v] has been unexpectedly marked as inactive",
				memberIndex,
			)
		} else if ia == activeMemberByte && inactiveExpected {
			t.Errorf(
				"member [%v] has not been unexpectedly marked as inactive",
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
