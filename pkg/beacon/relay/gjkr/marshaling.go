package gjkr

import (
	"encoding/binary"
	"math/big"

	"github.com/keep-network/keep-core/pkg/beacon/relay/gjkr/gen/pb"
	"github.com/keep-network/keep-core/pkg/net/ephemeral"
)

// Type returns a string describing an EphemeralPublicKeyMessage type for
// marshaling purposes.
func (epkm *EphemeralPublicKeyMessage) Type() string {
	return "gjkr/ephemeral_public_key"
}

// Marshal converts this EphemeralPublicKeyMessage to a byte array suitable for
// network communication.
func (epkm *EphemeralPublicKeyMessage) Marshal() ([]byte, error) {
	return (&pb.EphemeralPublicKey{
		SenderID:           memberIDToBytes(epkm.senderID),
		ReceiverID:         memberIDToBytes(epkm.receiverID),
		EphemeralPublicKey: epkm.ephemeralPublicKey.Marshal(),
	}).Marshal()
}

// Unmarshal converts a byte array produced by Marshal to
// an EphemeralPublicKeyMessage
func (epkm *EphemeralPublicKeyMessage) Unmarshal(bytes []byte) error {
	pbMsg := pb.EphemeralPublicKey{}
	if err := pbMsg.Unmarshal(bytes); err != nil {
		return err
	}

	ephemeralPublicKey, err := ephemeral.UnmarshalPublicKey(
		pbMsg.EphemeralPublicKey,
	)
	if err != nil {
		return err
	}

	epkm.senderID = bytesToMemberID(pbMsg.SenderID)
	epkm.receiverID = bytesToMemberID(pbMsg.ReceiverID)
	epkm.ephemeralPublicKey = ephemeralPublicKey

	return nil
}

// Type returns a string describing a MemberCommitmentsMessage type
// for marshaling purposes.
func (mcm *MemberCommitmentsMessage) Type() string {
	return "gjkr/member_commitments"
}

// Marshal converts this MemberCommitmentsMessage to a byte array suitable for
// network communication.
func (mcm *MemberCommitmentsMessage) Marshal() ([]byte, error) {
	commitmentBytes := make([][]byte, 0, len(mcm.commitments))
	for _, commitment := range mcm.commitments {
		commitmentBytes = append(commitmentBytes, commitment.Bytes())
	}

	return (&pb.MemberCommitments{
		SenderID:    memberIDToBytes(mcm.senderID),
		Commitments: commitmentBytes,
	}).Marshal()
}

// Unmarshal converts a byte array produced by Marshal to
// a MemberCommitmentsMessage
func (mcm *MemberCommitmentsMessage) Unmarshal(bytes []byte) error {
	pbMsg := pb.MemberCommitments{}
	if err := pbMsg.Unmarshal(bytes); err != nil {
		return err
	}

	mcm.senderID = bytesToMemberID(pbMsg.SenderID)

	var commitments []*big.Int
	for _, commitmentBytes := range pbMsg.Commitments {
		commitment := new(big.Int).SetBytes(commitmentBytes)
		commitments = append(commitments, commitment)
	}
	mcm.commitments = commitments

	return nil
}

// Type returns a string describing a OtherMemberSharesMessage type for marshaling
// purposes
func (omsm *OtherMemberSharesMessage) Type() string {
	return "gjkr/other_member_shares"
}

// Marshal converts this OtherMemberSharesMessage to a byte array suitable for
// network communication.
func (omsm *OtherMemberSharesMessage) Marshal() ([]byte, error) {
	return (&pb.OtherMemberShares{
		SenderID:        memberIDToBytes(omsm.senderID),
		ReceiverID:      memberIDToBytes(omsm.receiverID),
		EncryptedShareS: omsm.encryptedShareS,
		EncryptedShareT: omsm.encryptedShareT,
	}).Marshal()
}

// Unmarshal converts a byte array produced by Marshal to a OtherMemberSharesMessage.
func (omsm *OtherMemberSharesMessage) Unmarshal(bytes []byte) error {
	pbMsg := pb.OtherMemberShares{}
	if err := pbMsg.Unmarshal(bytes); err != nil {
		return err
	}

	omsm.senderID = bytesToMemberID(pbMsg.SenderID)
	omsm.receiverID = bytesToMemberID(pbMsg.ReceiverID)
	omsm.encryptedShareS = pbMsg.EncryptedShareS
	omsm.encryptedShareT = pbMsg.EncryptedShareT

	return nil
}

// Type returns a string describing a SecretSharesAccusationsMessage type
// for marshalling purposes.
func (ssam *SecretSharesAccusationsMessage) Type() string {
	return "gjkr/secret_shares_accusations"
}

