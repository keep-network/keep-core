package dkg2

import (
	"github.com/keep-network/keep-core/pkg/beacon/relay/chain"
)

// DKGResultHashSignatureMessage is a message payload that carries a hash of
// the DKG result and a signature over this hash for a DKG result.
//
// It is expected to be broadcast within the group.
type DKGResultHashSignatureMessage struct {
	senderIndex uint32 // i

	resultHash chain.DKGResultHash
	signature  []byte
}
