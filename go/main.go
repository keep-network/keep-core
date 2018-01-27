package main

import (
	"fmt"

	"github.com/dfinity/go-dfinity-crypto/bls"
	"github.com/dfinity/go-dfinity-crypto/groupsig"
)

type GroupMember struct {
	ID             bls.ID            // globally shared
	Commitments    [][]bls.PublicKey // globally shared
	GroupPublicKey bls.PublicKey     // globally computed once sharing is complete
	// created locally, shared only with
	// corresponding member unless justifying
	// an accusation
	shares              []bls.SecretKey
	receivedShares      []bls.SecretKey // received privately, stored locally
	groupSecretKeyShare bls.SecretKey   // computed from receivedShares, stored locally
}

func idsFromGroupMembers(members []GroupMember) []*bls.ID {
	memberIds := make([]*bls.ID, 0)
	for _, member := range members {
		memberIds = append(memberIds, &member.ID)
	}
	return memberIds
}

// Note: Each player generates as a set of shares. Each player's secret key
// shares can combine into a single secret key for the player. Each
// player's secret key plays a part in generating shares of the group's
// secret key. Each player's share of the group secret key is assembled by
// combining the privately-communicated shares of all player's individual
// secret keys.
//
// P_0 -> S_0 split into S_00, S_01, S_02; C_00, C_01: commitments to S_0
// P_1 -> S_1 split into S_10, S_11, S_12; C_10, C_11: commitments to S_1
// P_2 -> S_2 split into S_20, S_21, S_22; C_20, C_21: commitments to S_2
//
// NOTE: commitments are from 0 to threshold, not from 0 to number of players!
//
// S_0 used to generate s_00, s_01, s_02, the shares of S_0 specifically for P_0, P_1, P_2.
// S_1 used to generate s_10, s_11, s_12, the shares of S_1 specifically for P_0, P_1, P_2.
// S_2 used to generate s_20, s_21, s_22, the shares of S_2 specifically for P_0, P_1, P_2.
//
// P_0 verifies s_00 against C_00, C_01
//     verifies s_10 against C_10, C_11
//     verifies s_20 against C_20, C_21
// P_1 verifies s_01 against C_00, C_01
//     verifies s_11 against C_10, C_11
//     verifies s_21 against C_20, C_21
// P_2 verifies s_02 against C_00, C_01
//     verifies s_12 against C_10, C_11
//     verifies s_22 against C_20, C_21
//
// If P_0 fails to verify s_10, broadcast accusation against P_1.
//                 	      s_20, broadcast accusation against P_2.
// If P_1 fails to verify s_01, broadcast accusation against P_0.
//                 	      s_21, broadcast accusation against P_2.
// If P_2 fails to verify s_02, broadcast accusation against P_0.
//                 	      s_12, broadcast accusation against P_1.
//
// If P_0 receives accusation from P_1, broadcast s_01.
//                                 P_2, broadcast s_02.
// If P_1 receives accusation from P_0, broadcast s_10.
//                                 P_2, broadcast s_12.
// If P_2 receives accusation from P_0, broadcast s_20.
//                                 P_1, broadcast s_21.
//
// If P_0 confirms invalid s_10, disqualify P_1 locally.
//                         s_20, disqualify P_2 locally.
// If P_1 confirms invalid s_01, disqualify P_0 locally.
//                         s_21, disqualify P_2 locally.
// If P_2 confirms invalid s_02, disqualify P_0 locally.
//                         s_12, disqualify P_1 locally.
//
// Everyone has an equal set QUAL of qualified P.

// Each player generates a commitment, which consists of a set of secret
// shares, which are used in turn to generate a verification vector alongside a
// set of contributions, one per member in the system.
//
// Note that there are t verifications, where t is the threshold of honest
// players. However, there are n contributions, where n is the total number of
// players.
func generateCommitmentsAndShares(memberIDs []*bls.ID, threshold int) ([]bls.SecretKey, []bls.PublicKey, []*bls.SecretKey) {
	commitments := make([]bls.PublicKey, threshold)
	dealerShares := make([]bls.SecretKey, len(memberIDs))
	memberShares := make([]*bls.SecretKey, 0)

	// Note: bls.SecretKey, before we call some sort of Set on it, can be
	// considered a zeroed *container* for a secret key.
	//
	// SetByCSPRNG initializes the zeroed secret key from a cryptographically
	// secure pseudo-random number generator.
	// Set instead initializes a key from an existing set of shares and a member
	// id.
	//
	// [GJKR 99], Fig 2, 1(a).
	// For this dealer, i, we generate t secret keys, which are equivalent to t
	// coefficients a_ik and b_ik, k in [0,t], in two polynomials a and b (?),
	// and store them in secretShares. We also generate their equivalent public
	// keys, ~= g^{a_ik}·h^{b_ik}, which are inserted into the verification
	// vector. Note that g and h are the generators for the unique subgroups
	// chosen by the group via distributed coin flipping protocol (?)
	for i := 0; i < threshold; i++ {
		dealerShares[i].SetByCSPRNG()
		commitments[i] = *dealerShares[i].GetPublicKey()
	}

	// The public keys for each share of the member's secret key represent a
	// public commitment to the underlying secret key. Another member cannot get
	// the secret key or secret key shares from the public keys, but they can
	// use them to verify that the shares of the group secret key sent from this
	// member are valid against that underlying secret key.

	// Shamir secret sharing vs Pedersen VSS? Shamir's isn't verifiable,
	// perhaps?

	// Create a share of the group secret
	// For each member (including the caller!), we create a share from our set
	// of secret shares, ~= our polynomials. Equivalent to (s_ij, s'_ij), but
	// carried in the envelope of the secret key (similar to (a_ik, b_ik)).
	for _, memberID := range memberIDs {
		memberShare := bls.SecretKey{}
		memberShare.Set(dealerShares, memberID)
		memberShares = append(memberShares, &memberShare)
	}

	return dealerShares, commitments, memberShares
}

