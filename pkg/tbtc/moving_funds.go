package tbtc

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"time"

	"github.com/ipfs/go-log/v2"

	"github.com/keep-network/keep-common/pkg/chain/ethereum"
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"go.uber.org/zap"
)

const (
	// movingFundsProposalValidityBlocks determines the moving funds proposal
	// validity time expressed in blocks. In other words, this is the worst-case
	// time for moving funds during which the wallet is busy and cannot take
	// another actions. The value of 650 blocks is roughly 2 hours and 10
	// minutes, assuming 12 seconds per block. It is a slightly longer validity
	// time than in case of redemptions as moving funds involves waiting for
	// target wallets commitment transaction to be confirmed.
	movingFundsProposalValidityBlocks = 650
	// movingFundsSigningTimeoutSafetyMarginBlocks determines the duration of
	// the safety margin that must be preserved between the signing timeout and
	// the timeout of the entire moving funds action. This safety margin
	// prevents against the case where signing completes late and there is not
	// enough time to broadcast the moving funds transaction properly.
	// In such a case, wallet signatures may leak and make the wallet subject
	// of fraud accusations. Usage of the safety margin ensures there is enough
	// time to perform post-signing steps of the moving funds action.
	// The value of 300 blocks is roughly 1 hour, assuming 12 seconds per block.
	movingFundsSigningTimeoutSafetyMarginBlocks = 300
	// movingFundsBroadcastTimeout determines the time window for moving funds
	// transaction broadcast. It is guaranteed that at least
	// movingFundsSigningTimeoutSafetyMarginBlocks is preserved for the broadcast
	// step. However, the happy path for the broadcast step is usually quick
	// and few retries are needed to recover from temporary problems. That
	// said, if the broadcast step does not succeed in a tight timeframe,
	// there is no point to retry for the entire possible time window.
	// Hence, the timeout for broadcast step is set as 25% of the entire
	// time widow determined by movingFundsSigningTimeoutSafetyMarginBlocks.
	movingFundsBroadcastTimeout = 15 * time.Minute
	// movingFundsBroadcastCheckDelay determines the delay that must
	// be preserved between transaction broadcast and the check that ensures
	// the transaction is known on the Bitcoin chain. This delay is needed
	// as spreading the transaction over the Bitcoin network takes time.
	movingFundsBroadcastCheckDelay = 1 * time.Minute
	// movingFundsCommitmentConfirmationBlocks determines the period used when
	// waiting for the moving funds commitment confirmations. This period
	// ensures the moving funds commitment has definitely entered the blockchain
	// and will not be removed by a chain reorganization.
	movingFundsCommitmentConfirmationBlocks = 32
)

// MovingFundsProposal represents a moving funds proposal issued by a wallet's
// coordination leader.
type MovingFundsProposal struct {
	TargetWallets    [][20]byte
	MovingFundsTxFee *big.Int
}

func (mfp *MovingFundsProposal) ActionType() WalletActionType {
	return ActionMovingFunds
}

func (mfp *MovingFundsProposal) ValidityBlocks() uint64 {
	return movingFundsProposalValidityBlocks
}

// movingFundsAction is a walletAction implementation handling moving funds
// requests from the wallet coordinator.
type movingFundsAction struct {
	logger   *zap.SugaredLogger
	chain    Chain
	btcChain bitcoin.Chain

	movingFundsWallet   wallet
	transactionExecutor *walletTransactionExecutor

	proposal                     *MovingFundsProposal
	proposalProcessingStartBlock uint64
	proposalExpiryBlock          uint64

	signingTimeoutSafetyMarginBlocks uint64
	broadcastTimeout                 time.Duration
	broadcastCheckDelay              time.Duration
}

func newMovingFundsAction(
	logger *zap.SugaredLogger,
	chain Chain,
	btcChain bitcoin.Chain,
	movingFundsWallet wallet,
	signingExecutor walletSigningExecutor,
	proposal *MovingFundsProposal,
	proposalProcessingStartBlock uint64,
	proposalExpiryBlock uint64,
	waitForBlockFn waitForBlockFn,
) *movingFundsAction {
	transactionExecutor := newWalletTransactionExecutor(
		btcChain,
		movingFundsWallet,
		signingExecutor,
		waitForBlockFn,
	)

	return &movingFundsAction{
		logger:                           logger,
		chain:                            chain,
		btcChain:                         btcChain,
		movingFundsWallet:                movingFundsWallet,
		transactionExecutor:              transactionExecutor,
		proposal:                         proposal,
		proposalProcessingStartBlock:     proposalProcessingStartBlock,
		proposalExpiryBlock:              proposalExpiryBlock,
		signingTimeoutSafetyMarginBlocks: movingFundsSigningTimeoutSafetyMarginBlocks,
		broadcastTimeout:                 movingFundsBroadcastTimeout,
		broadcastCheckDelay:              movingFundsBroadcastCheckDelay,
	}
}

