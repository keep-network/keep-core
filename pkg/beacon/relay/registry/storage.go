package registry

import (
	"fmt"

	ds "github.com/keep-network/keep-core/pkg/disk_storage"

	"encoding/hex"
)

// Handle is an interface to handle memberships on disk
type Handle interface {
	Save(membership *Membership) error
}

type storage struct {
	diskStorage ds.DiskHandler
}

// NewStorage creates a new storage.
func NewStorage(diskStorage ds.DiskHandler) Handle {
	return &storage{
		diskStorage: diskStorage,
	}
}

// Save converts a membership suitable for disk storage.
func (s *storage) Save(membership *Membership) error {
	membershipBytes, err := membership.Marshal()
	if err != nil {
		return fmt.Errorf("marshalling of the membership failed: [%v]", err)
	}
	hexGroupPublicKey := hex.EncodeToString(membership.Signer.GroupPublicKeyBytes())

	return s.diskStorage.Save(membershipBytes, "/membership_"+hexGroupPublicKey)
}
