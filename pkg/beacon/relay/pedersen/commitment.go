// Package pedersen implements a Commitment scheme which is a part of Verifiable
// Secret Sharing (VSS) scheme described by Torben Pryds Pedersen in the
// referenced [Ped91b] paper.
// It consists of VSS parameters structure and functions to calculate and verify
// a commitment to chosen value.
//
// Commitment scheme allows a party (Commiter) to commit to a chosen value while
// keeping the value hidden from the other party (Verifier).
// On verification stage Committer reveals the value along with a DecommitmentKey,
// so Verifier can confirm the revealed value matches the committed one.
//
// pedersen.NewVSS() initializes scheme with `g` and `h` values, which need to
// be randomly generated for each scheme execution.
// To stop an adversary Committer from changing the value they already committed
// to, the scheme requires that `log_g(h)` is unknown to the Committer.
//
// You may consult our documentation for more details:
// docs/cryptography/trapdoor-commitments.html#_pedersen_commitment
//
//     [Ped91b]: T. Pedersen. Non-interactive and information-theoretic secure
//         verifiable secret sharing. In: Advances in Cryptology — Crypto '91,
//         pages 129-140. LNCS No. 576.
//         https://www.cs.cornell.edu/courses/cs754/2001fa/129.PDF
//     [GJKR 99]: Gennaro R., Jarecki S., Krawczyk H., Rabin T. (1999) Secure
//         Distributed Key Generation for Discrete-Log Based Cryptosystems. In:
//         Stern J. (eds) Advances in Cryptology — EUROCRYPT ’99. EUROCRYPT 1999.
//         Lecture Notes in Computer Science, vol 1592. Springer, Berlin, Heidelberg
//         http://groups.csail.mit.edu/cis/pubs/stasio/vss.ps.gz
package pedersen

import (
	crand "crypto/rand"
	"fmt"
	"io"
	"math/big"

	"github.com/keep-network/keep-core/pkg/internal/byteutils"
)

// VSS scheme parameters
type VSS struct {
	// p, q are primes such that `p = 2q + 1`.
	p, q *big.Int

	// g and h are elements of a group of order q, and should be chosen such that
	// no one knows log_g(h).
	G, h *big.Int
}

// Commitment represents a single commitment to a single message. One is produced
// for each message we have committed to.
//
// It is usually shared with the verifier immediately after it has been produced
// and lets the recipient verify if the message revealed later by the committing
// party is really what that party has committed to.
//
// The commitment itself is not enough for a verification. In order to perform
// a verification, the interested party must receive the `DecommitmentKey`.
type Commitment struct {
	vss        *VSS
	commitment *big.Int
}

// DecommitmentKey represents the key that allows a recipient to open an
// already-received commitment and verify if the value is what the sender have
// really committed to.
type DecommitmentKey struct {
	t *big.Int
}

// NewVSS creates a new instance of VSS from the provided parameters. Input
// parameters are not validated.
//
// TODO: It's just a temporary solution. This function should be removed once
// we switch to elliptic curves.
func NewVSS(p, q, g, h *big.Int) *VSS {
	return &VSS{p: p, q: q, G: g, h: h}
}

// GenerateVSS generates parameters for a scheme execution.
//
// It has to be run by a verifier or a trusted party. Executing generation by
// commiter themself causes that binding property is not held. Commiter gets an
// ability to change the value they already committed to.
func GenerateVSS(rand io.Reader, p, q *big.Int) (*VSS, error) {
	if !p.ProbablyPrime(20) || !q.ProbablyPrime(20) {
		return nil, fmt.Errorf("p and q have to be primes")
	}

	// Check if `p = 2q + 1`.
	pForQ := new(big.Int).Add(new(big.Int).Mul(big.NewInt(2), q), big.NewInt(1))
	if p.Cmp(pForQ) != 0 {
		return nil, fmt.Errorf("incorrect p and q values")
	}

	// Generate random `g` in Z_p, such that Euler's criterion `g^q mod p = 1`
	// holds.
	var err error
	var g *big.Int
	for {
		g, err = randomFromZn(rand, big.NewInt(1), p) // randomZ(1, p - 1]
		if err != nil {
			return nil, fmt.Errorf("g generation failed [%v]", err)
		}
		if new(big.Int).Exp(g, q, p).Cmp(big.NewInt(1)) == 0 {
			break
		}
	}

	// h = (g ^ randomZ(1, q - 1]) % p
	var h *big.Int
	for {
		randomValue, err := randomFromZn(rand, big.NewInt(1), q) // randomZ(1, q - 1]
		if err != nil {
			return nil, fmt.Errorf("randomValue generation failed [%v]", err)
		}
		h = new(big.Int).Exp(g, randomValue, p)

		if h.Cmp(big.NewInt(1)) > 0 {
			break
		}
	}

	return NewVSS(p, q, g, h), nil
}

// CommitmentTo takes a secret message and a set of parameters and returns
// a commitment to that message and the associated decommitment key.
//
// First random `t` value is chosen as a Decommitment Key.
// Then commitment is calculated as `(g ^ s) * (h ^ t) mod q`, where digest
// is sha256 hash of the secret brought to big.Int.
func (vss *VSS) CommitmentTo(rand io.Reader, secret []byte) (*Commitment, *DecommitmentKey, error) {
	t, err := randomFromZn(rand, big.NewInt(1), vss.q) // t = randomZ(1, q - 1]
	if err != nil {
		return nil, nil, fmt.Errorf("t generation failed [%v]", err)
	}

	s := calculateDigest(secret, vss.q) // s = hash(m) % q
	commitment := vss.CalculateCommitment(s, t, vss.q)

	return &Commitment{vss, commitment}, &DecommitmentKey{t}, nil
}

// Verify checks the received commitment against the revealed secret message and
// decommitment key.
//
// It returns `true` if a commitment calculated for passed decommitment key and
// secret matches the commitment value received before. Otherwise it returns false.
func (c *Commitment) Verify(decommitmentKey *DecommitmentKey, secret []byte) bool {
	s := calculateDigest(secret, c.vss.q)
	expectedCommitment := c.vss.CalculateCommitment(s, decommitmentKey.t, c.vss.q)
	return expectedCommitment.Cmp(c.commitment) == 0
}

func calculateDigest(secret []byte, mod *big.Int) *big.Int {
	hash := byteutils.Sha256Sum(secret)
	digest := new(big.Int).Mod(hash, mod)
	return digest
}

// CalculateCommitment calculates a commitment with equation `(g ^ s) * (h ^ t) mod m`
// where:
// - `g`, `h` are scheme specific parameters passed in vss,
// - `s` is a message to which one is committing,
// - `t` is a decommitment key.
func (vss *VSS) CalculateCommitment(s, t, m *big.Int) *big.Int {
	return new(big.Int).Mod(
		new(big.Int).Mul(
			new(big.Int).Exp(vss.G, s, m),
			new(big.Int).Exp(vss.h, t, m),
		),
		m,
	)
}

// randomFromZn generates a random `big.Int` in a range (min, max).
func randomFromZn(rand io.Reader, min, max *big.Int) (*big.Int, error) {
	for {
		x, err := crand.Int(rand, max) // returns a value in [0, max)
		if err != nil {
			return nil, fmt.Errorf("failed to generate random number [%v]", err)
		}

		if x.Cmp(min) > 0 {
			return x, nil
		}
	}
}
