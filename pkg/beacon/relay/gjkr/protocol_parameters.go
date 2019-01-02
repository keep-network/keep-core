package gjkr

import (
	"math/big"

	"github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/keep-network/keep-core/pkg/altbn128"
)

// protocolParameters holds all cryptographic parameters that must be the same
// for all members in the group.
type protocolParameters struct {
	// `H = G*a` is a custom generator where `a` is unknown. It is used to
	// produce Pedersen commitments.
	H *bn256.G1
}

// newProtocolParameters creates a new instance of protocolParameters from the
// provided seed value which can be the previous random beacon's result.
// The seed is used to evaluate `H` parameter so that the discrete logarithm of
// `H` is unknown.
func newProtocolParameters(seed *big.Int) *protocolParameters {
	return &protocolParameters{
		H: altbn128.G1HashToPoint(seed.Bytes()),
	}
}