func (member GroupMember) isValidShare(shareHolderIndex int, share bls.SecretKey) bool {
	commitments := member.Commitments[shareHolderIndex]

	combinedCommitment := bls.PublicKey{}
	combinedCommitment.Set(commitments, &member.ID)

	comparisonShare := share.GetPublicKey()

	return combinedCommitment.IsEqual(comparisonShare)
}

func (member GroupMember) invalidShares() []int {
	invalidShares := make([]int, 0)
	for j, share := range member.receivedShares {
		commitments := member.Commitments[j]

		combinedCommitment := bls.PublicKey{}
		combinedCommitment.Set(commitments, &member.ID)

		comparisonShare := share.GetPublicKey()

		if !combinedCommitment.IsEqual(comparisonShare) {
			invalidShares = append(invalidShares, j)
		}
	}

	return invalidShares
}

func (member GroupMember) computeGroupKeyShares() {
	// [GJKR 99], Fig 2, 3
	member.groupSecretKeyShare = bls.SecretKey{}
	member.groupSecretKeyShare.Set(member.receivedShares, &member.ID)

	// [GJKR 99], Fig 2, 4(c)? There is an accusation flow around public key
	//            			   computation as well...
	combinedCommitments := make([]bls.PublicKey, len(member.Commitments[0]))
	for i, commitment := range member.Commitments[0] {
		combinedCommitments[i].Deserialize(commitment.Serialize())
	}
	for _, commitmentSet := range member.Commitments[1:] {
		for i, commitment := range commitmentSet {
			combinedCommitments[i].Add(&commitment)
		}
	}

	member.GroupPublicKey = combinedCommitments[0]
}

