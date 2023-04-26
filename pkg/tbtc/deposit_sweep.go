package tbtc

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ipfs/go-log/v2"
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/chain"
	"go.uber.org/zap"
	"time"
)

const (
	// depositSweepProposalConfirmationBlocks determines the block length of the
	// confirmation period on the host chain that is preserved after a deposit
	// sweep proposal submission.
	depositSweepProposalConfirmationBlocks = 20
	// depositSweepRequiredFundingTxConfirmations determines the minimum
	// number of confirmations that are needed for a deposit funding Bitcoin
	// transaction in order to consider it a valid part of the deposit sweep
	// proposal.
	depositSweepRequiredFundingTxConfirmations = 6
	// depositSweepSigningTimeoutSafetyMargin determines the duration of the
	// safety margin that must be preserved between the signing timeout
	// and the timeout of the entire deposit sweep action. This safety
	// margin prevents against the case where signing completes late and there
	// is not enough time to broadcast the sweep transaction properly.
	// In such a case, wallet signatures may leak and make the wallet subject
	// of fraud accusations. Usage of the safety margin ensures there is enough
	// time to perform post-signing steps of the deposit sweep action.
	depositSweepSigningTimeoutSafetyMargin = 1 * time.Hour
	// depositSweepSigningDelayBlocks determines the per-deposit delay in
	// blocks that must be preserved before starting the deposit sweep
	// transaction signing process. This delay aims to reflect the time taken
	// by pre-signing steps that must be done for a single deposit.
	// Multiplying this constant by the number of proposal's deposits
	// allows to determine the total delay that must be added to the proposal
	// processing start block in order to designate a sane signing start block
	// and maximize chances for a successful signing process.
	depositSweepSigningDelayBlocks = 10
	// depositSweepBroadcastTimeout determines the time window for deposit
	// sweep transaction broadcast. It is guaranteed that at least
	// depositSweepSigningTimeoutSafetyMargin is preserved for the broadcast
	// step. However, the happy path for the broadcast step is usually quick
	// and few retries are needed to recover from temporary problems. That
	// said, if the broadcast step does not succeed in a tight timeframe,
	// there is no point to retry for the entire possible time window.
	// Hence, the timeout for broadcast step is set as 25% of the entire
	// time widow determined by depositSweepSigningTimeoutSafetyMargin.
	depositSweepBroadcastTimeout = depositSweepSigningTimeoutSafetyMargin / 4
	// depositSweepBroadcastCheckDelay determines the delay that must
	// be preserved between transaction broadcast and the check that ensures
	// the transaction is known on the Bitcoin chain. This delay is needed
	// as spreading the transaction over the Bitcoin network takes time.
	depositSweepBroadcastCheckDelay = 1 * time.Minute
)

// depositSweepAction is a deposit sweep walletAction.
type depositSweepAction struct {
	chain    Chain
	btcChain bitcoin.Chain

	sweepingWallet  wallet
	signingExecutor *signingExecutor

	proposal                     *DepositSweepProposal
	proposalProcessingStartBlock uint64
	proposalExpiresAt            time.Time
}

func newDepositSweepAction(
	chain Chain,
	btcChain bitcoin.Chain,
	sweepingWallet wallet,
	signingExecutor *signingExecutor,
	proposal *DepositSweepProposal,
	proposalProcessingStartBlock uint64,
	proposalExpiresAt time.Time,
) *depositSweepAction {
	return &depositSweepAction{
		chain:                        chain,
		btcChain:                     btcChain,
		sweepingWallet:               sweepingWallet,
		signingExecutor:              signingExecutor,
		proposal:                     proposal,
		proposalProcessingStartBlock: proposalProcessingStartBlock,
		proposalExpiresAt:            proposalExpiresAt,
	}
}