func (mfa *movingFundsAction) execute() error {
	validateProposalLogger := mfa.logger.With(
		zap.String("step", "validateProposal"),
	)

	walletPublicKeyHash := bitcoin.PublicKeyHash(mfa.wallet().publicKey)

	walletMainUtxo, err := DetermineWalletMainUtxo(
		walletPublicKeyHash,
		mfa.chain,
		mfa.btcChain,
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
		return fmt.Errorf("moving funds wallet has no main UTXO")
	}

	// Perform initial validation of the moving funds proposal.
	err = ValidateMovingFundsProposal(
		validateProposalLogger,
		walletPublicKeyHash,
		walletMainUtxo,
		mfa.proposal,
		mfa.chain,
	)
	if err != nil {
		return fmt.Errorf("validate proposal step failed: [%v]", err)
	}

	// Wait until the moving funds commitment transaction has gathered a
	// significant number of confirmations. This ensures that even a deep reorg
	// will not remove the commitment from the chain.
	err = mfa.waitForCommitmentConfirmation(walletPublicKeyHash)
	if err != nil {
		return fmt.Errorf(
			"failed to ensure moving funds transaction confirmed: [%v]",
			err,
		)
	}

	// Validate the moving funds proposal again after waiting for the
	// transaction confirmations. This repeated validation ensures the moving
	// funds commitment was not changed during waiting for the confirmations.
	err = ValidateMovingFundsProposal(
		validateProposalLogger,
		walletPublicKeyHash,
		walletMainUtxo,
		mfa.proposal,
		mfa.chain,
	)
	if err != nil {
		return fmt.Errorf("validate proposal step failed: [%v]", err)
	}

	err = EnsureWalletSyncedBetweenChains(
		walletPublicKeyHash,
		walletMainUtxo,
		mfa.chain,
		mfa.btcChain,
	)
	if err != nil {
		return fmt.Errorf(
			"error while ensuring wallet state is synced between "+
				"BTC and host chain: [%v]",
			err,
		)
	}

	unsignedMovingFundsTx, err := assembleMovingFundsTransaction(
		mfa.btcChain,
		mfa.wallet().publicKey,
		walletMainUtxo,
		mfa.proposal.TargetWallets,
		mfa.proposal.MovingFundsTxFee.Int64(),
	)
	if err != nil {
		return fmt.Errorf(
			"error while assembling moving funds transaction: [%v]",
			err,
		)
	}

	signTxLogger := mfa.logger.With(
		zap.String("step", "signTransaction"),
	)

	// Just in case. This should never happen.
	if mfa.proposalExpiryBlock < mfa.signingTimeoutSafetyMarginBlocks {
		return fmt.Errorf("invalid proposal expiry block")
	}

	movingFundsTx, err := mfa.transactionExecutor.signTransaction(
		signTxLogger,
		unsignedMovingFundsTx,
		mfa.proposalProcessingStartBlock+movingFundsCommitmentConfirmationBlocks,
		mfa.proposalExpiryBlock-mfa.signingTimeoutSafetyMarginBlocks,
	)
	if err != nil {
		return fmt.Errorf("sign transaction step failed: [%v]", err)
	}

	broadcastTxLogger := mfa.logger.With(
		zap.String("step", "broadcastTransaction"),
		zap.String(
			"movingFundsTxHash",
			movingFundsTx.Hash().Hex(bitcoin.ReversedByteOrder),
		),
	)

	err = mfa.transactionExecutor.broadcastTransaction(
		broadcastTxLogger,
		movingFundsTx,
		mfa.broadcastTimeout,
		mfa.broadcastCheckDelay,
	)
	if err != nil {
		return fmt.Errorf("broadcast transaction step failed: [%v]", err)
	}

	return nil
}

