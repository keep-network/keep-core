package zkp

import (
	"crypto/rand"
	"fmt"
	"io"
	"math/big"

	"github.com/keep-network/paillier"

	"github.com/keep-network/keep-core/pkg/tecdsa/curve"
)

// EcdsaSignatureFactorRangeProof is an implementation of Gennaro's Π_2,i
// proof for the Paillier encryption scheme, as described in [GGN16], section 4.4.
//
// The proof is used in the fourth round of the T-ECDSA signing algorithm and operates
// on a random elliptic curve point `r_i = g^{k_i}` (signature factor) and value
// `w` evaluated from `k_i`, random integer `c_i` and value `ρ` (secret key
// factor) evaluated in rounds one and two: `w_i = E(k_i * ρ + c_i * q)`.
//
// Values `w` and `r` evaluated as a sum of all shares `w_i` and `r_i` are used
// in the fifth round to produce an encrypted ECDSA signature for a message.
//
// Because of the complexity of the proof, we use the same naming for values
// as in the paper in most cases. We do an exception for function parameters:
// - `k_i` in the paper is the `signatureFactorSecret`
// - `r_i` in the paper is the `signatureFactorPublic`
// - `c_i` in the paper is the `signatureFactorMask`
// - `u` in the paper is the `secretKeyFactor` evaluated in round 2
// - `w_i` in the paper is the `signatureUnmask`
// - `r_c` in the paper is the `paillierR` which is a Paillier randomness r
//
// Few notes:
// - `k_i` is named `signatureFactorSecret` because it is used as a
//   random factor when we produce a signature `s = k^-1 (m + xr)` in
//   rounds 5 and 6. Also, since it's kept secret, it has the appropriate suffix,
// - `r_i` is named `signatureFactorPublic` because it's a revealed
//   value of `k` (`k` is kept hidden) to which signer may commit, `r_i = g^{k_i}`,
// - `c_i` is named `signatureFactorMask` because it's used in the
//   value `w_i` to mask the value of `k_i`: `w_i = E(k_i * ρ + c_i * q),
// - `w_i` is named `signatureUnmask` because the way how it's constructed
//   lets to unmask the value of `k_i` and eliminate random multiple `ρ` in
//   rounds 5 and 6 - please consult their documentation for more details.
//
// The proof states that:
// ∃ η1 ∈ [-q^3, q^3], η2 ∈ [-q^8, q^8] such that
// g^η1 = r
// D(w) = η1*D(u) + q*η2
//
// Using the same naming as in round 3 and 4 when the values used here are
// constructed, we have:
//
// η1 = k_i
// η2 = c_i
//
// and we prove that:
// D(w) = k_i * D(u) + c_i * q
//
// which is how `w` in round 3 is constructed:
// w = E(k_i * ρ + c_i * q)
//
// since u = E(ρ)
//
// The way how `w` is constructed lets to eliminate `ρ` factor in the fifth round
// of signing and `k` which is a sum of all `k_i` factors can be then used as
// a random factor for the ECDSA signature.
// Please consult the documentation of round 4 and 5 for more details.
//
//
//     [GGN 16]: Gennaro R., Goldfeder S., Narayanan A. (2016) Threshold-Optimal
//          DSA/ECDSA Signatures and an Application to Bitcoin Wallet Security.
//          In: Manulis M., Sadeghi AR., Schneider S. (eds) Applied Cryptography
//          and Network Security. ACNS 2016. Lecture Notes in Computer Science,
//          vol 9696. Springer, Cham
type EcdsaSignatureFactorRangeProof struct {
	z1 *big.Int
	z2 *big.Int

	u1 *curve.Point
	u2 *big.Int
	u3 *big.Int

	v1 *big.Int
	v2 *big.Int
	v3 *big.Int

	e *big.Int

	s1 *big.Int
	s2 *big.Int

	t1 *big.Int
	t2 *big.Int
	t3 *big.Int
}

