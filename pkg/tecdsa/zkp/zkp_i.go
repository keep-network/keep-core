package zkp

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
)

// PIi is an implementation of Gennaro's ZKP PIi proof.
type PIi struct {
	z        *big.Int
	u1x, u1y *big.Int
	u2       *big.Int
	u3       *big.Int

	e *big.Int

	s1 *big.Int
	s2 *big.Int
	s3 *big.Int
}

// Commit to y and w
func (zkp *PIi) Commit(y, w, eta, r *big.Int, params *PublicParameters) error {
	g := new(big.Int).Add(params.N, big.NewInt(1))

	q3 := new(big.Int).Exp(params.q, big.NewInt(3), nil) // q^3
	qNTilde := new(big.Int).Mul(params.q, params.NTilde) // q * NTilde
	q3NTilde := new(big.Int).Mul(q3, params.NTilde)      // q^3 * nTilde

	alpha, err := rand.Int(rand.Reader, q3)
	if err != nil {
		return fmt.Errorf("could not construct ZKPi [%v]", err)
	}

	beta, err := rand.Int(rand.Reader, params.N)
	if err != nil {
		return fmt.Errorf("could not construct ZKPi [%v]", err)
	}

	rho, err := rand.Int(rand.Reader, qNTilde)
	if err != nil {
		return fmt.Errorf("could not construct ZKPi [%v]", err)
	}

	gamma, err := rand.Int(rand.Reader, q3NTilde)
	if err != nil {
		return fmt.Errorf("could not construct ZKPi [%v]", err)
	}

	z := new(big.Int).Mod(
		new(big.Int).Mul(
			new(big.Int).Exp(params.h1, eta, params.NTilde),
			new(big.Int).Exp(params.h2, rho, params.NTilde),
		),
		params.NTilde,
	)

	u1x, u1y := params.curve.ScalarBaseMult(alpha.Bytes())

	NPow2 := new(big.Int).Exp(params.N, big.NewInt(2), nil)
	u2 := new(big.Int).Mod(
		new(big.Int).Mul(
			new(big.Int).Exp(g, alpha, NPow2),
			new(big.Int).Exp(beta, params.N, NPow2),
		),
		NPow2,
	)

	u3 := new(big.Int).Mod(
		new(big.Int).Mul(
			new(big.Int).Exp(params.h1, alpha, params.NTilde),
			new(big.Int).Exp(params.h2, gamma, params.NTilde),
		),
		params.NTilde,
	)

	digest := sum256(
		g.Bytes(),
		y.Bytes(),
		w.Bytes(),
		z.Bytes(),
		u1x.Bytes(),
		u1y.Bytes(),
		u2.Bytes(),
		u3.Bytes(),
	)

	zkp.e = new(big.Int).SetBytes(digest[:])

	zkp.s1 = new(big.Int).Add(new(big.Int).Mul(zkp.e, eta), alpha)
	zkp.s2 = new(big.Int).Mod(
		new(big.Int).Mul(
			new(big.Int).Exp(r, zkp.e, params.N),
			beta,
		),
		params.N,
	)
	zkp.s3 = new(big.Int).Add(new(big.Int).Mul(zkp.e, rho), gamma)

	return nil
}

// Verify y and w commitment.
func (zkp *PIi) Verify(y, w *big.Int, params *PublicParameters) bool {
	u1x, u1y := zkp.u1Verification(y, params)
	u2 := zkp.u2Verification(w, params)
	u3 := zkp.u3Verification(params)

	// TODO: add hash verification
	return zkp.u1x == u1x &&
		zkp.u1y == u1y &&
		zkp.u2 == u2 &&
		zkp.u3 == u3
}

func (zkp *PIi) u1Verification(y *big.Int, params *PublicParameters) (*big.Int, *big.Int) {
	gs1x, gs1y := params.curve.ScalarBaseMult(zkp.s1.Bytes())
	ye := new(big.Int).Exp(y, new(big.Int).Neg(zkp.e), nil)
	return params.curve.ScalarMult(gs1x, gs1y, ye.Bytes())
}

func (zkp *PIi) u2Verification(w *big.Int, params *PublicParameters) *big.Int {
	g := new(big.Int).Add(params.N, big.NewInt(1))
	nSquare := new(big.Int).Exp(params.N, big.NewInt(2), nil)

	gs1 := new(big.Int).Exp(g, zkp.s1, nSquare)
	s2N := new(big.Int).Exp(zkp.s2, params.N, nSquare)
	we := discreteExp(w, new(big.Int).Neg(zkp.e), nSquare)

	return new(big.Int).Mod(
		new(big.Int).Mul(new(big.Int).Mul(gs1, s2N), we),
		nSquare,
	)
}

func (zkp *PIi) u3Verification(params *PublicParameters) *big.Int {
	h1s1 := discreteExp(params.h1, zkp.s1, params.NTilde)
	h2s3 := discreteExp(params.h2, zkp.s3, params.NTilde)
	ze := discreteExp(zkp.z, new(big.Int).Neg(zkp.e), params.NTilde)

	return new(big.Int).Mod(
		new(big.Int).Mul(new(big.Int).Mul(h1s1, h2s3), ze),
		params.NTilde,
	)
}

func sum256(data ...[]byte) [sha256.Size]byte {
	accumulator := make([]byte, 0)
	for _, d := range data {
		accumulator = append(accumulator, d...)
	}
	return sha256.Sum256(accumulator)
}

func discreteExp(a, b, c *big.Int) *big.Int { //TODO: use it in every single case?
	if b.Cmp(big.NewInt(0)) == -1 { // b < 0 ?
		ret := new(big.Int).Exp(a, new(big.Int).Neg(b), c)
		return new(big.Int).ModInverse(ret, c)
	}
	return new(big.Int).Exp(a, b, c)
}
