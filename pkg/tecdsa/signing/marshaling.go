package signing

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/keep-network/keep-core/pkg/crypto/ephemeral"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/tecdsa/signing/gen/pb"
)

// Marshal converts this ephemeralPublicKeyMessage to a byte array suitable for
// network communication.
func (epkm *ephemeralPublicKeyMessage) Marshal() ([]byte, error) {
	ephemeralPublicKeys, err := marshalPublicKeyMap(epkm.ephemeralPublicKeys)
	if err != nil {
		return nil, err
	}

	return proto.Marshal(&pb.EphemeralPublicKeyMessage{
		SenderID:            uint32(epkm.senderID),
		EphemeralPublicKeys: ephemeralPublicKeys,
		SessionID:           epkm.sessionID,
	})
}

// Unmarshal converts a byte array produced by Marshal to
// an ephemeralPublicKeyMessage
func (epkm *ephemeralPublicKeyMessage) Unmarshal(bytes []byte) error {
	pbMsg := pb.EphemeralPublicKeyMessage{}
	if err := proto.Unmarshal(bytes, &pbMsg); err != nil {
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
	epkm.sessionID = pbMsg.SessionID

	return nil
}

// Marshal converts this tssRoundOneCompositeMessage to a byte array suitable
// for network communication.
func (trocm *tssRoundOneCompositeMessage) Marshal() ([]byte, error) {
	pbTssRoundOneMessages := make(
		map[string]*pb.TSSRoundOneCompositeMessage_TSSRoundOneMessage,
		len(trocm.tssRoundOneMessages),
	)

	for messageToSign, tssRoundOneMessage := range trocm.tssRoundOneMessages {
		peersPayload := make(map[uint32][]byte, len(tssRoundOneMessage.peersPayload))
		for receiverID, payload := range tssRoundOneMessage.peersPayload {
			peersPayload[uint32(receiverID)] = payload
		}
		pbTssRoundOneMessages[messageToSign] = &pb.TSSRoundOneCompositeMessage_TSSRoundOneMessage{
			PeersPayload:     peersPayload,
			BroadcastPayload: tssRoundOneMessage.broadcastPayload,
		}
	}

	return proto.Marshal(&pb.TSSRoundOneCompositeMessage{
		SenderID:            uint32(trocm.senderID),
		SessionID:           trocm.sessionID,
		TssRoundOneMessages: pbTssRoundOneMessages,
	})
}

// Unmarshal converts a byte array produced by Marshal to
// a tssRoundOneCompositeMessage.
func (trocm *tssRoundOneCompositeMessage) Unmarshal(bytes []byte) error {
	pbMsg := pb.TSSRoundOneCompositeMessage{}
	if err := proto.Unmarshal(bytes, &pbMsg); err != nil {
		return err
	}

	if err := validateMemberIndex(pbMsg.SenderID); err != nil {
		return err
	}

	tssRoundOneMessages := make(
		map[string]*tssRoundOneMessage,
		len(pbMsg.TssRoundOneMessages),
	)

	for messageToSign, pbTssRoundOneMessage := range pbMsg.TssRoundOneMessages {
		peersPayload := make(
			map[group.MemberIndex][]byte,
			len(pbTssRoundOneMessage.PeersPayload),
		)
		for receiverID, payload := range pbTssRoundOneMessage.PeersPayload {
			if err := validateMemberIndex(receiverID); err != nil {
				return err
			}
			peersPayload[group.MemberIndex(receiverID)] = payload
		}
		tssRoundOneMessages[messageToSign] = &tssRoundOneMessage{
			peersPayload:     peersPayload,
			broadcastPayload: pbTssRoundOneMessage.BroadcastPayload,
		}
	}

	trocm.senderID = group.MemberIndex(pbMsg.SenderID)
	trocm.sessionID = pbMsg.SessionID
	trocm.tssRoundOneMessages = tssRoundOneMessages

	return nil
}

// Marshal converts this tssRoundTwoCompositeMessage to a byte array suitable
// for network communication.
func (trtcm *tssRoundTwoCompositeMessage) Marshal() ([]byte, error) {
	pbTssRoundTwoMessages := make(
		map[string]*pb.TSSRoundTwoCompositeMessage_TSSRoundTwoMessage,
		len(trtcm.tssRoundTwoMessages),
	)

	for messageToSign, tssRoundTwoMessage := range trtcm.tssRoundTwoMessages {
		peersPayload := make(map[uint32][]byte, len(tssRoundTwoMessage.peersPayload))
		for receiverID, payload := range tssRoundTwoMessage.peersPayload {
			peersPayload[uint32(receiverID)] = payload
		}
		pbTssRoundTwoMessages[messageToSign] = &pb.TSSRoundTwoCompositeMessage_TSSRoundTwoMessage{
			PeersPayload: peersPayload,
		}
	}

	return proto.Marshal(&pb.TSSRoundTwoCompositeMessage{
		SenderID:            uint32(trtcm.senderID),
		SessionID:           trtcm.sessionID,
		TssRoundTwoMessages: pbTssRoundTwoMessages,
	})
}

// Unmarshal converts a byte array produced by Marshal to
// a tssRoundTwoCompositeMessage.
func (trtcm *tssRoundTwoCompositeMessage) Unmarshal(bytes []byte) error {
	pbMsg := pb.TSSRoundTwoCompositeMessage{}
	if err := proto.Unmarshal(bytes, &pbMsg); err != nil {
		return err
	}

	if err := validateMemberIndex(pbMsg.SenderID); err != nil {
		return err
	}

	tssRoundTwoMessages := make(
		map[string]*tssRoundTwoMessage,
		len(pbMsg.TssRoundTwoMessages),
	)

	for messageToSign, pbTssRoundTwoMessage := range pbMsg.TssRoundTwoMessages {
		peersPayload := make(
			map[group.MemberIndex][]byte,
			len(pbTssRoundTwoMessage.PeersPayload),
		)
		for receiverID, payload := range pbTssRoundTwoMessage.PeersPayload {
			if err := validateMemberIndex(receiverID); err != nil {
				return err
			}
			peersPayload[group.MemberIndex(receiverID)] = payload
		}
		tssRoundTwoMessages[messageToSign] = &tssRoundTwoMessage{
			peersPayload: peersPayload,
		}
	}

	trtcm.senderID = group.MemberIndex(pbMsg.SenderID)
	trtcm.sessionID = pbMsg.SessionID
	trtcm.tssRoundTwoMessages = tssRoundTwoMessages

	return nil
}

// Marshal converts this tssRoundThreeCompositeMessage to a byte array suitable
// for network communication.
func (trtcm *tssRoundThreeCompositeMessage) Marshal() ([]byte, error) {
	pbTssRoundThreeMessages := make(
		map[string]*pb.TSSRoundThreeCompositeMessage_TSSRoundThreeMessage,
		len(trtcm.tssRoundThreeMessages),
	)

	for messageToSign, tssRoundThreeMessage := range trtcm.tssRoundThreeMessages {
		pbTssRoundThreeMessages[messageToSign] = &pb.TSSRoundThreeCompositeMessage_TSSRoundThreeMessage{
			BroadcastPayload: tssRoundThreeMessage.broadcastPayload,
		}
	}

	return proto.Marshal(&pb.TSSRoundThreeCompositeMessage{
		SenderID:              uint32(trtcm.senderID),
		SessionID:             trtcm.sessionID,
		TssRoundThreeMessages: pbTssRoundThreeMessages,
	})
}

// Unmarshal converts a byte array produced by Marshal to
// a tssRoundThreeCompositeMessage.
func (trtcm *tssRoundThreeCompositeMessage) Unmarshal(bytes []byte) error {
	pbMsg := pb.TSSRoundThreeCompositeMessage{}
	if err := proto.Unmarshal(bytes, &pbMsg); err != nil {
		return err
	}

	if err := validateMemberIndex(pbMsg.SenderID); err != nil {
		return err
	}

	tssRoundThreeMessages := make(
		map[string]*tssRoundThreeMessage,
		len(pbMsg.TssRoundThreeMessages),
	)
	for messageToSign, pbTssRoundThreeMessage := range pbMsg.TssRoundThreeMessages {
		tssRoundThreeMessages[messageToSign] = &tssRoundThreeMessage{
			broadcastPayload: pbTssRoundThreeMessage.BroadcastPayload,
		}
	}

	trtcm.senderID = group.MemberIndex(pbMsg.SenderID)
	trtcm.sessionID = pbMsg.SessionID
	trtcm.tssRoundThreeMessages = tssRoundThreeMessages

	return nil
}

// Marshal converts this tssRoundFourCompositeMessage to a byte array suitable
// for network communication.
func (trfcm *tssRoundFourCompositeMessage) Marshal() ([]byte, error) {
	pbTssRoundFourMessages := make(
		map[string]*pb.TSSRoundFourCompositeMessage_TSSRoundFourMessage,
		len(trfcm.tssRoundFourMessages),
	)

	for messageToSign, tssRoundFourMessage := range trfcm.tssRoundFourMessages {
		pbTssRoundFourMessages[messageToSign] = &pb.TSSRoundFourCompositeMessage_TSSRoundFourMessage{
			BroadcastPayload: tssRoundFourMessage.broadcastPayload,
		}
	}

	return proto.Marshal(&pb.TSSRoundFourCompositeMessage{
		SenderID:             uint32(trfcm.senderID),
		SessionID:            trfcm.sessionID,
		TssRoundFourMessages: pbTssRoundFourMessages,
	})
}

// Unmarshal converts a byte array produced by Marshal to
// a tssRoundFourCompositeMessage.
func (trfcm *tssRoundFourCompositeMessage) Unmarshal(bytes []byte) error {
	pbMsg := pb.TSSRoundFourCompositeMessage{}
	if err := proto.Unmarshal(bytes, &pbMsg); err != nil {
		return err
	}

	if err := validateMemberIndex(pbMsg.SenderID); err != nil {
		return err
	}

	tssRoundFourMessages := make(
		map[string]*tssRoundFourMessage,
		len(pbMsg.TssRoundFourMessages),
	)
	for messageToSign, pbTssRoundFourMessage := range pbMsg.TssRoundFourMessages {
		tssRoundFourMessages[messageToSign] = &tssRoundFourMessage{
			broadcastPayload: pbTssRoundFourMessage.BroadcastPayload,
		}
	}

	trfcm.senderID = group.MemberIndex(pbMsg.SenderID)
	trfcm.sessionID = pbMsg.SessionID
	trfcm.tssRoundFourMessages = tssRoundFourMessages

	return nil
}

// Marshal converts this tssRoundFiveCompositeMessage to a byte array suitable
// for network communication.
func (trfcm *tssRoundFiveCompositeMessage) Marshal() ([]byte, error) {
	pbTssRoundFiveMessages := make(
		map[string]*pb.TSSRoundFiveCompositeMessage_TSSRoundFiveMessage,
		len(trfcm.tssRoundFiveMessages),
	)

	for messageToSign, tssRoundFiveMessage := range trfcm.tssRoundFiveMessages {
		pbTssRoundFiveMessages[messageToSign] = &pb.TSSRoundFiveCompositeMessage_TSSRoundFiveMessage{
			BroadcastPayload: tssRoundFiveMessage.broadcastPayload,
		}
	}

	return proto.Marshal(&pb.TSSRoundFiveCompositeMessage{
		SenderID:             uint32(trfcm.senderID),
		SessionID:            trfcm.sessionID,
		TssRoundFiveMessages: pbTssRoundFiveMessages,
	})
}

// Unmarshal converts a byte array produced by Marshal to
// a tssRoundFiveCompositeMessage.
func (trfcm *tssRoundFiveCompositeMessage) Unmarshal(bytes []byte) error {
	pbMsg := pb.TSSRoundFiveCompositeMessage{}
	if err := proto.Unmarshal(bytes, &pbMsg); err != nil {
		return err
	}

	if err := validateMemberIndex(pbMsg.SenderID); err != nil {
		return err
	}

	tssRoundFiveMessages := make(
		map[string]*tssRoundFiveMessage,
		len(pbMsg.TssRoundFiveMessages),
	)
	for messageToSign, pbTssRoundFiveMessage := range pbMsg.TssRoundFiveMessages {
		tssRoundFiveMessages[messageToSign] = &tssRoundFiveMessage{
			broadcastPayload: pbTssRoundFiveMessage.BroadcastPayload,
		}
	}

	trfcm.senderID = group.MemberIndex(pbMsg.SenderID)
	trfcm.sessionID = pbMsg.SessionID
	trfcm.tssRoundFiveMessages = tssRoundFiveMessages

	return nil
}

// Marshal converts this tssRoundSixCompositeMessage to a byte array suitable
// for network communication.
func (trscm *tssRoundSixCompositeMessage) Marshal() ([]byte, error) {
	pbTssRoundSixMessages := make(
		map[string]*pb.TSSRoundSixCompositeMessage_TSSRoundSixMessage,
		len(trscm.tssRoundSixMessages),
	)

	for messageToSign, tssRoundSixMessage := range trscm.tssRoundSixMessages {
		pbTssRoundSixMessages[messageToSign] = &pb.TSSRoundSixCompositeMessage_TSSRoundSixMessage{
			BroadcastPayload: tssRoundSixMessage.broadcastPayload,
		}
	}

	return proto.Marshal(&pb.TSSRoundSixCompositeMessage{
		SenderID:            uint32(trscm.senderID),
		SessionID:           trscm.sessionID,
		TssRoundSixMessages: pbTssRoundSixMessages,
	})
}

// Unmarshal converts a byte array produced by Marshal to
// a tssRoundSixCompositeMessage.
func (trscm *tssRoundSixCompositeMessage) Unmarshal(bytes []byte) error {
	pbMsg := pb.TSSRoundSixCompositeMessage{}
	if err := proto.Unmarshal(bytes, &pbMsg); err != nil {
		return err
	}

	if err := validateMemberIndex(pbMsg.SenderID); err != nil {
		return err
	}

	tssRoundSixMessages := make(
		map[string]*tssRoundSixMessage,
		len(pbMsg.TssRoundSixMessages),
	)
	for messageToSign, pbTssRoundSixMessage := range pbMsg.TssRoundSixMessages {
		tssRoundSixMessages[messageToSign] = &tssRoundSixMessage{
			broadcastPayload: pbTssRoundSixMessage.BroadcastPayload,
		}
	}

	trscm.senderID = group.MemberIndex(pbMsg.SenderID)
	trscm.sessionID = pbMsg.SessionID
	trscm.tssRoundSixMessages = tssRoundSixMessages

	return nil
}

// Marshal converts this tssRoundSevenCompositeMessage to a byte array suitable
// for network communication.
func (trscm *tssRoundSevenCompositeMessage) Marshal() ([]byte, error) {
	pbTssRoundSevenMessages := make(
		map[string]*pb.TSSRoundSevenCompositeMessage_TSSRoundSevenMessage,
		len(trscm.tssRoundSevenMessages),
	)

	for messageToSign, tssRoundSevenMessage := range trscm.tssRoundSevenMessages {
		pbTssRoundSevenMessages[messageToSign] = &pb.TSSRoundSevenCompositeMessage_TSSRoundSevenMessage{
			BroadcastPayload: tssRoundSevenMessage.broadcastPayload,
		}
	}

	return proto.Marshal(&pb.TSSRoundSevenCompositeMessage{
		SenderID:              uint32(trscm.senderID),
		SessionID:             trscm.sessionID,
		TssRoundSevenMessages: pbTssRoundSevenMessages,
	})
}

// Unmarshal converts a byte array produced by Marshal to
// a tssRoundSevenCompositeMessage.
func (trscm *tssRoundSevenCompositeMessage) Unmarshal(bytes []byte) error {
	pbMsg := pb.TSSRoundSevenCompositeMessage{}
	if err := proto.Unmarshal(bytes, &pbMsg); err != nil {
		return err
	}

	if err := validateMemberIndex(pbMsg.SenderID); err != nil {
		return err
	}

	tssRoundSevenMessages := make(
		map[string]*tssRoundSevenMessage,
		len(pbMsg.TssRoundSevenMessages),
	)
	for messageToSign, pbTssRoundSevenMessage := range pbMsg.TssRoundSevenMessages {
		tssRoundSevenMessages[messageToSign] = &tssRoundSevenMessage{
			broadcastPayload: pbTssRoundSevenMessage.BroadcastPayload,
		}
	}

	trscm.senderID = group.MemberIndex(pbMsg.SenderID)
	trscm.sessionID = pbMsg.SessionID
	trscm.tssRoundSevenMessages = tssRoundSevenMessages

	return nil
}

// Marshal converts this tssRoundEightCompositeMessage to a byte array suitable
// for network communication.
func (trecm *tssRoundEightCompositeMessage) Marshal() ([]byte, error) {
	pbTssRoundEightMessages := make(
		map[string]*pb.TSSRoundEightCompositeMessage_TSSRoundEightMessage,
		len(trecm.tssRoundEightMessages),
	)
	for messageToSign, tssRoundEightMessage := range trecm.tssRoundEightMessages {
		pbTssRoundEightMessages[messageToSign] = &pb.TSSRoundEightCompositeMessage_TSSRoundEightMessage{
			BroadcastPayload: tssRoundEightMessage.broadcastPayload,
		}
	}

	return proto.Marshal(&pb.TSSRoundEightCompositeMessage{
		SenderID:              uint32(trecm.senderID),
		SessionID:             trecm.sessionID,
		TssRoundEightMessages: pbTssRoundEightMessages,
	})
}

// Unmarshal converts a byte array produced by Marshal to
// a tssRoundEightCompositeMessage.
func (trecm *tssRoundEightCompositeMessage) Unmarshal(bytes []byte) error {
	pbMsg := pb.TSSRoundEightCompositeMessage{}
	if err := proto.Unmarshal(bytes, &pbMsg); err != nil {
		return err
	}

	if err := validateMemberIndex(pbMsg.SenderID); err != nil {
		return err
	}

	tssRoundEightMessages := make(
		map[string]*tssRoundEightMessage,
		len(pbMsg.TssRoundEightMessages),
	)
	for messageToSign, pbTssRoundEightMessage := range pbMsg.TssRoundEightMessages {
		tssRoundEightMessages[messageToSign] = &tssRoundEightMessage{
			broadcastPayload: pbTssRoundEightMessage.BroadcastPayload,
		}
	}

	trecm.senderID = group.MemberIndex(pbMsg.SenderID)
	trecm.sessionID = pbMsg.SessionID
	trecm.tssRoundEightMessages = tssRoundEightMessages

	return nil
}

// Marshal converts this tssRoundNineCompositeMessage to a byte array suitable
// for network communication.
func (trncm *tssRoundNineCompositeMessage) Marshal() ([]byte, error) {
	pbTssRoundNineMessages := make(
		map[string]*pb.TSSRoundNineCompositeMessage_TSSRoundNineMessage,
		len(trncm.tssRoundNineMessages),
	)
	for messageToSign, tssRoundNineMessage := range trncm.tssRoundNineMessages {
		pbTssRoundNineMessages[messageToSign] = &pb.TSSRoundNineCompositeMessage_TSSRoundNineMessage{
			BroadcastPayload: tssRoundNineMessage.broadcastPayload,
		}
	}

	return proto.Marshal(&pb.TSSRoundNineCompositeMessage{
		SenderID:             uint32(trncm.senderID),
		SessionID:            trncm.sessionID,
		TssRoundNineMessages: pbTssRoundNineMessages,
	})
}

// Unmarshal converts a byte array produced by Marshal to a tssRoundNineMessage.
func (trncm *tssRoundNineCompositeMessage) Unmarshal(bytes []byte) error {
	pbMsg := pb.TSSRoundNineCompositeMessage{}
	if err := proto.Unmarshal(bytes, &pbMsg); err != nil {
		return err
	}

	if err := validateMemberIndex(pbMsg.SenderID); err != nil {
		return err
	}

	tssRoundNineMessages := make(
		map[string]*tssRoundNineMessage,
		len(pbMsg.TssRoundNineMessages),
	)
	for messageToSign, pbTssRoundNineMessage := range pbMsg.TssRoundNineMessages {
		tssRoundNineMessages[messageToSign] = &tssRoundNineMessage{
			broadcastPayload: pbTssRoundNineMessage.BroadcastPayload,
		}
	}

	trncm.senderID = group.MemberIndex(pbMsg.SenderID)
	trncm.sessionID = pbMsg.SessionID
	trncm.tssRoundNineMessages = tssRoundNineMessages

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
