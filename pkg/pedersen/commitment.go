package pedersen

import (
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"math/big"
)

// Generates commitment based on [TC]

// TODO Update link when the file is in master
// [TC] https://github.com/keep-network/keep-core/blob/971ee24079e49385b4c957282770e1261a7bd74e/docs/cryptography/trapdoor-commitments.adoc#elliptic-curve-vss

var (
	// Curve - an Elliptic Curve. Points of the curve forms group `G` in [TC]
	Curve = elliptic.P256()
	// cardinality - cardinality of the group (N). `q` in [TC]
	cardinality = Curve.Params().N
	// randomPoint - randomly selected point on the curve. `randomPoint` in [TC]
	randomPoint *Point
	// basePoint - elliptic curve base point. `basePoint` in [TC]
	basePoint = &Point{x: Curve.Params().Gx, y: Curve.Params().Gy}
)

// Point is a structure of `x` and `y` coordinates
type Point struct {
	x, y *big.Int
}

// GenerateCommitment generates a Commitment for passed message.
// Returns:
// decommitmentKey
// publicKey
// commitment
// error
func GenerateCommitment(msg *[]byte) (*big.Int, *big.Int, []byte, error) {
	// Generate Randoms - another way of random generation is used in elliptic.GenerateKey
	// [TC]: `pkey = (randomFromZn[0, q - 1])`
	publicKey, err := RandomBigInt()
	if err != nil {
		return nil, nil, nil, err
	}

	// decommitmentKey - used to commitment validation. `r` in [TC]
	// [TC]: `r = (randomFromZn[0, q - 1])`
	decommitmentKey, err := RandomBigInt()
	if err != nil {
		return nil, nil, nil, err
	}

	// [TC]: `digest = sha256(secret) mod q`
	hash := sha256.New()
	_, err = hash.Write(*msg)
	if err != nil {
		return nil, nil, nil, err
	}
	hashBigInt := new(big.Int).SetBytes(hash.Sum(nil))

	digest := new(big.Int).Mod(hashBigInt, cardinality)

	// [TC]: `he = h + g * pkey`
	he := curveAdd(randomPoint, curveBaseMult(publicKey))

	// [TC]: `commitment = g * digest + he * r`
	commitmentPoint := curveAdd(curveBaseMult(digest), curveMult(he, decommitmentKey))
	commitment := elliptic.Marshal(Curve, commitmentPoint.x, commitmentPoint.y)

	// [TC]: `return (r, pkey, commitment)`
	return decommitmentKey, publicKey, commitment, nil
}

// ValidateCommitment TODO
func ValidateCommitment(commitment *[]byte, r *big.Int) bool {
	return false
}

// curveBaseMult returns result of `k` multiplications of the base point of the `Curve`.
func curveBaseMult(k *big.Int) *Point {
	var result *Point
	result.x, result.y = Curve.ScalarBaseMult(k.Bytes())
	return result
}

// curveMult returns result of `k` multiplications of a point on the `Curve`.
func curveMult(point1 *Point, k *big.Int) *Point {
	var result *Point
	result.x, result.y = Curve.ScalarMult(point1.x, point1.y, k.Bytes())
	return result
}

// curveAdd returns result of addition of two points on the `Curve`.
func curveAdd(point1, point2 *Point) *Point {
	var result *Point
	result.x, result.y = Curve.Add(point1.x, point1.y, point2.x, point2.y)
	return result
}

// RandomBigInt returns random `big.Int` value from `[0,q)`.
// It is required by [TC] to return a value from `[0, q-1]`, although
// `rand.Int` function operates in a half-open interval `[0,q)`.
func RandomBigInt() (*big.Int, error) {
	return rand.Int(rand.Reader, cardinality)
}
