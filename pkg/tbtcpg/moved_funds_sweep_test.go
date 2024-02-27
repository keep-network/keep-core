package tbtcpg_test

import (
	"testing"

	"github.com/keep-network/keep-core/internal/testutils"
	"github.com/keep-network/keep-core/pkg/tbtcpg"
)

func TestEstimateMovedFundsSweepFee(t *testing.T) {
	var tests = map[string]struct {
		sweepTxMaxTotalFee uint64
		hasMainUtxo        bool
		expectedFee        uint64
		expectedError      error
	}{
		"estimated fee correct, one input": {
			sweepTxMaxTotalFee: 3000,
			hasMainUtxo:        false,
			expectedFee:        1760,
			expectedError:      nil,
		},
		"estimated fee correct, two inputs": {
			sweepTxMaxTotalFee: 3000,
			hasMainUtxo:        true,
			expectedFee:        2848,
			expectedError:      nil,
		},
		"estimated fee too high": {
			sweepTxMaxTotalFee: 2500,
			hasMainUtxo:        true,
			expectedFee:        0,
			expectedError:      tbtcpg.ErrSweepTxFeeTooHigh,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			btcChain := tbtcpg.NewLocalBitcoinChain()
			btcChain.SetEstimateSatPerVByteFee(1, 16)

			actualFee, err := tbtcpg.EstimateMovedFundsSweepFee(
				btcChain,
				test.hasMainUtxo,
				test.sweepTxMaxTotalFee,
			)

			testutils.AssertUintsEqual(
				t,
				"fee",
				test.expectedFee,
				uint64(actualFee),
			)

			testutils.AssertAnyErrorInChainMatchesTarget(
				t,
				test.expectedError,
				err,
			)
		})
	}
}
