package wallet_test

import (
	"encoding/hex"
	"github.com/go-test/deep"
	"github.com/keep-network/keep-core/internal/testutils"
	"github.com/keep-network/keep-core/pkg/bitcoin"
	walletmtr "github.com/keep-network/keep-core/pkg/maintainer/wallet"
	"github.com/keep-network/keep-core/pkg/maintainer/wallet/internal/test"
	"github.com/keep-network/keep-core/pkg/tbtc"
	"math/big"
	"testing"
)

// Test based on example testnet redemption transaction:
// https://live.blockcypher.com/btc-testnet/tx/2724545276df61f43f1e92c4b9f1dd3c9109595c022dbd9dc003efbad8ded38b
func TestEstimateRedemptionFee(t *testing.T) {
	fromHex := func(hexString string) []byte {
		bytes, err := hex.DecodeString(hexString)
		if err != nil {
			t.Fatal(err)
		}
		return bytes
	}

	btcChain := walletmtr.NewLocalBitcoinChain()
	btcChain.SetEstimateSatPerVByteFee(1, 16)

	redeemersOutputScripts := []bitcoin.Script{
		fromHex("76a9142cd680318747b720d67bf4246eb7403b476adb3488ac"),                   // P2PKH
		fromHex("0014e6f9d74726b19b75f16fe1e9feaec048aa4fa1d0"),                         // P2WPKH
		fromHex("a914011beb6fb8499e075a57027fb0a58384f2d3f78487"),                       // P2SH
		fromHex("0020ef0b4d985752aa5ef6243e4c6f6bebc2a007e7d671ef27d4b1d0db8dcc93bc1c"), // P2WSH
	}

	actualFee, err := walletmtr.EstimateRedemptionFee(btcChain, redeemersOutputScripts)
	if err != nil {
		t.Fatal(err)
	}

	expectedFee := 4000 // transactionVirtualSize * satPerVByteFee = 250 * 16 = 4000
	testutils.AssertIntsEqual(t, "fee", expectedFee, int(actualFee))
}

func TestFindPendingRedemptions(t *testing.T) {
	scenarios, err := test.LoadFindPendingRedemptionsTestScenario()
	if err != nil {
		t.Fatal(err)
	}

	for _, scenario := range scenarios {
		t.Run(scenario.Title, func(t *testing.T) {
			tbtcChain := walletmtr.NewLocalChain()

			// Set the average block time enforced by the scenario.
			tbtcChain.SetAverageBlockTime(scenario.ChainParameters.AverageBlockTime)

			// Set the scenario's current block using a mock block counter.
			// This is needed to build a proper filter for the
			// `PastRedemptionRequestedEvents` call.
			blockCounter := walletmtr.NewMockBlockCounter()
			blockCounter.SetCurrentBlock(scenario.ChainParameters.CurrentBlock)
			tbtcChain.SetBlockCounter(blockCounter)

			// Set relevant governable parameters based on values provided by
			// the scenario.
			tbtcChain.SetRedemptionParameters(
				0,
				0,
				0,
				0,
				scenario.ChainParameters.RequestTimeout,
				nil,
				0,
			)
			tbtcChain.SetRedemptionRequestMinAge(scenario.ChainParameters.RequestMinAge)

			requestTimeoutBlocks := uint64(scenario.ChainParameters.RequestTimeout) /
				uint64(scenario.ChainParameters.AverageBlockTime.Seconds())

			// Record scenario wallets to the local chain.
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

			// Record scenario pending redemptions to the local chain.
			for _, pendingRedemption := range scenario.PendingRedemptions {
				// Record the corresponding event. Set only relevant fields.
				err = tbtcChain.AddPastRedemptionRequestedEvent(
					&tbtc.RedemptionRequestedEventFilter{
						// Remember about including the constant factor
						// of 1000 blocks.
						StartBlock:          scenario.ChainParameters.CurrentBlock - requestTimeoutBlocks - 1000,
						WalletPublicKeyHash: [][20]byte{pendingRedemption.WalletPublicKeyHash},
					},
					&tbtc.RedemptionRequestedEvent{
						WalletPublicKeyHash:  pendingRedemption.WalletPublicKeyHash,
						RedeemerOutputScript: pendingRedemption.RedeemerOutputScript,
						RequestedAmount:      pendingRedemption.RequestedAmount,
						BlockNumber:          pendingRedemption.RequestBlock,
					},
				)

				// Record the corresponding request object. Set only relevant fields.
				tbtcChain.SetPendingRedemptionRequest(
					pendingRedemption.WalletPublicKeyHash,
					&tbtc.RedemptionRequest{
						RedeemerOutputScript: pendingRedemption.RedeemerOutputScript,
						RequestedAmount:      pendingRedemption.RequestedAmount,
						RequestedAt:          pendingRedemption.RequestedAt,
					},
				)
			}

			walletsPendingRedemptions, err := walletmtr.FindPendingRedemptions(
				tbtcChain,
				scenario.Filter,
			)
			if err != nil {
				t.Fatal(err)
			}

			if diff := deep.Equal(
				scenario.ExpectedWalletsPendingRedemptions,
				walletsPendingRedemptions,
			); diff != nil {
				t.Errorf("invalid wallets pending redemptions: %v", diff)
			}
		})
	}
}

func TestProposeRedemption(t *testing.T) {
	fromHex := func(hexString string) []byte {
		bytes, err := hex.DecodeString(hexString)
		if err != nil {
			t.Fatal(err)
		}
		return bytes
	}

	var walletPublicKeyHash [20]byte
	copy(walletPublicKeyHash[:], fromHex(""))

	redeemersOutputScripts := []bitcoin.Script{
		fromHex("00140000000000000000000000000000000000000001"),
		fromHex("00140000000000000000000000000000000000000002"),
	}

	var tests = map[string]struct {
		fee              int64
		expectedProposal *tbtc.RedemptionProposal
	}{
		"fee provided": {
			fee: 10000,
			expectedProposal: &tbtc.RedemptionProposal{
				WalletPublicKeyHash:    walletPublicKeyHash,
				RedeemersOutputScripts: redeemersOutputScripts,
				RedemptionTxFee:        big.NewInt(10000),
			},
		},
		"fee estimated": {
			fee: 0, // trigger fee estimation
			expectedProposal: &tbtc.RedemptionProposal{
				WalletPublicKeyHash:    walletPublicKeyHash,
				RedeemersOutputScripts: redeemersOutputScripts,
				RedemptionTxFee:        big.NewInt(4300),
			},
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			tbtcChain := walletmtr.NewLocalChain()
			btcChain := walletmtr.NewLocalBitcoinChain()

			btcChain.SetEstimateSatPerVByteFee(1, 25)

			for _, script := range redeemersOutputScripts {
				tbtcChain.SetPendingRedemptionRequest(
					walletPublicKeyHash,
					&tbtc.RedemptionRequest{
						RedeemerOutputScript: script,
					},
				)
			}

			err := tbtcChain.SetRedemptionProposalValidationResult(
				test.expectedProposal,
				true,
			)
			if err != nil {
				t.Fatal(err)
			}

			err = walletmtr.ProposeRedemption(
				tbtcChain,
				btcChain,
				walletPublicKeyHash,
				test.fee,
				redeemersOutputScripts,
				false,
			)
			if err != nil {
				t.Fatal(err)
			}

			if diff := deep.Equal(
				tbtcChain.RedemptionProposals(),
				[]*tbtc.RedemptionProposal{test.expectedProposal},
			); diff != nil {
				t.Errorf("invalid deposits: %v", diff)
			}
		})
	}
}
