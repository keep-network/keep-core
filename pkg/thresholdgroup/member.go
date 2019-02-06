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
// and/or networking setups. This also means that none of the underlying
// implementation in this package is thread-safe.
//
// The distributed key generation approach is based on [GJKR 99], which in turn
// relies partially on the verifiable secret sharing approach in [Ped91b].
// References throughout the code are to these papers.
//
//     [GJKR 99]: Gennaro R., Jarecki S., Krawczyk H., Rabin T. (1999) Secure
//         Distributed Key Generation for Discrete-Log Based Cryptosystems. In:
//         Stern J. (eds) Advances in Cryptology — EUROCRYPT ’99. EUROCRYPT 1999.
//         Lecture Notes in Computer Science, vol 1592. Springer, Berlin, Heidelberg
//         http://groups.csail.mit.edu/cis/pubs/stasio/vss.ps.gz
//     [Ped91b]: T. Pedersen. Non-interactive and information-theoretic secure
//         verifiable secret sharing. In: Advances in Cryptology — Crypto '91,
//         pages 129-140. LNCS No. 576.
//         https://www.cs.cornell.edu/courses/cs754/2001fa/129.PDF
package thresholdgroup

import (
	"fmt"

	"github.com/dfinity/go-dfinity-crypto/bls"
)

// High-level description with 3 players. Note that in practice 3 players would
// mean our threshold could be at most 1, so we would only operate with 1
// coefficient. We do 3 players and threshold equals 2 for explanatory purposes:
//
// Take 3 players. Each generates a random polynomial of `dishonestThreshold` degree:
//
// f_1(x) = a_10 + a_11 x + a_12 x^2
// f_2(x) = a_20 + a_21 x + a_22 x^2
// f_3(x) = a_30 + a_31 x + a_32 x^2
//
// Each player broadcasts an array of commitments `g^<coefficient>` to all players:
//
// 1: [g^{a_10}, g^{a_11}, g^{a_12}]
// 2: [g^{a_20}, g^{a_22}, g^{a_22}]
// 3: [g^{a_30}, g^{a_32}, g^{a_32}]
//
// Each player i then sends a private share s_ij to each other player j,
// s_ij = f_i(j):
//
//        |             1                |             2                |             3                |
//    1   | a_10 + a_11 (1) + a_12 (1)^2 | a_10 + a_11 (2) + a_12 (2)^2 | a_10 + a_11 (3) + a_12 (3)^2 |
//    2   | a_20 + a_21 (1) + a_22 (1)^2 | a_20 + a_21 (2) + a_22 (2)^2 | a_20 + a_21 (3) + a_22 (3)^2 |
//    3   | a_30 + a_31 (1) + a_32 (1)^2 | a_30 + a_31 (2) + a_32 (2)^2 | a_30 + a_31 (3) + a_32 (3)^2 |
//
// Each player can verify their private share against the commitments by raising
// the g to the private share s_ij and checking that:
//
//   g^{s_ij} = (g^{a_i0})^j^0 · (g^{a_i1})^j^1 · (g^{a_i2})^j^2
//
// This is because:
//
//   g^{s_ij} = g^{a_i0 + a_i1 j + a_i2 j^2}
//            = g^{a_i0} · g^{a_i1 j} · g^{a_i2 j^2}
//
// Unverified shares are accused, and the players in question have the
// opportunity to justify themselves. All players with valid shares at the end
// are part of the qualified set, and are the only players considered after this
// point.
//
// With each share `g^{s_ij}`, each player can now recover g raised to the power
// of the zero coefficient of every other player, `g^{a_i0}`. These are the
// shares of the group public key, and can be combined by multiplying all of
// them into the complete group public key.
//
// Each player i's final share of the group private key is then the sum of the
// private shares they received from other players. For player 1, for example,
// this is s_11 + s_12 + s_13 = f_1(1) + f_2(1) + f_3(1). Note that the group
// private key is the sum of all of the players' zero coefficients, `a_i0`, but
// we don't have enough information at this stage for any given player to be
// able to recover this group private key, so only by collaborating can the
// players make use of it.
//
// PS: All of this is mod math. Private shares and keys are done mod q, public
// shares and keys are done mod p.

// BaseMember is a common interface implemented by all stages of threshold group
// members. It provides access to the member ID as a string irrespective of
// current DKG phase.
type BaseMember interface {
	MemberID() string
}

type BlsID bls.ID

func (id *BlsID) String() string {
	return (*bls.ID)(id).GetDecString()
}

func (id *BlsID) Raw() *bls.ID {
	return (*bls.ID)(id)
}

