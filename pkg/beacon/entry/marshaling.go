package entry

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/keep-network/keep-core/pkg/beacon/entry/gen/pb"
	"github.com/keep-network/keep-core/pkg/protocol/group"
)

// MemberIndex is represented as uint8 in gjkr. Protobuf does not have uint8
// type so we are using uint32. When unmarshalling message, we need to make
// sure we do not overflow.
const maxMemberIndex = 255

func validateMemberIndex(protoIndex uint32) error {
	if protoIndex > maxMemberIndex {
		return fmt.Errorf("invalid member index value: [%v]", protoIndex)
	}
	return nil
}

// Type returns a string describing a SignatureShareMessage's type.
func (*SignatureShareMessage) Type() string {
	return "relay/signature/share"
}

// Marshal converts this SignatureShareMessage to a byte array suitable for
// network communication.
func (ssm *SignatureShareMessage) Marshal() ([]byte, error) {
	pbSignatureShare := pb.SignatureShare{
		SenderID:  uint32(ssm.senderID),
		Share:     ssm.shareBytes,
		SessionID: ssm.sessionID,
	}

	return proto.Marshal(&pbSignatureShare)
}

// Unmarshal converts a byte array produced by Marshal to a
// SignatureShareMessage.
func (ssm *SignatureShareMessage) Unmarshal(bytes []byte) error {
	pbSignatureShare := pb.SignatureShare{}
	err := proto.Unmarshal(bytes, &pbSignatureShare)
	if err != nil {
		return err
	}

	if err := validateMemberIndex(pbSignatureShare.SenderID); err != nil {
		return err
	}
	ssm.senderID = group.MemberIndex(pbSignatureShare.SenderID)
	ssm.shareBytes = pbSignatureShare.Share
	ssm.sessionID = pbSignatureShare.SessionID

	return nil
}
