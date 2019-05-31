package registry

import (
	"fmt"

	"github.com/keep-network/keep-core/pkg/persistence"

	"encoding/hex"
)

type storage interface {
	save(membership *Membership) error
}

type persistentStorage struct {
	handle persistence.Handle
}

func newStorage(persistence persistence.Handle) storage {
	return &persistentStorage{
		handle: persistence,
	}
}

// Save converts a membership suitable for disk storage.
func (ps *persistentStorage) save(membership *Membership) error {
	membershipBytes, err := membership.Marshal()
	if err != nil {
		return fmt.Errorf("marshalling of the membership failed: [%v]", err)
	}
	hexGroupPublicKey := hex.EncodeToString(membership.Signer.GroupPublicKeyBytes())

	return ps.handle.Save(membershipBytes, hexGroupPublicKey, "/membership_"+fmt.Sprint(membership.Signer.MemberID()))
}