// CommitEcdsaSignatureFactorRangeProof generates `EcdsaSignatureFactorRangeProof`
// for the specified parameters.
// It's required to use the same randomness `paillierR` to generate this proof as
// the one used for Paillier encryption of `secretKeyFactor`.
func CommitEcdsaSignatureFactorRangeProof(
	signatureFactorPublic *curve.Point, // r_i = g^{k_i}
	signatureUnmask *paillier.Cypher, // w = E(k * ρ + c_i * q)
	secretKeyFactor *paillier.Cypher, // u = E(ρ)
	signatureFactorSecret *big.Int, // η1 = k_i
	signatureFactorMask *big.Int, // η2 = c_i
	paillierR *big.Int, // Paillier randomness r
	params *PublicParameters,
	random io.Reader,
) (*EcdsaSignatureFactorRangeProof, error) {
	alpha, err := rand.Int(random, params.QCube())
	if err != nil {
		return nil, fmt.Errorf("could not construct ZKP2i [%v]", err)
	}
	beta, err := paillier.GetRandomNumberInMultiplicativeGroup(params.N, random)
	if err != nil {
		return nil, fmt.Errorf("could not construct ZKP2i [%v]", err)
	}
	gamma, err := rand.Int(random, params.QCubeNTilde())
	if err != nil {
		return nil, fmt.Errorf("could not construct ZKP2i [%v]", err)
	}
	delta, err := rand.Int(random, params.QCube())
	if err != nil {
		return nil, fmt.Errorf("could not construct ZKP2i [%v]", err)
	}
	mu, err := paillier.GetRandomNumberInMultiplicativeGroup(params.N, random)
	if err != nil {
		return nil, fmt.Errorf("could not construct ZKP2i [%v]", err)
	}
	nu, err := rand.Int(random, params.QCubeNTilde())
	if err != nil {
		return nil, fmt.Errorf("could not construct ZKP2i [%v]", err)
	}
	theta, err := rand.Int(random, params.QPow8())
	if err != nil {
		return nil, fmt.Errorf("could not construct ZKP2i [%v]", err)
	}
	tau, err := rand.Int(random, params.QPow8NTilde())
	if err != nil {
		return nil, fmt.Errorf("could not construct ZKP2i [%v]", err)
	}
	rho1, err := rand.Int(random, params.QNTilde())
	if err != nil {
		return nil, fmt.Errorf("could not construct ZKP2i [%v]", err)
	}
	rho2, err := rand.Int(random, params.QPow6NTilde())
	if err != nil {
		return nil, fmt.Errorf("could not construct ZKP2i [%v]", err)
	}

	// z1 = (h1^η1) * (h2^ρ1) mod N ̃
	z1 := new(big.Int).Mod(
		new(big.Int).Mul(
			new(big.Int).Exp(params.h1, signatureFactorSecret, params.NTilde),
			new(big.Int).Exp(params.h2, rho1, params.NTilde),
		),

		params.NTilde,
	)
	// z2 = (h1^η2) * (h2^ρ2) mod N ̃
	z2 := new(big.Int).Mod(
		new(big.Int).Mul(
			new(big.Int).Exp(params.h1, signatureFactorMask, params.NTilde),
			new(big.Int).Exp(params.h2, rho2, params.NTilde),
		),

		params.NTilde,
	)
	// u1 = g^α in G
	u1 := curve.NewPoint(params.curve.ScalarBaseMult(
		new(big.Int).Mod(alpha, params.q).Bytes(),
	))
	// u2 = (Γ^α) * (β^N) mod N^2
	u2 := new(big.Int).Mod(
		new(big.Int).Mul(
			new(big.Int).Exp(params.G(), alpha, params.NSquare()),
			new(big.Int).Exp(beta, params.N, params.NSquare()),
		),
		params.NSquare(),
	)
	// u3 = (h1^α) * (h2^γ) mod N ̃
	u3 := new(big.Int).Mod(
		new(big.Int).Mul(
			new(big.Int).Exp(params.h1, alpha, params.NTilde),
			new(big.Int).Exp(params.h2, gamma, params.NTilde),
		),
		params.NTilde,
	)
	// v1 = (u^α) * Γ^(q*θ) * (μ^N) mod N^2
	v1 := new(big.Int).Mod(
		new(big.Int).Mul(
			new(big.Int).Mul(
				new(big.Int).Exp(secretKeyFactor.C, alpha, params.NSquare()),
				new(big.Int).Exp(params.G(),
					new(big.Int).Mul(params.q, theta),
					params.NSquare()),
			),
			new(big.Int).Exp(mu, params.N, params.NSquare()),
		),
		params.NSquare(),
	)
	// v2 = (h1^δ) * (h2^ν) mod N ̃
	v2 := new(big.Int).Mod(
		new(big.Int).Mul(
			new(big.Int).Exp(params.h1, delta, params.NTilde),
			new(big.Int).Exp(params.h2, nu, params.NTilde),
		),
		params.NTilde,
	)
	// v3 = (h1^θ) * (h2^τ) mod N ̃
	v3 := new(big.Int).Mod(
		new(big.Int).Mul(
			new(big.Int).Exp(params.h1, theta, params.NTilde),
			new(big.Int).Exp(params.h2, tau, params.NTilde),
		),
		params.NTilde,
	)

	// In the original paper, elliptic curve generator point is also hashed.
	// However, since g is a constant in go-ethereum, we don't include it in
	// the sum256.
	// e = hash(w, u, z1, z2, u1, u2, u3, v1, v2, v3)
	digest := sum256(
		signatureUnmask.C.Bytes(),
		secretKeyFactor.C.Bytes(),
		z1.Bytes(),
		z2.Bytes(),
		u1.Bytes(),
		u2.Bytes(),
		u3.Bytes(),
		v1.Bytes(),
		v2.Bytes(),
		v3.Bytes(),
	)
	e := new(big.Int).SetBytes(digest[:])

	// s1 = e*η1 + α
	s1 := new(big.Int).Add(
		new(big.Int).Mul(e, signatureFactorSecret),
		alpha,
	)
	// s2 = e*ρ1 + γ
	s2 := new(big.Int).Add(
		new(big.Int).Mul(e, rho1),
		gamma,
	)
	// t1 = (rc^e) * μ mod N
	t1 := new(big.Int).Mod(
		new(big.Int).Mul(
			discreteExp(paillierR, e, params.N),
			mu,
		),
		params.N,
	)
	// t2 = e*η2 + θ
	t2 := new(big.Int).Add(
		new(big.Int).Mul(e, signatureFactorMask),
		theta,
	)
	// t3 = e*ρ2 + τ
	t3 := new(big.Int).Add(
		new(big.Int).Mul(e, rho2),
		tau,
	)

	return &EcdsaSignatureFactorRangeProof{z1, z2, u1, u2, u3, v1, v2, v3, e, s1, s2, t1, t2, t3}, nil
}

