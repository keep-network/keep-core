package gjkr

import (
	"fmt"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
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
		SenderID: uint32(jm.senderID),
	}).Marshal()
}

// Unmarshal converts a byte array produced by Marshal to a JoinMessage.
func (jm *JoinMessage) Unmarshal(bytes []byte) error {
	pbMsg := pb.Join{}
	if err := pbMsg.Unmarshal(bytes); err != nil {
		return err
	}
	jm.senderID = MemberID(pbMsg.SenderID)
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
		SenderID:            uint32(epkm.senderID),
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

	epkm.senderID = MemberID(pbMsg.SenderID)

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
		commitmentBytes = append(commitmentBytes, commitment.Marshal())
	}

	return (&pb.MemberCommitments{
		SenderID:    uint32(mcm.senderID),
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

	mcm.senderID = MemberID(pbMsg.SenderID)

	var commitments []*bn256.G1
	for _, commitmentBytes := range pbMsg.Commitments {
		commitment := new(bn256.G1)
		_, err := commitment.Unmarshal(commitmentBytes)
		if err != nil {
			return fmt.Errorf(
				"could not unmarshal member's commitment [%v]",
				err,
			)
		}
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
	pbShares := make(map[uint32]*pb.PeerShares_Shares)
	for memberID, shares := range psm.shares {
		pbShares[uint32(memberID)] = &pb.PeerShares_Shares{
			EncryptedShareS: shares.encryptedShareS,
			EncryptedShareT: shares.encryptedShareT,
		}
	}

	return (&pb.PeerShares{
		SenderID: uint32(psm.senderID),
		Shares:   pbShares,
	}).Marshal()
}

// Unmarshal converts a byte array produced by Marshal to a PeerSharesMessage.
func (psm *PeerSharesMessage) Unmarshal(bytes []byte) error {
	pbMsg := pb.PeerShares{}
	if err := pbMsg.Unmarshal(bytes); err != nil {
		return err
	}

	psm.senderID = MemberID(pbMsg.SenderID)

	shares := make(map[MemberID]*peerShares)
	for memberID, pbShares := range pbMsg.Shares {
		shares[MemberID(memberID)] = &peerShares{
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
		SenderID:           uint32(ssam.senderID),
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

	ssam.senderID = MemberID(pbMsg.SenderID)

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
		keySharePoints = append(keySharePoints, keySharePoint.Marshal())
	}

	return (&pb.MemberPublicKeySharePoints{
		SenderID:             uint32(mpspm.senderID),
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

	mpspm.senderID = MemberID(pbMsg.SenderID)

	var keySharePoints []*bn256.G2
	for _, keySharePointBytes := range pbMsg.PublicKeySharePoints {
		keySharePoint := new(bn256.G2)
		_, err := keySharePoint.Unmarshal(keySharePointBytes)
		if err != nil {
			return fmt.Errorf(
				"could not unmarshal member's key share point [%v]",
				err,
			)
		}
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
		SenderID:           uint32(pam.senderID),
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

	pam.senderID = MemberID(pbMsg.SenderID)

	accusedMembersKeys, err := unmarshalPrivateKeyMap(pbMsg.AccusedMembersKeys)
	if err != nil {
		return nil
	}

	pam.accusedMembersKeys = accusedMembersKeys

	return nil
}

// Type returns a string describing DisqualifiedEphemeralKeysMessage type for
// marshalling purposes.
func (dekm *DisqualifiedEphemeralKeysMessage) Type() string {
	return "gjkr/disqualified_ephemeral_keys_message"
}

// Marshal converts this DisqualifiedEphemeralKeysMessage to a byte array
// suitable for network communication.
func (dekm *DisqualifiedEphemeralKeysMessage) Marshal() ([]byte, error) {
	return (&pb.DisqualifiedEphemeralKeys{
		SenderID:    uint32(dekm.senderID),
		PrivateKeys: marshalPrivateKeyMap(dekm.privateKeys),
	}).Marshal()
}

// Unmarshal converts a byte array produced by Marshal to
// a DisqualifiedEphemeralKeysMessage.
func (dekm *DisqualifiedEphemeralKeysMessage) Unmarshal(bytes []byte) error {
	pbMsg := pb.DisqualifiedEphemeralKeys{}
	if err := pbMsg.Unmarshal(bytes); err != nil {
		return err
	}

	dekm.senderID = MemberID(pbMsg.SenderID)

	privateKeys, err := unmarshalPrivateKeyMap(pbMsg.PrivateKeys)
	if err != nil {
		return err
	}

	dekm.privateKeys = privateKeys

	return nil
}

func marshalPublicKeyMap(
	publicKeys map[MemberID]*ephemeral.PublicKey,
) map[uint32][]byte {
	marshalled := make(map[uint32][]byte, len(publicKeys))
	for id, publicKey := range publicKeys {
		marshalled[uint32(id)] = publicKey.Marshal()
	}
	return marshalled
}

func unmarshalPublicKeyMap(
	publicKeys map[uint32][]byte,
) (map[MemberID]*ephemeral.PublicKey, error) {
	var unmarshalled = make(map[MemberID]*ephemeral.PublicKey, len(publicKeys))
	for memberID, publicKeyBytes := range publicKeys {
		publicKey, err := ephemeral.UnmarshalPublicKey(publicKeyBytes)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal public key [%v]", err)
		}

		unmarshalled[MemberID(memberID)] = publicKey

	}

	return unmarshalled, nil
}

func marshalPrivateKeyMap(
	privateKeys map[MemberID]*ephemeral.PrivateKey,
) map[uint32][]byte {
	marshalled := make(map[uint32][]byte, len(privateKeys))
	for id, privateKey := range privateKeys {
		marshalled[uint32(id)] = privateKey.Marshal()
	}
	return marshalled
}

func unmarshalPrivateKeyMap(
	privateKeys map[uint32][]byte,
) (map[MemberID]*ephemeral.PrivateKey, error) {
	var unmarshalled = make(map[MemberID]*ephemeral.PrivateKey, len(privateKeys))
	for memberID, privateKeyBytes := range privateKeys {
		unmarshalled[MemberID(memberID)] = ephemeral.UnmarshalPrivateKey(privateKeyBytes)
	}

	return unmarshalled, nil
}
