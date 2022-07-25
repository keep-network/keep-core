package dkg

import (
	"math/big"

	"github.com/keep-network/keep-core/pkg/chain"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/keep-network/keep-core/pkg/altbn128"
	"github.com/keep-network/keep-core/pkg/bls"
	"github.com/keep-network/keep-core/pkg/protocol/group"
)

// ThresholdSigner is created from GJKR group Member when DKG protocol completed
// successfully and each group member is ready to sign. ThresholdSigner contains
// its own private key share of group public key that should never be publicly
// revealed. It also contains group's public key and ID of GJKR Member
// represented by this ThresholdSigner instance.
type ThresholdSigner struct {
	memberIndex          group.MemberIndex
	groupPublicKey       *bn256.G2
	groupPrivateKeyShare *big.Int
	groupPublicKeyShares map[group.MemberIndex]*bn256.G2
	groupOperators       []chain.Address
}

// NewThresholdSigner returns a new ThresholdSigner
func NewThresholdSigner(
	memberIndex group.MemberIndex,
	groupPublicKey *bn256.G2,
	groupPrivateKeyShare *big.Int,
	groupPublicKeyShares map[group.MemberIndex]*bn256.G2,
	groupOperators []chain.Address,
) *ThresholdSigner {
	return &ThresholdSigner{
		memberIndex:          memberIndex,
		groupPublicKey:       groupPublicKey,
		groupPrivateKeyShare: groupPrivateKeyShare,
		groupPublicKeyShares: groupPublicKeyShares,
		groupOperators:       groupOperators,
	}
}

// MemberID returns GJKR MemberID represented by this ThresholdSigner.
func (ts *ThresholdSigner) MemberID() group.MemberIndex {
	return ts.memberIndex
}

// GroupPublicKeyBytes returns group public key bytes in an uncompressed form.
func (ts *ThresholdSigner) GroupPublicKeyBytes() []byte {
	return ts.groupPublicKey.Marshal()
}

// GroupPublicKeyBytesCompressed returns group public key bytes in a compressed
// form.
func (ts *ThresholdSigner) GroupPublicKeyBytesCompressed() []byte {
	altbn128GroupPublicKey := altbn128.G2Point{G2: ts.groupPublicKey}
	return altbn128GroupPublicKey.Compress()
}

// CalculateSignatureShare takes the message and calculates signer's signature
// share over that message.
func (ts *ThresholdSigner) CalculateSignatureShare(message *bn256.G1) *bn256.G1 {
	return bls.SignG1(ts.groupPrivateKeyShare, message)
}

// CompleteSignature accepts signature shares from all group threshold signers
// and produces a final group signature from them. Input slice should contain
// signature of the current signer as well. We parameterize the threshold (number
// of honest members we require), as we recover a threshold signature, not an
// aggregate signature (one which would require all members).
func (ts *ThresholdSigner) CompleteSignature(
	signatureShares []*bls.SignatureShare,
	honestThreshold int,
) (*bn256.G1, error) {
	return bls.RecoverSignature(signatureShares, honestThreshold)
}

// GroupPublicKeyShares returns group public key shares for each
// individual member of the group.
func (ts *ThresholdSigner) GroupPublicKeyShares() map[group.MemberIndex]*bn256.G2 {
	return ts.groupPublicKeyShares
}

// GroupOperators returns operators being members of the group.
func (ts *ThresholdSigner) GroupOperators() []chain.Address {
	return ts.groupOperators
}
