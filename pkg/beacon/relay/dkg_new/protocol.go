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
	"github.com/keep-network/keep-core/pkg/beacon/relay/pedersen"
)

// CalculateSharesAndCommitments starts with generating coefficients for two
// polynomials. Then it is calculating shares for all group member and packing
// them in individual messages for each peer member. Additionally, it calculates
// commitments to `a` coefficients of first polynomial using second's polynomial
// `b` coefficients.
//
// See http://docs.keep.network/cryptography/beacon_dkg.html#_phase_3_polynomial_generation
func (cm *CommittingMember) CalculateSharesAndCommitments() ([]*PeerSharesMessage, *MemberCommitmentsMessage, error) {
	var err error
	coefficientsSize := cm.group.dishonestThreshold + 1
	coefficientsA := make([]*big.Int, coefficientsSize)
	coefficientsB := make([]*big.Int, coefficientsSize)
	for j := 0; j < coefficientsSize; j++ {
		coefficientsA[j], err = rand.Int(rand.Reader, cm.ProtocolConfig().Q)
		if err != nil {
			return nil, nil, err
		}
		coefficientsB[j], err = rand.Int(rand.Reader, cm.ProtocolConfig().Q)
		if err != nil {
			return nil, nil, err
		}
	}

	cm.coefficientsA = coefficientsA

	// Calculate shares for other group members by evaluating polynomials defined
	// by coefficients `a_i` and `b_i`
	var sharesMessages []*PeerSharesMessage
	for _, id := range cm.group.MemberIDs() {
		// s_j = f_(j) mod q
		secretShare := calculateShare(id, coefficientsA, cm.ProtocolConfig().Q)
		// t_j = g_(j) mod q
		randomShare := calculateShare(id, coefficientsB, cm.ProtocolConfig().Q)

		// Check if calculated shares for the current member. If true store them
		// without sharing in a message.
		if cm.ID.Cmp(id) == 0 {
			cm.secretShares[cm.ID] = secretShare
			cm.randomShares[cm.ID] = randomShare
			continue
		}

		sharesMessages = append(sharesMessages,
			&PeerSharesMessage{
				senderID:    cm.ID,
				receiverID:  id,
				secretShare: secretShare,
				randomShare: randomShare,
			})
	}

	commitments := make([]*big.Int, coefficientsSize)
	for k := 0; k < coefficientsSize; k++ {
		// C_k = g^a_k * h^b_k mod p
		commitments[k] = pedersen.CalculateCommitment(cm.vss, coefficientsA[k], coefficientsB[k])
	}
	commitmentsMessage := &MemberCommitmentsMessage{
		senderID:    cm.ID,
		commitments: commitments,
	}

	return sharesMessages, commitmentsMessage, nil
}

