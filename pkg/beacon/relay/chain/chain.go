package chain

import (
	"math/big"

	"github.com/keep-network/keep-core/pkg/beacon/relay"
	"github.com/keep-network/keep-core/pkg/beacon/relay/config"
)

// Interface represents the interface that the relay expects to interact
// with the anchoring blockchain on.
type Interface interface {
	// GetConfig returns the expected configuration of the threshold relay.
	GetConfig() (config.Chain, error)
	// SubmitGroupPublicKey submits a 96-byte BLS public key to the blockchain,
	// associated with a string groupID. An error is generally only returned in
	// case of connectivity issues; on-chain errors are reported through event
	// callbacks.
	SubmitGroupPublicKey(groupID string, key [96]byte) error
	// OnGroupPublicKeySubmissionFailed takes a callback that is invoked when
	// an attempted group public key submission has failed. The provided groupID
	// is the id of the group for which the public key submission was attempted,
	// while the errorMsg is the on-chain error message indicating what went
	// wrong.
	OnGroupPublicKeySubmissionFailed(func(groupID string, errorMsg string)) error
	// OnGroupPublicKeySubmitted takes a callback that is invoked when a group
	// public key is submitted successfully. The provided groupID is the id of
	// the group for which the public key was submitted, and the activationBlock
	// is the block at which the group will be considered active in the relay.
	//
	// TODO activation delay may be unnecessary, we'll see.
	OnGroupPublicKeySubmitted(func(groupID string, activationBlock *big.Int)) error
	// SubmitRelayEntryCandidate takes an entry (a group signature and
	// timestamp), constructs a relay response from the elements,
	// and attempts to submit the relay reponse to the chain for consideration.
	// If the chain accepts the response as valid...
	SubmitRelayEntryCandidate(relay.Entry) error
}
