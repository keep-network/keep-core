package dkg

import (
	"github.com/dfinity/go-dfinity-crypto/bls"
	"github.com/keep-network/keep-core/pkg/beacon/relay/dkg/gen/pb"
)

// Marshal converts this JoinMessage to a byte array suitable for network
// communication.
func (m *JoinMessage) Marshal() ([]byte, error) {
	return (&pb.Join{Id: m.id.GetLittleEndian()}).Marshal()
}

// Unmarshal converts a byte array produced by Marshal to a JoinMessage.
func (m *JoinMessage) Unmarshal(bytes []byte) error {
	pbJoin := pb.Join{}
	err := pbJoin.Unmarshal(bytes)
	if err != nil {
		return err
	}

	m.id = &bls.ID{}
	err = m.id.SetLittleEndian(pbJoin.Id)
	if err != nil {
		return err
	}

	return nil
}

// Marshal converts this MemberCommitmentsMessage to a byte array suitable for network
// communication.
func (m *MemberCommitmentsMessage) Marshal() ([]byte, error) {
	commitmentBytes := make([][]byte, 0, len(m.Commitments))
	for _, commitment := range m.Commitments {
		commitmentBytes = append(commitmentBytes, commitment.Serialize())
	}

	pbCommitments :=
		pb.Commitments{
			Id:          m.id.GetLittleEndian(),
			Commitments: commitmentBytes}

	return pbCommitments.Marshal()
}

// Unmarshal converts a  byte array produced by Marshal to a
// MemberCommitmentsMessage.
func (m *MemberCommitmentsMessage) Unmarshal(bytes []byte) error {
	pbCommitments := pb.Commitments{}
	err := pbCommitments.Unmarshal(bytes)
	if err != nil {
		return err
	}

	m.id = &bls.ID{}
	err = m.id.SetLittleEndian(pbCommitments.Id)
	if err != nil {
		return err
	}

	m.Commitments = make([]bls.PublicKey, 0, len(pbCommitments.Commitments))
	for _, commitmentBytes := range pbCommitments.Commitments {
		pk := bls.PublicKey{}
		err = pk.Deserialize(commitmentBytes)
		if err != nil {
			return err
		}
		m.Commitments = append(m.Commitments, pk)
	}

	return nil
}

// Marshal converts this MemberShareMessage to a byte array suitable for network
// communication.
func (m *MemberShareMessage) Marshal() ([]byte, error) {
	pbShare :=
		pb.Share{
			Id:         m.id.GetLittleEndian(),
			ReceiverID: m.receiverID.GetLittleEndian(),
			Share:      m.Share.GetLittleEndian()}

	return pbShare.Marshal()
}

// Unmarshal converts a byte array produced by Marshal to a MemberShareMessage.
func (m *MemberShareMessage) Unmarshal(bytes []byte) error {
	pbShare := pb.Share{}
	err := pbShare.Unmarshal(bytes)
	if err != nil {
		return err
	}

	m.id = &bls.ID{}
	err = m.id.SetLittleEndian(pbShare.Id)
	if err != nil {
		return err
	}

	m.receiverID = &bls.ID{}
	err = m.receiverID.SetLittleEndian(pbShare.ReceiverID)
	if err != nil {
		return err
	}

	m.Share = &bls.SecretKey{}
	err = m.Share.SetLittleEndian(pbShare.Share)
	if err != nil {
		return err
	}

	return nil
}

// Marshal converts this AccusationsMessage to a byte array suitable for network
// communication.
func (m *AccusationsMessage) Marshal() ([]byte, error) {
	accusedIDBytes := make([][]byte, 0, len(m.accusedIDs))
	for _, accusedID := range m.accusedIDs {
		accusedIDBytes = append(accusedIDBytes, accusedID.GetLittleEndian())
	}

	pbAccusations :=
		pb.Accusations{
			Id:         m.id.GetLittleEndian(),
			AccusedIDs: accusedIDBytes}

	return pbAccusations.Marshal()
}

// Unmarshal converts a  byte array produced by Marshal to an
// AccusationsMessage.
func (m *AccusationsMessage) Unmarshal(bytes []byte) error {
	pbAccusations := pb.Accusations{}
	err := pbAccusations.Unmarshal(bytes)
	if err != nil {
		return err
	}

	m.id = &bls.ID{}
	err = m.id.SetLittleEndian(pbAccusations.Id)
	if err != nil {
		return err
	}

	m.accusedIDs = make([]bls.ID, 0, len(pbAccusations.AccusedIDs))
	for _, accusedIDBytes := range pbAccusations.AccusedIDs {
		id := bls.ID{}
		err = id.SetLittleEndian(accusedIDBytes)
		if err != nil {
			return err
		}
		m.accusedIDs = append(m.accusedIDs, id)
	}

	return nil
}

// Marshal converts this JustificationsMessage to a byte array suitable for
// network communication.
func (m *JustificationsMessage) Marshal() ([]byte, error) {
	justificationsMap := make(map[string][]byte, len(m.justifications))
	for id, sk := range m.justifications {
		justificationsMap[id.GetHexString()] = sk.GetLittleEndian()
	}

	pbJustifications :=
		pb.Justifications{
			Id:                 m.id.GetLittleEndian(),
			JustificationsByID: justificationsMap}

	return pbJustifications.Marshal()
}

// Unmarshal converts a byte array produced by Marshal to a
// JustificationsMessage.
func (m *JustificationsMessage) Unmarshal(bytes []byte) error {
	pbJustifications := pb.Justifications{}
	err := pbJustifications.Unmarshal(bytes)
	if err != nil {
		return err
	}

	m.id = &bls.ID{}
	err = m.id.SetLittleEndian(pbJustifications.Id)
	if err != nil {
		return err
	}

	m.justifications = make(map[bls.ID]bls.SecretKey, len(pbJustifications.JustificationsByID))
	for hexID, skBytes := range pbJustifications.JustificationsByID {
		id := bls.ID{}
		err = id.SetHexString(hexID)
		if err != nil {
			return err
		}
		sk := bls.SecretKey{}
		err = sk.SetLittleEndian(skBytes)
		if err != nil {
			return err
		}

		m.justifications[id] = sk
	}

	return nil
}
