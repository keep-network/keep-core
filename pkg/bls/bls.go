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
// using the private key `privateKey`.
func Sign(privateKey *big.Int, message []byte) *bn256.G1 {
	return new(bn256.G1).ScalarMult(altbn128.G1HashToPoint(message), privateKey)
}

// Verify performs the pairing operation to check if the signature is correct
// for the provided message and the corresponding public key.
func Verify(publicKey *bn256.G2, message []byte, signature *bn256.G1) bool {

	// Generator point of G2 group.
	p2 := new(bn256.G2).ScalarBaseMult(big.NewInt(1))

	a := []*bn256.G1{new(bn256.G1).Neg(signature), altbn128.G1HashToPoint(message)}
	b := []*bn256.G2{p2, publicKey}

	return bn256.PairingCheck(a, b)
}

// Recover reconstructs the full BLS signature from a threshold number of
// signature shares using Lagrange interpolation.
func Recover(shares []*bn256.G1, threshold int) *bn256.G1 {

	// x holds id's of the threshold amount of shares required to recover signature.
	x := make(map[int]*big.Int)

	for i, s := range shares {
		if s == nil {
			continue
		}
		x[i] = big.NewInt(1 + int64(i))
		if len(x) == threshold {
			break
		}
	}

	result := new(bn256.G1)

	for i, xi := range x {

		a := big.NewInt(1)
		b := big.NewInt(1)

		for j, xj := range x {
			if i == j {
				continue
			}
			a = new(big.Int).Mod(new(big.Int).Mul(a, xj), bn256.Order)
			b = new(big.Int).Mod(new(big.Int).Mul(b, new(big.Int).Sub(xj, xi)), bn256.Order)
		}

		// Perform modular division of a by b.
		modInv := new(big.Int).ModInverse(b, bn256.Order)
		div := new(big.Int).Mod(new(big.Int).Mul(a, modInv), bn256.Order)

		result.Add(result, new(bn256.G1).ScalarMult(shares[i], div))
	}

	return result
}

// SecretKeyShare computes private share by evaluating a polynomial with
// coefficients taken from masterSecretKey. This is based on Shamir's Secret
// Sharing scheme and 'i' represents participant enumeration.
func SecretKeyShare(masterSecretKey []*big.Int, i int64) *big.Int {

	xi := big.NewInt(int64(1 + i))
	share := big.NewInt(0)

	for j := len(masterSecretKey) - 1; j >= 0; j-- {
		share = new(big.Int).Mul(share, xi)
		share = new(big.Int).Add(share, masterSecretKey[j])
	}

	return share
}
