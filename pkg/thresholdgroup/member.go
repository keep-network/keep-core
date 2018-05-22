// Package thresholdgroup contains the code that implements threshold key
// generation and signing using BLS for members of groups of arbitrary size. It
// consists of a series of structs that represent the various states of a member
// of a threshold group, from an uninitialized member to a member that has a
// valid private share of a group key and can perform threshold signatures.
//
// thresholdgroup.NewMember creates a new LocalMember with a given id, which
// will participate in a group of a set size with a specified threshold. Each
// following struct (SharingMember, JustifyingMember, and Member) represents a
// following phase of distributed key generation.
//
// This package does not implement any synchronization or network operations;
// instead, it is meant to be the core implementation of distributed key
// generation and threshold signing, and can be plugged into separate sync
// and/or networking setups.
//
// The distributed key generation approach is based on [GJKR 99]. References
// throughout the code are to this paper.
//
//     [GJKR 99]: Gennaro R., Jarecki S., Krawczyk H., Rabin T. (1999) Secure
//         Distributed Key Generation for Discrete-Log Based Cryptosystems. In:
//         Stern J. (eds) Advances in Cryptology — EUROCRYPT ’99. EUROCRYPT 1999.
//         Lecture Notes in Computer Science, vol 1592. Springer, Berlin, Heidelberg
//         http://groups.csail.mit.edu/cis/pubs/stasio/vss.ps.gz
package thresholdgroup

import (
	"fmt"

	"github.com/dfinity/go-dfinity-crypto/bls"
)

// BaseMember is a common interface implemented by all stages of threshold group
// members. It provides access to the member ID as a string irrespective of
// current DKG phase.
type BaseMember interface {
	MemberID() string
}

type memberCore struct {
	// ID of this group member. Hex number in string form.
	ID string
	// The BLS ID of this group member, computed from the ID.
	BlsID bls.ID
	// The number of members in the complete group.
	groupSize int
	// The threshold of group members who must be honest in order for the
	// generated key to be uncompromised. Corresponds to the number of secret
	// shares and public commitments of this group member.
	threshold int
	// The BLS IDs of all members of this member's group, including the member
	// itself. Initially empty, populated as each other member announces its
	// presence.
	memberIDs []*bls.ID
}

