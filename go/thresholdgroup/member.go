package thresholdgroup

import "github.com/keep-network/go-dfinity-crypto/bls"

// Member represents one member in a threshold key sharing group.
// Publicly, it exposes only an ID.
type Member struct {
	// ID of this group member.
	ID string
	// The BLS ID of this group member, computed from the ID.
	blsID bls.ID
	// The threshold of group members who must be honest in order for the
	// generated key to be uncompromised. Corresponds to the number of secret
	// shares and public commitments of this group member.
	threshold int
	// The public commitments received from each other group member. For each
	// other group member, we track their list of public commitments to their
	// private secrets. This allows us to verify the share of their private
	// secret that they send us.
	commitments map[bls.ID][]bls.PublicKey
	// Created locally, these are the t secret components that, combined,
	// represent this group member's share of the group secret key. They are
	// publicly committed to via their public keys, which are broadcast to all
	// other members. They are used to generate shares of this member's group
	// secret key share for other members, which can be verified against the
	// public commitments from this member.
	secretShares []bls.SecretKey
	// Shares of this group member's secret, one per group member. The group
	// member generates a share of its own secret as well! Note that a share for
	// a given member m is shared privately with that member in the secret
	// sharing phase and only shared publicly in the justification phase if this
	// member receives an accusation from m in the accusation phase.
	shares map[bls.ID]bls.SecretKey
	// For each other group member, the share of that member's secret that the
	// member sent this member.
	receivedShares map[bls.ID]bls.SecretKey
	// Public key for the group; nil if not yet computed.
	groupPublicKey *bls.PublicKey
	// This group member's share of the group secret key; nil if not yet
	// computed.
	groupSecretKeyShare *bls.SecretKey
	// A list of ids of other group members who have accused this group member
	// of sending them an invalid share.
	accuserIDs []bls.ID
	// Received via broadcast, tracks IDs of all members who accused this
	// member so their shares can be broadcast publicly during justification.
	accusedIDs map[bls.ID]bool
	// A (spares) map from the ID of other group members to a boolean indicating
	// if they were disqualified. Players are disqualified when they are accused
	// and then broadcast a justification that fails to verify against their
	// public commitments.
	disqualifiedPlayers map[bls.ID]bool // all players disqualified during justification
	// The final list of qualified group members; empty if not yet computed.
	qualifiedPlayers []bls.ID
}

// NewMember creates a new member with the given id. The id should be a base-10
// string and is encoded into a bls.ID for use with the built-in secret sharing.
// The id should be unique per group member.
func NewMember(id string) Member {
	blsID := bls.ID{}
	blsID.SetDecString(id)

	return Member{
		ID:                  id,
		blsID:               blsID,
		commitments:         map[bls.ID][]bls.PublicKey{},
		shares:              map[bls.ID]bls.SecretKey{},
		receivedShares:      map[bls.ID]bls.SecretKey{},
		accusedIDs:          map[bls.ID]bool{},
		disqualifiedPlayers: map[bls.ID]bool{},
	}
}
