// Package gjkr conatins code that implements Distributed Key Generation protocol
// described in [GJKR 99].
//
// See http://docs.keep.network/cryptography/beacon_dkg.html#_protocol
//
//     [GJKR 99]: Gennaro R., Jarecki S., Krawczyk H., Rabin T. (1999) Secure
//         Distributed Key Generation for Discrete-Log Based Cryptosystems. In:
//         Stern J. (eds) Advances in Cryptology — EUROCRYPT ’99. EUROCRYPT 1999.
//         Lecture Notes in Computer Science, vol 1592. Springer, Berlin, Heidelberg
//         http://groups.csail.mit.edu/cis/pubs/stasio/vss.ps.gz
package gjkr

import (
	"fmt"
	"math"
	"math/big"
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
		// s_j = f_(j)
		memberShareS := evaluateMemberShare(receiverID, coefficientsA)
		// t_j = g_(j)
		memberShareT := evaluateMemberShare(receiverID, coefficientsB)

		// Check if calculated shares for the current member. If true store them
		// without sharing in a message.
		if cm.ID == receiverID {
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
	var accusedMembersIDs []int

	// `commitmentsProduct = Π (commitments_j[k] ^ (i^k)) mod p` for k in [0..T],
	// where: j is sender's ID, i is current member ID, T is threshold.
	for _, commitmentsMessage := range commitmentsMessages {
		commitmentsProduct := big.NewInt(1)
		for k, c := range commitmentsMessage.commitments {
			commitmentsProduct = new(big.Int).Mod(
				new(big.Int).Mul(
					commitmentsProduct,
					new(big.Int).Exp(
						c,
						pow(cm.ID, k),
						cm.protocolConfig.P,
					),
				),
				cm.protocolConfig.P,
			)
		}
		// Find share message sent by the same member who sent commitment message
		sharesMessageFound := false
		for _, sharesMessage := range sharesMessages {
			if sharesMessage.senderID == commitmentsMessage.senderID {
				sharesMessageFound = true
				// `expectedProduct = (g ^ s_ji) * (h ^ t_ji) mod p`
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
			return nil, fmt.Errorf("cannot find shares message from member %d",
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
// It calculates `s_j = Σ a_k * j^k`for k in [0..T], where:
// - `a_k` is k coefficient
// - `j` is memberID
// - `T` is threshold
//
// Note: [GJKR] fig. 2 pt. 1.a. states that calculation should be done `mod q`.
// Our tests gave unstable results if doing so. We decided not to be using modulo
// operation here.
func evaluateMemberShare(memberID int, coefficients []*big.Int) *big.Int {
	result := big.NewInt(0)
	for k, a := range coefficients {
		result = new(big.Int).Add(
			result,
			new(big.Int).Mul(
				a,
				pow(memberID, k),
			),
		)
	}
	return result
}

// generatePolynomial generates a random polynomial over Z_q of a given degree.
// This function will generate a slice of `degree + 1` coefficients. Each value
// will be a random `big.Int` in range (0, q).
func generatePolynomial(degree int, dkg *DKG) ([]*big.Int, error) {
	coefficients := make([]*big.Int, degree+1)
	var err error
	for i := range coefficients {
		coefficients[i], err = dkg.RandomQ()
		if err != nil {
			return nil, err
		}
	}
	return coefficients, nil
}

func pow(x, y int) *big.Int {
	return big.NewInt(int64(math.Pow(float64(x), float64(y))))
}

// CombineReceivedShares sums up all shares received from peer group members.
//
// See Phase 6 of the protocol specification.
func (sm *SharingMember) CombineReceivedShares() {
	shareS := big.NewInt(0)
	for _, s := range sm.receivedSharesS {
		shareS = new(big.Int).Mod(
			new(big.Int).Add(shareS, s),
			sm.protocolConfig.Q,
		)
	}

	shareT := big.NewInt(0)
	for _, t := range sm.receivedSharesT {
		shareT = new(big.Int).Mod(
			new(big.Int).Add(shareT, t),
			sm.protocolConfig.Q,
		)
	}

	sm.shareS = shareS
	sm.shareT = shareT
}

// CalculatePublicCoefficients calculates public values for member's coefficients.
// It calculates `A_k = g^a_k mod p` for k in [0..T].
//
// See Phase 7 of the protocol specification.
func (sm *SharingMember) CalculatePublicCoefficients() *MemberPublicCoefficientsMessage {
	var publicCoefficients []*big.Int
	for _, a := range sm.secretCoefficients {
		publicA := new(big.Int).Exp(
			sm.vss.G,
			a,
			sm.protocolConfig.P,
		)
		publicCoefficients = append(publicCoefficients, publicA)
	}
	sm.publicCoefficients = publicCoefficients

	return &MemberPublicCoefficientsMessage{
		senderID:           sm.ID,
		publicCoefficients: publicCoefficients,
	}
}

// VerifyPublicCoefficients validates public key shares received in messages from
// peer group members.
// It returns accusation message with ID of members for which verification failed.
//
// See Phase 8 of the protocol specification.
func (sm *SharingMember) VerifyPublicCoefficients(messages []*MemberPublicCoefficientsMessage) (*CoefficientsAccusationsMessage, error) {
	var accusedMembersIDs []int
	// `product = Π (A_jk ^ (i^k)) mod p` for k in [0..T],
	// where: j is sender's ID, i is current member ID, T is threshold.
	for _, message := range messages {
		product := big.NewInt(1)
		for k, a := range message.publicCoefficients {
			product = new(big.Int).Mod(
				new(big.Int).Mul(
					product,
					new(big.Int).Exp(
						a,
						pow(sm.ID, k),
						sm.protocolConfig.P,
					),
				),
				sm.protocolConfig.P,
			)
		}
		// `expectedProduct = g^s_ji`
		expectedProduct := new(big.Int).Exp(
			sm.vss.G,
			sm.receivedSharesS[message.senderID],
			sm.protocolConfig.P)

		if expectedProduct.Cmp(product) != 0 {
			accusedMembersIDs = append(accusedMembersIDs, message.senderID)
		}
	}

	return &CoefficientsAccusationsMessage{
		senderID:   sm.ID,
		accusedIDs: accusedMembersIDs,
	}, nil
}
