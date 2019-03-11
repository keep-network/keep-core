package resultsubmission

import (
	"github.com/keep-network/keep-core/pkg/beacon/relay/chain"
)

// Signature is a signature of the DKG result hash.
// TODO: When we have signing implemented update aliased type if needed.
type Signature []byte

// DKGResultHashSignatureMessage is a message payload that carries a hash of
// the DKG result and a signature over this hash for a DKG result.
//
// It is expected to be broadcast within the group.
type DKGResultHashSignatureMessage struct {
	// Index of the sender in the group.
	senderIndex uint32
	// Hash of the DKG result preferred by the sender.
	resultHash chain.DKGResultHash
	// Signature over the DKG result hash calculated by the sender.
	signature Signature
}
