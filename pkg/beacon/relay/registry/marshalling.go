package registry

import (
	"math/big"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/keep-network/keep-core/pkg/beacon/relay/dkg"
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
	"github.com/keep-network/keep-core/pkg/beacon/relay/registry/gen/pb"
)

// Marshal converts Membership to a byte array suitable for network communication.
func (m *Membership) Marshal() ([]byte, error) {
	return (&pb.MembershipMessage{
		MemberIndex:          uint32(m.Signer.MemberIndex),
		GroupPublicKey:       m.Signer.GroupPublicKey.Marshal(),
		GroupPrivateKeyShare: m.Signer.GroupPrivateKeyShare.String(),
		Channel:              "test channel",
	}).Marshal()

}

// Unmarshal converts a byte array produced by Marshal to Membership
func (m *Membership) Unmarshal(bytes []byte) error {
	protoBuffMembershipMessage := pb.MembershipMessage{}
	if err := protoBuffMembershipMessage.Unmarshal(bytes); err != nil {
		return err
	}

	groupPublicKeyBn256 := new(bn256.G2)
	_, err := groupPublicKeyBn256.Unmarshal(protoBuffMembershipMessage.GroupPublicKey)
	if err != nil {
		return err
	}

	privateKeyShare := new(big.Int)
	privateKeyShare, ok := privateKeyShare.SetString(protoBuffMembershipMessage.GroupPrivateKeyShare, 10)
	if !ok {
		return nil
	}

	thresholdSigner := dkg.NewThresholdSigner(
		group.MemberIndex(protoBuffMembershipMessage.MemberIndex),
		groupPublicKeyBn256,
		privateKeyShare,
	)

	m.Signer = thresholdSigner
	// m.Channel = protoBuffMembershipMessage.Channel; //TODO: will be implemented later

	return nil
}
