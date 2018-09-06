package tecdsa

import (
	"math/big"

	"github.com/keep-network/keep-core/pkg/tecdsa/curve"
	"github.com/keep-network/paillier"
)

// ThresholdEcdsaKey represents DSA key for a fully initialized Signer.
//
// For (secretKey, publicKey = g^secretKey) DSA key pair, publicKey
// and E(secretKey) are made public, where E is an additively homomorphic
// encryption scheme. This is an implicit (t, n) secret sharing of secretKey
// since the decryption key of E is shared among the n Signers.
type ThresholdEcdsaKey struct {
	secretKey *paillier.Cypher
	PublicKey *curve.Point
}

// ecdsaKeyShare represents a share of DSA key owned by LocalSigner before
// it's fully initialized into Signer.
//
// Each `LocalSigner` generates a share of secret and public DSA key.
// `publicKeyShare` is broadcasted to other signers along with
// `E(secretKeyShare)` where E is an additively homomorphic encryption scheme.
// It lets to compute:
//
// E(secretKey) = E(secretKeyShare_1) + E(secretKeyShare_2) + ... + E(secretKeyShare_n)
// publicKey = publicKeyShare_1 + publicKeyShare_2 + ... + publicKeyShare_n
//
// to create a `ThresholdEcdsaKey`.
type ecdsaKeyShare struct {
	secretKeyShare *big.Int
	publicKeyShare *curve.Point
}
