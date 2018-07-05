package commitment

import (
	"crypto/rand"
	"crypto/sha256"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
)

// Generates and validates commitments based on [TC]
// Utilizes group of points on bn256 elliptic curve for calculations.

// [TC] https://github.com/keep-network/keep-core/blob/master/docs/cryptography/trapdoor-commitments.adoc

// TrapdoorSecret - parameters of commitment generation process
type TrapdoorSecret struct {
	// Secret message
	secret *[]byte
	// Decommitment key; used to commitment validation.
	r *big.Int
}

// TrapdoorCommitment - parameters which can be revealed
type TrapdoorCommitment struct {
	// Public key for a specific trapdoor commitment.
	pubKey *big.Int
	// Master public key for the trapdoor commitment family.
	h *bn256.G2
	// Calculated trapdoor commitment.
	commitment *bn256.G2
}

// GenerateCommitment generates a Commitment for passed `secret` message.
// Returns:
// `commitmentParams` - Commitment generation process parameters
// `error` - If generation failed
func GenerateCommitment(secret *[]byte) (*TrapdoorCommitment, *TrapdoorSecret, error) {
	// Generate random public key.
	// pubKey = (randomFromZ[0, q - 1])
	pubKey, _, err := bn256.RandomG2(rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	// Generate decommitment key.
	// r = (randomFromZ[0, q - 1])
	r, _, err := bn256.RandomG1(rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	// Generate random point.
	_, h, err := bn256.RandomG2(rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	// Calculate `secret`'s hash and it's `digest`.
	// digest = sha256(secret) mod q
	hash := hash256BigInt(secret)
	digest := new(big.Int).Mod(hash, bn256.Order)

	// he = h + g * privKey
	he := new(bn256.G2).Add(h, new(bn256.G2).ScalarBaseMult(pubKey))

	// commitment = g * digest + he * r
	commitment := new(bn256.G2).Add(new(bn256.G2).ScalarBaseMult(digest), new(bn256.G2).ScalarMult(he, r))

	return &TrapdoorCommitment{
			pubKey:     pubKey,
			h:          h,
			commitment: commitment,
		}, &TrapdoorSecret{
			secret: secret,
			r:      r,
		}, nil
}

// ValidateCommitment validates received commitment against revealed secret.
func ValidateCommitment(commitment *TrapdoorCommitment, secret *TrapdoorSecret) bool {
	// Calculate `secret`'s hash and it's `digest`.
	// digest = sha256(secret) mod q
	hash := hash256BigInt(secret.secret)
	digest := new(big.Int).Mod(hash, bn256.Order)

	// a = g * r
	a := new(bn256.G1).ScalarBaseMult(secret.r)

	// b = h + g * privKey
	b := new(bn256.G2).Add(commitment.h, new(bn256.G2).ScalarBaseMult(commitment.pubKey))

	// c = commitment - g * digest
	c := new(bn256.G2).Add(commitment.commitment, new(bn256.G2).Neg(new(bn256.G2).ScalarBaseMult(digest)))

	// Get base point `g`
	g := new(bn256.G1).ScalarBaseMult(big.NewInt(1))

	// Compare pairings
	// pairing(a, b) == pairing(g, c)
	if bn256.Pair(a, b).String() != bn256.Pair(g, c).String() {
		return false
	}
	return true
}

// hash256BigInt - calculates 256-bit hash for passed `secret` and converts it
// to `big.Int`.
func hash256BigInt(secret *[]byte) *big.Int {
	hash := sha256.Sum256(*secret)

	return new(big.Int).SetBytes(hash[:])
}
