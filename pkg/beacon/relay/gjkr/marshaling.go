package gjkr

import (
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/keep-network/keep-core/pkg/beacon/relay/gjkr/gen/pb"
	"github.com/keep-network/keep-core/pkg/net/ephemeral"
)

// Type returns a string describing a JoinMessage type for marshalling purposes.
func (jm *JoinMessage) Type() string {
	return "gjkr/join_message"
}

// Marshal converts this JoinMessage to a byte array suitable for network
// communication.
func (jm *JoinMessage) Marshal() ([]byte, error) {
	return (&pb.Join{
		SenderID: memberIDToBytes(jm.SenderID),
	}).Marshal()
}

// Unmarshal converts a byte array produced by Marshal to a JoinMessage.
func (jm *JoinMessage) Unmarshal(bytes []byte) error {
	pbMsg := pb.Join{}
	if err := pbMsg.Unmarshal(bytes); err != nil {
		return err
	}
	jm.SenderID = bytesToMemberID(pbMsg.SenderID)
	return nil
}

// Type returns a string describing an EphemeralPublicKeyMessage type for
// marshaling purposes.
func (epkm *EphemeralPublicKeyMessage) Type() string {
	return "gjkr/ephemeral_public_key"
}

// Marshal converts this EphemeralPublicKeyMessage to a byte array suitable for
// network communication.
func (epkm *EphemeralPublicKeyMessage) Marshal() ([]byte, error) {
	return (&pb.EphemeralPublicKey{
		SenderID:            memberIDToBytes(epkm.senderID),
		EphemeralPublicKeys: marshalPublicKeyMap(epkm.ephemeralPublicKeys),
	}).Marshal()
}

// Unmarshal converts a byte array produced by Marshal to
// an EphemeralPublicKeyMessage
func (epkm *EphemeralPublicKeyMessage) Unmarshal(bytes []byte) error {
	pbMsg := pb.EphemeralPublicKey{}
	if err := pbMsg.Unmarshal(bytes); err != nil {
		return err
	}

	epkm.senderID = bytesToMemberID(pbMsg.SenderID)

	ephemeralPublicKeys, err := unmarshalPublicKeyMap(pbMsg.EphemeralPublicKeys)
	if err != nil {
		return err
	}

	epkm.ephemeralPublicKeys = ephemeralPublicKeys

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

// Type returns a string describing a PeerSharesMessage type for marshaling
// purposes
func (psm *PeerSharesMessage) Type() string {
	return "gjkr/peer_shares"
}

// Marshal converts this PeerSharesMessage to a byte array suitable for
// network communication.
func (psm *PeerSharesMessage) Marshal() ([]byte, error) {
	pbShares := make(map[string]*pb.PeerShares_Shares)
	for member, shares := range psm.shares {
		pbShares[member.HexString()] = &pb.PeerShares_Shares{
			EncryptedShareS: shares.encryptedShareS,
			EncryptedShareT: shares.encryptedShareT,
		}
	}

	return (&pb.PeerShares{
		SenderID: memberIDToBytes(psm.senderID),
		Shares:   pbShares,
	}).Marshal()
}

// Unmarshal converts a byte array produced by Marshal to a PeerSharesMessage.
func (psm *PeerSharesMessage) Unmarshal(bytes []byte) error {
	pbMsg := pb.PeerShares{}
	if err := pbMsg.Unmarshal(bytes); err != nil {
		return err
	}

	psm.senderID = bytesToMemberID(pbMsg.SenderID)

	shares := make(map[MemberID]*peerShares)
	for memberIDHex, pbShares := range pbMsg.Shares {
		memberID, err := MemberIDFromHex(memberIDHex)
		if err != nil {
			return fmt.Errorf("could not unmarshal member's ID [%v]", err)
		}
		shares[memberID] = &peerShares{
			encryptedShareS: pbShares.EncryptedShareS,
			encryptedShareT: pbShares.EncryptedShareT,
		}
	}

	psm.shares = shares

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
	return (&pb.SecretSharesAccusations{
		SenderID:           memberIDToBytes(ssam.senderID),
		AccusedMembersKeys: marshalPrivateKeyMap(ssam.accusedMembersKeys),
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

	accusedMembersKeys, err := unmarshalPrivateKeyMap(pbMsg.AccusedMembersKeys)
	if err != nil {
		return nil
	}

	ssam.accusedMembersKeys = accusedMembersKeys

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
	return (&pb.PointsAccusations{
		SenderID:           memberIDToBytes(pam.senderID),
		AccusedMembersKeys: marshalPrivateKeyMap(pam.accusedMembersKeys),
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

	accusedMembersKeys, err := unmarshalPrivateKeyMap(pbMsg.AccusedMembersKeys)
	if err != nil {
		return nil
	}

	pam.accusedMembersKeys = accusedMembersKeys

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

func marshalPublicKeyMap(
	publicKeys map[MemberID]*ephemeral.PublicKey,
) map[string][]byte {
	marshalled := make(map[string][]byte, len(publicKeys))
	for id, publicKey := range publicKeys {
		marshalled[id.HexString()] = publicKey.Marshal()
	}
	return marshalled
}

func unmarshalPublicKeyMap(
	publicKeys map[string][]byte,
) (map[MemberID]*ephemeral.PublicKey, error) {
	var unmarshalled = make(map[MemberID]*ephemeral.PublicKey, len(publicKeys))
	for memberIDHex, publicKeyBytes := range publicKeys {
		memberID, err := MemberIDFromHex(memberIDHex)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal member's ID [%v]", err)
		}

		publicKey, err := ephemeral.UnmarshalPublicKey(publicKeyBytes)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal public key [%v]", err)
		}

		unmarshalled[memberID] = publicKey

	}

	return unmarshalled, nil
}

func marshalPrivateKeyMap(
	privateKeys map[MemberID]*ephemeral.PrivateKey,
) map[string][]byte {
	marshalled := make(map[string][]byte, len(privateKeys))
	for id, privateKey := range privateKeys {
		marshalled[id.HexString()] = privateKey.Marshal()
	}
	return marshalled
}

func unmarshalPrivateKeyMap(
	privateKeys map[string][]byte,
) (map[MemberID]*ephemeral.PrivateKey, error) {
	var unmarshalled = make(map[MemberID]*ephemeral.PrivateKey, len(privateKeys))
	for memberIDHex, privateKeyBytes := range privateKeys {
		memberID, err := MemberIDFromHex(memberIDHex)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal member's ID [%v]", err)
		}
		unmarshalled[memberID] = ephemeral.UnmarshalPrivateKey(privateKeyBytes)
	}

	return unmarshalled, nil
}
