// Package dkg conatins code that implements Distributed Key Generation protocol
// described in [GJKR 99].
//
// See http://docs.keep.network/cryptography/beacon_dkg.html#_protocol
//
//     [GJKR 99]: Gennaro R., Jarecki S., Krawczyk H., Rabin T. (1999) Secure
//         Distributed Key Generation for Discrete-Log Based Cryptosystems. In:
//         Stern J. (eds) Advances in Cryptology — EUROCRYPT ’99. EUROCRYPT 1999.
//         Lecture Notes in Computer Science, vol 1592. Springer, Berlin, Heidelberg
//         http://groups.csail.mit.edu/cis/pubs/stasio/vss.ps.gz
package dkg

import (
	"math/big"

	"github.com/golang/go/src/crypto/rand"
	"github.com/keep-network/keep-core/pkg/beacon/relay/pedersen"
)

// CalculateSharesAndCommitments starts with generating coefficients for two
// polynomials. Then it is calculating shares for all group members. Additionally,
// it calculates commitments to coefficients of first polynomial using second's
// polynomial coefficients B.
//
// See http://docs.keep.network/cryptography/beacon_dkg.html#_phase_3_polynomial_generation
func (cm *CommittingMember) CalculateSharesAndCommitments() (*MemberCommitmentsMessage, error) {
	var err error
	coefficientsSize := cm.group.dishonestThreshold + 1
	coefficientsA := make([]*big.Int, coefficientsSize)
	coefficientsB := make([]*big.Int, coefficientsSize)
	for j := 0; j < coefficientsSize; j++ {
		coefficientsA[j], err = rand.Int(rand.Reader, cm.ProtocolConfig().Q)
		if err != nil {
			return nil, err
		}
		coefficientsB[j], err = rand.Int(rand.Reader, cm.ProtocolConfig().Q)
		if err != nil {
			return nil, err
		}
	}

	cm.secretKeyShare = coefficientsA[0] // z_i = a_i0

	commitments := make([]*big.Int, coefficientsSize)
	for k := 0; k < coefficientsSize; k++ {
		// C_k = g^a_k * h^b_k mod p
		commitments[k] = pedersen.CalculateCommitment(cm.vss, coefficientsA[k], coefficientsB[k])
	}

	sharesS := make(map[*big.Int]*big.Int, cm.group.groupSize)
	sharesT := make(map[*big.Int]*big.Int, cm.group.groupSize)
	for _, id := range cm.group.MemberIDs() {
		// s_j = f_(j) mod q
		sharesS[id] = calculateShare(id, coefficientsA, cm.ProtocolConfig().Q)
		// t_j = g_(j) mod q
		sharesT[id] = calculateShare(id, coefficientsB, cm.ProtocolConfig().Q)
	}

	return &MemberCommitmentsMessage{
		senderID:    cm.ID,
		sharesS:     sharesS,
		sharesT:     sharesT,
		commitments: commitments,
	}, nil
}

// VerifySharesAndCommitments verifies member's shares and commitments received
// in messages from all group members.
// It returns accusation message with ID of members for which verification failed.
//
// See http://docs.keep.network/cryptography/beacon_dkg.html#_phase_4_share_verification
func (cm *CommittingMember) VerifySharesAndCommitments(messages []*MemberCommitmentsMessage) (*FirstAccusationsMessage, error) {
	var accusedMembersIDs []*big.Int

	// `commitmentsProduct = Π (commitments_j[k] ^ (i^k)) mod p` for k in [0..T],
	// where: j is sender's ID, i is current member ID, T is threshold.
	for _, message := range messages {
		commitmentsProduct := big.NewInt(1)
		for k, c := range message.commitments {
			commitmentsProduct = new(big.Int).Mod(
				new(big.Int).Mul(
					commitmentsProduct,
					new(big.Int).Exp(
						c,
						new(big.Int).Exp(
							cm.ID,
							big.NewInt(int64(k)),
							cm.ProtocolConfig().P,
						),
						cm.ProtocolConfig().P,
					),
				),
				cm.ProtocolConfig().P,
			)
		}

		// `expectedProduct = (g ^ s_j) * (h ^ t_j)`
		expectedProduct := pedersen.CalculateCommitment(cm.vss, message.sharesS[cm.ID], message.sharesT[cm.ID])

		if expectedProduct.Cmp(commitmentsProduct) != 0 {
			accusedMembersIDs = append(accusedMembersIDs, message.senderID)
		}
	}
	return &FirstAccusationsMessage{
		senderID:   cm.ID,
		accusedIDs: accusedMembersIDs,
	}, nil
}

// calculateShare calculates a share for given memberID.
//
// It calculates `Σ a_j * z^j mod q`for j in [0..T], where:
// - `a_j` is j coefficient
// - `z` is memberID
// - `T` is threshold
func calculateShare(memberID *big.Int, coefficients []*big.Int, mod *big.Int) *big.Int {
	result := big.NewInt(0)
	for j, a := range coefficients {
		result = new(big.Int).Mod(
			new(big.Int).Add(
				result,
				new(big.Int).Mul(
					a,
					new(big.Int).Exp(memberID, big.NewInt(int64(j)), mod),
				),
			),
			mod,
		)
	}
	return result
}
