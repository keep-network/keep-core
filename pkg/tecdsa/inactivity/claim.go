package inactivity

import (
	"crypto/ecdsa"
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
