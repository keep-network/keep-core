package zkp

import (
	"crypto/rand"
	"fmt"
	"io"
	"math/big"

	"github.com/keep-network/paillier"
)

// DsaPaillierSecretKeyFactorRangeProof is an implementation of Gennaro's PI_1,i
// proof for the Paillier encryption scheme, as described in [GGN16], section 4.4.
//
// The proof is used in the first and second round of the T-ECDSA signing algorithm
// and operates on DSA secret key encrypted with an additively homomorphic
// encryption scheme.
//
// Because of the complexity of the proof, we use the same naming for values
// as in the paper in most cases. We do an exception for function parameters:
// - `c1` in the paper represents encrypted DSA secret key multiplied by a factor η,
// - `c2` in the paper represents encrypted DSA secret key,
// - `c3` represents encrypted factor η
//
// The proof states that:
// ∃ η ∈ [-q^3, g^3] such that
// D(c3) = η
// D(c1) = η * D(c2)
//
// In other words, for the Elliptic Curve of cardinality `q`, random integer
// `η ∈ Z_q`, an additively homomorphic encryption scheme `E`, `u`
// representing encrypted random η `u = E(η)` and the value `v` which is
// a multiplication of the encoded secret ECDSA key `E(x)` by `η`,
// `v = η (*) E(x) = E(ηx)`,
// exists such η ∈ [-q^3, g^3] that:
// D(u) = D(E(η)) = η
// D(v) = η * D(E(x)) = ηx
//
//
//     [GGN 16]: Gennaro R., Goldfeder S., Narayanan A. (2016) Threshold-Optimal
//          DSA/ECDSA Signatures and an Application to Bitcoin Wallet Security.
//          In: Manulis M., Sadeghi AR., Schneider S. (eds) Applied Cryptography
//          and Network Security. ACNS 2016. Lecture Notes in Computer Science,
//          vol 9696. Springer, Cham
type DsaPaillierSecretKeyFactorRangeProof struct {
	z *big.Int
	v *big.Int

	u1 *big.Int
	u2 *big.Int

	e *big.Int

	s1 *big.Int
	s2 *big.Int
	s3 *big.Int
}

// CommitDsaPaillierSecretKeyFactorRange generates
// `DsaPaillierSecretKeyFactorRangeProof ` for the specified DSA secret
// key and multiplication factor. It's required to use the same randomness
// `r` to generate this proof as the one used for Paillier encryption of
// `secretDsaKeyShare` into `encryptedSecretDsaKeyShare`.
func CommitDsaPaillierSecretKeyFactorRange(
	encryptedSecretDsaKeyMultiple *paillier.Cypher, // = c1 = E(ηx)
	encryptedSecretDsaKey *paillier.Cypher, // = c2 = E(x)
	encryptedFactor *paillier.Cypher, // = c3 = E(η)
	factor *big.Int, // = η
	r *big.Int,
	params *PublicParameters,
	random io.Reader,
) (*DsaPaillierSecretKeyFactorRangeProof, error) {
	alpha, err := rand.Int(random, params.QCube())
	if err != nil {
		return nil, fmt.Errorf("could not construct ZKP1i [%v]", err)
	}

	beta, err := randomFromMultiplicativeGroup(random, params.N)
	if err != nil {
		return nil, fmt.Errorf("could not construct ZKP1i [%v]", err)
	}

	rho, err := rand.Int(random, params.QNTilde())
	if err != nil {
		return nil, fmt.Errorf("could not construct ZKP1i [%v]", err)
	}

	gamma, err := rand.Int(random, params.QCubeNTilde())
	if err != nil {
		return nil, fmt.Errorf("could not construct ZKP1i [%v]", err)
	}

	// u_1 = ((h1)^η)*((h2)^ρ) mod N ̃
	u1 := new(big.Int).Mod(
		new(big.Int).Mul(
			new(big.Int).Exp(params.h1, factor, params.NTilde),
			new(big.Int).Exp(params.h2, rho, params.NTilde),
		),
		params.NTilde,
	)

	// u_2 = ((h1)^α)*((h2)^γ) mod N ̃
	u2 := new(big.Int).Mod(
		new(big.Int).Mul(
			new(big.Int).Exp(params.h1, alpha, params.NTilde),
			new(big.Int).Exp(params.h2, gamma, params.NTilde),
		),
		params.NTilde,
	)

	// z = ((Γ)^α)*((β)^N) mod N^2
	z := new(big.Int).Mod(
		new(big.Int).Mul(
			new(big.Int).Exp(params.G(), alpha, params.NSquare()),
			new(big.Int).Exp(beta, params.N, params.NSquare()),
		),
		params.NSquare(),
	)

	// v = (c2)^α mod N^2
	v := new(big.Int).Exp(encryptedSecretDsaKey.C, alpha, params.NSquare())

	// e = hash(c1, c2, c3, z, u1, u2, v)
	digest := sum256(
		encryptedSecretDsaKeyMultiple.C.Bytes(),
		encryptedSecretDsaKey.C.Bytes(),
		encryptedFactor.C.Bytes(),
		z.Bytes(),
		u1.Bytes(),
		u2.Bytes(),
		v.Bytes(),
	)
	e := new(big.Int).Abs(new(big.Int).SetBytes(digest[:]))

	// s1 = e*η+α
	s1 := new(big.Int).Add(
		new(big.Int).Mul(e, factor),
		alpha,
	)

	// s2 = (r^e)*β mod N
	s2 := new(big.Int).Mod(
		new(big.Int).Mul(
			new(big.Int).Exp(r, e, params.N),
			beta,
		),
		params.N,
	)

	// s3 = e*ρ+γ
	s3 := new(big.Int).Add(
		new(big.Int).Mul(e, rho),
		gamma,
	)

	return &DsaPaillierSecretKeyFactorRangeProof{z, v, u1, u2, e, s1, s2, s3}, nil
}

