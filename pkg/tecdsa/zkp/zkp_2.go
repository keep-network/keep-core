package zkp

import (
	"crypto/rand"
	"fmt"
	"io"
	"math/big"

	"github.com/keep-network/keep-core/pkg/tecdsa"
)

// TODO Find better name for this ZKP

// PI2 is an implementation of Gennaro's PI_2,i proof for the
// Paillier encryption scheme, as described in [GGN16], section 4.4.
//
// The proof states that:
// η1 (eta1) ∈ [−q3, q3], η1 (eta1) ∈ [−q8, q8] such that:
//   g^η1 = r
//   D(w) = η1*D(u) + q*η2
//
// This struct contains values computed by the prover.
//
//
//     [GGN 16]: Gennaro R., Goldfeder S., Narayanan A. (2016) Threshold-Optimal
//          DSA/ECDSA Signatures and an Application to Bitcoin Wallet Security.
//          In: Manulis M., Sadeghi AR., Schneider S. (eds) Applied Cryptography
//          and Network Security. ACNS 2016. Lecture Notes in Computer Science,
//          vol 9696. Springer, Cham
type PI2 struct {
	z1 *big.Int
	z2 *big.Int

	u1 *tecdsa.CurvePoint
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

// CommitZkpPi2 to the Proof PI_1,i
//
// Because of the complexity of the proof, we use the same naming for values
// as in the paper in most cases. We do an exception for function parameters:
// - `η` in the paper represents DSA secret key share,
// - `c1` in the paper represents ...... (`c2 = η ×E E(xi) = E(η*x)`), TODO Check name for this one
// - `c2` in the paper represents encrypted secret message share,
// - `c3` in the paper represents encrypted DSA secret key share (`c3 = E(η)`),
//
// We assume the Prover knows the value r ∈ Z_N∗ used to encrypt η (eta)
// such that c3 = (Γ^η)*(r^N) mod N2.
//
// First the prover chooses uniformly at random four values:
// * α(alpha) ∈ Z_q^3,
// * β(beta) ∈ Z_N∗,
// * γ(gamma) ∈ Z_((q^3)*N ̃),
// * δ(delta) ∈ Z_q^3,
// * μ(mu) ∈ Z_N∗,
// * ν(nu) ∈ Z_((q^3)*N ̃),
// * θ(theta) ∈ Z_q^8,
// * τ(tau) ∈ Z_((q^8)*N ̃),
// * ρ1(rho1) ∈ Z_q^N ̃,
// * ρ2(rho2) ∈ Z_((q^6)*N ̃).
//
// Then the prover computes u1, u2, z, v, e, s1, s2,s3. This values will be sent
// by the prover to the verifier.
func CommitZkpPi2(r,
	g *tecdsa.CurvePoint,
	w,
	u,
	eta1,
	eta2,
	rc *big.Int,
	params *PublicParameters,
	random io.Reader,
) (*PI2, error) {
	// α(alpha) ∈ Z_q^3
	alpha, err := rand.Int(random, params.QCube())
	if err != nil {
		return nil, fmt.Errorf("could not construct ZKP2i [%v]", err)
	}
	// β(beta) ∈ Z_N∗
	beta, err := randomFromMultiplicativeGroup(random, params.N)
	if err != nil {
		return nil, fmt.Errorf("could not construct ZKP2i [%v]", err)
	}
	// γ(gamma) ∈ Z_((q^3)*N ̃)
	gamma, err := rand.Int(random, params.QCubeNTilde())
	if err != nil {
		return nil, fmt.Errorf("could not construct ZKP2i [%v]", err)
	}
	// δ(delta) ∈ Z_q^3
	delta, err := rand.Int(random, params.QCube())
	if err != nil {
		return nil, fmt.Errorf("could not construct ZKP2i [%v]", err)
	}
	// μ(mu) ∈ Z_N∗
	mu, err := randomFromMultiplicativeGroup(random, params.N)
	if err != nil {
		return nil, fmt.Errorf("could not construct ZKP2i [%v]", err)
	}
	// ν(nu) ∈ Z_((q^3)*N ̃)
	nu, err := rand.Int(random, params.QCubeNTilde())
	if err != nil {
		return nil, fmt.Errorf("could not construct ZKP2i [%v]", err)
	}
	// θ(theta) ∈ Z_q^8
	theta, err := rand.Int(random, params.QEight())
	if err != nil {
		return nil, fmt.Errorf("could not construct ZKP2i [%v]", err)
	}
	// τ(tau) ∈ Z_((q^8)*N ̃)
	tau, err := rand.Int(random, params.QEightNTilde())
	if err != nil {
		return nil, fmt.Errorf("could not construct ZKP2i [%v]", err)
	}
	// ρ1(rho1) ∈ Z_q^N ̃
	rho1, err := rand.Int(random, params.QNTilde())
	if err != nil {
		return nil, fmt.Errorf("could not construct ZKP2i [%v]", err)
	}
	// ρ2(rho2) ∈ Z_((q^6)*N ̃)
	rho2, err := rand.Int(random, params.QSixNTilde())
	if err != nil {
		return nil, fmt.Errorf("could not construct ZKP2i [%v]", err)
	}

	// z1 = (h1^η1) * (h2^ρ1) mod N ̃
	z1 := new(big.Int).Mod(
		new(big.Int).Mul(
			new(big.Int).Exp(params.h1, eta1, params.NTilde),
			new(big.Int).Exp(params.h2, rho1, params.NTilde),
		),

		params.NTilde,
	)
	// z2 = (h1^η2) * (h2^ρ2) mod N ̃
	z2 := new(big.Int).Mod(
		new(big.Int).Mul(
			new(big.Int).Exp(params.h1, eta2, params.NTilde),
			new(big.Int).Exp(params.h2, rho2, params.NTilde),
		),

		params.NTilde,
	)
	// u1 = g^α in G
	u1 := tecdsa.NewCurvePoint(params.curve.ScalarBaseMult(
		new(big.Int).Mod(alpha, params.q).Bytes(),
	))
	// u2 = (Γ^α) * (β^N) mod N2
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
	// v1 = (u^α) * (Γ^(q*θ) * (μ^N) mod N2
	v1 := new(big.Int).Mod(
		new(big.Int).Mul(
			new(big.Int).Mul(
				new(big.Int).Exp(u, alpha, params.NSquare()),
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

	// e = hash(g, w, u, z1, z2, u1, u2, u3, v1, v2, v3)
	digest := sum256(
		g.Bytes(),
		w.Bytes(),
		u.Bytes(),
		z1.Bytes(),
		z2.Bytes(),
		u1.Bytes(),
		u2.Bytes(),
		u3.Bytes(),
		v1.Bytes(),
		v2.Bytes(),
		v3.Bytes(),
	)
	e := new(big.Int).Abs(new(big.Int).SetBytes(digest[:]))

	// s1 = e*η1 + α
	s1 := new(big.Int).Add(
		new(big.Int).Mul(e, eta1),
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
			new(big.Int).Exp(rc, e, params.N),
			mu,
		),
		params.N,
	)
	// t2 = e*η2 + θ
	t2 := new(big.Int).Add(
		new(big.Int).Mul(e, eta2),
		theta,
	)
	// t3 = e*ρ2 + τ
	t3 := new(big.Int).Add(
		new(big.Int).Mul(e, rho2),
		tau,
	)

	return &PI2{z1, z2, u1, u2, u3, v1, v2, v3, e, s1, s2, t1, t2, t3}, nil
}

// Verify checks the `PI2` against the provided secret message and secret key
// shares.
// If they match values used to generate the proof, function returns `true`.
// Otherwise, `false` is returned.
func (zkp *PI2) Verify(
	g *tecdsa.CurvePoint,
	r *tecdsa.CurvePoint,
	w *big.Int,
	u *big.Int,
	params *PublicParameters,
) bool {
	if !zkp.allParametersInRange(params) {
		return false
	}

	u1 := zkp.evaluateVerificationU1(r, params)
	u3 := zkp.evaluateVerificationU3(params)
	v1 := zkp.evaluateVerificationV1(w, u, params)
	v3 := zkp.evaluateVerificationV3(params)

	// e = hash(g,w,u,z1,z2,u1,u2,u3,v1,v2,v3)
	digest := sum256(
		g.Bytes(),
		w.Bytes(),
		u.Bytes(),
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
func (zkp *PI2) allParametersInRange(params *PublicParameters) bool {
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

// evaluateVerificationU1 computes z verification value and returns it for
// further comparison with the expected one, evaluated during the commitment
// phase.
//
// u1 = (c)s1 (r)−e in G
func (zkp *PI2) evaluateVerificationU1(r *tecdsa.CurvePoint, params *PublicParameters) *tecdsa.CurvePoint {
	cs1x, cs1y := params.curve.ScalarBaseMult(
		new(big.Int).Mod(zkp.s1, params.q).Bytes(),
	)
	rex, rey := params.curve.ScalarMult(
		r.X, r.Y, zkp.e.Bytes(),
	)

	// For a Weierstrass elliptic curve form, the additive inverse of
	// (x, y) is (x, -y)
	return tecdsa.NewCurvePoint(params.curve.Add(
		cs1x, cs1y, rex, new(big.Int).Neg(rey),
	))
}

// evaluateVerificationU3 computes u3 verification value and returns it for
// further comparison with the expected one, evaluated during the commitment
// phase.
//
// u2 = (h1^s1( * (h2^s2) (z1^(−e)) mod N ̃
func (zkp *PI2) evaluateVerificationU3(params *PublicParameters) *big.Int {
	return new(big.Int).Mod(
		new(big.Int).Mul(
			new(big.Int).Mul(
				new(big.Int).Exp(params.h1, zkp.s1, params.NTilde),
				new(big.Int).Exp(params.h2, zkp.s2, params.NTilde),
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
// v1 = (u^s1) * (Γ^(q*t2)) * (t1^N) * (w^(−e)) mod N2
func (zkp *PI2) evaluateVerificationV1(w, u *big.Int, params *PublicParameters) *big.Int {
	return new(big.Int).Mod(
		new(big.Int).Mul(
			new(big.Int).Mul(
				new(big.Int).Exp(u, zkp.s1, params.NSquare()),
				new(big.Int).Exp(
					params.G(),
					new(big.Int).Mul(params.q, zkp.t2), params.NSquare()),
			),
			new(big.Int).Mul(
				new(big.Int).Exp(zkp.t1, params.N, params.NSquare()),
				discreteExp(w, new(big.Int).Neg(zkp.e), params.NSquare()),
			),
		),
		params.NSquare(),
	)
}

// evaluateVerificationV3 computes v3 verification value and returns it for
// further comparison with the expected one, evaluated during the commitment
// phase.
//
// v3 = (h1^t2) * (h2^t3) * (z2^(−e)) mod N ̃
func (zkp *PI2) evaluateVerificationV3(params *PublicParameters) *big.Int {
	return new(big.Int).Mod(
		new(big.Int).Mul(
			new(big.Int).Mul(
				new(big.Int).Exp(params.h1, zkp.t2, params.NTilde),
				new(big.Int).Exp(params.h2, zkp.t3, params.NTilde),
			),
			discreteExp(zkp.z2, new(big.Int).Neg(zkp.e), params.NTilde),
		),
		params.NTilde,
	)
}
