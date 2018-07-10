package zkp

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// PI1 is an implementation of Gennaro's ZKP PI_1,i proof.
// This struct contains values computed by the prover.
type PI1 struct {
	z *big.Int
	v *big.Int

	u1 *big.Int
	u2 *big.Int

	e *big.Int

	s1 *big.Int
	s2 *big.Int
	s3 *big.Int
}

// Commit to the Proof PI_1,i, which states that there exists
// η (eta) ∈ [−q3, q3] such that:
//   D(c1) = η*D(c2)
//   D(c3) = η
//
// We assume the Prover knows the value r ∈ Z_N∗ used to encrypt η (eta)
// such that c3 = (Γ^η)*(r^N) mod N2.
//
// First the prover chooses uniformly at random four values:
// * α(alpha) ∈ Z_q^3,
// * β(beta) ∈ Z_N∗,
// * ρ(rho) ∈ Z_q^N ̃,
// * γ(gamma) ∈ Z_((q^3)*N ̃).
//
// Then the prover computes u1, u2, z, v, e, s1, s2,s3. This values will be sent
// by the prover to the verifier.
func (zkp *PI1) Commit(eta, r, c1, c2, c3 *big.Int, params *ZKPPublicParameters) error {
	q3 := new(big.Int).Exp(params.q, big.NewInt(3), nil) // q^3

	alpha, err := rand.Int(rand.Reader, q3)
	if err != nil {
		return fmt.Errorf("could not construct ZKP1i [%v]", err)
	}

	beta, err := randomFromMultiplicativeGroup(rand.Reader, params.N)
	if err != nil {
		return fmt.Errorf("could not construct ZKP1i [%v]", err)
	}

	rho, err := rand.Int(rand.Reader, new(big.Int).Mul(params.q, params.NTilde))
	if err != nil {
		return fmt.Errorf("could not construct ZKP1i [%v]", err)
	}

	gamma, err := rand.Int(rand.Reader, new(big.Int).Mul(q3, params.NTilde))
	if err != nil {
		return fmt.Errorf("could not construct ZKP1i [%v]", err)
	}

	// u_1 = ((h1)^η)*((h2)^ρ) mod N ̃
	zkp.u1 = new(big.Int).Mod(
		new(big.Int).Mul(
			new(big.Int).Exp(params.h1, eta, params.NTilde),
			new(big.Int).Exp(params.h2, rho, params.NTilde),
		),
		params.NTilde,
	)

	// u_2 = ((h1)^α)*((h2)^γ) mod N ̃
	zkp.u2 = new(big.Int).Mod(
		new(big.Int).Mul(
			new(big.Int).Exp(params.h1, alpha, params.NTilde),
			new(big.Int).Exp(params.h2, gamma, params.NTilde),
		),
		params.NTilde,
	)

	// z = ((Γ)^α)*((β)^N) mod N^2
	zkp.z = new(big.Int).Mod(
		new(big.Int).Mul(
			new(big.Int).Exp(params.G, alpha, params.N2),
			new(big.Int).Exp(beta, params.N, params.N2),
		),
		params.N2,
	)

	// v = (c2)^α mod N^2
	zkp.v = new(big.Int).Exp(c2, alpha, params.N2)

	// e = hash(c1, c2, c3, z, u1, u2, v)
	digest := sum256(
		c1.Bytes(),
		c2.Bytes(),
		c3.Bytes(),
		zkp.z.Bytes(),
		zkp.u1.Bytes(),
		zkp.u2.Bytes(),
		zkp.v.Bytes(),
	)
	zkp.e = new(big.Int).Abs(new(big.Int).SetBytes(digest[:]))

	// s1 = e*η+α
	zkp.s1 = new(big.Int).Add(
		new(big.Int).Mul(zkp.e, eta),
		alpha,
	)

	// s2 = (r^e)*β mod N
	zkp.s2 = new(big.Int).Mod(
		new(big.Int).Mul(
			new(big.Int).Exp(r, zkp.e, params.N),
			beta,
		),
		params.N,
	)

	// s3 = e*ρ+γ
	zkp.s3 = new(big.Int).Add(
		new(big.Int).Mul(zkp.e, rho),
		gamma,
	)

	return nil
}

// Verify c1, c2, c3 commitment.
func (zkp *PI1) Verify(c1, c2, c3 *big.Int, params *ZKPPublicParameters) bool {
	// Check that all the values are in the correct ranges
	if zkp.u1.CmpAbs(params.NTilde) >= 0 ||
		zkp.u2.CmpAbs(params.NTilde) >= 0 ||
		zkp.z.CmpAbs(params.N2) >= 0 ||
		zkp.v.CmpAbs(params.N2) >= 0 ||
		zkp.s2.CmpAbs(params.N) >= 0 {
		return false
	}

	// z = (Γ^s1)*(s2^N)*(c3^(−e)) mod N^2
	z := new(big.Int).Mod(
		new(big.Int).Mul(
			new(big.Int).Mul(
				new(big.Int).Exp(params.G, zkp.s1, params.N2),
				new(big.Int).Exp(zkp.s2, params.N, params.N2),
			),
			// Exp function doesn't support negative exponential
			// for y < 0: x**y mod m == (x**(-1))**|y| mod m
			new(big.Int).Exp(
				new(big.Int).ModInverse(c3, params.N2),
				new(big.Int).Abs(zkp.e),
				params.N2,
			),
		),
		params.N2,
	)

	// v = (c2^s1)*(c1^(−e)) mod N^2
	v := new(big.Int).Mod(
		new(big.Int).Mul(
			new(big.Int).Exp(c2, zkp.s1, params.N2),
			// Exp function doesn't support negative exponential
			// for y < 0: x**y mod m == (x**(-1))**|y| mod m
			new(big.Int).Exp(
				new(big.Int).ModInverse(c1, params.N2),
				new(big.Int).Abs(zkp.e),
				params.N2,
			),
		),
		params.N2,
	)

	// u1 was calculated by the prover
	u1 := zkp.u1

	// u2 =(h1^s1)*(h2^s3)*(u1^(−e)) mod N ̃
	u2 := new(big.Int).Mod(
		new(big.Int).Mul(
			new(big.Int).Mul(
				new(big.Int).Exp(params.h1, zkp.s1, params.NTilde),
				new(big.Int).Exp(params.h2, zkp.s3, params.NTilde),
			),
			// Exp function doesn't support negative exponential
			// for y < 0: x**y mod m == (x**(-1))**|y| mod m
			new(big.Int).Exp(
				new(big.Int).ModInverse(u1, params.NTilde),
				new(big.Int).Abs(zkp.e),
				params.NTilde,
			),
		),
		params.NTilde,
	)

	// e = hash(c1,c2,c3,z,u1,u2,v)
	digest := sum256(
		c1.Bytes(),
		c2.Bytes(),
		c3.Bytes(),
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
