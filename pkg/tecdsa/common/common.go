// Package common holds some common tools that can be used across multiple
// tECDSA protocols, e.g. DKG and signing.
package common

import (
	"fmt"
	"github.com/bnb-chain/tss-lib/tss"
	"github.com/keep-network/keep-core/pkg/crypto/ephemeral"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"math/big"
)

// IdentityConverter takes care of conversions between protocol-agnostic
// member indexes and TSS-specific party IDs.
type IdentityConverter interface {
	// MemberIndexToTssPartyID converts a single group member index to a
	// detached party ID instance. Such a party ID has an unset index since
	// it does not yet belong to a sorted parties IDs set.
	MemberIndexToTssPartyID(memberIndex group.MemberIndex) *tss.PartyID

	// MemberIndexToTssPartyIDKey converts a single group member index to a
	// TSS party ID key.
	MemberIndexToTssPartyIDKey(memberIndex group.MemberIndex) *big.Int

	// TssPartyIDToMemberIndex converts a single TSS party ID to a group
	// member index.
	TssPartyIDToMemberIndex(partyID *tss.PartyID) group.MemberIndex
}

// GenerateTssPartiesIDs converts group member indexes to parties ID suitable
// for the TSS protocol execution.
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
// IDs set.
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
// - tssMessages MUST hold 0 or 1 broadcast message
// - tssMessages MUST hold 0 or N point-to-point messages. If point-to-point
//   messages count is >0 then:
//     - each point-to-point message MUST target exactly 1 unique receiver
//     - each point-to-point message receiver MUST have a corresponding entry
//       in the symmetricKeys map
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
