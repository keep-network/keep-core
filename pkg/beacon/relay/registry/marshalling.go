package registry

import (
	"fmt"

	"github.com/keep-network/keep-core/pkg/beacon/relay/dkg"
	"github.com/keep-network/keep-core/pkg/beacon/relay/registry/gen/pb"
)

// Marshal converts Membership to a byte array.
func (m *Membership) Marshal() ([]byte, error) {
	signer, err := m.Signer.Marshal()
	if err != nil {
		return nil, err
	}

	return (&pb.Membership{
		Signer: signer,
		// Channel: "test channel",
	}).Marshal()
}

// Unmarshal converts a byte array produced by Marshal to Membership.
func (m *Membership) Unmarshal(bytes []byte) error {
	pbMembership := pb.Membership{}
	if err := pbMembership.Unmarshal(bytes); err != nil {
		return err
	}

	signer := &dkg.ThresholdSigner{}

	err := signer.Unmarshal(pbMembership.Signer)
	if err != nil {
		return fmt.Errorf("Unexpected error occured [%v]", err)
	}

	m.Signer = signer
	// m.Channel = protoBuffMembership.Channel; //TODO: will be implemented later

	return nil
}
