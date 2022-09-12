package signing

import (
	"github.com/keep-network/keep-core/pkg/tecdsa"
)

// Result of the tECDSA signing protocol.
type Result struct {
	// Signature is the tECDSA signature produced as result of the tECDSA
	// signing process.
	Signature *tecdsa.Signature
}