// Verify checks the `EcdsaSignatureFactorRangeProof` against the provided
// signature's factor, encrypted signature unmask and encrypted secret key
// factor.
// If they match values used to generate the proof, function returns `true`.
// Otherwise, `false` is returned.
func (zkp *EcdsaSignatureFactorRangeProof) Verify(
	signatureFactorPublic *curve.Point, // r_i = g^{k_i}
	signatureUnmask *paillier.Cypher, // w_i = E(k_i * ρ + c_i * q)
	secretKeyFactor *paillier.Cypher, // u = E(ρ)
	params *PublicParameters,
) bool {
	if !zkp.allParametersInRange(params) {
		return false
	}

	u1 := zkp.evaluateVerificationU1(signatureFactorPublic, params)
	u3 := zkp.evaluateVerificationU3(params)
	v1 := zkp.evaluateVerificationV1(signatureUnmask, secretKeyFactor, params)
	v3 := zkp.evaluateVerificationV3(params)

	// In the original paper, elliptic curve generator point is also hashed.
	// However, since g is a constant in go-ethereum, we don't include it in
	// the sum256.
	// e = hash(w,u,z1,z2,u1,u2,u3,v1,v2,v3)
	digest := sum256(
		signatureUnmask.C.Bytes(),
		secretKeyFactor.C.Bytes(),
		zkp.z1.Bytes(),
		zkp.z2.Bytes(),
		u1.Bytes(),
		zkp.u2.Bytes(),
		u3.Bytes(),
		v1.Bytes(),
		zkp.v2.Bytes(),
		v3.Bytes(),
	)
	e := new(big.Int).SetBytes(digest[:])

	return zkp.u1.X.Cmp(u1.X) == 0 &&
		zkp.u1.Y.Cmp(u1.Y) == 0 &&
		zkp.u3.Cmp(u3) == 0 &&
		zkp.v1.Cmp(v1) == 0 &&
		zkp.v3.Cmp(v3) == 0 &&
		zkp.e.Cmp(e) == 0
}

