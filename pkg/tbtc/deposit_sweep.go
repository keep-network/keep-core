package tbtc

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/chain"
	"go.uber.org/zap"
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
)

// depositSweepAction is a deposit sweep walletAction.
type depositSweepAction struct {
	chain           Chain
	btcChain        bitcoin.Chain
	sweepingWallet  wallet
	signingExecutor *signingExecutor
	proposal        *DepositSweepProposal
}

func newDepositSweepAction(
	chain Chain,
	btcChain bitcoin.Chain,
	sweepingWallet wallet,
	signingExecutor *signingExecutor,
	proposal *DepositSweepProposal,
) *depositSweepAction {
	return &depositSweepAction{
		chain:           chain,
		btcChain:        btcChain,
		sweepingWallet:  sweepingWallet,
		signingExecutor: signingExecutor,
		proposal:        proposal,
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

		// Worth mentioning this should be treated as an estimation, giving
		// the caveats mentioned in the docstring of GetBlockNumberByTimestamp.
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

	actionLogger.Infof("deposit sweep proposal is valid")

	// TODO: Do the following:
	//       - Construct the transaction
	//       - Sign
	//       - Broadcast
	//       - Monitor

	return nil
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
