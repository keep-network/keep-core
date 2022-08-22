package dkg

import (
	"github.com/keep-network/keep-core/pkg/crypto/ephemeral"
	"github.com/keep-network/keep-core/pkg/protocol/group"
)

const messageTypePrefix = "tecdsa_dkg/"

// ephemeralPublicKeyMessage is a message payload that carries the sender's
// ephemeral public keys generated for all other group members.
//
// The receiver performs ECDH on a sender's ephemeral public key intended for
// the receiver and on the receiver's private ephemeral key, creating a symmetric
// key used for encrypting a conversation between the sender and the receiver.
type ephemeralPublicKeyMessage struct {
	senderID group.MemberIndex

	ephemeralPublicKeys map[group.MemberIndex]*ephemeral.PublicKey
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

// tssRoundOneMessage is a message payload that carries the sender's TSS
// commitments and the Paillier public key.
type tssRoundOneMessage struct {
	senderID group.MemberIndex

	payload   []byte
	sessionID string
}

// SenderID returns protocol-level identifier of the message sender.
func (trom *tssRoundOneMessage) SenderID() group.MemberIndex {
	return trom.senderID
}

// Type returns a string describing an tssRoundOneMessage type for
// marshaling purposes.
func (trom *tssRoundOneMessage) Type() string {
	return messageTypePrefix + "tss_round_one_message"
}

// tssRoundTwoMessage is a message payload that carries the sender's TSS
// shares and de-commitments.
type tssRoundTwoMessage struct {
	senderID group.MemberIndex

	broadcastPayload []byte
	peersPayload     map[group.MemberIndex][]byte
	sessionID        string
}

// SenderID returns protocol-level identifier of the message sender.
func (trtm *tssRoundTwoMessage) SenderID() group.MemberIndex {
	return trtm.senderID
}

// Type returns a string describing an tssRoundTwoMessage type for
// marshaling purposes.
func (trtm *tssRoundTwoMessage) Type() string {
	return messageTypePrefix + "tss_round_two_message"
}

// tssRoundThreeMessage is a message payload that carries the sender's TSS
// Paillier proof.
type tssRoundThreeMessage struct {
	senderID group.MemberIndex

	payload   []byte
	sessionID string
}

// SenderID returns protocol-level identifier of the message sender.
func (trtm *tssRoundThreeMessage) SenderID() group.MemberIndex {
	return trtm.senderID
}

// Type returns a string describing an tssRoundThreeMessage type for
// marshaling purposes.
func (trtm *tssRoundThreeMessage) Type() string {
	return messageTypePrefix + "tss_round_three_message"
}

// resultSignatureMessage is a message payload that carries a hash of
// the DKG result and a signature over this hash for a DKG result.
//
// It is expected to be broadcast within the group.
type resultSignatureMessage struct {
	senderID group.MemberIndex

	resultHash ResultHash
	signature  []byte
	publicKey  []byte
	sessionID  string
}

// SenderID returns protocol-level identifier of the message sender.
func (rsm *resultSignatureMessage) SenderID() group.MemberIndex {
	return rsm.senderID
}

// Type returns a string describing an resultSignatureMessage type for
// marshaling purposes.
func (rsm *resultSignatureMessage) Type() string {
	return messageTypePrefix + "result_signature_message"
}
