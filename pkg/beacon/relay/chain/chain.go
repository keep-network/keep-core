package chain

import (
	"github.com/keep-network/keep-core/pkg/beacon/relay/config"
	"github.com/keep-network/keep-core/pkg/gen/async"
)

// Interface represents the interface that the relay expects to interact
// with the anchoring blockchain on.
type Interface interface {
	// SubmitGroupPublicKey submits a 96-byte BLS public key to the blockchain,
	// associated with a string groupID. An error is generally only returned in
	// case of connectivity issues; on-chain errors are reported through promise.
	SubmitGroupPublicKey(groupID string, key [96]byte) (*async.GroupPublicKeyPromise, error)

	// GetConfig returns the expected configuration of the threshold relay.
	GetConfig() (config.Chain, error)
}