// Checks whether parameters are in the expected range.
// It's a preliminary step to check if proof is not corrupted.
func (zkp *EcdsaSignatureFactorRangeProof) allParametersInRange(params *PublicParameters) bool {
	zero := big.NewInt(0)

	return isInRange(zkp.z1, zero, params.NTilde) &&
		isInRange(zkp.z2, zero, params.NTilde) &&
		isInRange(zkp.u2, zero, params.NSquare()) &&
		isInRange(zkp.u3, zero, params.NTilde) &&
		isInRange(zkp.v1, zero, params.NSquare()) &&
		isInRange(zkp.v2, zero, params.NTilde) &&
		isInRange(zkp.v3, zero, params.NTilde) &&
		isInRange(zkp.t1, zero, params.N) &&
		params.curve.IsOnCurve(zkp.u1.X, zkp.u1.Y)
}

// evaluateVerificationU1 computes u1 verification value and returns it for
// further comparison with the expected one, evaluated during the commitment
// phase.
//
// We want to verify whether u1 = (g^s1) * (r^(−e)) in G
// is equal to u1 = g^α in G
// we evaluated in the commitment phase.
//
// Since:
// s1 = e*η1 + α
// r = g^η1
//
// We can do:
// u1 = (g^s1) * (r^(−e)) = g^{e*η1+α} * (g^η1)^(-e) =
// g^{e*η1+α} * g^{-e*η1} = g^α
//
// which is exactly how u1 is evaluated during the commitment phase.
func (zkp *EcdsaSignatureFactorRangeProof) evaluateVerificationU1(
	signatureMultiplePublic *curve.Point,
	params *PublicParameters,
) *curve.Point {
	cs1x, cs1y := params.curve.ScalarBaseMult(
		new(big.Int).Mod(zkp.s1, params.q).Bytes(),
	)
	rex, rey := params.curve.ScalarMult(
		signatureMultiplePublic.X, signatureMultiplePublic.Y, zkp.e.Bytes(),
	)

	// For a Weierstrass elliptic curve form, the additive inverse of
	// (x, y) is (x, -y)
	return curve.NewPoint(params.curve.Add(
		cs1x, cs1y, rex, new(big.Int).Neg(rey),
	))
}

// evaluateVerificationU3 computes u3 verification value and returns it for
// further comparison with the expected one, evaluated during the commitment
// phase.
//
// We want to verify whether u3 = (h1^s1) * (h2^s2) * (z1^(−e)) mod N ̃
// is equal to u3 = (h1^α) * (h2^γ) mod N ̃
// we evaluated in the commitment phase.
//
// Since:
// s1 = e*η1 + α
// s2 = e*ρ1 + γ
// z1 = (h1^η1) * (h2^ρ1) mod N ̃
//
// We can do:
// u3 = (h1^s1) * (h2^s2) * (z1^(−e)) =
// h1^{e*η1+α} * h2^{e*ρ1 + γ} * [(h1^η1) * (h2^ρ1)]^(-e) =
// h1^{e*η1+α} * h2^{e*ρ1 + γ} * h1^{-e*η1} * h2^{-e*ρ1} =
// h1^α * h2^γ
//
// which is exactly how u3 is evaluated during the commitment phase.
func (zkp *EcdsaSignatureFactorRangeProof) evaluateVerificationU3(params *PublicParameters) *big.Int {
	return new(big.Int).Mod(
		new(big.Int).Mul(
			new(big.Int).Mul(
				discreteExp(params.h1, zkp.s1, params.NTilde),
				discreteExp(params.h2, zkp.s2, params.NTilde),
			),
			discreteExp(zkp.z1, new(big.Int).Neg(zkp.e), params.NTilde),
		),
		params.NTilde,
	)
}

