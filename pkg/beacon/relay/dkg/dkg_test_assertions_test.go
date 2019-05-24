package dkg

import (
	"testing"

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

func assertSamePublicKey(
	t *testing.T,
	result *dkgTestResult,
) {
	for _, signer := range result.signers {
		testutils.AssertBytesEqual(
			t,
			result.result.GroupPublicKey,
			signer.GroupPublicKeyBytes(),
		)
	}
}
