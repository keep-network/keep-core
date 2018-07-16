// Package commitment generates and validates trapdoor commitments using
// bn256 elliptic curve pairings as described in docs/cryptography/trapdoor-commitments.adoc
package commitment

import (
	"crypto/rand"
	"crypto/sha256"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
)

// Secret is what we have committed to. Secret is usually kept private until
// the receiving party receives the message and wants to validate the commitment.
// Along with the message we have committed to, Secret contains
// a decommitment key r needed for the commitment verification.
type Secret struct {
	// Secret message
	message []byte
	// Decommitment key; used to commitment validation.
	r *big.Int
}

// TrapdoorCommitment is produced for each message we have committed to.
// It is usually revealed to the receiver immediately after it has been produced.
// Commitment lets to verify if the message revealed later by the sending party
// is really what that party has committed to.
// However, the Commitment itself is not enough for a verification.
// In order to perform verification, the interested party must receive
// the Secret too.
// Usually the process happens in two phases:
// first Commitment is evaluated and sent to receiver and then, after some time,
// secret is revealed and the receiver can check the secret message against
// the Commitment received earlier.
type TrapdoorCommitment struct {
	// Public key for a specific trapdoor commitment.
	pubKey *big.Int
	// Master trapdoor public key for the commitment family.
	h *bn256.G2
	// Calculated trapdoor commitment.
	commitment *bn256.G2
}

// GenerateCommitment generates a Commitment for passed `secret` message.
func GenerateCommitment(secret []byte) (*TrapdoorCommitment, *Secret, error) {
	// Generate a random public key.
	pubKey, _, err := bn256.RandomG2(rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	// Generate a decommitment key.
	r, _, err := bn256.RandomG1(rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	// Generate random point.
	_, h, err := bn256.RandomG2(rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	hash := sha256Sum(secret)
	digest := new(big.Int).Mod(hash, bn256.Order)

	// he = h + g * pubKey
	he := new(bn256.G2).Add(h, new(bn256.G2).ScalarBaseMult(pubKey))

	// commitment = g * digest + he * r
	commitment := new(bn256.G2).Add(new(bn256.G2).ScalarBaseMult(digest), new(bn256.G2).ScalarMult(he, r))

	return &TrapdoorCommitment{
			pubKey:     pubKey,
			h:          h,
			commitment: commitment,
		},
		&Secret{
			message: secret,
			r:       r,
		},
		nil
}

// ValidateCommitment validates received commitment against revealed secret.
func (tc *TrapdoorCommitment) ValidateCommitment(secret *Secret) bool {
	hash := sha256Sum(secret.message)
	digest := new(big.Int).Mod(hash, bn256.Order)

	// a = g * r
	a := new(bn256.G1).ScalarBaseMult(secret.r)

	// b = h + g * pubKey
	b := new(bn256.G2).Add(tc.h, new(bn256.G2).ScalarBaseMult(tc.pubKey))

	// c = commitment - g * digest
	c := new(bn256.G2).Add(tc.commitment, new(bn256.G2).Neg(new(bn256.G2).ScalarBaseMult(digest)))

	// Get base point `g`
	g := new(bn256.G1).ScalarBaseMult(big.NewInt(1))

	if bn256.Pair(a, b).String() != bn256.Pair(g, c).String() {
		return false
	}
	return true
}

// sha256Sum calculates sha256 hash for the passed `secret`
// and converts it to `big.Int`.
func sha256Sum(secret []byte) *big.Int {
	hash := sha256.Sum256(secret)

	return new(big.Int).SetBytes(hash[:])
}
