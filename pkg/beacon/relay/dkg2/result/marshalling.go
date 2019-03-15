package result

import (
	"github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/dkg2/result/gen/pb"
	"github.com/keep-network/keep-core/pkg/beacon/relay/gjkr"
)

// Marshal converts this DKGResultHashSignatureMessage to a byte array suitable
// for network communication.
func (d *DKGResultHashSignatureMessage) Marshal() ([]byte, error) {
	return (&pb.DKGResultHashSignature{
		SenderIndex: uint32(d.senderIndex),
		ResultHash:  d.resultHash[:],
		Signature:   d.signature,
		// PublicKey:   , // TODO: Add public key marshalling when static.PublicKey is ready
	}).Marshal()
}

// Unmarshal converts a byte array produced by Marshal to a
// DKGResultHashSignatureMessage.
func (d *DKGResultHashSignatureMessage) Unmarshal(bytes []byte) error {
	pbMsg := pb.DKGResultHashSignature{}
	if err := pbMsg.Unmarshal(bytes); err != nil {
		return err
	}
	d.senderIndex = gjkr.MemberID(pbMsg.SenderIndex)

	resultHash, err := chain.DKGResultHashFromBytes(pbMsg.ResultHash)
	if err != nil {
		return err
	}
	d.resultHash = resultHash

	d.signature = pbMsg.Signature

	// d.publicKey =    // TODO: Add public key unmarshalling when static.PublicKey is ready

	return nil
}
