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
//     [FO 97]: Fujisaki E., Okamoto T. (1997) Statistical zero knowledge
//          protocols to prove modular polynomial relations. In: Kaliski B.S.
//          (eds) Advances in Cryptology — CRYPTO '97. CRYPTO 1997. Lecture
//          Notes in Computer Science, vol 1294. Springer, Berlin, Heidelberg
type PublicParameters struct {
	// Paillier modulus used for generating T-ECDSA key and signing.
	N *big.Int

	// Paillier base for plaintext message used during encryption.
	// We always set it to N+1 to be compatible with `keep-network/paillier`.
	G *big.Int

	// Auxilliary RSA modulus which is the product of two safe primes.
	// It's uniquely generated for each new instance of `PublicParameters`.
	NTilde *big.Int

	// Two numbers from `(0, Z_NTilde*)` used to construct range
	// commitments. They are uniquely generated for each new instance of
	// `PublicParameters` according to the algorithm presented in [FO 97],
	// section 3.1.
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
	//TODO: make G a function
	G := new(big.Int).Add(paillierModulus, big.NewInt(1))

	pTilde, pTildePrime, err := paillier.GenerateSafePrimes(
		safePrimeBitLength, rand.Reader,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"could not generate ZKP PublicParameters [%v]",
			err,
		)
	}

	qTilde, qTildePrime, err := paillier.GenerateSafePrimes(
		safePrimeBitLength, rand.Reader,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"could not generate ZKP PublicParameters [%v]",
			err,
		)
	}

	NTilde := new(big.Int).Mul(pTilde, qTilde)
	h2, err := randomFromMultiplicativeGroup(rand.Reader, NTilde)

	NTildePrime := new(big.Int).Mul(pTildePrime, qTildePrime)
	x, err := rand.Int(rand.Reader, NTildePrime)
	if err != nil {
		return nil, fmt.Errorf(
			"could not generate ZKP PublicParameters [%v]",
			err,
		)
	}

	h1 := new(big.Int).Exp(h2, x, NTilde)

	return &PublicParameters{
		N:      paillierModulus,
		G:      G,
		NTilde: NTilde,
		h1:     h1,
		h2:     h2,
		q:      curve.Params().N,
		curve:  curve,
	}, nil
}
