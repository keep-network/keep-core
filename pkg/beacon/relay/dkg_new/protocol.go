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

	"github.com/golang/go/src/crypto/rand"
)

// CalculateMembersSharesAndCommitments starts with generating coefficients for two
// polynomials. Then it is calculating shares for all group member and packing
// them in individual messages for each peer member. Additionally, it calculates
// commitments to `a` coefficients of first polynomial using second's polynomial
// `b` coefficients.
//
// See http://docs.keep.network/cryptography/beacon_dkg.html#_phase_3_polynomial_generation
func (cm *CommittingMember) CalculateMembersSharesAndCommitments() ([]*PeerSharesMessage, *MemberCommitmentsMessage, error) {
	var err error
	coefficientsSize := cm.group.dishonestThreshold + 1
	coefficientsA := make([]*big.Int, coefficientsSize)
	coefficientsB := make([]*big.Int, coefficientsSize)
	for j := 0; j < coefficientsSize; j++ {
		coefficientsA[j], err = rand.Int(rand.Reader, cm.protocolConfig.Q)
		if err != nil {
			return nil, nil, err
		}
		coefficientsB[j], err = rand.Int(rand.Reader, cm.protocolConfig.Q)
		if err != nil {
			return nil, nil, err
		}
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

	commitments := make([]*big.Int, coefficientsSize)
	for k := 0; k < coefficientsSize; k++ {
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

// VerifyReceivedSharesAndCommitmentsMessages verifies shares and commitments received in
// messages from peer group members.
// It returns accusation message with ID of members for which verification failed.
//
// If cannot match commitments message with shares message for given sender then
// error is returned.
//
// See http://docs.keep.network/cryptography/beacon_dkg.html#_phase_4_share_verification
func (cm *CommittingMember) VerifyReceivedSharesAndCommitmentsMessages(
	sharesMessages []*PeerSharesMessage,
	commitmentsMessages []*MemberCommitmentsMessage,
) (*FirstAccusationsMessage, error) {
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
			return nil, fmt.Errorf("cannot find shares message from member %s", commitmentsMessage.senderID)
		}
	}

	return &FirstAccusationsMessage{
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

// CombinePublicKeyShares calculates group public key from public key shares.
// Public key is calculated as a product of zeroth public shares (coefficients)
// `A_j0 = z_j` for all group members including member themself.
func (sm *SharingMember) CombinePublicKeyShares() {
	// Member's zeroth coefficient.
	memberGroupPublicKeyShare := sm.publicCoefficientsA[0]
	groupPublicKey := memberGroupPublicKeyShare
	// Multiply peer group members' zeroth coefficients.
	for _, publicKeyShare := range sm.receivedGroupPublicKeyShares {
		groupPublicKey = new(big.Int).Mod(
			new(big.Int).Mul(groupPublicKey, publicKeyShare),
			sm.ProtocolConfig().P,
		)
	}
	sm.groupPublicKey = groupPublicKey
}
