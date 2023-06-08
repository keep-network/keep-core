package tbtc

import (
	"testing"

	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/internal/tbtctest"
	"github.com/keep-network/keep-core/pkg/internal/testutils"
)

func TestAssembleRedemptionTransaction(t *testing.T) {
	scenarios, err := tbtctest.LoadRedemptionTestScenarios()
	if err != nil {
		t.Fatal(err)
	}

	for _, scenario := range scenarios {
		t.Run(scenario.Title, func(t *testing.T) {
			bitcoinChain := newLocalBitcoinChain()

			err := bitcoinChain.BroadcastTransaction(scenario.InputTransaction)
			if err != nil {
				t.Fatal(err)
			}

			requests := make([]*RedemptionRequest, len(scenario.RedemptionRequests))
			for i, r := range scenario.RedemptionRequests {
				requests[i] = &RedemptionRequest{
					Redeemer:             r.Redeemer,
					RedeemerOutputScript: r.RedeemerOutputScript,
					RequestedAmount:      r.RequestedAmount,
					TreasuryFee:          r.TreasuryFee,
					TxMaxFee:             r.TxMaxFee,
					RequestedAt:          r.RequestedAt,
				}
			}

			feeDistribution := func(requests []*RedemptionRequest) []int64 {
				return scenario.FeeShares
			}

			builder, err := assembleRedemptionTransaction(
				bitcoinChain,
				scenario.WalletPublicKey,
				scenario.WalletMainUtxo,
				requests,
				feeDistribution,
				ChangeLastRedemption,
			)
			if err != nil {
				t.Fatal(err)
			}

			sigHashes, err := builder.ComputeSignatureHashes()
			if err != nil {
				t.Fatal(err)
			}

			testutils.AssertIntsEqual(
				t,
				"sighash count",
				1,
				len(sigHashes),
			)

			testutils.AssertBigIntsEqual(
				t,
				"sighash",
				scenario.ExpectedSigHash,
				sigHashes[0],
			)

			transaction, err := builder.AddSignatures(
				[]*bitcoin.SignatureContainer{scenario.Signature},
			)
			if err != nil {
				t.Fatal(err)
			}

			testutils.AssertBytesEqual(
				t,
				scenario.ExpectedRedemptionTransaction.Serialize(),
				transaction.Serialize(),
			)
			testutils.AssertStringsEqual(
				t,
				"redemption transaction hash",
				scenario.ExpectedRedemptionTransactionHash.Hex(bitcoin.InternalByteOrder),
				transaction.Hash().Hex(bitcoin.InternalByteOrder),
			)
			testutils.AssertStringsEqual(
				t,
				"redemption transaction witness hash",
				scenario.ExpectedRedemptionTransactionWitnessHash.Hex(bitcoin.InternalByteOrder),
				transaction.WitnessHash().Hex(bitcoin.InternalByteOrder),
			)
		})
	}
}
