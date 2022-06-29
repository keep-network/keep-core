package chain

import (
	"math/big"

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
	) subscription.EventSubscription
	// OnRelayEntryRequested is a callback that is invoked when an on-chain
	// notification of a new, valid relay request is seen.
	OnRelayEntryRequested(
		func(request *event.Request),
	) subscription.EventSubscription
	// ReportRelayEntryTimeout notifies the chain when a selected group which was
	// supposed to submit a relay entry, did not deliver it within a specified
	// time frame (relayEntryTimeout) counted in blocks.
	ReportRelayEntryTimeout() error
	// IsEntryInProgress checks if a new relay entry is currently in progress.
	IsEntryInProgress() (bool, error)
	// CurrentRequestStartBlock returns a start block of a current entry.
	CurrentRequestStartBlock() (*big.Int, error)
	// CurrentRequestPreviousEntry returns previous entry of a current request.
	CurrentRequestPreviousEntry() ([]byte, error)
	// CurrentRequestGroupPublicKey returns group public key for the current request.
	CurrentRequestGroupPublicKey() ([]byte, error)
}

// GroupSelectionInterface defines the subset of the relay chain interface that
// pertains to relay group selection activities.
type GroupSelectionInterface interface {
	// TODO: Expose a function allowing to check the selected group for the
	//       given seed.
}

// GroupRegistrationInterface defines the subset of the relay chain interface
// that pertains to relay group registration activities.
type GroupRegistrationInterface interface {
	// OnGroupRegistered is a callback that is invoked when an on-chain
	// notification of a new, valid group being registered is seen.
	OnGroupRegistered(
		func(groupRegistration *event.GroupRegistration),
	) subscription.EventSubscription
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
	) subscription.EventSubscription
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
	GetConfig() *Config
	// GetKeys returns the key pair used to attest for messages being sent to
	// the chain.
	GetKeys() (*operator.PrivateKey, *operator.PublicKey, error)
	// MinimumStake returns the current on-chain value representing the minimum
	// necessary amount of KEEP a client must lock up to participate in the
	// threshold relay. This value can change over time according to the minimum
	// stake schedule.
	MinimumStake() (*big.Int, error)

	GroupInterface
	RelayEntryInterface
	DistributedKeyGenerationInterface
}

// Config contains the config data needed for the relay to operate.
type Config struct {
	// GroupSize is the size of a group in the threshold relay.
	GroupSize int
	// HonestThreshold is the minimum number of active participants behaving
	// according to the protocol needed to generate a new relay entry.
	HonestThreshold int
	// ResultPublicationBlockStep is the duration (in blocks) that has to pass
	// before group member with the given index is eligible to submit the
	// result.
	// Nth player becomes eligible to submit the result after
	// T_dkg + (N-1) * T_step
	// where T_dkg is time for phases 1-12 to complete and T_step is the result
	// publication block step.
	ResultPublicationBlockStep uint64
	// RelayEntryTimeout is a timeout in blocks on-chain for a relay
	// entry to be published by the selected group. Blocks are
	// counted from the moment relay request occur.
	RelayEntryTimeout uint64
}

// DishonestThreshold is the maximum number of misbehaving participants for
// which it is still possible to generate a new relay entry.
// Misbehaviour is any misconduct to the protocol, including inactivity.
func (c *Config) DishonestThreshold() int {
	return c.GroupSize - c.HonestThreshold
}
