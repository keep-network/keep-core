package thresholdsignature

import (
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
		SenderID: ssm.senderID.Bytes(),
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

	ssm.senderID = gjkr.MemberIDFromBytes(pbSignatureShare.SenderID)
	ssm.ShareBytes = pbSignatureShare.Share

	return nil
}
