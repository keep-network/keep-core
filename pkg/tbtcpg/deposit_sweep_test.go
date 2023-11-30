package tbtcpg_test

import (
	"reflect"
	"testing"

	"github.com/go-test/deep"
	"github.com/ipfs/go-log"
	"github.com/keep-network/keep-core/internal/hexutils"
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/tbtc"
	"github.com/keep-network/keep-core/pkg/tbtcpg"
	"github.com/keep-network/keep-core/pkg/tbtcpg/internal/test"
)

func TestFindDepositsToSweep(t *testing.T) {
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

			expectedWallet := scenario.ExpectedWalletPublicKeyHash

			// Chain setup.
			for _, wallet := range scenario.Wallets {
				err := tbtcChain.AddPastNewWalletRegisteredEvent(
					nil,
					&tbtc.NewWalletRegisteredEvent{
						WalletPublicKeyHash: wallet.WalletPublicKeyHash,
						BlockNumber:         wallet.RegistrationBlockNumber,
					},
				)
				if err != nil {
					t.Fatal(err)
				}
			}

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

			// Test execution.
			actualWallet, actualDeposits, err := tbtcpg.FindDepositsToSweep(
				tbtcChain,
				btcChain,
				scenario.WalletPublicKeyHash,
				scenario.MaxNumberOfDeposits,
			)

			if err != nil {
				t.Fatal(err)
			}

			if actualWallet != expectedWallet {
				t.Errorf(
					"invalid wallet public key hash\nexpected: %s\nactual:   %s",
					hexutils.Encode(expectedWallet[:]),
					hexutils.Encode(actualWallet[:]),
				)
			}

			if diff := deep.Equal(actualDeposits, scenario.ExpectedUnsweptDeposits); diff != nil {
				t.Errorf("invalid deposits: %v", diff)
			}
		})
	}
}

func TestProposeDepositsSweep(t *testing.T) {
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

				btcChain.SetTransaction(deposit.FundingTxHash, &bitcoin.Transaction{})
				btcChain.SetTransactionConfirmations(deposit.FundingTxHash, tbtc.DepositSweepRequiredFundingTxConfirmations)
			}

			if scenario.ExpectedDepositSweepProposal != nil {
				err := tbtcChain.SetDepositSweepProposalValidationResult(
					scenario.ExpectedDepositSweepProposal,
					nil,
					true,
				)
				if err != nil {
					t.Fatal(err)
				}
			}

			btcChain.SetEstimateSatPerVByteFee(1, scenario.EstimateSatPerVByteFee)

			// Test execution.
			err = tbtcpg.ProposeDepositsSweep(
				tbtcChain,
				btcChain,
				scenario.WalletPublicKeyHash,
				scenario.SweepTxFee,
				scenario.DepositsReferences(),
				false,
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

			var expectedDepositSweepProposals []*tbtc.DepositSweepProposal
			if p := scenario.ExpectedDepositSweepProposal; p != nil {
				expectedDepositSweepProposals = append(expectedDepositSweepProposals, p)
			}

			if diff := deep.Equal(
				tbtcChain.DepositSweepProposals(),
				expectedDepositSweepProposals,
			); diff != nil {
				t.Errorf("invalid deposits: %v", diff)
			}
		})
	}
}
