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

// Marshal converts the wallet to a byte array.
func (w *wallet) Marshal() ([]byte, error) {
	publicKey := elliptic.Marshal(
		w.publicKey.Curve,
		w.publicKey.X,
		w.publicKey.Y,
	)

	signingGroupOperators := make([]string, len(w.signingGroupOperators))
	for i := range signingGroupOperators {
		signingGroupOperators[i] = w.signingGroupOperators[i].String()
	}

	return (&pb.Wallet{
		PublicKey:             publicKey,
		SigningGroupOperators: signingGroupOperators,
	}).Marshal()
}

// Unmarshal converts a byte array back to the wallet.
func (w *wallet) Unmarshal(bytes []byte) error {
	pbWallet := &pb.Wallet{}
	if err := pbWallet.Unmarshal(bytes); err != nil {
		return fmt.Errorf("cannot unmarshal wallet: [%w]", err)
	}

	publicKeyX, publicKeyY := elliptic.Unmarshal(tecdsa.Curve, pbWallet.PublicKey)
	publicKey := &ecdsa.PublicKey{
		Curve: tecdsa.Curve,
		X:     publicKeyX,
		Y:     publicKeyY,
	}

	signingGroupOperators := make(
		[]chain.Address,
		len(pbWallet.SigningGroupOperators),
	)
	for i := range signingGroupOperators {
		signingGroupOperators[i] = chain.Address(pbWallet.SigningGroupOperators[i])
	}

	w.publicKey = publicKey
	w.signingGroupOperators = signingGroupOperators

	return nil
}

// Marshal converts the signer to a byte array.
func (s *signer) Marshal() ([]byte, error) {
	wallet, err := s.wallet.Marshal()
	if err != nil {
		return nil, fmt.Errorf("cannot marshal wallet: [%w]", err)
	}

	privateKeyShare, err := s.privateKeyShare.Marshal()
	if err != nil {
		return nil, fmt.Errorf("cannot marshal private key share: [%w]", err)
	}

	return (&pb.Signer{
		Wallet:                  wallet,
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

	wallet := &wallet{}
	if err := wallet.Unmarshal(pbSigner.Wallet); err != nil {
		return fmt.Errorf("cannot unmarshal wallet: [%w]", err)
	}

	privateKeyShare := &tecdsa.PrivateKeyShare{}
	if err := privateKeyShare.Unmarshal(pbSigner.PrivateKeyShare); err != nil {
		return fmt.Errorf("cannot unmarshal private key share: [%w]", err)
	}

	s.wallet = wallet
	s.signingGroupMemberIndex = group.MemberIndex(pbSigner.SigningGroupMemberIndex)
	s.privateKeyShare = privateKeyShare

	return nil
}