// evaluateVerificationV1 computes v1 verification value and returns it for
// further comparison with the expected one, evaluated during the commitment
// phase.
//
// We want to verify whether v1 = (u^s1) * (Γ^(q*t2)) * (t1^N) * (w^(−e)) mod N^2
// is equal to v1 = (u^α) * (Γ^(q*θ) * (μ^N) mod N^2
// we evaluated in the commitment phase.
//
// Since:
// s1 = e*η1 + α
// t1 = (rc^e) * μ mod N
// t2 = e*η2 + θ
// w = u^η1 * Γ^(q*η2) * rc^N mod N^2
//
// We can do:
// v1 = (u^s1) * (Γ^(q*t2)) * (t1^N) * (w^(−e)) =
// u^{e*η1+α} * Γ^{q*e*η2 + q*θ} * rc^{e*N} * μ^N * u^{-e*η1} * Γ^{-e*q*η2} * rc^{-e*N} =
// u^α * Γ^(q*θ) * μ^N
//
// which is exactly how v1 is evaluated during the commitment phase.
func (zkp *EcdsaSignatureFactorRangeProof) evaluateVerificationV1(
	signatureUnmask *paillier.Cypher,
	secretKeyFactor *paillier.Cypher,
	params *PublicParameters,
) *big.Int {
	return new(big.Int).Mod(
		new(big.Int).Mul(
			new(big.Int).Mul(
				discreteExp(secretKeyFactor.C, zkp.s1, params.NSquare()),
				discreteExp(
					params.G(),
					new(big.Int).Mul(params.q, zkp.t2),
					params.NSquare(),
				),
			),
			new(big.Int).Mul(
				new(big.Int).Exp(zkp.t1, params.N, params.NSquare()),
				discreteExp(signatureUnmask.C, new(big.Int).Neg(zkp.e), params.NSquare()),
			),
		),
		params.NSquare(),
	)
}

// evaluateVerificationV3 computes v3 verification value and returns it for
// further comparison with the expected one, evaluated during the commitment
// phase.
//
// We want to verify whether v3 = (h1^t2) * (h2^t3) * (z2^(−e)) mod N ̃
// is equal to v3 = (h1^θ) * (h2^τ) mod N ̃
// we evaluated in the commitment phase.
//
// Since:
// t2 = e*η2 + θ
// t3 = e*ρ2 + τ
// z2 = (h1^η2) * (h2^ρ2) mod N ̃
//
// We can do:
// v3 = (h1^t2) * (h2^t3) * (z2^(−e)) =
// h1^{e*η2 + θ} * h2^{e*ρ2 + τ} * [(h1^η2) * (h2^ρ2)]^{-e} =
// h1^{e*η2 + θ} * h2^{e*ρ2 + τ} * h1^{-e*η2} * h2^{-e*ρ2} =
// h1^θ * h2^τ
//
// which is exactly how v3 is evaluated during the commitment phase.
func (zkp *EcdsaSignatureFactorRangeProof) evaluateVerificationV3(params *PublicParameters) *big.Int {
	return new(big.Int).Mod(
		new(big.Int).Mul(
			new(big.Int).Mul(
				discreteExp(params.h1, zkp.t2, params.NTilde),
				discreteExp(params.h2, zkp.t3, params.NTilde),
			),
			discreteExp(zkp.z2, new(big.Int).Neg(zkp.e), params.NTilde),
		),
		params.NTilde,
	)
}
