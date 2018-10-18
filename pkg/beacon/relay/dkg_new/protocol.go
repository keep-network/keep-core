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
	"fmt"
	"math/big"

	"github.com/keep-network/keep-core/pkg/beacon/relay/config"
)

// CalculateMembersSharesAndCommitments starts with generating coefficients for
// two polynomials. It then calculates shares for all group member and packs them
// in individual messages for each peer member. Additionally, it calculates
// commitments to `a` coefficients of first polynomial using second's polynomial
// `b` coefficients.
//
// See Phase 3 of the protocol specification.
func (cm *CommittingMember) CalculateMembersSharesAndCommitments() (
	[]*PeerSharesMessage,
	*MemberCommitmentsMessage,
	error,
) {
	polynomialDegree := cm.group.dishonestThreshold
	coefficientsA, err := generatePolynomial(polynomialDegree, cm.protocolConfig)
	if err != nil {
		return nil, nil, err
	}
	coefficientsB, err := generatePolynomial(polynomialDegree, cm.protocolConfig)
	if err != nil {
		return nil, nil, err
	}

	cm.secretCoefficients = coefficientsA

	// Calculate shares for other group members by evaluating polynomials defined
	// by coefficients `a_i` and `b_i`
	var sharesMessages []*PeerSharesMessage
	for _, receiverID := range cm.group.MemberIDs() {
		// s_j = f_(j) mod q
		memberShareS := evaluateMemberShare(receiverID, coefficientsA, cm.protocolConfig.Q)
		// t_j = g_(j) mod q
		memberShareT := evaluateMemberShare(receiverID, coefficientsB, cm.protocolConfig.Q)

		// Check if calculated shares for the current member. If true store them
		// without sharing in a message.
		if cm.ID.Cmp(receiverID) == 0 {
			cm.selfSecretShareS = memberShareS
			cm.selfSecretShareT = memberShareT
			continue
		}

		sharesMessages = append(sharesMessages,
			&PeerSharesMessage{
				senderID:   cm.ID,
				receiverID: receiverID,
				shareS:     memberShareS,
				shareT:     memberShareT,
			})
	}

	commitments := make([]*big.Int, len(coefficientsA))
	for k := range commitments {
		// C_k = g^a_k * h^b_k mod p
		commitments[k] = cm.vss.CalculateCommitment(
			coefficientsA[k],
			coefficientsB[k],
			cm.protocolConfig.P,
		)
	}
	commitmentsMessage := &MemberCommitmentsMessage{
		senderID:    cm.ID,
		commitments: commitments,
	}

	return sharesMessages, commitmentsMessage, nil
}

// VerifyReceivedSharesAndCommitmentsMessages verifies shares and commitments
// received in messages from peer group members.
// It returns accusation message with ID of members for which verification failed.
//
// If cannot match commitments message with shares message for given sender then
// error is returned.
//
// See Phase 4 of the protocol specification.
func (cm *CommittingMember) VerifyReceivedSharesAndCommitmentsMessages(
	sharesMessages []*PeerSharesMessage,
	commitmentsMessages []*MemberCommitmentsMessage,
) (*SecretSharesAccusationsMessage, error) {
	var accusedMembersIDs []*big.Int

	// `commitmentsProduct = Π (commitments_j[k] ^ (i^k)) mod p` for k in [0..T],
	// where: j is sender's ID, i is current member ID, T is threshold.
	for _, commitmentsMessage := range commitmentsMessages {
		commitmentsProduct := commitmentsMessage.commitments[0]
		for k, c := range commitmentsMessage.commitments[1:] {
			commitmentsProduct = new(big.Int).Mod(
				new(big.Int).Mul(
					commitmentsProduct,
					new(big.Int).Exp(
						c,
						new(big.Int).Exp(
							cm.ID,
							big.NewInt(int64(k)),
							cm.protocolConfig.P,
						),
						cm.protocolConfig.P,
					),
				),
				cm.protocolConfig.P,
			)
		}
		// Find share message sent by the same member who sent commitment message
		sharesMessageFound := false
		for _, sharesMessage := range sharesMessages {
			if sharesMessage.senderID.Cmp(commitmentsMessage.senderID) == 0 {
				sharesMessageFound = true
				// `expectedProduct = (g ^ s_ji) * (h ^ t_ji)`
				// where: j is sender's ID, i is current member ID.
				expectedProduct := cm.vss.CalculateCommitment(
					sharesMessage.shareS,
					sharesMessage.shareT,
					cm.protocolConfig.P,
				)

				if expectedProduct.Cmp(commitmentsProduct) != 0 {
					accusedMembersIDs = append(accusedMembersIDs,
						commitmentsMessage.senderID)
					break
				}
				cm.receivedSharesS[commitmentsMessage.senderID] = sharesMessage.shareS
				cm.receivedSharesT[commitmentsMessage.senderID] = sharesMessage.shareT
			}
		}
		if !sharesMessageFound {
			return nil, fmt.Errorf("cannot find shares message from member %s",
				commitmentsMessage.senderID,
			)
		}
	}

	return &SecretSharesAccusationsMessage{
		senderID:   cm.ID,
		accusedIDs: accusedMembersIDs,
	}, nil
}

// evaluateMemberShare calculates a share for given memberID.
//
// It calculates `Σ a_j * z^j mod q`for j in [0..T], where:
// - `a_j` is j coefficient
// - `z` is memberID
// - `T` is threshold
func evaluateMemberShare(memberID *big.Int, coefficients []*big.Int, mod *big.Int) *big.Int {
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

// generatePolynomial generates a random polynomial over Z_q of a given degree.
// This function will generate a slice of `degree + 1` coefficients. Each value
// will be a random `big.Int` in range (0, q).
func generatePolynomial(degree int, protocolConfig *config.DKG) ([]*big.Int, error) {
	coefficients := make([]*big.Int, degree+1)
	var err error
	for i := range coefficients {
		coefficients[i], err = protocolConfig.RandomQ()
		if err != nil {
			return nil, err
		}
	}
	return coefficients, nil
}
