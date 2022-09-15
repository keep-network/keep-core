// Package common holds some common tools that can be used across multiple
// tECDSA protocols, e.g. DKG and signing.
package common

import (
	"fmt"
	"math/big"

	"github.com/bnb-chain/tss-lib/tss"
	"github.com/keep-network/keep-core/pkg/crypto/ephemeral"
	"github.com/keep-network/keep-core/pkg/protocol/group"
)

// IdentityConverter takes care of conversions between protocol-agnostic
// member indexes and TSS-specific party IDs.
//
// A different strategy for obtaining party IDs have to be used for key
// generation and signing.
//
// Distributed Key Generation:
//  1. We receive an array of 100 operators selected to the signing group by the
//     sortition pool.
//  2. Each operator has its member index assigned, based on the position in the
//     array returned from the sortition pool. selected[0] has member index = 1,
//     selected[1] has member index = 2,…, selected[99] has member index = 100.
//  3. When the DKG attempt fails, some members are excluded from the next
//     attempt. The exclusion algorithm provides indexes of members excluded from
//     the next run.
//  4. During the execution of a retry of a DKG protocol, indexes are not getting
//     shifted. For example, if the original selection was [0xA, 0xB, 0xC, 0xD, 0xE]
//     and operator 0xC has been excluded from the attempt, the member indexes
//     stay as [1 (0xA), 2 (0xB), 4 (0xD), 5 (0xE)]. This allows filtering out
//     messages from excluded members.
//  5. Once DKG completes successfully, indexes of members excluded from the
//     execution must be set in EcdsaDkg.Result.misbehavedMembersIndices and
//     EcdsaDkg.Result.membersHash should have the misbehaving members filtered
//     out. When the DKG result gets approved, EcdsaDkg.dkgResult.membersHash is
//     stored on-chain.
//  6. When establishing the transaction submission order when publishing the
//     result to the chain, we take into account the original selection result
//     and we do not skip misbehaving members.
//  7. Before saving information to disk, misbehaving members are excluded from
//     the result and the member indexes are shifted. We do it to avoid
//     unnecessary delays when determining the order for publishing a signature
//     and all other responsibilities. Also, this way member indexes stored on
//     disk are the same as the ones expected on-chain by the members array hash
//     stored on-chain.
//
// Signing:
//  1. We load signing group members with the key shares from disk storage.
//  2. In case there were any misbehaving members during the key generation,
//     member indexes stored on disk are shifted (see point 7 of key generation).
//     If the successful DKG run was with member indexes
//     [1 (0xA), 2 (0xB), 4 (0xD), 5 (0xE)] as in the example in key generation
//     step 4, the member indexes stored on disk will be [1,2,3,4] but the TSS
//     party IDs stored on disk will be [PartyID(1), PartyID(2), PartyID(4),
//     PartyID(5)].
//  3. For the signing protocol, the same Party IDs as used in the key generation
//     must be used. It means that in case there were misbehaving members during
//     DKG, member index != TSS party ID.
type IdentityConverter interface {
	// MemberIndexToTssPartyID converts a single group member index to a
	// detached party ID instance. Such a party ID has an unset index since
	// it does not yet belong to a sorted parties IDs set.
	//
	// This function is used from GenerateTssPartiesIDs that is the very first
	// step of key-generation and signing protocol. Using party ID without the
	// internal index set is fine only at that step. For further steps, please
	// use ResolveSortedTssPartyID that uses sorted parties IDs type from
	// tss-lib to obtain the party ID with an internal index value set.
	// Tss-lib's UpdateFromBytes function will _not_ accept a party ID without
	// an internal index set.
	//
	// Please note that the member index value may not be equal to TSS party
	// ID's key value.
	MemberIndexToTssPartyID(memberIndex group.MemberIndex) *tss.PartyID

	// MemberIndexToTssPartyIDKey converts a single group member index to a
	// TSS party ID key.
	//
	// The TSS party ID is created based on the key and TSS party ID can be
	// looked up in the sorted parties IDs collection based on the key.
	// This function is used in ResolveSortedTssPartyID to resolve party ID
	// from the sorted parties' IDs. Such a party ID has an internal index value
	// set.
	//
	// Please note that the member index value may not be equal to the key value.
	MemberIndexToTssPartyIDKey(memberIndex group.MemberIndex) *big.Int

	// TssPartyIDToMemberIndex converts a single TSS party ID to a group
	// member index.
	//
	// If the provided Party ID does not map to any known member index,
	// MemberIndex(0) is returned.
	TssPartyIDToMemberIndex(partyID *tss.PartyID) group.MemberIndex
}

