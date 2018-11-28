package gjkr

import (
	"encoding/binary"

	"github.com/keep-network/keep-core/pkg/beacon/relay/gjkr/gen/pb"
	"github.com/keep-network/keep-core/pkg/net/ephemeral"
)

func (epkm *EphemeralPublicKeyMessage) Type() string {
	return "dkg/ephemeral_public_key"
}

func (epkm *EphemeralPublicKeyMessage) Marshal() ([]byte, error) {
	return (&pb.EphemeralPublicKey{
		SenderID:           memberIDToBytes(epkm.senderID),
		ReceiverID:         memberIDToBytes(epkm.receiverID),
		EphemeralPublicKey: epkm.ephemeralPublicKey.Marshal(),
	}).Marshal()
}

func (epkm *EphemeralPublicKeyMessage) Unmarshal(bytes []byte) error {
	pbMsg := pb.EphemeralPublicKey{}
	err := pbMsg.Unmarshal(bytes)
	if err != nil {
		return err
	}

	ephemeralPublicKey, err := ephemeral.UnmarshalPublicKey(
		pbMsg.EphemeralPublicKey,
	)
	if err != nil {
		return err
	}

	epkm.senderID = bytesToMemberID(pbMsg.SenderID)
	epkm.receiverID = bytesToMemberID(pbMsg.ReceiverID)
	epkm.ephemeralPublicKey = ephemeralPublicKey

	return nil
}

func memberIDToBytes(memberID MemberID) []byte {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, uint32(memberID))
	return bytes
}

func bytesToMemberID(bytes []byte) MemberID {
	return MemberID(binary.LittleEndian.Uint32(bytes))
}
