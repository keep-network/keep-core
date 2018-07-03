package pedersen

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
)

// Generates commitment based on [TC]

// [TC] https://github.com/keep-network/keep-core/blob/master/docs/cryptography/trapdoor-commitments.adoc

// cardinality - cardinality of groups G1, G2, GT. `q` in [TC]
var cardinality = bn256.Order

// GenerateCommitment generates a Commitment for passed `secret` message.
// Returns:
// `r` - Decommitment Key.
// `pubKey` - Public Key
// `h` - Random Point.
// `commitment` - Calculated commitment
func GenerateCommitment(secret *[]byte) (*big.Int, *bn256.G2, *bn256.G2, *bn256.G2, error) {
	// Generate random private and public keys.
	// [TC]: `privKey = (randomFromZ[0, q - 1])`
	// [TC]: `pubKey = g * privKey`
	_, pubKey, err := bn256.RandomG2(rand.Reader)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// Generate decommitment key.
	// `r` - decommitment key; used to commitment validation.
	// [TC]: `r = (randomFromZ[0, q - 1])`
	r, _, err := bn256.RandomG1(rand.Reader)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// Generate random point.
	// `h` - randomPoint; randomly selected point on the curve.
	_, h, err := bn256.RandomG2(rand.Reader)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// Calculate `secret`'s hash and it's `digest`.
	// [TC]: `digest = sha256(secret) mod q`
	hash := hash256BigInt(secret)
	digest := new(big.Int).Mod(hash, cardinality)

	// Calculate `he`
	// [TC]: `he = h + g * privKey`
	he := new(bn256.G2).Add(h, pubKey)

	// Calculate `commitment`
	// [TC]: `commitment = g * digest + he * r`
	commitment := new(bn256.G2).Add(new(bn256.G2).ScalarBaseMult(digest), new(bn256.G2).ScalarMult(he, r))

	// [TC]: `return (r, pubKey, h, commitment)`
	return r, pubKey, h, commitment, nil
}

// ValidateCommitment validates received commitment against revealed secret.
func ValidateCommitment(secret *[]byte, r *big.Int, commitment *bn256.G2, pubKey *bn256.G2, randomPoint *bn256.G2) (bool, error) {
	// Hash `secret` and calculate `digest`.
	// [TC]: `digest = sha256(secret) mod q`
	hash := hash256BigInt(secret)
	digest := new(big.Int).Mod(hash, cardinality)

	// Calculate `a`
	// [TC]: `a = g * r`
	a := new(bn256.G1).ScalarBaseMult(r)

	// Calculate `b`
	// [TC]: `b = h + g * privKey`
	b := new(bn256.G2).Add(randomPoint, pubKey)

	// Calculate `c`
	// [TC]: `c = commitment - g * digest)`
	c := new(bn256.G2).Add(commitment, new(bn256.G2).Neg(new(bn256.G2).ScalarBaseMult(digest)))

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
