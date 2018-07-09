package zkp

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"math/big"

	"github.com/keep-network/keep-core/pkg/tecdsa"
	"github.com/keep-network/paillier"
)

// DsaPaillierKeyRangeProof is an implementation of Gennaro's PI_i proof for the
// Paillier encryption scheme, as described in [GGN16], section 4.4.
//
// The proof is used in the DSA key initialization phase of T-ECDSA and operates
// on DSA secret and public key shares produced by each signer in the group.
//
// Because of the complexity of the proof, we use the same naming for values
// as in the paper in most cases. We do an exception for function parameters:
// - `η` in the paper represents DSA secret key share,
// - `w` in the paper represents encrypted DSA secret key share (`w = E(η)`),
// - `y` in the paper represents public DSA key share.
//
// The proof states that:
// ∃ η ∈ [-q^3, g^3] such that
// g^η = y
// D(w) = η
//
// In other words, for the Elliptic Curve of cardinality `q`, secret DSA key
// share η, public DSA key share y = g^η (`g` is the generator of the Elliptic
// Curve), Paillier encryption scheme used to encrypt DSA secret key share
// `w = E(η)`, we state that the secret DSA key share η is in range [-q^3, g^3]
// without revealing it. We also state that `w` is Paillier-encrypted secret
// DSA key share `η`.
//
//
//     [GGN 16]: Gennaro R., Goldfeder S., Narayanan A. (2016) Threshold-Optimal
//          DSA/ECDSA Signatures and an Application to Bitcoin Wallet Security.
//          In: Manulis M., Sadeghi AR., Schneider S. (eds) Applied Cryptography
//          and Network Security. ACNS 2016. Lecture Notes in Computer Science,
//          vol 9696. Springer, Cham
type DsaPaillierKeyRangeProof struct {
	z  *big.Int
	u1 *tecdsa.CurvePoint
	u2 *big.Int
	u3 *big.Int

	e *big.Int

	s1 *big.Int
	s2 *big.Int
	s3 *big.Int
}

// CommitDsaPaillierKeyRange generates `DsaPaillierKeyRangeProof` for the
// specified DSA key shares. It's required to use the same randomness `r`
// to generate this proof as the one used for Paillier encryption of
// `secretDsaKeyShare` into `encryptedSecretDsaKeyShare`.
func CommitDsaPaillierKeyRange(
	secretDsaKeyShare *big.Int,
	publicDsaKeyShare *tecdsa.CurvePoint,
	encryptedSecretDsaKeyShare *paillier.Cypher,
	r *big.Int,
	params *PublicParameters,
	random io.Reader,
) (*DsaPaillierKeyRangeProof, error) {
	// TODO: move g to PublicParameters
	g := new(big.Int).Add(params.N, big.NewInt(1))

	q3 := new(big.Int).Exp(params.q, big.NewInt(3), nil)    // q^3
	qNTilde := new(big.Int).Mul(params.q, params.NTilde)    // q * NTilde
	q3NTilde := new(big.Int).Mul(q3, params.NTilde)         // q^3 * nTilde
	NPow2 := new(big.Int).Exp(params.N, big.NewInt(2), nil) // N^2

	alpha, err := rand.Int(random, q3)
	if err != nil {
		return nil, fmt.Errorf("could not construct the proof [%v]", err)
	}

	beta, err := randomFromMultiplicativeGroup(random, params.N)
	if err != nil {
		return nil, fmt.Errorf("could not construct the proof [%v]", err)
	}

	rho, err := rand.Int(random, qNTilde)
	if err != nil {
		return nil, fmt.Errorf("could not construct the proof [%v]", err)
	}

	gamma, err := rand.Int(random, q3NTilde)
	if err != nil {
		return nil, fmt.Errorf("could not construct the proof [%v]", err)
	}

	z := new(big.Int).Mod(
		new(big.Int).Mul(
			new(big.Int).Exp(params.h1, secretDsaKeyShare, params.NTilde),
			new(big.Int).Exp(params.h2, rho, params.NTilde),
		),
		params.NTilde,
	)

	u1 := tecdsa.NewCurvePoint(params.curve.ScalarBaseMult(
		new(big.Int).Mod(alpha, params.q).Bytes(),
	))
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
		g.Bytes(), publicDsaKeyShare.X.Bytes(), publicDsaKeyShare.Y.Bytes(),
		encryptedSecretDsaKeyShare.C.Bytes(), z.Bytes(),
		u1.X.Bytes(), u1.Y.Bytes(), u2.Bytes(), u3.Bytes(),
	)
	e := new(big.Int).SetBytes(digest[:])

	s1 := new(big.Int).Add(new(big.Int).Mul(e, secretDsaKeyShare), alpha)
	s2 := new(big.Int).Mod(
		new(big.Int).Mul(
			new(big.Int).Exp(r, e, params.N),
			beta,
		),
		params.N,
	)
	s3 := new(big.Int).Add(new(big.Int).Mul(e, rho), gamma)

	return &DsaPaillierKeyRangeProof{z, u1, u2, u3, e, s1, s2, s3}, nil
}

