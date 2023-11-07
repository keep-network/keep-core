package tbtc

import (
	"crypto/ecdsa"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/ipfs/go-log/v2"

	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/chain"
)

const (
	// redemptionProposalConfirmationBlocks determines the block length of the
	// confirmation period on the host chain that is preserved after a
	// redemption proposal submission.
	redemptionProposalConfirmationBlocks = 20
	// redemptionSigningTimeoutSafetyMargin determines the duration of the
	// safety margin that must be preserved between the signing timeout
	// and the timeout of the entire redemption action. This safety
	// margin prevents against the case where signing completes late and there
	// is not enough time to broadcast the redemption transaction properly.
	// In such a case, wallet signatures may leak and make the wallet subject
	// of fraud accusations. Usage of the safety margin ensures there is enough
	// time to perform post-signing steps of the redemption action.
	redemptionSigningTimeoutSafetyMargin = 1 * time.Hour
	// redemptionBroadcastTimeout determines the time window for redemption
	// transaction broadcast. It is guaranteed that at least
	// redemptionSigningTimeoutSafetyMargin is preserved for the broadcast
	// step. However, the happy path for the broadcast step is usually quick
	// and few retries are needed to recover from temporary problems. That
	// said, if the broadcast step does not succeed in a tight timeframe,
	// there is no point to retry for the entire possible time window.
	// Hence, the timeout for broadcast step is set as 25% of the entire
	// time widow determined by redemptionSigningTimeoutSafetyMargin.
	redemptionBroadcastTimeout = redemptionSigningTimeoutSafetyMargin / 4
	// redemptionBroadcastCheckDelay determines the delay that must
	// be preserved between transaction broadcast and the check that ensures
	// the transaction is known on the Bitcoin chain. This delay is needed
	// as spreading the transaction over the Bitcoin network takes time.
	redemptionBroadcastCheckDelay = 1 * time.Minute
)

// RedemptionTransactionShape is an enum describing the shape of
// a Bitcoin redemption transaction.
type RedemptionTransactionShape uint8

const (
	// RedemptionChangeFirst is a shape where the change output is the first one
	// in the transaction output vector. This shape makes the change's position
	// fixed and leverages some SPV proof cost optimizations made in the Bridge
	// implementation.
	RedemptionChangeFirst RedemptionTransactionShape = iota
	// RedemptionChangeLast is a shape where the change output is the last one
	// in the transaction output vector.
	RedemptionChangeLast
)

