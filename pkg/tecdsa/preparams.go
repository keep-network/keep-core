package tecdsa

import "github.com/binance-chain/tss-lib/ecdsa/keygen"

// PreParams represents tECDSA DKG pre-parameters that were not yet consumed
// by DKG protocol execution.
type PreParams struct {
	data keygen.LocalPreParams
}

// NewPreParams constructs a new instance of tECDSA DKG pre-parameters based on
// the generated numbers.
func NewPreParams(data keygen.LocalPreParams) *PreParams {
	return &PreParams{data}
}