// TODO: Cover this function with unit tests once everything is completed.
func (dsa *depositSweepAction) execute() error {
	walletPublicKeyBytes, err := marshalPublicKey(dsa.wallet().publicKey)
	if err != nil {
		return fmt.Errorf("cannot marshal wallet public key: [%v]", err)
	}

	actionLogger := logger.With(
		zap.String("wallet", fmt.Sprintf("0x%x", walletPublicKeyBytes)),
		zap.String("action", dsa.actionType().String()),
		zap.Uint64("startBlock", dsa.proposalProcessingStartBlock),
	)

	depositExtraInfo := make(
		[]struct {
			*Deposit
			FundingTx *bitcoin.Transaction
		},
		len(dsa.proposal.DepositsKeys),
	)

	depositsCount := len(dsa.proposal.DepositsKeys)

	actionLogger.Infof("gathering prerequisites for proposal validation")

	for i, depositKey := range dsa.proposal.DepositsKeys {
		depositDisplayIndex := fmt.Sprintf("%v/%v", i+1, depositsCount)

		depositLogger := actionLogger.With(
			zap.String(
				"depositFundingTxHash",
				fmt.Sprintf("0x%x", depositKey.FundingTxHash),
			),
			zap.String("depositIndex", depositDisplayIndex),
		)

		depositLogger.Infof("checking confirmations count for deposit's funding tx")

		confirmations, err := dsa.btcChain.GetTransactionConfirmations(
			depositKey.FundingTxHash,
		)
		if err != nil {
			return fmt.Errorf(
				"cannot get funding tx confirmations count "+
					"for deposit [%v]: [%v]",
				depositDisplayIndex,
				err,
			)
		}

		if confirmations < depositSweepRequiredFundingTxConfirmations {
			return fmt.Errorf(
				"funding tx of deposit [%v] has only [%v/%v] of "+
					"required confirmations",
				depositDisplayIndex,
				confirmations,
				depositSweepRequiredFundingTxConfirmations,
			)
		}

		depositLogger.Infof("fetching deposit's extra data")

		fundingTx, err := dsa.btcChain.GetTransaction(depositKey.FundingTxHash)
		if err != nil {
			return fmt.Errorf(
				"cannot get funding tx data for deposit [%v]: [%v]",
				depositDisplayIndex,
				err,
			)
		}

		depositRequest, err := dsa.chain.GetDepositRequest(
			depositKey.FundingTxHash,
			depositKey.FundingOutputIndex,
		)
		if err != nil {
			return fmt.Errorf(
				"cannot get on-chain request data for deposit [%v]: [%v]",
				depositDisplayIndex,
				err,
			)
		}

		revealBlock, err := dsa.chain.GetBlockNumberByTimestamp(
			uint64(depositRequest.RevealedAt.Unix()),
		)
		if err != nil {
			return fmt.Errorf(
				"cannot estimate reveal block for deposit [%v]: [%v]",
				depositDisplayIndex,
				err,
			)
		}

		// We need to fetch the past DepositRevealed event for the given deposit.
		// It may be tempting to fetch such events for all deposit keys
		// in the proposal using a single call, however, this solution has
		// serious downsides. Popular chain clients have limitations
		// for fetching past chain events regarding the requested block
		// range and/or returned data size. In this context, it is better to
		// do several well-tailored calls than a single general one.
		// We estimated the revealBlock so, we know the event was emitted
		// at this block or somewhere close. It makes sense to establish
		// a small margin and fetch past DepositRevealed events from range
		// [revealBlock - margin, revealBlock + margin] in order to handle
		// possible inaccuracies of revealBlock estimation. Moreover,
		// we use the depositor address and wallet PKH as additional filters
		// to limit the size of returned data.
		revealBlockMargin := uint64(10)
		startBlock := revealBlock - revealBlockMargin
		endBlock := revealBlock + revealBlockMargin

		events, err := dsa.chain.PastDepositRevealedEvents(&DepositRevealedEventFilter{
			StartBlock:          startBlock,
			EndBlock:            &endBlock,
			Depositor:           []chain.Address{depositRequest.Depositor},
			WalletPublicKeyHash: [][20]byte{dsa.proposal.WalletPublicKeyHash},
		})
		if err != nil {
			return fmt.Errorf(
				"cannot get on-chain DepositRevealed events for deposit [%v]: [%v]",
				depositDisplayIndex,
				err,
			)
		}

		// There may be multiple events returned for the provided filter.
		// Find the one matching our depositKey.
		var matchingEvent *DepositRevealedEvent
		for _, event := range events {
			if event.FundingTxHash == depositKey.FundingTxHash &&
				event.FundingOutputIndex == depositKey.FundingOutputIndex {
				matchingEvent = event
				break
			}
		}

		if matchingEvent == nil {
			return fmt.Errorf(
				"no matching DepositRevealed event for deposit [%v]: [%v]",
				depositDisplayIndex,
				err,
			)
		}

		depositExtraInfo[i] = struct {
			*Deposit
			FundingTx *bitcoin.Transaction
		}{
			Deposit:   matchingEvent.unpack(),
			FundingTx: fundingTx,
		}
	}

	actionLogger.Infof("calling chain for proposal validation")

	err = dsa.chain.ValidateDepositSweepProposal(dsa.proposal, depositExtraInfo)
	if err != nil {
		return fmt.Errorf("deposit sweep proposal is invalid: [%v]", err)
	}

	actionLogger.Infof(
		"deposit sweep proposal is valid; assembling deposit sweep transaction",
	)

	walletMainUtxo, err := DetermineWalletMainUtxo(
		bitcoin.PublicKeyHash(dsa.wallet().publicKey),
		dsa.chain,
		dsa.btcChain,
	)
	if err != nil {
		return fmt.Errorf(
			"error while determining wallet's main UTXO: [%v]",
			err,
		)
	}

	deposits := make([]*Deposit, len(depositExtraInfo))
	for i, dei := range depositExtraInfo {
		deposits[i] = dei.Deposit
	}

	unsignedSweepTx, err := assembleDepositSweepTransaction(
		dsa.btcChain,
		dsa.wallet().publicKey,
		walletMainUtxo,
		deposits,
		dsa.proposal.SweepTxFee.Int64(),
	)
	if err != nil {
		return fmt.Errorf(
			"error while bulding deposit sweep transaction: [%v]",
			err,
		)
	}

	actionLogger.Infof("computing deposit sweep transaction's sig hashes")

	sigHashes, err := unsignedSweepTx.ComputeSignatureHashes()
	if err != nil {
		return fmt.Errorf(
			"error while computing deposit sweep transaction's "+
				"sig hashes: [%v]",
			err,
		)
	}

	actionLogger.Infof("signing deposit sweep transaction's sig hashes")

	// Make sure signing times out far before the entire action.
	signingTimesOutAt := dsa.proposalExpiresAt.Add(-depositSweepSigningTimeoutSafetyMargin)
	signingCtx, cancelSigningCtx := context.WithTimeout(
		context.Background(),
		time.Until(signingTimesOutAt),
	)
	defer cancelSigningCtx()

	// Make sure the signing start block takes into account the time elapsed
	// during pre-signing steps, i.e. gathering deposit data and proposal
	// validation.
	signingStartBlock := dsa.proposalProcessingStartBlock +
		uint64(depositsCount*depositSweepSigningDelayBlocks)

	signatures, err := dsa.signingExecutor.signBatch(
		signingCtx,
		sigHashes,
		signingStartBlock,
	)
	if err != nil {
		return fmt.Errorf(
			"error while signing deposit sweep transaction's "+
				"sig hashes: [%v]",
			err,
		)
	}

	actionLogger.Infof("applying deposit sweep transaction's signatures")

	containers := make([]*bitcoin.SignatureContainer, len(signatures))
	for i, signature := range signatures {
		containers[i] = &bitcoin.SignatureContainer{
			R:         signature.R,
			S:         signature.S,
			PublicKey: dsa.wallet().publicKey,
		}
	}

	sweepTx, err := unsignedSweepTx.AddSignatures(containers)
	if err != nil {
		return fmt.Errorf(
			"error while applying deposit sweep transaction's "+
				"signatures: [%v]",
			err,
		)
	}

	broadcastLogger := actionLogger.With(
		zap.String("sweepTxHash", fmt.Sprintf("0x%x", sweepTx.Hash())),
	)

	return dsa.broadcastTransaction(
		broadcastLogger,
		sweepTx,
		depositSweepBroadcastTimeout,
		depositSweepBroadcastCheckDelay,
	)
}

