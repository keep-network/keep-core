package result

import (
	"github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
	"github.com/keep-network/keep-core/pkg/operator"
)

// DKGResultHashSignatureMessage is a message payload that carries a hash of
// the DKG result and a signature over this hash for a DKG result.
//
// It is expected to be broadcast within the group.
type DKGResultHashSignatureMessage struct {
	// Index of the sender in the group.
	senderIndex group.MemberIndex
	// Hash of the DKG result preferred by the sender.
	resultHash chain.DKGResultHash
	// Signature over the DKG result hash calculated by the sender.
	signature operator.Signature
	// Public key of the sender. It will be used to verify the signature by
	// the receiver.
	publicKey *operator.PublicKey
}

// SenderID returns protocol-level identifier of the message sender.
func (m *DKGResultHashSignatureMessage) SenderID() group.MemberIndex {
	return m.senderIndex
}
