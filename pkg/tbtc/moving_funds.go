package tbtc

import (
	"fmt"
	"math/big"

	"github.com/ipfs/go-log/v2"

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

	validateProposalLogger.Infof("moving funds proposal is valid")

	return proposal.TargetWallets, nil
}

func (mfa *movingFundsAction) wallet() wallet {
	return mfa.movingFundsWallet
}

func (mfa *movingFundsAction) actionType() WalletActionType {
	return ActionMovingFunds
}
