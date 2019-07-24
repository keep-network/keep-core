package operator

import (
	"crypto/ecdsa"
)

// PrivateKey represents peer's static key associated with an on-chain
// stake. It is used to authenticate the peer and for attributability (signing).
type PrivateKey = ecdsa.PrivateKey

// PublicKey represents peer's static key associated with an on-chain
// stake. It is used to authenticate the peer and for attributability
// (verification).
type PublicKey = ecdsa.PublicKey
