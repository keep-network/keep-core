package pedersen

import (
	"crypto/rand"
	"crypto/sha256"
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

// ValidateCommitment TODO
func ValidateCommitment(commitment *[]byte, r *big.Int) bool {
	return false
}

// hash256BigInt - calculates 256-bit hash for passed `secret` and converts it
// to `big.Int`.
func hash256BigInt(secret *[]byte) *big.Int {
	hash := sha256.Sum256(*secret)
	hashBigInt := new(big.Int).SetBytes(hash[:])

	return hashBigInt
}
