package wallet_test

import (
	"github.com/go-test/deep"
	"github.com/ipfs/go-log"
	"github.com/keep-network/keep-core/internal/hexutils"
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/coordinator"
	mtrwallet "github.com/keep-network/keep-core/pkg/maintainer/wallet"
	"github.com/keep-network/keep-core/pkg/maintainer/wallet/internal/test"
	"github.com/keep-network/keep-core/pkg/tbtc"
	"testing"
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
			tbtcChain := mtrwallet.NewLocalChain()
			btcChain := mtrwallet.NewLocalBitcoinChain()

			expectedWallet := scenario.ExpectedWalletPublicKeyHash

			// Chain setup.
			for _, wallet := range scenario.Wallets {
				err := tbtcChain.AddPastNewWalletRegisteredEvent(
					nil,
					&coordinator.NewWalletRegisteredEvent{
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
			actualWallet, actualDeposits, err := mtrwallet.FindDepositsToSweep(
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
			tbtcChain := mtrwallet.NewLocalChain()
			btcChain := mtrwallet.NewLocalBitcoinChain()

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

			err := tbtcChain.SetDepositSweepProposalValidationResult(
				scenario.ExpectedDepositSweepProposal,
				nil,
				true,
			)
			if err != nil {
				t.Fatal(err)
			}

			btcChain.SetEstimateSatPerVByteFee(1, scenario.EstimateSatPerVByteFee)

			// Test execution.
			err = mtrwallet.ProposeDepositsSweep(
				tbtcChain,
				btcChain,
				scenario.WalletPublicKeyHash,
				scenario.SweepTxFee,
				scenario.DepositsReferences(),
				false,
			)
			if err != nil {
				t.Fatal(err)
			}

			if diff := deep.Equal(
				tbtcChain.DepositSweepProposals(),
				[]*tbtc.DepositSweepProposal{scenario.ExpectedDepositSweepProposal},
			); diff != nil {
				t.Errorf("invalid deposits: %v", diff)
			}
		})
	}
}
