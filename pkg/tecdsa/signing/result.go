package signing

import (
	"fmt"
	"math/big"
)

// Result of the tECDSA signing protocol.
type Result struct {
	R          *big.Int
	S          *big.Int
	RecoveryID int
}

func (r *Result) String() string {
	return fmt.Sprintf(
		"R: %#x, S: %#x, RecoveryID: %d",
		r.R,
		r.S,
		r.RecoveryID,
	)
}
