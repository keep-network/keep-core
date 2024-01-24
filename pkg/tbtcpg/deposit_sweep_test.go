package tbtcpg_test

import (
	"reflect"
	"testing"

	"github.com/go-test/deep"
	"github.com/ipfs/go-log"
	"github.com/keep-network/keep-core/internal/testutils"
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/tbtc"
	"github.com/keep-network/keep-core/pkg/tbtcpg"
	"github.com/keep-network/keep-core/pkg/tbtcpg/internal/test"
)

func TestDepositSweepTask_FindDepositsToSweep(t *testing.T) {
	err := log.SetLogLevel("*", "DEBUG")
	if err != nil {
		t.Fatal(err)
	}

	scenarios, err := test.LoadFindDepositsToSweepTestScenario()
	if err != nil {
		t.Fatal(err)
	}

	for _, scenario := range scenarios {
		t.Run(scenario.Title, func(t *testing.T) {
			tbtcChain := tbtcpg.NewLocalChain()
			btcChain := tbtcpg.NewLocalBitcoinChain()

			// Chain setup.
			for _, deposit := range scenario.Deposits {
				tbtcChain.SetDepositRequest(
					deposit.FundingTxHash,
					deposit.FundingOutputIndex,
					&tbtc.DepositChainRequest{SweptAt: deposit.SweptAt},
				)
				btcChain.SetTransaction(deposit.FundingTxHash, deposit.FundingTx)
				btcChain.SetTransactionConfirmations(
					deposit.FundingTxHash,
					deposit.FundingTxConfirmations,
				)

				err := tbtcChain.AddPastDepositRevealedEvent(
					&tbtc.DepositRevealedEventFilter{WalletPublicKeyHash: [][20]byte{deposit.WalletPublicKeyHash}},
					&tbtc.DepositRevealedEvent{
						BlockNumber:         deposit.RevealBlockNumber,
						WalletPublicKeyHash: deposit.WalletPublicKeyHash,
						FundingTxHash:       deposit.FundingTxHash,
						FundingOutputIndex:  deposit.FundingOutputIndex,
					},
				)
				if err != nil {
					t.Fatal(err)
				}
			}

			task := tbtcpg.NewDepositSweepTask(tbtcChain, btcChain)

			// Test execution.
			actualDeposits, err := task.FindDepositsToSweep(
				&testutils.MockLogger{},
				scenario.WalletPublicKeyHash,
				scenario.MaxNumberOfDeposits,
			)

			if err != nil {
				t.Fatal(err)
			}

			if diff := deep.Equal(
				scenario.ExpectedUnsweptDeposits,
				actualDeposits,
			); diff != nil {
				t.Errorf("invalid deposits: %v", diff)
			}
		})
	}
}

func TestDepositSweepTask_ProposeDepositsSweep(t *testing.T) {
	err := log.SetLogLevel("*", "DEBUG")
	if err != nil {
		t.Fatal(err)
	}

	scenarios, err := test.LoadProposeSweepTestScenario()
	if err != nil {
		t.Fatal(err)
	}

	for _, scenario := range scenarios {
		t.Run(scenario.Title, func(t *testing.T) {
			tbtcChain := tbtcpg.NewLocalChain()
			btcChain := tbtcpg.NewLocalBitcoinChain()

			// Chain setup.
			tbtcChain.SetDepositParameters(0, 0, scenario.DepositTxMaxFee, 0)

			for _, deposit := range scenario.Deposits {
				err := tbtcChain.AddPastDepositRevealedEvent(
					&tbtc.DepositRevealedEventFilter{
						StartBlock:          deposit.RevealBlock,
						EndBlock:            &deposit.RevealBlock,
						WalletPublicKeyHash: [][20]byte{scenario.WalletPublicKeyHash},
					},
					&tbtc.DepositRevealedEvent{
						WalletPublicKeyHash: scenario.WalletPublicKeyHash,
						FundingTxHash:       deposit.FundingTxHash,
						FundingOutputIndex:  deposit.FundingOutputIndex,
					},
				)
				if err != nil {
					t.Fatal(err)
				}

				tbtcChain.SetDepositRequest(
					deposit.FundingTxHash,
					deposit.FundingOutputIndex,
					&tbtc.DepositChainRequest{
						// Set only relevant fields.
						ExtraData: nil,
					},
				)

				btcChain.SetTransaction(deposit.FundingTxHash, &bitcoin.Transaction{})
				btcChain.SetTransactionConfirmations(deposit.FundingTxHash, tbtc.DepositSweepRequiredFundingTxConfirmations)
			}

			if scenario.ExpectedDepositSweepProposal != nil {
				err := tbtcChain.SetDepositSweepProposalValidationResult(
					scenario.WalletPublicKeyHash,
					scenario.ExpectedDepositSweepProposal,
					nil,
					true,
				)
				if err != nil {
					t.Fatal(err)
				}
			}

			btcChain.SetEstimateSatPerVByteFee(1, scenario.EstimateSatPerVByteFee)

			task := tbtcpg.NewDepositSweepTask(tbtcChain, btcChain)

			// Test execution.
			proposal, err := task.ProposeDepositsSweep(
				&testutils.MockLogger{},
				scenario.WalletPublicKeyHash,
				scenario.DepositsReferences(),
				scenario.SweepTxFee,
			)

			if !reflect.DeepEqual(scenario.ExpectedErr, err) {
				t.Errorf(
					"unexpected error\n"+
						"expected: [%+v]\n"+
						"actual:   [%+v]",
					scenario.ExpectedErr,
					err,
				)
			}

			var actualDepositSweepProposals []*tbtc.DepositSweepProposal
			if proposal != nil {
				actualDepositSweepProposals = append(actualDepositSweepProposals, proposal)
			}

			var expectedDepositSweepProposals []*tbtc.DepositSweepProposal
			if p := scenario.ExpectedDepositSweepProposal; p != nil {
				expectedDepositSweepProposals = append(expectedDepositSweepProposals, p)
			}

			if diff := deep.Equal(
				actualDepositSweepProposals,
				expectedDepositSweepProposals,
			); diff != nil {
				t.Errorf("invalid deposit sweep proposal: %v", diff)
			}
		})
	}
}