// Verify checks the `PI1` against the provided secret message and secret key
// shares.
// If they match values used to generate the proof, function returns `true`.
// Otherwise, `false` is returned.
func (zkp *DsaPaillierSecretKeyFactorRangeProof) Verify(
	encryptedSecretDsaKeyMultiple *paillier.Cypher,
	encryptedSecretDsaKey *paillier.Cypher,
	encryptedFactor *paillier.Cypher,
	params *PublicParameters,
) bool {
	if !zkp.allParametersInRange(params) {
		return false
	}

	z := evaluateVerificationZ(encryptedFactor, zkp.s1, zkp.s2, zkp.e, params)
	v := evaluateVerificationV(encryptedSecretDsaKeyMultiple, encryptedSecretDsaKey, zkp.s1, zkp.e, params)
	u1 := zkp.u1 // u1 was calculated by the prover
	u2 := evaluateVerificationU2(zkp.u1, zkp.s1, zkp.s3, zkp.e, params)

	// e = hash(c1,c2,c3,z,u1,u2,v)
	digest := sum256(
		encryptedSecretDsaKeyMultiple.C.Bytes(),
		encryptedSecretDsaKey.C.Bytes(),
		encryptedFactor.C.Bytes(),
		z.Bytes(),
		u1.Bytes(),
		u2.Bytes(),
		v.Bytes(),
	)
	e := new(big.Int).SetBytes(digest[:])

	return zkp.u2.Cmp(u2) == 0 &&
		zkp.v.Cmp(v) == 0 &&
		zkp.z.Cmp(z) == 0 &&
		zkp.e.Cmp(e) == 0
}

// Checks whether parameters are in the expected range.
// It's a preliminary step to check if proof is not corrupted.
func (zkp *DsaPaillierSecretKeyFactorRangeProof) allParametersInRange(params *PublicParameters) bool {
	zero := big.NewInt(0)

	return isInRange(zkp.z, zero, params.NSquare()) &&
		isInRange(zkp.v, zero, params.NSquare()) &&
		isInRange(zkp.u1, zero, params.NTilde) &&
		isInRange(zkp.u2, zero, params.NTilde) &&
		isInRange(zkp.s2, zero, params.N)
}

// evaluateVerificationZ computes z verification value and returns it for
// further comparison with the expected one, evaluated during the commitment
// phase.
//
// z = (Γ^s1)*(s2^N)*(c3^(−e)) mod N^2
func evaluateVerificationZ(encryptedFactor *paillier.Cypher, s1, s2, e *big.Int,
	params *PublicParameters,
) *big.Int {
	return new(big.Int).Mod(
		new(big.Int).Mul(
			new(big.Int).Mul(
				new(big.Int).Exp(params.G(), s1, params.NSquare()),
				new(big.Int).Exp(s2, params.N, params.NSquare()),
			),
			discreteExp(encryptedFactor.C, new(big.Int).Neg(e), params.NSquare()),
		),
		params.NSquare(),
	)
}

// evaluateVerificationV computes v verification value and returns it for
// further comparison with the expected one, evaluated during the commitment
// phase.
//
// v = (c2^s1)*(c1^(−e)) mod N^2
func evaluateVerificationV(encryptedSecretDsaKeyMultiple, encryptedSecretDsaKey *paillier.Cypher,
	s1, e *big.Int, params *PublicParameters,
) *big.Int {
	return new(big.Int).Mod(
		new(big.Int).Mul(
			new(big.Int).Exp(encryptedSecretDsaKey.C, s1, params.NSquare()),
			discreteExp(encryptedSecretDsaKeyMultiple.C, new(big.Int).Neg(e), params.NSquare()),
		),
		params.NSquare(),
	)
}

// evaluateVerificationU2 computes u2 verification value and returns it for
// further comparison with the expected one, evaluated during the commitment
// phase.
//
// u2 =(h1^s1)*(h2^s3)*(u1^(−e)) mod N ̃
func evaluateVerificationU2(u1, s1, s3, e *big.Int, params *PublicParameters) *big.Int {
	return new(big.Int).Mod(
		new(big.Int).Mul(
			new(big.Int).Mul(
				new(big.Int).Exp(params.h1, s1, params.NTilde),
				new(big.Int).Exp(params.h2, s3, params.NTilde),
			),
			discreteExp(u1, new(big.Int).Neg(e), params.NTilde),
		),
		params.NTilde,
	)
}
