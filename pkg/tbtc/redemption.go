package tbtc

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/chain"
	"time"
)

// RedemptionRequest represents a tBTC redemption request.
type RedemptionRequest struct {
	// Redeemer is the redeemer's address on the host chain.
	Redeemer chain.Address
	// RedeemerOutputScript is the output script the redeemed Bitcoin funds are
	// locked to. This field is not prepended with the byte-length of the script.
	RedeemerOutputScript []byte
	// RequestedAmount is the TBTC amount (in satoshi) requested for redemption.
	RequestedAmount uint64
	// TreasuryFee is the treasury TBTC fee (in satoshi) at the moment of
	// request creation.
	TreasuryFee uint64
	// TxMaxFee is the maximum value of the per-redemption BTC tx fee (in satoshi)
	// that can be incurred by this request, determined at the moment of
	// request creation.
	TxMaxFee uint64
	// RequestedAt is the time the request was created at.
	RequestedAt time.Time
}

// assembleRedemptionTransaction constructs an unsigned redemption Bitcoin
// transaction.
//
// Regarding input arguments, the walletMainUtxo parameter is mandatory and
// must be set accordingly. The requests slice must contain at least one element.
// The fee argument is not validated in any way so must be chosen with respect
// to the system limitations.
//
// The resulting bitcoin.TransactionBuilder instance holds all the data
// necessary to sign the transaction and obtain a bitcoin.Transaction instance
// ready to be spread across the Bitcoin network.
func assembleRedemptionTransaction(
	bitcoinChain bitcoin.Chain,
	walletPublicKey *ecdsa.PublicKey,
	walletMainUtxo *bitcoin.UnspentTransactionOutput,
	requests []*RedemptionRequest,
	fee int64,
) (*bitcoin.TransactionBuilder, error) {
	if walletMainUtxo == nil {
		return nil, fmt.Errorf("wallet main UTXO is required")
	}

	if len(requests) < 1 {
		return nil, fmt.Errorf("at least one redemption request is required")
	}

	builder := bitcoin.NewTransactionBuilder(bitcoinChain)

	err := builder.AddPublicKeyHashInput(walletMainUtxo)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot add input pointing to wallet main UTXO: [%v]",
			err,
		)
	}

	redemptionsCount := int64(len(requests))

	feePerRedemptionRemainder := fee % redemptionsCount
	feePerRedemption := (fee - feePerRedemptionRemainder) / redemptionsCount

	redemptionOutputs := make([]*bitcoin.TransactionOutput, redemptionsCount)
	totalRedemptionOutputsValue := int64(0)

	// Build a list of redemption outputs based on the provided redemption
	// requests but do not add them to the transaction builder yet.
	// We want to put the change output (constructed in the next step) at the
	// first place in order to make its position predictable and leverage
	// some SPV proof cost optimizations made in the Bridge implementation.
	// That means no outputs can be added to the builder until the change
	// output is there.
	for i, request := range requests {
		// The redeemable amount for a redemption request is the difference
		// between the requested amount and treasury fee computed upon
		// request creation.
		redeemableAmount := int64(request.RequestedAmount - request.TreasuryFee)
		// The actual value of the redemption output is the difference between
		// the request's redeemable amount and fee per redemption.
		redemptionOutputValue := redeemableAmount - feePerRedemption
		// Make the last redemption incur the fee remainder.
		if i == len(requests)-1 {
			redemptionOutputValue -= feePerRedemptionRemainder
		}

		totalRedemptionOutputsValue += redemptionOutputValue

		redemptionOutputs[i] = &bitcoin.TransactionOutput{
			Value:           redemptionOutputValue,
			PublicKeyScript: request.RedeemerOutputScript,
		}
	}

	// We know that the total fee of a Bitcoin transaction is the difference
	// between the sum of inputs and the sum of outputs. In the case of a
	// redemption transaction, that translates to the following formula:
	// fee = main_utxo_input_value - (redemption_outputs_value + change_value)
	// That means we can calculate the change's value using:
	// change_value = main_utxo_input_value - redemption_outputs_value - fee
	changeOutputValue := builder.TotalInputsValue() -
		totalRedemptionOutputsValue -
		fee

	if changeOutputValue > 0 {
		changeOutputScript, err := bitcoin.PayToWitnessPublicKeyHash(
			bitcoin.PublicKeyHash(walletPublicKey),
		)
		if err != nil {
			return nil, fmt.Errorf(
				"cannot compute change output script: [%v]",
				err,
			)
		}

		builder.AddOutput(&bitcoin.TransactionOutput{
			Value:           changeOutputValue,
			PublicKeyScript: changeOutputScript,
		})
	}

	for _, redemptionOutput := range redemptionOutputs {
		builder.AddOutput(redemptionOutput)
	}

	return builder, nil
}
