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

// Commit to c1, c2, c3
func (zkp *PI1) Commit(c1, c2, c3, eta, r *big.Int, params *ZKPPublicParameters) error {
	g := new(big.Int).Add(params.N, big.NewInt(1)) // Γ - capital gamma

	q3 := new(big.Int).Exp(params.q, big.NewInt(3), nil) // q^3
	qNTilde := new(big.Int).Mul(params.q, params.NTilde) // q * NTilde ̃
	q3NTilde := new(big.Int).Mul(q3, params.NTilde)      // q^3 * nTilde

	// α ∈ Z_q^3
	alpha, err := rand.Int(rand.Reader, q3)
	if err != nil {
		return fmt.Errorf("could not construct ZKPi [%v]", err)
	}

	// β ∈ Z_N∗
	// TODO Check if it's correct random
	beta, err := rand.Int(rand.Reader, params.N)
	if err != nil {
		return fmt.Errorf("could not construct ZKPi [%v]", err)
	}

	// ρ ∈ Z_q^N ̃
	rho, err := rand.Int(rand.Reader, qNTilde)
	if err != nil {
		return fmt.Errorf("could not construct ZKPi [%v]", err)
	}

	// γ ∈ Z_((q^3)*N ̃)
	gamma, err := rand.Int(rand.Reader, q3NTilde)
	if err != nil {
		return fmt.Errorf("could not construct ZKPi [%v]", err)
	}

	// u_1 = ((h1)^η)*((h2)^ρ) mod N ̃
	u1 := new(big.Int).Mod(
		new(big.Int).Mul(
			new(big.Int).Exp(params.h1, eta, params.NTilde),
			new(big.Int).Exp(params.h2, rho, params.NTilde),
		),
		params.NTilde)

	// u_2 = ((h1)^α)*((h2)^γ) mod N ̃
	u2 := new(big.Int).Mod(
		new(big.Int).Mul(
			new(big.Int).Exp(params.h1, alpha, params.NTilde),
			new(big.Int).Exp(params.h2, gamma, params.NTilde),
		),
		params.NTilde)

	NPow2 := new(big.Int).Exp(params.N, big.NewInt(2), nil)
	// z = ((Γ)^α)*((β)^N) mod N2
	z := new(big.Int).Mod(
		new(big.Int).Mul(
			new(big.Int).Exp(g, alpha, NPow2),
			new(big.Int).Exp(beta, params.N, NPow2),
		),
		NPow2,
	)

	// v = (c2)^α mod N^2
	v := new(big.Int).Exp(c2, alpha, NPow2)

	// e = hash(c1, c2, c3, z, u1, u2, v)
	digest := sum256(
		c1.Bytes(),
		c2.Bytes(),
		c3.Bytes(),
		z.Bytes(),
		u1.Bytes(),
		u2.Bytes(),
		v.Bytes(),
	)

	zkp.e = new(big.Int).SetBytes(digest[:])

	// s1 = eη+α
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

	// s3 =eρ+γ
	zkp.s3 = new(big.Int).Add(
		new(big.Int).Mul(zkp.e, rho),
		gamma,
	)

	return nil
}

// Verify c1, c2, c3 commitment.
func (zkp *PI1) Verify(c1, c2, c3 *big.Int, params *ZKPPublicParameters) bool {
	return false
}
