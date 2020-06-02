package bls

import (
	"errors"
	"fmt"
	"github.com/ipfs/go-log"
	"math/big"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/keep-network/keep-core/pkg/altbn128"
)

var logger = log.Logger("keep-entry")

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
// using the provided secret key.
func Sign(secretKey *big.Int, message []byte) *bn256.G1 {
	return SignG1(secretKey, altbn128.G1HashToPoint(message))
}

// SignG1 creates a point on a curve G1 by signing the provided
// G1 point message using the provided secret key.
func SignG1(secretKey *big.Int, message *bn256.G1) *bn256.G1 {
	return new(bn256.G1).ScalarMult(message, secretKey)
}

// Verify performs the pairing operation to check if the signature is correct
// for the provided message and the corresponding public key.
func Verify(publicKey *bn256.G2, message []byte, signature *bn256.G1) bool {
	return VerifyG1(publicKey, altbn128.G1HashToPoint(message), signature)
}

// VerifyG1 performs the pairing operation to check if the signature is correct
// for the provided G1 point message and the corresponding public key.
func VerifyG1(publicKey *bn256.G2, message *bn256.G1, signature *bn256.G1) bool {
	// Generator point of G2 group.
	p2 := new(bn256.G2).ScalarBaseMult(big.NewInt(1))

	a := []*bn256.G1{new(bn256.G1).Neg(signature), message}
	b := []*bn256.G2{p2, publicKey}

	return bn256.PairingCheck(a, b)
}

// RecoverSignature reconstructs the full BLS signature from a threshold number of
// signature shares using Lagrange interpolation.
func RecoverSignature(shares []*SignatureShare, threshold int) (*bn256.G1, error) {

	// Indexes of participants that have valid shares.
	var validParticipants []*big.Int

	// Get sufficient number of participants with valid shares.
	for _, s := range shares {
		if len(validParticipants) == threshold {
			break
		}
		if s == nil || s.V == nil || s.I < 0 {
			continue
		}
		logger.Debugf(
			"recovering signature share index [%v] with the value [%v]", 
			s.I,
			s.V,
		)
		validParticipants = append(validParticipants, big.NewInt(int64(s.I)))
	}

	if len(validParticipants) < threshold {
		return nil, fmt.Errorf(
			"not enough shares to reconstruct signature: has [%v] shares, threshold is [%v]",
			len(validParticipants),
			threshold,
		)
	}

	result := new(bn256.G1)
	for i := range validParticipants {
		basis := lagrangeBasis(i, validParticipants)
		result.Add(result, new(bn256.G1).ScalarMult(shares[i].V, basis))
	}

	return result, nil
}

// GetSecretKeyShare computes secret share by evaluating a polynomial with
// coefficients taken from masterSecretKey. This is based on Shamir's Secret
// Sharing scheme and 'i' represents participant index.
func GetSecretKeyShare(masterSecretKey []*big.Int, i int) *SecretKeyShare {
	xi := big.NewInt(int64(i))
	share := big.NewInt(0)

	for j := 0; j < len(masterSecretKey); j++ {
		share = new(big.Int).Mul(share, xi)
		share = new(big.Int).Add(share, masterSecretKey[len(masterSecretKey)-1-j])
	}

	return &SecretKeyShare{i, share}
}

// PublicKeyShare returns public key share from the current secret key share.
func (s *SecretKeyShare) PublicKeyShare() *PublicKeyShare {
	return &PublicKeyShare{s.I, new(bn256.G2).ScalarBaseMult(s.V)}
}

// RecoverPublicKey reconstructs public key from a threshold number of
// public key shares using Lagrange interpolation.
func RecoverPublicKey(shares []*PublicKeyShare, threshold int) (*bn256.G2, error) {

	// Indexes of participants that have valid shares.
	var validParticipants []*big.Int

	// Get sufficient number of participants with valid shares.
	for _, s := range shares {
		if s == nil || s.V == nil || s.I < 0 {
			continue
		}
		validParticipants = append(validParticipants, big.NewInt(int64(s.I)))
		if len(validParticipants) == threshold {
			break
		}
	}

	if len(validParticipants) < threshold {
		return nil, errors.New("not enough shares to reconstruct public key")
	}

	result := new(bn256.G2)

	for i := range validParticipants {
		basis := lagrangeBasis(i, validParticipants)

		result.Add(result, new(bn256.G2).ScalarMult(shares[i].V, basis))
	}

	return result, nil
}

func lagrangeBasis(i int, validParticipants []*big.Int) *big.Int {

	// Prepare numerator and denominator as part of Lagrange interpolation.
	num := big.NewInt(1)
	den := big.NewInt(1)

	for j, xj := range validParticipants {
		if i == j {
			continue
		}
		num = new(big.Int).Mod(new(big.Int).Mul(num, xj), bn256.Order)
		den = new(big.Int).Mod(new(big.Int).Mul(den, new(big.Int).Sub(xj, validParticipants[i])), bn256.Order)
	}

	// Perform modular division of num by den.
	modInv := new(big.Int).ModInverse(den, bn256.Order)
	result := new(big.Int).Mod(new(big.Int).Mul(num, modInv), bn256.Order)

	return result
}
