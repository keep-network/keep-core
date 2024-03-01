package tbtc

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ipfs/go-log/v2"
)

// MovedFundsSweepRequestState represents the state of a moved funds request.
type MovedFundsSweepRequestState uint8

const (
	MovedFundsStateUnknown MovedFundsSweepRequestState = iota
	MovedFundsStatePending
	MovedFundsStateProcessed
	MovedFundsStateTimedOut
)

// MovedFundsSweepRequest represents a moved funds sweep request.
type MovedFundsSweepRequest struct {
	WalletPublicKeyHash [20]byte
	Value               uint64
	CreatedAt           time.Time
	State               MovedFundsSweepRequestState
}

const (
	// movedFundsSweepProposalValidityBlocks determines the moving funds
	// proposal validity time expressed in blocks. In other words, this is the
	// worst-case time for a moving funds during which the wallet is busy and
	// cannot take another actions. The value of 600 blocks is roughly 2 hours,
	// assuming 12 seconds per block.
	movedFundsSweepProposalValidityBlocks = 600
)

// MovedFundsSweepProposal represents a moved funds sweep proposal issued by a
// wallet's coordination leader.
type MovedFundsSweepProposal struct {
	MovingFundsTxHash        [32]byte
	MovingFundsTxOutputIndex uint32
	SweepTxFee               *big.Int
}

func (mfsp *MovedFundsSweepProposal) ActionType() WalletActionType {
	return ActionMovedFundsSweep
}

func (mfsp *MovedFundsSweepProposal) ValidityBlocks() uint64 {
	return movedFundsSweepProposalValidityBlocks
}

// ValidateMovedFundsSweepProposal checks the moved funds sweep proposal with
// on-chain validation rules.
func ValidateMovedFundsSweepProposal(
	validateProposalLogger log.StandardLogger,
	walletPublicKeyHash [20]byte,
	proposal *MovedFundsSweepProposal,
	chain interface {
		// ValidateMovedFundsSweepProposal validates the given moved funds sweep
		// proposal against the chain. Returns an error if the proposal is not
		// valid or nil otherwise.
		ValidateMovedFundsSweepProposal(
			walletPublicKeyHash [20]byte,
			proposal *MovedFundsSweepProposal,
		) error
	},
) error {
	validateProposalLogger.Infof("calling chain for proposal validation")

	err := chain.ValidateMovedFundsSweepProposal(
		walletPublicKeyHash,
		proposal,
	)

	if err != nil {
		return fmt.Errorf("moved funds sweep proposal is invalid: [%v]", err)
	}

	validateProposalLogger.Infof("moved funds sweep proposal is valid")

	return nil
}
