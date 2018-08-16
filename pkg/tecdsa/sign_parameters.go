// Package tecdsa contains the code that implements Threshold ECDSA signatures.
// The approach is based on [GGN 16].
//
//     [GGN 16]: Gennaro R., Goldfeder S., Narayanan A. (2016) Threshold-Optimal
//          DSA/ECDSA Signatures and an Application to Bitcoin Wallet Security.
//          In: Manulis M., Sadeghi AR., Schneider S. (eds) Applied Cryptography
//          and Network Security. ACNS 2016. Lecture Notes in Computer Science,
//          vol 9696. Springer, Cham
package tecdsa

import (
	"crypto/elliptic"
	"math/big"
)

// PublicParameters for T-ECDSA key generation and signing protocol.
// Defines how many Signers are in the group, what is the group signing
// threshold, which curve is used and what's the bit length of Paillier key.
type PublicParameters struct {

	// GroupSize defines how many signers are in the group.
	GroupSize int

	// Threshold defines a group signing threshold.
	//
	// If we consider an honest-but-curious adversary, i.e. an adversary that
	// learns all the secret data of compromised server but does not change
	// their code, then [GGN 16] protocol produces signature with `n = t + 1`
	// players in the network (since all players will behave honestly, even the
	// corrupted ones).
	// But in the presence of a malicious adversary, who can force corrupted
	// players to shut down or send incorrect messages, one needs at least
	// `n = 2t + 1` players in total to guarantee robustness, i.e. the ability
	// to generate signatures even in the presence of malicious faults.
	//
	// Threshold is just for signing. If anything goes wrong during key
	// generation, e.g. one of ZKPs fails or any commitment opens incorrectly,
	// key generation protocol terminates without an output.
	Threshold int

	// Curve defines the Elliptic Curve that is used for key generation and
	// signing protocols.
	Curve elliptic.Curve

	// PaillierKeyBitLength is the length of Paillier public key.
	//
	// In order for the [GGN 16] protocol to be correct, all the homomorphic
	// operations over the ciphertexts (which are modulo `N`) must not conflict
	// with the operations modulo `q` of the DSA algorithms. Because of that,
	// [GGN 16] requires that `N > q^8`, where `N` is a paillier modulus from
	// a Paillier public key and `q` is the elliptic curve cardinality.
	//
	// For instance, secp256k1 cardinality `q` is a 256 bit number, so we must
	// have at least 2048 bit PaillierKeyBitLength.
	PaillierKeyBitLength int
}

func (pp *PublicParameters) curveCardinality() *big.Int {
	return pp.Curve.Params().N
}

// BTC and ETH require that the S value inside ECDSA signatures is at most
// the curve order divided by 2 (essentially restricting this value to its
// lower half range). `halfCurveCardinality` helps to test if S is at most
// the curve order divided by 2.
func (pp *PublicParameters) halfCurveCardinality() *big.Int {
	return new(big.Int).Rsh(pp.curveCardinality(), 1)
}
