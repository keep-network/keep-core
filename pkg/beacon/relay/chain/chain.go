package chain

import (
	"github.com/keep-network/keep-core/pkg/beacon/relay/config"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/gen/async"
	"github.com/keep-network/keep-core/pkg/operator"
	"github.com/keep-network/keep-core/pkg/subscription"
)

// StakerAddress represents chain-specific address of the staker.
type StakerAddress []byte

// GroupMemberIndex is an index of a threshold relay group member.
// Maximum value accepted by the chain is 255.
type GroupMemberIndex = uint8

// RelayEntryInterface defines the subset of the relay chain interface that
// pertains specifically to submission and retrieval of relay requests and
// entries.
type RelayEntryInterface interface {
	// SubmitRelayEntry submits an entry in the threshold relay and returns a
	// promise to track the submission progress. The promise is fulfilled when
	// the entry has been successfully submitted to the on-chain, or failed if
	// the entry submission failed.
	SubmitRelayEntry(entry []byte) *async.EventEntrySubmittedPromise
	// OnRelayEntrySubmitted is a callback that is invoked when an on-chain
	// notification of a new, valid relay entry is seen.
	OnRelayEntrySubmitted(
		func(entry *event.EntrySubmitted),
	) (subscription.EventSubscription, error)
	// OnRelayEntryRequested is a callback that is invoked when an on-chain
	// notification of a new, valid relay request is seen.
	OnRelayEntryRequested(
		func(request *event.Request),
	) (subscription.EventSubscription, error)
	// ReportRelayEntryTimeout notifies the chain when a selected group which was
	// supposed to submit a relay entry, did not deliver it within a specified
	// time frame (relayEntryTimeout) counted in blocks.
	ReportRelayEntryTimeout() error
	// IsEntryInProgress checks if a new relay entry is currently in progress.
	IsEntryInProgress() (bool, error)
}

// GroupSelectionInterface defines the subset of the relay chain interface that
// pertains to relay group selection activities.
type GroupSelectionInterface interface {
	// OnGroupSelectionStarted is a callback that is invoked when an on-chain
	// group selection started and the contract is ready to accept tickets.
	OnGroupSelectionStarted(
		func(groupSelectionStarted *event.GroupSelectionStart),
	) (subscription.EventSubscription, error)
	// SubmitTicket submits a ticket corresponding to the virtual staker to
	// the chain, and returns a promise to track the submission. The promise
	// is fulfilled with the entry as seen on-chain, or failed if there is an
	// error submitting the entry.
	SubmitTicket(ticket *Ticket) *async.EventGroupTicketSubmissionPromise
	// GetSubmittedTickets gets the submitted group candidate tickets so far.
	GetSubmittedTickets() ([]uint64, error)
	// GetSelectedParticipants returns `GroupSize` slice of addresses of
	// candidates which have been selected to the currently assembling group.
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
	// Checks if a group with the given public key is considered as
	// stale on-chain. Group is considered as stale if it is expired and when
	// its expiration time and potentially executed operation timeout are both
	// in the past. Stale group is never selected by the chain to any new
	// operation.
	IsStaleGroup(groupPublicKey []byte) (bool, error)
	// GetGroupMembers returns `GroupSize` slice of addresses of
	// participants which have been selected to the group with given public key.
	GetGroupMembers(groupPublicKey []byte) ([]StakerAddress, error)
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
	// Signatures over DKG result hash are collected in a map keyed by signer's
	// member index.
	SubmitDKGResult(
		participantIndex GroupMemberIndex,
		dkgResult *DKGResult,
		signatures map[GroupMemberIndex][]byte,
	) *async.EventDKGResultSubmissionPromise
	// OnDKGResultSubmitted registers a callback that is invoked when an on-chain
	// notification of a new, valid submitted result is seen.
	OnDKGResultSubmitted(
		func(event *event.DKGResultSubmission),
	) (subscription.EventSubscription, error)
	// IsGroupRegistered checks if group with the given public key is registered
	// on-chain.
	IsGroupRegistered(groupPublicKey []byte) (bool, error)
	// CalculateDKGResultHash calculates 256-bit hash of DKG result in standard
	// specific for the chain. Operation is performed off-chain.
	CalculateDKGResultHash(dkgResult *DKGResult) (DKGResultHash, error)
}

// Interface represents the interface that the relay expects to interact with
// the anchoring blockchain on.
type Interface interface {
	// GetConfig returns the expected configuration of the threshold relay.
	GetConfig() (*config.Chain, error)
	// GetKeys returns the key pair used to attest for messages being sent to
	// the chain.
	GetKeys() (*operator.PrivateKey, *operator.PublicKey)

	GroupInterface
	RelayEntryInterface
	DistributedKeyGenerationInterface
}
