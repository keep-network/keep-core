package coordinator

import (
	"math/big"

	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/tbtc"
)

// NewWalletRegisteredEvent represents a new wallet registered event.
type NewWalletRegisteredEvent struct {
	EcdsaWalletID       [32]byte
	WalletPublicKeyHash [20]byte
	BlockNumber         uint64
}

// NewWalletRegisteredEventFilter is a component allowing to filter NewWalletRegisteredEvent.
type NewWalletRegisteredEventFilter struct {
	StartBlock          uint64
	EndBlock            *uint64
	EcdsaWalletID       [][32]byte
	WalletPublicKeyHash [][20]byte
}

// Chain represents the interface that the coordinator module expects to interact
// with the anchoring blockchain on.
type Chain interface {
	// GetDepositRequest gets the on-chain deposit request for the given
	// funding transaction hash and output index. Returns an error if the
	// deposit was not found.
	GetDepositRequest(
		fundingTxHash bitcoin.Hash,
		fundingOutputIndex uint32,
	) (*tbtc.DepositChainRequest, error)

	// PastNewWalletRegisteredEvents fetches past new wallet registered events
	// according to the provided filter or unfiltered if the filter is nil. Returned
	// events are sorted by the block number in the ascending order, i.e. the
	// latest event is at the end of the slice.
	PastNewWalletRegisteredEvents(
		filter *NewWalletRegisteredEventFilter,
	) ([]*NewWalletRegisteredEvent, error)

	// BuildDepositKey calculates a deposit key for the given funding transaction
	// which is a unique identifier for a deposit on-chain.
	BuildDepositKey(fundingTxHash bitcoin.Hash, fundingOutputIndex uint32) *big.Int

	// GetDepositParameters gets the current value of parameters relevant
	// for the depositing process.
	GetDepositParameters() (
		dustThreshold uint64,
		treasuryFeeDivisor uint64,
		txMaxFee uint64,
		revealAheadPeriod uint32,
		err error,
	)

	// SubmitDepositSweepProposalWithReimbursement submits a deposit sweep
	// proposal to the chain. It reimburses the gas cost to the caller.
	SubmitDepositSweepProposalWithReimbursement(
		proposal *tbtc.DepositSweepProposal,
	) error

	// GetDepositSweepMaxSize gets the maximum number of deposits that can
	// be part of a deposit sweep proposal.
	GetDepositSweepMaxSize() (uint16, error)

	// PastRedemptionRequestedEvents fetches past redemption requested events according
	// to the provided filter or unfiltered if the filter is nil. Returned
	// events are sorted by the block number in the ascending order, i.e. the
	// latest event is at the end of the slice.
	PastRedemptionRequestedEvents(
		filter *tbtc.RedemptionRequestedEventFilter,
	) ([]*tbtc.RedemptionRequestedEvent, error)

	// BuildRedemptionKey calculates a redemption key for the given redemption
	// request which is an identifier for a redemption at the given time
	// on-chain.
	BuildRedemptionKey(
		walletPublicKeyHash [20]byte,
		redeemerOutputScript bitcoin.Script,
	) (*big.Int, error)

	// GetRedemptionParameters gets the current value of parameters relevant
	// for the redemption process.
	GetRedemptionParameters() (
		dustThreshold uint64,
		treasuryFeeDivisor uint64,
		txMaxFee uint64,
		txMaxTotalFee uint64,
		timeout uint32,
		timeoutSlashingAmount *big.Int,
		timeoutNotifierRewardMultiplier uint32,
		err error,
	)

	// SubmitRedemptionProposalWithReimbursement submits a redemption proposal
	// to the chain. It reimburses the gas cost to the caller.
	SubmitRedemptionProposalWithReimbursement(
		proposal *tbtc.RedemptionProposal,
	) error

	// GetRedemptionMaxSize gets the maximum number of redemption requests that
	// can be a part of a redemption sweep proposal.
	GetRedemptionMaxSize() (uint16, error)

	// GetRedemptionRequestMinAge get the  minimum time that must elapse since
	// the redemption request creation before a request becomes eligible for
	// a processing.
	GetRedemptionRequestMinAge() (uint32, error)

	tbtc.ValidateDepositSweepProposalChain
}
