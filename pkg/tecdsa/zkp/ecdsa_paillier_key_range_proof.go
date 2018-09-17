package zkp

import (
	"crypto/rand"
	"fmt"
	"io"
	"math/big"

	"github.com/keep-network/keep-core/pkg/tecdsa/curve"
	"github.com/keep-network/paillier"
)

// EcdsaPaillierKeyRangeProof is an implementation of Gennaro's Π_i proof for the
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
type EcdsaPaillierKeyRangeProof struct {
	z  *big.Int
	u1 *curve.Point
	u2 *big.Int
	u3 *big.Int

	e *big.Int

	s1 *big.Int
	s2 *big.Int
	s3 *big.Int
}

// CommitEcdsaPaillierKeyRange generates `EcdsaPaillierKeyRangeProof` for the
// specified DSA key shares. It's required to use the same randomness `r`
// to generate this proof as the one used for Paillier encryption of
// `secretEcdsaKeyShare` into `encryptedSecretEcdsaKeyShare`.
func CommitEcdsaPaillierKeyRange(
	secretEcdsaKeyShare *big.Int,
	publicEcdsaKeyShare *curve.Point,
	encryptedSecretEcdsaKeyShare *paillier.Cypher,
	r *big.Int,
	params *PublicParameters,
	random io.Reader,
) (*EcdsaPaillierKeyRangeProof, error) {
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

	z := new(big.Int).Mod(
		new(big.Int).Mul(
			new(big.Int).Exp(params.h1, secretEcdsaKeyShare, params.NTilde),
			new(big.Int).Exp(params.h2, rho, params.NTilde),
		),
		params.NTilde,
	)

	u1 := curve.NewPoint(params.curve.ScalarBaseMult(
		new(big.Int).Mod(alpha, params.q).Bytes(),
	))
	u2 := new(big.Int).Mod(
		new(big.Int).Mul(
			new(big.Int).Exp(params.G(), alpha, params.NSquare()),
			new(big.Int).Exp(beta, params.N, params.NSquare()),
		),
		params.NSquare(),
	)
	u3 := new(big.Int).Mod(
		new(big.Int).Mul(
			new(big.Int).Exp(params.h1, alpha, params.NTilde),
			new(big.Int).Exp(params.h2, gamma, params.NTilde),
		),
		params.NTilde,
	)

	// In the original paper, elliptic curve generator point is also hashed.
	// However, since g is a constant in go-ethereum, we don't include it in
	// the sum256.
	digest := sum256(
		publicEcdsaKeyShare.Bytes(), encryptedSecretEcdsaKeyShare.C.Bytes(),
		z.Bytes(), u1.Bytes(), u2.Bytes(), u3.Bytes(),
	)
	e := new(big.Int).SetBytes(digest[:])

	s1 := new(big.Int).Add(new(big.Int).Mul(e, secretEcdsaKeyShare), alpha)
	s2 := new(big.Int).Mod(
		new(big.Int).Mul(
			new(big.Int).Exp(r, e, params.N),
			beta,
		),
		params.N,
	)
	s3 := new(big.Int).Add(new(big.Int).Mul(e, rho), gamma)

	return &EcdsaPaillierKeyRangeProof{z, u1, u2, u3, e, s1, s2, s3}, nil
}

// Verify checks the `EcdsaPaillierKeyRangeProof` against the provided DSA secret
// key shares. If they match values used to generate the proof, function returns
// `true`. Otherwise, `false` is returned.
func (zkp *EcdsaPaillierKeyRangeProof) Verify(
	encryptedSecretEcdsaKeyShare *paillier.Cypher,
	publicEcdsaKeyShare *curve.Point,
	params *PublicParameters,
) bool {
	if !zkp.allParametersInRange(params) {
		return false
	}

	u1 := zkp.evaluateU1Verification(publicEcdsaKeyShare, params)
	u2 := zkp.evaluateU2Verification(encryptedSecretEcdsaKeyShare.C, params)
	u3 := zkp.evaluateU3Verification(params)

	digest := sum256(
		publicEcdsaKeyShare.Bytes(), encryptedSecretEcdsaKeyShare.C.Bytes(),
		zkp.z.Bytes(), u1.Bytes(), u2.Bytes(), u3.Bytes(),
	)

	e := new(big.Int).SetBytes(digest[:])

	return zkp.e.Cmp(e) == 0 &&
		zkp.u1.X.Cmp(u1.X) == 0 &&
		zkp.u1.Y.Cmp(u1.Y) == 0 &&
		zkp.u2.Cmp(u2) == 0 &&
		zkp.u3.Cmp(u3) == 0
}

