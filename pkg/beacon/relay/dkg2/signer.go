package dkg2

import (
	"math/big"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/keep-network/keep-core/pkg/bls"
)

type thresholdSigner struct {
	groupPrivateKeyShare *big.Int
}

func (ts *thresholdSigner) signatureShare(message []byte) *bn256.G1 {
	return bls.Sign(ts.groupPrivateKeyShare, message)
}

func (ts *thresholdSigner) completeSignature(signatureShares []*bn256.G1) *bn256.G1 {
	return bls.AggregateG1Points(signatureShares)
}
