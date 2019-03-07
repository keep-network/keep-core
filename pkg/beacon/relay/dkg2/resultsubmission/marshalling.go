package resultsubmission

import (
	"github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/resultsubmission/gen/pb"
)

// Marshal converts this DKGResultHashSignatureMessage to a byte array suitable
// for network communication.
func (d *DKGResultHashSignatureMessage) Marshal() ([]byte, error) {
	return (&pb.DKGResultHashSignature{
		SenderID:   d.senderIndex,
		ResultHash: d.resultHash[:],
		Signature:  d.signature,
	}).Marshal()
}

// Unmarshal converts a byte array produced by Marshal to a
// DKGResultHashSignatureMessage.
func (d *DKGResultHashSignatureMessage) Unmarshal(bytes []byte) error {
	pbMsg := pb.DKGResultHashSignature{}
	if err := pbMsg.Unmarshal(bytes); err != nil {
		return err
	}
	d.senderIndex = pbMsg.SenderID

	resultHash, err := chain.DKGResultHashFromBytes(pbMsg.ResultHash)
	if err != nil {
		return err
	}
	d.resultHash = resultHash

	d.signature = pbMsg.Signature

	return nil
}
