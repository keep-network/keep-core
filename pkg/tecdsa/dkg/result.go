package dkg

import (
	"crypto/elliptic"
	"fmt"
	"sort"

	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/tecdsa"
)

// Result of distributed key generation protocol.
type Result struct {
	// Group represents the group state, including members, disqualified,
	// and inactive members.
	Group *group.Group
	// PrivateKeyShare is the tECDSA private key share required to operate
	// in the signing group generated as result of the DKG protocol.
	PrivateKeyShare *tecdsa.PrivateKeyShare
}

// GroupPublicKeyBytes returns the public key corresponding to the private
// key share generated during the DKG protocol execution.
func (r *Result) GroupPublicKeyBytes() ([]byte, error) {
	if r.PrivateKeyShare == nil {
		return nil, fmt.Errorf(
			"cannot retrieve group public key as private key share is nil",
		)
	}

	publicKey := r.PrivateKeyShare.PublicKey()

	return elliptic.Marshal(
		publicKey.Curve,
		publicKey.X,
		publicKey.Y,
	), nil
}

// MisbehavedMembersIndexes returns the indexes of group members that misbehaved
// during the DKG procedure. The indexes are sorted.
func (r *Result) MisbehavedMembersIndexes() []group.MemberIndex {
	// Merge inactive and disqualified member indexes into 'misbehaved' set.
	misbehaving := make(map[group.MemberIndex]bool)
	for _, inactiveMemberIndex := range r.Group.InactiveMemberIndexes() {
		misbehaving[inactiveMemberIndex] = true
	}
	for _, disqualifiedMemberIndex := range r.Group.DisqualifiedMemberIndexes() {
		misbehaving[disqualifiedMemberIndex] = true
	}

	// Convert misbehaving member indexes set into sorted list.
	var sorted []group.MemberIndex
	for m := range misbehaving {
		sorted = append(sorted, m)
	}
	sort.Slice(sorted[:], func(i, j int) bool {
		return sorted[i] < sorted[j]
	})

	return sorted
}

const ResultHashByteSize = 32

// ResultHash is a 256-bit hash of DKG Result. The hashing algorithm used
// depends on the client code.
type ResultHash [ResultHashByteSize]byte

// ResultHashFromBytes converts bytes slice to DKG Result Hash. It requires
// provided bytes slice size to be exactly 32 bytes.
func ResultHashFromBytes(bytes []byte) (ResultHash, error) {
	var hash ResultHash

	if len(bytes) != ResultHashByteSize {
		return hash, fmt.Errorf(
			"bytes length is not equal %v", ResultHashByteSize,
		)
	}
	copy(hash[:], bytes[:])

	return hash, nil
}
