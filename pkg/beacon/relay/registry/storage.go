package registry

import (
	"fmt"

	"github.com/keep-network/keep-core/pkg/persistence"

	"encoding/hex"
)

type storage interface {
	save(membership *Membership) error
	readAll() ([]*Membership, error)
	archive(groupName string) error
}

type persistentStorage struct {
	handle persistence.Handle
}

func newStorage(persistence persistence.Handle) storage {
	return &persistentStorage{
		handle: persistence,
	}
}

func (ps *persistentStorage) save(membership *Membership) error {
	membershipBytes, err := membership.Marshal()
	if err != nil {
		return fmt.Errorf("marshalling of the membership failed: [%v]", err)
	}

	hexGroupPublicKey := hex.EncodeToString(membership.Signer.GroupPublicKeyBytes())

	return ps.handle.Save(membershipBytes, hexGroupPublicKey, "/membership_"+fmt.Sprint(membership.Signer.MemberID()))
}

func (ps *persistentStorage) readAll() ([]*Membership, error) {
	memberships := []*Membership{}

	bytesMemberships, err := ps.handle.ReadAll()
	if err != nil {
		return nil, err
	}

	for _, byteMembership := range bytesMemberships {
		membership := &Membership{}
		err := membership.Unmarshal(byteMembership)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal membership: [%v]", err)
		}
		memberships = append(memberships, membership)
	}

	return memberships, nil
}

func (ps *persistentStorage) archive(groupName string) error {
	return ps.handle.Archive(groupName)
}