// RedemptionRequest represents a tBTC redemption request.
type RedemptionRequest struct {
	// Redeemer is the redeemer's address on the host chain.
	Redeemer chain.Address
	// RedeemerOutputScript is the output script the redeemed Bitcoin funds are
	// locked to. As stated in the bitcoin.Script docstring, this field is not
	// prepended with the byte-length of the script.
	RedeemerOutputScript bitcoin.Script
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

// redemptionAction is a redemption walletAction.
type redemptionAction struct {
	logger   *zap.SugaredLogger
	chain    Chain
	btcChain bitcoin.Chain

	redeemingWallet     wallet
	transactionExecutor *walletTransactionExecutor

	proposal                     *RedemptionProposal
	proposalProcessingStartBlock uint64
	proposalExpiresAt            time.Time

	signingTimeoutSafetyMargin time.Duration
	broadcastTimeout           time.Duration
	broadcastCheckDelay        time.Duration

	feeDistribution  redemptionFeeDistributionFn
	transactionShape RedemptionTransactionShape
}

func newRedemptionAction(
	logger *zap.SugaredLogger,
	chain Chain,
	btcChain bitcoin.Chain,
	redeemingWallet wallet,
	signingExecutor walletSigningExecutor,
	proposal *RedemptionProposal,
	proposalProcessingStartBlock uint64,
	proposalExpiresAt time.Time,
) *redemptionAction {
	transactionExecutor := newWalletTransactionExecutor(
		btcChain,
		redeemingWallet,
		signingExecutor,
	)

	feeDistribution := withRedemptionTotalFee(proposal.RedemptionTxFee.Int64())

	return &redemptionAction{
		logger:                       logger,
		chain:                        chain,
		btcChain:                     btcChain,
		redeemingWallet:              redeemingWallet,
		transactionExecutor:          transactionExecutor,
		proposal:                     proposal,
		proposalProcessingStartBlock: proposalProcessingStartBlock,
		proposalExpiresAt:            proposalExpiresAt,
		signingTimeoutSafetyMargin:   redemptionSigningTimeoutSafetyMargin,
		broadcastTimeout:             redemptionBroadcastTimeout,
		broadcastCheckDelay:          redemptionBroadcastCheckDelay,
		feeDistribution:              feeDistribution,
		transactionShape:             RedemptionChangeFirst,
	}
}

func (ra *redemptionAction) execute() error {
	validateProposalLogger := ra.logger.With(
		zap.String("step", "validateProposal"),
	)

	validatedRequests, err := ValidateRedemptionProposal(
		validateProposalLogger,
		ra.proposal,
		ra.chain,
	)
	if err != nil {
		return fmt.Errorf("validate proposal step failed: [%v]", err)
	}

	walletPublicKeyHash := bitcoin.PublicKeyHash(ra.wallet().publicKey)

	walletMainUtxo, err := DetermineWalletMainUtxo(
		walletPublicKeyHash,
		ra.chain,
		ra.btcChain,
	)
	if err != nil {
		return fmt.Errorf(
			"error while determining wallet's main UTXO: [%v]",
			err,
		)
	}

	// Proposal validation should detect this but let's make a check just
	// in case.
	if walletMainUtxo == nil {
		return fmt.Errorf("redeeming wallet has no main UTXO")
	}

	err = EnsureWalletSyncedBetweenChains(
		walletPublicKeyHash,
		walletMainUtxo,
		ra.chain,
		ra.btcChain,
	)
	if err != nil {
		return fmt.Errorf(
			"error while ensuring wallet state is synced between "+
				"BTC and host chain: [%v]",
			err,
		)
	}

	unsignedRedemptionTx, err := assembleRedemptionTransaction(
		ra.btcChain,
		ra.wallet().publicKey,
		walletMainUtxo,
		validatedRequests,
		ra.feeDistribution,
		ra.transactionShape,
	)
	if err != nil {
		return fmt.Errorf(
			"error while assembling redemption transaction: [%v]",
			err,
		)
	}

	signTxLogger := ra.logger.With(
		zap.String("step", "signTransaction"),
	)

	redemptionTx, err := ra.transactionExecutor.signTransaction(
		signTxLogger,
		unsignedRedemptionTx,
		ra.proposalProcessingStartBlock,
		ra.proposalExpiresAt.Add(-ra.signingTimeoutSafetyMargin),
	)
	if err != nil {
		return fmt.Errorf("sign transaction step failed: [%v]", err)
	}

	broadcastTxLogger := ra.logger.With(
		zap.String("step", "broadcastTransaction"),
		zap.String("redemptionTxHash", redemptionTx.Hash().Hex(bitcoin.ReversedByteOrder)),
	)

	err = ra.transactionExecutor.broadcastTransaction(
		broadcastTxLogger,
		redemptionTx,
		ra.broadcastTimeout,
		ra.broadcastCheckDelay,
	)
	if err != nil {
		return fmt.Errorf("broadcast transaction step failed: [%v]", err)
	}

	return nil
}

// ValidateRedemptionProposal checks the redemption proposal with on-chain
// validation rules.
func ValidateRedemptionProposal(
	validateProposalLogger log.StandardLogger,
	proposal *RedemptionProposal,
	chain interface {
		// GetPendingRedemptionRequest gets the on-chain pending redemption request
		// for the given wallet public key hash and redeemer output script.
		// The returned bool value indicates whether the request was found or not.
		GetPendingRedemptionRequest(
			walletPublicKeyHash [20]byte,
			redeemerOutputScript bitcoin.Script,
		) (*RedemptionRequest, bool, error)

		// ValidateRedemptionProposal validates the given redemption proposal
		// against the chain. Returns an error if the proposal is not valid or
		// nil otherwise.
		ValidateRedemptionProposal(proposal *RedemptionProposal) error
	},
) ([]*RedemptionRequest, error) {
	validateProposalLogger.Infof("calling chain for proposal validation")

	err := chain.ValidateRedemptionProposal(proposal)
	if err != nil {
		return nil, fmt.Errorf("redemption proposal is invalid: [%v]", err)
	}

	validateProposalLogger.Infof(
		"redemption proposal is valid",
	)

	requests := make([]*RedemptionRequest, len(proposal.RedeemersOutputScripts))
	for i, script := range proposal.RedeemersOutputScripts {
		requestDisplayIndex := fmt.Sprintf(
			"%v/%v",
			i+1,
			len(proposal.RedeemersOutputScripts),
		)

		request, found, err := chain.GetPendingRedemptionRequest(
			proposal.WalletPublicKeyHash,
			script,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"cannot get pending redemption request data for request [%v]: [%v]",
				requestDisplayIndex,
				err,
			)
		}
		if !found {
			return nil, fmt.Errorf(
				"request [%v] is not a pending redemption request",
				requestDisplayIndex,
			)
		}

		requests[i] = request
	}

	return requests, nil
}

func (ra *redemptionAction) wallet() wallet {
	return ra.redeemingWallet
}

func (ra *redemptionAction) actionType() WalletActionType {
	return ActionRedemption
}

