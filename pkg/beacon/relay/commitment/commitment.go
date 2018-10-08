// Package commitment implements a commitment scheme described by
// Torben Pryds Pedersen in the referenced [Ped] paper.
//
// [Ped] Pedersen T.P. (1992) Non-Interactive and Information-Theoretic Secure
// Verifiable Secret Sharing. In: Feigenbaum J. (eds) Advances in Cryptology —
// CRYPTO ’91. CRYPTO 1991. Lecture Notes in Computer Science, vol 576. Springer,
// Berlin, Heidelberg
package commitment

import (
	"crypto/rand"
	"math/big"
	"time"

	"github.com/keep-network/paillier"
)

// Parameters specific to the scheme
type Parameters struct {
	p, q, g, h *big.Int
}

// Commitment is produced for each message we have committed to.
type Commitment struct {
	commitment *big.Int
}

// DecommitmentKey allows to open a commitment and verify if the value is what
// we have really committed to.
type DecommitmentKey struct {
	r *big.Int
}

// Generate evaluates a commitment and a decommitment key with specific master
// public key for the secret messages provided as an argument.
func Generate(parameters *Parameters, secret []byte) (*Commitment, *DecommitmentKey, error) {
	return nil, nil, nil
}

// Verify checks the received commitment against the revealed secret message.
func (c *Commitment) Verify(decommitmentKey *DecommitmentKey, secret []byte) bool {
	return false
}

func generateSafePrimes() (*big.Int, *big.Int, error) {
	concurrencyLevel := 4
	timeout := 120 * time.Second
	safePrimeBitLength := 512

	return paillier.GenerateSafePrime(safePrimeBitLength, concurrencyLevel, timeout, rand.Reader)
}