// GenerateTssPartiesIDs converts group member indexes to parties ID suitable
// for the TSS protocol execution. Parties IDs returned from this function do
// not have the internal index value set. Using them is only allowed for the
// first steps of key-generation and signing protocols. For further steps
// requiring a call to UpdateFromBytes, please use ResolveSortedTssPartyID.
func GenerateTssPartiesIDs(
	memberIndex group.MemberIndex,
	groupMembersIndexes []group.MemberIndex,
	converter IdentityConverter,
) (*tss.PartyID, []*tss.PartyID) {
	var partyID *tss.PartyID
	groupPartiesIDs := make([]*tss.PartyID, len(groupMembersIndexes))

	for i, groupMemberIndex := range groupMembersIndexes {
		newPartyID := converter.MemberIndexToTssPartyID(groupMemberIndex)

		if memberIndex == groupMemberIndex {
			partyID = newPartyID
		}

		groupPartiesIDs[i] = newPartyID
	}

	return partyID, groupPartiesIDs
}

// ResolveSortedTssPartyID resolves the TSS party ID for the given member index
// based on the sorted parties IDs stored in the given TSS parameters set. Such
// a resolved party ID has an index which indicates its position in the parties
// IDs set and can be used for UpdateFromBytes call.
func ResolveSortedTssPartyID(
	tssParameters *tss.Parameters,
	memberIndex group.MemberIndex,
	converter IdentityConverter,
) *tss.PartyID {
	sortedPartiesIDs := tssParameters.Parties().IDs()
	partyIDKey := converter.MemberIndexToTssPartyIDKey(memberIndex)
	return sortedPartiesIDs.FindByKey(partyIDKey)
}

// AggregateTssMessages takes a list of TSS messages and build an aggregate
// consisting of unencrypted broadcast part and encrypted point-to-point parts
// intended for specific receivers. The encryption of a specific point-to-point
// part is done using a symmetric key taken from the provided symmetricKeys map
// using the receiver member index as key.
//
// This function has also the following requirements regarding the input
// tssMessages list and symmetricKeys map:
//  1. tssMessages MUST hold 0 or 1 broadcast message
//  2. tssMessages MUST hold 0 or N point-to-point messages. If point-to-point
//     messages count is >0 then:
//     a) each point-to-point message MUST target exactly 1 unique receiver
//     b) each point-to-point message receiver MUST have a corresponding entry
//     in the symmetricKeys map
func AggregateTssMessages(
	tssMessages []tss.Message,
	symmetricKeys map[group.MemberIndex]ephemeral.SymmetricKey,
	converter IdentityConverter,
) (
	[]byte,
	map[group.MemberIndex][]byte,
	error,
) {
	var broadcastPayload []byte
	peersPayload := make(map[group.MemberIndex][]byte)

	for _, tssMessage := range tssMessages {
		tssMessageBytes, tssMessageRouting, err := tssMessage.WireBytes()
		if err != nil {
			return nil, nil, fmt.Errorf(
				"failed to unpack TSS message: [%v]",
				err,
			)
		}

		if tssMessageRouting.IsBroadcast {
			// We expect only one broadcast message. Any other case is an error.
			if len(broadcastPayload) > 0 {
				return nil, nil, fmt.Errorf(
					"multiple TSS broadcast messages detected",
				)
			}

			broadcastPayload = tssMessageBytes
		} else {
			// We expect that each P2P message targets only a single member.
			// Any other case is an error.
			if len(tssMessageRouting.To) != 1 {
				return nil, nil, fmt.Errorf(
					"multi-receiver TSS P2P message detected",
				)
			}
			// Get the single receiver ID.
			receiverID := converter.TssPartyIDToMemberIndex(tssMessageRouting.To[0])
			// Get the symmetric key with the receiver. If the symmetric key
			// cannot be found, something awful happened.
			symmetricKey, ok := symmetricKeys[receiverID]
			if !ok {
				return nil, nil, fmt.Errorf(
					"cannot get symmetric key with member [%v]",
					receiverID,
				)
			}
			// Encrypt the payload using the receiver symmetric key.
			encryptedTssMessageBytes, err := symmetricKey.Encrypt(
				tssMessageBytes,
			)
			if err != nil {
				return nil, nil, fmt.Errorf(
					"cannot encrypt TSS P2P message for member [%v]: [%v]",
					receiverID,
					err,
				)
			}

			if _, exists := peersPayload[receiverID]; exists {
				return nil, nil, fmt.Errorf(
					"duplicate TSS P2P message for member [%v]",
					receiverID,
				)
			}

			peersPayload[receiverID] = encryptedTssMessageBytes
		}
	}

	return broadcastPayload, peersPayload, nil
}
