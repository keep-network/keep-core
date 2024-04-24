package inactivity

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/keep-network/keep-core/pkg/protocol/group"
)

// Claim represents an inactivity claim.
type Claim struct {
	Nonce                  *big.Int
	WalletPublicKey        *ecdsa.PublicKey
	InactiveMembersIndexes []group.MemberIndex
	HeartbeatFailed        bool
}

const ClaimSignatureHashByteSize = 32

// ClaimSignatureHash is a signature hash of the inactivity claim. The hashing
// algorithm used depends on the client code.
type ClaimSignatureHash [ClaimSignatureHashByteSize]byte

// ClaimSignatureHashFromBytes converts bytes slice to ClaimSignatureHash.
// It requires provided bytes slice size to be exactly
// ClaimSignatureHashByteSize.
func ClaimSignatureHashFromBytes(bytes []byte) (ClaimSignatureHash, error) {
	var hash ClaimSignatureHash

	if len(bytes) != ClaimSignatureHashByteSize {
		return hash, fmt.Errorf(
			"bytes length is not equal %v", ClaimSignatureHashByteSize,
		)
	}
	copy(hash[:], bytes[:])

	return hash, nil
}
