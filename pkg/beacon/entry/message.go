package entry

import (
	"github.com/keep-network/keep-core/pkg/beacon/group"
)

// SignatureShareMessage is a message payload that carries the sender's
// signature share for the given message.
type SignatureShareMessage struct {
	senderID   group.MemberIndex
	shareBytes []byte
}

// NewSignatureShareMessage creates new SignatureShareMessage.
func NewSignatureShareMessage(
	senderID group.MemberIndex,
	shareBytes []byte,
) *SignatureShareMessage {
	return &SignatureShareMessage{senderID, shareBytes}
}

// SenderID returns protocol-level identifier of the message sender.
func (ssm *SignatureShareMessage) SenderID() group.MemberIndex {
	return ssm.senderID
}