// Marshal converts this SecretSharesAccusationsMessage to a byte array
// suitable for network communication.
func (ssam *SecretSharesAccusationsMessage) Marshal() ([]byte, error) {
	accusedIDsBytes := make([][]byte, 0, len(ssam.accusedIDs))
	for _, accusedID := range ssam.accusedIDs {
		accusedIDsBytes = append(accusedIDsBytes, memberIDToBytes(accusedID))
	}

	return (&pb.SecretSharesAccusations{
		SenderID:   memberIDToBytes(ssam.senderID),
		AccusedIDs: accusedIDsBytes,
	}).Marshal()
}

// Unmarshal converts a byte array produced by Marshal to
// a SecretSharesAccusationsMessage.
func (ssam *SecretSharesAccusationsMessage) Unmarshal(bytes []byte) error {
	pbMsg := pb.SecretSharesAccusations{}
	if err := pbMsg.Unmarshal(bytes); err != nil {
		return err
	}

	ssam.senderID = bytesToMemberID(pbMsg.SenderID)

	var accusedIDs []MemberID
	for _, accusedIDBytes := range pbMsg.AccusedIDs {
		accusedIDs = append(accusedIDs, bytesToMemberID(accusedIDBytes))
	}
	ssam.accusedIDs = accusedIDs

	return nil
}

// Type returns a string describing MemberPublicKeySharePointsMessage type for
// marshaling purposes
func (mpspm *MemberPublicKeySharePointsMessage) Type() string {
	return "gjkr/member_public_key_share_points"
}

// Marshal converts this MemberPublicKeySharePointsMessage to a byte array
// suitable for network communication.
func (mpspm *MemberPublicKeySharePointsMessage) Marshal() ([]byte, error) {
	keySharePoints := make([][]byte, 0, len(mpspm.publicKeySharePoints))
	for _, keySharePoint := range mpspm.publicKeySharePoints {
		keySharePoints = append(keySharePoints, keySharePoint.Bytes())
	}

	return (&pb.MemberPublicKeySharePoints{
		SenderID:             memberIDToBytes(mpspm.senderID),
		PublicKeySharePoints: keySharePoints,
	}).Marshal()
}

// Unmarshal converts a byte array produced by Marshal to
// a MemberPublicKeySharePointsMessage.
func (mpspm *MemberPublicKeySharePointsMessage) Unmarshal(bytes []byte) error {
	pbMsg := pb.MemberPublicKeySharePoints{}
	if err := pbMsg.Unmarshal(bytes); err != nil {
		return err
	}

	mpspm.senderID = bytesToMemberID(pbMsg.SenderID)

	var keySharePoints []*big.Int
	for _, keySharePointBytes := range pbMsg.PublicKeySharePoints {
		keySharePoint := new(big.Int).SetBytes(keySharePointBytes)
		keySharePoints = append(keySharePoints, keySharePoint)
	}
	mpspm.publicKeySharePoints = keySharePoints

	return nil
}

// Type returns a string describing PointsAccusationsMessage type for
// marshaling purposes.
func (pam *PointsAccusationsMessage) Type() string {
	return "gjkr/points_accusations_message"
}

// Marshal converts this PointsAccusationsMessage to a byte array suitable
// for network communication.
func (pam *PointsAccusationsMessage) Marshal() ([]byte, error) {
	accusedIDsBytes := make([][]byte, 0, len(pam.accusedIDs))
	for _, accusedID := range pam.accusedIDs {
		accusedIDsBytes = append(accusedIDsBytes, memberIDToBytes(accusedID))
	}

	return (&pb.PointsAccusations{
		SenderID:   memberIDToBytes(pam.senderID),
		AccusedIDs: accusedIDsBytes,
	}).Marshal()
}

// Unmarshal converts a byte array produced by Marshal to
// a PointsAccusationsMessage.
func (pam *PointsAccusationsMessage) Unmarshal(bytes []byte) error {
	pbMsg := pb.PointsAccusations{}
	if err := pbMsg.Unmarshal(bytes); err != nil {
		return err
	}

	pam.senderID = bytesToMemberID(pbMsg.SenderID)

	var accusedIDs []MemberID
	for _, accusedIDBytes := range pbMsg.AccusedIDs {
		accusedIDs = append(accusedIDs, bytesToMemberID(accusedIDBytes))
	}
	pam.accusedIDs = accusedIDs

	return nil
}

func memberIDToBytes(memberID MemberID) []byte {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, uint32(memberID))
	return bytes
}

func bytesToMemberID(bytes []byte) MemberID {
	return MemberID(binary.LittleEndian.Uint32(bytes))
}
