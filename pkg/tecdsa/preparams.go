package tecdsa

import "github.com/binance-chain/tss-lib/ecdsa/keygen"

type PreParams struct {
	data keygen.LocalPreParams
}

func NewPreParams(data keygen.LocalPreParams) *PreParams {
	return &PreParams{data}
}
