package thresholdgroup

import (
	"github.com/dfinity/go-dfinity-crypto/bls"
)

// [GJKR 99]: Gennaro R., Jarecki S., Krawczyk H., Rabin T. (1999) Secure
//     Distributed Key Generation for Discrete-Log Based Cryptosystems. In:
//     Stern J. (eds) Advances in Cryptology — EUROCRYPT ’99. EUROCRYPT 1999.
//     Lecture Notes in Computer Science, vol 1592. Springer, Berlin, Heidelberg
//     http://groups.csail.mit.edu/cis/pubs/stasio/vss.ps.gz

// LocalMember represents one member in a threshold key sharing group, prior to
// any sharing or key generation process.
type LocalMember struct {
	// ID of this group member.
	ID string
	// The BLS ID of this group member, computed from the ID.
	BlsID bls.ID
	// The threshold of group members who must be honest in order for the
	// generated key to be uncompromised. Corresponds to the number of secret
	// shares and public commitments of this group member.
	threshold int
	// Created locally, these are the `threshold` secret components that,
	// combined, represent this group member's share of the group secret key.
	// They are used to generate shares of this member's group secret key share
	// for other members, which can be verified against the public commitments
	// from this member.
	secretShares []bls.SecretKey
	// Created locally from secretShares, these are the `threshold` public
	// commitments to this group member's secret shares, which are broadcast to
	// all other members.
	shareCommitments []bls.PublicKey
}

// SharingMember represents one member in a threshold key sharing group, after
// it has a full list of `memberIDs` that belong to its threshold group. A
// member in this state has a set of `memberShares`, one for each member of the
// group, which can be accessed per member using `SecretShareForID()`. A member
// in this state also has a set of public commitments, accessible via
// `Commitments()`.
//
// As public commitments come in from other members, they can be added using
// `AddCommitmentsFromID`. Similarly, as private shares come in from other
// members, they can be added using `AddShareFromID`.
//
// Once all commitments and shares have been received, `Accusations()` will
// return a full list of members who sent invalid private shares. These can then
// be broadcast to the group, and the member can be transitioned to
// the justification phase using `InitializeJustification()`.
type SharingMember struct {
	LocalMember

	// A list of the ids of all members in the threshold group, including this
	// one.
	memberIDs []bls.ID

	// Shares of this group member's secret, one per member of the overall
	// group. The group member generates a share of its own secret as well! Note
	// that a share for a given member m is shared privately with that member in
	// the secret sharing phase. It is only shared publicly this member receives
	// an accusation from m in the accusation phase; this public sharing takes
	// place in the justification phase.
	memberShares map[bls.ID]bls.SecretKey

	// The public commitments received from each other group member. For each
	// other group member, we track their list of public commitments to their
	// private secrets. This allows us to verify the share of their private
	// secret that they send us.
	commitments map[bls.ID][]bls.PublicKey
	// For each other group member m, the share of that member's secret that m
	// sent this group member. A share is only added if it is valid; a member
	// with no entry for their received share has either not sent their share
	// or has sent an invalid share; they are therefore subject to an accusation
	// requiring them to reveal their share to all group members.
	receivedShares map[bls.ID]bls.SecretKey
}

// JustifyingMember represents a threshold group member that has entered the
// justification phase. In this phase, the member will receive a set of
// accusations broadcast to the group from other members via
// `AddAccusationFromID`. Once all accuations have been received, the member
// provides access to a set of justifications for those accusers via
// `Justifications()`, which should be broadcast to all members. Finally, as
// justifications are received they can be recorded using
// `RecordJustificationFromID`. Once all justifications have been received and
// recorded, call `FinalizeMember()` to get the final `Member`. See [GJKR 99],
// Fig. 2 (c).
type JustifyingMember struct {
	SharingMember

	// A list of ids of other group members who have accused this group member
	// of sending them an invalid share.
	accuserIDs []bls.ID
	// A map of accuser IDs to a "set" of the IDs they accused.
	pendingJustificationIDs map[bls.ID]map[bls.ID]bool
}

