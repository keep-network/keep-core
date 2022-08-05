package tecdsa

import (
	"crypto/ecdsa"
	"github.com/binance-chain/tss-lib/ecdsa/keygen"
)

// PrivateKeyShare represents a private key share used to produce tECDSA
// signatures. Private key shares are generated as result of the tECDSA
// distributed key generation (DKG) process.
type PrivateKeyShare struct {
	data keygen.LocalPartySaveData
}

// NewPrivateKeyShare constructs a new instance of the tECDSA public key
// share based on the DKG result.
func NewPrivateKeyShare(data keygen.LocalPartySaveData) *PrivateKeyShare {
	return &PrivateKeyShare{data}
}

// PublicKey returns the ECDSA public key corresponding to the given tECDSA
// private key share.
func (pks *PrivateKeyShare) PublicKey() *ecdsa.PublicKey {
	return &ecdsa.PublicKey{
		Curve: pks.data.ECDSAPub.Curve(),
		X:     pks.data.ECDSAPub.X(),
		Y:     pks.data.ECDSAPub.Y(),
	}
}