func (mfa *movingFundsAction) waitForCommitmentConfirmation(
	walletPublicKeyHash [20]byte,
) error {
	blockCounter, err := mfa.chain.BlockCounter()
	if err != nil {
		return fmt.Errorf("error getting block counter [%w]", err)
	}

	currentBlock, err := blockCounter.CurrentBlock()
	if err != nil {
		return fmt.Errorf("error getting current block [%w]", err)
	}

	// To verify the commitment transaction is in the Ethereum blockchain check
	// that the commitment hash is not zero.
	stateCheck := func() (bool, error) {
		walletData, err := mfa.chain.GetWallet(walletPublicKeyHash)
		if err != nil {
			return false, err
		}

		return walletData.MovingFundsTargetWalletsCommitmentHash != [32]byte{}, nil
	}

	// Wait a significant number of blocks to make sure the transaction has not
	// been reverted for some reason, e.g. due to a chain reorganization.
	result, err := ethereum.WaitForBlockConfirmations(
		blockCounter,
		currentBlock,
		movingFundsCommitmentConfirmationBlocks,
		stateCheck,
	)
	if err != nil {
		return fmt.Errorf(
			"error while waiting for transaction confirmation [%w]",
			err,
		)
	}

	if !result {
		return fmt.Errorf("transaction not included in blockchain")
	}

	return nil
}

// ValidateMovingFundsProposal checks the moving funds proposal with on-chain
// validation rules.
func ValidateMovingFundsProposal(
	validateProposalLogger log.StandardLogger,
	walletPublicKeyHash [20]byte,
	mainUTXO *bitcoin.UnspentTransactionOutput,
	proposal *MovingFundsProposal,
	chain interface {
		// ValidateMovingFundsProposal validates the given moving funds proposal
		// against the chain. Returns an error if the proposal is not valid or
		// nil otherwise.
		ValidateMovingFundsProposal(
			walletPublicKeyHash [20]byte,
			mainUTXO *bitcoin.UnspentTransactionOutput,
			proposal *MovingFundsProposal,
		) error
	},
) error {
	validateProposalLogger.Infof("calling chain for proposal validation")

	err := chain.ValidateMovingFundsProposal(
		walletPublicKeyHash,
		mainUTXO,
		proposal,
	)
	if err != nil {
		return fmt.Errorf("moving funds proposal is invalid: [%v]", err)
	}

	validateProposalLogger.Infof("moving funds proposal is valid")

	return nil
}

func (mfa *movingFundsAction) wallet() wallet {
	return mfa.movingFundsWallet
}

func (mfa *movingFundsAction) actionType() WalletActionType {
	return ActionMovingFunds
}

func assembleMovingFundsTransaction(
	bitcoinChain bitcoin.Chain,
	walletPublicKey *ecdsa.PublicKey,
	walletMainUtxo *bitcoin.UnspentTransactionOutput,
	targetWallets [][20]byte,
	fee int64,
) (*bitcoin.TransactionBuilder, error) {
	if len(targetWallets) < 1 {
		return nil, fmt.Errorf("at least one target wallet is required")
	}

	if walletMainUtxo == nil {
		return nil, fmt.Errorf("wallet main UTXO is required")
	}

	builder := bitcoin.NewTransactionBuilder(bitcoinChain)
	err := builder.AddPublicKeyHashInput(walletMainUtxo)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot add input pointing to wallet main UTXO: [%v]",
			err,
		)
	}

	// The number of target wallets. It determines the number of outputs.
	targetWalletsCount := int64(len(targetWallets))

	// The sum of all the outputs equal to the input value minus fee.
	totalOutputValue := walletMainUtxo.Value - fee

	// The remainder left from distributing the total output value evenly
	// among the outputs. It should be added to the value of the last output.
	// The value of this remainder is at most `len(targetWallets) - 1` satoshis.
	remainder := totalOutputValue % targetWalletsCount

	// The amount of Bitcoin allocated to each output. Additionally, the amount
	// for the last output should be increased by the remainder left from
	// distributing the total output amount across all the outputs.
	singleOutputValue := (totalOutputValue - remainder) / targetWalletsCount

	// Add one output for each target wallets. The order of target wallets
	// should be the same as during the commitment. We don't to check it,
	// as the order of target wallets has already been validated by the on-chain
	// contract.
	for i, targetWalletPublicKeyHash := range targetWallets {
		outputScript, err := bitcoin.PayToWitnessPublicKeyHash(
			targetWalletPublicKeyHash,
		)
		if err != nil {
			return nil, fmt.Errorf("cannot compute output script: [%v]", err)
		}

		outputValue := singleOutputValue
		if i == len(targetWallets)-1 {
			// If we are at the last output, increase its value by the remaining
			// satoshis. If the total output amount is divisible by the number
			// of target wallets, the increase will be `0`.
			outputValue += remainder
		}

		builder.AddOutput(&bitcoin.TransactionOutput{
			Value:           outputValue,
			PublicKeyScript: outputScript,
		})
	}

	return builder, nil
}
