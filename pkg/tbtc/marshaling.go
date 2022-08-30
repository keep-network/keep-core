package tbtc

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/tbtc/gen/pb"
	"github.com/keep-network/keep-core/pkg/tecdsa"
)

var errIncompatiblePublicKey = fmt.Errorf(
	"public key is not tECDSA compatible and will cause unmarshaling error",
)

// Marshal converts the signer to a byte array.
func (s *signer) Marshal() ([]byte, error) {
	if s.wallet.publicKey.Curve.Params().Name != tecdsa.Curve.Params().Name {
		return nil, errIncompatiblePublicKey
	}

	walletPublicKey := elliptic.Marshal(
		s.wallet.publicKey.Curve,
		s.wallet.publicKey.X,
		s.wallet.publicKey.Y,
	)

	walletSigningGroupOperators := make(
		[]string,
		len(s.wallet.signingGroupOperators),
	)
	for i := range walletSigningGroupOperators {
		walletSigningGroupOperators[i] =
			s.wallet.signingGroupOperators[i].String()
	}

	pbWallet := &pb.Wallet{
		PublicKey:             walletPublicKey,
		SigningGroupOperators: walletSigningGroupOperators,
	}

	privateKeyShare, err := s.privateKeyShare.Marshal()
	if err != nil {
		return nil, fmt.Errorf("cannot marshal private key share: [%w]", err)
	}

	return (&pb.Signer{
		Wallet:                  pbWallet,
		SigningGroupMemberIndex: uint32(s.signingGroupMemberIndex),
		PrivateKeyShare:         privateKeyShare,
	}).Marshal()
}

// Unmarshal converts a byte array back to the signer.
func (s *signer) Unmarshal(bytes []byte) error {
	pbSigner := &pb.Signer{}
	if err := pbSigner.Unmarshal(bytes); err != nil {
		return fmt.Errorf("cannot unmarshal signer: [%w]", err)
	}

	walletPublicKey := unmarshalPublicKey(pbSigner.Wallet.PublicKey)

	walletSigningGroupOperators := make(
		[]chain.Address,
		len(pbSigner.Wallet.SigningGroupOperators),
	)
	for i := range walletSigningGroupOperators {
		walletSigningGroupOperators[i] =
			chain.Address(pbSigner.Wallet.SigningGroupOperators[i])
	}

	privateKeyShare := &tecdsa.PrivateKeyShare{}
	if err := privateKeyShare.Unmarshal(pbSigner.PrivateKeyShare); err != nil {
		return fmt.Errorf("cannot unmarshal private key share: [%w]", err)
	}

	s.wallet = wallet{
		publicKey:             walletPublicKey,
		signingGroupOperators: walletSigningGroupOperators,
	}
	s.signingGroupMemberIndex = group.MemberIndex(pbSigner.SigningGroupMemberIndex)
	s.privateKeyShare = privateKeyShare

	return nil
}

// unmarshalPublicKey converts a byte array to an ECDSA public key.
func unmarshalPublicKey(bytes []byte) *ecdsa.PublicKey {
	x, y := elliptic.Unmarshal(
		tecdsa.Curve,
		bytes,
	)

	return &ecdsa.PublicKey{
		Curve: tecdsa.Curve,
		X:     x,
		Y:     y,
	}
}