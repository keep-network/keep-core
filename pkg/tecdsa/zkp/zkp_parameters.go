package zkp

import (
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

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

	// Auxiliary RSA modulus which is the product of two safe primes.
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

// GeneratePublicParameters generates a new instance of `PublicParameters`.
// The auxiliary ZKP `NTilde` modulus will have the same bit length as the
// `paillierModulus` passed as a parameter.
func GeneratePublicParameters(
	paillierModulus *big.Int,
	curve elliptic.Curve,
) (*PublicParameters, error) {
	// Concurrency configuration for safe prime generator.
	safePrimeGenConcurrencyLevel := 4
	safePrimeGenTimeout := 120 * time.Second

	// Bit length of safe prime numbers used to generate NTilde.
	//
	// Bit length of `NTilde` will be the same as the bit length of
	// `paillierModulus`. This is not a protocol requirement but our
	// implementation choice - we use the same security guarantees for ZKPs
	// as they are used for the Paillier key.
	safePrimeBitLength := paillierModulus.BitLen() / 2

	pTilde, _, err := paillier.GenerateSafePrime(
		safePrimeBitLength,
		safePrimeGenConcurrencyLevel,
		safePrimeGenTimeout,
		rand.Reader,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"could not generate ZKP PublicParameters [%v]",
			err,
		)
	}

	qTilde := big.NewInt(0)
	for {
		qTilde, _, err = paillier.GenerateSafePrime(
			safePrimeBitLength,
			safePrimeGenConcurrencyLevel,
			safePrimeGenTimeout,
			rand.Reader,
		)
		if err != nil || pTilde.Cmp(qTilde) != 0 {
			break
		}
	}
	if err != nil {
		return nil, fmt.Errorf(
			"could not generate ZKP PublicParameters [%v]",
			err,
		)
	}

	return GeneratePublicParametersFromSafePrimes(
		paillierModulus, pTilde, qTilde, curve,
	)
}

// GeneratePublicParametersFromSafePrimes generates a new instance of
// `PublicParameters` from the provided safe primes of equal bit length.
func GeneratePublicParametersFromSafePrimes(
	paillierModulus *big.Int,
	pTilde *big.Int,
	qTilde *big.Int,
	curve elliptic.Curve,
) (*PublicParameters, error) {
	pTildePrime := new(big.Int).Div(
		new(big.Int).Sub(pTilde, big.NewInt(1)),
		big.NewInt(2),
	)

	qTildePrime := new(big.Int).Div(
		new(big.Int).Sub(qTilde, big.NewInt(1)),
		big.NewInt(2),
	)

	if pTilde.BitLen() != qTilde.BitLen() {
		return nil, fmt.Errorf(
			"safe primes must have the same bit length, got %d and %d bits",
			pTilde.BitLen(),
			qTilde.BitLen(),
		)
	}

	if pTilde.Cmp(qTilde) == 0 {
		return nil, fmt.Errorf(
			"safe primes must not be equal, got %v and %v", pTilde, qTilde,
		)
	}

	if !pTilde.ProbablyPrime(20) ||
		!qTilde.ProbablyPrime(20) ||
		!pTildePrime.ProbablyPrime(20) ||
		!qTildePrime.ProbablyPrime(20) {
		return nil, errors.New("both numbers must be safe primes")
	}

	NTilde := new(big.Int).Mul(pTilde, qTilde)
	h2, err := paillier.GetRandomNumberInMultiplicativeGroup(NTilde, rand.Reader)

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
		NTilde: NTilde,
		h1:     h1,
		h2:     h2,
		q:      curve.Params().N,
		curve:  curve,
	}, nil
}

// G is an auxiliary function returning Paillier base for plaintext message used
// during encryption. We always set it to N+1 to be compatible with
// `keep-network/paillier`.
// G is needed to generate and validate internal commitment parameters.
func (p *PublicParameters) G() *big.Int {
	return new(big.Int).Add(p.N, big.NewInt(1))
}

// NSquare is an auxiliary function returning N^2 which is needed to generate
// and validate internal commitment parameters.
func (p *PublicParameters) NSquare() *big.Int {
	return new(big.Int).Exp(p.N, big.NewInt(2), nil)
}

// QCube is an auxiliary function returning q^3 which is needed to generate
// and validate internal commitment parameters.
func (p *PublicParameters) QCube() *big.Int {
	return new(big.Int).Exp(p.q, big.NewInt(3), nil)
}

// QPow6 is an auxiliary function returning q^6 which is needed to generate
// and validate internal commitment parameters.
func (p *PublicParameters) QPow6() *big.Int {
	return new(big.Int).Exp(p.q, big.NewInt(6), nil)
}

// QPow8 is an auxiliary function returning q^8 which is needed to generate
// and validate internal commitment parameters.
func (p *PublicParameters) QPow8() *big.Int {
	return new(big.Int).Exp(p.q, big.NewInt(8), nil)
}

// QNTilde is an auxiliary function returning q * NTilde which is needed to
// generate and validate internal commitment parameters.
func (p *PublicParameters) QNTilde() *big.Int {
	return new(big.Int).Mul(p.q, p.NTilde)
}

// QCubeNTilde is an auxiliary function returning q^3 * NTilde which is needed
// to generate and validate internal commitment parameters.
func (p *PublicParameters) QCubeNTilde() *big.Int {
	return new(big.Int).Mul(p.QCube(), p.NTilde)
}

// QPow6NTilde is an auxiliary function returning q^6 * NTilde which is needed
// to generate and validate internal commitment parameters.
func (p *PublicParameters) QPow6NTilde() *big.Int {
	return new(big.Int).Mul(p.QPow6(), p.NTilde)
}

// QPow8NTilde is an auxiliary function returning q^8 * NTilde which is needed
// to generate and validate internal commitment parameters.
func (p *PublicParameters) QPow8NTilde() *big.Int {
	return new(big.Int).Mul(p.QPow8(), p.NTilde)
}
