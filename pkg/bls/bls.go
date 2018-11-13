package bls

import (
	"math/big"

	"github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/keep-network/keep-core/pkg/altbn128"
)

// SecretKeyShare represents secret key share and its index.
type SecretKeyShare struct {
	I int      // Index of secret key share
	V *big.Int // Value of secret key share
}

// PublicKeyShare represents public key share and its index.
type PublicKeyShare struct {
	I int       // Index of public key share
	V *bn256.G2 // Value of public key share
}

// SignatureShare represents signature share and its index.
type SignatureShare struct {
	I int       // Index of signature share
	V *bn256.G1 // Value of signature share
}

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
func Recover(shares []*SignatureShare, threshold int) *bn256.G1 {

	// x holds id's of the threshold amount of shares required to recover signature.
	var x []*big.Int

	for _, s := range shares {
		if s == nil {
			continue
		}
		x = append(x, big.NewInt(1+int64(s.I)))
		if len(x) == threshold {
			break
		}
	}

	result := new(bn256.G1)

	for i, xi := range x {

		// Prepare numerator and denominator as part of Lagrange interpolation.
		num := big.NewInt(1)
		den := big.NewInt(1)

		for j, xj := range x {
			if i == j {
				continue
			}
			num = new(big.Int).Mod(new(big.Int).Mul(num, xj), bn256.Order)
			den = new(big.Int).Mod(new(big.Int).Mul(den, new(big.Int).Sub(xj, xi)), bn256.Order)
		}

		// Perform modular division of num by den.
		modInv := new(big.Int).ModInverse(den, bn256.Order)
		div := new(big.Int).Mod(new(big.Int).Mul(num, modInv), bn256.Order)

		result.Add(result, new(bn256.G1).ScalarMult(shares[i].V, div))
	}

	return result
}

// GetSecretKeyShare computes private share by evaluating a polynomial with
// coefficients taken from masterSecretKey. This is based on Shamir's Secret
// Sharing scheme and 'i' represents participant enumeration.
func GetSecretKeyShare(masterSecretKey []*big.Int, i int) *SecretKeyShare {

	xi := big.NewInt(int64(1 + i))
	share := big.NewInt(0)

	for j := len(masterSecretKey) - 1; j >= 0; j-- {
		share = new(big.Int).Mul(share, xi)
		share = new(big.Int).Add(share, masterSecretKey[j])
	}

	return &SecretKeyShare{i, share}
}

func (s SecretKeyShare) publicKeyShare() *PublicKeyShare {
	return &PublicKeyShare{s.I, new(bn256.G2).ScalarBaseMult(s.V)}
}

// RecoverPublicKey reconstructs public key from a threshold number of
// public key shares using Lagrange interpolation.
func RecoverPublicKey(shares []*PublicKeyShare, threshold int) *bn256.G2 {

	var x []*big.Int

	for _, s := range shares {
		if s == nil {
			continue
		}
		x = append(x, big.NewInt(1+int64(s.I)))
		if len(x) == threshold {
			break
		}
	}

	result := new(bn256.G2)

	for i, xi := range x {

		// Prepare numerator and denominator as part of Lagrange interpolation.
		num := big.NewInt(1)
		den := big.NewInt(1)

		for j, xj := range x {
			if i == j {
				continue
			}
			num = new(big.Int).Mod(new(big.Int).Mul(num, xj), bn256.Order)
			den = new(big.Int).Mod(new(big.Int).Mul(den, new(big.Int).Sub(xj, xi)), bn256.Order)
		}

		// Perform modular division of num by den.
		modInv := new(big.Int).ModInverse(den, bn256.Order)
		div := new(big.Int).Mod(new(big.Int).Mul(num, modInv), bn256.Order)

		result.Add(result, new(bn256.G2).ScalarMult(shares[i].V, div))
	}

	return result
}