// LocalMember represents one member in a threshold group, prior to the
// initiation of the distributed key generation process.
type LocalMember struct {
	memberCore
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
// member in this state has a map of `memberShares`, one for each member of the
// group, which can be accessed per member using `SecretShareForID()`. A member
// in this state also has a slice of public commitments, accessible via
// `Commitments()`.
//
// As public commitments come in from other members, they can be added using
// `AddCommitmentsFromID`. Similarly, as private shares come in from other
// members, they can be added using `AddShareFromID`.
//
// Once all commitments and shares have been received, `SharesComplete()` will
// return true, and `Accusations()` will return a full list of members who sent
// invalid private shares. These can then be broadcast to the group, and the
// member can be transitioned to the justification phase using
// `InitializeJustification()`.
//
// See [GJKR 99], Fig. 2, 1(a) and 1(b).
type SharingMember struct {
	LocalMember

	// Shares of this group member's secret, one per member of the overall
	// group. The group member generates a share of its own secret as well! Note
	// that a share for a given member m is shared *privately* with that member
	// in the secret sharing phase. It is only shared publicly if this member
	// receives an accusation from m in the accusation phase; in this case, the
	// public sharing takes place during the justification phase.
	memberShares map[bls.ID]bls.SecretKey

	// The public commitments received from each other group member. For each
	// other group member, we track their list of public commitments to their
	// private secrets. This allows this group member to verify the share of
	// each other member's private secret that is shared secretly with this
	// member.
	commitments map[bls.ID][]bls.PublicKey
	// For each other group member m, the share of m's secret that they sent
	// this group member. A share is only added if it is valid; a member with no
	// entry for their received share has either not sent their share or has
	// sent an invalid share. As such, they are subject to an accusation
	// requiring them to publicly reveal their private share for this member to
	// all group members.
	receivedShares map[bls.ID]bls.SecretKey
}

// JustifyingMember represents a threshold group member that has entered the
// justification phase. In this phase, the member will receive a set of
// accusations broadcast to the group from each other member via
// `AddAccusationFromID`. Once all accuations have been received, the member
// provides access to a map of all justifications for seen accusers via
// `Justifications()`, which should be broadcast publicly to all members.
// Finally, as justifications are received they can be recorded using
// `RecordJustificationFromID`. Once all justifications have been received and
// recorded, call `FinalizeMember()` to get the final `Member`.
//
// See [GJKR 99], Fig. 2, 1(c).
type JustifyingMember struct {
	SharingMember

	// A list of ids of other group members who have accused this group member
	// of sending them an invalid share.
	accuserIDs []bls.ID
	// A map of accuser IDs to a "set" of the other member IDs they accused
	// (excluding this member).
	pendingJustificationIDs map[bls.ID]map[bls.ID]bool
}

// Member represents a fully initialized threshold group member that is ready to
// participate in group threshold signatures and signature validation.
//
// See [GJKR 99], Fig. 2 (3).
type Member struct {
	memberCore

	// Public key for the group; nil if not yet computed.
	groupPublicKey *bls.PublicKey
	// This group member's share of the group secret key; nil if not yet
	// computed.
	groupSecretKeyShare *bls.SecretKey
	// The final list of qualified group members; empty if not yet computed.
	qualifiedMembers []bls.ID
}

func (mc *memberCore) MemberID() string {
	return mc.ID
}

// NewMember creates a new member with the given id for a threshold group with
// the given threshold and group size. The id should be a base-16 string and is
// encoded into a bls.ID for use with the built-in secret sharing. The id should
// be unique per group member.
//
// Returns an error if the id fails to be read as a valid hex string, or if the
// threshold >= groupSize / 2, as the distributed key generation and threshold
// signature algorithm security breaks down past that point.
func NewMember(id string, threshold int, groupSize int) (*LocalMember, error) {
	if threshold >= groupSize/2 {
		return nil, fmt.Errorf(
			"threshold %v >= %v / 2, so group security cannot be guaranteed",
			threshold,
			groupSize,
		)
	}

	blsID := bls.ID{}
	err := blsID.SetHexString(id)
	if err != nil {
		return nil, err
	}

	// Note: bls.SecretKey, before we call some sort of `Set` on it, can be
	// considered a zeroed *container* for a secret key.
	//
	//  - `SetByCSPRNG` initializes the zeroed secret key from a
	//    cryptographically secure pseudo-random number generator.
	//  - `Set` instead initializes a key from an existing set of shares and a
	//    group member bls.ID.
	secretShares := make([]bls.SecretKey, threshold)
	shareCommitments := make([]bls.PublicKey, threshold)

	// Alternate description from original Pedersen VSS paper, [Ped91b]
	// reference in [GJKR 99]:
	// F(x) and G(x) are polynomials of degree k with F_i/G_i being random
	// coefficients in those polynomials, i ∈ [1,k-1]. 0 coefficient for F is
	// s, for G is t. s is the secret, t here is a random value chosen by the
	// current player (used to mask s). k is the threshold.
	// F_i = coefficient i in F(x) = s + F_1·x + F_2·x^2 + ... + F_{k-1}·x^{k-1}
	// G_i = coefficient i in G(x) = t + G_1·x + G_2·x^2 + ... + G_{k-1}·x^{k-1}
	// E_i = E(F_i, G_i)
	// Commitmnent to s is E_0 = E(s, t) = g^s·h^t.
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

	return &LocalMember{
		memberCore: memberCore{
			ID:        fmt.Sprintf("0x%010s", id),
			BlsID:     blsID,
			groupSize: groupSize,
			threshold: threshold,
			memberIDs: make([]*bls.ID, 0, groupSize),
		},
		secretShares:     secretShares,
		shareCommitments: shareCommitments,
	}, nil
}

// RegisterMemberID adds a member to the list of group members the local member
// knows about.
func (lm *LocalMember) RegisterMemberID(id *bls.ID) {
	lm.memberIDs = append(lm.memberIDs, id)
}

// MemberListComplete returns true if the member has a complete local member
// list and is ready to move into sharing (via InitializeSharing), false
// otherwise.
func (lm *LocalMember) MemberListComplete() bool {
	return len(lm.memberIDs) >= lm.groupSize
}

// Commitments returns the `threshold` public commitments this group member has
// generated corresponding to the `threshold` shares of its secret key.
func (lm *LocalMember) Commitments() []bls.PublicKey {
	return lm.shareCommitments
}

// InitializeSharing initializes a LocalMember with a list of the memberIDs of
// all members in the threshold group it is operating in, producing a
// SharingMember ready to participate in secret sharing.
func (lm *LocalMember) InitializeSharing() *SharingMember {
	// [GJKR 99], Fig 2, 1(a).
	// For each member (including the caller!), we create a share from our set
	// of secret shares (that is, our polynomials). Equivalent to (s_ij, s'_ij),
	// but carried in the envelope of a bls.SecretKey (similar to (a_ik, b_ik)).
	shares := make(map[bls.ID]bls.SecretKey)
	for _, memberID := range lm.memberIDs {
		memberShare := bls.SecretKey{}
		memberShare.Set(lm.secretShares, memberID)
		shares[*memberID] = memberShare
	}

	return &SharingMember{
		LocalMember:    *lm,
		memberShares:   shares,
		commitments:    make(map[bls.ID][]bls.PublicKey),
		receivedShares: make(map[bls.ID]bls.SecretKey),
	}
}

// OtherMemberIDs returns the BLS IDs of all members in the group except this
// one.
func (sm *SharingMember) OtherMemberIDs() []*bls.ID {
	otherIDs := make([]*bls.ID, 0, len(sm.memberIDs)-1)
	for _, memberID := range sm.memberIDs {
		if !memberID.IsEqual(&sm.BlsID) {
			otherIDs = append(otherIDs, memberID)
		}
	}

	return otherIDs
}

// SecretShareForID returns the secret share this member has generated for the
// given `memberID`.
func (sm *SharingMember) SecretShareForID(memberID *bls.ID) bls.SecretKey {
	return sm.memberShares[*memberID]
}

// AddCommitmentsFromID associates the given commitments with the given
// memberID. These will later be used to verify the validity of the member
// shares sent by the member with that id.
func (sm *SharingMember) AddCommitmentsFromID(memberID bls.ID, commitments []bls.PublicKey) {
	sm.commitments[memberID] = commitments
}

// CommitmentsComplete returns true if all commitments expected by this member
// have been seen, false otherwise.
func (sm SharingMember) CommitmentsComplete() bool {
	return len(sm.commitments) == sm.groupSize-1
}

// AddShareFromID associates the given secret share with the given `senderID`,
// if and only if the share is valid with respect to the public commitments the
// sharing member gave.
func (sm *SharingMember) AddShareFromID(senderID bls.ID, share bls.SecretKey) {
	if sm.isValidShare(senderID, share) {
		sm.receivedShares[senderID] = share
	}
}

// SharesComplete returns true if all shares expected by this member have been
// seen, false otherwise.
func (sm *SharingMember) SharesComplete() bool {
	// FIXME If a member sent an invalid share, we'll never hit the right len.
	return len(sm.receivedShares) == len(sm.memberIDs)-1
}

// Check whether the given share is valid with respect to the sender's public
// commitvments as seen by this member.
func (sm *SharingMember) isValidShare(shareSenderID bls.ID, share bls.SecretKey) bool {
	commitments := sm.commitments[shareSenderID]

	combinedCommitment := bls.PublicKey{}
	combinedCommitment.Set(commitments, &sm.BlsID)

	comparisonShare := share.GetPublicKey()

	return combinedCommitment.IsEqual(comparisonShare)
}

// AccusedIDs returns the list of member IDs that this member will accuse. These
// are the members who have either not sent their shares to this group member,
// or who sent their shares but the shares were invalid with respect to their
// public commitments.
func (sm *SharingMember) AccusedIDs() []*bls.ID {
	accusedIDs := make([]*bls.ID, 0, len(sm.memberIDs)-len(sm.receivedShares))
	for _, memberID := range sm.OtherMemberIDs() {
		if _, found := sm.receivedShares[*memberID]; !found {
			accusedIDs = append(accusedIDs, memberID)
		}
	}

	return accusedIDs
}

// InitializeJustification switches a member from sharing mode to justifying
// mode.
func (sm *SharingMember) InitializeJustification() *JustifyingMember {
	return &JustifyingMember{
		*sm,
		make([]bls.ID, 0),
		make(map[bls.ID]map[bls.ID]bool),
	}
}

// AddAccusationFromID registers an accusation sent by the member with the given
// `senderID` against the member with id `accusedID`, claiming the accused sent
// an invalid share to the sender.
func (jm *JustifyingMember) AddAccusationFromID(senderID *bls.ID, accusedID *bls.ID) {
	if accusedID.IsEqual(&jm.BlsID) {
		jm.accuserIDs = append(jm.accuserIDs, *senderID)
	} else {
		existingAccusedIDs, found := jm.pendingJustificationIDs[*senderID]
		if !found {
			existingAccusedIDs = make(map[bls.ID]bool)
			jm.pendingJustificationIDs[*senderID] = existingAccusedIDs
		}
		existingAccusedIDs[*accusedID] = true
	}
}

// Justifications returns a map from accuser ID to their secret share that is
// to be broadcast to justify against an accusation. A given accuser will have
// accused this member of providing an invalid secret share with respect to this
// member's public commitments, and this justification publishes that share for
// all other members to verify against the same public commitments.
func (jm *JustifyingMember) Justifications() map[bls.ID]bls.SecretKey {
	justifications := make(map[bls.ID]bls.SecretKey, len(jm.accuserIDs))
	for _, accuserID := range jm.accuserIDs {
		justifications[accuserID] = jm.memberShares[accuserID]
	}
	return justifications
}

// RecordJustificationFromID records, from this member's perspective, a
// justification from accusedID regarding an accusation from accuserID, in the
// form of the secretShare that was privately exchanged between accusedID and
// accuserID.
func (jm *JustifyingMember) RecordJustificationFromID(accusedID bls.ID, accuserID bls.ID, secretShare bls.SecretKey) {
	if !jm.isValidShare(accusedID, secretShare) {
		// If the member broadcast an invalid justification, we immediately
		// remove them from our shares as they have proven dishonest.
		delete(jm.receivedShares, accusedID)
	} else {
		if pendingAccusedIDs, found := jm.pendingJustificationIDs[accuserID]; found {
			delete(pendingAccusedIDs, accusedID)
			if len(pendingAccusedIDs) == 0 {
				delete(jm.pendingJustificationIDs, accuserID)
			}
		}

		if accuserID.IsEqual(&jm.BlsID) {
			// If we originally accused, and the justification is valid, then we
			// can add the valid entry to our received shares.
			jm.receivedShares[accuserID] = secretShare
		}
	}
}

func (jm *JustifyingMember) deleteUnjustifiedShares() {
	// At this point any entry in pendingJustificationIDs is a member who was
	// accused but whose justification we did not see. Those members are invalid
	// from our perspective. For each accuser that remains, go through the IDs
	// they accused. For each of those IDs, clear out their received shares, as
	// their failure to justify means they are not eligible players.
	for _, accusedIDs := range jm.pendingJustificationIDs {
		for accusedID := range accusedIDs {
			delete(jm.receivedShares, accusedID)
		}
	}
}

// FinalizeMember initializes a member that has finished the justification phase
// into a fully functioning Member that knows the group public key and can sign
// with a share of the private key.
//
// Returns an error if, during finalization, the final set of qualified members
// is less than the `threshold`.
func (jm *JustifyingMember) FinalizeMember() (*Member, error) {
	jm.deleteUnjustifiedShares()

	if len(jm.receivedShares) < jm.threshold {
		return nil, fmt.Errorf(
			"required %v qualified members but only had %v",
			jm.threshold,
			len(jm.receivedShares),
		)
	}

	// [GJKR 99], Fig 2, 3
	initialShare := jm.SecretShareForID(&jm.BlsID)
	groupSecretKeyShare := &initialShare
	for _, share := range jm.receivedShares {
		groupSecretKeyShare.Add(&share)
	}

	// [GJKR 99], Fig 2, 4(c)? There is an accusation flow around public key
	//            			   computation as well...
	combinedCommitments := make([]bls.PublicKey, jm.threshold)
	for i, commitment := range jm.shareCommitments {
		combinedCommitments[i] = commitment
	}
	for _, commitmentSet := range jm.commitments {
		for i, commitment := range commitmentSet {
			combinedCommitments[i].Add(&commitment)
		}
	}
	groupPublicKey := combinedCommitments[0]

	// Qualified players are the players who ended up with entries in
	// receivedShares; other players were removed.
	qualifiedMembers := make([]bls.ID, 0, len(jm.receivedShares))
	for memberID := range jm.receivedShares {
		qualifiedMembers = append(qualifiedMembers, memberID)
	}

	return &Member{
		memberCore:          jm.memberCore,
		groupSecretKeyShare: groupSecretKeyShare,
		groupPublicKey:      &groupPublicKey,
		qualifiedMembers:    qualifiedMembers,
	}, nil
}

// GroupPublicKeyBytes returns a fixed-length 96-byte array containing the value
// of the group public key.
func (m *Member) GroupPublicKeyBytes() [96]byte {
	keyBytes := [96]byte{}
	copy(keyBytes[:], m.groupPublicKey.Serialize())

	return keyBytes
}

// SignatureShare returns this member's serialized share of the threshold
// signature for the given message. It can be combined with `threshold` other
// signatures to produce a valid group signature. This group signature will be
// the same no matter which other group members' signatures are combined, as
// long as there are at least `threshold` of them.
func (m *Member) SignatureShare(message string) []byte {
	return m.groupSecretKeyShare.Sign(message).Serialize()
}

// CompleteSignature takes a set of signature shares, bls.IDs associated with
// the bytes of each member's signature, and combines them into one complete
// signature. Returns an error if the number of signature shares is less than
// the group threshold.
func (m *Member) CompleteSignature(signatureShares map[bls.ID][]byte) (*bls.Sign, error) {
	if len(signatureShares) < m.threshold {
		return nil, fmt.Errorf(
			"%v shares are insufficient for a complete signature; need %v",
			len(signatureShares),
			m.threshold,
		)
	}

	availableIDs := make([]bls.ID, 0, len(signatureShares))
	deserializedShares := make([]bls.Sign, 0, len(signatureShares))
	for _, memberID := range m.memberIDs {
		if serializedShare, found := signatureShares[*memberID]; found {
			share := bls.Sign{}
			share.Deserialize(serializedShare)

			availableIDs = append(availableIDs, *memberID)
			deserializedShares = append(deserializedShares, share)
		}
	}

	fullSignature := bls.Sign{}
	fullSignature.Recover(deserializedShares, availableIDs)

	return &fullSignature, nil
}

// VerifySignature takes a message and a set of serialized signature shares by
// member ID, and verifies that the signature shares combine to a group
// signature that is valid for the given message. Returns true if so, false if
// not.
func (m *Member) VerifySignature(signatureShares map[bls.ID][]byte, message string) bool {
	fullSignature, err := m.CompleteSignature(signatureShares)
	if err != nil {
		return false
	}

	return fullSignature.Verify(m.groupPublicKey, message)
}
