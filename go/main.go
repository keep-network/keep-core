package main

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/dfinity/go-dfinity-crypto/bls"
	"github.com/dfinity/go-dfinity-crypto/groupsig"
)

type GroupMember struct {
	ID             bls.ID                     // globally shared
	Commitments    map[bls.ID][]bls.PublicKey // globally shared
	GroupPublicKey bls.PublicKey              // globally computed once sharing is complete
	// Created locally, shared only with corresponding member unless justifying
	// an accusation.
	Shares map[bls.ID]bls.SecretKey
	// Created locally, used to generate Shares, never shared.
	commitmentSecrets   []bls.SecretKey
	receivedShares      map[bls.ID]bls.SecretKey // received privately, stored locally
	GroupSecretKeyShare *bls.SecretKey           // computed from receivedShares, stored locally
	// Received via broadcast, tracks all accused IDs from this member's
	// perspective.
	accuserIDs []bls.ID
	// Received via broadcast, tracks IDs of all members who accused this
	// member.
	accusedIDs          map[bls.ID]bool
	disqualifiedPlayers map[bls.ID]bool // all players disqualified during justification
	qualifiedPlayers    []bls.ID        // a list of ids of all qualified players
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
func (member *GroupMember) generateCommitmentsAndShares(memberIDs []bls.ID, threshold int) []bls.PublicKey {
	commitments := make([]bls.PublicKey, threshold)
	member.commitmentSecrets = make([]bls.SecretKey, threshold)

	// Note: bls.SecretKey, before we call some sort of Set on it, can be
	// considered a zeroed *container* for a secret key.
	//
	// SetByCSPRNG initializes the zeroed secret key from a cryptographically
	// secure pseudo-random number generator.
	// Set instead initializes a key from an existing set of shares and a member
	// id.
	//
	// Commitmnent to s is E_0 = E(s, t) = g^s·h^t.
	// E_i = E(F_i, G_i)
	// F_i = coefficient i in F(x) = s + F_1·x + F_2·x^2 + ... + F_{k-1}·x^{k-1}
	// s_i = F(i)
	// G_i = coefficient i in G(x) = t + G_1·x + G_2·x^2 + ... + G_{k-1}·x^{k-1}
	// t_i = G(i)
	// Broadcast commitment is E_i = E(F_i, G_i) for i = 1, ..., k - 1
	//
	// [GJKR 99], Fig 2, 1(a).
	// For this dealer, i, we generate t secret keys, which are equivalent to t
	// coefficients a_ik and b_ik, k in [0,t], in two polynomials a and b (?),
	// and store them in secretShares. We also generate their equivalent public
	// keys, ~= g^{a_ik}·h^{b_ik}, which are inserted into the verification
	// vector. Note that g and h are the generators for the unique subgroups
	// chosen by the group via distributed coin flipping protocol (?)
	for i := 0; i < threshold; i++ {
		member.commitmentSecrets[i].SetByCSPRNG()
		commitments[i] = *member.commitmentSecrets[i].GetPublicKey()
	}

	// The public keys for each share of the member's secret key represent a
	// public commitment to the underlying secret key. Another member cannot get
	// the secret key or secret key shares from the public keys, but they can
	// use them to verify that the shares of the group secret key sent from this
	// member are valid against that underlying secret key.

	// Shamir secret sharing vs Pedersen VSS? Shamir's isn't verifiable.

	// Create a share of the group secret
	// For each member (including the caller!), we create a share from our set
	// of secret shares, ~= our polynomials. Equivalent to (s_ij, s'_ij), but
	// carried in the envelope of the secret key (similar to (a_ik, b_ik)).
	for _, memberID := range memberIDs {
		memberShare := bls.SecretKey{}
		memberShare.Set(member.commitmentSecrets, &memberID)
		member.Shares[memberID] = memberShare
	}

	return commitments
}

func (member GroupMember) isValidShare(shareHolderID bls.ID, share bls.SecretKey) bool {
	commitments := member.Commitments[shareHolderID]

	combinedCommitment := bls.PublicKey{}
	combinedCommitment.Set(commitments, &member.ID)

	comparisonShare := share.GetPublicKey()

	return combinedCommitment.IsEqual(comparisonShare)
}

func (member GroupMember) buildAccusations() []bls.ID {
	accusedIDs := make([]bls.ID, 0)
	for id, share := range member.receivedShares {
		commitments := member.Commitments[id]

		combinedCommitment := bls.PublicKey{}
		combinedCommitment.Set(commitments, &member.ID)

		comparisonShare := share.GetPublicKey()

		if !combinedCommitment.IsEqual(comparisonShare) {
			accusedIDs = append(accusedIDs, id)
		}
	}

	return accusedIDs
}

func (member *GroupMember) recordAccusations(from bls.ID, accusedIDs []bls.ID) {
	for _, id := range accusedIDs {
		if member.ID.IsEqual(&id) {
			member.accuserIDs = append(member.accuserIDs, member.ID)
		} else {
			member.accusedIDs[id] = true
		}
	}
}

func (member GroupMember) needsJustification() bool {
	return len(member.accuserIDs) > 0
}

func (member GroupMember) buildJustifications() map[bls.ID]bls.SecretKey {
	justifications := map[bls.ID]bls.SecretKey{}
	for _, accuserID := range member.accuserIDs {
		justifications[accuserID] = member.Shares[accuserID]
	}

	return justifications
}

func (member *GroupMember) considerJustification(justifyingID bls.ID, justification bls.SecretKey) {
	if !member.isValidShare(justifyingID, justification) {
		member.disqualifiedPlayers[justifyingID] = true
	}
}

// [GJKR 99], Fig 2, 2
func (member GroupMember) qualifiedShares() []bls.SecretKey {
	shares := make([]bls.SecretKey, 0)
	for id, share := range member.receivedShares {
		if !member.disqualifiedPlayers[id] {
			shares = append(shares, share)
		}
	}

	return shares
}

func (member *GroupMember) computeGroupKeyShares() {
	// [GJKR 99], Fig 2, 3
	initialShare := member.receivedShares[member.ID]
	member.GroupSecretKeyShare = &initialShare
	for id, share := range member.receivedShares {
		if !id.IsEqual(&member.ID) {
			member.GroupSecretKeyShare.Add(&share)
		}
	}

	// [GJKR 99], Fig 2, 4(c)? There is an accusation flow around public key
	//            			   computation as well...
	combinedCommitments := make([]bls.PublicKey, len(member.Commitments[member.ID]))
	for i, commitment := range member.Commitments[member.ID] {
		combinedCommitments[i] = commitment
	}
	for id, commitmentSet := range member.Commitments {
		if !id.IsEqual(&member.ID) { // we handled this above
			for i, commitment := range commitmentSet {
				combinedCommitments[i].Add(&commitment)
			}
		}
	}

	member.GroupPublicKey = combinedCommitments[0]
}

func createGroupMember() GroupMember {
	return GroupMember{
		Commitments:         map[bls.ID][]bls.PublicKey{},
		Shares:              map[bls.ID]bls.SecretKey{},
		receivedShares:      map[bls.ID]bls.SecretKey{},
		accuserIDs:          make([]bls.ID, 0),
		accusedIDs:          map[bls.ID]bool{},
		disqualifiedPlayers: map[bls.ID]bool{},
	}
}

func main() {
	fmt.Println("Starting!")
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

	memberNumbers := []int{1, 2, 3, 4, 5, 6, 7}
	memberIDs := make([]bls.ID, 0)
	membersByID := map[bls.ID]*GroupMember{}
	for _, number := range memberNumbers {
		member := createGroupMember()
		member.ID.SetDecString(fmt.Sprintf("%v", number))
		membersByID[member.ID] = &member
		memberIDs = append(memberIDs, member.ID)
	}

	type MemberCommitments struct {
		commitments []bls.PublicKey
	}
	type PrivateShare struct {
		encrypted  bool
		receiverID bls.ID
		share      bls.SecretKey
	}
	type Accusations struct {
		accusedIDs []bls.ID
	}
	type Justifications struct {
		// For one sender, justifications for each accusing member.
		justifications map[bls.ID]bls.SecretKey
	}

	encryptShare := func(share PrivateShare) PrivateShare {
		// Actually we just want to encrypt the share itself, not the whole
		// envelope, since we still need the receiver to know they're receiving
		// it.
		share.encrypted = true // 🙈🙉🙊

		return share
	}

	broadcast := func(senderID bls.ID, msg interface{}) {
		switch message := msg.(type) {
		case MemberCommitments:
			for _, member := range membersByID {
				member.Commitments[senderID] = message.commitments
			}

		case PrivateShare:
			membersByID[message.receiverID].receivedShares[senderID] = message.share

		case Accusations:
			for _, member := range membersByID {
				member.recordAccusations(senderID, message.accusedIDs)
			}

		case Justifications:
			// [GJKR 99], Fig 2, 1(d)
			// Validate justifications, build disqualified list. This should be happening
			// on each client, based on the same view of all accusations and
			// justifications.
			for accuserID, justification := range message.justifications {
				accused := membersByID[accuserID]
				accused.considerJustification(senderID, justification)
			}
		}
	}

	fmt.Println("Generating and broadcasting commitments and shares...")
	for _, dealer := range membersByID { // inner loop is in each client
		memberCommitments := dealer.generateCommitmentsAndShares(memberIDs, threshold)

		broadcast(dealer.ID, MemberCommitments{memberCommitments})
		for memberID, share := range dealer.Shares {
			// Timing attack-wise, we may want to encrypt all shares and then
			// ship them?
			broadcast(dealer.ID, encryptShare(PrivateShare{false, memberID, share}))
		}
	}

	// In network comms, wait until all commitments and shares are received.

	fmt.Println("Handling accusations...")
	// Validate and accuse: each member broadcasts accusations against a member
	// whose share they failed to validate.
	// [GJKR 99], Fig 2, 1(b)
	for _, member := range membersByID {
		accusations := member.buildAccusations()
		broadcast(member.ID, Accusations{accusations})
	}

	type justification struct {
		index int
		proof bls.SecretKey
	}

	fmt.Println("Handling justifications...")
	for _, member := range membersByID {
		if member.needsJustification() {
			// [GJKR 99], Fig 2, 1(c)
			// Justify against accusations; optional, can also just immediately
			// fail if there are any accusations.
			// Justifications: for each accused ID, we broadcast the accused's
			// secret key share.
			broadcast(member.ID, Justifications{member.buildJustifications()})
		}
	}

	fmt.Println("Computing key shares...")
	// [GJKR 99], Fig 2, 3
	// [GJKR 99], Fig 2, 4(c)? There is an accusation flow around public key
	//            			   computation as well...
	// Key computation.
	for _, member := range membersByID {
		// NOTE: We now have a group public key based on these shares.
		member.computeGroupKeyShares()
	}

	fmt.Println("Validating group public keys...")
	var key *bls.PublicKey
	for _, member := range membersByID {
		if key == nil {
			key = &member.GroupPublicKey
		} else if !key.IsEqual(&member.GroupPublicKey) {
			panic(fmt.Sprintf("Public key for %v is bad.", key))
		}
	}

	fmt.Println("Generating signatures...")
	messageToSign := "This is a test message!"
	randomizedMemberIndices := rand.Perm(len(memberIDs))
	thresholdSignatures := make([]bls.Sign, threshold)
	signerIDs := make([]bls.ID, threshold)
	for i := 0; i < threshold; i++ {
		memberID := memberIDs[randomizedMemberIndices[i]]
		signerIDs[i] = memberID
		thresholdSignatures[i] = *membersByID[memberID].GroupSecretKeyShare.Sign(messageToSign)
	}

	finalSignature := &bls.Sign{}
	finalSignature.Recover(thresholdSignatures, signerIDs)
	if finalSignature.Verify(key, messageToSign) {
		fmt.Println("Verified!")
		os.Exit(0)
	} else {
		fmt.Println("Verification failed :(")
		os.Exit(1)
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
}