type memberCore struct {
	// ID of this group member. Hex number in string form.
	ID string
	// The BLS ID of this group member, computed from the ID.
	BlsID *BlsID
	// The number of members in the complete group.
	groupSize int
	// The BLS IDs of all members of this member's group, including the member
	// itself. Initially empty, populated as each other member announces its
	// presence.
	memberIDs []*bls.ID
}

// LocalMember represents one member in a threshold group, prior to the
// initiation of the distributed key generation process.
type LocalMember struct {
	memberCore
	// The maximum number of group members who could be dishonest in order for the
	// generated key to be uncompromised.
	dishonestThreshold int
	// Created locally, these are the `dishonestThreshold + 1` secret components that,
	// combined, represent this group member's share of the group secret key.
	// They are used to generate shares of this member's group secret key share
	// for other members, which can be verified against the public commitments
	// from this member.
	secretShares []bls.SecretKey
	// Created locally from secretShares, these are the `dishonestThreshold + 1` public
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
	memberShares map[bls.ID]*bls.SecretKey

	// The public commitments received from each other group member. For each
	// other group member, we track their list of public commitments to their
	// private secrets. This allows this group member to verify the share of
	// each other member's private secret that is shared secretly with this
	// member.
	commitments map[bls.ID][]bls.PublicKey
	// For each other group member m, the share of m's secret that they sent
	// this group member. A share is only added if it is valid; a member with no
	// entry for their received share has not sent their share, while a member
	// with a nil entry has sent an invalid share. As such, they are subject to
	// an accusation requiring them to publicly reveal their private share for
	// this member to all group members.
	receivedShares map[bls.ID]*bls.SecretKey
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
	// The minimum number of participants needed to produce a valid signature.
	signingThreshold int
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
// `dishonestThreshold` >= `groupSize` / 2, as the distributed key generation and threshold
// signature algorithm security breaks down past that point.
func NewMember(id string, dishonestThreshold int, groupSize int) (*LocalMember, error) {
	if float64(dishonestThreshold) >= float64(groupSize)/2 {
		return nil, fmt.Errorf(
			"threshold %v >= %v / 2, so group security cannot be guaranteed",
			dishonestThreshold,
			groupSize,
		)
	}

	blsID := new(BlsID)
	err := blsID.Raw().SetHexString(id)
	if err != nil {
		return nil, err
	}

	// According to [GJKR 99] the polynomial is of degree `dishonestThreshold`,
	// it means that we have `dishonestThreshold + 1` coefficients in the polynomial,
	// which is also the number of `secretShares` and `shareCommitments`
	//
	// Note: bls.SecretKey, before we call some sort of `Set` on it, can be
	// considered a zeroed *container* for a secret key.
	//
	//  - `SetByCSPRNG` initializes the zeroed secret key from a
	//    cryptographically secure pseudo-random number generator.
	//  - `Set` instead initializes a key from an existing set of shares and a
	//    group member bls.ID.
	secretSharesCount := dishonestThreshold + 1
	secretShares := make([]bls.SecretKey, secretSharesCount)
	shareCommitments := make([]bls.PublicKey, secretSharesCount)

	// Alternate description from original Pedersen VSS paper, [Ped91b]
	// reference in [GJKR 99]:
	// F(x) and G(x) are polynomials of degree k with F_i/G_i being random
	// coefficients in those polynomials, i ∈ [1,k-1]. 0 coefficient for F is
	// s, for G is t. s is the secret, t here is a random value chosen by the
	// current player (used to mask s). k is the threshold.
	// F_i = coefficient i in F(x) = s + F_1·x + F_2·x^2 + ... + F_{k-1}·x^{k-1}
	// G_i = coefficient i in G(x) = t + G_1·x + G_2·x^2 + ... + G_{k-1}·x^{k-1}
	// E_i = E(F_i, G_i) = g^{F_i}·h^{G_i}
	// g is a generator of the group G_q, h is another element in G_q, such that
	// no one knows log_g(h)
	// Commitmnent to s is E_0 = E(s, t) = g^s·h^t.
	// Broadcast commitment is E_i = E(F_i, G_i) for i = 1, ..., k - 1
	//
	// [GJKR 99], Fig 2, 1(a).
	// For this dealer, i, we generate t secret keys, which are equivalent to t
	// coefficients a_ik and b_ik, k in [0,t], in two polynomials A and B,
	// and store them in secretShares. We also generate the equivalent public
	// keys, C_ik = g^{a_ik}·h^{b_ik} mod p, which are stored as the commitments
	// to those shares.
	for i := 0; i < secretSharesCount; i++ {
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
			memberIDs: make([]*bls.ID, 0, groupSize),
		},
		dishonestThreshold: dishonestThreshold,
		secretShares:       secretShares,
		shareCommitments:   shareCommitments,
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

// Commitments returns the `dishonestThreshold + 1` public commitments this group member has
// generated corresponding to the `dishonestThreshold + 1` shares of its secret key.
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
	shares := make(map[bls.ID]*bls.SecretKey)
	for _, memberID := range lm.memberIDs {
		memberShare := &bls.SecretKey{}
		memberShare.Set(lm.secretShares, memberID)
		shares[*memberID] = memberShare
	}

	return &SharingMember{
		LocalMember:    *lm,
		memberShares:   shares,
		commitments:    make(map[bls.ID][]bls.PublicKey),
		receivedShares: make(map[bls.ID]*bls.SecretKey),
	}
}

// OtherMemberIDs returns the BLS IDs of all members in the group except this
// one.
func (sm *SharingMember) OtherMemberIDs() []*bls.ID {
	otherIDs := make([]*bls.ID, 0, len(sm.memberIDs)-1)
	for _, memberID := range sm.memberIDs {
		if !memberID.IsEqual(sm.BlsID.Raw()) {
			otherIDs = append(otherIDs, memberID)
		}
	}

	return otherIDs
}

// SecretShareForID returns the secret share this member has generated for the
// given `memberID`.
func (sm *SharingMember) SecretShareForID(memberID *bls.ID) *bls.SecretKey {
	return sm.memberShares[*memberID]
}

// AddCommitmentsFromID associates the given commitments with the given
// memberID. These will later be used to verify the validity of the member
// shares sent by the member with that id.
func (sm *SharingMember) AddCommitmentsFromID(memberID *bls.ID, commitments []bls.PublicKey) {
	sm.commitments[*memberID] = commitments
}

// CommitmentsComplete returns true if all commitments expected by this member
// have been seen, false otherwise.
func (sm *SharingMember) CommitmentsComplete() bool {
	return len(sm.commitments) == sm.groupSize-1
}

// AddShareFromID associates the given secret share with the given `senderID`,
// if and only if the share is valid with respect to the public commitments the
// sharing member gave.
func (sm *SharingMember) AddShareFromID(senderID *bls.ID, share *bls.SecretKey) {
	if share != nil && sm.isValidShare(senderID, share) {
		sm.receivedShares[*senderID] = share
	} else {
		sm.receivedShares[*senderID] = nil
	}
}

// SharesComplete returns true if all shares expected by this member have been
// seen, false otherwise.
func (sm *SharingMember) SharesComplete() bool {
	return len(sm.receivedShares) == len(sm.memberIDs)-1
}

// Check whether the given share from the sender to this member is valid with
// respect to the sender's public commitments as seen by this member.
func (sm *SharingMember) isValidShare(shareSenderID *bls.ID, share *bls.SecretKey) bool {
	return sm.isValidShareFor(shareSenderID, sm.BlsID.Raw(), share)
}

// Check whether the given share from the sender to the receiver is valid with
// respect to the sender's public commitments as seen by this member.
func (sm *SharingMember) isValidShareFor(
	shareSenderID *bls.ID,
	shareReceiverID *bls.ID,
	share *bls.SecretKey,
) bool {
	commitments := sm.commitments[*shareSenderID]
	// Empty commitments cause `Set` below to panic, so quit early.
	if len(commitments) <= 0 {
		return false
	}

	combinedCommitment := bls.PublicKey{}
	combinedCommitment.Set(commitments, shareReceiverID)

	comparisonShare := share.GetPublicKey()

	return combinedCommitment.IsEqual(comparisonShare)
}

// AccusedIDs returns the list of member IDs that this member will accuse. These
// are the members who have either not sent their shares to this group member,
// or who sent their shares but the shares were invalid with respect to their
// public commitments.
func (sm *SharingMember) AccusedIDs() []*bls.ID {
	accusedIDs := make([]*bls.ID, 0)
	for _, memberID := range sm.OtherMemberIDs() {
		if share, found := sm.receivedShares[*memberID]; !found || share == nil {
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
// an invalid share to the sender. The accusation may be against this member.
func (jm *JustifyingMember) AddAccusationFromID(senderID *bls.ID, accusedID *bls.ID) {
	if accusedID.IsEqual(jm.BlsID.Raw()) {
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
func (jm *JustifyingMember) Justifications() map[bls.ID]*bls.SecretKey {
	justifications := make(map[bls.ID]*bls.SecretKey, len(jm.accuserIDs))
	for _, accuserID := range jm.accuserIDs {
		justifications[accuserID] = jm.memberShares[accuserID]
	}
	return justifications
}

// RecordJustificationFromID records, from this member's perspective, a
// justification from accusedID regarding an accusation from accuserID, in the
// form of the secretShare that was privately exchanged between accusedID and
// accuserID.
func (jm *JustifyingMember) RecordJustificationFromID(accusedID *bls.ID, accuserID *bls.ID, secretShare *bls.SecretKey) {
	if !jm.isValidShareFor(accusedID, accuserID, secretShare) {
		// If the member broadcast an invalid justification, we immediately
		// remove them from our shares as they have proven dishonest.
		jm.receivedShares[*accusedID] = nil
	} else {
		if pendingAccusedIDs, found := jm.pendingJustificationIDs[*accuserID]; found {
			delete(pendingAccusedIDs, *accusedID)
			if len(pendingAccusedIDs) == 0 {
				delete(jm.pendingJustificationIDs, *accuserID)
			}
		}

		if accuserID.IsEqual(jm.BlsID.Raw()) {
			// If we originally accused, and the justification is valid, then we
			// can add the valid entry to our received shares.
			jm.receivedShares[*accuserID] = secretShare
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

	// Also clear nil shares, which are shares that were invalid and never
	// justified.
	for id, share := range jm.receivedShares {
		if share == nil {
			delete(jm.receivedShares, id)
		}
	}
}

// FinalizeMember initializes a member that has finished the justification phase
// into a fully functioning Member that knows the group public key and can sign
// with a share of the private key.
//
// Returns an error if, during finalization, the final set of qualified members
// (including this member) is less than the `dishonestThreshold+1`.
func (jm *JustifyingMember) FinalizeMember() (*Member, error) {
	jm.deleteUnjustifiedShares()

	// Note: this member is counted as a qualified member.
	if len(jm.receivedShares) < jm.dishonestThreshold {
		return nil, fmt.Errorf(
			"required %v qualified members but only had %v",
			jm.dishonestThreshold+1,
			len(jm.receivedShares)+1,
		)
	}

	// [GJKR 99], Fig 2, 3
	groupSecretKeyShare := &bls.SecretKey{}
	groupSecretKeyShare.SetLittleEndian(jm.SecretShareForID(jm.BlsID.Raw()).GetLittleEndian())
	for _, share := range jm.receivedShares {
		groupSecretKeyShare.Add(share)
	}

	// [GJKR 99], Fig 2, 4(c)? There is an accusation flow around public key
	//            			   computation as well...
	combinedCommitments := make([]bls.PublicKey, len(jm.shareCommitments))
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
		signingThreshold:    jm.dishonestThreshold + 1,
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
// signature for the given message. It can be combined with `signingThreshold` other
// signatures to produce a valid group signature. This group signature will be
// the same no matter which other group members' signatures are combined, as
// long as there are at least `signingThreshold` of them.
func (m *Member) SignatureShare(message string) []byte {
	return m.groupSecretKeyShare.Sign(message).Serialize()
}

// CompleteSignature takes a set of signature shares, bls.IDs associated with
// the bytes of each member's signature, and combines them into one complete
// signature. Returns an error if the number of signature shares is less then
// the `signingThreshold`.
func (m *Member) CompleteSignature(signatureShares map[bls.ID][]byte) (*bls.Sign, error) {
	if len(signatureShares) < m.signingThreshold {
		return nil, fmt.Errorf(
			"%v shares are insufficient for a complete signature; need %v",
			len(signatureShares),
			m.signingThreshold,
		)
	}

	availableIDs := make([]bls.ID, 0, len(signatureShares))
	deserializedShares := make([]bls.Sign, 0, len(signatureShares))
	for _, memberID := range m.memberIDs {
		if serializedShare, found := signatureShares[*memberID]; found {
			share := bls.Sign{}
			if err := share.Deserialize(serializedShare); err != nil {
				return nil, fmt.Errorf("failed to deserliaze share %v with err: %v", serializedShare, err)
			}

			availableIDs = append(availableIDs, *memberID)
			deserializedShares = append(deserializedShares, share)
		}
	}

	fullSignature := bls.Sign{}
	if err := fullSignature.Recover(deserializedShares, availableIDs); err != nil {
		return nil, fmt.Errorf("failed to recover the fullsignature with err: %v", err)
	}

	return &fullSignature, nil
}

// VerifySignature takes a message and a set of serialized signature shares by
// member ID, and verifies that the signature shares combine to a group
// signature that is valid for the given message. Returns true if so, false if
// not.
func (m *Member) VerifySignature(signatureShares map[bls.ID][]byte, message string) (bool, error) {
	fullSignature, err := m.CompleteSignature(signatureShares)
	if err != nil {
		return false, err
	}

	return fullSignature.Verify(m.groupPublicKey, message), nil
}
