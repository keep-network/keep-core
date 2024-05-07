package inactivity

import (
	"github.com/keep-network/keep-core/pkg/protocol/group"
)

const messageTypePrefix = "protocol_inactivity/"

// message holds common traits of all inactivity protocol messages.
type message interface {
	// SenderID returns protocol-level identifier of the message sender.
	SenderID() group.MemberIndex
	// SessionID returns the session identifier of the message.
	SessionID() string
	// Type returns the exact type of the message.
	Type() string
}

type claimSignatureMessage struct {
	senderID group.MemberIndex

	claimHash ClaimSignatureHash
	signature []byte
	publicKey []byte
	sessionID string
}

// SenderID returns protocol-level identifier of the message sender.
func (csm *claimSignatureMessage) SenderID() group.MemberIndex {
	return csm.senderID
}

// SessionID returns the session identifier of the message.
func (csm *claimSignatureMessage) SessionID() string {
	return csm.sessionID
}

// Type returns a string describing an claimSignatureMessage type for
// marshaling purposes.
func (csm *claimSignatureMessage) Type() string {
	return messageTypePrefix + "claim_signature_message"
}