// Member represents a fully initialized threshold group member that is ready to
// participate in group threshold signatures and signature validation.
type Member struct {
	JustifyingMember

	// Public key for the group; nil if not yet computed.
	groupPublicKey *bls.PublicKey
	// This group member's share of the group secret key; nil if not yet
	// computed.
	groupSecretKeyShare *bls.SecretKey
	// The final list of qualified group members; empty if not yet computed.
	qualifiedMembers []bls.ID
}

// NewMember creates a new member with the given id for a threshold group with
// the given threshold. The id should be a base-10 string and is encoded into a
// bls.ID for use with the built-in secret sharing. The id should be unique per
// group member.
//
// Note that the returned member is not initialized; you will need to call
// `Initialize` on it once the full list of member IDs for the group is available,
// at which time it will be promoted to an `InitializedMember`.
func NewMember(id string, threshold int) LocalMember {
	blsID := bls.ID{}
	blsID.SetHexString(id)

	// Note: bls.SecretKey, before we call some sort of `Set` on it, can be
	// considered a zeroed *container* for a secret key.
	//
	//  - `SetByCSPRNG` initializes the zeroed secret key from a
	//    cryptographically secure pseudo-random number generator.
	//  - `Set` instead initializes a key from an existing set of shares and a
	//    group member bls.ID.
	secretShares := make([]bls.SecretKey, threshold)
	shareCommitments := make([]bls.PublicKey, threshold)

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
	// coefficients a_ik and b_ik, k in [0,t], in two polynomials A and B,
	// and store them in secretShares. We also generate the equivalent public
	// keys, C_ik = g^{a_ik}·h^{b_ik} mod p, which are stored as the commitments
	// to those shares.
	for i := 0; i < threshold; i++ {
		secretShares[i].SetByCSPRNG()

		// The public keys for each share of this group member's secret key
		// represent a public commitment to the underlying secret key shares.
		// Another member cannot get the secret key or secret key shares from
		// the public keys, but they can use them to verify that the shares of
		// the group secret key sent from this member were validly generated
		// from the same secret data.
		shareCommitments[i] = *secretShares[i].GetPublicKey()
	}

	return LocalMember{
		ID:               id,
		BlsID:            blsID,
		threshold:        threshold,
		secretShares:     secretShares,
		shareCommitments: shareCommitments,
	}
}

// InitializeSharing initializes a LocalMember with a list of the memberIDs of
// all members in the threshold group it is operating in, producing a
// SharingMember ready to participate in secret sharing.
func (member *LocalMember) InitializeSharing(otherMemberIDs []bls.ID) SharingMember {
	memberIDs := append(otherMemberIDs, member.BlsID)

	// [GJKR 99], Fig 2, 1(a).
	// For each member (including the caller!), we create a share from our set
	// of secret shares (that is, our polynomials). Equivalent to (s_ij, s'_ij),
	// but carried in the envelope of a bls.SecretKey (similar to (a_ik, b_ik)).
	shares := make(map[bls.ID]bls.SecretKey)
	for _, memberID := range memberIDs {
		memberShare := bls.SecretKey{}
		memberShare.Set(member.secretShares, &memberID)
		shares[memberID] = memberShare
	}

	return SharingMember{
		LocalMember:    *member,
		memberIDs:      memberIDs,
		memberShares:   shares,
		commitments:    make(map[bls.ID][]bls.PublicKey),
		receivedShares: make(map[bls.ID]bls.SecretKey),
	}
}

// Commitments returns the `threshold` public commitments this group member has
// generated corresponding to the `threshold` shares of its secret key.
func (member LocalMember) Commitments() []bls.PublicKey {
	return member.shareCommitments
}

func (member SharingMember) OtherMemberIDs() []bls.ID {
	otherIDs := make([]bls.ID, 0, len(member.memberIDs)-1)
	for _, memberID := range member.memberIDs {
		if !memberID.IsEqual(&member.BlsID) {
			otherIDs = append(otherIDs, memberID)
		}
	}

	return otherIDs
}

// SecretShareForID returns the secret share this member has generated for the
// given `memberID`.
func (member *SharingMember) SecretShareForID(memberID bls.ID) bls.SecretKey {
	return member.memberShares[memberID]
}

