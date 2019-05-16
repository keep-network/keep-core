package dkg

import (
	"fmt"
	"math/big"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
	"github.com/keep-network/keep-core/pkg/beacon/relay/registry/gen/pb"
)

// Marshal converts ThresholdSigner to byte array.
func (ts *ThresholdSigner) Marshal() ([]byte, error) {
	return (&pb.ThresholdSigner{
		MemberIndex:          uint32(ts.memberIndex),
		GroupPublicKey:       ts.groupPublicKey.Marshal(),
		GroupPrivateKeyShare: ts.groupPrivateKeyShare.String(),
	}).Marshal()
}

// Unmarshal converts a byte array back to ThresholdSigner.
func (ts *ThresholdSigner) Unmarshal(bytes []byte) error {
	pbThresholdSigner := pb.ThresholdSigner{}
	if err := pbThresholdSigner.Unmarshal(bytes); err != nil {
		return err
	}

	groupPublicKey := new(bn256.G2)
	_, err := groupPublicKey.Unmarshal(pbThresholdSigner.GroupPublicKey)
	if err != nil {
		return err
	}

	privateKeyShare := new(big.Int)
	privateKeyShare, ok := privateKeyShare.SetString(pbThresholdSigner.GroupPrivateKeyShare, 10)
	if !ok {
		return fmt.Errorf("Error occured while converting a private key share to string")
	}

	ts.memberIndex = group.MemberIndex(pbThresholdSigner.MemberIndex)
	ts.groupPublicKey = groupPublicKey
	ts.groupPrivateKeyShare = privateKeyShare

	return nil
}
