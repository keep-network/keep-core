package zkp

import (
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
)

// PI1 is an implementation of Gennaro's PI_1,i proof for the
// Paillier encryption scheme, as described in [GGN16], section 4.4.
//
// The proof states that:
// η (eta) ∈ [−q3, q3] such that:
//   D(c1) = η*D(c2)
//   D(c3) = η
//
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

// CommitZkpPi1 to the Proof PI_1,i
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
// * ρ(rho) ∈ Z_q^N ̃,
// * γ(gamma) ∈ Z_((q^3)*N ̃).
//
// Then the prover computes u1, u2, z, v, e, s1, s2,s3. This values will be sent
// by the prover to the verifier.
func CommitZkpPi1(secretKeyShare,
	c1,
	encryptedMessageShare,
	encryptedSecretKeyShare,
	r *big.Int,
	params *PublicParameters,
	random io.Reader,
) (*PI1, error) {
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
			new(big.Int).Exp(params.h1, secretKeyShare, params.NTilde),
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
	v := new(big.Int).Exp(encryptedMessageShare, alpha, params.NSquare())

	// e = hash(c1, c2, c3, z, u1, u2, v)
	digest := sum256(
		c1.Bytes(),
		encryptedMessageShare.Bytes(),
		encryptedSecretKeyShare.Bytes(),
		z.Bytes(),
		u1.Bytes(),
		u2.Bytes(),
		v.Bytes(),
	)
	e := new(big.Int).Abs(new(big.Int).SetBytes(digest[:]))

	// s1 = e*η+α
	s1 := new(big.Int).Add(
		new(big.Int).Mul(e, secretKeyShare),
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

	return &PI1{z, v, u1, u2, e, s1, s2, s3}, nil
}

// Verify c1, c2, c3 commitment.
func (zkp *PI1) Verify(c1,
	encryptedMessageShare,
	encryptedSecretKeyShare *big.Int,
	params *PublicParameters,
) bool {
	if !zkp.allParametersInRange(params) {
		return false
	}

	z := evaluateVerificationZ(encryptedSecretKeyShare, zkp.s1, zkp.s2, zkp.e, params)
	v := evaluateVerificationV(c1, encryptedMessageShare, zkp.s1, zkp.e, params)
	u1 := zkp.u1 // u1 was calculated by the prover
	u2 := evaluateVerificationU2(zkp.u1, zkp.s1, zkp.s3, zkp.e, params)

	// e = hash(c1,c2,c3,z,u1,u2,v)
	digest := sum256(
		c1.Bytes(),
		encryptedMessageShare.Bytes(),
		encryptedSecretKeyShare.Bytes(),
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
func (zkp *PI1) allParametersInRange(params *PublicParameters) bool {
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
func evaluateVerificationZ(encryptedSecretKeyShare, s1, s2, e *big.Int, params *PublicParameters) *big.Int {
	return new(big.Int).Mod(
		new(big.Int).Mul(
			new(big.Int).Mul(
				new(big.Int).Exp(params.G(), s1, params.NSquare()),
				new(big.Int).Exp(s2, params.N, params.NSquare()),
			),
			discreteExp(encryptedSecretKeyShare, new(big.Int).Neg(e), params.NSquare()),
		),
		params.NSquare(),
	)
}

// evaluateVerificationV computes v verification value and returns it for
// further comparison with the expected one, evaluated during the commitment
// phase.
//
// v = (c2^s1)*(c1^(−e)) mod N^2
func evaluateVerificationV(c1, encryptedMessageShare, s1, e *big.Int, params *PublicParameters) *big.Int {
	return new(big.Int).Mod(
		new(big.Int).Mul(
			new(big.Int).Exp(encryptedMessageShare, s1, params.NSquare()),
			discreteExp(c1, new(big.Int).Neg(e), params.NSquare()),
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
