package dkg

import (
	"testing"

	"github.com/keep-network/keep-core/pkg/altbn128"
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
	"github.com/keep-network/keep-core/pkg/internal/testutils"
)

func assertSuccessfulSignersCount(
	t *testing.T,
	result *dkgTestResult,
	expectedCount int,
) {
	if len(result.signers) != expectedCount {
		t.Errorf(
			"Unexpected number of successful signers\nExpected: [%v]\nActual:   [%v]",
			expectedCount,
			len(result.signers),
		)
	}
}

func assertMemberFailuresCount(
	t *testing.T,
	result *dkgTestResult,
	expectedCount int,
) {
	if len(result.memberFailures) != expectedCount {
		t.Errorf(
			"Unexpected number of member failures\nExpected: [%v]\nActual:   [%v]",
			expectedCount,
			len(result.memberFailures),
		)
	}
}

func assertSamePublicKey(t *testing.T, result *dkgTestResult) {
	for _, signer := range result.signers {
		testutils.AssertBytesEqual(
			t,
			result.result.GroupPublicKey,
			signer.GroupPublicKeyBytes(),
		)
	}
}

func assertNoDisqualifiedMembers(t *testing.T, result *dkgTestResult) {
	disqualifiedMemberByte := byte(0x01)

	for i, dq := range result.result.Disqualified {
		if dq == disqualifiedMemberByte {
			t.Errorf("Member [%v] has been disqualified", i)
		}
	}
}

func assertNoInactiveMembers(t *testing.T, result *dkgTestResult) {
	assertInactiveMembers(t, result)
}

func assertInactiveMembers(
	t *testing.T,
	result *dkgTestResult,
	expectedInactive ...group.MemberIndex,
) {
	inactiveMemberByte := byte(0x01)
	activeMemberByte := byte(0x00)

	containsIndex := func(
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

	for i, ia := range result.result.Inactive {
		index := i + 1 // member indexes starts from 1
		inactiveExpected := containsIndex(group.MemberIndex(index), expectedInactive)

		if ia == inactiveMemberByte && !inactiveExpected {
			t.Errorf("Member [%v] has been marked as inactive", index)
		} else if ia == activeMemberByte && inactiveExpected {
			t.Errorf("Member [%v] has not been marked as inactive", index)
		}
	}
}

func assertValidGroupPublicKey(t *testing.T, result *dkgTestResult) {
	_, err := altbn128.DecompressToG2(result.result.GroupPublicKey)
	if err != nil {
		t.Errorf("Invalid group public key: [%v]", err)
	}
}
