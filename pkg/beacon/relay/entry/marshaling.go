package entry

import (
	"github.com/keep-network/keep-core/pkg/beacon/relay/entry/gen/pb"
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
)

// Type returns a string describing a SignatureShareMessage's type.
func (*SignatureShareMessage) Type() string {
	return "relay/signature/share"
}

// Marshal converts this SignatureShareMessage to a byte array suitable for
// network communication.
func (ssm *SignatureShareMessage) Marshal() ([]byte, error) {
	pbSignatureShare := pb.SignatureShare{
		SenderID: uint32(ssm.senderID),
		Share:    ssm.shareBytes,
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

	ssm.senderID = group.MemberIndex(pbSignatureShare.SenderID)
	ssm.shareBytes = pbSignatureShare.Share

	return nil
}
