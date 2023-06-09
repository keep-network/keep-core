package tbtc

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"time"

	"github.com/ipfs/go-log/v2"
	"go.uber.org/zap"

	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/tecdsa"
)

const (
	// depositSweepProposalConfirmationBlocks determines the block length of the
	// confirmation period on the host chain that is preserved after a deposit
	// sweep proposal submission.
	depositSweepProposalConfirmationBlocks = 20
	// DepositSweepRequiredFundingTxConfirmations determines the minimum
	// number of confirmations that are needed for a deposit funding Bitcoin
	// transaction in order to consider it a valid part of the deposit sweep
	// proposal.
	DepositSweepRequiredFundingTxConfirmations = 6
	// depositSweepSigningTimeoutSafetyMargin determines the duration of the
	// safety margin that must be preserved between the signing timeout
	// and the timeout of the entire deposit sweep action. This safety
	// margin prevents against the case where signing completes late and there
	// is not enough time to broadcast the sweep transaction properly.
	// In such a case, wallet signatures may leak and make the wallet subject
	// of fraud accusations. Usage of the safety margin ensures there is enough
	// time to perform post-signing steps of the deposit sweep action.
	depositSweepSigningTimeoutSafetyMargin = 1 * time.Hour
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

// depositSweepSigningExecutor is an interface meant to decouple the
// specific implementation of the signing executor from the deposit sweep
// action
type depositSweepSigningExecutor interface {
	signBatch(
		ctx context.Context,
		messages []*big.Int,
		startBlock uint64,
	) ([]*tecdsa.Signature, error)
}

// depositSweepAction is a deposit sweep walletAction.
type depositSweepAction struct {
	logger   *zap.SugaredLogger
	chain    Chain
	btcChain bitcoin.Chain

	sweepingWallet  wallet
	signingExecutor depositSweepSigningExecutor

	proposal                     *DepositSweepProposal
	proposalProcessingStartBlock uint64
	proposalExpiresAt            time.Time

	requiredFundingTxConfirmations uint
	signingTimeoutSafetyMargin     time.Duration
	broadcastTimeout               time.Duration
	broadcastCheckDelay            time.Duration
}

func newDepositSweepAction(
	logger *zap.SugaredLogger,
	chain Chain,
	btcChain bitcoin.Chain,
	sweepingWallet wallet,
	signingExecutor depositSweepSigningExecutor,
	proposal *DepositSweepProposal,
	proposalProcessingStartBlock uint64,
	proposalExpiresAt time.Time,
) *depositSweepAction {
	return &depositSweepAction{
		logger:                         logger,
		chain:                          chain,
		btcChain:                       btcChain,
		sweepingWallet:                 sweepingWallet,
		signingExecutor:                signingExecutor,
		proposal:                       proposal,
		proposalProcessingStartBlock:   proposalProcessingStartBlock,
		proposalExpiresAt:              proposalExpiresAt,
		requiredFundingTxConfirmations: DepositSweepRequiredFundingTxConfirmations,
		signingTimeoutSafetyMargin:     depositSweepSigningTimeoutSafetyMargin,
		broadcastTimeout:               depositSweepBroadcastTimeout,
		broadcastCheckDelay:            depositSweepBroadcastCheckDelay,
	}
}

func (dsa *depositSweepAction) execute() error {
	validateProposalLogger := dsa.logger.With(
		zap.String("step", "validateProposal"),
	)

	validatedDeposits, err := ValidateDepositSweepProposal(
		validateProposalLogger,
		dsa.proposal,
		dsa.requiredFundingTxConfirmations,
		dsa.chain,
		dsa.btcChain,
	)
	if err != nil {
		return fmt.Errorf("validate proposal step failed: [%v]", err)
	}

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

	err = dsa.ensureWalletSyncedBetweenChains(walletMainUtxo)
	if err != nil {
		return fmt.Errorf(
			"error while ensuring wallet state is synced between "+
				"BTC and host chain: [%v]",
			err,
		)
	}

	createTxLogger := dsa.logger.With(
		zap.String("step", "createTransaction"),
	)

	sweepTx, err := dsa.createTransaction(
		createTxLogger,
		walletMainUtxo,
		validatedDeposits,
	)
	if err != nil {
		return fmt.Errorf("create transaction step failed: [%v]", err)
	}

	broadcastTxLogger := dsa.logger.With(
		zap.String("step", "broadcastTransaction"),
		zap.String("sweepTxHash", sweepTx.Hash().Hex(bitcoin.ReversedByteOrder)),
	)

	err = dsa.broadcastTransaction(broadcastTxLogger, sweepTx)
	if err != nil {
		return fmt.Errorf("broadcast transaction step failed: [%v]", err)
	}

	return nil
}

// ValidateDepositSweepProposal checks the deposit sweep proposal with on-chain
// validation rules and verifies transactions on the Bitcoin chain.
func ValidateDepositSweepProposal(
	validateProposalLogger log.StandardLogger,
	proposal *DepositSweepProposal,
	requiredFundingTxConfirmations uint,
	tbtcChain Chain,
	btcChain bitcoin.Chain,
) ([]*Deposit, error) {
	depositExtraInfo := make(
		[]struct {
			*Deposit
			FundingTx *bitcoin.Transaction
		},
		len(proposal.DepositsKeys),
	)

	validateProposalLogger.Infof("gathering prerequisites for proposal validation")

	if len(proposal.DepositsKeys) != len(proposal.DepositsRevealBlocks) {
		return nil, fmt.Errorf("proposal's reveal blocks list has a wrong length")
	}

	for i, depositKey := range proposal.DepositsKeys {
		depositDisplayIndex := fmt.Sprintf("%v/%v", i+1, len(proposal.DepositsKeys))

		validateProposalLogger.Infof(
			"deposit [%v] - checking confirmations count for funding tx",
			depositDisplayIndex,
		)

		confirmations, err := btcChain.GetTransactionConfirmations(
			depositKey.FundingTxHash,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"cannot get funding tx confirmations count "+
					"for deposit [%v]: [%v]",
				depositDisplayIndex,
				err,
			)
		}

		if confirmations < requiredFundingTxConfirmations {
			return nil, fmt.Errorf(
				"funding tx of deposit [%v] has only [%v/%v] of "+
					"required confirmations",
				depositDisplayIndex,
				confirmations,
				requiredFundingTxConfirmations,
			)
		}

		validateProposalLogger.Infof(
			"deposit [%v] - fetching deposit's extra data",
			depositDisplayIndex,
		)

		fundingTx, err := btcChain.GetTransaction(depositKey.FundingTxHash)
		if err != nil {
			return nil, fmt.Errorf(
				"cannot get funding tx data for deposit [%v]: [%v]",
				depositDisplayIndex,
				err,
			)
		}

		revealBlock := proposal.DepositsRevealBlocks[i].Uint64()

		// We need to fetch the past DepositRevealed event for the given deposit.
		// It may be tempting to fetch such events for all deposit keys
		// in the proposal using a single call, however, this solution has
		// serious downsides. Popular chain clients have limitations
		// for fetching past chain events regarding the requested block
		// range and/or returned data size. In this context, it is better to
		// do several well-tailored calls than a single general one.
		// We have the revealBlock passed by the coordinator within the proposal
		// so, we can use it to make a narrow call. Moreover, we use the
		// wallet PKH as additional filter to limit the size of returned data.
		events, err := tbtcChain.PastDepositRevealedEvents(&DepositRevealedEventFilter{
			StartBlock:          revealBlock,
			EndBlock:            &revealBlock,
			WalletPublicKeyHash: [][20]byte{proposal.WalletPublicKeyHash},
		})
		if err != nil {
			return nil, fmt.Errorf(
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
			return nil, fmt.Errorf(
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

	validateProposalLogger.Infof("calling chain for proposal validation")

	err := tbtcChain.ValidateDepositSweepProposal(proposal, depositExtraInfo)
	if err != nil {
		return nil, fmt.Errorf("deposit sweep proposal is invalid: [%v]", err)
	}

	validateProposalLogger.Infof(
		"deposit sweep proposal is valid",
	)

	deposits := make([]*Deposit, len(depositExtraInfo))
	for i, dei := range depositExtraInfo {
		deposits[i] = dei.Deposit
	}

	return deposits, nil
}

// ensureWalletSyncedBetweenChains makes sure all actions taken by the wallet
// on the Bitcoin chain are reflected in the host chain Bridge. This translates
// to two conditions that must be met:
// - The wallet main UTXO registered in the host chain Bridge comes from the
//   latest BTC transaction OR wallet main UTXO is unset and wallet's BTC
//   transaction history is empty. This condition ensures that all expected SPV
//   proofs of confirmed BTC transactions were submitted to the host chain Bridge
//   thus the wallet state held known to the Bridge matches the actual state
//   on the BTC chain.
// - There are no pending BTC transactions in the mempool. This condition
//   ensures the wallet doesn't currently perform any action on the BTC chain.
//   Such a transactions indicate a possible state change in the future
//   but their outcome cannot be determined at this stage so, the wallet
//   should not perform new actions at the moment.
func (dsa *depositSweepAction) ensureWalletSyncedBetweenChains(
	walletMainUtxo *bitcoin.UnspentTransactionOutput,
) error {
	walletPublicKeyHash := bitcoin.PublicKeyHash(dsa.wallet().publicKey)

	// Take the recent transactions history for the wallet.
	history, err := dsa.btcChain.GetTransactionsForPublicKeyHash(walletPublicKeyHash, 5)
	if err != nil {
		return fmt.Errorf("cannot get transactions history: [%v]", err)
	}

	if walletMainUtxo != nil {
		// If the wallet main UTXO exists, the transaction history must
		// contain at least one item. If it is empty, something went
		// really wrong. This should never happen but check this scenario
		// just in case.
		if len(history) == 0 {
			return fmt.Errorf(
				"wallet main UTXO exists but there are no BTC " +
					"transactions produced by the wallet",
			)
		}

		// The transaction history is not empty for sure. Take the latest BTC
		// transaction from the history.
		latestTransaction := history[len(history)-1]

		// Make sure the wallet main UTXO comes from the latest transaction.
		// That means all expected SPV proofs were submitted to the Bridge.
		// If the wallet main UTXO transaction hash doesn't match the latest
		// transaction, that means the SPV proof for the latest transaction was
		// not submitted to the Bridge yet.
		//
		// Note that it is enough to check that the wallet main UTXO transaction
		// hash matches the latest transaction hash. There is no way the main
		// UTXO changes and the transaction hash stays the same. The Bridge
		// enforces that all wallet transactions form a sequence and refer
		// each other.
		if walletMainUtxo.Outpoint.TransactionHash != latestTransaction.Hash() {
			return fmt.Errorf(
				"wallet main UTXO doesn't come from the latest BTC transaction",
			)
		}
	} else {
		// If the wallet main UTXO doesn't exist, the transaction history must
		// be empty. If it is not, that could mean there is a Bitcoin transaction
		// produced by the wallet whose SPV proof was not submitted to
		// the Bridge yet.
		if len(history) != 0 {
			return fmt.Errorf(
				"wallet main UTXO doesn't exist but there are BTC " +
					"transactions produced by the wallet",
			)
		}
	}

	// Regardless of the main UTXO state, we need to make sure that
	// no pending wallet transactions exist in the mempool. That way,
	// we are handling a plenty of corner cases like transactions races
	// that could potentially lead to fraudulent transactions and funds loss.
	mempool, err := dsa.btcChain.GetMempoolForPublicKeyHash(walletPublicKeyHash)
	if err != nil {
		return fmt.Errorf("cannot get mempool: [%v]", err)
	}

	if len(mempool) != 0 {
		return fmt.Errorf("unconfirmed transactions exist in the mempool")
	}

	return nil
}

func (dsa *depositSweepAction) createTransaction(
	createTxLogger log.StandardLogger,
	walletMainUtxo *bitcoin.UnspentTransactionOutput,
	deposits []*Deposit,
) (*bitcoin.Transaction, error) {
	createTxLogger.Infof("creating deposit sweep transaction")

	unsignedSweepTx, err := assembleDepositSweepTransaction(
		dsa.btcChain,
		dsa.wallet().publicKey,
		walletMainUtxo,
		deposits,
		dsa.proposal.SweepTxFee.Int64(),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error while assembling deposit sweep transaction: [%v]",
			err,
		)
	}

	createTxLogger.Infof("computing deposit sweep transaction's sig hashes")

	sigHashes, err := unsignedSweepTx.ComputeSignatureHashes()
	if err != nil {
		return nil, fmt.Errorf(
			"error while computing deposit sweep transaction's "+
				"sig hashes: [%v]",
			err,
		)
	}

	createTxLogger.Infof("signing deposit sweep transaction's sig hashes")

	// Make sure signing times out far before the entire action.
	signingTimesOutAt := dsa.proposalExpiresAt.Add(-dsa.signingTimeoutSafetyMargin)
	signingCtx, cancelSigningCtx := context.WithTimeout(
		context.Background(),
		time.Until(signingTimesOutAt),
	)
	defer cancelSigningCtx()

	signatures, err := dsa.signingExecutor.signBatch(
		signingCtx,
		sigHashes,
		dsa.proposalProcessingStartBlock,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error while signing deposit sweep transaction's "+
				"sig hashes: [%v]",
			err,
		)
	}

	createTxLogger.Infof("applying deposit sweep transaction's signatures")

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
		return nil, fmt.Errorf(
			"error while applying deposit sweep transaction's "+
				"signatures: [%v]",
			err,
		)
	}

	createTxLogger.Infof("deposit sweep transaction created successfully")

	return sweepTx, nil
}

func (dsa *depositSweepAction) broadcastTransaction(
	broadcastTxLogger log.StandardLogger,
	sweepTx *bitcoin.Transaction,
) error {
	sweepTxHash := sweepTx.Hash()

	broadcastCtx, cancelBroadcastCtx := context.WithTimeout(
		context.Background(),
		dsa.broadcastTimeout,
	)
	defer cancelBroadcastCtx()

	broadcastAttempt := 0

	for {
		select {
		case <-broadcastCtx.Done():
			return fmt.Errorf("broadcast timeout exceeded")
		default:
			broadcastAttempt++

			broadcastTxLogger.Infof(
				"broadcasting deposit sweep transaction on "+
					"the Bitcoin chain - attempt [%v]",
				broadcastAttempt,
			)

			err := dsa.btcChain.BroadcastTransaction(sweepTx)
			if err != nil {
				broadcastTxLogger.Warnf(
					"broadcasting failed: [%v]; transaction could be "+
						"broadcasted by another wallet operators though",
					err,
				)
			} else {
				broadcastTxLogger.Infof("broadcasting completed")
			}

			broadcastTxLogger.Infof(
				"waiting [%v] before checking whether the "+
					"transaction is known on Bitcoin chain",
				dsa.broadcastCheckDelay,
			)

			select {
			case <-time.After(dsa.broadcastCheckDelay):
			case <-broadcastCtx.Done():
				return fmt.Errorf("broadcast timeout exceeded")
			}

			broadcastTxLogger.Infof(
				"checking whether the transaction is known on Bitcoin chain",
			)

			_, err = dsa.btcChain.GetTransactionConfirmations(sweepTxHash)
			if err != nil {
				broadcastTxLogger.Warnf(
					"cannot say whether the transaction is known "+
						"on Bitcoin chain; check returned an error: [%v]",
					err,
				)
				continue
			}

			broadcastTxLogger.Infof("transaction is known on Bitcoin chain")
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
