package tbtcpg

import (
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/chain"
	"math/big"
	"time"

	"github.com/keep-network/keep-core/pkg/tbtc"
)

// Chain represents the interface that the wallet maintainer module expects
// to interact with the anchoring blockchain on.
type Chain interface {
	// GetDepositRequest gets the on-chain deposit request for the given
	// funding transaction hash and output index.The returned values represent:
	// - deposit request which is non-nil only when the deposit request was
	//   found,
	// - boolean value which is true if the deposit request was found, false
	//   otherwise,
	// - error which is non-nil only when the function execution failed. It will
	//   be nil if the deposit request was not found, but the function execution
	//   succeeded.
	GetDepositRequest(
		fundingTxHash bitcoin.Hash,
		fundingOutputIndex uint32,
	) (*tbtc.DepositChainRequest, bool, error)

	// PastNewWalletRegisteredEvents fetches past new wallet registered events
	// according to the provided filter or unfiltered if the filter is nil. Returned
	// events are sorted by the block number in the ascending order, i.e. the
	// latest event is at the end of the slice.
	PastNewWalletRegisteredEvents(
		filter *tbtc.NewWalletRegisteredEventFilter,
	) ([]*tbtc.NewWalletRegisteredEvent, error)

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

	// GetRedemptionMaxSize gets the maximum number of redemption requests that
	// can be a part of a redemption sweep proposal.
	GetRedemptionMaxSize() (uint16, error)

	// GetRedemptionRequestMinAge get the  minimum time that must elapse since
	// the redemption request creation before a request becomes eligible for
	// a processing.
	GetRedemptionRequestMinAge() (uint32, error)

	// PastDepositRevealedEvents fetches past deposit reveal events according
	// to the provided filter or unfiltered if the filter is nil. Returned
	// events are sorted by the block number in the ascending order, i.e. the
	// latest event is at the end of the slice.
	PastDepositRevealedEvents(
		filter *tbtc.DepositRevealedEventFilter,
	) ([]*tbtc.DepositRevealedEvent, error)

	// GetPendingRedemptionRequest gets the on-chain pending redemption request
	// for the given wallet public key hash and redeemer output script.
	// The returned bool value indicates whether the request was found or not.
	GetPendingRedemptionRequest(
		walletPublicKeyHash [20]byte,
		redeemerOutputScript bitcoin.Script,
	) (*tbtc.RedemptionRequest, bool, error)

	// ValidateDepositSweepProposal validates the given deposit sweep proposal
	// against the chain. It requires some additional data about the deposits
	// that must be fetched externally. Returns an error if the proposal is
	// not valid or nil otherwise.
	ValidateDepositSweepProposal(
		proposal *tbtc.DepositSweepProposal,
		depositsExtraInfo []struct {
			*tbtc.Deposit
			FundingTx *bitcoin.Transaction
		},
	) error

	// ValidateRedemptionProposal validates the given redemption proposal
	// against the chain. Returns an error if the proposal is not valid or
	// nil otherwise.
	ValidateRedemptionProposal(proposal *tbtc.RedemptionProposal) error

	// GetDepositSweepMaxSize gets the maximum number of deposits that can
	// be part of a deposit sweep proposal.
	GetDepositSweepMaxSize() (uint16, error)

	BlockCounter() (chain.BlockCounter, error)

	AverageBlockTime() time.Duration

	// ValidateHeartbeatProposal validates the given heartbeat proposal
	// against the chain. Returns an error if the proposal is not valid or
	// nil otherwise.
	ValidateHeartbeatProposal(proposal *tbtc.HeartbeatProposal) error
}