func (dsa *depositSweepAction) broadcastTransaction(
	broadcastLogger log.StandardLogger,
	sweepTx *bitcoin.Transaction,
	timeout time.Duration,
	checkDelay time.Duration,
) error {
	sweepTxHash := sweepTx.Hash()

	broadcastCtx, cancelBroadcastCtx := context.WithTimeout(
		context.Background(),
		timeout,
	)
	defer cancelBroadcastCtx()

	broadcastAttempt := 0

	for {
		select {
		case <-broadcastCtx.Done():
			return fmt.Errorf("broadcast timeout exceeded")
		default:
			broadcastAttempt++

			broadcastLogger.Infof(
				"broadcasting deposit sweep transaction on "+
					"the Bitcoin chain - attempt [%v]",
				broadcastAttempt,
			)

			err := dsa.btcChain.BroadcastTransaction(sweepTx)
			if err != nil {
				broadcastLogger.Warnf(
					"broadcasting failed: [%v]; transaction could be "+
						"broadcasted by another wallet operators though",
					err,
				)
			} else {
				broadcastLogger.Infof("broadcasting completed")
			}

			broadcastLogger.Infof(
				"waiting [%v] before checking whether the "+
					"transaction is known on Bitcoin chain",
				checkDelay,
			)

			select {
			case <-time.After(checkDelay):
			case <-broadcastCtx.Done():
				return fmt.Errorf("broadcast timeout exceeded")
			}

			broadcastLogger.Infof(
				"checking whether the transaction is known on Bitcoin chain",
			)

			_, err = dsa.btcChain.GetTransactionConfirmations(sweepTxHash)
			if err != nil {
				broadcastLogger.Warnf(
					"cannot say whether the transaction is known "+
						"on Bitcoin chain; check returned an error: [%v]",
					err,
				)
				continue
			}

			broadcastLogger.Infof("transaction is known on Bitcoin chain")
			return nil
		}
	}
}