// redemptionFeeDistributionFn calculates the redemption transaction fee
// distribution for the given redemption requests. The resulting list
// contains the fee shares ordered in the same way as the input requests, i.e.
// the first fee share corresponds to the first request and so on.
type redemptionFeeDistributionFn func([]*RedemptionRequest) []int64

// withRedemptionTotalFee is a fee distribution function that takes a
// total transaction fee and distributes it evenly over all redemption requests.
// If the fee cannot be divided evenly, the last request incurs the remainder.
func withRedemptionTotalFee(totalFee int64) redemptionFeeDistributionFn {
	return func(requests []*RedemptionRequest) []int64 {
		requestsCount := int64(len(requests))
		remainder := totalFee % requestsCount
		feePerRequest := (totalFee - remainder) / requestsCount

		feeShares := make([]int64, requestsCount)
		for i := range requests {
			feeShare := feePerRequest

			if i == len(requests)-1 {
				feeShare += remainder
			}

			feeShares[i] = feeShare
		}

		return feeShares
	}
}

// assembleRedemptionTransaction constructs an unsigned redemption Bitcoin
// transaction.
//
// Regarding input arguments, the requests slice must contain at least one element.
// The fee shares applied to specific requests according to the provided
// feeDistribution function are not validated in any way so must be chosen with
// respect to the system limitations. The shape argument is optional - if not
// provided the RedemptionChangeFirst value is used by default.
//
// The resulting bitcoin.TransactionBuilder instance holds all the data
// necessary to sign the transaction and obtain a bitcoin.Transaction instance
// ready to be spread across the Bitcoin network.
func assembleRedemptionTransaction(
	bitcoinChain bitcoin.Chain,
	walletPublicKey *ecdsa.PublicKey,
	walletMainUtxo *bitcoin.UnspentTransactionOutput,
	requests []*RedemptionRequest,
	feeDistribution redemptionFeeDistributionFn,
	shape ...RedemptionTransactionShape,
) (*bitcoin.TransactionBuilder, error) {
	resolvedShape := RedemptionChangeFirst
	if len(shape) == 1 {
		resolvedShape = shape[0]
	}

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

	// Calculate the transaction fee shares for all redemption requests.
	feeShares := feeDistribution(requests)
	// Helper variable that will hold the total Bitcoin transaction fee.
	totalFee := int64(0)
	// Helper variable that will hold the summarized value of all redemption
	// outputs. The change value will not be counted in here.
	totalRedemptionOutputsValue := int64(0)
	// List that will hold all transaction outputs, i.e. redemption outputs
	// and the possible change output.
	outputs := make([]*bitcoin.TransactionOutput, 0)

	// Create redemption outputs based on the provided redemption requests but
	// do not add them to the transaction builder yet. The builder cannot be
	// filled right now due to the change output that will be constructed in the
	// next step and whose position in the transaction output vector depends on
	// the requested RedemptionTransactionShape.
	for i, request := range requests {
		// The redeemable amount for a redemption request is the difference
		// between the requested amount and treasury fee computed upon
		// request creation.
		redeemableAmount := int64(request.RequestedAmount - request.TreasuryFee)
		// The actual value of the redemption output is the difference between
		// the request's redeemable amount and share of the transaction fee
		// incurred by the given request.
		feeShare := feeShares[i]
		redemptionOutputValue := redeemableAmount - feeShare

		totalFee += feeShare
		totalRedemptionOutputsValue += redemptionOutputValue

		redemptionOutput := &bitcoin.TransactionOutput{
			Value:           redemptionOutputValue,
			PublicKeyScript: request.RedeemerOutputScript,
		}

		outputs = append(outputs, redemptionOutput)
	}

	// We know that the total fee of a Bitcoin transaction is the difference
	// between the sum of inputs and the sum of outputs. In the case of a
	// redemption transaction, that translates to the following formula:
	// fee = main_utxo_input_value - (redemption_outputs_value + change_value)
	// That means we can calculate the change's value using:
	// change_value = main_utxo_input_value - redemption_outputs_value - fee
	changeOutputValue := builder.TotalInputsValue() -
		totalRedemptionOutputsValue -
		totalFee

	// If we can have a non-zero change, construct it.
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

		changeOutput := &bitcoin.TransactionOutput{
			Value:           changeOutputValue,
			PublicKeyScript: changeOutputScript,
		}

		switch resolvedShape {
		case RedemptionChangeFirst:
			outputs = append([]*bitcoin.TransactionOutput{changeOutput}, outputs...)
		case RedemptionChangeLast:
			outputs = append(outputs, changeOutput)
		default:
			panic("unknown redemption transaction shape")
		}
	}

	// Finally, fill the builder with outputs constructed so far.
	for _, output := range outputs {
		builder.AddOutput(output)
	}

	return builder, nil
}
