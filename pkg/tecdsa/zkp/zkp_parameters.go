package zkp

import (
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/keep-network/paillier"
)

// PublicParameters for all Gennaro's Zero Knowledge Proofs used for T-ECDSA
// algorithm as described in [GGN 16] section 4.4.
//
// New PublicParameters should be generated for each key generation and signing
// process.
//
//
//     [GGN 16]: Gennaro R., Goldfeder S., Narayanan A. (2016) Threshold-Optimal
//          DSA/ECDSA Signatures and an Application to Bitcoin Wallet Security.
//          In: Manulis M., Sadeghi AR., Schneider S. (eds) Applied Cryptography
//          and Network Security. ACNS 2016. Lecture Notes in Computer Science,
//          vol 9696. Springer, Cham
type PublicParameters struct {
	// Paillier modulus used for generating T-ECDSA key and signing.
	N *big.Int

	// Square of paillier modulus
	NSquare *big.Int

	// Paillier base for plaintext message used during encryption.
	// We always set it to N+1 to be compatible with `keep-network/paillier`.
	G *big.Int

	// Auxilliary RSA modulus which is the product of two safe primes.
	// It's uniquely generated for each new instance of `PublicParameters`.
	NTilde *big.Int

	// Two random numbers from `(0, Z_NTilde)` used to construct range
	// commitments. They are uniquely generated for each new instance of
	// `PublicParameters`.
	h1 *big.Int
	h2 *big.Int

	// Cardinality of Elliptic Curve used for T-ECDSA.
	q *big.Int

	// Elliptic Curve reference, the same one which is used for key generation
	// and signing in T-ECDSA.
	curve elliptic.Curve
}

//TODO: Increase once safe prime generator performance is bosted
const safePrimeBitLength = 256

// GeneratePublicParameters creates a new instance of `PublicParameters` with
// all field set to appropriate values.
func GeneratePublicParameters(
	paillierModulus *big.Int,
	curve elliptic.Curve,
) (*PublicParameters, error) {
	G := new(big.Int).Add(paillierModulus, big.NewInt(1))

	pTilde, qTilde, err := paillier.GenerateSafePrimes(
		safePrimeBitLength, rand.Reader,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"could not generate ZKP PublicParameters [%v]",
			err,
		)
	}

	NTilde := new(big.Int).Mul(pTilde, qTilde)

	h1 := big.NewInt(0)
	for h1, _ = rand.Int(rand.Reader, NTilde); h1.Cmp(big.NewInt(0)) == 0; {
	}

	h2 := big.NewInt(0)
	for h2, _ = rand.Int(rand.Reader, NTilde); h2.Cmp(big.NewInt(0)) == 0; {
	}

	return &PublicParameters{
		N:       paillierModulus,
		NSquare: new(big.Int).Exp(paillierModulus, big.NewInt(2), nil),
		G:       G,
		NTilde:  NTilde,
		h1:      h1,
		h2:      h2,
		q:       curve.Params().N,
		curve:   curve,
	}, nil
}
