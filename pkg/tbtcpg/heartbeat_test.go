package tbtcpg

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/keep-network/keep-core/internal/testutils"
	"github.com/keep-network/keep-core/pkg/tbtc"
)

func TestHeartbeatTask_Run(t *testing.T) {
	tests := map[string]struct {
		validationResult bool
		expectedProposal tbtc.CoordinationProposal
		expectedOk       bool
		expectedErr      error
	}{
		"valid proposal": {
			validationResult: true,
			expectedProposal: &tbtc.HeartbeatProposal{
				Message: [16]byte{
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xe0, 0xd7, 0x5a, 0xec, 0xd2, 0x9e, 0x5b, 0xca,
				},
			},
			expectedOk:  true,
			expectedErr: nil,
		},
		"invalid proposal": {
			validationResult: false,
			expectedProposal: nil,
			expectedOk:       false,
			expectedErr: fmt.Errorf(
				"failed to verify heartbeat proposal: [validation failed]",
			),
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			tbtcChain := NewLocalChain()
			blockCounter := NewMockBlockCounter()

			blockCounter.SetCurrentBlock(900)
			tbtcChain.SetBlockCounter(blockCounter)

			tbtcChain.SetHeartbeatProposalValidationResult(
				&tbtc.HeartbeatProposal{
					Message: [16]byte{
						0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
						0xe0, 0xd7, 0x5a, 0xec, 0xd2, 0x9e, 0x5b, 0xca,
					},
				},
				test.validationResult,
			)

			walletPublicKeyHash := [20]byte{0x01, 0x02}

			task := NewHeartbeatTask(tbtcChain)

			proposal, ok, err := task.Run(
				&tbtc.CoordinationProposalRequest{
					// Set only relevant fields.
					WalletPublicKeyHash: walletPublicKeyHash,
				},
			)

			if !reflect.DeepEqual(test.expectedErr, err) {
				t.Errorf(
					"unexpected error\nexpected: [%v]\nactual:   [%v]",
					test.expectedErr,
					err,
				)
			}

			testutils.AssertBoolsEqual(t, "boolean flag", test.expectedOk, ok)

			if !reflect.DeepEqual(test.expectedProposal, proposal) {
				t.Errorf(
					"unexpected proposal\nexpected: [%v]\nactual:   [%v]",
					test.expectedProposal,
					proposal,
				)
			}
		})
	}
}
