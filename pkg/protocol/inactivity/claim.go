package inactivity

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"sort"

	"github.com/keep-network/keep-core/pkg/protocol/group"
)

// ClaimPreimage represents an inactivity claim preimage.
type ClaimPreimage struct {
	Nonce                  *big.Int
	WalletPublicKey        *ecdsa.PublicKey
	InactiveMembersIndexes []group.MemberIndex
	HeartbeatFailed        bool
}

// GetInactiveMembersIndexes returns the indexes of inactive members.
// The original slice is copied to avoid concurrency issues if the claim object
// is shared between many goroutines. The returned indexes are sorted.
func (c *ClaimPreimage) GetInactiveMembersIndexes() []group.MemberIndex {
	sortedIndexes := make([]group.MemberIndex, len(c.InactiveMembersIndexes))

	copy(sortedIndexes, c.InactiveMembersIndexes)

	sort.Slice(sortedIndexes, func(i, j int) bool {
		return sortedIndexes[i] < sortedIndexes[j]
	})

	return sortedIndexes
}

const ClaimHashByteSize = 32

// ClaimHash is a hash of the inactivity claim. The hashing algorithm used
// depends on the client code.
type ClaimHash [ClaimHashByteSize]byte

// ClaimHashFromBytes converts bytes slice to ClaimHash. It requires provided
// bytes slice size to be exactly ClaimHashByteSize.
func ClaimHashFromBytes(bytes []byte) (ClaimHash, error) {
	var hash ClaimHash

	if len(bytes) != ClaimHashByteSize {
		return hash, fmt.Errorf(
			"bytes length is not equal %v", ClaimHashByteSize,
		)
	}
	copy(hash[:], bytes)

	return hash, nil
}
