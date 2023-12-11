package tbtc

import (
	"github.com/ipfs/go-log/v2"

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
	// TODO: Remove WalletPublicKeyHash field.
	WalletPublicKeyHash [20]byte
	TargetWallets       [][20]byte
}

func (mfp *MovingFundsProposal) ActionType() WalletActionType {
	return ActionMovingFunds
}

func (mfp *MovingFundsProposal) ValidityBlocks() uint64 {
	return movingFundsProposalValidityBlocks
}

type MovingFundsRequest struct {
	TargetWalletPublicKeyHash [20]byte
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
	// TODO: Before proceeding with creation of the Bitcoin transaction, wait
	//       32 blocks to ensure the commitment transaction has accumulated
	//       enough confirmations in the Ethereum chain and will not be reverted
	//       even if a reorg occurs.
	return nil
}

// ValidateMovingFundsProposal checks the moving funds proposal with on-chain
// validation rules.
func ValidateMovingFundsProposal(
	validateProposalLogger log.StandardLogger,
	proposal *MovingFundsProposal,
) ([]*MovingFundsRequest, error) {
	// TODO: Implement
	return nil, nil
}

func (mfa *movingFundsAction) wallet() wallet {
	return mfa.movingFundsWallet
}

func (mfa *movingFundsAction) actionType() WalletActionType {
	return ActionMovingFunds
}
