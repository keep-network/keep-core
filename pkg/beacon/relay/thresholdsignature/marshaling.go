package thresholdsignature

import (
	"github.com/dfinity/go-dfinity-crypto/bls"
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
		Id:    ssm.ID.GetLittleEndian(),
		Share: ssm.ShareBytes,
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

	id := &bls.ID{}
	err = id.SetLittleEndian(pbSignatureShare.Id)
	if err != nil {
		return err
	}

	ssm.ID = id
	ssm.ShareBytes = pbSignatureShare.Share

	return nil
}
