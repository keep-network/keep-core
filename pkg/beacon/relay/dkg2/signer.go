package dkg2

import (
	"math/big"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/keep-network/keep-core/pkg/beacon/relay/gjkr"
	"github.com/keep-network/keep-core/pkg/bls"
)

type ThresholdSigner struct {
	memberID             gjkr.MemberID
	groupPublicKey       *bn256.G1
	groupPrivateKeyShare *big.Int
}

func (ts *ThresholdSigner) MemberID() gjkr.MemberID {
	return ts.memberID
}

func (ts *ThresholdSigner) GroupPublicKeyBytes() []byte {
	return ts.groupPublicKey.Marshal()
}

func (ts *ThresholdSigner) SignatureShare(message []byte) *bn256.G1 {
	return bls.Sign(ts.groupPrivateKeyShare, message)
}

func (ts *ThresholdSigner) CompleteSignature(signatureShares []*bn256.G1) *bn256.G1 {
	return bls.AggregateG1Points(signatureShares)
}
