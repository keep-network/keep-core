package tbtc

import (
	"fmt"
	"math/big"

	"github.com/ipfs/go-log/v2"

	"github.com/keep-network/keep-common/pkg/chain/ethereum"
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"go.uber.org/zap"
)

const (
	// TODO: Determine what the value should be.
	movingFundsProposalValidityBlocks = 600
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
		logger:                       logger,
		chain:                        chain,
		btcChain:                     btcChain,
		movingFundsWallet:            movingFundsWallet,
		transactionExecutor:          transactionExecutor,
		proposal:                     proposal,
		proposalProcessingStartBlock: proposalProcessingStartBlock,
		proposalExpiryBlock:          proposalExpiryBlock,
	}
}

func (mfa *movingFundsAction) execute() error {
	walletPublicKeyHash := bitcoin.PublicKeyHash(mfa.wallet().publicKey)

	// Make sure the moving funds commitment transaction has entered the
	// Ethereum blockchain and has accumulated enough confirmations.
	err := mfa.ensureCommitmentTxConfirmed(walletPublicKeyHash)
	if err != nil {
		return fmt.Errorf(
			"failed to ensure moving funds transaction confirmed: [%v]",
			err,
		)
	}

	// TODO: Continue with the implementation.
	return nil
}

func (mfa *movingFundsAction) ensureCommitmentTxConfirmed(
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

	// Wait `32` blocks to make sure the transaction has not been reverted for
	// some reason, e.g. due to a chain reorganization.
	result, err := ethereum.WaitForBlockConfirmations(
		blockCounter,
		currentBlock,
		32,
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
) ([][20]byte, error) {
	validateProposalLogger.Infof("calling chain for proposal validation")

	err := chain.ValidateMovingFundsProposal(
		walletPublicKeyHash,
		mainUTXO,
		proposal,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"moving funds proposal is invalid: [%v]",
			err,
		)
	}

	return proposal.TargetWallets, nil
}

func (mfa *movingFundsAction) wallet() wallet {
	return mfa.movingFundsWallet
}

func (mfa *movingFundsAction) actionType() WalletActionType {
	return ActionMovingFunds
}
