package chain

import (
	"math/big"

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
	RequestRelayEntry(blockReward, seed *big.Int) *async.RelayRequestPromise
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
	// SubmitGroupPublicKey submits a public key to the blockchain, associated
	// with a request with id requestID. On-chain errors are reported through
	// the promise.
	SubmitGroupPublicKey(
		requestID *big.Int,
		groupPublicKey []byte,
	) *async.GroupRegistrationPromise
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
	// SubmitDKGResult sends DKG result to a chain.
	SubmitDKGResult(
		requestID *big.Int,
		dkgResult *DKGResult,
	) *async.DKGResultPublicationPromise
	// OnDKGResultPublished is a callback that is invoked when an on-chain
	// notification of a new, valid published result is seen.
	OnDKGResultPublished(
		func(dkgResultPublication *event.DKGResultPublication),
	) (subscription.EventSubscription, error)
	// IsDKGResultPublished checks if any DKG result has already been published
	// to a chain for the given request ID.
	IsDKGResultPublished(requestID *big.Int) (bool, error)
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
