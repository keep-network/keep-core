package thresholdsignature

import (
	"encoding/binary"

	"github.com/keep-network/keep-core/pkg/beacon/relay/gjkr"
	"github.com/keep-network/keep-core/pkg/beacon/relay/thresholdsignature/gen/pb"
)

// Type returns a string describing a SignatureShareMessage's type.
func (*SignatureShareMessage) Type() string {
	return "relay/signature/share"
}

// Marshal converts this JustificationsMessage to a byte array suitable for
// network communication.
func (ssm *SignatureShareMessage) Marshal() ([]byte, error) {
	pbSignatureShare := pb.SignatureShare{
		SenderID: memberIDToBytes(ssm.senderID),
		Share:    ssm.ShareBytes,
	}

	return pbSignatureShare.Marshal()
}

// Unmarshal converts a byte array produced by Marshal to a
// SignatureShareMessage.
func (ssm *SignatureShareMessage) Unmarshal(bytes []byte) error {
	pbSignatureShare := pb.SignatureShare{}
	err := pbSignatureShare.Unmarshal(bytes)
	if err != nil {
		return err
	}

	ssm.senderID = bytesToMemberID(pbSignatureShare.SenderID)
	ssm.ShareBytes = pbSignatureShare.Share

	return nil
}

// TODO: CODE DUPLICATION! MOVE THOSE FUNCTIONS TO MEMBER ID

func memberIDToBytes(memberID gjkr.MemberID) []byte {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, uint32(memberID))
	return bytes
}

func bytesToMemberID(bytes []byte) gjkr.MemberID {
	return gjkr.MemberID(binary.LittleEndian.Uint32(bytes))
}
