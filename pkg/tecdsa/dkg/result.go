package dkg

import (
	"crypto/elliptic"
	"fmt"

	"github.com/keep-network/keep-core/pkg/tecdsa"
)

// Result of distributed key generation protocol.
type Result struct {
	// Misbehaved members are all members either inactive or disqualified.
	// Misbehaved members are represented as a slice of bytes for optimizing
	// on-chain storage. Each byte is an inactive or disqualified member index.
	Misbehaved []byte
	// PrivateKeyShare is the tECDSA private key share required to operate
	// in the signing group generated as result of the DKG protocol.
	PrivateKeyShare *tecdsa.PrivateKeyShare
}

// GetGroupPublicKeyBytes returns the public key corresponding to the private
// key share generated during the DKG protocol execution.
func (r *Result) GetGroupPublicKeyBytes() ([]byte, error) {
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
