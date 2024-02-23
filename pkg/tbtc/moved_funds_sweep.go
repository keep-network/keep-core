package tbtc

import (
	"math/big"
	"time"
)

// MovedFundsSweepRequestState represents the state of a moved funds request.
type MovedFundsSweepRequestState uint8

const (
	MovedFundsStateUnknown MovedFundsSweepRequestState = iota
	MovedFundsStatePending
	MovedFundsStateProcessed
	MovedFundsStateTimedOut
)

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
