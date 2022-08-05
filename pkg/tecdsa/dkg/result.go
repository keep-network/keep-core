package dkg

import "github.com/keep-network/keep-core/pkg/tecdsa"

// Result of distributed key generation protocol.
type Result struct {
	// PrivateKeyShare is the tECDSA private key share required to operate
	// in the signing group generated as result of the DKG protocol.
	PrivateKeyShare *tecdsa.PrivateKeyShare
}
