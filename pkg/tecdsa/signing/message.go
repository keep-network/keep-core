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

// tssRoundOneMessage is a message payload that carries the sender's
// TSS round one components.
type tssRoundOneMessage struct {
	senderID group.MemberIndex

	broadcastPayload []byte
	peersPayload     map[group.MemberIndex][]byte
	sessionID        string
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

// tssRoundTwoMessage is a message payload that carries the sender's
// TSS round two components.
type tssRoundTwoMessage struct {
	senderID group.MemberIndex

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

// tssRoundThreeMessage is a message payload that carries the sender's
// TSS round three components.
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

// tssRoundFourMessage is a message payload that carries the sender's
// TSS round four components.
type tssRoundFourMessage struct {
	senderID group.MemberIndex

	payload   []byte
	sessionID string
}

// SenderID returns protocol-level identifier of the message sender.
func (trfm *tssRoundFourMessage) SenderID() group.MemberIndex {
	return trfm.senderID
}

// Type returns a string describing an tssRoundFourMessage type for
// marshaling purposes.
func (trfm *tssRoundFourMessage) Type() string {
	return messageTypePrefix + "tss_round_four_message"
}

// tssRoundFiveMessage is a message payload that carries the sender's
// TSS round five components.
type tssRoundFiveMessage struct {
	senderID group.MemberIndex

	payload   []byte
	sessionID string
}

// SenderID returns protocol-level identifier of the message sender.
func (trfm *tssRoundFiveMessage) SenderID() group.MemberIndex {
	return trfm.senderID
}

// Type returns a string describing an tssRoundFiveMessage type for
// marshaling purposes.
func (trfm *tssRoundFiveMessage) Type() string {
	return messageTypePrefix + "tss_round_five_message"
}

// tssRoundSixMessage is a message payload that carries the sender's
// TSS round six components.
type tssRoundSixMessage struct {
	senderID group.MemberIndex

	payload   []byte
	sessionID string
}

// SenderID returns protocol-level identifier of the message sender.
func (trsm *tssRoundSixMessage) SenderID() group.MemberIndex {
	return trsm.senderID
}

// Type returns a string describing an tssRoundSixMessage type for
// marshaling purposes.
func (trfm *tssRoundSixMessage) Type() string {
	return messageTypePrefix + "tss_round_six_message"
}