func main() {
	fmt.Printf("Starting!")
	groupsig.Init(bls.CurveFp382_1)

	// Network substrate: need to be able to:
	//  - Encryptedly send messages to another staker in the group, via their
	//    group key in registry (possibility: a _derivation_ based on that public
	//    key).
	//  - Send all messages with an HMAC or other signature that can prove the
	//    message is in fact from the sender.
	//  - Messages from the same sender arrive in order.
	//  - Message loss?
	//
	// Want to be able to:
	//  - Use bls.ID for everything (e.g., C_ij, i and j are actually bls.ID).
	/// - Alternatively, everyone needs to keep an ordered list of members,
	//    the output of our sortition.

	threshold := 4

	memberNumbers := []int{0, 1, 2, 3, 4, 5, 6}
	members := []GroupMember{}
	for _, number := range memberNumbers {
		member := GroupMember{}
		member.ID.SetDecString(string(number))
		members = append(members, member)
	}

	for dealerI, dealer := range members { // inner loop is per client
		ownShares, memberCommitments, memberShares :=
			generateCommitmentsAndShares(idsFromGroupMembers(members), threshold)

		dealer.shares = ownShares

		// Broadcast commitments to all members and send them each their share
		// of this member's secret.
		for playerJ, member := range members {
			member.Commitments = append(member.Commitments, memberCommitments)
			// In network comms, this needs to be encrypted so it's only visible
			// by `member`.
			member.receivedShares[dealerI] = *memberShares[playerJ]
		}
	}

	// In network comms, wait until all commitments and shares are received.

	// Validate and accuse: each member broadcasts accusations against a member
	// whose share they failed to validate.
	// [GJKR 99], Fig 2, 1(b)
	allAccusationIndices := make([][]int, 0)  // "broadcast channel"
	for _, validatorMember := range members { // inner loop is per client
		accusationIndices := make([]int, 0)
		for _, invalidShareIndex := range validatorMember.invalidShares() {
			accusationIndices = append(accusationIndices, invalidShareIndex)
		}

		allAccusationIndices = append(allAccusationIndices, accusationIndices)
	}

	type justification struct {
		index int
		proof bls.SecretKey
	}

	// [GJKR 99], Fig 2, 1(c)
	// Justify against accusations; optional, can also just immediately fail if
	// there are any accusations.
	// Justifications: for each accused index, for each accuser index, we
	// broadcast the accused's secret key share.
	allJustifications := make([][]justification, 0) // "broadcast channel"
	// Handle accusations. This should be happening on each client, based on the
	// same view of all accusations.
	for accuserI, allAccused := range allAccusationIndices {
		accuserJustifications := make([]justification, len(allAccused))
		for accusedJ, accusedIndex := range allAccused {
			accused := members[accusedIndex]
			accuserJustifications[accusedJ] = justification{accusedIndex, accused.shares[accuserI]}
		}
		allJustifications = append(allJustifications, accuserJustifications)
	}

	// [GJKR 99], Fig 2, 1(d)
	// Validate justifications, build disqualified list. This should be happening
	// on each client, based on the same view of all accusations and
	// justifications.
	disqualifiedMemberIndices := make([]int, 0)
	for accuserI, justifications := range allJustifications {
		memberView := members[accuserI] // this client's view of member i
		for justificationI, justification := range justifications {
			if !memberView.isValidShare(justificationI, justification.proof) {
				disqualifiedMemberIndices = append(disqualifiedMemberIndices, justification.index)
			}
		}
	}

	// [GJKR 99], Fig 2, 2
	// Build qualified list.
	qualifiedMembers := make([]GroupMember, len(members)-len(disqualifiedMemberIndices))
	for i, member := range members {
		disqualified := false
		for index := range disqualifiedMemberIndices {
			if index == i {
				disqualified = true
				break
			}
		}

		if !disqualified {
			qualifiedMembers = append(qualifiedMembers, member)
		}
	}

	// [GJKR 99], Fig 2, 3
	// [GJKR 99], Fig 2, 4(c)? There is an accusation flow around public key
	//            			   computation as well...
	// Key computation.
	for _, member := range qualifiedMembers {
		// NOTE: We now have a group public key based on these shares.
		member.computeGroupKeyShares()
	}

	// Public key submission: what happens if a bad key is submitted?
	//  - Group won't be able to sign, bad key holder will.
	//  - How do we validate the key submitted is indeed for the group in
	//    question?
	//  - What prevents one client from submitting their own private/public key
	//    as the group key at the very beginning?
	//  -> Aggregate signature from the group of the public key message?
	//     -> Each member generates the public key, and broadcasts their
	//     -> sig of the key.
	//     -> Key publishing publishes key + group aggregated signature of key.

	// =~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=

	message := []byte("Booyan booyanescu")
	seckeys := []groupsig.Seckey{
		*groupsig.NewSeckeyFromInt(1),
		*groupsig.NewSeckeyFromInt(2),
		*groupsig.NewSeckeyFromInt(3),
		*groupsig.NewSeckeyFromInt(4),
		*groupsig.NewSeckeyFromInt(5),
	}
	pubkeys := []groupsig.Pubkey{
		*groupsig.NewPubkeyFromSeckey(seckeys[0]),
		*groupsig.NewPubkeyFromSeckey(seckeys[1]),
		*groupsig.NewPubkeyFromSeckey(seckeys[2]),
		*groupsig.NewPubkeyFromSeckey(seckeys[3]),
		*groupsig.NewPubkeyFromSeckey(seckeys[4]),
	}
	signatures := []groupsig.Signature{
		groupsig.Sign(seckeys[0], message),
		groupsig.Sign(seckeys[1], message),
		groupsig.Sign(seckeys[2], message),
		groupsig.Sign(seckeys[3], message),
		groupsig.Sign(seckeys[4], message),
	}

	master := groupsig.AggregateSeckeys(seckeys)
	masterPub := groupsig.NewPubkeyFromSeckey(*master)
	signature := groupsig.Sign(*master, message)
	aggregatedSig := groupsig.AggregateSigs(signatures)

	verification := groupsig.VerifySig(*masterPub, message, signature) == groupsig.VerifySig(*masterPub, message, aggregatedSig)
	aggregateVerification := groupsig.VerifyAggregateSig(pubkeys, message, signature)
	batchVerification := groupsig.BatchVerify(pubkeys, message, signatures)

	fmt.Printf(
		"%v = %v\nVerified: %v\nVerified in aggregate: %v\nVerified by batch: %v\n",
		signature.GetHexString(),
		aggregatedSig.GetHexString(),
		verification,
		aggregateVerification,
		batchVerification,
	)
}
