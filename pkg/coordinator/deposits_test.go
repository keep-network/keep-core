package coordinator_test

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/ipfs/go-log"

	"github.com/keep-network/keep-core/pkg/coordinator"
	"github.com/keep-network/keep-core/pkg/coordinator/internal/test"

	"github.com/keep-network/keep-core/internal/hexutils"
	"github.com/keep-network/keep-core/pkg/tbtc"
)

func TestFindDepositsToSweep(t *testing.T) {
	log.SetLogLevel("*", "DEBUG")

	scenarios, err := test.LoadFindDepositsToSweepTestScenario()
	if err != nil {
		t.Fatal(err)
	}

	for _, scenario := range scenarios {
		t.Run(scenario.Title, func(t *testing.T) {
			tbtcChain := newLocalTbtcChain()
			btcChain := newLocalBitcoinChain()

			expectedWallet := scenario.ExpectedWalletPublicKeyHash

			// Chain setup.
			for _, wallet := range scenario.Wallets {
				tbtcChain.addPastNewWalletRegisteredEvent(
					nil,
					&tbtc.NewWalletRegisteredEvent{
						WalletPublicKeyHash: wallet.WalletPublicKeyHash,
						BlockNumber:         wallet.RegistrationBlockNumber,
					},
				)

			}

			for _, deposit := range scenario.Deposits {
				tbtcChain.setDepositRequest(
					deposit.FundingTxHash,
					deposit.FundingOutputIndex,
					&tbtc.DepositChainRequest{SweptAt: deposit.SweptAt},
				)
				btcChain.setTransactionConfirmations(
					deposit.FundingTxHash,
					deposit.FundingTxConfirmations,
				)

				tbtcChain.addPastDepositRevealedEvent(
					&tbtc.DepositRevealedEventFilter{WalletPublicKeyHash: [][20]byte{deposit.WalletPublicKeyHash}},
					&tbtc.DepositRevealedEvent{
						BlockNumber:         deposit.RevealBlockNumber,
						WalletPublicKeyHash: deposit.WalletPublicKeyHash,
						FundingTxHash:       deposit.FundingTxHash,
						FundingOutputIndex:  deposit.FundingOutputIndex,
					},
				)
			}

			// Test execution.
			actualWallet, actualDeposits, err := coordinator.FindDepositsToSweep(
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
