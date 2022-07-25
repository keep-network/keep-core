package dkg

import (
	"fmt"
	"github.com/keep-network/keep-core/pkg/crypto/ephemeral"
	"github.com/keep-network/keep-core/pkg/ecdsa/dkg/gen/pb"
	"github.com/keep-network/keep-core/pkg/protocol/group"
)

const messageTypePrefix = "ecdsa_dkg/"

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

// Marshal converts this ephemeralPublicKeyMessage to a byte array suitable for
// network communication.
func (epkm *ephemeralPublicKeyMessage) Marshal() ([]byte, error) {
	ephemeralPublicKeys, err := marshalPublicKeyMap(epkm.ephemeralPublicKeys)
	if err != nil {
		return nil, err
	}

	return (&pb.EphemeralPublicKeyMessage{
		SenderID:            uint32(epkm.senderID),
		EphemeralPublicKeys: ephemeralPublicKeys,
	}).Marshal()
}

// Unmarshal converts a byte array produced by Marshal to
// an ephemeralPublicKeyMessage
func (epkm *ephemeralPublicKeyMessage) Unmarshal(bytes []byte) error {
	pbMsg := pb.EphemeralPublicKeyMessage{}
	if err := pbMsg.Unmarshal(bytes); err != nil {
		return err
	}

	if err := validateMemberIndex(pbMsg.SenderID); err != nil {
		return err
	}
	epkm.senderID = group.MemberIndex(pbMsg.SenderID)

	ephemeralPublicKeys, err := unmarshalPublicKeyMap(pbMsg.EphemeralPublicKeys)
	if err != nil {
		return err
	}

	epkm.ephemeralPublicKeys = ephemeralPublicKeys

	return nil
}

// tssRoundOneMessage is a message payload that carries the sender's TSS
// commitments and Paillier public keys generated for all other group members.
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

// Marshal converts this tssRoundOneMessage to a byte array suitable for
// network communication.
func (trom *tssRoundOneMessage) Marshal() ([]byte, error) {
	return (&pb.TSSRoundOneMessage{
		SenderID:  uint32(trom.senderID),
		Payload:   trom.payload,
		SessionID: trom.sessionID,
	}).Marshal()
}

// Unmarshal converts a byte array produced by Marshal to an tssRoundOneMessage.
func (trom *tssRoundOneMessage) Unmarshal(bytes []byte) error {
	pbMsg := pb.TSSRoundOneMessage{}
	if err := pbMsg.Unmarshal(bytes); err != nil {
		return err
	}

	if err := validateMemberIndex(pbMsg.SenderID); err != nil {
		return err
	}

	trom.senderID = group.MemberIndex(pbMsg.SenderID)
	trom.payload = pbMsg.Payload
	trom.sessionID = pbMsg.SessionID

	return nil
}

func validateMemberIndex(protoIndex uint32) error {
	// Protobuf does not have uint8 type, so we are using uint32. When
	// unmarshalling message, we need to make sure we do not overflow.
	if protoIndex > group.MaxMemberIndex {
		return fmt.Errorf("invalid member index value: [%v]", protoIndex)
	}
	return nil
}

func marshalPublicKeyMap(
	publicKeys map[group.MemberIndex]*ephemeral.PublicKey,
) (map[uint32][]byte, error) {
	marshalled := make(map[uint32][]byte, len(publicKeys))
	for id, publicKey := range publicKeys {
		if publicKey == nil {
			return nil, fmt.Errorf("nil public key for member [%v]", id)
		}

		marshalled[uint32(id)] = publicKey.Marshal()
	}
	return marshalled, nil
}

func unmarshalPublicKeyMap(
	publicKeys map[uint32][]byte,
) (map[group.MemberIndex]*ephemeral.PublicKey, error) {
	var unmarshalled = make(map[group.MemberIndex]*ephemeral.PublicKey, len(publicKeys))
	for memberID, publicKeyBytes := range publicKeys {
		if err := validateMemberIndex(memberID); err != nil {
			return nil, err
		}

		publicKey, err := ephemeral.UnmarshalPublicKey(publicKeyBytes)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal public key [%v]", err)
		}

		unmarshalled[group.MemberIndex(memberID)] = publicKey

	}

	return unmarshalled, nil
}
