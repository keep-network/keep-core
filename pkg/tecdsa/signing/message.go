package signing

import (
	"github.com/keep-network/keep-core/pkg/crypto/ephemeral"
	"github.com/keep-network/keep-core/pkg/protocol/group"
)

const messageTypePrefix = "tecdsa_signing/"

// message holds common traits of all signing protocol messages.
type message interface {
	// SenderID returns protocol-level identifier of the message sender.
	SenderID() group.MemberIndex
	// SessionID returns the session identifier of the message.
	SessionID() string
	// Type returns the exact type of the message.
	Type() string
}

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

// SessionID returns the session identifier of the message.
func (epkm *ephemeralPublicKeyMessage) SessionID() string {
	return epkm.sessionID
}

// Type returns a string describing an ephemeralPublicKeyMessage type for
// marshaling purposes.
func (epkm *ephemeralPublicKeyMessage) Type() string {
	return messageTypePrefix + "ephemeral_public_key_message"
}

// tssRoundOneCompositeMessage is a composite of tssRoundOneMessages keyed by
// the message being signed.
type tssRoundOneCompositeMessage struct {
	senderID  group.MemberIndex
	sessionID string

	tssRoundOneMessages map[string]*tssRoundOneMessage
}

// tssRoundOneMessage is a message payload that carries the sender's
// TSS round one components.
type tssRoundOneMessage struct {
	broadcastPayload []byte
	peersPayload     map[group.MemberIndex][]byte
}

// SenderID returns protocol-level identifier of the message sender.
func (trocm *tssRoundOneCompositeMessage) SenderID() group.MemberIndex {
	return trocm.senderID
}

// SessionID returns the session identifier of the message.
func (trocm *tssRoundOneCompositeMessage) SessionID() string {
	return trocm.sessionID
}

// Type returns a string describing a tssRoundOneCompositeMessage type for
// marshaling purposes.
func (trocm *tssRoundOneCompositeMessage) Type() string {
	return messageTypePrefix + "tss_round_one_composite_message"
}

// tssRoundTwoCompositeMessage is a composite of tssRoundTwoMessages keyed by the
// message being signed.
type tssRoundTwoCompositeMessage struct {
	senderID  group.MemberIndex
	sessionID string

	tssRoundTwoMessages map[string]*tssRoundTwoMessage
}

// tssRoundTwoMessage is a message payload that carries the sender's
// TSS round two components.
type tssRoundTwoMessage struct {
	peersPayload map[group.MemberIndex][]byte
}

// SenderID returns protocol-level identifier of the message sender.
func (trtcm *tssRoundTwoCompositeMessage) SenderID() group.MemberIndex {
	return trtcm.senderID
}

// SessionID returns the session identifier of the message.
func (trtcm *tssRoundTwoCompositeMessage) SessionID() string {
	return trtcm.sessionID
}

// Type returns a string describing a tssRoundTwoCompositeMessage type for
// marshaling purposes.
func (trtcm *tssRoundTwoCompositeMessage) Type() string {
	return messageTypePrefix + "tss_round_two_composite_message"
}

// tssRoundThreeCompositeMessage is a composite of tssRoundThreeMessages keyed
// by the message being signed.
type tssRoundThreeCompositeMessage struct {
	senderID  group.MemberIndex
	sessionID string

	tssRoundThreeMessages map[string]*tssRoundThreeMessage
}

// tssRoundThreeMessage is a message payload that carries the sender's
// TSS round three components.
type tssRoundThreeMessage struct {
	broadcastPayload []byte
}

// SenderID returns protocol-level identifier of the message sender.
func (trtcm *tssRoundThreeCompositeMessage) SenderID() group.MemberIndex {
	return trtcm.senderID
}

// SessionID returns the session identifier of the message.
func (trtcm *tssRoundThreeCompositeMessage) SessionID() string {
	return trtcm.sessionID
}

// Type returns a string describing a tssRoundThreeCompositeMessage type for
// marshaling purposes.
func (trtcm *tssRoundThreeCompositeMessage) Type() string {
	return messageTypePrefix + "tss_round_three_composite_message"
}

// tssRoundFourCompositeMessage is a composite of tssRoundFourMessages keyed by
// the message being signed.
type tssRoundFourCompositeMessage struct {
	senderID  group.MemberIndex
	sessionID string

	tssRoundFourMessages map[string]*tssRoundFourMessage
}

// tssRoundFourMessage is a message payload that carries the sender's
// TSS round four components.
type tssRoundFourMessage struct {
	broadcastPayload []byte
}

// SenderID returns protocol-level identifier of the message sender.
func (trfcm *tssRoundFourCompositeMessage) SenderID() group.MemberIndex {
	return trfcm.senderID
}

// SessionID returns the session identifier of the message.
func (trfcm *tssRoundFourCompositeMessage) SessionID() string {
	return trfcm.sessionID
}

// Type returns a string describing a tssRoundFourCompositeMessage type for
// marshaling purposes.
func (trfcm *tssRoundFourCompositeMessage) Type() string {
	return messageTypePrefix + "tss_round_four_composite_message"
}

