package tbtc

import (
	"crypto/ecdsa"
	"fmt"
	"time"

	"github.com/ipfs/go-log/v2"
	"go.uber.org/zap"

	"github.com/keep-network/keep-core/pkg/bitcoin"
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

// depositSweepAction is a deposit sweep walletAction.
type depositSweepAction struct {
	logger   *zap.SugaredLogger
	chain    Chain
	btcChain bitcoin.Chain

	sweepingWallet      wallet
	transactionExecutor *walletTransactionExecutor

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
	signingExecutor walletSigningExecutor,
	proposal *DepositSweepProposal,
	proposalProcessingStartBlock uint64,
	proposalExpiresAt time.Time,
) *depositSweepAction {
	transactionExecutor := newWalletTransactionExecutor(
		btcChain,
		sweepingWallet,
		signingExecutor,
	)

	return &depositSweepAction{
		logger:                         logger,
		chain:                          chain,
		btcChain:                       btcChain,
		sweepingWallet:                 sweepingWallet,
		transactionExecutor:            transactionExecutor,
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

	walletPublicKeyHash := bitcoin.PublicKeyHash(dsa.wallet().publicKey)

	walletMainUtxo, err := DetermineWalletMainUtxo(
		walletPublicKeyHash,
		dsa.chain,
		dsa.btcChain,
	)
	if err != nil {
		return fmt.Errorf(
			"error while determining wallet's main UTXO: [%v]",
			err,
		)
	}

	err = EnsureWalletSyncedBetweenChains(
		walletPublicKeyHash,
		walletMainUtxo,
		dsa.chain,
		dsa.btcChain,
	)
	if err != nil {
		return fmt.Errorf(
			"error while ensuring wallet state is synced between "+
				"BTC and host chain: [%v]",
			err,
		)
	}

	unsignedSweepTx, err := assembleDepositSweepTransaction(
		dsa.btcChain,
		dsa.wallet().publicKey,
		walletMainUtxo,
		validatedDeposits,
		dsa.proposal.SweepTxFee.Int64(),
	)
	if err != nil {
		return fmt.Errorf(
			"error while assembling deposit sweep transaction: [%v]",
			err,
		)
	}

	signTxLogger := dsa.logger.With(
		zap.String("step", "signTransaction"),
	)

	sweepTx, err := dsa.transactionExecutor.signTransaction(
		signTxLogger,
		unsignedSweepTx,
		dsa.proposalProcessingStartBlock,
		dsa.proposalExpiresAt.Add(-dsa.signingTimeoutSafetyMargin),
	)
	if err != nil {
		return fmt.Errorf("sign transaction step failed: [%v]", err)
	}

	broadcastTxLogger := dsa.logger.With(
		zap.String("step", "broadcastTransaction"),
		zap.String("sweepTxHash", sweepTx.Hash().Hex(bitcoin.ReversedByteOrder)),
	)

	err = dsa.transactionExecutor.broadcastTransaction(
		broadcastTxLogger,
		sweepTx,
		dsa.broadcastTimeout,
		dsa.broadcastCheckDelay,
	)
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
	chain interface {
		// PastDepositRevealedEvents fetches past deposit reveal events according
		// to the provided filter or unfiltered if the filter is nil. Returned
		// events are sorted by the block number in the ascending order, i.e. the
		// latest event is at the end of the slice.
		PastDepositRevealedEvents(
			filter *DepositRevealedEventFilter,
		) ([]*DepositRevealedEvent, error)

		// ValidateDepositSweepProposal validates the given deposit sweep proposal
		// against the chain. It requires some additional data about the deposits
		// that must be fetched externally. Returns an error if the proposal is
		// not valid or nil otherwise.
		ValidateDepositSweepProposal(
			proposal *DepositSweepProposal,
			depositsExtraInfo []struct {
				*Deposit
				FundingTx *bitcoin.Transaction
			},
		) error
	},
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
		events, err := chain.PastDepositRevealedEvents(&DepositRevealedEventFilter{
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

	err := chain.ValidateDepositSweepProposal(proposal, depositExtraInfo)
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

func (dsa *depositSweepAction) wallet() wallet {
	return dsa.sweepingWallet
}

func (dsa *depositSweepAction) actionType() WalletActionType {
	return ActionDepositSweep
}

// assembleDepositSweepTransaction constructs an unsigned deposit sweep Bitcoin
// transaction.
//
// Regarding input arguments, the walletMainUtxo parameter is optional and
// can be set as nil if the wallet does not have a main UTXO at the moment.
// The deposits slice must contain at least one element. The fee argument
// is not validated in any way so must be chosen with respect to the system
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
