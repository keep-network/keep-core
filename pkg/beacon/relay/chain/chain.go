package chain

import (
	"github.com/keep-network/keep-core/pkg/beacon/relay"
	"github.com/keep-network/keep-core/pkg/beacon/relay/config"
	"github.com/keep-network/keep-core/pkg/gen/async"
)

// Interface represents the interface that the relay expects to interact
// with the anchoring blockchain on.
type Interface interface {
	// GetConfig returns the expected configuration of the threshold relay.
	GetConfig() (config.Chain, error)
	// SubmitGroupPublicKey submits a 96-byte BLS public key to the blockchain,
	// associated with a string groupID. On-chain errors are
	// are reported through the promise.
	SubmitGroupPublicKey(groupID string, key [96]byte) *async.GroupPublicKeyPromise
	// SubmitRelayEntry submits an entry in the threshold relay and returns a
	// promise to track the submission result. The promise is fulfilled with
	// the entry as seen on-chain, or failed if there is an error submitting
	// the entry.
	SubmitRelayEntry(entry *relay.Entry) *async.RelayEntryPromise
	// OnRelayEntryGenerated is a callback that is only invoked when a chain
	// notifies clients of a new, valid relay entry. On invocation, the
	// client determines if they are responsible for generating a signature
	// for the next relay entry.
	OnRelayEntryGenerated(handle func(entry relay.Entry))
}
