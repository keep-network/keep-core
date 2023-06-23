package coordinator

import (
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/tbtc"
	"math/big"
)

// Chain represents the interface that the coordinator module expects to interact
// with the anchoring blockchain on.
type Chain interface {
	// TODO: Change to something more specific once https://github.com/keep-network/keep-core/issues/3632
	// is handled.
	tbtc.Chain

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
}
