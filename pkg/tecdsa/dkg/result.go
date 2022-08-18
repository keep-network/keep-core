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

// MisbehavedMembersBytes returns the IDs of group members that misbehaved during
// the DKG procedure. The IDs are sorted and returned as bytes.
func (r *Result) MisbehavedMembersBytes() []byte {
	// Merge IA and DQ into 'misbehaved' set.
	misbehaving := make(map[group.MemberIndex]bool)
	for _, inactiveMemberID := range r.Group.InactiveMemberIDs() {
		misbehaving[inactiveMemberID] = true
	}
	for _, disqualifiedMemberID := range r.Group.DisqualifiedMemberIDs() {
		misbehaving[disqualifiedMemberID] = true
	}

	// Convert misbehaving member IDs set into sorted list.
	var sorted []group.MemberIndex
	for m := range misbehaving {
		sorted = append(sorted, m)
	}
	sort.Slice(sorted[:], func(i, j int) bool {
		return sorted[i] < sorted[j]
	})

	// Convert sorted list of member indexes into bytes.
	bytes := make([]byte, len(sorted))
	for i, m := range sorted {
		bytes[i] = byte(m)
	}

	return bytes
}

const hashByteSize = 32

// ResultHash is a 256-bit hash of DKG Result. The hashing algorithm should
// be the same as the one used on-chain.
type ResultHash [hashByteSize]byte

// ResultsVotes is a map of votes for each DKG Result.
type ResultsVotes map[ResultHash]int

// ResultHashFromBytes converts bytes slice to DKG Result Hash. It requires
// provided bytes slice size to be exactly 32 bytes.
func ResultHashFromBytes(bytes []byte) (ResultHash, error) {
	var hash ResultHash

	if len(bytes) != hashByteSize {
		return hash, fmt.Errorf("bytes length is not equal %v", hashByteSize)
	}
	copy(hash[:], bytes[:])

	return hash, nil
}
