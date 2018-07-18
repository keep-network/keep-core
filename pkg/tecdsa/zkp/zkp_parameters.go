package zkp

import (
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/keep-network/keep-core/pkg/tecdsa"

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
	pTilde, pTildePrime, err := paillier.GenerateSafePrimes(
		safePrimeBitLength, rand.Reader,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"could not generate ZKP PublicParameters [%v]",
			err,
		)
	}

	qTilde, qTildePrime := big.NewInt(0), big.NewInt(0)
	for qTilde, qTildePrime, err = paillier.GenerateSafePrimes(
		safePrimeBitLength, rand.Reader,
	); err == nil && pTilde.Cmp(qTilde) == 0; {
	}
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

// QSix is an auxiliary function returning q^8 which is needed to generate
// and validate internal commitment parameters.
func (p *PublicParameters) QSix() *big.Int {
	return new(big.Int).Exp(p.q, big.NewInt(6), nil)
}

// QEight is an auxiliary function returning q^8 which is needed to generate
// and validate internal commitment parameters.
func (p *PublicParameters) QEight() *big.Int {
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

// QSixNTilde is an auxiliary function returning q^8 * NTilde which is needed
// to generate and validate internal commitment parameters.
func (p *PublicParameters) QSixNTilde() *big.Int {
	return new(big.Int).Mul(p.QSix(), p.NTilde)
}

// QEightNTilde is an auxiliary function returning q^8 * NTilde which is needed
// to generate and validate internal commitment parameters.
func (p *PublicParameters) QEightNTilde() *big.Int {
	return new(big.Int).Mul(p.QEight(), p.NTilde)
}

// CurveBasePoint returns base point of the Elliptic Curve used in `curve` reference
func (p *PublicParameters) CurveBasePoint() *tecdsa.CurvePoint {
	return tecdsa.NewCurvePoint(p.curve.Params().Gx, p.curve.Params().Gy)
}
