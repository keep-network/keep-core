package chain

import (
	"math/big"

	"github.com/keep-network/keep-core/pkg/beacon/relay/member"

	"github.com/keep-network/keep-core/pkg/beacon/relay/config"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/gen/async"
	"github.com/keep-network/keep-core/pkg/subscription"
)

// StakerAddress represents chain-specific address of the staker.
type StakerAddress []byte

// RelayEntryInterface defines the subset of the relay chain interface that
// pertains specifically to submission and retrieval of relay requests and
// entries.
type RelayEntryInterface interface {
	// RequestRelayEntry makes an on-chain request to start generation of a
	// random signature.  An event is generated.
	RequestRelayEntry(seed *big.Int) *async.RelayRequestPromise
	// SubmitRelayEntry submits an entry in the threshold relay and returns a
	// promise to track the submission result. The promise is fulfilled with
	// the entry as seen on-chain, or failed if there is an error submitting
	// the entry.
	SubmitRelayEntry(entry *event.Entry) *async.RelayEntryPromise
	// OnRelayEntryGenerated is a callback that is invoked when an on-chain
	// notification of a new, valid relay entry is seen.
	OnRelayEntryGenerated(
		func(entry *event.Entry),
	) (subscription.EventSubscription, error)
	// OnRelayEntryRequested is a callback that is invoked when an on-chain
	// notification of a new, valid relay request is seen.
	OnRelayEntryRequested(
		func(request *event.Request),
	) (subscription.EventSubscription, error)
}

// GroupSelectionInterface defines the subset of the relay chain interface that
// pertains to relay group selection activities.
type GroupSelectionInterface interface {
	// SubmitTicket submits a ticket corresponding to the virtual staker to
	// the chain, and returns a promise to track the submission. The promise
	// is fulfilled with the entry as seen on-chain, or failed if there is an
	// error submitting the entry.
	SubmitTicket(ticket *Ticket) *async.GroupTicketPromise
	// GetSelectedParticipants returns `GroupSize` slice of addresses of
	// candidates which have been selected to the group.
	GetSelectedParticipants() ([]StakerAddress, error)
}

// GroupRegistrationInterface defines the subset of the relay chain interface
// that pertains to relay group registration activities.
type GroupRegistrationInterface interface {
	// OnGroupRegistered is a callback that is invoked when an on-chain
	// notification of a new, valid group being registered is seen.
	OnGroupRegistered(
		func(groupRegistration *event.GroupRegistration),
	) (subscription.EventSubscription, error)
}

// GroupInterface defines the subset of the relay chain interface that pertains
// specifically to relay group management.
type GroupInterface interface {
	GroupSelectionInterface
	GroupRegistrationInterface
}

// DistributedKeyGenerationInterface defines the subset of the relay chain
// interface that pertains specifically to group formation's distributed key
// generation process.
type DistributedKeyGenerationInterface interface {
	// SubmitDKGResult sends DKG result to a chain, along with signatures over
	// result hash from group participants supporting the result.
	// Signatures are the ordered, collected, and concatenated (65 bytes each)
	// signer hashes.
	// A slice of the participants index (starting with 1), corresponding to
	// the signatures slice is included as well.
	SubmitDKGResult(
		requestID *big.Int,
		participantIndex member.Index,
		dkgResult *DKGResult,
		signatures []byte,
		membersIndex []member.Index,
	) *async.DKGResultSubmissionPromise
	// OnDKGResultSubmitted registers a callback that is invoked when an on-chain
	// notification of a new, valid submitted result is seen.
	OnDKGResultSubmitted(
		func(event *event.DKGResultSubmission),
	) (subscription.EventSubscription, error)
	// IsDKGResultSubmitted checks if a DKG result hash has already been
	// submitted to the chain for the given request ID.
	IsDKGResultSubmitted(
		requestID *big.Int,
	) (bool, error)
	// CalculateDKGResultHash calculates 256-bit hash of DKG result in standard
	// specific for the chain. Operation is performed off-chain.
	CalculateDKGResultHash(dkgResult *DKGResult) (DKGResultHash, error)
}

// Interface represents the interface that the relay expects to interact with
// the anchoring blockchain on.
type Interface interface {
	// GetConfig returns the expected configuration of the threshold relay.
	GetConfig() (*config.Chain, error)

	GroupInterface
	RelayEntryInterface
	DistributedKeyGenerationInterface
}