// Checks whether parameters are in the expected range.
// It's a preliminary step to check if proof is not corrupted.
func (zkp *EcdsaPaillierKeyRangeProof) allParametersInRange(
	params *PublicParameters,
) bool {
	zero := big.NewInt(0)

	return isInRange(zkp.z, zero, params.NTilde) &&
		isInRange(zkp.u2, zero, params.NSquare()) &&
		isInRange(zkp.u3, zero, params.NTilde) &&
		isInRange(zkp.s2, zero, params.N)
}

// evaluateU1Verification computes u1 verification value and returns it for
// further comparison with the expected one, evaluated during the commitment
// phase.
//
// We want to verify whether u1 = g^s1 * y^-e
// is equal to u1 = g^α
// we evaluated in the commitment phase.
//
// Since:
// s1 = e * η + α
// y = g^η
//
// We can do:
// g^s1 * y^-e =
// g^{e * η + α} * (g^{-e})^η =
// g^{e * η + α} * g^{-e * η} =
// g^α
//
// which is exactly how u1 is evaluated during the commitment phase.
func (zkp *EcdsaPaillierKeyRangeProof) evaluateU1Verification(
	publicEcdsaKeyShare *curve.Point,
	params *PublicParameters,
) *curve.Point {
	gs1x, gs1y := params.curve.ScalarBaseMult(
		new(big.Int).Mod(zkp.s1, params.q).Bytes(),
	)
	yex, yey := params.curve.ScalarMult(
		publicEcdsaKeyShare.X, publicEcdsaKeyShare.Y, zkp.e.Bytes(),
	)

	// For a Weierstrass elliptic curve form, the additive inverse of
	// (x, y) is (x, -y)
	return curve.NewPoint(params.curve.Add(
		gs1x, gs1y, yex, new(big.Int).Neg(yey),
	))
}

// evaluateU2Verification computes u2 verification value and returns it for
// further comparison with the expected one, evaluated during the commitment
// phase.
//
// We want to verify whether u2 = G^s1 * (s2)^N * (w)^-e
// is equal to u2 = G^α * β^N
// we evaluated in the commitment phase.
//
// Since:
// s1 = eη + α
// s2 = r^e * β
// w = E(η) = G^η * r^N
//
// We can do:
// G^s1 * (s2)^N * (w)^-e =
// G^{eη + α} * (r^e * β)^N * (G^η * r^N)^-e =
// G^{eη + α} * G^{-eη} * β^N * r^{eN} * r^{-eN} =
// G^α * β^N
//
// which is exactly how u2 is evaluated during the commitment phase.
func (zkp *EcdsaPaillierKeyRangeProof) evaluateU2Verification(
	encryptedSecretEcdsaKeyShare *big.Int,
	params *PublicParameters,
) *big.Int {
	gs1 := new(big.Int).Exp(params.G(), zkp.s1, params.NSquare())
	s2N := new(big.Int).Exp(zkp.s2, params.N, params.NSquare())
	we := discreteExp(
		encryptedSecretEcdsaKeyShare,
		new(big.Int).Neg(zkp.e),
		params.NSquare(),
	)

	return new(big.Int).Mod(
		new(big.Int).Mul(new(big.Int).Mul(gs1, s2N), we),
		params.NSquare(),
	)
}

// evaluateU3Verification computes u3 verification value and returns it for
// further comparison with the expected one, evaluated during the commitment
// phase.
//
// We want to verify whether u3 = (h1)^{s1} * (h2)^{s3} * (z)^-e
// is equal to u3 = (h1)^alpha * (h2)^γ
// we evaluated in the commitment phase.
//
// Since:
// s1 = eη + α
// s3 = e*ρ + γ
// z = (h1)^η * (h2)^ρ
//
// We can do:
// (h1)^{s1} * (h2)^{s3} * (z)^-e =
// (h1)^{e^η + α} * h2^{e*ρ + γ} * [(h1)^η * (h2)^ρ]^-e =
// (h1)^{eη} * (h1)^{α} * (h2)^{e*ρ} * (h2)^{γ} * (h1)^{-eη} * (h2)^{-e * ρ} =
// (h1)^α * (h2)^γ
//
// which is exactly how u3 is evaluated during the commitment phase.
func (zkp *EcdsaPaillierKeyRangeProof) evaluateU3Verification(
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