// AddCommitmentsFromID associates the given commitments with the given
// memberID. These will later be used to verify the validity of the member
// shares sent by the member with that id.
func (member *SharingMember) AddCommitmentsFromID(memberID bls.ID, commitments []bls.PublicKey) {
	member.commitments[memberID] = commitments
}

// CommitmentsComplete returns true if all commitments expected by this member
// have been seen, false otherwise.
func (member SharingMember) CommitmentsComplete() bool {
	return len(member.commitments) == len(member.memberIDs)-1
}

// AddShareFromID associates the given secret share with the given `senderID`,
// if and only if the share is valid with respect to the public commitments the
// sharing member gave.
func (member *SharingMember) AddShareFromID(senderID bls.ID, share bls.SecretKey) {
	if member.isValidShare(senderID, share) {
		member.receivedShares[senderID] = share
	}
}

// SharesComplete returns true if all shares expected by this member have been
// seen, false otherwise.
func (member SharingMember) SharesComplete() bool {
	// FIXME If a member sent an invalid share, we'll never hit the right len.
	return len(member.receivedShares) == len(member.memberIDs)-1
}

// Check whether the given share is valid with respect to the sender's public
// commitvments as seen by this member.
func (member SharingMember) isValidShare(shareSenderID bls.ID, share bls.SecretKey) bool {
	commitments := member.commitments[shareSenderID]

	combinedCommitment := bls.PublicKey{}
	combinedCommitment.Set(commitments, &member.BlsID)

	comparisonShare := share.GetPublicKey()

	return combinedCommitment.IsEqual(comparisonShare)
}

// AccusedIDs returns the list of member IDs that this member will accuse. These
// are the members who have either not sent their shares to this group member,
// or who sent their shares but the shares were invalid with respect to their
// public commitments.
func (member SharingMember) AccusedIDs() []bls.ID {
	accusedIDs := make([]bls.ID, 0, len(member.memberIDs)-len(member.receivedShares))
	for _, memberID := range member.OtherMemberIDs() {
		if _, found := member.receivedShares[memberID]; !found {
			accusedIDs = append(accusedIDs, memberID)
		}
	}

	return accusedIDs
}

// InitializeJustification switches a member from sharing mode to justifying
// mode.
func (member SharingMember) InitializeJustification() JustifyingMember {
	return JustifyingMember{
		member,
		make([]bls.ID, 0),
		make(map[bls.ID]map[bls.ID]bool),
	}
}

// AddAccusationFromID registers an accusation sent by the member with the given
// `senderID` against the member with id `accusedID`, claiming the accused sent
// an invalid share to the sender.
func (member *JustifyingMember) AddAccusationFromID(senderID bls.ID, accusedID bls.ID) {
	if accusedID.IsEqual(&member.BlsID) {
		member.accuserIDs = append(member.accuserIDs, senderID)
	} else {
		existingAccusedIDs, found := member.pendingJustificationIDs[senderID]
		if !found {
			existingAccusedIDs = make(map[bls.ID]bool)
			member.pendingJustificationIDs[senderID] = existingAccusedIDs
		}
		existingAccusedIDs[accusedID] = true
	}
}

// Justifications returns a map from accuser ID to their secret share that is
// to be broadcast to justify against an accusation. A given accuser will have
// accused this member of providing an invalid secret share with respect to this
// member's public commitments, and this justification publishes that share for
// all other members to verify against the same public commitments.
func (member JustifyingMember) Justifications() map[bls.ID]bls.SecretKey {
	justifications := make(map[bls.ID]bls.SecretKey, len(member.accuserIDs))
	for _, accuserID := range member.accuserIDs {
		justifications[accuserID] = member.memberShares[accuserID]
	}
	return justifications
}

