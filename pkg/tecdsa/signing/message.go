package signing

import (
	"github.com/keep-network/keep-core/pkg/crypto/ephemeral"
	"github.com/keep-network/keep-core/pkg/protocol/group"
)

const messageTypePrefix = "tecdsa_signing/"

// ephemeralPublicKeyMessage is a message payload that carries the sender's
// ephemeral public keys generated for all other group members.
//
// The receiver performs ECDH on a sender's ephemeral public key intended for
// the receiver and on the receiver's private ephemeral key, creating a symmetric
// key used for encrypting a conversation between the sender and the receiver.
type ephemeralPublicKeyMessage struct {
	senderID group.MemberIndex

	ephemeralPublicKeys map[group.MemberIndex]*ephemeral.PublicKey
	sessionID           string
}

// SenderID returns protocol-level identifier of the message sender.
func (epkm *ephemeralPublicKeyMessage) SenderID() group.MemberIndex {
	return epkm.senderID
}

// Type returns a string describing an ephemeralPublicKeyMessage type for
// marshaling purposes.
func (epkm *ephemeralPublicKeyMessage) Type() string {
	return messageTypePrefix + "ephemeral_public_key_message"
}
