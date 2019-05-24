package dkg

import (
	"testing"

	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
	"github.com/keep-network/keep-core/pkg/internal/testutils"
)

func assertSignersCount(
	t *testing.T,
	signers []*ThresholdSigner,
	expectedCount int,
) {
	if len(signers) != expectedCount {
		t.Errorf(
			"Unexpected number of signers\nExpected: [%v]\nActual:   [%v]",
			expectedCount,
			len(signers),
		)
	}
}

func assertSamePublicKey(
	t *testing.T,
	result *relaychain.DKGResult,
	signers []*ThresholdSigner,
) {
	for _, signer := range signers {
		testutils.AssertBytesEqual(
			t,
			result.GroupPublicKey,
			signer.GroupPublicKeyBytes(),
		)
	}
}

func filterOutMisbehavingSigners(
	signers []*ThresholdSigner,
	misbehavingSignersIDs ...group.MemberIndex,
) []*ThresholdSigner {
	var honestSigners []*ThresholdSigner
	for _, signer := range signers {
		isMisbehaving := false
		for _, misbehavingID := range misbehavingSignersIDs {
			if signer.MemberID() == misbehavingID {
				isMisbehaving = true
				break
			}
		}
		if !isMisbehaving {
			honestSigners = append(honestSigners, signer)
		}
	}
	return honestSigners
}