// RecordJustificationFromID records, from this member's perspective, a
// justification from accusedID regarding an accusation from accuserID, in the
// form of the secretShare that was privately exchanged between accusedID and
// accuserID.
func (member *JustifyingMember) RecordJustificationFromID(accusedID bls.ID, accuserID bls.ID, secretShare bls.SecretKey) {
	if !member.isValidShare(accusedID, secretShare) {
		// If the member broadcast an invalid justification, we immediately
		// remove them from our shares as they have proven dishonest.
		delete(member.receivedShares, accusedID)
	} else {
		if pendingAccusedIDs, found := member.pendingJustificationIDs[accuserID]; found {
			delete(pendingAccusedIDs, accusedID)
			if len(pendingAccusedIDs) == 0 {
				delete(member.pendingJustificationIDs, accuserID)
			}
		}

		if accuserID.IsEqual(&member.BlsID) {
			// If we originally accused, and the justification is valid, then we
			// can add the valid entry to our received shares.
			member.receivedShares[accuserID] = secretShare
		}
	}
}

func (member *JustifyingMember) deleteUnjustifiedShares() {
	// At this point any entry in pendingJustificationIDs is a member who was
	// accused but whose justification we did not see. Those members are invalid
	// from our perspective. For each accuser that remains, go through the IDs
	// they accused. For each of those IDs, clear out their received shares, as
	// their failure to justify means they are not eligible players.
	for _, accusedIDs := range member.pendingJustificationIDs {
		for accusedID := range accusedIDs {
			delete(member.receivedShares, accusedID)
		}
	}
}

// FinalizeMember initializes a member that has finished the justification phase
// into a fully functioning Member that knows the group public key and can sign
// with a share of the private key.
func (member JustifyingMember) FinalizeMember() Member {
	member.deleteUnjustifiedShares()

	// [GJKR 99], Fig 2, 3
	initialShare := member.SecretShareForID(member.BlsID)
	groupSecretKeyShare := &initialShare
	for _, share := range member.receivedShares {
		groupSecretKeyShare.Add(&share)
	}

	// [GJKR 99], Fig 2, 4(c)? There is an accusation flow around public key
	//            			   computation as well...
	combinedCommitments := make([]bls.PublicKey, member.threshold)
	for i, commitment := range member.shareCommitments {
		combinedCommitments[i] = commitment
	}
	for _, commitmentSet := range member.commitments {
		for i, commitment := range commitmentSet {
			combinedCommitments[i].Add(&commitment)
		}
	}
	groupPublicKey := combinedCommitments[0]

	// Qualified players are the players who ended up with entries in
	// receivedShares; other players were removed.
	qualifiedMembers := make([]bls.ID, 0, len(member.receivedShares))
	for memberID := range member.receivedShares {
		qualifiedMembers = append(qualifiedMembers, memberID)
	}

	return Member{
		JustifyingMember:    member,
		groupSecretKeyShare: groupSecretKeyShare,
		groupPublicKey:      &groupPublicKey,
		qualifiedMembers:    qualifiedMembers,
	}
}

// GroupPublicKeyBytes returns a fixed-length 96-byte array containing the value
// of the group public key.
func (member *Member) GroupPublicKeyBytes() [96]byte {
	keyBytes := [96]byte{}
	copy(keyBytes[:], member.groupPublicKey.Serialize())

	return keyBytes
}

// SignatureShare returns this member's serialized share of the threshold
// signature for the given message. It can be combined with `threshold` other
// signatures to produce a valid group signature (that is the same no matter
// which other members participate).
func (member Member) SignatureShare(message string) []byte {
	return member.groupSecretKeyShare.Sign(message).Serialize()
}

// VerifySignature takes a message and a set of serialized signature shares by
// member ID, and verifies that the signature shares combine to a group
// signature that is valid for the given message. Returns true if so, false if
// not.
func (member Member) VerifySignature(signatureShares map[bls.ID][]byte, message string) bool {
	availableIDs := make([]bls.ID, 0, len(signatureShares))
	deserializedShares := make([]bls.Sign, 0, len(signatureShares))
	for _, memberID := range member.memberIDs {
		if serializedShare, found := signatureShares[memberID]; found {
			share := bls.Sign{}
			share.Deserialize(serializedShare)

			availableIDs = append(availableIDs, memberID)
			deserializedShares = append(deserializedShares, share)
		}
	}

	fullSignature := bls.Sign{}
	fullSignature.Recover(deserializedShares, availableIDs)

	return fullSignature.Verify(member.groupPublicKey, message)
}
