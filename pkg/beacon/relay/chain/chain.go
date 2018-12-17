package chain

import (
	"math/big"

	"github.com/keep-network/keep-core/pkg/beacon/relay/config"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/gen/async"
)

// Interface represents the interface that the relay expects to interact
// with the anchoring blockchain on.
type Interface interface {
	// GetConfig returns the expected configuration of the threshold relay.
	GetConfig() (config.Chain, error)
	// SubmitGroupPublicKey submits a 96-byte BLS public key to the blockchain,
	// associated with a request with id requestID. On-chain errors are reported
	// through the promise.
	SubmitGroupPublicKey(requestID *big.Int, key [96]byte) *async.GroupRegistrationPromise
	// SubmitRelayEntry submits an entry in the threshold relay and returns a
	// promise to track the submission result. The promise is fulfilled with
	// the entry as seen on-chain, or failed if there is an error submitting
	// the entry.
	SubmitRelayEntry(entry *event.Entry) *async.RelayEntryPromise
	// OnRelayEntryGenerated is a callback that is invoked when an on-chain
	// notification of a new, valid relay entry is seen.
	OnRelayEntryGenerated(func(entry *event.Entry))
	// OnRelayEntryRequested is a callback that is invoked when an on-chain
	// notification of a new, valid relay request is seen.
	OnRelayEntryRequested(func(request *event.Request))
	// OnStakerAdded is a callback that is invoked when an on-chain
	// notification of a new, valid staker is seen.
	OnStakerAdded(func(staker *event.StakerRegistration))
	// AddStaker is a temporary function for Milestone 1 that
	// adds a staker to the group contract.
	AddStaker(groupMemberID string) *async.StakerRegistrationPromise
	// GetStakerList is a temporary function for Milestone 1 that
	// gets back the list of stakers.
	GetStakerList() ([]string, error)
	// OnGroupRegistered is a callback that is invoked when an on-chain
	// notification of a new, valid group being registered is seen.
	OnGroupRegistered(func(key *event.GroupRegistration))
	// RequestRelayEntry makes an on-chain request to start generation of a
	// random signature.  An event is generated.
	RequestRelayEntry(blockReward, seed *big.Int) *async.RelayRequestPromise
	// SubmitDKGResult sends DKG result to a chain.
	SubmitDKGResult(requestID *big.Int, dkgResult *DKGResult) *async.DKGResultPublicationPromise
	// OnDKGResultPublished is a callback that is invoked when an on-chain
	// notification of a new, valid published result is seen.
	OnDKGResultPublished(func(dkgResultPublication *event.DKGResultPublication))
	// IsDKGResultPublished checks if the specific DKG result has already been
	// published to a chain for given request ID.
	IsDKGResultPublished(requestID *big.Int, dkgResult *DKGResult) bool
	// GetDKGSubmissions returns the entire set of submissions that is related
	// to a particular requestID.
	GetDKGSubmissions(requestID *big.Int) *Submissions
	// Vote increases the number of votes for a dkgResultHash
	Vote(requestID *big.Int, dkgResultHash []byte)
	// OnDKGResultVote registers a callback when a vote occurs.
	OnDKGResultVote(func(dkgResultVote *event.DKGResultVote))
	// CurrentBlock rturns the current block number
	CurrentBlock() (int, error)
}
