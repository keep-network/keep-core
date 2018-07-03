package tecdsa

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
)

// Generates and validates commitments based on [TC]
// Utilizes group of points on bn256 elliptic curve for calculations.

// [TC] https://github.com/keep-network/keep-core/blob/master/docs/cryptography/trapdoor-commitments.adoc

// cardinality - number of points in each of groups G1 and G2. `q` in [TC]
var cardinality = bn256.Order

// commitmentParams - parameters of commitment generation process
type commitmentParams struct {
	publicParams
	// Secret message
	secret *[]byte
	// `r` - decommitment key; used to commitment validation.
	r *big.Int
}

// publicParams - parameters which can be revealed
type publicParams struct {
	// Random public key.
	pubKey *bn256.G2
	// Randomly selected point on the curve.
	h *bn256.G2
	// Calculated commitment.
	commitment *bn256.G2
}

// GenerateCommitment generates a Commitment for passed `secret` message.
// Returns:
// `commitmentParams` - Commitment generation process parameters
// `error` - If generation failed
func GenerateCommitment(secret *[]byte) (*commitmentParams, error) {
	// Generate random private and public keys.
	// [TC]: `privKey = (randomFromZ[0, q - 1])`
	// [TC]: `pubKey = g * privKey`
	_, pubKey, err := bn256.RandomG2(rand.Reader)
	if err != nil {
		return nil, err
	}

	// Generate decommitment key.
	// [TC]: `r = (randomFromZ[0, q - 1])`
	r, _, err := bn256.RandomG1(rand.Reader)
	if err != nil {
		return nil, err
	}

	// Generate random point.
	_, h, err := bn256.RandomG2(rand.Reader)
	if err != nil {
		return nil, err
	}

	// Calculate `secret`'s hash and it's `digest`.
	// [TC]: `digest = sha256(secret) mod q`
	hash := hash256BigInt(secret)
	cardinality := bn256.Order
	digest := new(big.Int).Mod(hash, cardinality)

	// Calculate `he`
	// [TC]: `he = h + g * privKey`
	he := new(bn256.G2).Add(h, pubKey)

	// Calculate `commitment`
	// [TC]: `commitment = g * digest + he * r`
	commitment := new(bn256.G2).Add(new(bn256.G2).ScalarBaseMult(digest), new(bn256.G2).ScalarMult(he, r))

	// [TC]: `return (r, pubKey, h, commitment)`
	return &commitmentParams{
		publicParams: publicParams{
			pubKey:     pubKey,
			h:          h,
			commitment: commitment,
		},
		secret: secret,
		r:      r,
	}, nil
}

// ValidateCommitment validates received commitment against revealed secret.
func ValidateCommitment(privateParams *commitmentParams) (bool, error) {
	// Hash `secret` and calculate `digest`.
	// [TC]: `digest = sha256(secret) mod q`
	hash := hash256BigInt(privateParams.secret)
	digest := new(big.Int).Mod(hash, cardinality)

	// Calculate `a`
	// [TC]: `a = g * r`
	a := new(bn256.G1).ScalarBaseMult(privateParams.r)

	// Calculate `b`
	// [TC]: `b = h + g * privKey`
	b := new(bn256.G2).Add(privateParams.h, privateParams.pubKey)

	// Calculate `c`
	// [TC]: `c = commitment - g * digest`
	c := new(bn256.G2).Add(privateParams.commitment, new(bn256.G2).Neg(new(bn256.G2).ScalarBaseMult(digest)))

	// Get base point `g`
	g := new(bn256.G1).ScalarBaseMult(big.NewInt(1))

	// Compare pairings
	// [TC]: pairing(a, b) == pairing(g, c)
	if bn256.Pair(a, b).String() != bn256.Pair(g, c).String() {
		return false, fmt.Errorf("pairings doesn't match")
	}
	return true, nil
}

// hash256BigInt - calculates 256-bit hash for passed `secret` and converts it
// to `big.Int`.
func hash256BigInt(secret *[]byte) *big.Int {
	hash := sha256.Sum256(*secret)
	hashBigInt := new(big.Int).SetBytes(hash[:])

	return hashBigInt
}
