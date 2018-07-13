package tecdsa

import (
	"github.com/keep-network/keep-core/pkg/tecdsa/curve"
	"github.com/keep-network/keep-core/pkg/tecdsa/zkp"
	"github.com/keep-network/paillier"
)

// InitMessage is a message payload that carries the sender's share of
// `dsaKeyShare` during T-ECDSA distributed DSA key generation as well as
// proofs of correctness for the shares. Sender's share is encrypted with (t, n)
// Paillier threshold key. The message is expected to be broadcast publicly.
type InitMessage struct {
	secretKeyShare *paillier.Cypher
	publicKeyShare *curve.Point

	rangeProof *zkp.DsaPaillierKeyRangeProof
}

// IsValid checks secret and public key share against zero knowledge range proof
// shipped alongside them. This function should be called for each received
// InitMessage before it's combined to a final key.
func (im *InitMessage) IsValid(zkpParams *zkp.PublicParameters) bool {
	return im.rangeProof.Verify(im.secretKeyShare, im.publicKeyShare, zkpParams)
}