// tssRoundFiveCompositeMessage is a composite of tssRoundFiveMessage keyed by
// the message being signed.
type tssRoundFiveCompositeMessage struct {
	senderID  group.MemberIndex
	sessionID string

	tssRoundFiveMessages map[string]*tssRoundFiveMessage
}

// tssRoundFiveMessage is a message payload that carries the sender's
// TSS round five components.
type tssRoundFiveMessage struct {
	broadcastPayload []byte
}

// SenderID returns protocol-level identifier of the message sender.
func (trfcm *tssRoundFiveCompositeMessage) SenderID() group.MemberIndex {
	return trfcm.senderID
}

// SessionID returns the session identifier of the message.
func (trfcm *tssRoundFiveCompositeMessage) SessionID() string {
	return trfcm.sessionID
}

// Type returns a string describing a tssRoundFiveCompositeMessage type for
// marshaling purposes.
func (trfcm *tssRoundFiveCompositeMessage) Type() string {
	return messageTypePrefix + "tss_round_five_composite_message"
}

// tssRoundSixCompositeMessage is a composite of tssRoundSixMessage keyed by
// the message being signed.
type tssRoundSixCompositeMessage struct {
	senderID  group.MemberIndex
	sessionID string

	tssRoundSixMessages map[string]*tssRoundSixMessage
}

// tssRoundSixMessage is a message payload that carries the sender's
// TSS round six components.
type tssRoundSixMessage struct {
	broadcastPayload []byte
}

// SenderID returns protocol-level identifier of the message sender.
func (trscm *tssRoundSixCompositeMessage) SenderID() group.MemberIndex {
	return trscm.senderID
}

// SessionID returns the session identifier of the message.
func (trscm *tssRoundSixCompositeMessage) SessionID() string {
	return trscm.sessionID
}

// Type returns a string describing a tssRoundSixCompositeMessage type for
// marshaling purposes.
func (trscm *tssRoundSixCompositeMessage) Type() string {
	return messageTypePrefix + "tss_round_six_composite_message"
}

// tssRoundSevenCompositeMessage is a composite of tssRoundSevenMessages keyed
// by the message being signed.
type tssRoundSevenCompositeMessage struct {
	senderID  group.MemberIndex
	sessionID string

	tssRoundSevenMessages map[string]*tssRoundSevenMessage
}

// tssRoundSevenMessage is a message payload that carries the sender's
// TSS round seven components.
type tssRoundSevenMessage struct {
	broadcastPayload []byte
}

// SenderID returns protocol-level identifier of the message sender.
func (trscm *tssRoundSevenCompositeMessage) SenderID() group.MemberIndex {
	return trscm.senderID
}

// SessionID returns the session identifier of the message.
func (trscm *tssRoundSevenCompositeMessage) SessionID() string {
	return trscm.sessionID
}

// Type returns a string describing a tssRoundSevenCompositeMessage type for
// marshaling purposes.
func (trscm *tssRoundSevenCompositeMessage) Type() string {
	return messageTypePrefix + "tss_round_seven_composite_message"
}

// tssRoundEightCompositeMessage is a composite of tssRoundEightMessages keyed
// by the message being signed.
type tssRoundEightCompositeMessage struct {
	senderID  group.MemberIndex
	sessionID string

	tssRoundEightMessages map[string]*tssRoundEightMessage
}

// tssRoundEightMessage is a message payload that carries the sender's
// TSS round eight components.
type tssRoundEightMessage struct {
	broadcastPayload []byte
}

// SenderID returns protocol-level identifier of the message sender.
func (trecm *tssRoundEightCompositeMessage) SenderID() group.MemberIndex {
	return trecm.senderID
}

// SessionID returns the session identifier of the message.
func (trecm *tssRoundEightCompositeMessage) SessionID() string {
	return trecm.sessionID
}

// Type returns a string describing a tssRoundEightCompositeMessage type for
// marshaling purposes.
func (trecm *tssRoundEightCompositeMessage) Type() string {
	return messageTypePrefix + "tss_round_eight_composite_message"
}

// tssRoundNineCompositeMessage is a composite of tssRoundNineMessage keyed by
// the message being signed.
type tssRoundNineCompositeMessage struct {
	senderID  group.MemberIndex
	sessionID string

	tssRoundNineMessages map[string]*tssRoundNineMessage
}

// tssRoundNineMessage is a message payload that carries the sender's
// TSS round nine components.
type tssRoundNineMessage struct {
	broadcastPayload []byte
}

// SenderID returns protocol-level identifier of the message sender.
func (trncm *tssRoundNineCompositeMessage) SenderID() group.MemberIndex {
	return trncm.senderID
}

// SessionID returns the session identifier of the message.
func (trncm *tssRoundNineCompositeMessage) SessionID() string {
	return trncm.sessionID
}

// Type returns a string describing a tssRoundNineCompositeMessage type for
// marshaling purposes.
func (trncm *tssRoundNineCompositeMessage) Type() string {
	return messageTypePrefix + "tss_round_nine_composite_message"
}
