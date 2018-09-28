package zkp

import (
	"crypto/rand"
	"fmt"
	"io"
	"math/big"

	"github.com/keep-network/paillier"
)

// EcdsaPaillierSecretKeyFactorRangeProof is an implementation of Gennaro's Π_1,i
// proof for the Paillier encryption scheme, as described in [GGN16], section 4.4.
//
// The proof is used in the first and second round of the T-ECDSA signing algorithm
// and operates on ECDSA secret key encrypted with an additively homomorphic
// encryption scheme.
//
// Because of the complexity of the proof, we use the same naming for values
// as in the paper in most cases. We do an exception for function parameters:
// - `c1` in the paper represents encrypted ECDSA secret key multiplied by a factor η,
// - `c2` in the paper represents encrypted ECDSA secret key,
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
type EcdsaPaillierSecretKeyFactorRangeProof struct {
	z *big.Int
	v *big.Int

	u1 *big.Int
	u2 *big.Int

	e *big.Int

	s1 *big.Int
	s2 *big.Int
	s3 *big.Int
}

// CommitEcdsaPaillierSecretKeyFactorRange generates
// `EcdsaPaillierSecretKeyFactorRangeProof` for the specified ECDSA secret
// key and multiplication factor. It's required to use the same randomness
// `r` to generate this proof as the one used for Paillier encryption of
// `secretEcdsaKeyShare` into `encryptedSecretEcdsaKeyShare`.
func CommitEcdsaPaillierSecretKeyFactorRange(
	secretEcdsaKeyMultiple *paillier.Cypher, // = c1 = E(ηx)
	secretEcdsaKey *paillier.Cypher, // = c2 = E(x)
	secretEcdsaKeyFactor *paillier.Cypher, // = c3 = E(η)
	factor *big.Int, // = η
	r *big.Int,
	params *PublicParameters,
	random io.Reader,
) (*EcdsaPaillierSecretKeyFactorRangeProof, error) {
	alpha, err := rand.Int(random, params.QCube())
	if err != nil {
		return nil, fmt.Errorf("could not construct the proof [%v]", err)
	}

	beta, err := paillier.GetRandomNumberInMultiplicativeGroup(params.N, random)
	if err != nil {
		return nil, fmt.Errorf("could not construct the proof [%v]", err)
	}

	rho, err := rand.Int(random, params.QNTilde())
	if err != nil {
		return nil, fmt.Errorf("could not construct the proof [%v]", err)
	}

	gamma, err := rand.Int(random, params.QCubeNTilde())
	if err != nil {
		return nil, fmt.Errorf("could not construct the proof [%v]", err)
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
	v := new(big.Int).Exp(secretEcdsaKey.C, alpha, params.NSquare())

	// e = hash(c1, c2, c3, z, u1, u2, v)
	digest := sum256(
		secretEcdsaKeyMultiple.C.Bytes(),
		secretEcdsaKey.C.Bytes(),
		secretEcdsaKeyFactor.C.Bytes(),
		z.Bytes(),
		u1.Bytes(),
		u2.Bytes(),
		v.Bytes(),
	)
	e := new(big.Int).SetBytes(digest[:])

	// s1 = e*η+α
	s1 := new(big.Int).Add(
		new(big.Int).Mul(e, factor),
		alpha,
	)

	// s2 = (r^e)*β mod N
	s2 := new(big.Int).Mod(
		new(big.Int).Mul(
			discreteExp(r, e, params.N),
			beta,
		),
		params.N,
	)

	// s3 = e*ρ+γ
	s3 := new(big.Int).Add(
		new(big.Int).Mul(e, rho),
		gamma,
	)

	return &EcdsaPaillierSecretKeyFactorRangeProof{z, v, u1, u2, e, s1, s2, s3}, nil
}

// Verify checks the `EcdsaPaillierSecretKeyFactorRangeProof` against the provided
// ECDSA secret key and the multiplication factor.
// If they match values used to generate the proof, function returns `true`.
// Otherwise, `false` is returned.
func (zkp *EcdsaPaillierSecretKeyFactorRangeProof) Verify(
	secretEcdsaKeyMultiple *paillier.Cypher, // = c1 = E(ηx)
	secretEcdsaKey *paillier.Cypher, // = c2 = E(x)
	secretEcdsaKeyFactor *paillier.Cypher, // = c3 = E(η)
	params *PublicParameters,
) bool {
	if !zkp.allParametersInRange(params) {
		return false
	}

	z := evaluateVerificationZ(secretEcdsaKeyFactor, zkp.s1, zkp.s2, zkp.e, params)
	v := evaluateVerificationV(secretEcdsaKeyMultiple, secretEcdsaKey, zkp.s1, zkp.e, params)
	u2 := evaluateVerificationU2(zkp.u1, zkp.s1, zkp.s3, zkp.e, params)

	// e = hash(c1,c2,c3,z,u1,u2,v)
	digest := sum256(
		secretEcdsaKeyMultiple.C.Bytes(),
		secretEcdsaKey.C.Bytes(),
		secretEcdsaKeyFactor.C.Bytes(),
		z.Bytes(),
		zkp.u1.Bytes(),
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
func (zkp *EcdsaPaillierSecretKeyFactorRangeProof) allParametersInRange(params *PublicParameters) bool {
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
// We want to verify whether z = (Γ^s1)*(s2^N)*(c3^(−e))
// is equal to z = Γ^α * β^N
// we evaluated in the commitment phase.
//
// Since:
// s1 = eη + α
// s2 = r^e * β
//
// and since we know from Paillier that:
// c3 = E(η) = Γ^η * r^N
//
// We can do:
// z = (Γ^s1)*(s2^N)*(c3^(−e)) =
// Γ^{eη + α} * r^eN * β^N * Γ^{-eη} * r^{-eN} =
// Γ^α * β^N
//
// which is exactly how z is evaluated during the commitment phase.
func evaluateVerificationZ(encryptedFactor *paillier.Cypher, s1, s2, e *big.Int,
	params *PublicParameters,
) *big.Int {
	return new(big.Int).Mod(
		new(big.Int).Mul(
			new(big.Int).Mul(
				discreteExp(params.G(), s1, params.NSquare()),
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
// We want to verify whether  v = (c2^s1)*(c1^(−e))
// is equal to v = (c2)^α
// we evaluated in the commitment phase.
//
// Since:
// s1 = eη + α
//
// and since we know from Paillier that:
// c1 = E(ηx) = Γ^{ηx} * r^{ηN}
// c2 = E(x) = Γ^x * r^N
//
// We can do:
// v = (c2^s1)*(c1^(−e)) =
// (c2)^{eη} * (c2)^α * (c1)^{-e} =
// (c2)^α * Γ^{eηx} * r^{eηN} * Γ^{-eηx} * r^{-eηN} =
// (c2)^α
//
// which is exactly how v is evaluated during the commitment phase.
func evaluateVerificationV(encryptedSecretDsaKeyMultiple, encryptedSecretDsaKey *paillier.Cypher,
	s1, e *big.Int, params *PublicParameters,
) *big.Int {
	return new(big.Int).Mod(
		new(big.Int).Mul(
			discreteExp(encryptedSecretDsaKey.C, s1, params.NSquare()),
			discreteExp(encryptedSecretDsaKeyMultiple.C, new(big.Int).Neg(e), params.NSquare()),
		),
		params.NSquare(),
	)
}

// evaluateVerificationU2 computes u2 verification value and returns it for
// further comparison with the expected one, evaluated during the commitment
// phase.
//
// We want to verify whether u2 = (h1)^{s1} * (h2)^{s3} * (u1)^-e
// is equal to u2 = (h1)^alpha * (h2)^γ
// we evaluated in the commitment phase.
//
// Since:
// s1 = eη + α
// s3 = e*ρ + γ
// u1 = (h1)^η * (h2)^ρ
//
// We can do:
// u2 = (h1)^{s1} * (h2)^{s3} * (u1)^-e =
// (h1)^{e^η + α} * h2^{e*ρ + γ} * [(h1)^η * (h2)^ρ]^-e =
// (h1)^{eη} * (h1)^{α} * (h2)^{e*ρ} * (h2)^{γ} * (h1)^{-eη} * (h2)^{-e * ρ} =
// (h1)^α * (h2)^γ
//
// which is exactly how u2 is evaluated during the commitment phase.
func evaluateVerificationU2(u1, s1, s3, e *big.Int, params *PublicParameters) *big.Int {
	return new(big.Int).Mod(
		new(big.Int).Mul(
			new(big.Int).Mul(
				discreteExp(params.h1, s1, params.NTilde),
				discreteExp(params.h2, s3, params.NTilde),
			),
			discreteExp(u1, new(big.Int).Neg(e), params.NTilde),
		),
		params.NTilde,
	)
}
