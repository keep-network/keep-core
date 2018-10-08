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
	"fmt"
	"math/big"
	"time"

	"github.com/keep-network/keep-core/pkg/internal/byteutils"
	"github.com/keep-network/paillier"
)

// Parameters specific to the scheme
type Parameters struct {
	p    *big.Int // Prime such that p = 2q + 1
	q    *big.Int // Sophie Germain prime
	g, h *big.Int // Elements of a subgroup of quadratic residues of order q
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

// GenerateParameters generates parameters for a scheme execution
func GenerateParameters() (*Parameters, error) {
	p, q, err := generateSafePrimes()
	if err != nil {
		return nil, fmt.Errorf("p,q generation failed [%s]", err)
	}

	randomG, err := randomFromZn(p)
	if err != nil {
		return nil, fmt.Errorf("g generation failed [%s]", err)
	}
	g := new(big.Int).Exp(randomG, big.NewInt(2), nil) // (randomZ(0, 2^p - 1]) ^2

	randomH, err := randomFromZn(p) // (randomZ(0, 2^p - 1]) ^2
	if err != nil {
		return nil, fmt.Errorf("h generation failed [%s]", err)
	}
	h := new(big.Int).Exp(randomH, big.NewInt(2), nil) // (randomZ(0, 2^p - 1]) ^2

	return &Parameters{p: p, q: q, g: g, h: h}, nil
}

// Generate evaluates a commitment and a decommitment key with specific master
// public key for the secret messages provided as an argument.
func Generate(parameters *Parameters, secret []byte) (*Commitment, *DecommitmentKey, error) {
	r, err := randomFromZn(parameters.q) // randomZ(0, 2^q - 1]
	if err != nil {
		return nil, nil, fmt.Errorf("r generation failed [%s]", err)
	}

	digest := calculateDigest(secret, parameters.q)

	commitment := calculateCommitment(parameters, digest, r)

	return &Commitment{commitment},
		&DecommitmentKey{r},
		nil
}

// Verify checks the received commitment against the revealed secret message.
func (c *Commitment) Verify(parameters *Parameters, decommitmentKey *DecommitmentKey, secret []byte) bool {
	digest := calculateDigest(secret, parameters.q)
	expectedCommitment := calculateCommitment(parameters, digest, decommitmentKey.r)
	return expectedCommitment.Cmp(c.commitment) == 0
}

func calculateDigest(secret []byte, mod *big.Int) *big.Int {
	hash := byteutils.Sha256Sum(secret)
	digest := new(big.Int).Mod(hash, mod)
	return digest
}

func calculateCommitment(parameters *Parameters, digest, r *big.Int) *big.Int {
	// ((g ^ digest) % p) * ((h ^ r) % p)
	return new(big.Int).Mul(
		new(big.Int).Exp(parameters.g, digest, parameters.p),
		new(big.Int).Exp(parameters.h, r, parameters.p),
	)
}

func generateSafePrimes() (*big.Int, *big.Int, error) {
	concurrencyLevel := 4
	timeout := 120 * time.Second
	safePrimeBitLength := 512

	return paillier.GenerateSafePrime(safePrimeBitLength, concurrencyLevel, timeout, rand.Reader)
}

// randomFromZn generates a random `big.Int` in a range (0, 2^n - 1]
// TODO check if this is what we really need for g,h and r
func randomFromZn(n *big.Int) (*big.Int, error) {
	x := big.NewInt(0)
	var err error
	// 2^n - 1
	max := new(big.Int).Sub(
		new(big.Int).Exp(
			big.NewInt(2),
			n,
			nil,
		),
		big.NewInt(1),
	)
	for x.Sign() == 0 {
		x, err = rand.Int(rand.Reader, max)
		if err != nil {
			return nil, fmt.Errorf("failed to generate random number [%s]", err)
		}
	}
	return x, nil
}
