package zkp

import (
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
	g *big.Int
	r *tecdsa.CurvePoint
	w *big.Int
	u *big.Int

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
