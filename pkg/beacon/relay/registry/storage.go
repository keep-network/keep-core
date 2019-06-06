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

var (
	currentDir = "current"
	archiveDir = "archive"
)

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

	err = ps.createStorageDir(hexGroupPublicKey)
	if err != nil {
		return err
	}

	return ps.handle.Save(membershipBytes, hexGroupPublicKey, "/membership_"+fmt.Sprint(membership.Signer.MemberID()))
}

func (ps *persistentStorage) createStorageDir(groupPublicKey string) error {
	err := ps.handle.CreateDir(ps.handle.GetDataDir(), currentDir)
	if err != nil {
		return err
	}

	currentPath := fmt.Sprintf("%s/%s", ps.handle.GetDataDir(), currentDir)
	err = ps.handle.CreateDir(currentPath, groupPublicKey)
	if err != nil {
		return err
	}

	return nil
}

func (ps *persistentStorage) readAll() ([]*Membership, error) {
	memberships := []*Membership{}

	currentPath := fmt.Sprintf("%s/%s", ps.handle.GetDataDir(), currentDir)
	bytesMemberships, err := ps.handle.ReadAll(currentPath)
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
	from := fmt.Sprintf("%s/%s/%s", ps.handle.GetDataDir(), currentDir, groupName)
	to := fmt.Sprintf("%s/%s/%s", ps.handle.GetDataDir(), archiveDir, groupName)

	err := ps.handle.CreateDir(ps.handle.GetDataDir(), archiveDir)
	if err != nil {
		return err
	}

	return ps.handle.Archive(from, to)
}