func (dsa *depositSweepAction) wallet() wallet {
	return dsa.sweepingWallet
}

func (dsa *depositSweepAction) actionType() WalletActionType {
	return DepositSweep
}

// assembleDepositSweepTransaction constructs an unsigned deposit sweep Bitcoin
// transaction.
//
// Regarding input arguments, the walletPublicKey parameter is optional and
// can be set as nil if the wallet does not have a main UTXO at the moment.
// The deposits slice must contain at least one element. The fee argument
// is not validated anyway so must be chosen with respect to the system
// limitations.
//
// The resulting bitcoin.TransactionBuilder instance holds all the data
// necessary to sign the transaction and obtain a bitcoin.Transaction instance
// ready to be spread across the Bitcoin network.
func assembleDepositSweepTransaction(
	bitcoinChain bitcoin.Chain,
	walletPublicKey *ecdsa.PublicKey,
	walletMainUtxo *bitcoin.UnspentTransactionOutput,
	deposits []*Deposit,
	fee int64,
) (*bitcoin.TransactionBuilder, error) {
	if len(deposits) < 1 {
		return nil, fmt.Errorf("at least one deposit is required")
	}

	builder := bitcoin.NewTransactionBuilder(bitcoinChain)

	if walletMainUtxo != nil {
		err := builder.AddPublicKeyHashInput(walletMainUtxo)
		if err != nil {
			return nil, fmt.Errorf(
				"cannot add input pointing to wallet main UTXO: [%v]",
				err,
			)
		}
	}

	for i, deposit := range deposits {
		depositScript, err := deposit.Script()
		if err != nil {
			return nil, fmt.Errorf(
				"cannot get script for deposit [%v]: [%v]",
				i,
				err,
			)
		}

		err = builder.AddScriptHashInput(deposit.Utxo, depositScript)
		if err != nil {
			return nil, fmt.Errorf(
				"cannot add input pointing to deposit [%v] UTXO: [%v]",
				i,
				err,
			)
		}
	}

	walletPublicKeyHash := bitcoin.PublicKeyHash(walletPublicKey)
	outputScript, err := bitcoin.PayToWitnessPublicKeyHash(walletPublicKeyHash)
	if err != nil {
		return nil, fmt.Errorf("cannot compute output script: [%v]", err)
	}

	outputValue := builder.TotalInputsValue() - fee

	builder.AddOutput(&bitcoin.TransactionOutput{
		Value:           outputValue,
		PublicKeyScript: outputScript,
	})

	return builder, nil
}
