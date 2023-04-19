package tbtc

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"go.uber.org/zap"
)

const (
	// depositSweepProposalConfirmationBlocks determines the block length of the
	// confirmation period that is preserved after a deposit sweep proposal
	// submission.
	depositSweepProposalConfirmationBlocks = 20
	// depositSweepRequiredFundingTxConfirmations determines the minimum
	// number of confirmations that are needed for a deposit funding transaction
	// in order to consider it a valid part of the deposit sweep proposal.
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

		// TODO: Fetch deposit extra info.

		depositExtraInfo[i] = struct {
			*Deposit
			FundingTx *bitcoin.Transaction
		}{
			Deposit:   nil,
			FundingTx: fundingTx,
		}
	}

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
