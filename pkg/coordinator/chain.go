package coordinator

import (
	"errors"
	"math/big"
	"time"

	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/tbtc"
)

var (
	// ErrPendingRedemptionRequestNotFound is an error that is returned if
	// a pending redemption request was not found on-chain for the given redemption
	// key.
	ErrPendingRedemptionRequestNotFound = errors.New(
		"no pending redemption requests for the given key",
	)
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
		filter *RedemptionRequestedEventFilter,
	) ([]*RedemptionRequestedEvent, error)

	// GetPendingRedemptionRequest gets the on-chain pending redemption request
	// for the given redeemer output script and wallet public key hash.
	// Returns an ErrPendingRedemptionRequestNotFound error if a redemption request
	// was not found.
	GetPendingRedemptionRequest(
		redeemerOutputScript []byte,
		walletPublicKeyHash [20]byte,
	) (*RedemptionChainRequest, error)

	// BuildRedemptionKey calculates a redemption key for the given redemption
	// request which is an identifier for a redemption at the given time
	// on-chain.
	BuildRedemptionKey(redeemerOutputScript []byte, walletPublicKeyHash [20]byte) *big.Int

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
		proposal *RedemptionProposal,
	) error

	// GetRedemptionMaxSize gets the maximum number of redemption requests that
	// can be a part of a redemption sweep proposal.
	GetRedemptionMaxSize() (uint16, error)

	// GetRedemptionRequestMinAge get the  minimum time that must elapse since
	// the redemption request creation before a request becomes eligible for
	// a processing.
	GetRedemptionRequestMinAge() (uint32, error)
}

// RedemptionRequestedEvent represents a redemption requested event.
type RedemptionRequestedEvent struct {
	WalletPublicKeyHash  [20]byte
	RedeemerOutputScript []byte
	Redeemer             chain.Address
	RequestedAmount      uint64
	TreasuryFee          uint64
	TxMaxFee             uint64
	BlockNumber          uint64
}

// RedemptionRequestedEventFilter is a component allowing to filter RedemptionRequestedEvent.
type RedemptionRequestedEventFilter struct {
	StartBlock          uint64
	EndBlock            *uint64
	WalletPublicKeyHash [][20]byte
	Redeemer            []chain.Address
}

// RedemptionChainRequest represents a redemption request stored on-chain.
type RedemptionChainRequest struct {
	Redeemer        chain.Address
	RequestedAmount uint64
	TreasuryFee     uint64
	TxMaxFee        uint64
	RequestedAt     time.Time
}

// RedemptionProposal represents a redemption proposal submitted to the chain.
type RedemptionProposal struct {
	WalletPublicKeyHash    [20]byte
	RedeemersOutputScripts [][]byte
	RedemptionTxFee        *big.Int
}
