// Package gjkr contains code that implements Distributed Key Generation protocol
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
		return nil, nil, fmt.Errorf("cannot generate polynomial [%v]", err)
	}
	coefficientsB, err := generatePolynomial(polynomialDegree, cm.protocolConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot generate hiding polynomial [%v]", err)
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

	for _, commitmentsMessage := range commitmentsMessages {
		// Find share message sent by the same member who sent commitment message
		sharesMessageFound := false
		for _, sharesMessage := range sharesMessages {
			if sharesMessage.senderID == commitmentsMessage.senderID {
				sharesMessageFound = true

				// Check if `commitmentsProduct == expectedProduct`
				// `commitmentsProduct = Π (C_j[k] ^ (i^k)) mod p` for k in [0..T]
				// `expectedProduct = (g ^ s_ji) * (h ^ t_ji) mod p`
				// where: j is sender's ID, i is current member ID, T is threshold.
				if !cm.areSharesValidAgainstCommitments(
					sharesMessage.shareS, sharesMessage.shareT, // s_ji, t_ji
					commitmentsMessage.commitments, // C_j
					cm.ID,                          // i
				) {
					accusedMembersIDs = append(accusedMembersIDs,
						commitmentsMessage.senderID)
					break
				}
				cm.receivedValidSharesS[commitmentsMessage.senderID] = sharesMessage.shareS
				cm.receivedValidSharesT[commitmentsMessage.senderID] = sharesMessage.shareT
				cm.receivedValidPeerCommitments[commitmentsMessage.senderID] = commitmentsMessage.commitments
				break
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

// areSharesValidAgainstCommitments verifies if commitments are valid for passed
// shares.
//
// The `j` member generated a polynomial with `k` coefficients before. Then they
// calculated a commitments to the polynomial's coefficients `C_j` and individual
// shares `s_ji` and `t_ji` with a polynomial for a member `i`. In this function
// the verifier checks if the shares are valid against the commitments.
//
// The verifier checks that:
// `commitmentsProduct == expectedProduct`
// where:
// `commitmentsProduct = Π (C_j[k] ^ (i^k)) mod p` for k in [0..T],
// and
// `expectedProduct = (g ^ s_ji) * (h ^ t_ji) mod p`:
func (cm *CommittingMember) areSharesValidAgainstCommitments(
	shareS, shareT *big.Int, // s_ji, t_ji
	commitments []*big.Int, // C_j
	memberID int, // i
) bool {
	// `commitmentsProduct = Π (C_j[k] ^ (i^k)) mod p`
	commitmentsProduct := big.NewInt(1)
	for k, c := range commitments {
		commitmentsProduct = new(big.Int).Mod(
			new(big.Int).Mul(
				commitmentsProduct,
				new(big.Int).Exp(
					c,
					pow(memberID, k),
					cm.protocolConfig.P,
				),
			),
			cm.protocolConfig.P,
		)
	}

	// `expectedProduct = (g ^ s_ji) * (h ^ t_ji) mod p`, where:
	expectedProduct := cm.vss.CalculateCommitment(
		shareS,
		shareT,
		cm.protocolConfig.P,
	)

	return expectedProduct.Cmp(commitmentsProduct) == 0
}

// ResolveSecretSharesAccusations resolves a complaint received from a sender
// against a member accused in the shares and commitments verification phase.
// A member is calling this function to judge which party of the dispute is lying.
//
// The function requires shares `s_mj` and `t_mj` calculated by the accused
// member (`m`) for the sender (`j`). These values are expected to be broadcast
// before in encrypted form. On accusation, the shares should be decrypted and
// the revealed value should be passed to this function.
//
// A current member cannot be a part of a dispute. If the current member is
// either an accuser or is accused the function will return an error. The accused
// party cannot be a judge in its own case. From the other hand, the accuser has
// already performed the calculation in the previous phase which resulted in the
// accusation and waits now for a judgment from other players.
//
// The returned value is an ID of the member who should be slashed. It will be
// an accuser ID if the validation shows that shares and commitments are valid,
// so the accusation was unfounded. Else it confirms that accused member misbehaved
// and their ID is returned.
//
// See Phase 5 of the protocol specification.
func (sjm *SharesJustifyingMember) ResolveSecretSharesAccusations(
	senderID, accusedID int, // j, m
	shareS, shareT *big.Int, // s_mj, t_mj
) (int, error) {
	if sjm.ID == senderID || sjm.ID == accusedID {
		return 0, fmt.Errorf("current member cannot be a part of a dispute")
	}

	// Check if `commitmentsProduct == expectedProduct`
	// `commitmentsProduct = Π (C_m[k] ^ (j^k)) mod p` for k in [0..T]
	// `expectedProduct = (g ^ s_mj) * (h ^ t_mj) mod p`
	// where: m is accused member's ID, j is sender's ID, T is threshold.
	if sjm.areSharesValidAgainstCommitments(
		shareS, shareT, // s_mj, t_mj
		sjm.receivedValidPeerCommitments[accusedID], // C_m
		senderID, // j
	) {
		return senderID, nil
	}
	return accusedID, nil
}

// CombineMemberShares sums up all `s` and `t` shares intended for this member.
// Combines secret shares calculated by current member `i` for itself `s_ii` with
// shares calculated by peer members `j` for this member `s_ji`.
//
// `x_i = Σ s_ji mod q` and `x'_i = Σ t_ji mod q` for `j` in a group of players
// who passed secret shares accusations stage.
//
// See Phase 6 of the protocol specification.
func (qm *QualifiedMember) CombineMemberShares() {
	combinedSharesS := qm.selfSecretShareS // s_ii
	for _, s := range qm.receivedValidSharesS {
		combinedSharesS = new(big.Int).Mod(
			new(big.Int).Add(combinedSharesS, s),
			qm.protocolConfig.Q,
		)
	}

	combinedSharesT := qm.selfSecretShareT // t_ii
	for _, t := range qm.receivedValidSharesT {
		combinedSharesT = new(big.Int).Mod(
			new(big.Int).Add(combinedSharesT, t),
			qm.protocolConfig.Q,
		)
	}

	qm.masterPrivateKeyShare = combinedSharesS
	qm.shareT = combinedSharesT
}

// CalculatePublicKeySharePoints calculates public values for member's coefficients.
// It calculates `A_k = g^a_k mod p` for k in [0..T].
//
// See Phase 7 of the protocol specification.
func (sm *SharingMember) CalculatePublicKeySharePoints() *MemberPublicKeySharePointsMessage {
	sm.publicKeySharePoints = make([]*big.Int, len(sm.secretCoefficients))
	for i, a := range sm.secretCoefficients {
		sm.publicKeySharePoints[i] = new(big.Int).Exp(
			sm.vss.G,
			a,
			sm.protocolConfig.P,
		)
	}

	return &MemberPublicKeySharePointsMessage{
		senderID:             sm.ID,
		publicKeySharePoints: sm.publicKeySharePoints,
	}
}

// VerifyPublicKeySharePoints validates public key share points received in
// messages from peer group members.
// It returns accusation message with ID of members for which the verification
// failed.
//
// See Phase 8 of the protocol specification.
func (sm *SharingMember) VerifyPublicKeySharePoints(
	messages []*MemberPublicKeySharePointsMessage,
) (*PointsAccusationsMessage, error) {
	var accusedMembersIDs []int
	// `product = Π (A_jk ^ (i^k)) mod p` for k in [0..T],
	// where: j is sender's ID, i is current member ID, T is threshold.
	for _, message := range messages {
		product := big.NewInt(1)
		for k, a := range message.publicKeySharePoints {
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
			sm.receivedValidSharesS[message.senderID],
			sm.protocolConfig.P)

		if expectedProduct.Cmp(product) != 0 {
			accusedMembersIDs = append(accusedMembersIDs, message.senderID)
			continue
		}
		sm.receivedValidPeerPublicKeySharePoints[message.senderID] = message.publicKeySharePoints
	}

	return &PointsAccusationsMessage{
		senderID:   sm.ID,
		accusedIDs: accusedMembersIDs,
	}, nil
}

// ResolvePublicKeySharePointsAccusations resolves a complaint received from a sender
// against a member accused in public key share points verification.
//
// Current member cannot be a part of a dispute, if the member is either a sender
// or accused the function will return an error.
//
// The function requires share `s_mj` calculated by the accused member (`m`) for
// the sender (`j`). This value was shared privately by member `m` with member
// `j` in the previous phase. On accusation, this value is revealed publicly to
// resolve the dispute between `m` and `j` and is an input parameter to this function.
//
// The returned value is an ID of the member who should be slashed. It will be
// an accuser ID if the validation shows that coefficients are valid, so the
// accusation was unfounded. Else it confirms that accused member misbehaved
// and their ID is returned.
//
// See Phase 9 of the protocol specification.
func (cjm *PointsJustifyingMember) ResolvePublicKeySharePointsAccusations(
	senderID, accusedID int,
	shareS *big.Int,
) (int, error) {
	if cjm.ID == senderID || cjm.ID == accusedID {
		return 0, fmt.Errorf("current member cannot be a part of a dispute")
	}

	// `product = Π (A_mk ^ (j^k)) mod p` for k in [0..T],
	// where: m is accused member's ID, j is sender's ID, T is threshold.
	product := big.NewInt(1)
	for k, a := range cjm.receivedValidPeerPublicKeySharePoints[accusedID] {
		product = new(big.Int).Mod(
			new(big.Int).Mul(
				product,
				new(big.Int).Exp(
					a,
					pow(senderID, k),
					cjm.protocolConfig.P,
				),
			),
			cjm.protocolConfig.P,
		)
	}

	// `expectedProduct = g^s_mj mod p`, where:
	// m is accused member's ID, j is sender's ID.
	expectedProduct := new(big.Int).Exp(
		cjm.vss.G,
		shareS,
		cjm.protocolConfig.P,
	)

	if expectedProduct.Cmp(product) == 0 {
		return senderID, nil
	}
	return accusedID, nil
}
