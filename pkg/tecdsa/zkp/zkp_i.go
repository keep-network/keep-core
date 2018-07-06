package zkp

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"math/big"

	"github.com/keep-network/keep-core/pkg/tecdsa"
)

// PIi is an implementation of Gennaro's ZKP PIi proof.
type PIi struct {
	z  *big.Int
	u1 *tecdsa.CurvePoint
	u2 *big.Int
	u3 *big.Int

	e *big.Int

	s1 *big.Int
	s2 *big.Int
	s3 *big.Int
}

// CommitPIi to y and w
func CommitPIi(
	w, eta, r *big.Int,
	y *tecdsa.CurvePoint,
	params *PublicParameters,
	random io.Reader,
) (*PIi, error) {
	g := new(big.Int).Add(params.N, big.NewInt(1))

	q3 := new(big.Int).Exp(params.q, big.NewInt(3), nil) // q^3
	qNTilde := new(big.Int).Mul(params.q, params.NTilde) // q * NTilde
	q3NTilde := new(big.Int).Mul(q3, params.NTilde)      // q^3 * nTilde

	alpha, err := rand.Int(random, q3)
	if err != nil {
		return nil, fmt.Errorf("could not construct ZKPi [%v]", err)
	}

	beta, err := randomFromMultiplicativeGroup(random, params.N)
	if err != nil {
		return nil, fmt.Errorf("could not construct ZKPi [%v]", err)
	}

	rho, err := rand.Int(random, qNTilde)
	if err != nil {
		return nil, fmt.Errorf("could not construct ZKPi [%v]", err)
	}

	gamma, err := rand.Int(random, q3NTilde)
	if err != nil {
		return nil, fmt.Errorf("could not construct ZKPi [%v]", err)
	}

	z := new(big.Int).Mod(
		new(big.Int).Mul(
			new(big.Int).Exp(params.h1, eta, params.NTilde),
			new(big.Int).Exp(params.h2, rho, params.NTilde),
		),
		params.NTilde,
	)

	u1 := tecdsa.NewCurvePoint(params.curve.ScalarBaseMult(
		new(big.Int).Mod(alpha, params.q).Bytes(),
	))

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
		g.Bytes(), y.X.Bytes(), y.Y.Bytes(), w.Bytes(), z.Bytes(),
		u1.X.Bytes(), u1.Y.Bytes(), u2.Bytes(), u3.Bytes(),
	)

	e := new(big.Int).SetBytes(digest[:])

	s1 := new(big.Int).Add(new(big.Int).Mul(e, eta), alpha)
	s2 := new(big.Int).Mod(
		new(big.Int).Mul(
			new(big.Int).Exp(r, e, params.N),
			beta,
		),
		params.N,
	)
	s3 := new(big.Int).Add(new(big.Int).Mul(e, rho), gamma)

	return &PIi{z, u1, u2, u3, e, s1, s2, s3}, nil
}

// Verify y and w commitment.
func (zkp *PIi) Verify(
	w *big.Int,
	y *tecdsa.CurvePoint,
	params *PublicParameters,
) bool {
	u1 := zkp.u1Verification(y, params)
	u2 := zkp.u2Verification(w, params)
	u3 := zkp.u3Verification(params)

	// TODO: add hash verification
	return zkp.u1 == u1 &&
		zkp.u2 == u2 &&
		zkp.u3 == u3
}

func (zkp *PIi) u1Verification(y *tecdsa.CurvePoint, params *PublicParameters) *tecdsa.CurvePoint {
	gs1x, gs1y := params.curve.ScalarBaseMult(zkp.s1.Bytes())
	yx, yy := params.curve.ScalarMult(y.X, y.Y, new(big.Int).Neg(zkp.e).Bytes())

	return tecdsa.NewCurvePoint(params.curve.Add(gs1x, gs1y, yx, yy))
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

func randomFromMultiplicativeGroup(random io.Reader, n *big.Int) (*big.Int, error) {
	for {
		r, err := rand.Int(random, n)

		if err != nil {
			return nil, err
		}

		if new(big.Int).GCD(nil, nil, r, n).Cmp(big.NewInt(1)) == 0 {
			return r, nil
		}
	}
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