// VerifySharesAndCommitments verifies shares and commitments received in messages
// from peer group members.
// It returns accusation message with ID of members for which verification failed.
//
// If cannot match commitments message with shares message for given sender then
// error is returned.
//
// See http://docs.keep.network/cryptography/beacon_dkg.html#_phase_4_share_verification
func (cm *CommittingMember) VerifySharesAndCommitments(
	sharesMessages []*PeerSharesMessage,
	commitmentsMessages []*MemberCommitmentsMessage,
) (*FirstAccusationsMessage, error) {
	var accusedMembersIDs []*big.Int
	// `commitmentsProduct = Π (commitments_j[k] ^ (i^k)) mod p` for k in [0..T],
	// where: j is sender's ID, i is current member ID, T is threshold.

	for _, commitmentMessage := range commitmentsMessages {
		commitmentsProduct := big.NewInt(1)
		for k, c := range commitmentMessage.commitments {
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
		// Find share message sent by the same member who sent commitment message
		shareMessageFound := false
		for _, shareMessage := range sharesMessages {
			if shareMessage.senderID.Cmp(commitmentMessage.senderID) == 0 {
				shareMessageFound = true
				// `expectedProduct = (g ^ s_j) * (h ^ t_j)`
				expectedProduct := pedersen.CalculateCommitment(cm.vss, shareMessage.secretShare, shareMessage.randomShare)

				if expectedProduct.Cmp(commitmentsProduct) != 0 {
					accusedMembersIDs = append(accusedMembersIDs, commitmentMessage.senderID)
					break
				}
				// Phase 6
				cm.secretShares[commitmentMessage.senderID] = shareMessage.secretShare
				cm.randomShares[commitmentMessage.senderID] = shareMessage.randomShare
			}
		}
		if !shareMessageFound {
			return nil, fmt.Errorf("cannot find shares message from member %s", commitmentMessage.senderID)
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

// CombineReceivedShares sums up all shares received from peer group members.
//
// See http://docs.keep.network/cryptography/beacon_dkg.html#_phase_6_share_calculation
func (sm *SharingMember) CombineReceivedShares() {
	secretShare := big.NewInt(0)
	for _, s := range sm.receivedSecretShares {
		secretShare = new(big.Int).Mod(
			new(big.Int).Add(secretShare, s),
			sm.ProtocolConfig().Q,
		)
	}

	randomShare := big.NewInt(0)
	for _, t := range sm.receivedRandomShares {
		randomShare = new(big.Int).Mod(
			new(big.Int).Add(randomShare, t),
			sm.ProtocolConfig().Q,
		)
	}

	sm.privateKeyShare = secretShare
	sm.privateRandomShare = randomShare
}

// CalculatePublicKeyShares calculates public values for member's shares.
// It calculates `A_k = g^a_k mod p` for k in [0..T].
//
// See http://docs.keep.network/cryptography/beacon_dkg.html#_phase_7_public_key_share_points
func (sm *SharingMember) CalculatePublicKeyShares() *MemberPublicKeySharesMessage {
	var publicShares []*big.Int
	for _, share := range sm.secretShares {
		publicShare := new(big.Int).Exp(
			sm.vss.G(),
			share,
			sm.ProtocolConfig().P,
		)
		publicShares = append(publicShares, publicShare)
	}
	sm.publicShares = publicShares

	return &MemberPublicKeySharesMessage{
		senderID:        sm.ID,
		publicKeyShares: publicShares,
	}
}

// VerifyPublicKeyShares validates public key shares received in messages from
// peer group members.
// It returns accusation message with ID of members for which verification failed.
//
// See http://docs.keep.network/cryptography/beacon_dkg.html#_phase_8_public_key_share_validation
func (sm *SharingMember) VerifyPublicKeyShares(messages []*MemberPublicKeySharesMessage) (*SecondAccusationsMessage, error) {
	var accusedMembersIDs []*big.Int
	// `product = Π (A_jk ^ (i^k)) mod p` for k in [0..T],
	// where: j is sender's ID, i is current member ID, T is threshold.
	for _, message := range messages {
		product := big.NewInt(1)
		for k, c := range message.publicKeyShares {
			product = new(big.Int).Mod(
				new(big.Int).Mul(
					product,
					new(big.Int).Exp(
						c,
						new(big.Int).Exp(
							sm.ID,
							big.NewInt(int64(k)),
							sm.ProtocolConfig().P,
						),
						sm.ProtocolConfig().P,
					),
				),
				sm.ProtocolConfig().P,
			)
		}
		// `expectedProduct = g^s_ji`
		expectedProduct := new(big.Int).Exp(
			sm.vss.G(),
			sm.receivedSecretShares[message.senderID],
			sm.ProtocolConfig().P)

		if expectedProduct.Cmp(product) != 0 {
			accusedMembersIDs = append(accusedMembersIDs, message.senderID)
		}
	}

	return &SecondAccusationsMessage{
		senderID:   sm.ID,
		accusedIDs: accusedMembersIDs,
	}, nil
}
