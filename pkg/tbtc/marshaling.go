package tbtc

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"math/big"

	"google.golang.org/protobuf/proto"

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
	walletPublicKey, err := marshalPublicKey(s.wallet.publicKey)
	if err != nil {
		return nil, err
	}

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

	return proto.Marshal(&pb.Signer{
		Wallet:                  pbWallet,
		SigningGroupMemberIndex: uint32(s.signingGroupMemberIndex),
		PrivateKeyShare:         privateKeyShare,
	})
}

// Unmarshal converts a byte array back to the signer.
func (s *signer) Unmarshal(bytes []byte) error {
	pbSigner := pb.Signer{}
	if err := proto.Unmarshal(bytes, &pbSigner); err != nil {
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

// Marshal converts the signingSyncMessage to a byte array.
func (ssm *signingSyncMessage) Marshal() ([]byte, error) {
	signatureBytes, err := ssm.signature.Marshal()
	if err != nil {
		return nil, err
	}

	return proto.Marshal(&pb.SigningSyncMessage{
		SenderID:      uint32(ssm.senderID),
		Message:       ssm.message.Bytes(),
		AttemptNumber: uint64(ssm.attemptNumber),
		Signature:     signatureBytes,
		EndBlock:      ssm.endBlock,
	})
}

// Unmarshal converts a byte array back to the signingSyncMessage.
func (ssm *signingSyncMessage) Unmarshal(bytes []byte) error {
	pbMsg := pb.SigningSyncMessage{}
	if err := proto.Unmarshal(bytes, &pbMsg); err != nil {
		return fmt.Errorf("failed to unmarshal SigningSyncMessage: [%v]", err)
	}

	if err := validateMemberIndex(pbMsg.SenderID); err != nil {
		return err
	}

	signature := &tecdsa.Signature{}
	if err := signature.Unmarshal(pbMsg.Signature); err != nil {
		return fmt.Errorf("cannot unmarshal signature: [%v]", err)
	}

	ssm.senderID = group.MemberIndex(pbMsg.SenderID)
	ssm.message = new(big.Int).SetBytes(pbMsg.Message)
	ssm.attemptNumber = uint(pbMsg.AttemptNumber)
	ssm.signature = signature
	ssm.endBlock = pbMsg.EndBlock

	return nil
}

// marshalPublicKey converts an ECDSA public key to a byte
// array (uncompressed).
func marshalPublicKey(publicKey *ecdsa.PublicKey) ([]byte, error) {
	if publicKey.Curve.Params().Name != tecdsa.Curve.Params().Name {
		return nil, errIncompatiblePublicKey
	}

	return elliptic.Marshal(
		publicKey.Curve,
		publicKey.X,
		publicKey.Y,
	), nil
}

// unmarshalPublicKey converts a byte array (uncompressed) to an ECDSA
// public key.
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

func validateMemberIndex(protoIndex uint32) error {
	// Protobuf does not have uint8 type, so we are using uint32. When
	// unmarshalling message, we need to make sure we do not overflow.
	if protoIndex > group.MaxMemberIndex {
		return fmt.Errorf("invalid member index value: [%v]", protoIndex)
	}
	return nil
}
