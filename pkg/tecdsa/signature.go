package tecdsa

import (
	"fmt"
	"github.com/bnb-chain/tss-lib/common"
	"math/big"
)

// Signature holds a signature in a form of two big.Int `r` and `s` values and
// a recovery ID value in {0, 1, 2, 3}.
//
// The signature is chain-agnostic. Some chains (e.g. Ethereum and BTC)
// requires `v` to start from 27. Please consult the documentation about
// what the particular chain expects.
type Signature struct {
	R          *big.Int
	S          *big.Int
	RecoveryID int
}

// NewSignature constructs a new instance of the tECDSA signature based on
// the signing result.
func NewSignature(data *common.SignatureData) *Signature {
	// `SignatureData` contains recovery ID as a byte slice. Only the first
	// byte is relevant and is converted to `int`.
	recoveryBytes := data.GetSignatureRecovery()
	recoveryInt := 0
	recoveryInt = (recoveryInt << 8) | int(recoveryBytes[0])

	return &Signature{
		R:          new(big.Int).SetBytes(data.GetR()),
		S:          new(big.Int).SetBytes(data.GetS()),
		RecoveryID: recoveryInt,
	}
}

// String formats Signature to a string that contains R and S values
// as hexadecimals.
func (s *Signature) String() string {
	return fmt.Sprintf(
		"R: %#x, S: %#x, RecoveryID: %d",
		s.R,
		s.S,
		s.RecoveryID,
	)
}
