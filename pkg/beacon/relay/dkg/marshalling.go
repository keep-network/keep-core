package dkg

import (
	"fmt"
	"math/big"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
	"github.com/keep-network/keep-core/pkg/beacon/relay/registry/gen/pb"
)

// Marshal converts pb.ThresholdSigner to a dkg.ThresholdSigner.
func (ts *ThresholdSigner) Marshal() ([]byte, error) {
	return (&pb.ThresholdSigner{
		MemberIndex:          uint32(ts.memberIndex),
		GroupPublicKey:       ts.groupPublicKey.Marshal(),
		GroupPrivateKeyShare: ts.groupPrivateKeyShare.String(),
	}).Marshal()
}

// Unmarshal converts a byte array back to ThresholdSigner
func (ts *ThresholdSigner) Unmarshal(bytes []byte) error {
	protoBuffThresholdSigner := pb.ThresholdSigner{}
	if err := protoBuffThresholdSigner.Unmarshal(bytes); err != nil {
		return err
	}

	groupPublicKeyBn256 := new(bn256.G2)
	_, err := groupPublicKeyBn256.Unmarshal(protoBuffThresholdSigner.GroupPublicKey)
	if err != nil {
		return err
	}

	privateKeyShare := new(big.Int)
	privateKeyShare, ok := privateKeyShare.SetString(protoBuffThresholdSigner.GroupPrivateKeyShare, 10)
	if !ok {
		return fmt.Errorf("Error occured while converting a private key share to string")
	}

	ts.memberIndex = group.MemberIndex(protoBuffThresholdSigner.MemberIndex)
	ts.groupPublicKey = groupPublicKeyBn256
	ts.groupPrivateKeyShare = privateKeyShare

	return nil
}
