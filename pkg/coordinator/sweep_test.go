package coordinator_test

import (
	"github.com/keep-network/keep-core/pkg/maintainer/wallet"
	"testing"

	"github.com/go-test/deep"
	"github.com/ipfs/go-log"

	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/coordinator/internal/test"
	"github.com/keep-network/keep-core/pkg/tbtc"
)

func TestProposeDepositsSweep(t *testing.T) {
	log.SetLogLevel("*", "DEBUG")

	scenarios, err := test.LoadProposeSweepTestScenario()
	if err != nil {
		t.Fatal(err)
	}

	for _, scenario := range scenarios {
		t.Run(scenario.Title, func(t *testing.T) {
			tbtcChain := newLocalTbtcChain()
			btcChain := newLocalBitcoinChain()

			// Chain setup.
			tbtcChain.setDepositParameters(0, 0, scenario.DepositTxMaxFee, 0)

			for _, deposit := range scenario.Deposits {
				tbtcChain.addPastDepositRevealedEvent(
					&tbtc.DepositRevealedEventFilter{
						StartBlock:          deposit.RevealBlock,
						EndBlock:            &deposit.RevealBlock,
						WalletPublicKeyHash: [][20]byte{scenario.WalletPublicKeyHash}},
					&tbtc.DepositRevealedEvent{
						WalletPublicKeyHash: scenario.WalletPublicKeyHash,
						FundingTxHash:       deposit.FundingTxHash,
						FundingOutputIndex:  deposit.FundingOutputIndex,
					},
				)
				btcChain.setTransaction(deposit.FundingTxHash, &bitcoin.Transaction{})
				btcChain.setTransactionConfirmations(deposit.FundingTxHash, tbtc.DepositSweepRequiredFundingTxConfirmations)
			}

			tbtcChain.setDepositSweepProposalValidationResult(scenario.ExpectedDepositSweepProposal, nil, true)
			btcChain.setEstimateSatPerVByteFee(1, scenario.EstimateSatPerVByteFee)

			// Test execution.
			err := wallet.ProposeDepositsSweep(
				nil, // TODO: Set correct chain.
				btcChain,
				scenario.WalletPublicKeyHash,
				scenario.SweepTxFee,
				scenario.DepositsSweepDetails(),
				false,
			)

			if err != nil {
				t.Fatal(err)
			}

			if diff := deep.Equal(
				tbtcChain.depositSweepProposals,
				[]*tbtc.DepositSweepProposal{scenario.ExpectedDepositSweepProposal},
			); diff != nil {
				t.Errorf("invalid deposits: %v", diff)
			}
		})
	}
}
