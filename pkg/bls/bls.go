package bls

import (
	"math/big"

	"github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/keep-network/keep-core/pkg/altbn128"
)

// AggregateG1Points aggregates array of G1 points into a single G1 point.
func AggregateG1Points(points []*bn256.G1) *bn256.G1 {
	result := new(bn256.G1)
	for _, point := range points {
		result.Add(result, point)
	}
	return result
}

// AggregateG2Points aggregates array of G2 points into a single G2 point.
func AggregateG2Points(points []*bn256.G2) *bn256.G2 {
	result := new(bn256.G2)
	for _, point := range points {
		result.Add(result, point)
	}
	return result
}

// Sign creates a point on a curve G1 by hashing and signing provided message
// using the private key `k`.
func Sign(k *big.Int, msg []byte) *bn256.G1 {
	return new(bn256.G1).ScalarMult(altbn128.G1HashToPoint(msg), k)
}

// Verify performs the pairing operation to check if the signature is correct
// for the provided message and the corresponding public key.
func Verify(pub *bn256.G2, msg []byte, sig *bn256.G1) bool {

	// Generator point of G2 group.
	p2 := new(bn256.G2).ScalarBaseMult(big.NewInt(1))

	a := []*bn256.G1{new(bn256.G1).Neg(sig), altbn128.G1HashToPoint(msg)}
	b := []*bn256.G2{p2, pub}

	return bn256.PairingCheck(a, b)
}
