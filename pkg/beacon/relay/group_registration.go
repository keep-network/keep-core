package relay

import "math/big"

// GroupRegistration represents a registered group in the threshold relay with a
// public key, that is considered active at ActivationBlockHeight, and was
// spawned by the relay request with id, RequestID.
type GroupRegistration struct {
	GroupPublicKey        []byte
	RequestID             *big.Int
	ActivationBlockHeight *big.Int
}