// Verify checks the `DsaPaillierKeyRangeProof` against the provided DSA secret
// key shares. If they match values used to generate the proof, function returns
// `true`. Otherwise, `false` is returned.
func (zkp *DsaPaillierKeyRangeProof) Verify(
	encryptedSecretDsaKeyShare *paillier.Cypher,
	publicDsaKeyShare *tecdsa.CurvePoint,
	params *PublicParameters,
) bool {
	u1 := zkp.u1Verification(publicDsaKeyShare, params)
	u2 := zkp.u2Verification(encryptedSecretDsaKeyShare.C, params)
	u3 := zkp.u3Verification(params)

	g := new(big.Int).Add(params.N, big.NewInt(1))

	digest := sum256(
		g.Bytes(), publicDsaKeyShare.X.Bytes(), publicDsaKeyShare.Y.Bytes(),
		encryptedSecretDsaKeyShare.C.Bytes(), zkp.z.Bytes(),
		u1.X.Bytes(), u1.Y.Bytes(), u2.Bytes(), u3.Bytes(),
	)

	e := new(big.Int).SetBytes(digest[:])

	return zkp.e.Cmp(e) == 0 &&
		zkp.u1.X.Cmp(u1.X) == 0 &&
		zkp.u1.Y.Cmp(u1.Y) == 0 &&
		zkp.u2.Cmp(u2) == 0 &&
		zkp.u3.Cmp(u3) == 0
}

// We verify whether u1 = g^s1 * y^-e
//
// Since:
// s1 = e * eta + alpha
// y = g^eta
//
// We can do:
// g^s1 * y^-e =
// g^{e * eta + alpha} * (g^{-e})^eta =
// g^{e * eta + alpha} * g^{-e * eta} =
// g^alpha
//
// which is exactly how u1 is evaluated during the commitment phase
func (zkp *DsaPaillierKeyRangeProof) u1Verification(
	publicDsaKeyShare *tecdsa.CurvePoint,
	params *PublicParameters,
) *tecdsa.CurvePoint {
	gs1x, gs1y := params.curve.ScalarBaseMult(
		new(big.Int).Mod(zkp.s1, params.q).Bytes(),
	)
	yx, yy := params.curve.ScalarMult(
		publicDsaKeyShare.X, publicDsaKeyShare.Y, zkp.e.Bytes(),
	)

	// For a Weierstrass elliptic curve form, the additive inverse of
	// (x, y) is (x, -y)
	return tecdsa.NewCurvePoint(params.curve.Add(
		gs1x, gs1y, yx, new(big.Int).Neg(yy),
	))
}

func (zkp *DsaPaillierKeyRangeProof) u2Verification(
	encryptedSecretDsaKeyShare *big.Int,
	params *PublicParameters,
) *big.Int {
	g := new(big.Int).Add(params.N, big.NewInt(1)) // TODO: move to PublicParameters
	nSquare := new(big.Int).Exp(params.N, big.NewInt(2), nil)

	gs1 := new(big.Int).Exp(g, zkp.s1, nSquare)
	s2N := new(big.Int).Exp(zkp.s2, params.N, nSquare)
	we := discreteExp(
		encryptedSecretDsaKeyShare,
		new(big.Int).Neg(zkp.e),
		nSquare,
	)

	return new(big.Int).Mod(
		new(big.Int).Mul(new(big.Int).Mul(gs1, s2N), we),
		nSquare,
	)
}

func (zkp *DsaPaillierKeyRangeProof) u3Verification(
	params *PublicParameters,
) *big.Int {
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
